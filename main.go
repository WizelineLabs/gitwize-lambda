<<<<<<< HEAD
package main
=======
package main

import (
	"github.com/GitWize/gitwize-lambda/gogit"
	"log"
	"os"
)

func main() {
	log.Println("Test function locally")
	url := "https://github.com/go-git/go-git"
	token := os.Getenv("DEFAULT_GITHUB_TOKEN")
	gogit.GetRepo("go-git", url, token)
}
>>>>>>> add functions to update all repos and one repo
