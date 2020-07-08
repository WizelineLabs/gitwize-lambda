package main

import (
	"database/sql"
	"gitwize-lambda/db"
	"gitwize-lambda/github"
	"gitwize-lambda/gogit"
	"gitwize-lambda/utils"
	"log"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "Get Commit & PR Data All Repo")

	fields := []string{"id", "name", "url", "access_token"}
	repoRows := db.NewCommonOps().GetAllRepoRows(fields)
	if repoRows == nil {
		log.Println("No repositories found")
		return
	}

	var id int
	var name, url, accessToken string

	count := 0
	conn := db.SQLDBConn()
	defer conn.Close()
	for repoRows.Next() {
		err := repoRows.Scan(&id, &name, &url, &accessToken)
		token := utils.GetAccessToken(accessToken)
		if err != nil {
			log.Println(err)
		} else {
			count++
			getDataOneRepo(id, url, name, token, conn)
		}
	}

	log.Printf("Done. %d repo updated successfully", count)

}

func getDataOneRepo(id int, url, name, token string, conn *sql.DB) {
	defer func() {
		r := recover()
		if r != nil {
			log.Println("Recover: ", r)
		}
		return
	}()
	gogit.UpdateDataForRepo(id, url, name, token, "", gogit.GetFullGitDateRange(), conn)
	github.CollectPRsOfRepo(github.NewGithubPullRequestService(token), id, url, conn)
	db.NewCommonOps().UpdateRepoLastUpdated(id)
}
