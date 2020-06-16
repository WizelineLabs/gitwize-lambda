package main

import (
	"gitwize-lambda/utils"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	utils.SetupIntegrationTest()
	os.Args = []string{"first-arg", "1", "integration-test-mock-repo", "https://github.com/sang-d/mock-repo"}
	main()
}
