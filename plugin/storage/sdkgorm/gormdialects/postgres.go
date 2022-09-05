package gormdialects

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDB Get Postgres DB connection
// URI string
// Ex: host=myhost port=myport user=gorm dbname=gorm password=mypassword
func PostgresDB(uri string) (db *gorm.DB, err error) {
	return gorm.Open(postgres.Open(uri))
}
