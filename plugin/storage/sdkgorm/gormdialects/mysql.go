package gormdialects

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySqlDB Get MySQL DB connection
// URI string
// Ex: user:password@/db_name?charset=utf8&parseTime=True&loc=Local
func MySqlDB(uri string) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(uri), &gorm.Config{})
}
