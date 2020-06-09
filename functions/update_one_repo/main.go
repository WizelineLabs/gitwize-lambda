package main

import (
	"github.com/GitWize/gitwize-lambda/db"
	"github.com/GitWize/gitwize-lambda/github"
	"github.com/GitWize/gitwize-lambda/gogit"
	"github.com/GitWize/gitwize-lambda/utils"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

// Handler lambda function handler
func Handler(e gogit.RepoEvent) (string, error) {
	log.Println("Repo Event", e)
	conn := db.SQLDBConn()
	defer conn.Close()

	dateRange := gogit.GetLastNDayDateRange(360)
	token := utils.GetAccessToken(e.RepoPass)
	gogit.UpdateDataForRepo(e.RepoID, e.URL, e.RepoName, token, e.Branch, dateRange, conn)
	github.CollectPRsOfRepo(github.NewGithubPullRequestService(token), e.RepoID, e.URL, conn)
	resp := "Update Repo Completed"
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
