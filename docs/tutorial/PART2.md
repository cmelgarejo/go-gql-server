# Creating an opinionated Go GQL Server - Part 2

-[Part 1](PART1.md)

We'll be now add GQLGen's generated server into our project and start `gql`ing
away! And, we are going to move much faster than the [Part 1](PART1.md)

## Adding GQLGen into the project

Now, we can use **GQLGen** to initialize our server's GQLGen generated files
`go run github.com/99designs/gqlgen init`

This will initialize the gqlgen server, we have to make modification to fit our
project layout though, move a couple of files and edit `gqlgen.yml` file

- `generated.go`(file that holds the gql gen's generated graphql server code)
- `resolver.go` (holds the resolvers for the queries and mutations)
- `schema.graphql` (example schema that gqlgen generates, we'll modify this)
- `server/server.go` (stub server that we'll dump but use the handlers)
- `models_gen.go` (file with the generated models based on the schema.graphql)

Right now we are going to to something maybe counter intuitive, let's delete all
these new files! Except `gqlgen.yml`; trust me, we'll get them back and in a
proper place.

> For more on config for gqlgen: [go here](https://gqlgen.com/config/)

Next, what we need to do first is modify the `gqlgen.yml` file like so:

```yml
# go-gql-server gqlgen.yml file
# Refer to https://gqlgen.com/config/
# for detailed .gqlgen.yml documentation.

schema:
  - internal/gql/schemas/schema.graphql
# Let gqlgen know where to put the generated server
exec:
  filename: internal/gql/generated.go
  package: gql
# Let gqlgen know where to put the generated models (if any)
model:
  filename: internal/gql/models/generated.go
  package: models
# Let gqlgen know where to put the generated resolvers
resolver:
  filename: internal/gql/resolvers/generated.go
  type: Resolver
  package: resolvers
autobind: []
```

So, I want to automate the generation of the gqlgen files, let's create a script
named:

- `scripts/gqlgen.sh` (remember to `chmod +x` this one too)

```bash
#!/bin/bash
printf "\nRegenerating gqlgen files\n"
rm -f internal/gql/generated.go \
    internal/gql/models/generated.go \
    internal/gql/resolvers/generated.go
time go run -v github.com/99designs/gqlgen $1
printf "\nDone.\n\n"
```

Why am I deleting also the resolvers file? to regenerate it of course, if we are
to change anything on the schema most likely we'll be updating the resolvers too
and, when the project is small is good you have it all in one file, but we'll be
going to sort the resolvers in their own files, and have this file
`internal/gql/resolvers/generated.go` as a temporary file between gqlgen
generations. It's a little tedious but it will payoff later on.

Alright, now we need to define a gql schema file, so that we can use our script
to regenerate the gqlgen files in their rightful place.

> `$ mkdir -p internal/gql/schemas`

Then you can edit a `schema.graphql` file

> `$ vi internal/gql/schemas/schema.graphql`

and paste the following on it:

```graphql
# Types
type User {
  id: ID
  email: String
  userId: String
}

# Input Types
input UserInput {
  email: String
  userId: String
}

# Define mutations here
type Mutation {
  createUser(input: UserInput!): User!
  updateUser(input: UserInput!): User!
  deleteUser(userId: ID!): Boolean!
}

# Define queries here
type Query {
  users(userId: ID): [User]
}
```

Don't worry about the gql specifics, we are going to enhance this schema a lot
more in the next parts, using OpenID specifications.

For now, lets grab `internal/gql/resolvers/generated.go` and edit it a little to
return a mocked user:

```go
package resolvers

import (
    "context"

    "github.com/cmelgarejo/go-gql-server/internal/gql"
    "github.com/cmelgarejo/go-gql-server/internal/gql/models"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() gql.MutationResolver {
    return &mutationResolver{r}
}
func (r *Resolver) Query() gql.QueryResolver {
    return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
    panic("not implemented")
}
func (r *mutationResolver) UpdateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
    panic("not implemented")
}
func (r *mutationResolver) DeleteUser(ctx context.Context, userID string) (bool, error) {
    panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context, userID *string) ([]*models.User, error) {
    records := []*models.User{
        &models.User{
            ID:     "ec17af15-e354-440c-a09f-69715fc8b595",
            Email:  "your@email.com",
            UserID: "UserID-1",
        },
    }
    return records, nil
}
```

Now let's bind the Graphql Server middleware to our server! Create a file in
`internal/handlers/gql.go` an paste this in:

```go
package handlers

import (
    "github.com/99designs/gqlgen/handler"
    "github.com/cmelgarejo/go-gql-server/internal/gql"
    "github.com/cmelgarejo/go-gql-server/internal/gql/resolvers"
    "github.com/gin-gonic/gin"
)

// GraphqlHandler defines the GQLGen GraphQL server handler
func GraphqlHandler() gin.HandlerFunc {
    // NewExecutableSchema and Config are in the generated.go file
    c := gql.Config{
        Resolvers: &resolvers.Resolver{},
    }

    h := handler.GraphQL(gql.NewExecutableSchema(c))

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}

// PlaygroundHandler Defines the Playground handler to expose our playground
func PlaygroundHandler(path string) gin.HandlerFunc {
    h := handler.Playground("Go GraphQL Server", path)
    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}
```

Now we can modify the `pkg/server/main.go` like this:

```go
package server

import (
    "github.com/cmelgarejo/go-gql-server/internal/logger"

    "github.com/cmelgarejo/go-gql-server/internal/handlers"
    "github.com/cmelgarejo/go-gql-server/pkg/utils"
    "github.com/gin-gonic/gin"
)

var host, port, gqlPath, gqlPgPath string
var isPgEnabled bool

func init() {
    host = utils.MustGet("SERVER_HOST")
    port = utils.MustGet("SERVER_PORT")
    gqlPath = utils.MustGet("GQL_SERVER_GRAPHQL_PATH")
    gqlPgPath = utils.MustGet("GQL_SERVER_GRAPHQL_PLAYGROUND_PATH")
    isPgEnabled = utils.MustGetBool("GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED")
}

// Run spins up the server
func Run() {
    endpoint := "http://" + host + ":" + port

    r := gin.Default()

    // Handlers
    // Simple keep-alive/ping handler
    r.GET("/ping", handlers.Ping())

    // GraphQL handlers
    // Playground handler
    if isPgEnabled {
        r.GET(gqlPgPath, handlers.PlaygroundHandler(gqlPath))
        logger.Println("GraphQL Playground @ " + endpoint + gqlPgPath)
    }
    r.POST(gqlPath, handlers.GraphqlHandler())
    logger.Println("GraphQL @ " + endpoint + gqlPath)

    // Run the server
    // Inform the user where the server is listening
    logger.Println("Running @ " + endpoint)
    // Print out and exit(1) to the OS if the server cannot run
    logger.Fatalln(r.Run(host + ":" + port))
}
```

So with the new modifications, we set new ENV variables:

```bash .env
# Web framework config
GIN_MODE=debug
GQL_SERVER_HOST=localhost
GQL_SERVER_PORT=7777
# GQLGen config
GQL_SERVER_GRAPHQL_PATH=/graphql
GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED=true
GQL_SERVER_GRAPHQL_PLAYGROUND_PATH=/
```

that dictate where will our server listen with the graphql handler and serve
queries and mutations we already defined, and let's try and run the program now:

> `$ ./scripts/run/sh`

```bash
$ ./scripts/run.sh

Start running: gql-server
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> github.com/cmelgarejo/go-gql-server/internal/handlers.Ping.func1 (3 handlers)
[GIN-debug] GET    /                         --> github.com/cmelgarejo/go-gql-server/internal/handlers.PlaygroundHandler.func1 (3 handlers)
2019/07/13 23:28:38 GraphQL Playground @ http://localhost:7777/
[GIN-debug] POST   /graphql                  --> github.com/cmelgarejo/go-gql-server/internal/handlers.GraphqlHandler.func1 (3 handlers)
2019/07/13 23:28:38 GraphQL @ http://localhost:7777/graphql
2019/07/13 23:28:38 Running @ http://localhost:7777
[GIN-debug] Listening and serving HTTP on localhost:7777
[GIN] 2019/07/13 - 23:28:39 | 200 |     486.646µs |       127.0.0.1 | GET      /
[GIN] 2019/07/13 - 23:28:40 | 200 |    1.353992ms |       127.0.0.1 | POST     /graphql
```

We can see now GraphQL requests being redirected to its handler! Nice!

Let's navigate to: [http://localhost:7777](http://localhost:7777)

![gql is up!](images/p2-p1.png?raw=true 'Hey there, sexy!')

Now we got ourselves a functioning GQL server, let's try and query...

![show me the users!](images/p2-p2.png?raw=true 'User! mocked one, but COOL!')

Now we see that everything is working as intended, let's move on to better
things now, and refactor the code to be a little more ordered.

## Refactoring code (a.k.a my own personal filename hell)

Like I noted before, I have a contrived, yet effective way to organize
the code for GQLgen, it will help to keep code organized in small files:

- `{entity_plural}.go` where the resolves that are generated with our `gqlgen.sh`
  script will have to be copied individually for each `entity`
  (users, posts, comments, you name it)
- `transformations/{entity_plural}.go` now, this one is interesting, once we add
  **GORM** into our project and with that **database structs**, we'll have to
  _transform_ these GQL Input types into database representation to be stored in
  the _db_ and vice-versa to return queries from _db_ to _GQL_ for example.
- Then I move whatever was generated from
  `internal/gql/resolvers/generated/generated.go` to
  `internal/gql/resolvers/main.go` and trim it from the entity methods we might
  have and just leave it like this:

```go
package resolvers

import (
    "github.com/cmelgarejo/go-gql-server/internal/gql"
)

// Resolver is a modifable struct that can be used to pass on properties used
// in the resolvers, such as DB access
type Resolver struct{}

// Mutation exposes mutation methods
func (r *Resolver) Mutation() gql.MutationResolver {
    return &mutationResolver{r}
}

// Query exposes query methods
func (r *Resolver) Query() gql.QueryResolver {
    return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }
```

This file won't have any more generated changes from now on, we'll change a
couple of things yes, but that will be done in [Part 3](PART3.md)!

I haven't find a way to organize the folder structure into folders without
having to do recursive imports, so the prefixes will help to organize the
filenames for the entities in the `resolvers` folder. My TOC forces me to do
this 😅

> If you have a nicer way to arrange this mess, feel free to open an Issue/PR!

As the last part, all the code is available in the [repository here](https://github.com/cmelgarejo/go-gql-server/tree/tutorial/part-1)!
