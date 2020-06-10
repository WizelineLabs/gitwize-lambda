package gogit

import (
	"database/sql"
	"log"
	"strings"
)

type dtosInterface interface {
	generateSQLStatement() (string, []interface{})
}

func (this CommitDtos) generateSQLStatement() (string, []interface{}) {
	dtos := this.dtos
	statement := "INSERT INTO " + commitTable + " (repository_id, hash, author_email, message, num_files, addition_loc, deletion_loc, num_parents, total_loc, year, month, day, hour, commit_time_stamp) "
	values := make([]string, len(dtos))
	valArgs := []interface{}{}
	for i, dto := range dtos {
		values[i] = "(" + strings.Repeat("?, ", 13) + "?)"
		args := dto.getListValues()
		valArgs = append(valArgs, args...)
	}
	statement = statement + "VALUES" + strings.Join(values, ",") + " ON DUPLICATE KEY UPDATE repository_id=repository_id"
	return statement, valArgs
}

func (this FileStatDtos) generateSQLStatement() (string, []interface{}) {
	dtos := this.dtos
	statement := "INSERT INTO " + fileStatTable + " (repository_id, hash, author_email, file_name, addition_loc, deletion_loc, year, month, day, hour, commit_time_stamp) "
	values := make([]string, len(dtos))
	valArgs := []interface{}{}
	for i, dto := range dtos {
		values[i] = "(" + strings.Repeat("?, ", 10) + "?)"
		args := dto.getListValues()
		valArgs = append(valArgs, args...)
	}
	statement = statement + "VALUES" + strings.Join(values, ",") + " ON DUPLICATE KEY UPDATE repository_id=repository_id"
	return statement, valArgs
}

func executeBulkStatement(i dtosInterface, conn *sql.DB) {
	statement, valArgs := i.generateSQLStatement()
	result, err := conn.Exec(statement, valArgs...)
	if err != nil {
		log.Panicln(err.Error())
	}
	rows, _ := result.RowsAffected()
	log.Println("number rows affected ", rows)
}
