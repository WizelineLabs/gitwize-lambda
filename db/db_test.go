package db

import (
	"gitwize-lambda/utils"
	"testing"
)

func TestSQLDBConn(t *testing.T) {
	if utils.IntegrationTestEnabled() {
		db := SQLDBConn()
		if err := db.Ping(); err != nil {
			t.Error("Failed to connect db")
		}
		db.Close()
	}
}

func TestUpdateMetricTable(t *testing.T) {
	if utils.IntegrationTestEnabled() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Failed to update metric table")
			}
		}()
		UpdateMetricTable()
	}
}

func TestGetAllRepoRows(t *testing.T) {
	if utils.IntegrationTestEnabled() {
		fields := []string{"id", "name", "url", "access_token"}
		rows := NewCommonOps().GetAllRepoRows(fields)
		if rows == nil {
			t.Error("No repo found, check init query loaded repo to mysql docker")
		}
	}
}

func TestGetUpdateMetricTableQuery(t *testing.T) {
	expected := `
SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=4;
-- commit
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 4, COUNT(*), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE num_parents=1
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=2;
-- line added
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 2, SUM(addition_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE num_parents=1
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=3;
-- line removed
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 3, SUM(deletion_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE num_parents=1
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
	if expected != GetUpdateMetricTableQuery() {
		t.Errorf("expected query %s, got %s", expected, GetUpdateMetricTableQuery())
	}
}

func TestGetUpdateRepoQuery(t *testing.T) {
	expected := `
SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=4 AND repository_id=$repoID;
-- commit
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 4, COUNT(*), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE repository_id=$repoID
AND num_parents=1
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=2 AND repository_id=$repoID;
-- line added
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 2, SUM(addition_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE repository_id=$repoID
AND num_parents=1
GROUP BY repository_id, year, month, day, hour
;

SET SQL_SAFE_UPDATES = 0;
DELETE FROM metric WHERE branch='master' AND type=3 AND repository_id=$repoID;
-- line removed
INSERT INTO metric (repository_id, branch, type, value, year, month, day, hour)
SELECT repository_id, 'master', 3, SUM(deletion_loc), year, year*100+month, (year*100+month)*100+day, (year*10000+month*100+day)*100+hour
FROM commit_data
WHERE repository_id=$repoID
AND num_parents=1
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

	if expected != GetUpdateRepoQuery() {
		t.Errorf("expected query %s, got %s", expected, GetUpdateRepoQuery())
	}
}
