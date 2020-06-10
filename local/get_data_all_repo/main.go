package main

import (
	"github.com/GitWize/gitwize-lambda/db"
	"github.com/GitWize/gitwize-lambda/github"
	"github.com/GitWize/gitwize-lambda/gogit"
	"github.com/GitWize/gitwize-lambda/utils"
	"log"
	"time"
)

func main() {
	defer utils.TimeTrack(time.Now(), "Get Commit & PR Data All Repo")

	fields := []string{"id", "name", "url", "password"}
	repoRows := db.GetAllRepoRows(fields)
	if repoRows == nil {
		log.Println("No repositories found")
		return
	}

	var id int
	var name, url, password string

	count := 0
	for repoRows.Next() {
		count++
		err := repoRows.Scan(&id, &name, &url, &password)
		token := utils.GetAccessToken(password)
		if err != nil {
			log.Panicln(err)
		} else {
			conn := db.SQLDBConn()
			defer conn.Close()
			gogit.UpdateDataForRepo(id, url, name, token, "", gogit.GetLastNDayDateRange(360), conn)
			github.CollectPRsOfRepo(github.NewGithubPullRequestService(token), id, url, conn)
		}
	}
	log.Println("Completed update ", count, "repositories")
}
