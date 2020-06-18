package main

import (
	"database/sql"
	"github.com/wizeline/gitwize-lambda/db"
	"github.com/wizeline/gitwize-lambda/github"
	"github.com/wizeline/gitwize-lambda/gogit"
	"github.com/wizeline/gitwize-lambda/utils"
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
	conn := db.SQLDBConn()
	defer conn.Close()
	c := make(chan bool)
	for repoRows.Next() {
		err := repoRows.Scan(&id, &name, &url, &password)
		token := utils.GetAccessToken(password)
		if err != nil {
			log.Println(err)
		} else {
			count++
			go getDataOneRepo(c, id, url, name, token, conn)
		}
	}

	successCount, failCount := 0, 0
	for i := 0; i < count; i++ {
		if <-c {
			successCount++
		} else {
			failCount++
		}
	}
	log.Printf("Done. %d repo updated successfully. %d repo failed", successCount, failCount)

}

func getDataOneRepo(c chan bool, id int, url, name, token string, conn *sql.DB) {
	flag := false
	defer func() {
		r := recover()
		if r != nil {
			log.Println("Recover: ", r)
		}
		c <- flag
		return
	}()
	gogit.UpdateDataForRepo(id, url, name, token, "", gogit.GetLastNDayDateRange(360), conn)
	github.CollectPRsOfRepo(github.NewGithubPullRequestService(token), id, url, conn)
	db.UpdateRepoLastUpdated(id)
	flag = true
}
