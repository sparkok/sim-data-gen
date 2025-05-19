// filePath: sim_data_gen/services/maker.go
package services

import (
	"bytes"
	"context" // Added for GORM logger interface
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	boomGroupModel "sim_data_gen/models/boomGroup"
	locationGnssDataModel "sim_data_gen/models/locationGnssData"
	weighLoggerModel "sim_data_gen/models/weighLogger"
	yAnalyserModel "sim_data_gen/models/yAnalyser"
	"sim_data_gen/utils" // Your zap logger package
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ProductMaterialInfo
type ProductMaterialInfo struct {
	Name                 string
	Min                  float64
	Max                  float64
	Required             bool
	FieldNameInBoomGroup string
}

// SimulateDataMaker
type SimulateDataMaker struct {
	dataPath           string
	productName        string
	dateStr            string
	baseTime           time.Time
	pathConfig         *PathConfig
	outputFileIndex    int
	productMaterials   map[string]ProductMaterialInfo
	currentTime        time.Time
	allEvents          *EventCollection
	currentFileEndTime time.Time
	gormDB             *gorm.DB
	logger             *zap.Logger
}

// LorryState and DiggerState structs
type LorryState struct {
	ID                   string
	Name                 string
	AssignedDiggerID     string
	StatusString         string
	CurrentBusiNo        string
	CurrentPos           Point
	CurrentTime          time.Time
	IsLoaded             bool
	TripCount            int
	CurrentBoomGroupData *boomGroupModel.BoomGroup
	LastNetWeightKg      float64
	LastTareWeightKg     float64
}
type DiggerState struct {
	ID                  string
	Name                string
	AssignedBoomGroupID string
	StatusString        string
	CurrentPos          Point
	CurrentTime         time.Time
}

// NewSimulateDataMaker
func NewSimulateDataMaker(dataPath, productName, dateStr string) (*SimulateDataMaker, error) {
	logger := utils.Logger.Named("SimulateDataMaker")
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		logger.Error("Invalid date string format", zap.String("dateStr", dateStr), zap.Error(err))
		return nil, fmt.Errorf("invalid date string format '%s': %w", dateStr, err)
	}
	materials := map[string]ProductMaterialInfo{
		"CaO": {Name: "CaO", Min: 40, Max: 47, Required: true, FieldNameInBoomGroup: "Material1"},
		"MgO": {Name: "MgO", Min: 0, Max: 2, Required: true, FieldNameInBoomGroup: "Material2"},
	}
	if productName != "CP1" {
		logger.Warn("Product is not CP1, using default material definitions", zap.String("productName", productName))
	}

	gormZapLogger := gormlogger.New(
		NewZapGormLoggerAdapter(utils.Logger.Named("gorm")), // Use the adapter
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  gormlogger.Warn, // Set GORM log level (Info, Warn, Error, Silent) - Warn to reduce noise
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: gormZapLogger})
	if err != nil {
		logger.Error("Failed to initialize GORM DB", zap.Error(err))
		return nil, fmt.Errorf("failed to initialize GORM DB: %w", err)
	}
	err = db.AutoMigrate(&locationGnssDataModel.LocationGnssData{}, &weighLoggerModel.WeighLogger{}, &yAnalyserModel.YAnalyser{})
	if err != nil {
		logger.Error("Failed to auto-migrate GORM schemas", zap.Error(err))
		return nil, fmt.Errorf("failed to auto-migrate GORM schemas: %w", err)
	}

	return &SimulateDataMaker{
		dataPath: dataPath, productName: productName, dateStr: dateStr, baseTime: parsedDate,
		outputFileIndex: 0, productMaterials: materials, allEvents: NewEventCollection(),
		currentTime: parsedDate.Add(6 * time.Hour), gormDB: db, logger: logger,
	}, nil
}

// --- GORM SQL Generation Helpers ---
func formatVarAsSQLLiteral(v interface{}) string {
	if v == nil {
		return "NULL"
	}
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(val, "'", "''"))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32:
		return fmt.Sprintf("%g", val)
	case float64:
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "1"
		} else {
			return "0"
		}
	case time.Time:
		return fmt.Sprintf("'%s'", val.UTC().Format("2006-01-02 15:04:05.000"))
	case *string:
		if val == nil {
			return "NULL"
		}
		return formatVarAsSQLLiteral(*val)
	case *int:
		if val == nil {
			return "NULL"
		}
		return formatVarAsSQLLiteral(*val)
	case *int64:
		if val == nil {
			return "NULL"
		}
		return formatVarAsSQLLiteral(*val)
	case *float32:
		if val == nil {
			return "NULL"
		}
		return formatVarAsSQLLiteral(*val)
	case *float64:
		if val == nil {
			return "NULL"
		}
		return formatVarAsSQLLiteral(*val)
	case *bool:
		if val == nil {
			return "NULL"
		}
		return formatVarAsSQLLiteral(*val)
	case *time.Time:
		if val == nil {
			return "NULL"
		}
		return formatVarAsSQLLiteral(*val)
	default:
		utils.Logger.Warn("interpolateSQL: Unhandled type, attempting fmt.Sprintf", zap.Any("value", v), zap.String("type", fmt.Sprintf("%T", v)))
		return fmt.Sprintf("'%v'", v)
	}
}
func interpolateSQL(sqlTemplate string, vars []interface{}) string {
	interpolatedSQL := sqlTemplate
	for _, v := range vars {
		literal := formatVarAsSQLLiteral(v)
		interpolatedSQL = strings.Replace(interpolatedSQL, "?", literal, 1)
	}
	if strings.Contains(interpolatedSQL, "?") {
		utils.Logger.Warn("SQL interpolation might be incomplete", zap.String("sql", interpolatedSQL), zap.Any("vars", vars))
	}
	return interpolatedSQL
}
func (m *SimulateDataMaker) generateSQLWithGORM(modelInstance interface{}) string {
	stmt := m.gormDB.Session(&gorm.Session{DryRun: true}).Create(modelInstance).Statement
	if stmt.SQL.String() == "" && stmt.Error != nil {
		m.logger.Error("Error in GORM DryRun session", zap.Error(stmt.Error), zap.Any("model", modelInstance))
		return fmt.Sprintf("-- ERROR GENERATING SQL: %v", stmt.Error)
	}
	return interpolateSQL(stmt.SQL.String(), stmt.Vars)
}

