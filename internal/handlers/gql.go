package handlers

import (
	"github.com/99designs/gqlgen/handler"
	"github.com/cmelgarejo/go-gql-server/internal/gql"
	"github.com/cmelgarejo/go-gql-server/internal/gql/resolvers"
	"github.com/cmelgarejo/go-gql-server/internal/orm"
	"github.com/cmelgarejo/go-gql-server/pkg/utils"
	"github.com/gin-gonic/gin"
)

// GraphqlHandler defines the GQLGen GraphQL server handler
func GraphqlHandler(orm *orm.ORM, gqlConfig *utils.GQLConfig) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	c := gql.Config{
		Resolvers: &resolvers.Resolver{
			ORM: orm, // pass in the ORM instance in the resolvers to be used
		},
	}

	setProjectComplexity(&c)
	h := handler.GraphQL(gql.NewExecutableSchema(c), handler.ComplexityLimit(gqlConfig.ComplexityLimit))

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

func setProjectComplexity(c *gql.Config) {
	countComplexity := func(childComplexity int, limit *int, offset *int) int {
		return *limit
	}
	fixedComplexity := func(childComplexity int) int {
		return 100
	}
	// c.Complexity.Query.Users = func(childComplexity int, id *string, filters []*models.QueryFilter, limit *int, offset *int, orderBy *string, sortDirection *string) int {
	// 	return *limit
	// }
	c.Complexity.User.CreatedBy = fixedComplexity
	c.Complexity.User.UpdatedBy = fixedComplexity
	c.Complexity.User.Profiles = countComplexity
}
