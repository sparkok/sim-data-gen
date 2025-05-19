package services

/*
*
  - 自动调度:

1.根据挖机的位置,自动确定它正在开采那个配矿单元。
2.根据卡车的位置,自动确定它正在和哪个挖机进行搭配。
3.根据目前的开采情况,自动确定使用了那个采矿方案。
*/
/*
type SimulateDataMaker struct {
	mineProductEnableAsActive int64 // 1 - 使能状态重定义了代表的其实是激活的状态, 0 - 未重定义
	lorryId2Name              map[string]string
	diggerIds                 []string
	lorryIdsInSystem          []string
	lorryNamesInSystem        []string
	boomGroupIds              []string
	boomGroupNames            []string
	boomGroupId2Info          map[string]boomGroupModel.BoomGroup
	mineInfo                  *SystemBaseInfo
}
type DiggerNearby struct {
	DiggerId      string
	Utc           int
	BoomGroupId   string
	BoomGroupName string
	boomGroupId   string
}

func (t *SimulateDataMaker) Generate(productName string, dateAsStr string, excelPath string) {
	accessThreshold := 15.0
	t.mineInfo = NewMineInfo("mineFrom", productName, t.mineProductEnableAsActive)
	t.mineInfo.loadInfo(dateAsStr)

	Logger.Info("SimulateDataMaker.Execute")
	allSql := strings.Builder{}
	allSql.WriteString(t.makeClearSql(dateAsStr))
	allSql.WriteString(";\n")

	//获取产品相关的挖机
	diggerIds, diggerNames, diggerId2Name := t.queryDiggersOfProductName(productName)
	//打开excel文件
	var excelFile *excelize.File
	var err error
	if _, err = os.Stat(excelPath); !os.IsNotExist(err) {
		//存在则先删除
		os.Remove(excelPath)
	}
	if _, err = os.Stat(excelPath); os.IsNotExist(err) {
		excelFile = excelize.NewFile()
	}
	defer excelFile.Close()

	//create a shell
	// 检查是否存在名为 "Sheet1" 的工作表
	sheetLorryNearbyDigger := "Sheet1"
	t.makeSureSheet(excelFile, sheetLorryNearbyDigger)
	lineInLorryNearbyDigger := 1
	excelFile.SetColWidth(sheetLorryNearbyDigger, "A", "A", 18)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "B", "B", 40)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "C", "C", 8)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "D", "D", 40)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "E", "E", 20)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "F", "F", 20)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "G", "G", 40)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "H", "H", 25)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "I", "I", 20)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "K", "K", 40)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "J", "J", 25)
	excelFile.SetColWidth(sheetLorryNearbyDigger, "L", "L", 25)

	t.WriteRow(excelFile, sheetLorryNearbyDigger, lineInLorryNearbyDigger, []interface{}{
		"时间",
		"卡车Id",
		"卡车",
		"最近挖机Id",
		"最近挖机",
		fmt.Sprintf("最近距离(<%.1fm)", accessThreshold),
		"次近挖机Id",
		"次近挖机",
		fmt.Sprintf("次近距离(<%.1fm)", accessThreshold),
		"最近配矿单元Id",
		"最近配矿单元",
		"备注",
	})
	lineInLorryNearbyDigger += 1

	sheetLorrySwitchDigger := "卡车换挖机"
	t.makeSureSheet(excelFile, sheetLorrySwitchDigger)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "A", "A", 18)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "B", "B", 40)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "C", "C", 10)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "E", "E", 10)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "G", "G", 10)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "D", "D", 40)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "F", "F", 40)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "H", "H", 40)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "H", "H", 40)
	excelFile.SetColWidth(sheetLorrySwitchDigger, "I", "I", 150)

	lineInLorrySwitchDigger := 1
	t.WriteRow(excelFile, sheetLorrySwitchDigger, lineInLorrySwitchDigger, []interface{}{
		"时间",
		"卡车Id",
		"卡车",
		"之前挖机Id",
		"之前挖机",
		"挖机Id",
		"挖机",
		"配矿单元Id",
		"配矿单元",
		"说明",
	})
	lineInLorrySwitchDigger++
	t.boomGroupIds, t.boomGroupNames, t.boomGroupId2Info = t.getBoomGroupInfos()
	if len(t.boomGroupIds) == 0 {
		Logger.Error("No boomGroupIds found!")
		return
	}
	//因为挖机开采的时候不一定在配矿单元里面,所以需要留出大概10米的缓冲区,但是这样也就出现了一个问题，做模拟数据的时候,两个配矿单元必须间隔大于这个缓冲区，不然就不会认为切换了
	//因为这是做实验数据,所以姑且可以吧
	digBuffer := 10.0
	digBuffer = 60
	polygonSearchTree := NewSearchTreeByPolygons()
	for boomGroupId, boomGroupInfo := range t.boomGroupId2Info {
		if boomGroupInfo.Geom == nil {
			Logger.Error("boomGroupInfo.Geom is nil", zap.String("boomGroupId", boomGroupId))
			continue
		}
		polygonSearchTree.InsertWktAsPolygon(boomGroupId, *boomGroupInfo.Name, boomGroupInfo.Geom.Wkt, digBuffer)
	}

	t.lorryIdsInSystem, t.lorryNamesInSystem, t.lorryId2Name = t.getLorryInfos()
	if len(t.lorryIdsInSystem) == 0 {
		Logger.Error("No digger IDs found!")
		return
	}

	//获取dateAsStr所在日期的时间范围
	beginUtc, endUtc, err := GetDateAsStrBeginAndEndUtc(dateAsStr)
	if err != nil {
		Logger.Error("GetDateAsStrBeginAndEndUtc", zap.Error(err))
		return
	}

	var diggerLocationsList = make([]DiggerLocations, 0)
	for i, diggerId := range diggerIds {
		locations, err := locationGnssDataDao.ListObjDuring(diggerId, beginUtc, endUtc)
		if err != nil {
			Logger.Error("ListLorryDuring", zap.Error(err))
			return
		}

		diggerLocationsList = append(diggerLocationsList, NewDiggerLocations(diggerId, diggerNames[i], locations))
	}

	var lorryBoundByTheDigger = make(map[string]*DiggerNearby)
	for _, lorryId := range t.lorryIdsInSystem {
		lorryLocations, err := locationGnssDataDao.ListObjDuring(lorryId, beginUtc, endUtc)
		if err != nil {
			Logger.Error("ListLorryDuring", zap.Error(err))
			return
		}
		if len(lorryLocations) == 0 {
			Logger.Info("No lorry locations found!", zap.String("lorryId", lorryId))
			continue
		}

		var lorryLocation locationGnssDataModel.LocationGnssData
		var diggerLocations [][2]float64
		var nearbyDiggerIds []string
		var distances []float64
		var diggerUTCs []int
		for _, lorryLocation = range lorryLocations {
			nearbyDiggerIds, distances, diggerLocations, diggerUTCs = t.FindNearDiggerIdsByLocation(diggerLocationsList, lorryLocation, accessThreshold, 2)
			if len(nearbyDiggerIds) == 0 {
				continue
			}
			Logger.Info("lorry is by diggers", zap.String("lorryId", lorryId), zap.Strings("diggerIds", nearbyDiggerIds))
			//when code reaches here,there are only two conditions: 1. only one digger is near lorry; 2. two diggers are near lorry.
			if len(distances) == 1 {
				//polygonSearchTree的坐标是经纬度,diggerLocations是百万分之一度
				boomGroupId, boomGroupName := extractOneBoomGroupIdAndName(polygonSearchTree.SearchIntersect(diggerLocations[0][0], diggerLocations[0][1], true))
				timePoint := time.Unix(int64(*lorryLocation.Utc), 0).In(SystemZone).Format("01-02 15:04:05")
				if boomGroupId == "" {
					diggerTimePoint := time.Unix(int64(diggerUTCs[0]), 0).In(SystemZone).Format("01-02 15:04:05")
					Logger.Warn("boomGroupId is empty", zap.String("utc", timePoint), zap.String("diggerName", diggerId2Name[nearbyDiggerIds[0]]), zap.String("nearbyDiggerId", nearbyDiggerIds[0]), zap.Float64("x", diggerLocations[0][0]), zap.Float64("y", diggerLocations[0][1]), zap.String("diggerTime", diggerTimePoint))
					continue
				}

				t.WriteRow(excelFile, sheetLorryNearbyDigger, lineInLorryNearbyDigger, []interface{}{
					timePoint,
					*lorryLocation.Tid,
					t.lorryId2Name[*lorryLocation.Tid],
					nearbyDiggerIds[0],
					diggerId2Name[nearbyDiggerIds[0]],
					fmt.Sprintf("%1.f", distances[0]),
					boomGroupId,
					boomGroupName,
				})

				//只能临近一个挖机,如果同时邻近两个挖机时,无法判断切换关系
				//之前接近过挖机1,现在接近挖机2则认为是切换
				switchDigger := false
				var boundDigger *DiggerNearby
				var ok bool
				if boundDigger, ok = lorryBoundByTheDigger[lorryId]; !ok {
					//如果之前没有接近过,则记录
					switchDigger = true
				} else {
					//如果之前接近过,则判断是否切换
					if boundDigger.DiggerId != nearbyDiggerIds[0] || boundDigger.boomGroupId != boomGroupId {
						switchDigger = true
					}
				}
				if switchDigger {
					if boundDigger == nil {
						//首次使用挖机
						t.WriteRow(excelFile, sheetLorrySwitchDigger, lineInLorrySwitchDigger, []interface{}{
							timePoint,
							*lorryLocation.Tid,
							t.lorryId2Name[*lorryLocation.Tid],
							"",
							"",
							nearbyDiggerIds[0],
							diggerId2Name[nearbyDiggerIds[0]],
							boomGroupId,
							boomGroupName,
							fmt.Sprintf("使用挖机 %s", diggerId2Name[nearbyDiggerIds[0]]),
						})
						//allSql.WriteString(t.genSql(diggerId2Name, boundDigger, lorryLocation, dateAsStr, lorryId) + ";\n")
					} else {
						//切换为其他挖机
						t.WriteRow(excelFile, sheetLorrySwitchDigger, lineInLorrySwitchDigger, []interface{}{
							timePoint,
							*lorryLocation.Tid,
							t.lorryId2Name[*lorryLocation.Tid],
							boundDigger.DiggerId,
							diggerId2Name[boundDigger.DiggerId],
							nearbyDiggerIds[0],
							diggerId2Name[nearbyDiggerIds[0]],
							boomGroupId,
							boomGroupName,
							fmt.Sprintf("切换挖机,之前从 %s 到 %s 使用 %s,合计 %.1f 分钟",
								time.Unix(int64(boundDigger.Utc), 0).In(SystemZone).Format("01-02 15:04:05"),
								timePoint,
								diggerId2Name[boundDigger.DiggerId],
								math.Round(float64(*lorryLocation.Utc-boundDigger.Utc)/60.0)),
						})
						allSql.WriteString(t.genSql(diggerId2Name, boundDigger, lorryLocation, dateAsStr, lorryId, boomGroupId, boomGroupName) + ";\n")
					}
					lorryBoundByTheDigger[lorryId] = &DiggerNearby{
						DiggerId: nearbyDiggerIds[0],
						Utc:      *lorryLocation.Utc,
					}
					lineInLorrySwitchDigger++
				}
			} else if len(distances) == 2 {
				boomGroupId, boomGroupName := extractOneBoomGroupIdAndName(polygonSearchTree.SearchIntersect(diggerLocations[0][0], diggerLocations[0][1], true))
				t.WriteRow(excelFile, sheetLorryNearbyDigger, lineInLorryNearbyDigger, []interface{}{
					time.Unix(int64(*lorryLocation.Utc), 0).In(SystemZone).Format("01-02 15:04:05"),
					*lorryLocation.Tid,
					t.lorryId2Name[*lorryLocation.Tid],
					nearbyDiggerIds[0],
					diggerId2Name[nearbyDiggerIds[0]],
					fmt.Sprintf("%1.f", distances[0]),
					nearbyDiggerIds[1],
					diggerId2Name[nearbyDiggerIds[1]],
					fmt.Sprintf("%1.f", distances[1]),
					boomGroupId, boomGroupName,
				})
			}
			lineInLorryNearbyDigger++

		}
		//设置为和挖机绑定结束
		if diggerNearbyDigger, ok := lorryBoundByTheDigger[*lorryLocation.Tid]; !ok {
			//卡车从未接近过挖机,不合理
			timePoint := time.Unix(int64(*lorryLocation.Utc), 0).In(SystemZone).Format("01-02 15:04:05")
			t.WriteRow(excelFile, sheetLorrySwitchDigger, lineInLorrySwitchDigger, []interface{}{
				timePoint,
				*lorryLocation.Tid,
				t.lorryId2Name[*lorryLocation.Tid],
				"",
				"",
				"",
				"",
				"",
				"",
				"卡车从未接近过任何挖机",
			})
			lineInLorrySwitchDigger++
		} else {
			timePoint := time.Unix(int64(*lorryLocation.Utc), 0).In(SystemZone).Format("01-02 15:04:05")
			timeSpan := math.Round(float64(*lorryLocation.Utc-diggerNearbyDigger.Utc) / 60.0)
			//feature := geojson.NewPolygonFeature()
			//ExtractCoordsFromPolygonForFloats(geomObj.(*geom.Polygon), srid)
			t.WriteRow(excelFile, sheetLorrySwitchDigger, lineInLorrySwitchDigger, []interface{}{
				timePoint,
				*lorryLocation.Tid,
				t.lorryId2Name[*lorryLocation.Tid],
				"",
				"",
				diggerNearbyDigger.DiggerId,
				diggerId2Name[diggerNearbyDigger.DiggerId],
				diggerNearbyDigger.BoomGroupId, diggerNearbyDigger.BoomGroupName,
				fmt.Sprintf("挖机使用结束,之前从 %s 到 %s 使用 %s,合计 %.1f 分钟", time.Unix(int64(diggerNearbyDigger.Utc), 0).In(SystemZone).Format("01-02 15:04:05"), timePoint, diggerId2Name[diggerNearbyDigger.DiggerId], timeSpan),
			})
			var diggerName string
			var existed bool
			if diggerName, existed = diggerId2Name[diggerNearbyDigger.DiggerId]; !existed {
				diggerName = ""
			}
			var lorryName string
			if lorryName, existed = t.lorryId2Name[*lorryLocation.Tid]; !existed {
				lorryName = ""
			}

			bindingLog := lorryDiggerBindingLog.LorryDiggerBindingLog{
				DateFlag:         &dateAsStr,
				TokenOfLorry:     &lorryId,
				LorryName:        RefString(lorryName),
				TokenOfDigger:    &diggerNearbyDigger.DiggerId,
				DiggerName:       RefString(diggerName),
				StartUtc:         &diggerNearbyDigger.Utc,
				EndUtc:           lorryLocation.Utc,
				Token:            Uuid(),
				BoomGroupName:    RefString(""),
				TokenOfBoomGroup: RefString(""),
			}
			allSql.WriteString(t.makeCreateSql(bindingLog))
			allSql.WriteString(";\n")

			lineInLorrySwitchDigger++
		}
	}
	//将所以目前的绑定设置为终止
	err = excelFile.SetSheetName("Sheet1", "卡车靠近挖机")
	if err != nil {
		Logger.Error("重命名失败:", zap.Error(err))
		return
	}

	// 设置活动工作表
	excelFile.SetActiveSheet(1)
	// 保存文件
	if err := excelFile.SaveAs(excelPath); err != nil {
		Logger.Error("failed to save excel", zap.Error(err))
	}
	os.WriteFile(strings.Replace(excelPath, ".xlsx", ".sql", 1), []byte(allSql.String()), os.FileMode.Perm(0777))
}

func extractOneBoomGroupIdAndName(polygons []*Polygon) (string, string) {
	if polygons == nil || len(polygons) == 0 {
		return "", ""
	}
	return polygons[0].ID, polygons[0].Name
}
func extractTwoBoomGroupIdAndNames(polygons []*Polygon) (string, string, string, string) {
	if polygons == nil || len(polygons) == 0 {
		return "", "", "", ""
	}
	return polygons[0].ID, polygons[0].Name, polygons[1].ID, polygons[1].Name
}

func (t *SimulateDataMaker) genSql(diggerId2Name map[string]string, diggerNearbyDigger *DiggerNearby, lorryLocation locationGnssDataModel.LocationGnssData, dateAsStr string, lorryId string, boomGroupId string, boomGroupName string) string {
	//生成sql
	var diggerName string
	var existed bool
	if diggerNearbyDigger == nil {
		return ""
	}
	if diggerName, existed = diggerId2Name[diggerNearbyDigger.DiggerId]; !existed {
		diggerName = ""
	}
	var lorryName string
	if lorryName, existed = t.lorryId2Name[*lorryLocation.Tid]; !existed {
		lorryName = ""
	}

	bindingLog := lorryDiggerBindingLog.LorryDiggerBindingLog{
		DateFlag:         &dateAsStr,
		TokenOfLorry:     &lorryId,
		LorryName:        RefString(lorryName),
		TokenOfDigger:    &diggerNearbyDigger.DiggerId,
		DiggerName:       RefString(diggerName),
		StartUtc:         &diggerNearbyDigger.Utc,
		EndUtc:           lorryLocation.Utc,
		Token:            Uuid(),
		BoomGroupName:    RefString(""),
		TokenOfBoomGroup: RefString(""),
	}
	sql := t.makeCreateSql(bindingLog)
	return sql
}

func NilAsEmpty(value string, existed bool) string {
	if existed {
		return value
	}
	return ""
}
func (t *SimulateDataMaker) makeClearSql(dateAsStr string) string {
	sql := GetDb().ToSQL(func(tx *gorm.DB) *gorm.DB {
		condition := lorryDiggerBindingLog.LorryDiggerBindingLog{
			DateFlag: &dateAsStr,
		}
		return tx.Where(condition).Delete(condition)
	})
	return sql
}

func (t *SimulateDataMaker) makeCreateSql(bindingLog lorryDiggerBindingLog.LorryDiggerBindingLog) string {
	sql := GetDb().ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(bindingLog)
	})
	return sql
}

func (t *SimulateDataMaker) makeSureSheet(f *excelize.File, sheetName string) (error, int) {
	index, err := f.GetSheetIndex(sheetName)
	if err != nil {
		Logger.Error("GetSheetIndex", zap.Error(err))
		return err, -1
	}
	if index == -1 {
		// 如果不存在，则创建工作表
		index, err = f.NewSheet(sheetName)
	}
	return err, index
}

type LorryToDigger struct {
	Distance float64
	DiggerId string
	LorryId  string
}

func (t *SimulateDataMaker) FindNearDiggerIdsByLocation(diggers []DiggerLocations, lorry locationGnssDataModel.LocationGnssData, rByMeter float64, maxSize int) ([]string, []float64, [][2]float64, []int) {
	ids := make([]string, 0)
	distances := make([]float64, 0)
	locations := make([][2]float64, 0)
	diggerUTCs := make([]int, 0)
	for _, digger := range diggers {
		diggerLocation, found := digger.FindLatestLocation(int64(*lorry.Utc))
		if !found {
			continue
		}
		if *lorry.Utc-*diggerLocation.Utc > 60 {
			continue
		}

		distance := CalculateDistanceAsMeter(*diggerLocation.X, *diggerLocation.Y, *lorry.X, *lorry.Y)

		if distance < rByMeter {
			//Logger.Info(fmt.Sprintf("lorry(%s) is near digger(%s)", *lorry.Token, digger.DiggerId))
			ids = append(ids, digger.DiggerId)
			distances = append(distances, distance)
			diggerUTCs = append(diggerUTCs, *diggerLocation.Utc)
			locations = append(locations, [2]float64{*lorry.X, *lorry.Y})
		}
	}
	//对distances按从小到大排序,同时调整对应索引的ids
	if len(distances) > 1 {
		for i := 0; i < len(distances)-1; i++ {
			for j := 0; j < len(distances)-1-i; j++ {
				if distances[j] > distances[j+1] {
					distances[j], distances[j+1] = distances[j+1], distances[j]
					ids[j], ids[j+1] = ids[j+1], ids[j]
					locations[j], locations[j+1] = locations[j+1], locations[j]
					diggerUTCs[j], diggerUTCs[j+1] = diggerUTCs[j+1], diggerUTCs[j]
				}
			}
		}
	}
	if maxSize > 0 && len(ids) > maxSize {
		ids = ids[:maxSize]
		distances = distances[:maxSize]
		locations = locations[:maxSize]
		diggerUTCs = diggerUTCs[:maxSize]
	}
	return ids, distances, locations, diggerUTCs
}

func (t *SimulateDataMaker) HeightDiff(result *SearchResult, digger locationModel.Location) bool {
	if result.Elevation == nil || digger.Elevation == nil {
		return false
	}
	if math.Abs(*result.Elevation-*digger.Elevation) < 7 {
		return false
	}
	return true
}

func (t *SimulateDataMaker) queryDiggersOfProductName(ProductName string) (ids []string, names []string, id2Name map[string]string) {
	ids = make([]string, 0)
	names = make([]string, 0)
	id2Name = make(map[string]string)
	product, err := mineProduct.GetObjByName(&ProductName)
	if err != nil {
		Logger.Error("mineProduct.GetObjByName", zap.Error(err))
		return
	}
	includeNoBoundDigger := false
	if t.mineProductEnableAsActive == 1 && (product.Status == nil || *product.Status == PRODUCT_STATUS_ENABLED) {
		includeNoBoundDigger = true
		//0 激活的产品, 1 其他的产品
	} else if product.Status == nil || *product.Status != PRODUCT_STATUS_ENABLED {
		//0 正常, 1 关闭
		return
	}
	var bindings []diggerProductBindingModel.DiggerProductBindingFully1
	//这里并没有支持缺省产品,因为这个自动设置器未来可能会弃用,因此简单实现
	bindings, err = diggerProductBinding.ListDiggerOfMineProductExt(*product.Token, includeNoBoundDigger)
	if len(bindings) == 0 {
		Logger.Error("diggerProductBinding.ListDiggerOfMineProductExt", zap.Error(err))
		return
	}
	for _, binding := range bindings {
		ids = append(ids, *binding.TokenOfDigger)
		names = append(names, *binding.DiggerDesp)
		id2Name[*binding.TokenOfDigger] = *binding.DiggerDesp
	}
	return
}

func (t *SimulateDataMaker) getBoomGroupInfos() (ids []string, names []string, boomGroupMap map[string]boomGroupModel.BoomGroup) {
	Logger.Info("Fetching boomGroup IDs...")
	list, err := boomGroupDao.ListObj()
	if err != nil || len(list) == 0 {
		Logger.Error("boomGroupDao.ListObj", zap.Error(err))
		return
	}

	ids = make([]string, len(list))
	names = make([]string, len(list))
	boomGroupMap = make(map[string]boomGroupModel.BoomGroup)
	for i, boomGroup := range list {
		ids[i] = *boomGroup.Token
		names[i] = *boomGroup.Name
		boomGroupMap[*boomGroup.Token] = boomGroup
	}
	return
}

func (t *SimulateDataMaker) getLorryInfos() (ids []string, names []string, lorryMap map[string]string) {
	Logger.Info("Fetching lorry IDs...")
	list, err := lorryDao.ListObj()
	if err != nil {
		Logger.Error("lorryDiggerBindingLogDao.ListObj", zap.Error(err))
		return
	}
	if len(list) == 0 {
		Logger.Error("diggerProductBinding.ListDiggerOfMineProductExt", zap.Error(err))
		return
	}

	ids = make([]string, len(list))
	names = make([]string, len(list))
	lorryMap = make(map[string]string)
	for i, lorry := range list {
		ids[i] = *lorry.Token
		names[i] = *lorry.Name
		lorryMap[*lorry.Token] = *lorry.Name
	}
	return
}

func (t *SimulateDataMaker) GetValidLocationsByIds(ids []string, minUtc int64) []locationModel.Location {
	var locations []locationModel.Location
	var err error
	locations, err = locationDao.GetValidLocationsByIds(ids, minUtc)
	if err != nil {
		Logger.Error("locationDao.GetLocationsByIds", zap.Error(err))
		return locations
	}
	return locations
}

func (t *SimulateDataMaker) WriteRow(f *excelize.File, sheetName string, row int, values []interface{}) {
	// 定义要写入的数据
	// 写入一行数据
	cell := fmt.Sprintf("A%d", row)
	err := f.SetSheetRow(sheetName, cell, &values)
	if err != nil {
		Logger.Error("failed to set sheet row", zap.Error(err))
		return
	}
}

func NewSimulateDataMaker(varNamePrx string) *SimulateDataMaker {
	dataGenerator := new(SimulateDataMaker)
	dataGenerator.mineProductEnableAsActive = GetConfig().Int64("productions.mineProductEnableAsActive", 1)
	return dataGenerator
}

func GetDateAsStrBeginAndEndUtc(dateAsStr string) (int64, int64, error) {
	//获取指定日期的起始和结束时间戳
	var beginTime, endTime time.Time
	var err error
	beginTime, err = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 00:00:00", dateAsStr), SystemZone)
	if err != nil {
		return 0, 0, err
	}
	endTime, err = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 23:59:59", dateAsStr), SystemZone)
	if err != nil {
		return 0, 0, err
	}
	return beginTime.Unix(), endTime.Unix(), nil
}
*/
