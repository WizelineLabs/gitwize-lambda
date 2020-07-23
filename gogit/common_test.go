package gogit

import (
	"testing"
	"time"
)

func TestGetFullGitDateRange(t *testing.T) {
	dateRange := GetFullGitDateRange()
	expectedSince := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedUntil := time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC)
	if expectedSince != *dateRange.Since {
		t.Errorf("expected since %s, got %s", expectedSince, *dateRange.Since)
	}
	if expectedUntil != *dateRange.Until {
		t.Errorf("expected until %s, got %s", expectedUntil, *dateRange.Until)
	}
}

func TestGetRepoPath(t *testing.T) {
	repoName := "mock-repo"
	expectedPath := "/tmp/mock-repo"
	if expectedPath != getRepoPath(repoName) {
		t.Errorf("expected: %s, got: %s", expectedPath, getRepoPath(repoName))
	}
}
