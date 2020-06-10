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
	updateCommitAndFileStatData(commitIter, repoID, conn)
}

func updateCommitAndFileStatData(commitIter object.CommitIter, repoID int, conn *sql.DB) {
	defer utils.TimeTrack(time.Now(), "updateCommitAndFileStatData")
	cdtos := CommitDtos{}
	fdtos := FileStatDtos{}

	err := commitIter.ForEach(func(c *object.Commit) error {
		if len(cdtos.dtos) >= batchSize {
			executeBulkStatement(cdtos, conn)
			cdtos = CommitDtos{}
		} else {
			dto := getCommitDTO(c)
			dto.RepositoryID = repoID
			cdtos.append(dto)
		}

		if len(fdtos.dtos) >= batchSize {
			executeBulkStatement(fdtos, conn)
			fdtos = FileStatDtos{}
		} else {
			newFileDtos := getFileStatDTO(c, repoID)
			fdtos.append(newFileDtos)
		}

		return nil
	})
	if err != nil {
		log.Panicln(err.Error())
	}

	if len(cdtos.dtos) > 0 {
		executeBulkStatement(cdtos, conn)
	}
	if len(fdtos.dtos) > 0 {
		executeBulkStatement(fdtos, conn)
	}
}
