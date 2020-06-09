package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	lbd "github.com/aws/aws-lambda-go/lambda"

	"database/sql"
	"github.com/GitWize/gitwize-lambda/db"
	"github.com/GitWize/gitwize-lambda/gogit"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
)

const (
	updateOneRepoFuncName = "gitwize-lambda-dev-update_one_repo"
	awsRegion             = "ap-southeast-1"
)

// Ref: https://docs.aws.amazon.com/sdk-for-go/api/aws/
// Ref: https://godoc.org/github.com/aws/aws-sdk-go/service/lambda#
func updateAllRepos() {
	log.Println("Start updating all repositories")
	conn := db.SQLDBConn()
	defer conn.Close()

	rows, _ := conn.Query("SELECT id, name, url, password FROM repository")

	var id int
	var name, url string
	password := sql.NullString{
		String: "",
		Valid:  false,
	}
	if rows == nil {
		log.Printf("[WARN] No repositories found")
		return
	}
	count := 0
	for rows.Next() {
		count++
		err := rows.Scan(&id, &name, &url, &password)
		if err != nil {
			log.Panicln(err)
		} else {
			p := gogit.RepoEvent{
				RepoID:   id,
				URL:      url,
				RepoName: name,
				RepoPass: password.String,
				Branch:   "",
			}
			payload, err := json.Marshal(p)
			if err != nil {
				log.Panicln(err)
			}

			mySession := session.Must(session.NewSession())
			svc := lambda.New(mySession, aws.NewConfig().WithRegion(awsRegion))

			input := &lambda.InvokeInput{
				InvocationType: aws.String("Event"),
				FunctionName:   aws.String(updateOneRepoFuncName),
				Payload:        payload,
				LogType:        aws.String("Tail"),
			}

			_, err = svc.Invoke(input)
			if err != nil {
				log.Println("ERR invoke lambda", err)
			}
		}
	}
	log.Println("Completed update ", count, "repositories")
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