package main

import (
	"github.com/wizeline/gitwize-lambda/gogit"
	"log"
	"os"
	"time"
)

// set env var before execute
// export DEFAULT_GITHUB_TOKEN=
// export DB_CONN_STRING
func main() {
	log.Println("Test function locally")
	url := "https://github.com/go-git/go-git"
	token := os.Getenv("DEFAULT_GITHUB_TOKEN")
	go gogit.GetRepo("go-git", url, token)

	url = "https://github.com/wizeline/gitwize-be"
	go gogit.GetRepo("gitwize-be", url, token)
	time.Sleep(10 * time.Second)
}
