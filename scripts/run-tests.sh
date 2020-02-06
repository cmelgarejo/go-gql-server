#!/bin/sh
app="gql-server"
program="$buildPath/$app"

printf "\nStart testing: $app\n"
# Set all ENV vars for the server to run
export $(grep -v '^#' .env | xargs)
time go test ./...
# This should unset all the ENV vars, just in case.
unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)
printf "\nRun tests completed: $app\n\n"
