package gormdialects

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Get MySQL DB connection
// URI string
// Ex: user:password@/db_name?charset=utf8&parseTime=True&loc=Local
func MysqlDB(uri string) (db *gorm.DB, err error) {
	return gorm.Open("mysql", uri)
}
