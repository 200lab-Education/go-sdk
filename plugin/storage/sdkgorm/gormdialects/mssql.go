package gormdialects

import (
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// MSSqlDB Get MS SQL DB connection
// URI string
// Ex: sqlserver://username:password@localhost:1433?database=dbname
func MSSqlDB(uri string) (db *gorm.DB, err error) {
	return gorm.Open(sqlserver.Open(uri))
}
