#!/bin/sh
buildPath="build"
app="gql-server"
program="$buildPath/$app"

printf "\nStart app: $app\n"
# Set all ENV vars for the program to run
export $(grep -v '^#' .env | xargs)
time ./$program
# This should unset all the ENV vars, just in case.
unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)
printf "\nStopped app: $app\n\n"
