package location

/**
* 实体类 Location
 */
type Location struct {
	Alarm       *int
	Heading     *float32
	Name        *string
	Speed       *float32
	Status      *int
	Token       *string `gorm:"primaryKey;"`
	Utc         *int
	X           *float64
	Y           *float64
	LastCommUtc *int
	Elevation   *float64
}
