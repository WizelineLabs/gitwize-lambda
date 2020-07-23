package gogit

import (
	"database/sql"
	"log"
	"strings"
)

func generateSQLStatement(table string, dtos []dtoInterface) (string, []interface{}) {
	fields := dtos[0].getFieldNames()
	statement := "INSERT INTO " + table + " (" + strings.Join(fields, ", ") + ") "
	values := make([]string, len(dtos))
	valArgs := []interface{}{}
	for i, dto := range dtos {
		values[i] = "(" + strings.Repeat("?, ", len(fields)-1) + "?)"
		args := dto.getListValues()
		valArgs = append(valArgs, args...)
	}
	statement = statement + "VALUES" + strings.Join(values, ",") + " ON DUPLICATE KEY UPDATE repository_id=repository_id"
	return statement, valArgs
}

func executeBulkStatement(table string, dtos []dtoInterface, conn *sql.DB) {
	statement, valArgs := generateSQLStatement(table, dtos)
	_, err := conn.Exec(statement, valArgs...)
	if err != nil {
		log.Panicln(err.Error())
	}
}