// --- Mocked API Implementations ---
func (m *SimulateDataMaker) getBoomGroupInfos() (ids []string, names []string, boomGroupMap map[string]boomGroupModel.BoomGroup) {
	boomGroupMap = make(map[string]boomGroupModel.BoomGroup)
	numBoomGroups := 10
	statusActive := "Active"
	defaultNt := "自动生成"
	defaultTag := "Simulated"
	defaultPile := "default_pile"
	defaultUsed := 0.0
	defaultDist := 100.0
	baseX, baseY := 120.0010e6, 30.0010e6
	baseElevHigh, baseElevLow := 90.0, 80.0
	for i := 1; i <= numBoomGroups; i++ {
		id := fmt.Sprintf("BG%02d", i)
		name := fmt.Sprintf("配矿单元%02d", i)
		ids = append(ids, id)
		names = append(names, name)
		x := baseX + float64(i*500)
		y := baseY + float64((i%3)*300)
		elevH := baseElevHigh + float64(i)
		elevL := baseElevLow + float64(i)
		cao := 40.0 + rand.Float64()*(47.0-40.0)
		mgo := 0.0 + rand.Float64()*(2.0-0.0)
		boomGroupMap[id] = boomGroupModel.BoomGroup{
			Token: strPtr(id), Name: strPtr(name), X: &x, Y: &y, High: &elevH, Low: &elevL, Status: &statusActive, Nt: &defaultNt, Tag: &defaultTag, TokenOfPile: &defaultPile, Used: &defaultUsed, Distance: &defaultDist, Material1: &cao, Material2: &mgo,
		}
	}
	return ids, names, boomGroupMap
}
func (m *SimulateDataMaker) getDiggerInfosOfProduct(productName string) (ids []string, names []string, id2Name map[string]string) {
	ids = []string{"Digger1", "Digger2"}
	names = []string{"挖机1", "挖机2"}
	id2Name = map[string]string{"Digger1": "挖机1", "Digger2": "挖机2"}
	return
}
func (m *SimulateDataMaker) getLorryInfos() (ids []string, names []string, id2Name map[string]string) {
	id2Name = make(map[string]string)
	numLorries := 8
	for i := 1; i <= numLorries; i++ {
		id := fmt.Sprintf("Truck%d", i)
		name := fmt.Sprintf("卡车%d", i)
		ids = append(ids, id)
		names = append(names, name)
		id2Name[id] = name
	}
	return
}
func (m *SimulateDataMaker) getCenterOfBoomGroup(id string) (lng float64, lat float64) {
	_, _, bgMap := m.getBoomGroupInfos()
	if bg, ok := bgMap[id]; ok && bg.X != nil && bg.Y != nil {
		return *bg.X, *bg.Y
	}
	m.logger.Warn("Boom group not found or X/Y nil in getCenterOfBoomGroup", zap.String("boomGroupID", id))
	return 0, 0
}
func (m *SimulateDataMaker) getCenterOfUnloadSite(productName string) (lng float64, lat float64) {
	if productName == "CP1" {
		return 120.005e6, 30.005e6
	}
	m.logger.Warn("Unload site for product not defined, returning 0,0", zap.String("productName", productName))
	return 0, 0
}

// --- SQL Data Object Preparation and Generation ---
func (m *SimulateDataMaker) generateLocationGnssSQL(entityID, entityName string, point Point, utcTime time.Time, heading float32, speedKmh float32, status int, alarm int, nt *string) string {
	token := uuid.NewString()
	utcUnix := int(utcTime.Unix())
	lastCommUtc := utcUnix
	dbX := point.Lng / 1e6
	dbY := point.Lat / 1e6
	locData := locationGnssDataModel.LocationGnssData{
		Alarm: &alarm, Heading: &heading, Name: entityName, Speed: &speedKmh, Status: &status,
		Token: &token, Utc: &utcUnix, X: &dbX, Y: &dbY, LastCommUtc: &lastCommUtc, Tid: &entityID, Nt: nt,
	}
	return m.generateSQLWithGORM(&locData)
}
func (m *SimulateDataMaker) generateWeighLoggerSQL(lorryName, busiNo string, checkTime time.Time, direction string, grossW, netW, tareW float64, matCode, matName string, nt *string) string {
	token := uuid.NewString()
	updateAt := time.Now().UTC()
	checkTimeStr := checkTime.Format(DateTimeFormat)
	weighData := weighLoggerModel.WeighLogger{
		BusiNo: &busiNo, CheckTime: &checkTimeStr, Direction: &direction, GrossWeight: &grossW, NetWeight: &netW, Nt: nt,
		SiteCode: &m.productName, SiteName: &m.productName, TareWeight: &tareW, Token: &token, UpdateAt: &updateAt,
		VehicleNo: &lorryName, MaterialCode: &matCode, MaterialName: &matName,
	}
	return m.generateSQLWithGORM(&weighData)
}
func (m *SimulateDataMaker) generateYAnalyserSQL(testAtTime time.Time, flux, load, beltSpeed float64, status int, matGrades map[string]float64) string {
	token := uuid.NewString()
	createdAt := time.Now().UTC()
	testAtUnix := int(testAtTime.Unix())
	yData := yAnalyserModel.YAnalyser{
		AnalyserNum: &m.productName, CreatedAt: &createdAt, CrushingPlant: &m.productName, Flux: &flux,
		Load: &load, Speed: &beltSpeed, Status: &status, TestAt: &testAtUnix, Token: &token,
	}
	if val, ok := matGrades["CaO"]; ok {
		yData.Mat1 = &val
	}
	if val, ok := matGrades["MgO"]; ok {
		yData.Mat2 = &val
	}
	return m.generateSQLWithGORM(&yData)
}

