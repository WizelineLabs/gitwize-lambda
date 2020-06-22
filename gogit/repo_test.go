package gogit

import (
	"gitwize-lambda/utils"
	"os"
	"testing"
)

const (
	repoName = "mock-repo"
	repoURL  = "https://github.com/sang-d/mock-repo"
)

func TestIntegrationGetCommitIterFromBranch(t *testing.T) {
	if utils.IntegrationTestEnabled() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		repoPath := tmpDirectory + "/" + repoName
		os.RemoveAll(repoPath)
		r := GetRepo(repoName, repoURL, os.Getenv("DEFAULT_GITHUB_TOKEN"))
		obj := GetCommitIterFromBranch(r, "master", GetFullGitDateRange())
		if obj == nil {
			t.Errorf("Failed to get commit iter, check if test repo and branch exist %s", repoURL)
		}
		GetCommitIterFromBranch(r, "NonExistBranch ", GetFullGitDateRange()) // test panic
	}
}
