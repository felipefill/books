package utils

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/felipefill/books/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres dialect for GORM
)

var _db *gorm.DB

// GetDB gets DB connection, in case of failure it will panic
func GetDB() *gorm.DB {
	if _db != nil {
		return _db
	}

	host, name, user, pswd := getDatabaseInfo()
	db, err := gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", user, pswd, host, name))
	if err != nil {
		panic(fmt.Sprintf("Could not connect to database: %s", err.Error()))
	}

	_db = db
	migrateSchema(_db)

	return _db
}

// InjectDB injects given database
func InjectDB(db *sql.DB) {
	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		panic("Failed to convert sql.DB to gorm.DB")
	}

	_db = gormDB
}

func migrateSchema(db *gorm.DB) {
	db.AutoMigrate(&model.Book{})
}

func getDatabaseInfo() (host string, name string, user string, pswd string) {
	return mustGetEnvVar("DB_HOST"),
		mustGetEnvVar("DB_NAME"),
		mustGetEnvVar("DB_USER"),
		mustGetEnvVar("DB_PSWD")
}

func mustGetEnvVar(v string) string {
	envVar := os.Getenv(v)

	if envVar == "" {
		panic(fmt.Sprintf("Failed to retrieve %s environment variable", v))
	}

	return envVar
}
