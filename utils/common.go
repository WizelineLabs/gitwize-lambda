package utils

import (
	"github.com/GitWize/gitwize-lambda/cypher"
	"log"
	"os"
	"time"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("\n%s took %s", name, elapsed)
}

func GetAccessToken(repoPass string) (accessToken string) {
	if check := os.Getenv("USE_DEFAULT_API_TOKEN"); check != "" {
		accessToken = os.Getenv("DEFAULT_GITHUB_TOKEN")
	} else {
		accessToken = cypher.DecryptString(repoPass, os.Getenv("CYPHER_PASS_PHASE"))
	}
	return accessToken
}
