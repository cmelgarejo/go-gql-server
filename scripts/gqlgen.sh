#!/bin/bash
# Optional, delete the resolver to regenerate, only if there are new queries
# or mutations, if you are just changeing the input or type definition and
# doesn't impact the resolvers definitions, no need to do it
while [[ "$#" -gt 0 ]]; do case $1 in
  -r|--resolver)
    printf "\nRecreating 'internal/gql/resolvers/generated/resolver.go':
    Remember to delete all definitions and take what is needed to another
    file in 'internal/gql/resolvers/'\n"
    rm -f internal/gql/resolvers/generated/resolver.go
  ;;
  *) echo "Unknown parameter provided: $1"; exit 1;;
esac; shift; done
printf "\nRegenerating gqlgen files\n"
time go run -v github.com/99designs/gqlgen
printf "\nDone.\n\n"