package entity

import (
	. "sim_data_gen/models/diggerSwitchBoomGroupLog"
	. "sim_data_gen/utils"
	"go.uber.org/zap"
)

type DiggerLogger struct {
	DiggerId          string
	DiggerName        string
	boomGroupSwitches *TimeSeries
	MaxUtc            int
}

func NewDiggerLogger(diggerId, diggerName string, boomGroupSwitches []DiggerSwitchBoomGroupLog) *DiggerLogger {
	ts := NewTimeSeries()
	var maxUtc int = 0
	for _, item := range boomGroupSwitches {
		if maxUtc < int(*item.ApplyUtc) {
			maxUtc = int(*item.ApplyUtc)
		}
		ts.Insert(TimePoint{Key: int64(*item.ApplyUtc), Value: item})
	}
	return &DiggerLogger{
		DiggerId:          diggerId,
		DiggerName:        diggerName,
		boomGroupSwitches: ts,
		MaxUtc:            maxUtc,
	}
}
func (t *DiggerLogger) FindLatest(utc int64) (DiggerSwitchBoomGroupLog, bool) {
	var tp TimePoint
	var found bool
	tp, found = t.boomGroupSwitches.FindClosest(utc)
	if !found {
		return DiggerSwitchBoomGroupLog{}, false
	}
	return tp.Value.(DiggerSwitchBoomGroupLog), true
}

func (t *DiggerLogger) AppendItemsBeyondTime(boomGroupSwitches []DiggerSwitchBoomGroupLog) {
	var maxUtc int = t.MaxUtc
	var beyondUtc int = t.MaxUtc
	for _, item := range boomGroupSwitches {
		if beyondUtc >= *item.ApplyUtc {
			Logger.Warn("AppendItemsBeyondTime", zap.String("condition", "applyUtc < beyondUtc"))
			continue
		}
		if maxUtc < *item.ApplyUtc {
			maxUtc = *item.ApplyUtc
		}
		t.boomGroupSwitches.Insert(TimePoint{Key: int64(*item.ApplyUtc), Value: item})
	}
	t.MaxUtc = maxUtc
}
