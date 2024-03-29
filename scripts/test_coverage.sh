#!/usr/bin/env bash

export GITWIZE_INTEGRATION_TEST="TRUE"

go test ./...  -coverprofile cover.out
ret=$?
if [[ "${ret}" -gt 0 ]]
then
    exit 1
fi

go tool cover -func cover.out | grep -o '[^,]\+$' | grep total |  awk '{print substr($3, 1, length($3)-1)}' > percentage.txt

percentage=$(cat percentage.txt)

echo "total coverage ${percentage}%"

if [[ "$(echo "${percentage} < ${GITWIZE_TEST_COVERAGE}" | bc)" -ne 0 ]]
then
    echo "Test coverage failed. Expected ${GITWIZE_TEST_COVERAGE}% , got ${percentage}%."
    exit 1
fi
