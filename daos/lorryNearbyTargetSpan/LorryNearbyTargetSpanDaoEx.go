package lorryNearbyTargetSpan

import (
	lorryNearbyTargetSpanModel "sim_data_gen/models/lorryNearbyTargetSpan"
	. "sim_data_gen/utils"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func FindLatestObjBeforeTime(lorryToken string, dateAsStr string, utc int, tx ...*gorm.DB) (*lorryNearbyTargetSpanModel.LorryNearbyTargetSpan, error) {
	var result lorryNearbyTargetSpanModel.LorryNearbyTargetSpan
	condition := fmt.Sprintf("end_utc < %d and lorry_id = '%s' and date_flag = '%s'", utc, lorryToken, dateAsStr)
	db := GetDb(tx...).Where(condition).Take(&result)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		Logger.Error("failed to FindLatestObjBeforeTime", zap.Error(db.Error))
	}
	return &result, db.Error
}

func ListLatestNUntil(lorryId string, dateAsStr string, someTime int64, n int, tx ...*gorm.DB) ([]lorryNearbyTargetSpanModel.LorryNearbyTargetSpan, error) {
	list := []lorryNearbyTargetSpanModel.LorryNearbyTargetSpan{}
	db := GetDb(tx...).Table("lorry_nearby_target_span").Where("lorry_id = ? and ( begin_utc < ? )  and date_flag = '%s'", lorryId, someTime, dateAsStr).Order("end_utc DESC").Limit(n).Find(&list)
	return list, db.Error
}

func FindObjIncludeTime(lorryToken string, dateAsStr string, target string, utc int, tx ...*gorm.DB) (*lorryNearbyTargetSpanModel.LorryNearbyTargetSpan, error) {
	var result lorryNearbyTargetSpanModel.LorryNearbyTargetSpan
	condition := fmt.Sprintf("begin_utc <= %d and end_utc >= %d and lorry_id = '%s' and date_flag = '%s'", utc, utc, lorryToken, dateAsStr)
	if len(target) > 0 {
		condition += fmt.Sprintf(" and nearby_obj = '%s'", target)
	}
	db := GetDb(tx...).Where(condition).Take(&result)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		Logger.Error("failed to FindObjIncludeTime", zap.Error(db.Error))
	}
	return &result, db.Error
}
