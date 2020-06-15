package main

import (
	"gitwize-lambda/db"
	"gitwize-lambda/github"
	"gitwize-lambda/gogit"
	"gitwize-lambda/utils"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	id, name, url := os.Args[1], os.Args[2], os.Args[3]
	if id == "" || name == "" || url == "" {
		log.Panic("missing arguments")
	}
	repoID, _ := strconv.Atoi(id)

	defer utils.TimeTrack(time.Now(), "Get Data Single Repo "+name)

	password := ""
	if len(os.Args) > 4 {
		password = os.Args[4]
	}
	token := utils.GetAccessToken(password)
	conn := db.SQLDBConn()
	defer conn.Close()
	gogit.UpdateDataForRepo(repoID, url, name, token, "", gogit.GetLastNDayDateRange(360), conn)
	github.CollectPRsOfRepo(github.NewGithubPullRequestService(token), repoID, url, conn)
}
