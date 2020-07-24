package gogit

import (
	"time"
)

const (
	commitTable   = "commit_data"
	fileStatTable = "file_stat_data"
	gitDateFormat = "2006-01-02"
	batchSize     = 1000
	tmpDirectory  = "/tmp"
)

type DateRange struct {
	Since *time.Time
	Until *time.Time
}

func GetFullGitDateRange() DateRange {
	since := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC)
	return DateRange{Since: &since, Until: &until}
}

func GetLastNDayDateRange(n int) DateRange {
	nDayBefore := time.Now().AddDate(0, 0, -n)
	tomorrow := time.Now().AddDate(0, 0, +1)
	since := time.Date(nDayBefore.Year(), nDayBefore.Month(), nDayBefore.Day(), 0, 0, 0, 0, time.UTC)
	until := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, time.UTC)
	return DateRange{Since: &since, Until: &until}
}

// RepoPayload event for update one repo
type RepoPayload struct {
	RepoID          int    `json:"RepoID"`
	URL             string `json:"URL"`
	RepoName        string `json:"RepoName"`
	RepoAccessToken string `json:"RepoAccessToken"`
	Branch          string `json:"Branch"`
}

func getRepoPath(repoName string) string {
	return tmpDirectory + "/" + repoName
}
