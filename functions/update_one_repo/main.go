package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"gitwize-lambda/db"
	"gitwize-lambda/github"
	"gitwize-lambda/gogit"
	"gitwize-lambda/utils"
	"log"
)

// Handler lambda function handler
func Handler(e gogit.RepoPayload) (string, error) {
	log.Println("Repo Event", e)
	conn := db.SQLDBConn()
	defer conn.Close()

	dateRange := gogit.GetLastNDayDateRange(360)
	token := utils.GetAccessToken(e.RepoPass)
	gogit.UpdateDataForRepo(e.RepoID, e.URL, e.RepoName, token, e.Branch, dateRange, conn)
	github.CollectPRsOfRepo(github.NewGithubPullRequestService(token), e.RepoID, e.URL, conn)
	db.NewCommonOps().UpdateRepoLastUpdated(e.RepoID)
	resp := "Update Repo " + e.RepoName + " Completed"
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
