package services

import (
	"crypto/md5"
	"sim_data_gen/daos/diggerProductBinding"
	"sim_data_gen/daos/mineProduct"
	. "sim_data_gen/mine_assignment"
	. "sim_data_gen/models/boomGroupInfo"
	"sim_data_gen/models/mineChangeStatus"
	mineProductModel "sim_data_gen/models/mineProduct"
	. "sim_data_gen/utils"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

/*
*
根据中子仪的读数,估计当前配单单元的含量
*/
type ProductionsContentService struct {
	products               map[string]*ProductionContent
	lock                   sync.Mutex
	productsMd5            string
	BreakSeconds           int64
	Crontab                string
	Cron                   *cron.Cron
	StatusSampleSeconds    int
	UseGrossOrNetWeigh     int    // 1 - 使用净重, 0 - 使用毛重
	UseNowAsEndTime        int    // 1 - 使用当前时间作为查询结束时间,使用模拟数据时采用,这样可以提前把之前的数据写入数据库,然后更改一下日期,就可以作为模拟数据 0 - 不使用结束时间,在实际环境中使用
	measureOfFindBoomGroup int    // 通过卡车号查找爆堆的对应关系的方法,  0 - 通过载货日志 , 1 - 通过历史挖机切换记录
	MineAssignmentApiUrl   string //配矿api
	//!!!:这个选项表示整个系统同时只会生成一个产品,则所有可用的资源都为这个产品准备
	//产品列表中只运行一个产品处于启用状态,必须禁用其他产品
	//系统同时有多个产品,如果该标志被设置为1,则表示Enable表示Active的意思
	//若为0则表示Enable表示启用禁用的含义
	MineProductEnableAsActive         int64
	PrimaryMaterialName               string            //主要物质索引,一般是氧化钙
	DesignatedDayAsTodayWhenFetchData string            //当获取地磅和中子仪数据的时候用某一天的日期作为日期,时间仍然用今天的,如果为空则表示使用今天作为日期
	DesignatedDayAsTodayWhenSavePlan  string            //当存储配矿方案的时候用某一天的日期作为日期,时间仍然用今天的,如果为空则表示使用今天作为日期
	Vehicle2BoomGroupCalcMode         CalcBoomGroupMode // 0 - 根据调度日志确定, 1 - 根据卡车和挖机的切换记录确定 ,2 - 根据挖机和配矿单元的切换记录确定(默认)
}

var ProductionsContentSvr *ProductionsContentService

const (
	OprtInsert             = "INSERT"
	OprtUpdate             = "UPDATE"
	OprtDelete             = "DELETE"
	PRODUCT_STATUS_ENABLED = "0"
)
const DEFAULT_MINE_ASSIGNMENT_API = "http://localhost:8080/assign/v4"

func StartProductionsContentService(varNamePrx string) {
	ProductionsContentSvr = new(ProductionsContentService)
	ProductionsContentSvr.productsMd5 = ""
	ProductionsContentSvr.MineProductEnableAsActive = GetConfig().Int64(varNamePrx+".mineProductEnableAsActive", 1)
	ProductionsContentSvr.UseGrossOrNetWeigh = GetConfig().Int(varNamePrx+".useGrossOrNetWeigh", 1)
	ProductionsContentSvr.UseNowAsEndTime = GetConfig().Int(varNamePrx+".useNowAsEndTime", 0)
	ProductionsContentSvr.DesignatedDayAsTodayWhenFetchData = GetConfig().String(varNamePrx+".designatedDayAsTodayWhenFetchData", "")
	ProductionsContentSvr.DesignatedDayAsTodayWhenSavePlan = GetConfig().String(varNamePrx+".designatedDayAsTodayWhenSavePlan", "")
	ProductionsContentSvr.PrimaryMaterialName = GetConfig().String(varNamePrx+".primaryMaterialName", "CaO")
	ProductionsContentSvr.measureOfFindBoomGroup = GetConfig().Int(varNamePrx+".measureOfFindBoomGroup", 0)
	ProductionsContentSvr.MineAssignmentApiUrl = GetConfig().String(varNamePrx+".mineAssignmentApiUrl", DEFAULT_MINE_ASSIGNMENT_API)
	ProductionsContentSvr.Vehicle2BoomGroupCalcMode = CalcBoomGroupMode(GetConfig().Int(varNamePrx+".Vehicle2BoomGroupCalcMode", int(ByDiggerBoomGroupLogger)))

	crontabEnabled := GetConfig().Int(varNamePrx+".crontabEnabled", 1)
	if crontabEnabled == 1 {
		if ProductionsContentSvr.Crontab = GetConfig().String(varNamePrx+".crontab", "0 2 * * ?"); len(ProductionsContentSvr.Crontab) == 0 {
			Logger.Info("Skip ProductionsContentSvr!")
			return
		}
		if ProductionsContentSvr.StatusSampleSeconds = GetConfig().Int(varNamePrx+".statusSampleMinutes", 10); ProductionsContentSvr.StatusSampleSeconds == 0 {
			Logger.Info("Skip StartCleanDataService!")
			return
		}
		ProductionsContentSvr.Cron = cron.New()
		if entityId, err := ProductionsContentSvr.Cron.AddFunc(ProductionsContentSvr.Crontab, ProductionsContentSvr.CalculateContentsOfProductionsFunc); err != nil {
			Logger.Error("Skip StartCleanDataService for error!", zap.Error(err))
		} else {
			Logger.Info(fmt.Sprintf("StartCleanDataService with Id = %d", entityId))
		}
		//refresh content at once.
		ProductionsContentSvr.CalculateContentsOfProductionsFunc()
		ProductionsContentSvr.Cron.Start()
	}
}
func (this *ProductionsContentService) CalculateContentsOfProductionsFunc() {
	now := time.Now()
	dateStrOfSavePlan := now.Format("2006-01-02")
	if this.DesignatedDayAsTodayWhenSavePlan != "" {
		dateStrOfSavePlan = this.DesignatedDayAsTodayWhenSavePlan
	}
	if this.DesignatedDayAsTodayWhenFetchData != "" {
		//如果这个值被设置,则表示要用历史数据进行测试,那么中子仪,车辆的操作记录,地磅都将使用历史日期,但是仍然使用今天的时间,这样的话看起来应该还是比较逼真的
		now = Change2TheDay(now, this.DesignatedDayAsTodayWhenFetchData)
	}
	dateStrOfFetchData := now.Format("2006-01-02")
	this.CalculateContentsOfProductions(dateStrOfFetchData, dateStrOfSavePlan, mineChangeStatus.LOW_PROBABILITY, false, true, now, ByLoadGoodsLogger, "", nil, "", nil, true)
}
func (this *ProductionsContentService) CalculateContentsOfProductions(dateStrOfFetchData, dateStrOfSavePlan string, minProbability int, forcePredict bool, isAsync bool, predictTime time.Time, calcBoomGroupMode CalcBoomGroupMode, snapName string, maxTimeAsUtc *int, debugProductName string, outputExcelPath *string, cacheInDb bool) {
	Logger.Info("CalculateContentsOfProductions")
	products, md5Sum, err, same := this.loadProductContents(this.productsMd5, this.StatusSampleSeconds, this.UseGrossOrNetWeigh, this.PrimaryMaterialName)
	if err != nil {
		return
	}
	//if md5 is changed,I will refresh product list.
	if !same {
		Logger.Info("product list change!")
		//如果产品的列表项发生了变化,则要更新产品列表
		this.loadProductList(products, md5Sum)
	}
	//函数调用是要传入 this.products,因为当same为true时,products是一个空map,所以需要传入this.products
	this.predictAndSaveMineAssignments4ProductList(dateStrOfFetchData, dateStrOfSavePlan, minProbability, forcePredict, isAsync, predictTime, calcBoomGroupMode, snapName, maxTimeAsUtc, debugProductName, outputExcelPath, cacheInDb)
}

func (this *ProductionsContentService) MeasureOfFindBoomGroup() int {
	return this.measureOfFindBoomGroup
}
func (this *ProductionsContentService) IsUseNowAsEndTime() bool {
	return ProductionsContentSvr.UseNowAsEndTime == 1
}

func (this *ProductionsContentService) loadProductContents(oldMd5 string, statusSampleSeconds int, weighUseInOrOutSite int, primaryMaterialName string) (map[string]*ProductionContent, string, error, bool) {
	var err error
	products := make(map[string]*ProductionContent)
	mineProductList, err := mineProduct.ListObjSortByName()
	if err != nil {
		Logger.Error("Failed to list mine products", zap.Error(err))
		return nil, "", err, false
	}
	md5Sum, err := calculateProductsMd5(mineProductList)
	if md5Sum == "" {
		Logger.Info("md5", zap.String("md5", md5Sum))
		return nil, "", err, false
	}
	if oldMd5 == md5Sum {
		Logger.Info("product list is not changed!")
		return products, md5Sum, nil, true
	}
	//读取count个配置项对应的参数
	defaultProductExisted := false
	for i := 0; i < len(mineProductList); i++ {
		if mineProductList[i].Status == nil || *mineProductList[i].Status != PRODUCT_STATUS_ENABLED {
			Logger.Info("product is not enabled,skip it!", zap.String("product", *mineProductList[i].Name))
			continue
		}
		if mineProductList[i].Name == nil || len(*mineProductList[i].Name) == 0 {
			Logger.Info("product has not name", zap.String("token", *mineProductList[i].Token))
			continue
		}
		if mineProductList[i].MatIndexes == nil || len(*mineProductList[i].MatIndexes) == 0 {
			Logger.Info("product has not matIndexes", zap.String("name", *mineProductList[i].Name))
			continue
		}
		if mineProductList[i].ContentLimits == nil || len(*mineProductList[i].ContentLimits) == 0 {
			Logger.Info("product has not ContentLimits", zap.String("name", *mineProductList[i].Name))
			continue
		}

		var productContent *ProductionContent
		if this.MineProductEnableAsActive == 1 {
			if defaultProductExisted {
				//如果已经有了默认产品,则不需要再搜索,直接退出
				break
			}
			productContent = NewProductionContent(*mineProductList[i].Token, *mineProductList[i].Name, *mineProductList[i].MatIndexes, *mineProductList[i].ContentLimits, statusSampleSeconds, weighUseInOrOutSite, primaryMaterialName, this.Vehicle2BoomGroupCalcMode, this.MineProductEnableAsActive)
			if defaultProductExisted == false && mineProductList[i].Status != nil || *mineProductList[i].Status == PRODUCT_STATUS_ENABLED {
				defaultProductExisted = true
				//这里代表的是这个产品是默认产品,所有产品中只能有一个默认产品
				productContent.IsDefault = true
			}
		} else {
			productContent = NewProductionContent(*mineProductList[i].Token, *mineProductList[i].Name, *mineProductList[i].MatIndexes, *mineProductList[i].ContentLimits, statusSampleSeconds, weighUseInOrOutSite, primaryMaterialName, this.Vehicle2BoomGroupCalcMode, this.MineProductEnableAsActive)
		}

		if productContent != nil {
			products[*mineProductList[i].Name] = productContent
		}
	}

	return products, md5Sum, nil, false
}

func calculateProductsMd5(list []mineProductModel.MineProduct) (string, error) {
	bytes, err := json.Marshal(list)
	if err != nil {
		Logger.Error("failed to calculateProductsMd5", zap.Error(err))
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(bytes)), nil
}
func GetProductionsContents() *ProductionsContentService {
	if ProductionsContentSvr == nil {
		ProductionsContentSvr = new(ProductionsContentService)
	}
	return ProductionsContentSvr
}

func (this *ProductionsContentService) OnProductChange(channelName string, msg string) {
	//!!!这里必须异步处理,避免影响数据库写入的效率
	go func() {
		var message struct {
			Oprt string        `json:"oprt"`
			Data BoomGroupInfo `json:"data"`
		}

		err := json.Unmarshal([]byte(msg), &message)
		if err != nil {
			Logger.Error("Failed to unmarshal production data",
				zap.String("channel", channelName),
				zap.String("msg", msg),
				zap.Error(err))
			return
		} else {
			Logger.Info("production data",
				zap.String("channel", channelName),
				zap.String("msg", msg),
				zap.Error(err))
		}

		// Execute different operations
		switch message.Oprt {
		case OprtInsert:
			this.handleInsert(message.Data)
		case OprtUpdate:
			this.handleUpdate(message.Data)
		case OprtDelete:
			this.handleDelete(message.Data)
		default:
			Logger.Warn("Unknown operation type",
				zap.String("channel", channelName),
				zap.String("operation", message.Oprt))
			return
		}

		Logger.Info("ProductionsContentService msg received",
			zap.String("channel", channelName),
			zap.String("operation", message.Oprt),
			zap.Any("data", message.Data))
	}()
}

// 处理插入操作
func (this *ProductionsContentService) handleInsert(data BoomGroupInfo) {
	Logger.Info("Handling insert operation", zap.Any("data", data))
	// TODO: 实现插入逻辑
	// 例如：将数据插入到数据库或其他存储系统
	// 根据配矿单元的id,判断那个productionContent应该被更新

}

// 处理更新操作
func (this *ProductionsContentService) handleUpdate(data BoomGroupInfo) {
	Logger.Info("Handling update operation", zap.Any("data", data))
	// TODO: 实现更新逻辑
	// 例如：更新现有记录

}

// 处理删除操作
func (this *ProductionsContentService) handleDelete(data BoomGroupInfo) {
	Logger.Info("Handling delete operation", zap.Any("data", data))
	// TODO: 实现删除逻辑
	// 例如：从存储系统中删除相关记录
}

func (this *ProductionsContentService) loadProductsStaticAttribs(products map[string]*ProductionContent, md5Sum string, onlyUpdateData bool) {
	Logger.Info("product list will be changed!")
	//if md5 is changed,I will refresh product list.
	if !onlyUpdateData {
		//如果产品的列表项发生了变化,则要更新产品列表
		this.loadProductList(products, md5Sum)
	}
}

func (this *ProductionsContentService) loadProductList(products map[string]*ProductionContent, md5Sum string) {
	defer this.lock.Unlock()
	this.lock.Lock()
	this.productsMd5 = md5Sum
	this.products = products

}

func (this *ProductionsContentService) predictAndSaveMineAssignments4ProductList(dateStrOfFetchData, dateStrOfSavePlan string, minProbability int, forcePredict bool, isAsync bool, predictTime time.Time, calcBoomGroupMode CalcBoomGroupMode, snapName string, maxTimeAsUtc *int, debugProductName string, outputExcelPath *string, cacheInDb bool) {
	if this.MineProductEnableAsActive == 1 {
		//多个产品时,只计算第一个被启用的产品,这个产品是被激活的产品,并且没有被其他产品占用的挖机都由这个产品来使用
		for _, productContent := range this.products {
			diggerObjs, err := diggerProductBinding.ListDiggerOfMineProductExt(productContent.TokenOfMineProduct, productContent.IsDefault)
			for i, diggerObj := range diggerObjs {
				if productContent.IsDefault && diggerObj.MineProductDesp == nil {
					diggerObj.MineProductDesp = RefString(productContent.Name)
					diggerObj.TokenOfMineProduct = RefString(productContent.TokenOfMineProduct)
					diggerObjs[i] = diggerObj
				}
			}
			if err != nil {
				Logger.Error("Failed to list valid diggers", zap.Error(err), zap.String("productContent.Name", productContent.Name))
				return
			}
			if diggerObjs == nil || len(diggerObjs) == 0 {
				Logger.Info("no digger found for product", zap.String("productContent.Name", productContent.Name))
				continue
			}
			if productContent.Name == debugProductName {
				productContent.predictAndSaveMineAssignments4Product(this.MineAssignmentApiUrl, dateStrOfFetchData, dateStrOfSavePlan, NewMineAssignmentReqFromDiggerProductBindings(diggerObjs), minProbability, forcePredict, isAsync, predictTime, calcBoomGroupMode, snapName, maxTimeAsUtc, debugProductName, outputExcelPath, false)
			} else {
				productContent.predictAndSaveMineAssignments4Product(this.MineAssignmentApiUrl, dateStrOfFetchData, dateStrOfSavePlan, NewMineAssignmentReqFromDiggerProductBindings(diggerObjs), minProbability, forcePredict, isAsync, predictTime, calcBoomGroupMode, "", maxTimeAsUtc, "", nil, false)
			}
		}
	} else {
		//多个产品时
		for _, productContent := range this.products {
			diggerObjs, err := diggerProductBinding.ListDiggerOfMineProductExt(productContent.TokenOfMineProduct, productContent.IsDefault)
			if err != nil {
				Logger.Error("Failed to list valid diggers", zap.Error(err), zap.String("productContent.Name", productContent.Name))
				return
			}
			if diggerObjs == nil || len(diggerObjs) == 0 {
				Logger.Info("no digger found for product", zap.String("productContent.Name", productContent.Name))
				continue
			}
			if productContent.Name == debugProductName {
				productContent.predictAndSaveMineAssignments4Product(this.MineAssignmentApiUrl, dateStrOfFetchData, dateStrOfSavePlan, NewMineAssignmentReqFromDiggerProductBindings(diggerObjs), minProbability, forcePredict, isAsync, predictTime, calcBoomGroupMode, snapName, maxTimeAsUtc, debugProductName, outputExcelPath, false)
			} else {
				productContent.predictAndSaveMineAssignments4Product(this.MineAssignmentApiUrl, dateStrOfFetchData, dateStrOfSavePlan, NewMineAssignmentReqFromDiggerProductBindings(diggerObjs), minProbability, forcePredict, isAsync, predictTime, calcBoomGroupMode, "", maxTimeAsUtc, "", nil, false)
			}
		}
	}
}

func (this *ProductionsContentService) GetNow4FetchData() time.Time {
	if this.DesignatedDayAsTodayWhenFetchData == "" {
		return time.Now().In(SystemZone)
	} else {
		return Change2TheDay(time.Now().In(SystemZone), this.DesignatedDayAsTodayWhenFetchData)
	}
}
func (this *ProductionsContentService) GetNow4SaveData() time.Time {
	if this.DesignatedDayAsTodayWhenSavePlan == "" {
		return time.Now().In(SystemZone)
	} else {
		return Change2TheDay(time.Now().In(SystemZone), this.DesignatedDayAsTodayWhenSavePlan)
	}
}
