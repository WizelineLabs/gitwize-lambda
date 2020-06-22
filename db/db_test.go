package db

import (
	"gitwize-lambda/utils"
	"testing"
	// "fmt"
	// "os"
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
		UpdateMetricTable("update_metric_table.sql")
	}
}

func TestGetAllRepoRows(t *testing.T) {
	if utils.IntegrationTestEnabled() {
		fields := []string{"id", "name", "url", "password"}
		rows := NewCommonOps().GetAllRepoRows(fields)
		if rows == nil {
			t.Error("No repo found, check init query loaded repo to mysql docker")
		}
	}
}

// func TestMain(m *testing.M) {
//     // call flag.Parse() here if TestMain uses flags
//     rc := m.Run()

//     // rc 0 means we've passed,
//     // and CoverMode will be non empty if run with -cover
//     if rc == 0 && testing.CoverMode() != "" {
//         c := testing.Coverage()
//         if c < 0.5 {
//             fmt.Println("Tests passed but coverage failed at", c)
//             rc = -1
//         }
//     }
//     os.Exit(rc)
// }
