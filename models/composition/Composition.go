package composition


/**
* 实体类 Composition
*/
type Composition struct {
                                              
	Name *string                                               
	Token *string `gorm:"primaryKey;"`
}
