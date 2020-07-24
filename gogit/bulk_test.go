package gogit

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var mockCommitTable = "mockCommitTable"
var mockFileTable = "mockFileTable"

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

func TestGenerateSQLStatementForCommit(t *testing.T) {
	dtos := []dtoInterface{
		dtoInterface(cdto),
		dtoInterface(cdto),
		dtoInterface(cdto),
	}
	commitFields := cdto.getFieldNames()
	value := "(" + strings.Repeat("?, ", len(commitFields)-1) + "?)"
	values := "VALUES" + value + "," + value + "," + value
	expectedCommitStatement := "INSERT INTO mockCommitTable (" + strings.Join(commitFields, ", ") + ") " + values + " ON DUPLICATE KEY UPDATE repository_id=repository_id"

	statement, valArgs := generateSQLStatement(mockCommitTable, dtos)
	assert.Equal(t, expectedCommitStatement, statement)
	assert.Equal(t, len(commitFields)*len(dtos), len(valArgs))
}

func TestGenerateSQLStatementForFile(t *testing.T) {
	dtos := []dtoInterface{
		dtoInterface(fdto),
		dtoInterface(fdto),
	}
	fileFields := fdto.getFieldNames()
	value := "(" + strings.Repeat("?, ", len(fileFields)-1) + "?)"
	values := "VALUES" + value + "," + value
	var expectedFileStatement = "INSERT INTO mockFileTable (" + strings.Join(fileFields, ", ") + ") " + values + " ON DUPLICATE KEY UPDATE repository_id=repository_id"

	statement, valArgs := generateSQLStatement(mockFileTable, dtos)
	assert.Equal(t, expectedFileStatement, statement)
	assert.Equal(t, len(fileFields)*len(dtos), len(valArgs))
}
