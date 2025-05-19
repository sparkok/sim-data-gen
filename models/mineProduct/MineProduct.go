package mineProduct

/**
* 实体类 MineProduct
 */
type MineProduct struct {
                                              
	ContentLimits *string                                               
	MatIndexes *string `gorm:"default:'[{"index":1,"threshold":0.1},{"index":2,"threshold":0.1},{"index":1,"threshold":0.1}]';"`
	Name       *string
	Status     *string
	Token      *string `gorm:"primaryKey;"`
}
