package main

import (
	"gitwize-lambda/db"
	"gitwize-lambda/utils"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "Load Metric All Repos")
	db.UpdateMetricTable("db/update_metric_table.sql")
}
