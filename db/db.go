package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
)

// SQLDBConn get database connection
func SQLDBConn() (db *sql.DB) {
	dbConnDSN := os.Getenv("DB_CONN_STRING")
	db, err := sql.Open("mysql", dbConnDSN)
	if err != nil {
		log.Fatalln("Failed to connect database", err)
	}
	return db
}
