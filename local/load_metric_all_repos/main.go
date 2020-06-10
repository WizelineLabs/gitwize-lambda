package main

import (
	"github.com/GitWize/gitwize-lambda/db"
	"github.com/GitWize/gitwize-lambda/utils"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "Load Metric All Repos")
	db.UpdateMetricTable()
}
