package sqlite_boost

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

const spatialite = "spatialite"

func SupportSqliteBoost(spatial bool) {
	extensions := make([]string, 0)
	if spatial {
		extensions = append(extensions, "mod_spatialite")
	}
	sql.Register(spatialite,
		&sqlite3.SQLiteDriver{
			Extensions: extensions,
		})
}
