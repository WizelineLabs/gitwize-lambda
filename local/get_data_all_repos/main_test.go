package main

import (
	"gitwize-lambda/utils"
	"testing"
)

func TestMain(t *testing.T) {
	if utils.IntegrationTestEnabled() {
		main()
	}
}
