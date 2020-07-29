package db

import (
	"database/sql"
	"gitwize-lambda/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

// CommonOps common operation
type CommonOps struct{}

// NewCommonOps constructor for CommonOps
func NewCommonOps() CommonOps {
	return CommonOps{}
}

// GetAllRepoRows get all repository from db
func (t CommonOps) GetAllRepoRows(fields []string) *sql.Rows {
	conn := SQLDBConn()
	defer conn.Close()
	query := "SELECT " + strings.Join(fields, ", ") + " FROM repository"
	rows, err := conn.Query(query)
	if err != nil {
		log.Panicln(err)
	}
	return rows
}

// UpdateRepoLastUpdated update ctl_last_metric_updated
func (t CommonOps) UpdateRepoLastUpdated(id int) {
	conn := SQLDBConn()
	defer conn.Close()
	query := "UPDATE repository SET ctl_last_metric_updated = ?, status = \"AVAILABLE\" WHERE id = ?"
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
	_, err = stmt.Exec(time.Now(), id)

	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
}

// UpdateMetricTable execute db/update_metric_table.sql
func UpdateMetricTable() {
	defer utils.TimeTrack(time.Now(), "UpdateMetricTable")
	queries := strings.Split(GetUpdateMetricTableQuery(), ";\n")
	conn := SQLDBConn()
	defer conn.Close()

	for _, query := range queries {
		processUpdateQuery(query, conn)
	}
}

func processUpdateQuery(query string, conn *sql.DB) {
	if query = strings.TrimSpace(query); query == "" {
		return
	}
	result, err := conn.Exec(query)
	if err != nil && err.Error() != "Error 1065: Query was empty" {
		log.Panic(err)
	}
	count, _ := result.RowsAffected()
	log.Println("Number of rows affected", count)
}

// UpdateMetricForRepo execute db/update_metric_table.sql
func UpdateMetricForRepo(repoID int) {
	defer utils.TimeTrack(time.Now(), "UpdateMetricForRepo")
	queries := strings.Split(GetUpdateRepoQuery(), ";\n")
	conn := SQLDBConn()
	defer conn.Close()

	for _, query := range queries {
		query = strings.Replace(query, "$repoID", strconv.Itoa(repoID), -1)
		processUpdateQuery(query, conn)
	}
}

func GetUpdateMetricTableQuery() string {
	return `
SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=4;
-- commit
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 4, COUNT(*), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE num_parents<2
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=2;
-- line added
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 2, SUM(addition_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE num_parents<2
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=3;
-- line removed
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 3, SUM(deletion_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE num_parents<2
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=5;
-- pr created
INSERT INTO metric(repository_id, branch, type, year, month, day, hour, value)
SELECT repository_id, 'master' as branch, 5 as type, created_year as year,
	created_month as month, created_day as day, created_hour as hour, COUNT(*) as value
FROM pull_request
GROUP BY repository_id, created_year, created_month, created_day, created_hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=6;
-- pr merged
INSERT INTO metric(repository_id, branch, type, year, month, day, hour, value)
SELECT repository_id, 'master' as branch, 6 as type, closed_year as year,
	closed_month as month, closed_day as day, closed_hour as hour, COUNT(*) as value
FROM pull_request
WHERE state = 'merged'
GROUP BY repository_id, closed_year, closed_month, closed_day, closed_hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=7;
-- pr rejected
INSERT INTO metric(repository_id, branch, type, year, month, day, hour, value)
SELECT repository_id, 'master' as branch, 7 as type, closed_year as year,
	closed_month as month, closed_day as day, closed_hour as hour, COUNT(*) as value
FROM pull_request
WHERE state = 'rejected'
GROUP BY repository_id, closed_year, closed_month, closed_day, closed_hour
;

CALL calculate_metric_open_pr_all_repos()
;`
}

func GetUpdateRepoQuery() string {
	return `
SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=4 AND repository_id=$repoID;
-- commit
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 4, COUNT(*), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE repository_id=$repoID
AND num_parents<2
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=2 AND repository_id=$repoID;
-- line added
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 2, SUM(addition_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE repository_id=$repoID
AND num_parents<2
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=3 AND repository_id=$repoID;
-- line removed
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 3, SUM(deletion_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE repository_id=$repoID
AND num_parents<2
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=5 AND repository_id=$repoID;
-- pr created
INSERT INTO metric(repository_id, branch, type, year, month, day, hour, value)
SELECT repository_id, 'master' as branch, 5 as type, created_year as year,
	created_month as month, created_day as day, created_hour as hour, COUNT(*) as value
FROM pull_request
WHERE repository_id=$repoID
GROUP BY repository_id, created_year, created_month, created_day, created_hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=6 AND repository_id=$repoID;
-- pr merged
INSERT INTO metric(repository_id, branch, type, year, month, day, hour, value)
SELECT repository_id, 'master' as branch, 6 as type, closed_year as year,
	closed_month as month, closed_day as day, closed_hour as hour, COUNT(*) as value
FROM pull_request
WHERE state = 'merged' AND repository_id=$repoID
GROUP BY repository_id, closed_year, closed_month, closed_day, closed_hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=7 AND repository_id=$repoID;
-- pr rejected
INSERT INTO metric(repository_id, branch, type, year, month, day, hour, value)
SELECT repository_id, 'master' as branch, 7 as type, closed_year as year,
	closed_month as month, closed_day as day, closed_hour as hour, COUNT(*) as value
FROM pull_request
WHERE state = 'rejected' AND repository_id=$repoID
GROUP BY repository_id, closed_year, closed_month, closed_day, closed_hour
;

CALL calculate_metric_open_pr($repoID)
;`
}
