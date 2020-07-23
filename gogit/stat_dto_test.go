package gogit

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestGetDTOFromCommitObject(t *testing.T) {
	repoPath := tmpDirectory + "/" + repoName
	os.RemoveAll(repoPath)
	r := GetRepo("mock-repo", "https://github.com/sang-d/mock-repo.git", os.Getenv("DEFAULT_GITHUB_TOKEN"))

	expectedHash := "e15d6dad1576edf08811cb1b85a80c23b6d91153"
	expectedEmail := "sang.dinh@wizeline.com"
	expectedName := "Sang Dinh"

	commit, _ := r.CommitObject(plumbing.NewHash(expectedHash))
	dto := getCommitDTO(commit)

	if dto.Hash != expectedHash {
		t.Errorf("expected hash %s, got %s", expectedHash, dto.Hash)
	}
	if dto.AuthorEmail != expectedEmail {
		t.Errorf("expected author email %s, got %s", expectedEmail, dto.AuthorEmail)
	}
	if dto.AuthorName != expectedName {
		t.Errorf("expected author name %s, got %s", expectedName, dto.AuthorName)
	}
	if dto.NumParents != 1 {
		t.Errorf("expected number parents %d, got %d", 1, dto.NumParents)
	}
	if dto.AdditionLOC != 2 {
		t.Errorf("expected addition loc %d, got %d", 2, dto.AdditionLOC)
	}
	if dto.DeletionLOC != 0 {
		t.Errorf("expected deletion loc %d, got %d", 0, dto.DeletionLOC)
	}
	if dto.NumFiles != 1 {
		t.Errorf("expected num files loc %d, got %d", 1, dto.NumFiles)
	}
}

func TestGetFileStatsFromCommitObject(t *testing.T) {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: "https://github.com/sang-d/mock-repo.git",
	})
	if err != nil {
		panic(err.Error())
	}
	expectedHash := "e15d6dad1576edf08811cb1b85a80c23b6d91153"

	commit, _ := r.CommitObject(plumbing.NewHash(expectedHash))
	dtos := getFileStatDTO(commit, 1)

	if len(dtos) != 1 {
		t.Errorf("expected num files changes %d, got %d", 1, len(dtos))
	}
	dto := dtos[0]
	if dto.Hash != expectedHash {
		t.Errorf("expected hash %s, got %s", expectedHash, dto.Hash)
	}
	if dto.FileName != "hello.txt" {
		t.Errorf("expected file name %s, got %s", "hello.txt", dto.FileName)
	}
	if dto.AdditionLOC != 2 {
		t.Errorf("expected addition loc %d, got %d", 2, dto.AdditionLOC)
	}
	if dto.DeletionLOC != 0 {
		t.Errorf("expected deletion loc %d, got %d", 0, dto.DeletionLOC)
	}
}

func TestGetFieldNames(t *testing.T) {
	commitDtoItem := commitDto{}
	fileDtoItem := fileStatDTO{}

	expectedCommitDtoNames := "repository_id,hash,author_email,author_name,message,num_files,addition_loc,deletion_loc,num_parents,insertion_point,total_loc,year,month,day,hour,commit_time_stamp"
	assert.Equal(t, expectedCommitDtoNames, strings.Join(getFieldNames(commitDtoItem), ","))
	assert.Equal(t, expectedCommitDtoNames, strings.Join(commitDtoItem.getFieldNames(), ","))

	expectedFileStatNames := "repository_id,hash,author_email,author_name,file_name,addition_loc,deletion_loc,year,month,day,hour,commit_time_stamp"
	assert.Equal(t, expectedFileStatNames, strings.Join(getFieldNames(fileDtoItem), ","))
	assert.Equal(t, expectedFileStatNames, strings.Join(fileDtoItem.getFieldNames(), ","))

}
