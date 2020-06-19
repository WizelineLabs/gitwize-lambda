package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"gitwize-lambda/db"
)

// Handler lambda function handler
func Handler() (string, error) {
	db.UpdateMetricTable()
	return "load  metric completed", nil
}

func main() {
	lambda.Start(Handler)
}
