package main

import (
	"encoding/json"
	lbd "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"gitwize-lambda/db"
	"gitwize-lambda/gogit"
	"gitwize-lambda/utils"
	"log"
)

const (
	awsRegion = "ap-southeast-1"
)

func triggerLambda(p interface{}, functionName string) {
	payload, err := json.Marshal(p)
	if err != nil {
		log.Println("ERR", err)
	}

	mySession := session.Must(session.NewSession())
	svc := lambda.New(mySession, aws.NewConfig().WithRegion(awsRegion))

	input := &lambda.InvokeInput{
		InvocationType: aws.String("Event"),
		FunctionName:   aws.String(functionName),
		Payload:        payload,
		LogType:        aws.String("Tail"),
	}

	_, err = svc.Invoke(input)
	if err != nil {
		log.Println("ERR invoke lambda", err)
	}
}

func updateAllRepos() {
	conn := db.SQLDBConn()
	defer conn.Close()

	fields := []string{"id", "name", "url", "password"}
	rows := db.GetAllRepoRows(fields)

	var id int
	var name, url, password string

	if rows == nil {
		log.Printf("[WARN] No repositories found")
		return
	}

	count := 0

	for rows.Next() {
		err := rows.Scan(&id, &name, &url, &password)
		if err != nil {
			log.Println("ERR", err)
		} else {
			count++
			payload := gogit.RepoPayload{
				RepoID:   id,
				URL:      url,
				RepoName: name,
				RepoPass: password,
				Branch:   "",
			}
			triggerLambda(payload, utils.GetUpdateOneRepoFuncName())
			db.UpdateRepoLastUpdated(id)
		}
	}
	log.Println("Completed trigger update ", count, "repositories")
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler() (string, error) {
	updateAllRepos()
	return "update all repositories triggered", nil
}

func main() {
	lbd.Start(Handler)
}
