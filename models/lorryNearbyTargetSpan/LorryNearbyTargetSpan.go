package lorryNearbyTargetSpan


/**
* 实体类 LorryNearbyTargetSpan
*/
type LorryNearbyTargetSpan struct {
                                              
	BeginUtc *int `gorm:"default:0;"`                                              
	DateFlag *string                                               
	EndUtc *int `gorm:"default:0;"`                                              
	LorryId *string                                               
	Name *string                                               
	NearbyObj *string                                               
	ObjType *int `gorm:"default:0;"`                                              
	ProductName *string                                               
	Token *string `gorm:"primaryKey;"`
}
