package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//  have a global DB and use it for all requests
var DB Database

type Database struct {
	*gorm.DB
}

// GetDB returns a Database connection to the postgres service with
// `databaseURL` as the data source.
// Example of `databaseURL` data source:
// 		"postgres://bob:secret@1.2.3.4:5432/mydb?sslmode=verify-full"
func GetDB(databaseURL string) (Database, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return Database{}, err
	}
	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return Database{}, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = Database{db}
	return DB, nil
}
