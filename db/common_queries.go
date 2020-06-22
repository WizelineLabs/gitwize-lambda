package db

import (
	"database/sql"
	"gitwize-lambda/utils"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// CommonOps common operation
type CommonOps struct{}

// NewCommonOps constructor for CommonOps
func NewCommonOps() CommonOps {
	return CommonOps{}
}

// GetAllRepoRows get all repository from db
func (t CommonOps) GetAllRepoRows(fields []string) *sql.Rows {
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
func (t CommonOps) UpdateRepoLastUpdated(id int) {
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

// UpdateMetricTable execute db/update_metric_table.sql
func UpdateMetricTable(sqlFile string) {
	defer utils.TimeTrack(time.Now(), "UpdateMetricTable")

	log.Println("start loading metrics")
	file, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Panicln("Failed to read sql script: " + err.Error())
	}

	queries := strings.Split(string(file), ";\n")
	conn := SQLDBConn()
	defer conn.Close()

	for _, query := range queries {
		processUpdateQuery(query, conn)
	}
}

func processUpdateQuery(query string, conn *sql.DB) {
	if query = strings.TrimSpace(query); query == "" {
		return
	}
	result, err := conn.Exec(query)
	if err != nil && err.Error() != "Error 1065: Query was empty" {
		log.Panic(err)
	}
	count, _ := result.RowsAffected()
	log.Println("Number of rows affected", count)
}
