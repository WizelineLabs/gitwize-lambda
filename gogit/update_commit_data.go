package gogit

import (
	"database/sql"
	"github.com/GitWize/gitwize-lambda/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"log"
	"time"
)

// UpdateDataForRepo update data for public/private remote repo using in memory clone
func UpdateDataForRepo(repoID int, repoURL, repoName, token, branch string, dateRange DateRange, conn *sql.DB) {
	defer utils.TimeTrack(time.Now(), "UpdateDataForRepo")
	var r *git.Repository
	r = GetRepo(repoName, repoURL, token)
	commitIter := GetCommitIterFromBranch(r, branch, dateRange)
	updateCommitData(commitIter, repoID, conn)
}

func updateCommitData(commitIter object.CommitIter, repoID int, conn *sql.DB) {
	defer utils.TimeTrack(time.Now(), "updateCommitData")
	dtos := []commitDto{}
	err := commitIter.ForEach(func(c *object.Commit) error {
		if len(dtos) == batchSize {
			executeBulkStatement(dtos, conn)
			dtos = []commitDto{}
		} else {
			dto := getCommitDTO(c)
			dto.RepositoryID = repoID
			dtos = append(dtos, dto)
		}
		return nil
	})
	if err != nil {
		log.Panicln(err.Error())
	}
	if len(dtos) > 0 {
		executeBulkStatement(dtos, conn)
	}
}