// --- Path and Movement Logic ---
func (m *SimulateDataMaker) generatePathSQLs(startTime time.Time, startPoint, endPoint Point, waypoints []Point, entityID, entityName string, speedMPS float64, statusMoving int) ([]string, Point, time.Time) {
	var generatedSQLs []string
	currentPos := startPoint
	currentTime := startTime
	fullPath := []Point{startPoint}
	if len(waypoints) > 0 {
		fullPath = append(fullPath, waypoints...)
	}
	fullPath = append(fullPath, endPoint)
	if speedMPS <= 0.01 {
		speedMPS = 1.0
		m.logger.Warn("Speed too low, defaulted to 1.0 m/s", zap.String("entityID", entityID))
	}
	if startPoint.Lng == endPoint.Lng && startPoint.Lat == endPoint.Lat && len(waypoints) == 0 {
		sql := m.generateLocationGnssSQL(entityID, entityName, currentPos, currentTime, 0, 0, statusMoving, 0, nil)
		generatedSQLs = append(generatedSQLs, sql)
		m.allEvents.AddSQLEvent(currentTime, sql)
		return generatedSQLs, currentPos, currentTime
	}
	for i := 0; i < len(fullPath)-1; i++ {
		p1 := fullPath[i]
		p2 := fullPath[i+1]
		segmentDistance := Distance(p1, p2)
		if segmentDistance < 0.1 {
			continue
		}
		segmentDurationSeconds := segmentDistance / speedMPS
		heading := CalculateHeading(p1, p2)
		speedKmh := float32(speedMPS * 3.6)
		numSteps := int(math.Ceil(segmentDurationSeconds))
		if numSteps == 0 && segmentDurationSeconds > 0 {
			numSteps = 1
		}
		if numSteps == 0 {
			continue
		}
		timeStepDuration := time.Duration(segmentDurationSeconds * float64(time.Second) / float64(numSteps))
		if numSteps == 1 {
			timeStepDuration = time.Duration(segmentDurationSeconds * float64(time.Second))
		}
		for step := 0; step < numSteps; step++ {
			var pointToLog Point
			var eventTimeForLog time.Time
			if step == 0 {
				pointToLog = p1
				eventTimeForLog = currentTime
			} else {
				ratio := float64(step) / float64(numSteps)
				pointToLog = InterpolatePoint(p1, p2, ratio)
				currentTime = currentTime.Add(timeStepDuration)
				eventTimeForLog = currentTime
			}
			currentPos = pointToLog
			sql := m.generateLocationGnssSQL(entityID, entityName, currentPos, eventTimeForLog, heading, speedKmh, statusMoving, 0, nil)
			generatedSQLs = append(generatedSQLs, sql)
			m.allEvents.AddSQLEvent(eventTimeForLog, sql)
		}
	}
	totalPathDistance := 0.0
	for k := 0; k < len(fullPath)-1; k++ {
		totalPathDistance += Distance(fullPath[k], fullPath[k+1])
	}
	totalPathDuration := totalPathDistance / speedMPS
	finalArrivalTime := startTime.Add(time.Duration(totalPathDuration * float64(time.Second)))
	finalSpeed := float32(0)
	finalStatus := statusMoving
	finalHeading := float32(0)
	if len(fullPath) > 1 {
		finalHeading = CalculateHeading(fullPath[len(fullPath)-2], fullPath[len(fullPath)-1])
	}
	sql := m.generateLocationGnssSQL(entityID, entityName, endPoint, finalArrivalTime, finalHeading, finalSpeed, finalStatus, 0, nil)
	generatedSQLs = append(generatedSQLs, sql)
	m.allEvents.AddSQLEvent(finalArrivalTime, sql)
	return generatedSQLs, endPoint, finalArrivalTime
}

