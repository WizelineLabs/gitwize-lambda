package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
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

// Ref: https://docs.aws.amazon.com/sdk-for-go/api/aws/
// Ref: https://godoc.org/github.com/aws/aws-sdk-go/service/lambda#
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
		}
	}
	log.Println("Completed trigger update ", count, "repositories")
}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	updateAllRepos()

	var buf bytes.Buffer
	body, _ := json.Marshal(map[string]interface{}{
		"message": "update all repositories completed",
	})
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "update-all-repos-handler",
		},
	}
	return resp, nil
}

func main() {
	lbd.Start(Handler)
}
