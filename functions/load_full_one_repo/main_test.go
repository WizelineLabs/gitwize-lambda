package main

import (
	"gitwize-lambda/gogit"
	"gitwize-lambda/utils"
	"testing"
)

func TestHandler(t *testing.T) {
	if utils.IntegrationTestEnabled() {
		payload := gogit.RepoPayload{
			RepoID:   1,
			URL:      "https://github.com/sang-d/mock-repo",
			RepoName: "mock-repo-one",
			RepoPass: "",
			Branch:   "",
		}
		_, err := Handler(payload)
		if err != nil {
			t.Errorf("Error %s", err)
		}
	}
}
