package services

const (
	TimeFormatSQLComment = "15:04:05"
	//DateTimeFormat       = "2006-01-02 15:04:05"
	DateFormat    = "2006-01-02"
	EarthRadiusKm = 6371.0
)

//// Point represents a geographic coordinate
//type Point struct {
//	Lng       float64 `json:"lng"`       // Millionth of a degree
//	Lat       float64 `json:"lat"`       // Millionth of a degree
//	Elevation float64 `json:"elevation"` // Meters
//}
//
//// PathRule defines a rule for waypoints between a source and destination
//type PathRule struct {
//	Src    string  `json:"Src"`
//	Dst    string  `json:"Dst"`
//	Points []Point `json:"Point"`
//}

// var pathRules []PathRule
//
// // LoadPathRules loads path rules from Path.json
//
//	func LoadPathRules(dataPath string) error {
//		filePath := filepath.Join(dataPath, "..", "Path.json") // Assuming Path.json is one level above dataPath (e.g. project root)
//		if _, err := os.Stat(filePath); os.IsNotExist(err) {
//			fmt.Printf("Path.json not found at %s, proceeding without path rules.\n", filePath)
//			pathRules = []PathRule{} // Ensure it's initialized
//			return nil
//		}
//
//		bytes, err := ioutil.ReadFile(filePath)
//		if err != nil {
//			return fmt.Errorf("failed to read Path.json: %w", err)
//		}
//		err = json.Unmarshal(bytes, &pathRules)
//		if err != nil {
//			return fmt.Errorf("failed to unmarshal Path.json: %w", err)
//		}
//		fmt.Printf("Loaded %d path rules from %s\n", len(pathRules), filePath)
//		return nil
//	}
//
// // GetWaypoints returns waypoints for a given src and dst
//
//	func GetWaypoints(srcName, dstName string) []Point {
//		for _, rule := range pathRules {
//			if rule.Src == srcName && rule.Dst == dstName {
//				return rule.Points
//			}
//		}
//		return nil
//	}
//
// // Helper to create a pointer to a value
//
//	func Ptr[T any](v T) *T {
//		return &v
//	}
//
// Event represents a SQL statement to be executed at a specific time
//type Event struct {
//	Timestamp time.Time
//	SQLs      []string
//}
//
//// EventBlock groups SQLs by their timestamp comment
//type EventBlock struct {
//	Timestamp time.Time
//	SQLs      []string
//}

