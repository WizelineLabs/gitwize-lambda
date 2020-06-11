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
	cdtos := []dtoInterface{}
	fdtos := []dtoInterface{}

	err := commitIter.ForEach(func(c *object.Commit) error {
		if len(cdtos) >= batchSize {
			executeBulkStatement(commitTable, getCommitFields(), cdtos, conn)
			cdtos = []dtoInterface{}
		} else {
			dto := getCommitDTO(c)
			dto.RepositoryID = repoID
			cdtos = append(cdtos, dtoInterface(dto))
		}

		if len(fdtos) >= batchSize {
			executeBulkStatement(fileStatTable, getFileStatFields(), fdtos, conn)
			fdtos = []dtoInterface{}
		} else {
			newFileDtos := getFileStatDTO(c, repoID)
			newDtos := convertFileDtosToDtoInterfaces(newFileDtos)
			fdtos = append(fdtos, newDtos...)
		}

		return nil
	})
	if err != nil {
		log.Panicln(err.Error())
	}

	if len(cdtos) > 0 {
		executeBulkStatement(commitTable, getCommitFields(), cdtos, conn)
	}
	if len(fdtos) > 0 {
		executeBulkStatement(fileStatTable, getFileStatFields(), fdtos, conn)
	}
}
