package gormdialects

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Get Postgres DB connection
// URI string
// Ex: host=myhost port=myport user=gorm dbname=gorm password=mypassword
func PostgresDB(uri string) (db *gorm.DB, err error) {
	return gorm.Open("postgres", uri)
}