//// ByTimestamp implements sort.Interface for []Event
//type ByTimestamp []Event
//
//func (a ByTimestamp) Len() int           { return len(a) }
//func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp.Before(a[j].Timestamp) }
//
//// ByEventBlockTimestamp implements sort.Interface for []EventBlock
//type ByEventBlockTimestamp []EventBlock
//
//func (a ByEventBlockTimestamp) Len() int           { return len(a) }
//func (a ByEventBlockTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
//func (a ByEventBlockTimestamp) Less(i, j int) bool { return a[i].Timestamp.Before(a[j].Timestamp) }
//
//// CreateSQLFile writes events to a .sql file
//func CreateSQLFile(filePath string, events []Event) error {
//	if len(events) == 0 {
//		// Create an empty file if no events
//		f, err := os.Create(filePath)
//		if err != nil {
//			return fmt.Errorf("failed to create empty file %s: %w", filePath, err)
//		}
//		f.Close()
//		return nil
//	}
//
//	sort.Sort(ByTimestamp(events))
//
//	var content strings.Builder
//	lastCommentTime := time.Time{}
//
//	for _, event := range events {
//		eventTimeOnly := event.Timestamp.Format(TimeFormatSQLComment)
//		lastCommentTimeOnly := lastCommentTime.Format(TimeFormatSQLComment)
//
//		// New time block if the HH:MM:SS part is different
//		if lastCommentTime.IsZero() || eventTimeOnly != lastCommentTimeOnly {
//			if !lastCommentTime.IsZero() {
//				content.WriteString("\n") // Add a newline before the next time block
//			}
//			content.WriteString(fmt.Sprintf("--- %s\n", eventTimeOnly))
//			lastCommentTime = event.Timestamp
//		}
//		for _, sql := range event.SQLs {
//			content.WriteString(sql)
//			content.WriteString(";\n")
//		}
//	}
//
//	return ioutil.WriteFile(filePath, []byte(content.String()), 0644)
//}
//
//// GenerateLocationInsertSQL creates an INSERT SQL string for LocationGnssData
//func GenerateLocationInsertSQL(loc locationGnssDataModel.LocationGnssData) string {
//	// Handle nil pointers for SQL formatting
//	alarmStr := "NULL"
//	if loc.Alarm != nil {
//		alarmStr = fmt.Sprintf("%d", *loc.Alarm)
//	}
//	headingStr := "NULL"
//	if loc.Heading != nil {
//		headingStr = fmt.Sprintf("%f", *loc.Heading)
//	}
//	speedStr := "NULL"
//	if loc.Speed != nil {
//		speedStr = fmt.Sprintf("%f", *loc.Speed)
//	}
//	statusStr := "NULL"
//	if loc.Status != nil {
//		statusStr = fmt.Sprintf("%d", *loc.Status)
//	}
//	tokenStr := "NULL"
//	if loc.Token != nil {
//		tokenStr = fmt.Sprintf("'%s'", *loc.Token)
//	}
//	xStr := "NULL"
//	if loc.X != nil {
//		xStr = fmt.Sprintf("%f", *loc.X)
//	}
//	yStr := "NULL"
//	if loc.Y != nil {
//		yStr = fmt.Sprintf("%f", *loc.Y)
//	}
//	lastCommUtcStr := "NULL"
//	if loc.LastCommUtc != nil {
//		lastCommUtcStr = fmt.Sprintf("%d", *loc.LastCommUtc)
//	}
//	ntStr := "NULL"
//	if loc.Nt != nil {
//		ntStr = fmt.Sprintf("'%s'", *loc.Nt)
//	}
//
//	return fmt.Sprintf("INSERT INTO location_gnss_data (utc, tid, name, x, y, speed, heading, alarm, status, token, last_comm_utc, nt) VALUES (%d,'%s', '%s', %s, %s, %s, %s, %s, %s, %s, %s, %s);",
//		*loc.Utc, *loc.Tid, loc.Name, xStr, yStr, speedStr, headingStr, alarmStr, statusStr, tokenStr, lastCommUtcStr, ntStr)
//}
//
//// Dummy CreateObj function as per requirement "already implemented"
//// In a real scenario, this would interact with a database.
//func CreateObj(obj *locationGnssDataModel.LocationGnssData, tx ...*gorm.DB) (int64, error) {
//	// fmt.Printf("Simulating CreateObj for Tid: %s, Utc: %d\n", *obj.Tid, *obj.Utc)
//	// In a real system, this would return the number of affected rows and an error.
//	return 1, nil
//}
//
//// calculateDistance calculates the approximate distance between two points (millionths of a degree)
//// This is a simplified Euclidean distance, not Haversine. For small areas, it's an approximation.
//// Output is in "units" of millionths of degrees.
//func calculateDistance(p1, p2 Point) float64 {
//	dx := p1.Lng - p2.Lng
//	dy := p1.Lat - p2.Lat
//	return math.Sqrt(dx*dx + dy*dy)
//}
//
//// estimateTravelTime estimates travel time in seconds
//// Speed is in km/h. Distance is in "millionths of degrees". This needs a conversion factor.
//// For simplicity, let's assume 100,000 millionths of a degree is roughly 1km (very rough, depends on latitude)
//// This function needs significant refinement for real-world accuracy.
//func estimateTravelTime(p1, p2 Point, speedKmh float64) time.Duration {
//	if speedKmh == 0 {
//		return time.Duration(math.MaxInt64) // effectively infinite
//	}
//	distMillionths := calculateDistance(p1, p2)
//	distKm := distMillionths / 100000.0 // Very rough conversion
//	timeHours := distKm / speedKmh
//	return time.Duration(timeHours * float64(time.Hour))
//}
//
//// Interpolate path if segment time > 1s
//func interpolatePath(startPoint, endPoint Point, startTime time.Time, vehicleTID, vehicleName string, speedKmh float64, events *[]Event) time.Time {
//	segmentTravelTime := estimateTravelTime(startPoint, endPoint, speedKmh)
//	currentTime := startTime
//
//	if segmentTravelTime.Seconds() <= 1.0 {
//		currentTime = startTime.Add(segmentTravelTime)
//		locData := locationGnssDataModel.LocationGnssData{
//			Tid:         Ptr(vehicleTID),
//			Name:        vehicleName,
//			Utc:         Ptr(int(currentTime.Unix())),
//			X:           Ptr(endPoint.Lng),
//			Y:           Ptr(endPoint.Lat),
//			Speed:       Ptr(float32(speedKmh / 3.6)), // kmh to m/s
//			Heading:     Ptr(calculateHeading(startPoint, endPoint)),
//			LastCommUtc: Ptr(int(currentTime.Unix())),
//			Token:       Ptr(uuid.NewString()),
//		}
//		*events = append(*events, Event{Timestamp: currentTime, SQLs: []string{GenerateLocationInsertSQL(locData)}})
//		// CreateObj(&locData) // Call if needed per requirements
//		return currentTime
//	}
//
//	numSteps := int(math.Floor(segmentTravelTime.Seconds()))
//	if numSteps == 0 {
//		numSteps = 1
//	} // ensure at least one step to reach endpoint
//
//	dx := (endPoint.Lng - startPoint.Lng) / float64(numSteps)
//	dy := (endPoint.Lat - startPoint.Lat) / float64(numSteps)
//	dt := segmentTravelTime / time.Duration(numSteps)
//
//	currentPoint := startPoint
//	for i := 0; i < numSteps; i++ {
//		currentTime = startTime.Add(time.Duration(i+1) * dt)
//		currentPoint.Lng = startPoint.Lng + dx*float64(i+1)
//		currentPoint.Lat = startPoint.Lat + dy*float64(i+1)
//
//		locData := locationGnssDataModel.LocationGnssData{
//			Tid:         Ptr(vehicleTID),
//			Name:        vehicleName,
//			Utc:         Ptr(int(currentTime.Unix())),
//			X:           Ptr(currentPoint.Lng),
//			Y:           Ptr(currentPoint.Lat),
//			Speed:       Ptr(float32(speedKmh / 3.6)),
//			Heading:     Ptr(calculateHeading(startPoint, endPoint)), // Could refine heading per segment
//			LastCommUtc: Ptr(int(currentTime.Unix())),
//			Token:       Ptr(uuid.NewString()),
//		}
//		*events = append(*events, Event{Timestamp: currentTime, SQLs: []string{GenerateLocationInsertSQL(locData)}})
//		// CreateObj(&locData)
//	}
//	return currentTime // Time of arrival at endPoint
//}
//
//func calculateHeading(p1, p2 Point) float32 {
//	// Calculate heading from p1 to p2 in degrees (0-360, 0 = North)
//	// This is a simplified calculation for Cartesian coordinates
//	deltaX := p2.Lng - p1.Lng
//	deltaY := p2.Lat - p1.Lat
//	angleRad := math.Atan2(deltaX, deltaY) // Note: Atan2(x,y) for navigation where y is North
//	angleDeg := angleRad * (180.0 / math.Pi)
//	if angleDeg < 0 {
//		angleDeg += 360
//	}
//	return float32(angleDeg)
//}
//
//// NeutronSensorLog represents data from a neutron instrument
//type NeutronSensorLog struct {
//	gorm.Model
//	Timestamp      time.Time `gorm:"index"`
//	ProductName    string
//	LorryTID       string
//	FlowRate       float64            // e.g., tons per hour
//	MaterialGrades map[string]float64 `gorm:"-"` // For easy use, not directly stored
//	Material1      *float64
//	Material2      *float64
//	Material3      *float64
//	// ... Add up to Material20 as in BoomGroup
//}
//
//func (NeutronSensorLog) TableName() string {
//	return "neutron_sensor_logs"
//}
//
//// GenerateNeutronInsertSQL creates an INSERT SQL string for NeutronSensorLog
//func GenerateNeutronInsertSQL(log NeutronSensorLog) string {
//	// This is a simplified example, real one would map MaterialGrades to individual columns
//	// For now, let's just log a few materials.
//	mat1Str, mat2Str, mat3Str := "NULL", "NULL", "NULL"
//	if log.Material1 != nil {
//		mat1Str = fmt.Sprintf("%f", *log.Material1)
//	}
//	if log.Material2 != nil {
//		mat2Str = fmt.Sprintf("%f", *log.Material2)
//	}
//	if log.Material3 != nil {
//		mat3Str = fmt.Sprintf("%f", *log.Material3)
//	}
//
//	// Using Unix timestamp for DB compatibility in a simple int column
//	// If DB column is proper datetime, use log.Timestamp.Format("2006-01-02 15:04:05")
//	return fmt.Sprintf("INSERT INTO neutron_sensor_logs (created_at, updated_at, deleted_at, timestamp, product_name, lorry_tid, flow_rate, material1, material2, material3) VALUES ('%s', '%s', NULL, '%s', '%s', '%s', %f, %s, %s, %s);",
//		log.Timestamp.Format(DateTimeFormat), // created_at
//		log.Timestamp.Format(DateTimeFormat), // updated_at
//		log.Timestamp.Format(DateTimeFormat), // actual timestamp
//		log.ProductName,
//		log.LorryTID,
//		log.FlowRate,
//		mat1Str, mat2Str, mat3Str,
//		// ... other material grades
//	)
//}
//
//// EnsureDir ensures a directory exists
//func EnsureDir(dirPath string) error {
//	return os.MkdirAll(dirPath, os.ModePerm)
//}
//
//// Placeholder definitions for entity names for Path.json
//const (
//	EntityParkingLot     = "PARKING_LOT"
//	EntityWeighbridge    = "WEIGHBRIDGE"
//	EntityDiggerParkBase = "DIGGER_PARK" // e.g. DIGGER_PARK_D001
//	EntityBoomGroupBase  = "BOOM_GROUP"  // e.g. BOOM_GROUP_BG001
//	EntityUnloadSiteBase = "UNLOAD_SITE" // e.g. UNLOAD_SITE_CP1
//)
//
//// Placeholder coordinates (millionths of a degree)
//var (
//	CoordParkingLot     = Point{Lng: 120123456, Lat: 30123456, Elevation: 100}
//	CoordWeighbridge    = Point{Lng: 120223456, Lat: 30223456, Elevation: 110}
//	CoordDiggerParkBase = Point{Lng: 120323456, Lat: 30323456, Elevation: 150} // Base, adjust per digger
//)
