package gormdialects

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Get SQLite DB connection
// URI string
// Ex: /tmp/gorm.db
func SQLiteDB(uri string) (db *gorm.DB, err error) {
	return gorm.Open("sqlite3", uri)
}
