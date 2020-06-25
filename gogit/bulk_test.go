package gogit

import (
	"testing"
)

var mockCommitTable = "mockCommitTable"
var mockFileTable = "mockFileTable"
var commitFields = getCommitFields()
var fileFields = getFileStatFields()
var cdto = commitDto{
	RepositoryID: 1,
	Hash:         "testhash",
	AuthorEmail:  "test@wizeline.com",
	AuthorName:   "test-user",
	Message:      "test message",
	NumFiles:     10,
	AdditionLOC:  100,
	DeletionLOC:  200,
	NumParents:   1,
	LOC:          10000,
	Year:         2020,
	Month:        1,
	Day:          3,
	Hour:         3,
	TimeStamp:    "2019-08-04 01:50:31",
}
var fdto = fileStatDTO{
	RepositoryID: 1,
	Hash:         "testhash",
	AuthorEmail:  "test@wizeline.com",
	AuthorName:   "test-user",
	FileName:     "test_file.go",
	AdditionLOC:  100,
	DeletionLOC:  200,
	Year:         2020,
	Month:        1,
	Day:          3,
	Hour:         3,
	TimeStamp:    "2019-08-04 01:50:31",
}

var expectedCommitStatement = "INSERT INTO mockCommitTable (repository_id, hash, author_email, author_name, message, num_files, addition_loc, deletion_loc, num_parents, total_loc, year, month, day, hour, commit_time_stamp) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE repository_id=repository_id"

func TestGenerateSQLStatementForCommit(t *testing.T) {
	dtos := []dtoInterface{
		dtoInterface(cdto),
		dtoInterface(cdto),
		dtoInterface(cdto),
	}

	statement, valArgs := generateSQLStatement(mockCommitTable, commitFields, dtos)
	if statement != expectedCommitStatement {
		t.Errorf("expected statement %s, got %s", expectedCommitStatement, statement)
	}
	if len(valArgs) != len(commitFields)*len(dtos) {
		t.Errorf("expected valArgs %d, got %d items", len(commitFields)*len(dtos), len(valArgs))
	}
}

var expectedFileStatement = "INSERT INTO mockFileTable (repository_id, hash, author_email, author_name, file_name, addition_loc, deletion_loc, year, month, day, hour, commit_time_stamp) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE repository_id=repository_id"

func TestGenerateSQLStatementForFile(t *testing.T) {
	dtos := []dtoInterface{
		dtoInterface(fdto),
		dtoInterface(fdto),
	}

	statement, valArgs := generateSQLStatement(mockFileTable, fileFields, dtos)
	if statement != expectedFileStatement {
		t.Errorf("expected statement %s got %s", expectedFileStatement, statement)
	}
	if len(valArgs) != len(fileFields)*len(dtos) {
		t.Errorf("expected valArgs %d, got %d items", len(fileFields)*len(dtos), len(valArgs))
	}
}
