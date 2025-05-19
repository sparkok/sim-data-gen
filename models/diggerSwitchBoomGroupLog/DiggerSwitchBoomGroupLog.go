package diggerSwitchBoomGroupLog

/**
* 实体类 DiggerSwitchBoomGroupLog
 */
type DiggerSwitchBoomGroupLog struct {
	ApplyUtc    *int `gorm:"default:0;"`
	BoomGroupId *string
	DateFlag    *string
	DiggerId    *string
	Name        *string
	Status      *int    `gorm:"default:0;"`
	SubmitUtc   *int    `gorm:"default:0;"`
	Token       *string `gorm:"primaryKey;"`
}
