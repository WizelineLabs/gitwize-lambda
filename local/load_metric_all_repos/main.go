package main

import (
	"github.com/wizeline/gitwize-lambda/db"
	"github.com/wizeline/gitwize-lambda/utils"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "Load Metric All Repos")
	db.UpdateMetricTable()
}
