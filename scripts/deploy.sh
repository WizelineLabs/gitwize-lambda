#!/usr/bin/env bash

if [[ -z "${APP_STAGE}" ]]; then
    export APP_STAGE=dev
fi

if [[ ${APP_STAGE} == 'qa' ]]; then
    export DB_CONN_STRING=${DB_CONN_STRING_QA}
    export DEFAULT_GITHUB_TOKEN=${DEFAULT_GITHUB_TOKEN_QA}
    export CYPHER_PASS_PHASE=${CYPHER_PASS_PHASE_QA}
    export USE_DEFAULT_API_TOKEN=${USE_DEFAULT_API_TOKEN_QA}
    echo "deploy to QA"
elif [[ ${APP_STAGE} == 'prod' ]]; then
    export DB_CONN_STRING=${DB_CONN_STRING_PROD}
    export DEFAULT_GITHUB_TOKEN=${DEFAULT_GITHUB_TOKEN_PROD}
    export CYPHER_PASS_PHASE=${CYPHER_PASS_PHASE_PROD}
    export USE_DEFAULT_API_TOKEN=${USE_DEFAULT_API_TOKEN_PROD}
    echo "deploy to PROD"
fi

sls deploy --stage ${APP_STAGE} --verbose
