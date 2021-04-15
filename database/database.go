package database

import (
	"database/sql"
	"github.com/bwmarrin/lit"
	// Possible drivers
	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

// DB is the database connection
var DB *sql.DB

// InitializeDatabase initializes db connection given the driver and datasource name.
// Supported driver are sqlite and mysql.
func InitializeDatabase(driver, dsn string) {
	var err error
	// Open database connection
	DB, err = sql.Open(driver, dsn)
	if err != nil {
		lit.Error("Error opening db connection: %v", err)
		return
	}
}
