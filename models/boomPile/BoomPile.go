package boomPile


/**
* 实体类 BoomPile
*/
type BoomPile struct {
                                              
	Bench *string                                               
	BoomDate *string                                               
	CostToGo *int `gorm:"default:0;"`                                              
	Geom *string                                               
	Material1 *float64 `gorm:"default:0;"`                                              
	Material10 *float64 `gorm:"default:0;"`                                              
	Material11 *float64 `gorm:"default:0;"`                                              
	Material12 *float64 `gorm:"default:0;"`                                              
	Material13 *float64 `gorm:"default:0;"`                                              
	Material14 *float64 `gorm:"default:0;"`                                              
	Material15 *float64 `gorm:"default:0;"`                                              
	Material16 *float64 `gorm:"default:0;"`                                              
	Material17 *float64 `gorm:"default:0;"`                                              
	Material18 *float64 `gorm:"default:0;"`                                              
	Material19 *float64 `gorm:"default:0;"`                                              
	Material2 *float64 `gorm:"default:0;"`                                              
	Material20 *float64 `gorm:"default:0;"`                                              
	Material3 *float64 `gorm:"default:0;"`                                              
	Material4 *float64 `gorm:"default:0;"`                                              
	Material5 *float64 `gorm:"default:0;"`                                              
	Material6 *float64 `gorm:"default:0;"`                                              
	Material7 *float64 `gorm:"default:0;"`                                              
	Material8 *float64 `gorm:"default:0;"`                                              
	Material9 *float64 `gorm:"default:0;"`                                              
	MineType *string                                               
	Name *string                                               
	Nt *string                                               
	Quantity *float64 `gorm:"default:0;"`                                              
	Status *string                                               
	Tag *string                                               
	Token *string `gorm:"primaryKey;"`                                              
	Used *float64 `gorm:"default:0;"`
}
