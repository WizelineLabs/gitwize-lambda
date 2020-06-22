package utils

import (
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
