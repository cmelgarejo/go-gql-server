package handlers

import (
	"time"

	"github.com/99designs/gqlgen/graphql/handler/lru"

	"github.com/99designs/gqlgen/graphql/handler"

	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
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
		Directives: gql.DirectiveRoot{},
		Complexity: gql.ComplexityRoot{},
	}

	// setProjectComplexity(&c)
	// h := handler.GraphQL(gql.NewExecutableSchema(c), handler.ComplexityLimit(gqlConfig.ComplexityLimit))

	srv := handler.New(gql.NewExecutableSchema(c))
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.Use(extension.FixedComplexityLimit(gqlConfig.ComplexityLimit))
	if gqlConfig.IsIntrospectionEnabled {
		srv.Use(extension.Introspection{})
	}
	srv.Use(apollotracing.Tracer{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	return func(c *gin.Context) {
		// h.ServeHTTP(c.Writer, c.Request)
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// PlaygroundHandler defines a handler to expose the Playground
func PlaygroundHandler(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		playground.Handler("Go GraphQL Server", path).ServeHTTP(c.Writer, c.Request)
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
