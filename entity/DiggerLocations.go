package entity

import . "sim_data_gen/models/locationGnssData"

type DiggerLocations struct {
	DiggerId   string
	DiggerName string
	Locations  *TimeSeries
}

func NewDiggerLocations(diggerId, diggerName string, locations []LocationGnssData) DiggerLocations {
	ts := NewTimeSeries()
	for _, location := range locations {
		ts.Insert(TimePoint{Key: int64(*location.Utc), Value: location})
	}
	return DiggerLocations{
		DiggerId:   diggerId,
		DiggerName: diggerName,
		Locations:  ts,
	}
}
func (t *DiggerLocations) FindLatestLocation(utc int64) (LocationGnssData, bool) {
	var tp TimePoint
	var found bool
	tp, found = t.Locations.FindClosest(utc)
	if !found {
		return LocationGnssData{}, false
	}
	return tp.Value.(LocationGnssData), true
}
