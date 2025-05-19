package utils

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

var Db *gorm.DB

func GetDb(newDb ...*gorm.DB) *gorm.DB {
	if newDb != nil {
		return newDb[0]
	}
	return Db
}

var CurrentDbType string
var ZERO_TIME_TXT string
var STANDARD_TIME_FORMAT string

func makeDialector(dbType string, dbUrl string) (gorm.Dialector, error) {
	CurrentDbType = dbType

	switch dbType {
	case "mysql":
		{
			STANDARD_TIME_FORMAT = "2006-01-02T15:04:05Z07:00"
			ZERO_TIME_TXT = time.Unix(0, 0).Format(STANDARD_TIME_FORMAT)
			return mysql.Open(dbUrl), nil
		}
	case "spatialite":
		//sqlite_boost.SupportSqliteBoost(true)
		STANDARD_TIME_FORMAT = "2006-01-02 15:04:05"
		ZERO_TIME_TXT = time.Unix(0, 0).UTC().Format(STANDARD_TIME_FORMAT)
		//return sql.Open("spatialite", dbUrl), nil
		//if err = testSpatialSupport(db); err != nil {
		//	return fmt.Errorf("test spatialite: %w", err)
		//}
		return sqlite.Open(dbUrl), nil
	case "sqlite":
		{
			STANDARD_TIME_FORMAT = "2006-01-02 15:04:05"
			ZERO_TIME_TXT = time.Unix(0, 0).UTC().Format(STANDARD_TIME_FORMAT)
			return sqlite.Open(dbUrl), nil
		}
	case "postgres":
		{
			STANDARD_TIME_FORMAT = "2006-01-02T15:04:05Z07:00"
			ZERO_TIME_TXT = time.Unix(0, 0).Format(STANDARD_TIME_FORMAT)
			return postgres.Open(dbUrl), nil
		}
	default:
		return nil, errors.New("unsupported driver")
	}

}
func ConnectDB(dsnVarName string, dbDriverVarName string) {
	var dsn string
	var err error
	var dbDriverName string
	if dsn = GetConfig().String(dsnVarName, ""); dsn == "" {
		Logger.Error(err.Error())
		return
	}
	if dbDriverName = GetConfig().String(dbDriverVarName, ""); dbDriverName == "" {
		Logger.Error(err.Error())
		return
	}
	logLevel := logger.Warn
	if dbDebug := GetConfig().Bool("dbDebug", false); dbDebug == true {
		logLevel = logger.Info
	}
	var dialector gorm.Dialector
	if dialector, err = makeDialector(dbDriverName, dsn); err != nil {
		Logger.Error(err.Error())
		return
	}
	config := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",    // table name prefix, table for `User` would be `t_users`
			SingularTable: true,  // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   false, // skip the snake_casing of names
			NameReplacer:  nil,   //strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
		Logger: logger.Default.LogMode(logLevel),
	}
	if Db, err = gorm.Open(dialector, &config); err != nil {
		Logger.Error(err.Error())
		return
	}
	if CurrentDbType == "spatialite" {
		// 执行SQL加载SpatiaLite扩展（需确保mod_spatialite路径正确）
		Db.Exec("SELECT load_extension('mod_spatialite')")
		// 初始化空间元数据（创建geometry_columns表等）
		Db.Exec("SELECT InitSpatialMetadata(1)")
	}
}
func AutoMigrate(autoBuildDbVarName string, dst []interface{}) {
	var auotBuild string
	if auotBuild = GetConfig().String(autoBuildDbVarName, "false"); auotBuild != "auto" {
		Logger.Info("跳过自动创建")
		return
	}
	for _, model := range dst {
		GetDb().AutoMigrate(model)
	}
}
