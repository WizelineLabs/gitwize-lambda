package main

import (
	"github.com/GitWize/gitwize-lambda/db"
	"github.com/GitWize/gitwize-lambda/github"
	"github.com/GitWize/gitwize-lambda/gogit"
	"github.com/GitWize/gitwize-lambda/utils"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	// need parameter of repo
	id, name, url, password := os.Args[1], os.Args[2], os.Args[3], os.Args[4]
	if id == "" || name == "" || url == "" {
		log.Panic("missing arguments")
	}
	repoID, _ := strconv.Atoi(id)

	defer utils.TimeTrack(time.Now(), "Get Data Single Repo"+name)

	token := utils.GetAccessToken(password)
	conn := db.SQLDBConn()
	gogit.UpdateDataForRepo(repoID, url, name, token, "", gogit.GetLastNDayDateRange(360), conn)
	github.CollectPRsOfRepo(github.NewGithubPullRequestService(token), repoID, url, conn)
}
