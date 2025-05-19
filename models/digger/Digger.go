package digger

/**
* 实体类 Digger
 */
type Digger struct {
	Attribs *string
	Name    *string  `gorm:"index:unique_lorry_name,unique"`
	Produce *float64 `gorm:"default:0;"`
	Speed   *float64 `gorm:"default:0;"`
	Status  *int     `gorm:"default:0;"`
	Token   *string  `gorm:"primaryKey;"`
	Utc     *int     `gorm:"default:0;"`
	X       *float64 `gorm:"default:0;"`
	Y       *float64 `gorm:"default:0;"`
}