// --- Main Generation Logic (GenSimulateData) ---
func (m *SimulateDataMaker) GenSimulateData() error {
	m.logger.Info("Starting data generation process", zap.String("dataPath", m.dataPath), zap.String("productName", m.productName), zap.String("dateStr", m.dateStr))
	rand.Seed(time.Now().UnixNano())
	var err error
	m.pathConfig, err = LoadPathConfig(m.dataPath, m.productName, m.dateStr)
	if err != nil {
		m.logger.Warn("Failed to load Path.json, proceeding with no explicit path rules.", zap.Error(err))
		m.pathConfig = &PathConfig{Rules: []PathRule{}}
	}

	boomGroupIDs, _, boomGroupMap := m.getBoomGroupInfos()
	diggerIDs, _, diggerIDToName := m.getDiggerInfosOfProduct(m.productName)
	lorryIDs, _, lorryIDToName := m.getLorryInfos()
	if len(diggerIDs) == 0 || len(lorryIDs) == 0 || len(boomGroupIDs) == 0 {
		err := fmt.Errorf("insufficient entities: D:%d, L:%d, BG:%d", len(diggerIDs), len(lorryIDs), len(boomGroupIDs))
		m.logger.Error("Cannot start simulation", zap.Error(err))
		return err
	}

	parkingLotPos := Point{Lng: 120.0000e6, Lat: 30.0000e6, Elevation: 10}
	diggerSafeZonePos := Point{Lng: 120.0080e6, Lat: 30.0080e6, Elevation: 90}
	weighbridgePos := Point{Lng: 120.0040e6, Lat: 30.0040e6, Elevation: 20}
	unloadSiteLng, unloadSiteLat := m.getCenterOfUnloadSite(m.productName)
	unloadSitePos := Point{Lng: unloadSiteLng, Lat: unloadSiteLat, Elevation: 25}

	initTime := m.baseTime.Add(5*time.Hour + 50*time.Minute)
	var initSQLs []string
	initialStatusParked, initialSpeed, initialHeading := 0, float32(0), float32(0)
	for _, lorryID := range lorryIDs {
		initSQLs = append(initSQLs, m.generateLocationGnssSQL(lorryID, lorryIDToName[lorryID], parkingLotPos, initTime, initialHeading, initialSpeed, initialStatusParked, 0, nil))
	}
	for _, diggerID := range diggerIDs {
		initSQLs = append(initSQLs, m.generateLocationGnssSQL(diggerID, diggerIDToName[diggerID], diggerSafeZonePos, initTime, initialHeading, initialSpeed, initialStatusParked, 0, nil))
	}
	m.allEvents.AddSQLsEvent(initTime, initSQLs)
	m.updateCurrentFileEndTime(initTime)
	m.writeEventsToSQLFile("init", true)

	m.currentTime = m.baseTime.Add(6 * time.Hour)
	m.outputFileIndex = 1
	m.currentFileEndTime = m.currentTime
	simulationEndTime := m.baseTime.Add(18 * time.Hour)
	lunchStartTime := m.baseTime.Add(12 * time.Hour)
	lunchEndTime := m.baseTime.Add(13 * time.Hour)

	diggerAssignments := make(map[string]string)
	diggerToBGs := make(map[string][]string)
	if len(diggerIDs) >= 1 {
		diggerToBGs[diggerIDs[0]] = []string{}
		for i := 0; i < len(boomGroupIDs)/2 && i < len(boomGroupIDs); i++ {
			diggerToBGs[diggerIDs[0]] = append(diggerToBGs[diggerIDs[0]], boomGroupIDs[i])
		}
	}
	if len(diggerIDs) >= 2 {
		diggerToBGs[diggerIDs[1]] = []string{}
		for i := len(boomGroupIDs) / 2; i < len(boomGroupIDs); i++ {
			diggerToBGs[diggerIDs[1]] = append(diggerToBGs[diggerIDs[1]], boomGroupIDs[i])
		}
	}
	for diggerID, bgs := range diggerToBGs {
		if len(bgs) > 0 {
			diggerAssignments[diggerID] = bgs[0]
		} else if len(boomGroupIDs) > 0 {
			diggerAssignments[diggerID] = boomGroupIDs[0]
		}
	}

	lorryStates := make(map[string]*LorryState)
	for _, id := range lorryIDs {
		lorryStates[id] = &LorryState{ID: id, Name: lorryIDToName[id], CurrentPos: parkingLotPos, CurrentTime: m.currentTime, StatusString: "idle_at_parking"}
	}
	diggerStates := make(map[string]*DiggerState)
	for _, id := range diggerIDs {
		assignedBG := diggerAssignments[id]
		if assignedBG == "" && len(boomGroupIDs) > 0 {
			assignedBG = boomGroupIDs[0]
			m.logger.Warn("Digger had no assigned BG, falling back", zap.String("diggerID", id), zap.String("fallbackBG", assignedBG))
		}
		diggerStates[id] = &DiggerState{ID: id, Name: diggerIDToName[id], CurrentPos: diggerSafeZonePos, CurrentTime: m.currentTime, AssignedBoomGroupID: assignedBG, StatusString: "idle_at_safe_zone"}
	}

	for _, dState := range diggerStates {
		if dState.CurrentTime.After(simulationEndTime) || dState.AssignedBoomGroupID == "" {
			continue
		}
		boomGroup := boomGroupMap[dState.AssignedBoomGroupID]
		boomGroupTargetPos := Point{Lng: boomGroup.X_(), Lat: boomGroup.Y_(), Elevation: (boomGroup.High_() + boomGroup.Low_()) / 2}
		waypoints := m.pathConfig.GetRoute(LocationDiggerStop, boomGroup.Name_())
		_, dState.CurrentPos, dState.CurrentTime = m.generatePathSQLs(dState.CurrentTime, dState.CurrentPos, boomGroupTargetPos, waypoints, dState.ID, dState.Name, DiggerSpeed, 1)
		dState.StatusString = "at_boom_group_idle"
		m.logDiggerStatusEvent(dState, boomGroup, 2, 0, dState.CurrentTime)
		m.updateCurrentFileEndTime(dState.CurrentTime)
	}
	m.writeEventsToSQLFile("diggers_setup", false)

	maxTripsPerLorry := 2
	diggerServiceIndex := 0
	for lorryIdxLoop := 0; lorryIdxLoop < len(lorryIDs); lorryIdxLoop++ {
		lState := lorryStates[lorryIDs[lorryIdxLoop]]
		for lState.TripCount < maxTripsPerLorry && lState.CurrentTime.Before(simulationEndTime) {
			var targetDiggerState *DiggerState
			var assignedDiggerInfo DiggerState
			attempts := 0
			for attempts < len(diggerIDs) {
				currentDiggerID := diggerIDs[diggerServiceIndex%len(diggerIDs)]
				diggerServiceIndex++
				potentialDigger := diggerStates[currentDiggerID]
				if potentialDigger.AssignedBoomGroupID != "" && potentialDigger.StatusString == "at_boom_group_idle" {
					targetDiggerState = potentialDigger
					assignedDiggerInfo = *targetDiggerState
					break
				}
				attempts++
			}
			if targetDiggerState == nil {
				m.logger.Info("Lorry: No digger available, waiting", zap.String("lorry", lState.Name), zap.Time("time", lState.CurrentTime))
				lState.CurrentTime = lState.CurrentTime.Add(10 * time.Minute)
				m.logLorryStatusEvent(lState, 0, 0, lState.CurrentTime)
				m.updateCurrentFileEndTime(lState.CurrentTime)
				continue
			}
			lState.AssignedDiggerID = targetDiggerState.ID
			boomGroupID := targetDiggerState.AssignedBoomGroupID
			boomGroup := boomGroupMap[boomGroupID]
			lState.CurrentBoomGroupData = &boomGroup
			diggerPosAtBoom := Point{Lng: boomGroup.X_(), Lat: boomGroup.Y_(), Elevation: (boomGroup.High_() + boomGroup.Low_()) / 2}
			if lState.CurrentTime.After(lunchStartTime) && lState.CurrentTime.Before(lunchEndTime) && lState.StatusString != "lunch_break" {
				if Distance(lState.CurrentPos, unloadSitePos) < 10 {
					previousStatus := lState.StatusString
					lState.StatusString = "lunch_break"
					m.logLorryStatusEvent(lState, 10, 0, lState.CurrentTime)
					m.updateCurrentFileEndTime(lState.CurrentTime)
					lState.CurrentTime = lunchEndTime
					lState.StatusString = previousStatus
					if lState.StatusString == "lunch_break" {
						lState.StatusString = "idle_after_lunch"
					}
					m.logLorryStatusEvent(lState, 0, 0, lState.CurrentTime)
					m.updateCurrentFileEndTime(lState.CurrentTime) /* m.writeEventsToSQLFile(fmt.Sprintf("lorry_%s_lunch", lState.ID), false) */
				}
			}
			if lState.CurrentTime.After(simulationEndTime) {
				break
			}
			lState.TripCount++
			lState.CurrentBusiNo = uuid.NewString()
			m.logger.Info("Lorry starting trip", zap.String("lorry", lState.Name), zap.Int("trip", lState.TripCount), zap.Time("startTime", lState.CurrentTime), zap.String("digger", assignedDiggerInfo.Name), zap.String("boomGroup", boomGroup.Name_()), zap.String("busiNo", lState.CurrentBusiNo))
			currentSpeed := TruckSpeedEmpty
			startLocNameForPathCalc := LocationParking
			if lState.TripCount > 1 || lState.StatusString == "weighed_tare_returning" {
				startLocNameForPathCalc = LocationWeighbridge
			}
			lState.StatusString = "to_digger"
			waypointsToDigger := m.pathConfig.GetRoute(startLocNameForPathCalc, boomGroup.Name_())
			_, lState.CurrentPos, lState.CurrentTime = m.generatePathSQLs(lState.CurrentTime, lState.CurrentPos, diggerPosAtBoom, waypointsToDigger, lState.ID, lState.Name, currentSpeed, 1)
			m.logLorryStatusEvent(lState, 2, 0, lState.CurrentTime)
			m.updateCurrentFileEndTime(lState.CurrentTime)
			lState.StatusString = "loading"
			loadingStartTime := lState.CurrentTime
			targetDiggerState.StatusString = "digging"
			m.logLorryStatusEvent(lState, 3, 0, loadingStartTime)
			m.logDiggerStatusEvent(targetDiggerState, boomGroup, 3, 0, loadingStartTime)
			lState.CurrentTime = lState.CurrentTime.Add(LoadingTimeDuration)
			lState.IsLoaded = true
			targetDiggerState.CurrentTime = lState.CurrentTime
			targetDiggerState.StatusString = "at_boom_group_idle"
			m.logLorryStatusEvent(lState, 4, 0, lState.CurrentTime)
			m.logDiggerStatusEvent(targetDiggerState, boomGroup, 2, 0, lState.CurrentTime)
			m.updateCurrentFileEndTime(lState.CurrentTime)

			lState.StatusString = "to_weigh_gross"
			waypointsToWeighGross := m.pathConfig.GetRoute(boomGroup.Name_(), LocationWeighbridge)
			_, lState.CurrentPos, lState.CurrentTime = m.generatePathSQLs(lState.CurrentTime, lState.CurrentPos, weighbridgePos, waypointsToWeighGross, lState.ID, lState.Name, TruckSpeedLoaded, 1)
			m.logLorryStatusEvent(lState, 2, 0, lState.CurrentTime)
			m.updateCurrentFileEndTime(lState.CurrentTime)
			lState.StatusString = "weighing_gross"
			weighingGrossStartTime := lState.CurrentTime
			m.logLorryStatusEvent(lState, 5, 0, weighingGrossStartTime)
			lState.CurrentTime = lState.CurrentTime.Add(WeighingDuration)
			lState.LastTareWeightKg = 20000.0
			lState.LastNetWeightKg = 60000.0
			grossWKg := lState.LastTareWeightKg + lState.LastNetWeightKg
			dirGross := "1"
			matCode := "0102010003"
			matName := "石灰石碎石"
			weighLoggerSQLGross := m.generateWeighLoggerSQL(lState.Name, lState.CurrentBusiNo, weighingGrossStartTime, dirGross, grossWKg, lState.LastNetWeightKg, lState.LastTareWeightKg, matCode, matName, nil)
			m.allEvents.AddSQLEvent(weighingGrossStartTime, weighLoggerSQLGross)
			m.updateCurrentFileEndTime(weighingGrossStartTime)
			m.logLorryStatusEvent(lState, 4, 0, lState.CurrentTime)
			m.updateCurrentFileEndTime(lState.CurrentTime)

			lState.StatusString = "to_unload"
			waypointsToUnload := m.pathConfig.GetRoute(LocationWeighbridge, LocationUnloadSite)
			_, lState.CurrentPos, lState.CurrentTime = m.generatePathSQLs(lState.CurrentTime, lState.CurrentPos, unloadSitePos, waypointsToUnload, lState.ID, lState.Name, TruckSpeedLoaded, 1)
			m.logLorryStatusEvent(lState, 2, 0, lState.CurrentTime)
			m.updateCurrentFileEndTime(lState.CurrentTime)

			lState.StatusString = "unloading"
			unloadingOpStartTime := lState.CurrentTime
			m.logLorryStatusEvent(lState, 6, 0, unloadingOpStartTime)
			unloadingOpEndTime := unloadingOpStartTime.Add(UnloadingTimeDuration)
			currentNetWeightKg := lState.LastNetWeightKg
			if currentNetWeightKg <= 0 {
				m.logger.Warn("Lorry LastNetWeightKg is not positive for YAnalyser.", zap.String("lorry", lState.Name), zap.Float64("netWeight", currentNetWeightKg))
			}
			fluxKgPerSec := 0.0
			if UnloadingTimeDuration.Seconds() > 0 {
				fluxKgPerSec = currentNetWeightKg / UnloadingTimeDuration.Seconds()
			} else {
				m.logger.Warn("Lorry UnloadingTimeDuration is zero.", zap.String("lorry", lState.Name))
			}
			yAnalyserTimeCursor := unloadingOpStartTime
			m.logger.Debug("YANALYSER: Lorry Unloading cycle started.", zap.String("lorry", lState.Name), zap.Time("start", unloadingOpStartTime), zap.Time("end", unloadingOpEndTime), zap.Float64("netWeightKg", currentNetWeightKg), zap.Float64("fluxKgPerSec", fluxKgPerSec))
			yanalyserRecordsGenerated := 0
			for yAnalyserTimeCursor.Before(unloadingOpEndTime) {
				m.logger.Debug("YANALYSER: Loop iteration.", zap.Time("cursor", yAnalyserTimeCursor), zap.Time("targetEnd", unloadingOpEndTime))
				matGrades := make(map[string]float64)
				if lState.CurrentBoomGroupData != nil {
					bg_data := lState.CurrentBoomGroupData
					if bg_data.Material1 != nil {
						matGrades["CaO"] = m.fluctuateMaterial(*bg_data.Material1, m.productMaterials["CaO"])
					}
					if bg_data.Material2 != nil {
						matGrades["MgO"] = m.fluctuateMaterial(*bg_data.Material2, m.productMaterials["MgO"])
					}
				} else {
					m.logger.Warn("Lorry CurrentBoomGroupData nil for YAnalyser.", zap.String("lorry", lState.Name))
					matGrades["CaO"] = m.fluctuateMaterial(43.0, m.productMaterials["CaO"])
					matGrades["MgO"] = m.fluctuateMaterial(1.0, m.productMaterials["MgO"])
				}
				m.logger.Debug("YANALYSER: Calculated material grades for current record.", zap.Time("forTime", yAnalyserTimeCursor), zap.Float64("CaO", matGrades["CaO"]), zap.Float64("MgO", matGrades["MgO"]))
				yAnalyserSQL := m.generateYAnalyserSQL(yAnalyserTimeCursor, fluxKgPerSec, currentNetWeightKg, 1.5, 1, matGrades)
				m.logger.Debug("YANALYSER: Generated YAnalyser SQL.", zap.Time("forTime", yAnalyserTimeCursor), zap.String("sqlPrefix", firstNChars(yAnalyserSQL, 120)))
				if strings.HasPrefix(yAnalyserSQL, "-- ERROR") || yAnalyserSQL == "" {
					m.logger.Error("YANALYSER: Failed to generate valid YAnalyser SQL string.", zap.Time("forTime", yAnalyserTimeCursor), zap.String("returnedSQL", yAnalyserSQL))
				} else {
					m.allEvents.AddSQLEvent(yAnalyserTimeCursor, yAnalyserSQL)
					yanalyserRecordsGenerated++
					m.logger.Debug("YANALYSER: Successfully added YAnalyser SQL to events.", zap.Time("forTime", yAnalyserTimeCursor))
				}
				m.updateCurrentFileEndTime(yAnalyserTimeCursor)
				yAnalyserTimeCursor = yAnalyserTimeCursor.Add(30 * time.Second)
			}
			m.logger.Debug("YANALYSER: Lorry - Total YAnalyser records generated in this unload cycle.", zap.String("lorry", lState.Name), zap.Int("count", yanalyserRecordsGenerated))
			lState.CurrentTime = unloadingOpEndTime
			lState.IsLoaded = false
			m.logLorryStatusEvent(lState, 0, 0, lState.CurrentTime)
			m.updateCurrentFileEndTime(lState.CurrentTime)

			lState.StatusString = "to_weigh_tare"
			waypointsToWeighTare := m.pathConfig.GetRoute(LocationUnloadSite, LocationWeighbridge)
			_, lState.CurrentPos, lState.CurrentTime = m.generatePathSQLs(lState.CurrentTime, lState.CurrentPos, weighbridgePos, waypointsToWeighTare, lState.ID, lState.Name, TruckSpeedEmpty, 1)
			m.logLorryStatusEvent(lState, 2, 0, lState.CurrentTime)
			m.updateCurrentFileEndTime(lState.CurrentTime)
			lState.StatusString = "weighing_tare"
			weighingTareStartTime := lState.CurrentTime
			m.logLorryStatusEvent(lState, 5, 0, weighingTareStartTime)
			lState.CurrentTime = lState.CurrentTime.Add(WeighingDuration)
			currentTareWeightKg := lState.LastTareWeightKg
			netWEmpty := 0.0
			dirTare := "0"
			tareMatCode := matCode
			tareMatName := matName
			weighLoggerSQLTare := m.generateWeighLoggerSQL(lState.Name, lState.CurrentBusiNo, weighingTareStartTime, dirTare, currentTareWeightKg, netWEmpty, currentTareWeightKg, tareMatCode, tareMatName, nil)
			m.allEvents.AddSQLEvent(weighingTareStartTime, weighLoggerSQLTare)
			m.updateCurrentFileEndTime(weighingTareStartTime)
			m.logLorryStatusEvent(lState, 0, 0, lState.CurrentTime)
			m.updateCurrentFileEndTime(lState.CurrentTime)
			lState.StatusString = "weighed_tare_returning"
			m.writeEventsToSQLFile(fmt.Sprintf("lorry_%s_trip%d", lState.ID, lState.TripCount), false)
			if lState.CurrentTime.After(simulationEndTime) {
				break
			}
		}
		if lState.StatusString != "returned_to_parking" {
			startLocForParkingPath := LocationWeighbridge
			if lState.StatusString == "unloading" || (Distance(lState.CurrentPos, unloadSitePos) < 10) {
				startLocForParkingPath = LocationUnloadSite
			}
			waypointsToParking := m.pathConfig.GetRoute(startLocForParkingPath, LocationParking)
			_, lState.CurrentPos, lState.CurrentTime = m.generatePathSQLs(lState.CurrentTime, lState.CurrentPos, parkingLotPos, waypointsToParking, lState.ID, lState.Name, TruckSpeedEmpty, 1)
			m.logLorryStatusEvent(lState, 0, 0, lState.CurrentTime)
			lState.StatusString = "returned_to_parking"
			m.updateCurrentFileEndTime(lState.CurrentTime)
			m.writeEventsToSQLFile(fmt.Sprintf("lorry_%s_to_parking", lState.ID), false)
		}
	}

	m.outputFileIndex = 100
	uninitEventTime := simulationEndTime.Add(1 * time.Hour)
	if m.currentFileEndTime.After(uninitEventTime) {
		uninitEventTime = m.currentFileEndTime.Add(10 * time.Minute)
	}
	m.currentFileEndTime = uninitEventTime
	var uninitSQLs []string
	offlineStatus, finalSpeed, finalHeading := 99, float32(0), float32(0)
	for _, lState := range lorryStates {
		uninitSQLs = append(uninitSQLs, m.generateLocationGnssSQL(lState.ID, lState.Name, lState.CurrentPos, uninitEventTime, finalHeading, finalSpeed, offlineStatus, 0, nil))
	}
	for _, dState := range diggerStates {
		uninitSQLs = append(uninitSQLs, m.generateLocationGnssSQL(dState.ID, dState.Name, dState.CurrentPos, uninitEventTime, finalHeading, finalSpeed, offlineStatus, 0, nil))
	}
	m.allEvents.AddSQLsEvent(uninitEventTime, uninitSQLs)
	m.writeEventsToSQLFile("uninit", false)
	m.logger.Info("SimulateDataMaker: Data generation process complete.")
	return nil
}

