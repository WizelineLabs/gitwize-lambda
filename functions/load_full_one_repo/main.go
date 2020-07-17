package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"gitwize-lambda/db"
	"gitwize-lambda/github"
	"gitwize-lambda/gogit"
	"gitwize-lambda/utils"
	"log"
)

/*
	For big repo such as `go` git@github.com:golang/go.git, it will not be able to load all data in one lambda
	It'll need to refactor to multiple call, for example each call in a range 1000 commit.
*/

// Handler lambda function handler
func Handler(e gogit.RepoPayload) (string, error) {
	log.Println("Start loading full data for repo", e)
	conn := db.SQLDBConn()
	defer conn.Close()

	dateRange := gogit.GetFullGitDateRange()
	token := utils.GetAccessToken(e.RepoAccessToken)
	github.CollectPRsOfRepo(github.NewGithubPullRequestService(token), e.RepoID, e.URL, conn)
	gogit.UpdateDataForRepo(e.RepoID, e.URL, e.RepoName, token, e.Branch, dateRange, conn)
	db.UpdateMetricForRepo(e.RepoID)
	db.NewCommonOps().UpdateRepoLastUpdated(e.RepoID)
	resp := "Load full repo " + e.RepoName + " Completed"
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
