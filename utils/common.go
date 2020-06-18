package utils

import (
	"github.com/wizeline/gitwize-lambda/cypher"
	"log"
	"os"
	"time"
)

const (
	functionPrefix = "github.com/wizeline/gitwize-lambda-"
)

// TimeTrack use with defer to track processing time of a function
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// GetAccessToken retrieve access token from db or environ
func GetAccessToken(repoPass string) (accessToken string) {
	check := os.Getenv("USE_DEFAULT_API_TOKEN")
	if check != "" || repoPass == "" {
		accessToken = os.Getenv("DEFAULT_GITHUB_TOKEN")
	} else {
		accessToken = cypher.DecryptString(repoPass, os.Getenv("CYPHER_PASS_PHASE"))
	}
	return accessToken
}

//GetAppStage return deployed environment dev/qa/prod...
func GetAppStage() string {
	stage := os.Getenv("APP_STAGE")
	if stage == "" {
		stage = "dev"
	}
	return stage
}

//GetUpdateOneRepoFuncName return update-one-repo function name
func GetUpdateOneRepoFuncName() string {
	return functionPrefix + GetAppStage() + "-update_one_repo"
}