// --- Logging and Helper Wrappers ---
func (m *SimulateDataMaker) logLorryStatusEvent(lState *LorryState, status int, speedKmh float32, eventTime time.Time) {
	sql := m.generateLocationGnssSQL(lState.ID, lState.Name, lState.CurrentPos, eventTime, 0, speedKmh, status, 0, nil)
	m.allEvents.AddSQLEvent(eventTime, sql)
}
func (m *SimulateDataMaker) logDiggerStatusEvent(dState *DiggerState, bg boomGroupModel.BoomGroup, status int, speedKmh float32, eventTime time.Time) {
	diggerPosAtBoom := Point{Lng: bg.X_(), Lat: bg.Y_(), Elevation: (bg.High_() + bg.Low_()) / 2}
	sql := m.generateLocationGnssSQL(dState.ID, dState.Name, diggerPosAtBoom, eventTime, 0, speedKmh, status, 0, nil)
	m.allEvents.AddSQLEvent(eventTime, sql)
}

// fluctuateMaterial
func (m *SimulateDataMaker) fluctuateMaterial(baseValue float64, pmi ProductMaterialInfo) float64 {
	var fr float64
	if pmi.Max-pmi.Min > 0 {
		fr = (pmi.Max - pmi.Min) * MaterialFluctuationPercentage
	} else {
		fr = baseValue * MaterialFluctuationPercentage
	}
	if math.IsNaN(fr) || math.IsInf(fr, 0) {
		fr = 0.01
	}
	fl := (rand.Float64()*2 - 1) * fr
	res := baseValue + fl
	if res < 0 && baseValue >= 0 {
		res = rand.Float64() * 0.01 * baseValue
	}
	return res
}

