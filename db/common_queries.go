package db

import (
	"database/sql"
	"log"
	"strings"
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
