package gogit

import (
	"strings"
	"testing"
	"time"
)

func TestGetFullGitDateRange(t *testing.T) {
	dateRange := GetFullGitDateRange()
	expectedSince := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedUntil := time.Date(2040, 1, 1, 0, 0, 0, 0, time.UTC)
	if expectedSince != *dateRange.Since {
		t.Errorf("expected since %s, got %s", expectedSince, *dateRange.Since)
	}
	if expectedUntil != *dateRange.Until {
		t.Errorf("expected until %s, got %s", expectedUntil, *dateRange.Until)
	}
}

func TestGetCommitFields(t *testing.T) {
	expectedFields := []string{"repository_id", "hash", "author_email", "author_name", "message", "num_files", "addition_loc", "deletion_loc", "num_parents", "total_loc", "year", "month", "day", "hour", "commit_time_stamp"}
	expectedString := strings.Join(expectedFields, ",")
	getFields := strings.Join(getCommitFields(), ",")
	if expectedString != getFields {
		t.Errorf("expected: %s, got: %s", expectedString, getFields)
	}
}

func TestGetFileStatFields(t *testing.T) {
	expectedFields := []string{"repository_id", "hash", "author_email", "author_name", "file_name", "addition_loc", "deletion_loc", "year", "month", "day", "hour", "commit_time_stamp"}
	expectedString := strings.Join(expectedFields, ",")
	getFields := strings.Join(getFileStatFields(), ",")
	if expectedString != getFields {
		t.Errorf("expected: %s, got: %s", expectedString, getFields)
	}
}
