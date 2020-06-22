package utils

import (
	"gitwize-lambda/cypher"
	"log"
	"os"
	"time"
)

const (
	functionPrefix = "gitwize-lambda-"
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

// GetAppStage return deployed environment dev/qa/prod...
func GetAppStage() string {
	stage := os.Getenv("APP_STAGE")
	if stage == "" {
		stage = "dev"
	}
	return stage
}

// GetUpdateOneRepoFuncName return update-one-repo function name
func GetUpdateOneRepoFuncName() string {
	return functionPrefix + GetAppStage() + "-update_one_repo"
}

// IntegrationTestEnabled check if integration test mode enabled
func IntegrationTestEnabled() bool {
	enabled := os.Getenv("GITWIZE_INTEGRATION_TEST")
	if enabled == "TRUE" {
		os.Setenv("DB_CONN_STRING", "gitwize_user:P@ssword123@(localhost:3306)/gitwize?parseTime=true")
		os.Setenv("DEFAULT_GITHUB_TOKEN", "555748599586519a1cc7ed638ff3fd2234dfebf5") // token test acc https://github.com/TestAccWZL
		os.Setenv("USE_DEFAULT_API_TOKEN", "True")
	}
	return enabled == "TRUE"
}
