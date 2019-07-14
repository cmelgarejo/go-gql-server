# Go-GQL-Server

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

## Deployment

Use docker, swarm or kubernetes, GCP, AWS, DO, you name it.
