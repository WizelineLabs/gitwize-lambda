package db

import (
	"github.com/GitWize/gitwize-lambda/utils"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// UpdateMetricTable execute db/update_metric_table.sql
func UpdateMetricTable() {
	defer utils.TimeTrack(time.Now(), "UpdateMetricTable")

	log.Println("start loading metrics")
	file, err := ioutil.ReadFile("db/update_metric_table.sql")
	if err != nil {
		log.Fatal("Failed to read sql script: " + err.Error())
	}

	requests := strings.Split(string(file), ";\n")
	conn := SQLDBConn()
	defer conn.Close()

	for _, request := range requests {
		request = strings.TrimSpace(request)
		if request == "" {
			continue
		}
		log.Println(request)
		result, err := conn.Exec(request)
		if err != nil {
			emptySqlError := "Error 1065: Query was empty"
			if emptySqlError != err.Error() {
				log.Panic(err)
			}
		}
		count, _ := result.RowsAffected()
		log.Println("Number of rows affected", count)
	}
}
