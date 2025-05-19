package services

/*
*
  - 自动调度:

1.根据挖机的位置,自动确定它正在开采那个配矿单元。
2.根据卡车的位置,自动确定它正在和哪个挖机进行搭配。
3.根据目前的开采情况,自动确定使用了那个采矿方案。
*/
//type DiggerSwitchBoomGroupLogSqlMaker struct {
//	mineProductEnableAsActive int64 // 1 - 使能状态重定义了代表的其实是激活的状态, 0 - 未重定义
//	lorryId2Name              map[string]string
//	diggerIds                 []string
//	lorryIdsInSystem          []string
//	lorryNamesInSystem        []string
//	boomGroupIds              []string
//	boomGroupNames            []string
//	boomGroupId2Info          map[string]boomGroupModel.BoomGroup
//	//mineInfo               *SystemBaseInfo
//}
//
//func (t *DiggerSwitchBoomGroupLogSqlMaker) queryDiggersOfProductName(ProductName string) (ids []string, names []string, id2Name map[string]string) {
//	ids = make([]string, 0)
//	names = make([]string, 0)
//	id2Name = make(map[string]string)
//	product, err := mineProduct.GetObjByName(&ProductName)
//	if err != nil {
//		Logger.Error("mineProduct.GetObjByName", zap.Error(err))
//		return
//	}
//	includeNoBoundDigger := false
//	if t.mineProductEnableAsActive == 1 && (product.Status == nil || *product.Status == PRODUCT_STATUS_ENABLED) {
//		includeNoBoundDigger = true
//		//0 激活的产品, 1 其他的产品
//	} else if product.Status == nil || *product.Status != PRODUCT_STATUS_ENABLED {
//		//0 正常, 1 关闭
//		return
//	}
//	var bindings []diggerProductBindingModel.DiggerProductBindingFully1
//	//这里并没有支持缺省产品,因为这个自动设置器未来可能会弃用,因此简单实现
//	bindings, err = diggerProductBinding.ListDiggerOfMineProductExt(*product.Token, includeNoBoundDigger)
//	if len(bindings) == 0 {
//		Logger.Error("diggerProductBinding.ListDiggerOfMineProductExt", zap.Error(err))
//		return
//	}
//	for _, binding := range bindings {
//		ids = append(ids, *binding.TokenOfDigger)
//		names = append(names, *binding.DiggerDesp)
//		id2Name[*binding.TokenOfDigger] = *binding.DiggerDesp
//	}
//	return
//}
//
//func (t *DiggerSwitchBoomGroupLogSqlMaker) Scan(productName string, dateFlag string, accessThreshold float64, excelPath string) {
//	//t.mineInfo = NewMineInfo("mineFrom", productName, accessThreshold, t.mineProductEnableAsActive, 300)
//	Logger.Info("DiggerSwitchBoomGroupLogSqlMaker.Execute")
//	allSql := strings.Builder{}
//	allSql.WriteString(t.makeClearSql(dateFlag))
//	allSql.WriteString(";\n")
//
//	//打开excel文件
//	var excelFile *excelize.File
//	var err error
//	if _, err = os.Stat(excelPath); !os.IsNotExist(err) {
//		//存在则先删除
//		os.Remove(excelPath)
//	}
//	if _, err = os.Stat(excelPath); os.IsNotExist(err) {
//		excelFile = excelize.NewFile()
//	}
//	defer excelFile.Close()
//
//	//weightBridges, err := t.queryWeightBridges(productName, dateFlag)
//	//if err != nil {
//	//	Logger.Error("queryWeightBridges", zap.Error(err))
//	//	return
//	//}
//
//	//create a shell
//	// 检查是否存在名为 "Sheet1" 的工作表
//	sheetDiggerSwitchBoomGroupLog := "Sheet1"
//	t.makeSureSheet(excelFile, sheetDiggerSwitchBoomGroupLog)
//	//将所以目前的绑定设置为终止
//	err = excelFile.SetSheetName("Sheet1", "挖机和配矿单元切换")
//	if err != nil {
//		Logger.Error("重命名失败:", zap.Error(err))
//		return
//	}
//
//	sheetDiggerSwitchBoomGroupLog = "挖机和配矿单元切换"
//	lineInWeighLogger := 1
//	excelFile.SetColWidth(sheetDiggerSwitchBoomGroupLog, "A", "A", 18)
//	excelFile.SetColWidth(sheetDiggerSwitchBoomGroupLog, "B", "B", 40)
//	excelFile.SetColWidth(sheetDiggerSwitchBoomGroupLog, "C", "C", 40)
//	excelFile.SetColWidth(sheetDiggerSwitchBoomGroupLog, "D", "D", 20)
//	excelFile.SetColWidth(sheetDiggerSwitchBoomGroupLog, "E", "E", 40)
//	excelFile.SetColWidth(sheetDiggerSwitchBoomGroupLog, "F", "F", 20)
//
//	t.WriteRow(excelFile, sheetDiggerSwitchBoomGroupLog, lineInWeighLogger, []interface{}{
//		"Id",
//		"时间",
//		"挖机Id",
//		"挖机编号",
//		"配矿单元Id",
//		"配矿单元",
//	})
//	lineInWeighLogger += 1
//
//	t.boomGroupIds, t.boomGroupNames, t.boomGroupId2Info = t.getBoomGroupInfos()
//	if len(t.boomGroupIds) == 0 {
//		Logger.Error("No boomGroupIds found!")
//		return
//	}
//	//因为挖机开采的时候不一定在配矿单元里面,所以需要留出大概10米的缓冲区,但是这样也就出现了一个问题，做模拟数据的时候,两个配矿单元必须间隔大于这个缓冲区，不然就不会认为切换了
//	//因为这是做实验数据,所以姑且可以吧
//	digBuffer := 10.0
//	digBuffer = 60
//	polygonSearchTree := NewSearchTreeByPolygons()
//
//	for boomGroupId, boomGroupInfo := range t.boomGroupId2Info {
//		if boomGroupInfo.Geom == nil {
//			Logger.Error("boomGroupInfo.Geom is nil", zap.String("boomGroupId", boomGroupId))
//			continue
//		}
//		Logger.Info("insert " + *boomGroupInfo.Name)
//		polygonSearchTree.InsertWktAsPolygon(boomGroupId, *boomGroupInfo.Name, boomGroupInfo.Geom.Wkt, digBuffer)
//	}
//
//	//获取dateAsStr所在日期的时间范围
//	beginUtc, endUtc, err := GetDateAsStrBeginAndEndUtc(dateFlag)
//	if err != nil {
//		Logger.Error("GetDateAsStrBeginAndEndUtc", zap.Error(err))
//		return
//	}
//
//	diggerIds, _, diggerId2Name := t.queryDiggersOfProductName(productName)
//
//	//var diggerLocationsList = make([]DiggerLocations, 0)
//	diggerLoc := make(map[string]string)
//	for _, diggerId := range diggerIds {
//		locations, err := locationGnssDataDao.ListObjDuring(diggerId, beginUtc, endUtc)
//		if err != nil {
//			Logger.Error("ListLorryDuring", zap.Error(err))
//			return
//		}
//
//		var provBoomGroup string = ""
//		var firstBoomGroup = false
//		for _, loc := range locations {
//			if *loc.Utc == 1744262644 {
//				Logger.Info("2025-04-10 13:24:04")
//			}
//			//这里是制作演示数据,必须让两个配矿单元离得足够远,根据挖机位置判断配矿单元时只允许计算到一个配矿单元
//			boomGroupPolygons := polygonSearchTree.SearchIntersect(*loc.X, *loc.Y, false)
//			if len(boomGroupPolygons) == 0 {
//				Logger.Info("boomGroupPolygons is empty", zap.String("diggerId", diggerId), zap.Int("len(boomGroupPolygons)", len(boomGroupPolygons)), zap.String("loc", fmt.Sprintf("%f,%f", *loc.X, *loc.Y)))
//				continue
//			} else if len(boomGroupPolygons) > 1 {
//				Logger.Info("boomGroupPolygons size > 1", zap.String("diggerId", diggerId), zap.Int("len(boomGroupPolygons)", len(boomGroupPolygons)), zap.String("loc", fmt.Sprintf("%f,%f", *loc.X, *loc.Y)))
//				continue
//			}
//			if provBoomGroup != boomGroupPolygons[0].ID {
//				if _, firstBoomGroup = diggerLoc[diggerId]; !firstBoomGroup {
//					diggerLoc[diggerId] = boomGroupPolygons[0].ID
//				}
//				locTime := *loc.Utc
//				if !firstBoomGroup {
//					locTime = int(beginUtc)
//				}
//				//记录切换日志
//				t.WriteRow(excelFile, sheetDiggerSwitchBoomGroupLog, lineInWeighLogger, []interface{}{
//					*loc.Token,
//					time.Unix(int64(locTime), 0).In(SystemZone).Format(STANDARD_TIME_FORMAT),
//					diggerId,
//					diggerId2Name[diggerId],
//					boomGroupPolygons[0].ID,
//					boomGroupPolygons[0].Name,
//				})
//				allSql.WriteString(t.makeInsertDiggerSwitchBoomGroupLogSql(*loc.Token, dateFlag, *loc.Utc, diggerId, diggerId2Name[diggerId], boomGroupPolygons[0].ID, boomGroupPolygons[0].Name))
//				allSql.WriteString(";\n")
//				provBoomGroup = boomGroupPolygons[0].ID
//				lineInWeighLogger++
//			}
//		}
//	}
//
//	// 设置活动工作表
//	excelFile.SetActiveSheet(1)
//	// 保存文件
//	if err := excelFile.SaveAs(excelPath); err != nil {
//		Logger.Error("failed to save excel", zap.Error(err))
//	}
//	os.WriteFile(strings.Replace(excelPath, ".xlsx", ".sql", 1), []byte(allSql.String()), os.FileMode.Perm(0777))
//}
//func (t *DiggerSwitchBoomGroupLogSqlMaker) getBoomGroupInfos() (ids []string, names []string, boomGroupMap map[string]boomGroupModel.BoomGroup) {
//	Logger.Info("Fetching boomGroup IDs...")
//	list, err := boomGroupDao.ListObj()
//	if err != nil || len(list) == 0 {
//		Logger.Error("boomGroupDao.ListObj", zap.Error(err))
//		return
//	}
//
//	ids = make([]string, len(list))
//	names = make([]string, len(list))
//	boomGroupMap = make(map[string]boomGroupModel.BoomGroup)
//	for i, boomGroup := range list {
//		ids[i] = *boomGroup.Token
//		names[i] = *boomGroup.Name
//		boomGroupMap[*boomGroup.Token] = boomGroup
//	}
//	return
//}
//func (t *DiggerSwitchBoomGroupLogSqlMaker) makeSureSheet(f *excelize.File, sheetName string) (error, int) {
//	index, err := f.GetSheetIndex(sheetName)
//	if err != nil {
//		Logger.Error("GetSheetIndex", zap.Error(err))
//		return err, -1
//	}
//	if index == -1 {
//		// 如果不存在，则创建工作表
//		index, err = f.NewSheet(sheetName)
//	}
//	return err, index
//}
//
//func (t *DiggerSwitchBoomGroupLogSqlMaker) WriteRow(f *excelize.File, sheetName string, row int, values []interface{}) {
//	// 定义要写入的数据
//	// 写入一行数据
//	cell := fmt.Sprintf("A%d", row)
//	err := f.SetSheetRow(sheetName, cell, &values)
//	if err != nil {
//		Logger.Error("failed to set sheet row", zap.Error(err))
//		return
//	}
//}
//
//func (t *DiggerSwitchBoomGroupLogSqlMaker) queryWeightBridges(siteCode, dateStr string) ([]weighLogger.WeighLogger, error) {
//	startFrom, endTo, err := t.GetDateAsStrBeginAndEnd(dateStr)
//	if err != nil {
//		return nil, err
//	}
//	return weighLoggerDao.ListObjDuring(siteCode, startFrom, endTo)
//}
//
//func NewDiggerSwitchBoomGroupLogSqlMaker(varNamePrx string) *DiggerSwitchBoomGroupLogSqlMaker {
//	lorrySwitchScanner := new(DiggerSwitchBoomGroupLogSqlMaker)
//	lorrySwitchScanner.mineProductEnableAsActive = GetConfig().Int64("productions.mineProductEnableAsActive", 1)
//	return lorrySwitchScanner
//}
//
//func (t *DiggerSwitchBoomGroupLogSqlMaker) GetDateAsStrBeginAndEnd(dateAsStr string) (time.Time, time.Time, error) {
//	//获取指定日期的起始和结束时间戳
//	var beginTime, endTime time.Time
//	var err error
//	beginTime, err = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", dateAsStr), SystemZone)
//	if err != nil {
//		return beginTime, endTime, err
//	}
//	endTime, err = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 23:59:59", dateAsStr), SystemZone)
//	if err != nil {
//		return beginTime, endTime, err
//	}
//	return beginTime, endTime, nil
//}
//
//func (t *DiggerSwitchBoomGroupLogSqlMaker) GetLorryId(vehicleNo string, lorryName2Id map[string]string) (string, error) {
//	if id, ok := lorryName2Id[vehicleNo]; ok {
//		return id, nil
//	} else {
//		return "", nil
//	}
//}
//
////	func (t *DiggerSwitchBoomGroupLogSqlMaker) CalcMineInfo(lorryId string, checkTime *string) (string, string, error) {
////		var diggerId, boomGroupId string
////		var found bool
////		checkTimeObj, err := time.ParseInLocation("2006-01-02 15:04:05", *checkTime, SystemZone)
////		if err != nil {
////			Logger.Error("failed to parse bridge time", zap.Error(err))
////			return "", "", err
////		}
////		diggerId, boomGroupId, found = t.mineInfo.calcBoomGroupDugByDigLogger(lorryId, checkTimeObj.Unix())
////		if found != true {
////			Logger.Error("CalcMineInfo return false")
////			return "", "", errors.New("CalcMineInfo return nothing")
////		}
////
////		return diggerId, boomGroupId, nil
////	}
//func (t *DiggerSwitchBoomGroupLogSqlMaker) makeClearSql(dateAsStr string) string {
//	sql := GetDb().ToSQL(func(tx *gorm.DB) *gorm.DB {
//		condition := diggerSwitchBoomGroupLog.DiggerSwitchBoomGroupLog{
//			DateFlag: &dateAsStr,
//		}
//		return tx.Where(condition).Delete(condition)
//	})
//	return sql
//}
//
//func (t *DiggerSwitchBoomGroupLogSqlMaker) makeInsertDiggerSwitchBoomGroupLogSql(token string, dateAsStr string, applyUtc int, diggerId, diggerName, boomGroupId, boomGroupName string) string {
//	return GetDb().ToSQL(func(tx *gorm.DB) *gorm.DB {
//		switchLog := diggerSwitchBoomGroupLog.DiggerSwitchBoomGroupLog{
//			Token:       RefString(token),
//			ApplyUtc:    RefInt(applyUtc),
//			DateFlag:    RefString(dateAsStr),
//			BoomGroupId: RefString(boomGroupId),
//			Name:        RefString(diggerName + " use " + boomGroupName),
//			DiggerId:    RefString(diggerId),
//			SubmitUtc:   RefInt(applyUtc),
//			Status:      RefInt(2),
//		}
//		return tx.Create(switchLog)
//	})
//}

//func (t *DiggerSwitchBoomGroupLogSqlMaker) makeCreateSql(switchLog diggerSwitchBoomGroupLog.DiggerSwitchBoomGroupLog) string {
//	sql := GetDb().ToSQL(func(tx *gorm.DB) *gorm.DB {
//		return tx.Create(switchLog)
//	})
//	return sql
//}
