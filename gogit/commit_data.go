package gogit

import (
	"database/sql"
	"gitwize-lambda/utils"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitData struct {
	cdto  commitDto
	fdtos []fileStatDTO
}

// UpdateDataForRepo update data for public/private remote repo using in memory clone
func UpdateDataForRepo(repoID int, repoURL, repoName, token, branch string, dateRange DateRange, conn *sql.DB) {
	defer utils.TimeTrack(time.Now(), "UpdateDataForRepo "+repoName)
	var r *git.Repository
	r = GetRepo(repoName, repoURL, token)
	commitIter := GetCommitIterFromBranch(r, branch, dateRange)
	updateCommitAndFileStatData(commitIter, repoID, conn)
}

func processCommit(repoID int, c *object.Commit, ch chan CommitData) {
	cdto := getCommitDTO(c)
	cdto.RepositoryID = repoID
	fdtos := getFileStatDTO(c, repoID)
	data := CommitData{
		cdto:  cdto,
		fdtos: fdtos,
	}
	ch <- data
	return
}

func iterateCommits(repoID int, commitIter object.CommitIter, ch chan CommitData) (counter int) {
	err := commitIter.ForEach(func(c *object.Commit) error {
		go processCommit(repoID, c, ch)
		counter++
		return nil
	})
	if err != nil {
		log.Panicln(err.Error())
	}
	return counter
}

func updateCommitAndFileStatData(commitIter object.CommitIter, repoID int, conn *sql.DB) {
	ch := make(chan CommitData)
	counter := iterateCommits(repoID, commitIter, ch)
	log.Println("Number go routine created", counter)
	cdtos := []dtoInterface{}
	fdtos := []dtoInterface{}
	for i := 0; i < counter; i++ {
		data := <-ch
		cdtos = append(cdtos, dtoInterface(data.cdto))
		newFDtos := convertFileDtosToDtoInterfaces(data.fdtos)
		fdtos = append(fdtos, newFDtos...)
		if len(cdtos) >= batchSize || i == counter-1 {
			executeBulkStatement(commitTable, getCommitFields(), cdtos, conn)
			cdtos = []dtoInterface{}
		}
		if len(fdtos) >= batchSize || i == counter-1 {
			executeBulkStatement(fileStatTable, getFileStatFields(), fdtos, conn)
			fdtos = []dtoInterface{}
		}
	}
}
