package unloadSite


/**
* 实体类 UnloadSite
*/
type UnloadSite struct {
                                              
	Capacity *string                                               
	Geom *string                                               
	Name *string                                               
	Nt *string                                               
	Token *string `gorm:"primaryKey;"`                                              
	X *float64                                               
	Y *float64 
}
