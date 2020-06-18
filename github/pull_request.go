/*
This module used to fetch Pull requests of git repositories (in repository table),
currently support github repositories only. Auth token is taken from `password` field.

How to run:
- call method: CollectPRs()
If you want to always use a default github token:
- Set USE_DEFAULT_API_TOKEN=true
- Set DEFAULT_GITHUB_TOKEN
*/

package github

import (
	"context"
	"database/sql"
	"github.com/wizeline/gitwize-lambda/utils"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

type PullRequestService interface {
	List(owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
}

type GithubPullRequestService struct {
	githubClient *github.Client
}

func NewGithubPullRequestService(token string) *GithubPullRequestService {
	return &GithubPullRequestService{
		githubClient: newGithubClient(token),
	}
}

func newGithubClient(token string) *github.Client {
	ctx := context.Background()

	useDefault := os.Getenv("USE_DEFAULT_API_TOKEN")
	if useDefault == "true" || token == "" {
		token = os.Getenv("DEFAULT_GITHUB_TOKEN")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func (s *GithubPullRequestService) List(owner string, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error) {
	return s.githubClient.PullRequests.List(context.Background(), owner, repo, opts)
}

func CollectPRsOfRepo(prSvc PullRequestService, id int, url string, conn *sql.DB) {
	defer utils.TimeTrack(time.Now(), "CollectPRsOfRepo")
	var owner, repo string
	if strings.HasPrefix(url, "git") {
		s := strings.Split(url, ":")[1]
		owner = strings.Split(s, "/")[0]
		repo = strings.Split(s, "/")[1]
	}
	if strings.HasPrefix(url, "https") {
		s := strings.Split(url, "/")
		owner = s[len(s)-2]
		repo = s[len(s)-1]
	}
	repo = strings.Replace(repo, ".git", "", -1)
	log.Printf("Collecting PRs: owner=%s, repo=%s", owner, repo)
	collectPRsOfRepo(prSvc, id, owner, repo, conn)
}

func collectPRsOfRepo(prSvc PullRequestService, id int, owner string, repo string, conn *sql.DB) {
	lastMetricUpdated := sql.NullTime{
		Time:  time.Unix(0, 0),
		Valid: false,
	}
	selectStmt, err := conn.Prepare("SELECT ctl_last_metric_updated FROM repository WHERE id = ?")
	defer selectStmt.Close()

	if err != nil {
		log.Printf("[ERROR] %s", err)
		return
	}

	err = selectStmt.QueryRow(id).Scan(&lastMetricUpdated)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return
	}

	// fetch PRs page by page, 100 per_page
	listOpt := github.ListOptions{
		Page:    1,
		PerPage: 100,
	}
	for {
		prs, _, err := prSvc.List(owner, repo, &github.PullRequestListOptions{
			State:       "all",
			Sort:        "updated",
			Direction:   "desc",
			ListOptions: listOpt,
		})

		if err != nil {
			log.Printf("[ERROR] %s", err)
			return
		}
		if len(prs) == 0 {
			break
		}

		filteredPrs := []*github.PullRequest{}
		stopFetching := false
		for _, pr := range prs {
			updatedAt := (*pr).UpdatedAt.Add(time.Duration(24) * time.Hour) // include prs updated on previous date in the run
			if updatedAt.After(lastMetricUpdated.Time) {
				filteredPrs = append(filteredPrs, pr)
			} else {
				stopFetching = true
				break
			}
		}
		insertPRs(prSvc, id, filteredPrs, conn)

		if stopFetching {
			break
		}
		listOpt.Page++
	}
}

func insertPRs(prSvc PullRequestService, repoID int, prs []*github.PullRequest, conn *sql.DB) {
	// Prepare statements
	sql := `INSERT INTO pull_request (repository_id, url, pr_no, title, body, head, base, state, created_by, created_year, created_month, created_day, created_hour, closed_year, closed_month, closed_day, closed_hour)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE state = ?, title = ?, body = ?, head = ?, base = ?, closed_year = ?, closed_month = ?, closed_day = ?, closed_hour = ?`
	insertStmt, err := conn.Prepare(sql)

	if err != nil {
		log.Printf("[ERROR] %s", err)
		return
	}
	defer insertStmt.Close()

	// insert again
	for _, pr := range prs {
		state := "open"
		if *pr.State == "closed" && pr.MergedAt != nil {
			state = "merged"
		} else if pr.State == nil {
			state = "rejected"
		}

		created := pr.CreatedAt.UTC()
		yearCreated := created.Year()
		monthCreated := yearCreated*100 + int(created.Month())
		dayCreated := monthCreated*100 + created.Day()
		hourCreated := dayCreated*100 + created.Hour()
		var yearClosed, monthClosed, dayClosed, hourClosed int
		if pr.ClosedAt != nil {
			yearClosed = pr.ClosedAt.Year()
			monthClosed = yearClosed*100 + int(pr.ClosedAt.Month())
			dayClosed = monthClosed*100 + pr.ClosedAt.Day()
			hourClosed = dayClosed*100 + pr.ClosedAt.Hour()
		}
		_, err := insertStmt.Exec(repoID, pr.HTMLURL, pr.Number, pr.Title, pr.Body, pr.Head.Ref, pr.Base.Ref, state, pr.User.Login,
			yearCreated, monthCreated, dayCreated, hourCreated,
			yearClosed, monthClosed, dayClosed, hourClosed,
			state, pr.Title, pr.Body, pr.Head.Ref, pr.Base.Ref, yearClosed, monthClosed, dayClosed, hourClosed)

		if err != nil {
			log.Printf("[ERROR] %s", err)
			return
		}
	}
}
