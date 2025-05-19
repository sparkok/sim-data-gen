package yAnalyser

import "time"

/**
* 实体类 YAnalyser
 */
type YAnalyser struct {
	AnalyserNum   *string    `json:"AnalyserNum"`
	CreatedAt     *time.Time `json:"CreatedAt"`
	CrushingPlant *string    `json:"CrushingPlant"`
	Flux          *float64   `gorm:"default:0;" json:"Flux"`
	Load          *float64   `gorm:"default:0;" json:"Load"`
	Mat1          *float64   `gorm:"default:0;" json:"Mat1"`
	Mat10         *float64   `gorm:"default:0;" json:"Mat10"`
	Mat11         *float64   `gorm:"default:0;" json:"Mat11"`
	Mat12         *float64   `gorm:"default:0;" json:"Mat12"`
	Mat13         *float64   `gorm:"default:0;" json:"Mat13"`
	Mat14         *float64   `gorm:"default:0;" json:"Mat14"`
	Mat15         *float64   `gorm:"default:0;" json:"Mat15"`
	Mat16         *float64   `gorm:"default:0;" json:"Mat16"`
	Mat17         *float64   `gorm:"default:0;" json:"Mat17"`
	Mat18         *float64   `gorm:"default:0;" json:"Mat18"`
	Mat19         *float64   `gorm:"default:0;" json:"Mat19"`
	Mat2          *float64   `gorm:"default:0;" json:"Mat2"`
	Mat20         *float64   `gorm:"default:0;" json:"Mat20"`
	Mat3          *float64   `gorm:"default:0;" json:"Mat3"`
	Mat4          *float64   `gorm:"default:0;" json:"Mat4"`
	Mat5          *float64   `gorm:"default:0;" json:"Mat5"`
	Mat6          *float64   `gorm:"default:0;" json:"Mat6"`
	Mat7          *float64   `gorm:"default:0;" json:"Mat7"`
	Mat8          *float64   `gorm:"default:0;" json:"Mat8"`
	Mat9          *float64   `gorm:"default:0;" json:"Mat9"`
	Speed         *float64   `gorm:"default:0;" json:"Speed"`
	Status        *int       `gorm:"default:0;" json:"Status"`
	TestAt        *int       `gorm:"default:0;" json:"TestAt"`
	Token         *string    `gorm:"primaryKey;" json:"Token"`
}
