package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"gitwize-lambda/gogit"
	"gitwize-lambda/utils"
	"log"
)

// Handler lambda function handler
func Handler() (string, error) {
	token := utils.GetAccessToken("")
	repoName := "mockRepo"
	path := "/tmp/" + repoName
	url := "https://github.com/sang-d/mock-repo"
	gogit.GetRepo(repoName, url, token)

	data, err := utils.ExecuteCommand(path, "git", "log", "-a", "--before=2020-10-10", "--after=2019-10-10")
	if err != nil {
		log.Panicln(err)
	}
	log.Println("ExecuteCommand \n", string(data))
	return "completed", nil
}

func main() {
	lambda.Start(Handler)
}
