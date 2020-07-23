package gogit

import (
	"database/sql"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gitwize-lambda/gitnative"
	"gitwize-lambda/utils"
	"log"
	"os/exec"
	"strconv"
	"time"
)

type CommitData struct {
	cdto  commitDto
	fdtos []fileStatDTO
}

// UpdateDataForRepo update data for public/private remote repo using in memory clone
func UpdateDataForRepo(repoID int, repoURL, repoName, token, branch string, dateRange DateRange, conn *sql.DB) {
	defer utils.TimeTrack(time.Now(), "UpdateDataForRepo "+repoName)
	r := GetRepo(repoName, repoURL, token)
	commitIter := GetCommitIterFromBranch(r, branch, dateRange)
	updateCommitAndFileStatData(commitIter, repoID, repoName, conn)
	updateFileStat(repoID, repoName, dateRange)
}

func updateFileStat(repoID int, repoName string, dateRange DateRange) {
	utils.SetDBConnSingleComponent()
	layout := "2006-01-02"
	vRepoPath := getRepoPath(repoName)
	command := fmt.Sprintf("./scripts/filestat.sh %s %s %s %s", strconv.Itoa(repoID), vRepoPath, dateRange.Since.Format(layout), dateRange.Until.Format(layout))
	out, err := exec.Command("/bin/sh", "-c", command).Output()
	if err != nil {
		log.Printf("updateFileStat failed: %s", err)
	}
	output := string(out[:])
	log.Println(output)
}

func processCommit(repoID int, repoName string, c *object.Commit, ch chan CommitData) {
	cdto := getCommitDTO(c)
	updateInsertionPoint(&cdto, repoName)
	cdto.RepositoryID = repoID
	fdtos := getFileStatDTO(c, repoID)
	data := CommitData{
		cdto:  cdto,
		fdtos: fdtos,
	}
	ch <- data
	return
}

func updateInsertionPoint(cdto *commitDto, repoName string) {
	insertionPoint, err := gitnative.GetInsertionPoint(getRepoPath(repoName), cdto.Hash)
	if err == nil {
		cdto.InsertionPoint = insertionPoint
	}
}

func iterateCommits(repoID int, repoName string, commitIter object.CommitIter, ch chan CommitData) (counter int) {
	err := commitIter.ForEach(func(c *object.Commit) error {
		go processCommit(repoID, repoName, c, ch)
		counter++
		return nil
	})
	if err != nil {
		log.Panicln(err.Error())
	}
	return counter
}

func updateCommitAndFileStatData(commitIter object.CommitIter, repoID int, repoName string, conn *sql.DB) {
	ch := make(chan CommitData)
	counter := iterateCommits(repoID, repoName, commitIter, ch)
	log.Println("Number go routine created", counter)
	cdtos := []dtoInterface{}
	fdtos := []dtoInterface{}
	for i := 0; i < counter; i++ {
		data := <-ch
		cdtos = append(cdtos, dtoInterface(data.cdto))
		newFDtos := convertFileDtosToDtoInterfaces(data.fdtos)
		fdtos = append(fdtos, newFDtos...)
		if len(cdtos) >= batchSize || i == counter-1 {
			executeBulkStatement(commitTable, cdtos, conn)
			cdtos = []dtoInterface{}
		}
		if len(fdtos) >= batchSize || i == counter-1 {
			executeBulkStatement(fileStatTable, fdtos, conn)
			fdtos = []dtoInterface{}
		}
	}
}
