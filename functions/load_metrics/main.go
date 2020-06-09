package main

import (
	"context"
	"github.com/GitWize/gitwize-lambda/db"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

// Handler lambda function handler
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	db.UpdateMetricTable()
	msg := "load  metric completed"
	return events.APIGatewayProxyResponse{Body: msg, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
