package sqlite_boost

import (
	"database/sql/driver"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

// // Time 自定义时间类型，专用于SQLite数据库交互
type Time struct {
	time.Time
}

// Value 实现Valuer接口，把自定义时间类型转成数据库能存储的值
func (t Time) Value() (driver.Value, error) {
	return t.Format("2006-01-01 15:04:05"), nil
}

// Scanner 实现Scanner接口，把数据库读取的值转成自定义时间类型
func (t *Time) Scanner(v interface{}) error {
	switch vt := v.(type) {
	case string:
		var err error
		t.Time, err = time.Parse("2006-01-01 15:04:05", vt)
		return err
	case []byte:
		var err error
		t.Time, err = time.Parse("2006-01-01 15:04:05", string(vt))
		return err
	default:
		return errors.Errorf("无法将%v转为时间类型", vt)
	}
}

// SQLiteConfig 用于给gorm配置SQLite数据库相关设置
func SQLiteConfig() *gorm.Config {
	return &gorm.Config{
		//NowFunc: func() time.Time {
		//	return time.Now()
		//},
		//CreateTime: func(t time.Time) interface{} {
		//	return Time{t}
		//},
		//UpdateTime: func(t time.Time) interface{} {
		//	return Time{t}
		//},
	}
}
