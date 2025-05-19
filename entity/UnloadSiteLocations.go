package entity

import (
	. "sim_data_gen/models/unloadSite"
	. "sim_data_gen/utils"
)

type UnloadSiteLocation struct {
	LocationList []UnloadSite
}

func NewUnloadSiteLocation(unloadSites []UnloadSite) UnloadSiteLocation {
	var locationList []UnloadSite
	for _, unloadSiteObj := range unloadSites {
		if unloadSiteObj.X == nil || unloadSiteObj.Y == nil || unloadSiteObj.Name == nil || unloadSiteObj.Token == nil {
			Logger.Error("UnloadSiteLocation.NewUnloadSiteLocation: unloadSiteObj data is invalid")
			continue
		}
		locationList = append(locationList, unloadSiteObj)
	}
	return UnloadSiteLocation{
		LocationList: locationList,
	}
}
func (t *UnloadSiteLocation) FindNearby(x float64, y float64, rByMeter float64) (site *UnloadSite) {
	for _, unloadSiteObj := range t.LocationList {
		if CalculateDistanceAsMeter(*unloadSiteObj.X, *unloadSiteObj.Y, x, y) < rByMeter {
			return &unloadSiteObj
		}
	}
	return nil
}
func (t *UnloadSiteLocation) Valid() bool {
	return len(t.LocationList) > 0
}
