package main

import (
	"context"
	"gitwize-lambda/db"
	"gitwize-lambda/github"
	"gitwize-lambda/gogit"
	"gitwize-lambda/utils"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler lambda function handler
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	log.Println("querystring", request.QueryStringParameters)
	log.Println("pathparam", request.PathParameters)
	conn := db.SQLDBConn()
	defer conn.Close()

	dateRange := gogit.GetLastNDayDateRange(360)
	token := utils.GetAccessToken(request.QueryStringParameters["pass"])
	repoID, _ := strconv.Atoi(request.QueryStringParameters["id"])
	url := request.QueryStringParameters["url"]
	name := request.QueryStringParameters["name"]
	gogit.UpdateDataForRepo(repoID, url, name, token, "", dateRange, conn)
	github.CollectPRsOfRepo(github.NewGithubPullRequestService(token), repoID, url, conn)
	msg := "update completed for repo " + name
	return events.APIGatewayProxyResponse{Body: msg, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
