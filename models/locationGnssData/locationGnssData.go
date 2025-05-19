package locationGnssData

type LocationGnssData struct {
	//
	Alarm *int
	//
	Heading *float32
	//
	Name string
	//
	Speed *float32
	//
	Status *int
	//
	Token *string
	//
	Utc *int `gorm:"primaryKey" `
	//
	X *float64
	//
	Y *float64
	//
	LastCommUtc *int
	//
	Tid *string `gorm:"primaryKey" `
	//
	Nt *string
}
