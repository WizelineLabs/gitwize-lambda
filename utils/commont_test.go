package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetUpdateOneRepoFuncName(t *testing.T) {
	expectedName := functionPrefix + GetAppStage() + "-update_one_repo"
	if expectedName != GetUpdateOneRepoFuncName() {
		t.Errorf("Test failed, expected %s, got %s", expectedName, GetUpdateOneRepoFuncName())
	}
}

func TestGetAccessTokenEmptyPass(t *testing.T) {
	expectedAccessToken := os.Getenv("DEFAULT_GITHUB_TOKEN")
	if expectedAccessToken != GetAccessToken("") {
		t.Errorf("Test failed, expected %s, got %s", expectedAccessToken, GetAccessToken(""))
	}
}

func TestGetAppStage(t *testing.T) {
	expected := "dev"
	if expected != GetAppStage() {
		t.Errorf("Test failed, expected %s, got %s", expected, GetAppStage())
	}
}

func TestSetDBConnSingleComponent(t *testing.T) {
	dbConnString := os.Getenv("DB_CONN_STRING")
	os.Setenv("DB_CONN_STRING", "gitwize_user:P@ssword123@(localhost:3306)/gitwize")
	SetDBConnSingleComponent()
	assert.Equal(t, "gitwize_user", os.Getenv("GW_DB_USER"))
	assert.Equal(t, "P@ssword123", os.Getenv("GW_DB_PASSWORD"))
	assert.Equal(t, "localhost", os.Getenv("GW_DB_HOST"))
	assert.Equal(t, "3306", os.Getenv("GW_DB_PORT"))
	assert.Equal(t, "gitwize", os.Getenv("GW_DB_NAME"))
	os.Setenv("DB_CONN_STRING", dbConnString)
}
