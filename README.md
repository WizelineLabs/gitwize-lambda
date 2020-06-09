### gitwize-lambda
- lambda functions to retrieve data for gitwize (commit data, PR data, file changes...)


### install pre-commit hook (linter, auto format)
`pre-commit install`


### build and deploy from local
- install serverless framework
https://www.serverless.com/framework/docs/providers/aws/guide/quick-start/

- install aws cli and config with aws credentials (user in aws iam `lambda`) and config region `ap-southeast-1`

- build and deploy: `make && sls deploy`
