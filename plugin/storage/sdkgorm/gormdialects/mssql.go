package gormdialects

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

// Get MS SQL DB connection
// URI string
// Ex: sqlserver://username:password@localhost:1433?database=dbname
func MSSQLDB(uri string) (db *gorm.DB, err error) {
	return gorm.Open("mssql", uri)
}
