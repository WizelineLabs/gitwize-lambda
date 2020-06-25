package main

import (
	"database/sql"
	"encoding/json"
	lambda2 "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"gitwize-lambda/db"
	"gitwize-lambda/gogit"
	"gitwize-lambda/utils"
	"log"
)

type awsLambda interface {
	Trigger(payload interface{}, funcName string, awsRegion string) error
}

type lambdaClient struct{}

func (t lambdaClient) Trigger(payloadValues interface{}, funcName string, awsRegion string) error {
	payload, err := json.Marshal(payloadValues)
	if err != nil {
		return err
	}
	mySession := session.Must(session.NewSession())
	svc := lambda.New(mySession, aws.NewConfig().WithRegion(awsRegion))
	input := &lambda.InvokeInput{
		InvocationType: aws.String("Event"),
		FunctionName:   aws.String(funcName),
		Payload:        payload,
		LogType:        aws.String("Tail"),
	}
	_, err = svc.Invoke(input)
	return err
}

type dbInterface interface {
	GetAllRepoRows(fields []string) *sql.Rows
}

func updateAllRepos(lbd awsLambda, awsRegion string, mydb dbInterface) {
	fields := []string{"id", "name", "url", "password"}
	rows := mydb.GetAllRepoRows(fields)

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
			payload := gogit.RepoPayload{
				RepoID:   id,
				URL:      url,
				RepoName: name,
				RepoPass: password,
				Branch:   "",
			}
			err := lbd.Trigger(payload, utils.GetUpdateOneRepoFuncName(), awsRegion)
			if err != nil {
				log.Println("ERR", err)
			} else {
				count++
			}
		}
	}
	log.Println("Completed update ", count, "repositories")
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler() (string, error) {
	updateAllRepos(lambdaClient{}, "ap-southeast-1", db.NewCommonOps())
	return "Update all repositories completed", nil
}

func main() {
	lambda2.Start(Handler)
}
