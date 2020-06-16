package main

import (
	"gitwize-lambda/utils"
	"testing"
)

func TestMain(t *testing.T) {
	utils.SetupIntegrationTest()
	main() //can disable as it might run quite long ~20s
}
