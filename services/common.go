// filePath: sim_data_gen/services/common.go
package services

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap" // For zap.String, zap.Time etc.
	"math"
	"os"
	"path/filepath"
	"sim_data_gen/utils" // Import your utils package
	"sort"
	"strconv"
	"strings"
	"time"
	// "log" // No longer needed if all logging goes through zap
)

// Constants ... (TimeFormat, DateTimeFormat, SqlCommentFormat, speeds, durations, locations etc.)
const (
	TimeFormat                    = "15:04:05"
	DateTimeFormat                = "2006-01-02 15:04:05"
	SqlCommentFormat              = "--- %s"
	TruckSpeedEmptyLoadedRatio    = 1.5
	TruckSpeedLoaded              = 5.0
	TruckSpeedEmpty               = TruckSpeedLoaded * TruckSpeedEmptyLoadedRatio
	DiggerSpeed                   = 2.0
	LoadingTimeDuration           = 5 * time.Minute
	UnloadingTimeDuration         = 3 * time.Minute
	WeighingDuration              = 30 * time.Second
	LocationParking               = "停车场"
	LocationDiggerStop            = "挖机停靠地"
	LocationWeighbridge           = "地磅"
	LocationUnloadSite            = "破碎站"
	MaterialFluctuationPercentage = 0.03
)

// Point, PathRule, PathConfig ...
type Point struct {
	Lng       float64 `json:"lng"`
	Lat       float64 `json:"lat"`
	Elevation float64 `json:"elevation"`
}
type PathRule struct {
	Src    string  `json:"Src"`
	Dst    string  `json:"Dst"`
	Points []Point `json:"Point"`
}
type PathConfig struct{ Rules []PathRule }

func LoadPathConfig(dataPath, productName, dateStr string) (*PathConfig, error) {
	pathFilePath := filepath.Join(dataPath, productName, dateStr, "Path.json")
	bytes, err := os.ReadFile(pathFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			utils.Logger.Warn("Path.json not found, proceeding with no explicit path rules.", zap.String("path", pathFilePath))
			return &PathConfig{Rules: []PathRule{}}, nil
		}
		utils.Logger.Error("Failed to read Path.json", zap.String("path", pathFilePath), zap.Error(err))
		return nil, fmt.Errorf("failed to read Path.json from %s: %w", pathFilePath, err)
	}
	var rules []PathRule
	if err := json.Unmarshal(bytes, &rules); err != nil {
		utils.Logger.Error("Failed to unmarshal Path.json", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal Path.json: %w", err)
	}
	return &PathConfig{Rules: rules}, nil
}
func (pc *PathConfig) GetRoute(srcName, dstName string) []Point { /* ... same as before ... */
	for _, rule := range pc.Rules {
		if rule.Src == srcName && (strings.HasPrefix(dstName, rule.Dst) || rule.Dst == dstName) {
			return rule.Points
		}
	}
	return nil
}

// Event and EventCollection
type Event struct {
	Timestamp time.Time
	SQLs      []string
}
type EventCollection struct{ Events []Event }

func NewEventCollection() *EventCollection { return &EventCollection{Events: []Event{}} }

