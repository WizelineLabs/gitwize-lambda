package db

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

// GetAllRepoRows get all repository from db
func GetAllRepoRows(fields []string) *sql.Rows {
	conn := SQLDBConn()
	defer conn.Close()
	query := "SELECT " + strings.Join(fields, ", ") + " FROM repository"
	rows, err := conn.Query(query)
	if err != nil {
		log.Panicln(err)
	}
	return rows
}

// UpdateRepoLastUpdated update ctl_last_metric_updated
func UpdateRepoLastUpdated(id int) {
	conn := SQLDBConn()
	defer conn.Close()
	query := "UPDATE repository SET ctl_last_metric_updated = ? WHERE id = ?"
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
	_, err = stmt.Exec(time.Now(), id)

	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
}
