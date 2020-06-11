#!/usr/bin/env bash

export APP_STAGE=dev
export PATTERN="v[0-9]+(\.[0-9]+).*"
if [[ ${CIRCLE_TAG} =~  ${PATTERN} ]]; then
    export APP_STAGE=qa
    export DB_CONN_STRING=${DB_CONN_STRING_QA}
    export DEFAULT_GITHUB_TOKEN=${DEFAULT_GITHUB_TOKEN_QA}
    export CYPHER_PASS_PHASE=${CYPHER_PASS_PHASE_QA}
    export USE_DEFAULT_API_TOKEN=${USE_DEFAULT_API_TOKEN_QA}
    echo "deploy to QA"
fi

sls deploy --stage ${APP_STAGE} --verbose
