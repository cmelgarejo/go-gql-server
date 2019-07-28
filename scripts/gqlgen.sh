#!/bin/bash
printf "\nRegenerating gqlgen files\n"
# Optional, delete the resolver to regenerate, only if there are new queries
# or mutations, if you are just chageing the input or type definition and
# doesn't impact the resolvers definitions, no need to do it
while [[ "$#" -gt 0 ]]; do case $1 in
  -r|--resolvers)
    rm -f internal/gql/resolvers/generated/resolver.go
  ;;
  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

time go run -v github.com/99designs/gqlgen
printf "\nDone.\n\n"