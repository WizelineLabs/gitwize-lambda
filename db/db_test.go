package db

import (
	"gitwize-lambda/utils"
	"testing"
)

func TestMain(t *testing.T) {
	utils.SetupIntegrationTest()
}

func TestSQLDBConn(t *testing.T) {
	db := SQLDBConn()
	if err := db.Ping(); err != nil {
		t.Error("Failed to connect db")
	}
	db.Close()
}

func TestUpdateMetricTable(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Failed to update metric table")
		}
	}()
	UpdateMetricTable("update_metric_table.sql")
}

func TestGetAllRepoRows(t *testing.T) {
	fields := []string{"id", "name", "url", "password"}
	rows := NewCommonOps().GetAllRepoRows(fields)
	if rows == nil {
		t.Error("No repo found, check init query loaded repo to mysql docker")
	}
}
