package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"gitwize-lambda/db"
)

// Handler lambda function handler
func Handler() (string, error) {
	db.UpdateMetricTable("db/update_metric_table.sql")
	return "load  metric completed", nil
}

func main() {
	lambda.Start(Handler)
}
