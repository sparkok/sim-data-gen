package mineChangeStatus

import "time"

/**
* 实体类 MineChangeStatus
 */
type MineChangeStatus struct {
	Changing           *string
	CreatedAt          *time.Time
	DataTimeOfBridge   *time.Time
	DataTimeOfNeutron  *time.Time
	DataType           *int `gorm:"default:0;"` //准确度: 普通的数据 0 ~ 50, 重要数据 50~100
	DateFlag           *string
	Name               *string
	ProductName        *string
	Token              *string `gorm:"primaryKey;"`
	TotalMassOfBridge *float64 `gorm:"default:0;"`                                              
	TotalMassOfNeutron *float64 `gorm:"default:0;"`
}