// updateCurrentFileEndTime
func (m *SimulateDataMaker) updateCurrentFileEndTime(eventTime time.Time) {
	if eventTime.After(m.currentFileEndTime) {
		m.currentFileEndTime = eventTime
	}
}

// writeEventsToSQLFile
func (m *SimulateDataMaker) writeEventsToSQLFile(nameHint string, isInitFile bool) {
	m.logger.Debug("Attempting to write SQL file", zap.String("hint", nameHint), zap.Bool("isInit", isInitFile), zap.Int("eventCount", len(m.allEvents.Events)))
	for idx, event := range m.allEvents.Events {
		var firstSQLPrefix string
		if len(event.SQLs) > 0 {
			if len(event.SQLs[0]) > 120 {
				firstSQLPrefix = event.SQLs[0][:120] + "..."
			} else {
				firstSQLPrefix = event.SQLs[0]
			}
		} else {
			firstSQLPrefix = "[No SQLs in this event]"
		}
		m.logger.Debug("Event for write", zap.Int("index", idx), zap.Time("timestamp", event.Timestamp), zap.Int("sqlCount", len(event.SQLs)), zap.String("firstSQLPrefix", firstSQLPrefix))
	}
	if len(m.allEvents.Events) == 0 {
		m.logger.Debug("No events to write for file hint.", zap.String("hint", nameHint))
		return
	}
	sort.Slice(m.allEvents.Events, func(i, j int) bool { return m.allEvents.Events[i].Timestamp.Before(m.allEvents.Events[j].Timestamp) })
	var fi int
	if isInitFile {
		fi = 0
		nameHint = "init"
	} else if m.outputFileIndex == 100 {
		fi = 100
		nameHint = "uninit"
	} else {
		fi = m.outputFileIndex
	}
	fn := fmt.Sprintf("%d_%s.sql", fi, nameHint)
	fp := filepath.Join(m.dataPath, m.productName, m.dateStr, fn)
	if err := os.MkdirAll(filepath.Dir(fp), os.ModePerm); err != nil {
		m.logger.Error("Error creating directory", zap.String("path", filepath.Dir(fp)), zap.Error(err))
		return
	}
	var c bytes.Buffer
	lts := ""
	for _, ev := range m.allEvents.Events {
		ts := ev.Timestamp.Format(TimeFormat)
		if ts != lts {
			if c.Len() > 0 {
				c.WriteString("\n")
			}
			c.WriteString(fmt.Sprintf(SqlCommentFormat, ts) + "\n")
			lts = ts
		}
		for _, s := range ev.SQLs {
			c.WriteString(s + "\n")
		}
	}
	if err := os.WriteFile(fp, c.Bytes(), 0644); err != nil {
		m.logger.Error("Error writing SQL file", zap.String("path", fp), zap.Error(err))
		return
	}
	m.logger.Info("Successfully wrote SQL file", zap.String("path", fp), zap.Int("eventEntries", len(m.allEvents.Events)))
	m.allEvents = NewEventCollection()
	m.logger.Debug("Cleared m.allEvents after writing file.", zap.String("filename", fn), zap.Int("newEventCount", len(m.allEvents.Events)))
	if !isInitFile && m.outputFileIndex < 99 {
		m.outputFileIndex++
	}
}

