package utils

import (
	"github.com/GitWize/gitwize-lambda/cypher"
	"log"
	"os"
	"time"
)

// TimeTrack use with defer to track processing time of a function
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("\n%s took %s", name, elapsed)
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