func (ec *EventCollection) AddSQLEvent(eventTime time.Time, sql string) {
	utils.Logger.Debug("AddSQLEvent called", zap.Time("eventTime", eventTime), zap.String("sqlPrefix", firstNChars(sql, 80)))
	if sql == "" {
		utils.Logger.Debug("Skipped adding empty SQL", zap.Time("eventTime", eventTime))
		return
	}
	for i, e := range ec.Events {
		if e.Timestamp.Equal(eventTime) {
			ec.Events[i].SQLs = append(ec.Events[i].SQLs, sql)
			return
		}
	}
	ec.Events = append(ec.Events, Event{Timestamp: eventTime, SQLs: []string{sql}})
	sort.Slice(ec.Events, func(i, j int) bool { return ec.Events[i].Timestamp.Before(ec.Events[j].Timestamp) })
}
func (ec *EventCollection) AddSQLsEvent(eventTime time.Time, sqls []string) {
	if len(sqls) == 0 {
		return
	}
	var validSqls []string
	for _, s := range sqls {
		if s != "" {
			validSqls = append(validSqls, s)
		}
	}
	if len(validSqls) == 0 {
		return
	}
	for i, e := range ec.Events {
		if e.Timestamp.Equal(eventTime) {
			ec.Events[i].SQLs = append(ec.Events[i].SQLs, validSqls...)
			return
		}
	}
	ec.Events = append(ec.Events, Event{Timestamp: eventTime, SQLs: validSqls})
	sort.Slice(ec.Events, func(i, j int) bool { return ec.Events[i].Timestamp.Before(ec.Events[j].Timestamp) })
}

// Utility Functions (Distance, InterpolatePoint, CalculateHeading, ParseTimeAnnotation)
func Distance(p1, p2 Point) float64 { /* ... same as before ... */
	const R = 6371000
	lat1Rad := p1.Lat / 1e6 * math.Pi / 180
	lon1Rad := p1.Lng / 1e6 * math.Pi / 180
	lat2Rad := p2.Lat / 1e6 * math.Pi / 180
	lon2Rad := p2.Lng / 1e6 * math.Pi / 180
	dlon := lon2Rad - lon1Rad
	dlat := lat2Rad - lat1Rad
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Asin(math.Sqrt(a))
	horizontalDistance := R * c
	verticalDistance := math.Abs(p1.Elevation - p2.Elevation)
	return math.Sqrt(math.Pow(horizontalDistance, 2) + math.Pow(verticalDistance, 2))
}
func InterpolatePoint(p1, p2 Point, ratio float64) Point { /* ... same as before ... */
	if ratio <= 0 {
		return p1
	}
	if ratio >= 1 {
		return p2
	}
	return Point{Lng: p1.Lng + (p2.Lng-p1.Lng)*ratio, Lat: p1.Lat + (p2.Lat-p1.Lat)*ratio, Elevation: p1.Elevation + (p2.Elevation-p1.Elevation)*ratio}
}
func CalculateHeading(p1, p2 Point) float32 { /* ... same as before ... */
	if p1.Lat == p2.Lat && p1.Lng == p2.Lng {
		return 0
	}
	lat1Rad := p1.Lat / 1e6 * math.Pi / 180
	lon1Rad := p1.Lng / 1e6 * math.Pi / 180
	lat2Rad := p2.Lat / 1e6 * math.Pi / 180
	lon2Rad := p2.Lng / 1e6 * math.Pi / 180
	dLon := lon2Rad - lon1Rad
	y := math.Sin(dLon) * math.Cos(lat2Rad)
	x := math.Cos(lat1Rad)*math.Sin(lat2Rad) - math.Sin(lat1Rad)*math.Cos(lat2Rad)*math.Cos(dLon)
	bearingRad := math.Atan2(y, x)
	bearingDegrees := math.Mod(bearingRad*180/math.Pi+360, 360)
	return float32(bearingDegrees)
}
func ParseTimeAnnotation(baseDate time.Time, timeStr string) (time.Time, error) { /* ... same as before ... */
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid time format: '%s', expected HH:MM:SS", timeStr)
	}
	hour, errH := strconv.Atoi(parts[0])
	min, errM := strconv.Atoi(parts[1])
	sec, errS := strconv.Atoi(parts[2])
	if errH != nil || errM != nil || errS != nil || hour < 0 || hour > 23 || min < 0 || min > 59 || sec < 0 || sec > 59 {
		return time.Time{}, fmt.Errorf("invalid time value in '%s'", timeStr)
	}
	return time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), hour, min, sec, 0, baseDate.Location()), nil
}

// Helper for logging SQL prefix
func firstNChars(s string, n int) string {
	if len(s) == 0 {
		return "[EMPTY SQL]"
	}
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}