// --- Pointer literal helpers ---
func intPtr(i int) *int             { return &i }
func float64Ptr(f float64) *float64 { return &f }
func strPtr(s string) *string       { return &s }

// Helper for logging SQL prefix (if not in common.go)
// func firstNChars(s string, n int) string { if len(s) > n { return s[:n] }; return s } // Assuming this is in common.go

// genSimulateData entry point function
func GenSimulateData(dataPath, productName, dateStr string) {
	maker, err := NewSimulateDataMaker(dataPath, productName, dateStr)
	if err != nil {
		utils.Logger.Fatal("Create Maker failed", zap.Error(err))
	}
	if err = maker.GenSimulateData(); err != nil {
		utils.Logger.Fatal("Gen Sim Data failed", zap.Error(err))
	}
}

// --- ZapGormLoggerAdapter (Corrected for GORM v2 logger.Interface) ---
type ZapGormLoggerAdapter struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

func (l *ZapGormLoggerAdapter) Printf(s string, i ...interface{}) {

}

func NewZapGormLoggerAdapter(zapLogger *zap.Logger) *ZapGormLoggerAdapter {
	return &ZapGormLoggerAdapter{
		ZapLogger:                 zapLogger,
		LogLevel:                  gormlogger.Warn, // Default GORM log level
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}
func (l *ZapGormLoggerAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}
func (l *ZapGormLoggerAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.ZapLogger.Sugar().Infow(msg, data...)
	}
}
func (l *ZapGormLoggerAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.ZapLogger.Sugar().Warnw(msg, data...)
	}
}
func (l *ZapGormLoggerAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.ZapLogger.Sugar().Errorw(msg, data...)
	}
}
func (l *ZapGormLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{zap.Duration("elapsed", elapsed), zap.String("sql", sql), zap.Int64("rows", rows)}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.IgnoreRecordNotFoundError) { // Use errors.Is for gorm.ErrRecordNotFound
		l.ZapLogger.Error("GORM Trace Error", append(fields, zap.Error(err))...)
	} else if elapsed > l.SlowThreshold && l.SlowThreshold != 0 {
		l.ZapLogger.Warn("GORM Slow SQL", fields...)
	} else if l.LogLevel >= gormlogger.Info { // Changed from Debug to Info for GORM Trace to align with common practice
		l.ZapLogger.Info("GORM SQL Trace", fields...) // GORM trace often logged at Info or Debug
	}
}
