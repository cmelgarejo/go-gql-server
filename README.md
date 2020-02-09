# Go-GQL-Server

[![Go Report Card](https://goreportcard.com/badge/github.com/cmelgarejo/go-gql-server)](https://goreportcard.com/report/github.com/cmelgarejo/go-gql-server)

Opinionated GraphQL server using:

- [Gin-gonic](https://gin-gonic.com) web framework
  - `go get -u github.com/gin-gonic/gin`
- [Goth](https://github.com/markbates/goth) for OAuth2 connections
  - `go get github.com/markbates/goth`
- [GORM](http://gorm.io) as DB ORM
  - `go get -u github.com/jinzhu/gorm`
  - [Gomigrate](https://gopkg.in/gormigrate.v1)
    - `go get gopkg.in/gormigrate.v1`
- [GQLGen](https://gqlgen.com/) for building GraphQL servers without any fuss
  - `go run github.com/99designs/gqlgen init`

## Development with docker

Just run it with `docker-compose`:

`$ docker-compose run dev`

And you'll have your server running with `realize` for your development joy.

## Deployment

Use docker, swarm or kubernetes, GCP, AWS, DO, you name it.

Running `prod.dockerfile` will build a multistaged build that will give you a slim image containing just the gql-server executable.

### With `docker-compose`

`$ docker-compose build prod`

or

`$ docker-compose run prod`

### Build from the `prod.dockerfile`

`docker build -f docker/prod.dockerfile -t go-gql-server.prod ./`
