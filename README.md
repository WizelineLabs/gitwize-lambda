## Gitwize Lambda
lambda functions to retrieve and process data for gitwize (commit data, PR data, file changes...)


## Pre-setup for development
`pre-commit install`

## Running function on local

#### Env Variables needed:

- `export USE_DEFAULT_API_TOKEN=***`

- `export DEFAULT_GITHUB_TOKEN=***`

- `export DB_CONN_STRING=***`

**Note that you can point DB_CONN_STRING to local/dev/prod directly, so be caution.**

#### Function to get commit & PR data for All repositories:
`go run local/get_data_all_repos/main.go`

#### Function to load metrics for All repositories:
`go run local/load_metric_all_repos/main.go`

#### Function to get commit & PR data for single repo
`go run local/get_data_single_repo/main.go [repo_id] [repo_name] [repo_url] [repo_pass]`

for example:
`go run local/get_data_single_repo/main.go 61 go-git https://github.com/go-git/go-git.git`

## Build and deploy to cloud
- install serverless framework
https://www.serverless.com/framework/docs/providers/aws/guide/quick-start/

- install aws cli and config with aws credentials (user in aws iam `lambda`) and config region `ap-southeast-1`

- build and deploy dev: `make && sls deploy --stage dev`

- build and deploy qa: `make && sls deploy --stage qa`
