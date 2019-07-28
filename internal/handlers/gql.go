package handlers

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/cmelgarejo/go-gql-server/internal/gql"
	"github.com/cmelgarejo/go-gql-server/internal/gql/resolvers"
	"github.com/cmelgarejo/go-gql-server/internal/orm"
	"github.com/gin-gonic/gin"
)

// GraphqlHandler defines the GQLGen GraphQL server handler
func GraphqlHandler(orm *orm.ORM) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	c := gql.Config{
		Resolvers: &resolvers.Resolver{
			ORM: orm, // pass in the ORM instance in the resolvers to be used
		},
	}

	h := handler.GraphQL(gql.NewExecutableSchema(c))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// PlaygroundHandler defines a handler to expose the Playground
func PlaygroundHandler(path string) gin.HandlerFunc {
	h := handler.Playground("Go GraphQL Server", path)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
