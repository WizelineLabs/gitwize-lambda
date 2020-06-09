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
	"github.com/GitWize/gitwize-lambda/utils"
	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strings"
	"time"
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
	deleteStmt, err := conn.Prepare("DELETE FROM pull_request WHERE repository_id = ?")
	defer deleteStmt.Close()

	// delete old data
	_, err = deleteStmt.Exec(id)
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
			Sort:        "created",
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
		insertPRs(prSvc, id, prs, conn)

		listOpt.Page++
	}
}

func insertPRs(prSvc PullRequestService, repoID int, prs []*github.PullRequest, conn *sql.DB) {
	// Prepare statements
	insertStmt, err := conn.Prepare("INSERT INTO pull_request (repository_id, url, pr_no, title, body, head, base, state, created_by, created_year, created_month, created_day, created_hour, closed_year, closed_month, closed_day, closed_hour) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
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
		if pr.ClosedAt != nil {
			yearClosed := pr.ClosedAt.Year()
			monthClosed := yearClosed*100 + int(pr.ClosedAt.Month())
			dayClosed := monthClosed*100 + pr.ClosedAt.Day()
			hourClosed := dayClosed*100 + pr.ClosedAt.Hour()
			insertStmt.Exec(repoID, pr.HTMLURL, pr.Number, pr.Title, pr.Body, pr.Head.Ref, pr.Base.Ref, state, pr.User.Login,
				yearCreated, monthCreated, dayCreated, hourCreated,
				yearClosed, monthClosed, dayClosed, hourClosed)
		} else {
			insertStmt.Exec(repoID, pr.HTMLURL, pr.Number, pr.Title, pr.Body, pr.Head.Ref, pr.Base.Ref, state, pr.User.Login,
				yearCreated, monthCreated, dayCreated, hourCreated,
				nil, nil, nil, nil)
		}
	}
}