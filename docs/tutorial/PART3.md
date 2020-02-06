# Creating an opinionated Go GQL Server - Part 3

This is part of a series, not suited for beginners, but welcome too! You can
check it out from the start

- [Part 1](PART1.md)
- [Part 2](PART2.md)

And review how to get up to this part, as always all the code is at github for
[part 3](https://github.com/cmelgarejo/go-gql-server/tree/tutorial/part-3)

---

Today we'll add GORM and a database to our project, and the RDBMS of my choice
is [PostgreSQL](https://www.postgresql.org) You can use MySQL, SQLite, whatever
you like as long as GORM supports it. And, also you can check out
[Part 2](PART2.md) in case you've missed it. I will go even on a faster pace in
this one, check out the comments in the code for more info

## Adding GORM into the project

Let's add GORM to our project:

> `$ go get -u github.com/jinzhu/gorm`

So, we need migrations, let's add [Gomigrate](https://gopkg.in/gormigrate.v1)

> `$ go get gopkg.in/gormigrate.v1`

I don't want the id's to be as easily discoverable as `int`s [1,2,3] right now
anyone could just send a mutation on `deleteUser(userId: 1)`, this should add
another layer on complexity to reduce the attack surface on our silly API.

Add UUID package from `go-frs`:

> `$ go get github.com/gofrs/uuid`

And we'll setup the db models to use UUID as primary keys.

## Diving into GORM

Now, let's setup `GORM` into our project, with:

- Migrations
- Models

First lets create an entry point for everything DB, where we can initialize and
setup database connection, let's put this code at `internal/orm/main.go`

```go
// Package orm provides `GORM` helpers for the creation, migration and access
// on the project's database
package orm

import (
    "github.com/cmelgarejo/go-gql-server/internal/logger"

    "github.com/cmelgarejo/go-gql-server/internal/orm/migration"

    "github.com/cmelgarejo/go-gql-server/pkg/utils"
    //Imports the database dialect of choice
    _ "github.com/jinzhu/gorm/dialects/postgres"

    "github.com/jinzhu/gorm"
)

var autoMigrate, logMode, seedDB bool
var dsn, dialect string

// ORM struct to holds the gorm pointer to db
type ORM struct {
    DB *gorm.DB
}

func init() {
    dialect = utils.MustGet("GORM_DIALECT")
    dsn = utils.MustGet("GORM_CONNECTION_DSN")
    seedDB = utils.MustGetBool("GORM_SEED_DB")
    logMode = utils.MustGetBool("GORM_LOGMODE")
    autoMigrate = utils.MustGetBool("GORM_AUTOMIGRATE")
}

// Factory creates a db connection with the selected dialect and connection string
func Factory() (*ORM, error) {
    db, err := gorm.Open(dialect, dsn)
    if err != nil {
        logger.Panic("[ORM] err: ", err)
    }
    orm := &ORM{
        DB: db,
    }
    // Log every SQL command on dev, @prod: this should be disabled?
    db.LogMode(logMode)
    // Automigrate tables
    if autoMigrate {
        err = migration.ServiceAutoMigration(orm.DB)
    }
    logger.Info("[ORM] Database connection initialized.")
    return orm, err
}
```

You might've noticed `internal/logger` log package, we'll dive into that in the
[bonus](#extra:-bonus-section) section, for now, you can even use `import "log"` package :)

### Migration

The migration service called in `migration.ServiceAutoMigration` that we'll
save in `internal/ocm/migration/main.go`:

```go
package migration

import (
    "fmt"

    "github.com/cmelgarejo/go-gql-server/internal/logger"

    "github.com/cmelgarejo/go-gql-server/internal/orm/migration/jobs"
    "github.com/cmelgarejo/go-gql-server/internal/orm/models"
    "github.com/jinzhu/gorm"
    "gopkg.in/gormigrate.v1"
)

func updateMigration(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
    ).Error
}

// ServiceAutoMigration migrates all the tables and modifications to the connected source
func ServiceAutoMigration(db *gorm.DB) error {
    // Keep a list of migrations here
    m := gormigrate.New(db, gormigrate.DefaultOptions, nil)
    m.InitSchema(func(db *gorm.DB) error {
        logger.Info("[Migration.InitSchema] Initializing database schema")
        switch db.Dialect().GetName() {
        case "postgres":
            // Let's create the UUID extension, the user has to ahve superuser
            // permission for now
            db.Exec("create extension \"uuid-ossp\";")
        }
        if err := updateMigration(db); err != nil {
            return fmt.Errorf("[Migration.InitSchema]: %v", err)
        }
        // Add more jobs, etc here
        return nil
    })
    m.Migrate()

    if err := updateMigration(db); err != nil {
        return err
    }
    m = gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
        jobs.SeedUsers,
    })
    return m.Migrate()
}
```

Now first we see a couple of packages more dependant on migration,
`internal/orm/migration/jobs` and `internal/orm/models`

### Defining models

Let's start by setting up the models that can be added in the `updateMigration`
func, so to define `user.go` but first, I have a special need for all my models,
I want some of them to be soft _deletable_ (not removed from the table itself)
and some can be utterly destroyed, since we are using GORM we can create structs
by using `gorm.Model` struct, or define our own, _we'll do just that_, I
need to be changed, theres another thing, the `gorm.Model` also uses
autoincremented numeric id's, but I like to complicate things so I want to
use UUID's as primary keys :D

Let's create a `internal/orm/models/base.go` model file with 2 versions of the
base model, one with `hard` delete and the other `soft` delete:

```go
package models

import (
  "time"

  "github.com/gofrs/uuid"
  "github.com/jinzhu/gorm"
)

// BaseModel defines the common columns that all db structs should hold, usually
// db structs based on this have no soft delete
type BaseModel struct {
  // ID should use uuid_generate_v4() for the pk's
  ID        uuid.UUID  `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
  CreatedAt time.Time  `gorm:"index;not null;default:CURRENT_TIMESTAMP"` // (My|Postgre)SQL
  UpdatedAt *time.Time `gorm:"index"`
}

// BaseModelSoftDelete defines the common columns that all db structs should
// hold, usually. This struct also defines the fields for GORM triggers to
// detect the entity should soft delete
type BaseModelSoftDelete struct {
  BaseModel

  DeletedAt *time.Time `sql:"index"`
}
```

Now, with these base models we can define any other one we need for our server,
so here is `internal/orm/models/user.go`:

```go
package models

// User defines a user for the app
type User struct {
  BaseModelSoftDelete // We don't to actually delete the users, maybe audit if we want to hard delete them? or wait x days to purge from the table, also

  Email               string  `gorm:"not null;unique_index"`
  UserID              *string // External user ID
  Name                *string
  NickName            *string
  FirstName           *string
  LastName            *string
  Location            *string `gorm:"size:1048"`
  Description         *string `gorm:"size:1048"`
}
```

Since we're going to use a external service for the authentication via `Goth`
and use it for the Authentication flow, no password column for now.

All right, now with this, we just need to seed the table with at least an user,
we'll use [Gomigrate](https://gopkg.in/gormigrate.v1) pkg, and prepare a
migration job file, let's call it: `internal/orm/migrations/jobs/seed_users.go`

```go
package jobs

import (
    "github.com/cmelgarejo/go-gql-server/internal/orm/models"
    "github.com/jinzhu/gorm"
    "gopkg.in/gormigrate.v1"
)

var (
    uname                    = "Test User"
    fname                    = "Test"
    lname                    = "User"
    nname                    = "Foo Bar"
    description              = "This is the first user ever!"
    location                 = "His house, maybe? Wouldn't know"
    firstUser   *models.User = &models.User{
        Email:       "test@test.com",
        Name:        &uname,
        FirstName:   &fname,
        LastName:    &lname,
        NickName:    &nname,
        Description: &description,
        Location:    &location,
    }
)

// SeedUsers inserts the first users
var SeedUsers *gormigrate.Migration = &gormigrate.Migration{
    ID: "SEED_USERS",
    Migrate: func(db *gorm.DB) error {
        return db.Create(&firstUser).Error
    },
    Rollback: func(db *gorm.DB) error {
        return db.Delete(&firstUser).Error
    },
}
```

Now that we have everything set-up we can modify out `cmd/gql-server/main.go`

```go
package main

import (
  "github.com/cmelgarejo/go-gql-server/internal/logger"

  "github.com/cmelgarejo/go-gql-server/internal/orm"
  "github.com/cmelgarejo/go-gql-server/pkg/server"
)

func main() {
  // Create a new ORM instance to send it to our
  orm, err := orm.Factory()
  if err != nil {
    logger.Panic(err)
  }
  // Send: ORM instance
  server.Run(orm)
}
```

Also `pkg/server/main.go`:

```go
package server

import (
  "github.com/cmelgarejo/go-gql-server/internal/logger"

  "github.com/cmelgarejo/go-gql-server/internal/orm"

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
func Run(orm *orm.ORM) {
  logger.Info("GORM_CONNECTION_DSN: ", utils.MustGet("GORM_CONNECTION_DSN"))

  endpoint := "http://" + host + ":" + port

  r := gin.Default()
  // Handlers
  // Simple keep-alive/ping handler
  r.GET("/ping", handlers.Ping())

  // GraphQL handlers
  // Playground handler
  if isPgEnabled {
    r.GET(gqlPgPath, handlers.PlaygroundHandler(gqlPath))
    logger.Info("GraphQL Playground @ " + endpoint + gqlPgPath)
  }
  // Pass in the ORM instance to the GraphqlHandler
  r.POST(gqlPath, handlers.GraphqlHandler(orm))
  logger.Info("GraphQL @ " + endpoint + gqlPath)

  // Run the server
  // Inform the user where the server is listening
  logger.Info("Running @ " + endpoint)
  // Print out and exit(1) to the OS if the server cannot run
  logger.Fatal(r.Run(host + ":" + port))
}
```

And the GraphQL Handler should also receive the ORM instance, so we can use the
database connection in the resolvers `internal/handlers/gql.go`:

```go
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
```

And by this point lets modify the `internal/gql/resolvers/main.go` for the
`Resolver` struct to use the ORM instance:

```go
package resolvers

import (
  "github.com/cmelgarejo/go-gql-server/internal/gql"
  "github.com/cmelgarejo/go-gql-server/internal/orm"
)

// Resolver is a modifiable struct that can be used to pass on properties used
// in the resolvers, such as DB access
type Resolver struct {
  ORM *orm.ORM
}

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

Oh right! Lets modify our `internal/schemas.graphql` to reflect the database
model:

```gql
scalar Time
# Types
type User {
  id: ID!
  email: String!
  userId: String
  name: String
  firstName: String
  lastName: String
  nickName: String
  description: String
  location: String
  createdAt: Time!
  updatedAt: Time
}

# Input Types
input UserInput {
  email: String
  userId: String
  displayName: String
  name: String
  firstName: String
  lastName: String
  nickName: String
  description: String
  location: String
}

# List Types
type Users {
  count: Int # You want to return count for a grid for example
  list: [User!]! # that is why we need to specify the users object this way
}

# Define mutations here
type Mutation {
  createUser(input: UserInput!): User!
  updateUser(id: ID!, input: UserInput!): User!
  deleteUser(id: ID!): Boolean!
}

# Define queries here
type Query {
  users(id: ID): Users!
}
```

You've might noticed the way that I'm returning the `users` query. That's
because we might need to return the `count` of records to a grid for example,
and I want to filter and search in the future, and that's where `list` comes in

And, there's no single user query, why is that? I want to `KISS` always (not you,
stop looking at me like that!) that's why we'll use the same query to retrieve
a single specific record or a bunch of them.

Now, let's modify out `scripts/gqlgen.sh` with an option to regenerate the
resolver functions at `internal/gql/resolvers/generated/resolver.go` so you can
copy any new query or mutation resolver into a file of its own

```bash
#!/bin/bash
printf "\nRegenerating gqlgen files\n"
# Optional, delete the resolver to regenerate, only if there are new queries
# or mutations, if you are just changing the input or type definition and
# doesn't impact the resolvers definitions, no need to do it
while [[ "$#" -gt 0 ]]; do case $1 in
  -r|--resolvers)
    rm -f internal/gql/resolvers/generated/resolver.go
  ;;
  *) echo "Unknown parameter passed: $1"; exit 1;;
esac; shift; done

time go run -v github.com/99designs/gqlgen
printf "\nDone.\n\n"
```

Take from `internal/gql/resolvers/generated/resolver.go` the parts needed, the
`func`'s that have a `"not implemented"` panic on them, and delete everything
except the `package resolvers` line, or else you will have impacting `func`
definitions.

So at last we can create the specific resolver for the users, so a new file with
all the `internal/gql/resolvers/users.go`:

```go
package resolvers

import (
  "context"

  "github.com/cmelgarejo/go-gql-server/internal/logger"

  "github.com/cmelgarejo/go-gql-server/internal/gql/models"
  tf "github.com/cmelgarejo/go-gql-server/internal/gql/resolvers/transformations"
  dbm "github.com/cmelgarejo/go-gql-server/internal/orm/models"
)

// CreateUser creates a record
func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
  return userCreateUpdate(r, input, false)
}

// UpdateUser updates a record
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input models.UserInput) (*models.User, error) {
  return userCreateUpdate(r, input, true, id)
}

// DeleteUser deletes a record
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
  return userDelete(r, id)
}

// Users lists records
func (r *queryResolver) Users(ctx context.Context, id *string) (*models.Users, error) {
  return userList(r, id)
}

// ## Helper functions

func userCreateUpdate(r *mutationResolver, input models.UserInput, update bool, ids ...string) (*models.User, error) {
  dbo, err := tf.GQLInputUserToDBUser(&input, update, ids...)
  if err != nil {
    return nil, err
  }
  // Create scoped clean db interface
  db := r.ORM.DB.Begin()
  if !update {
    db = db.Create(dbo).First(dbo) // Create the user
  } else {
    db = db.Model(&dbo).Update(dbo).First(dbo) // Or update it
  }
  gql, err := tf.DBUserToGQLUser(dbo)
  if err != nil {
    db.RollbackUnlessCommitted()
    return nil, err
  }
  db = db.Commit()
  return gql, db.Error
}

func userDelete(r *mutationResolver, id string) (bool, error) {
  return false, nil
}

func userList(r *queryResolver, id *string) (*models.Users, error) {
  entity := "users"
  whereID := "id = ?"
  record := &models.Users{}
  dbRecords := []*dbm.User{}
  db := r.ORM.DB
  if id != nil {
    db = db.Where(whereID, *id)
  }
  db = db.Find(&dbRecords).Count(&record.Count)
  for _, dbRec := range dbRecords {
    if rec, err := tf.DBUserToGQLUser(dbRec); err != nil {
      logger.Errorfn(entity, err)
    } else {
      record.List = append(record.List, rec)
    }
  }
  return record, db.Error
}
```

You sure will have wondered what `tf` package is: Remember the `trasformations`
file and folders I mentioned in Part 2? Ok, here is where we need to
_transform_ from the GQL Input to the database `user` struct to be easily saved
in the database and convert back from the database to the GQL return struct of
the resolver:

```go
package transformations

import (
  "errors"

  gql "github.com/cmelgarejo/go-gql-server/internal/gql/models"
  dbm "github.com/cmelgarejo/go-gql-server/internal/orm/models"
  "github.com/gofrs/uuid"
)

// DBUserToGQLUser transforms [user] db input to gql type
func DBUserToGQLUser(i *dbm.User) (o *gql.User, err error) {
  o = &gql.User{
    ID:          i.ID.String(),
    Email:       i.Email,
    UserID:      i.UserID,
    Name:        i.Name,
    FirstName:   i.FirstName,
    LastName:    i.LastName,
    NickName:    i.NickName,
    Description: i.Description,
    Location:    i.Location,
    CreatedAt:   i.CreatedAt,
    UpdatedAt:   i.UpdatedAt,
  }
  return o, err
}

// GQLInputUserToDBUser transforms [user] gql input to db model
func GQLInputUserToDBUser(i *gql.UserInput, update bool, ids ...string) (o *dbm.User, err error) {
  o = &dbm.User{
    UserID:      i.UserID,
    Name:        i.Name,
    FirstName:   i.FirstName,
    LastName:    i.LastName,
    NickName:    i.NickName,
    Description: i.Description,
    Location:    i.Location,
  }
  if i.Email == nil && !update {
    return nil, errors.New("field [email] is required")
  }
  if i.Email != nil {
    o.Email = *i.Email
  }
  if len(ids) > 0 {
    updID, err := uuid.FromString(ids[0])
    if err != nil {
      return nil, err
    }
    o.ID = updID
  }
  return o, err
}
```

Getting back to the `users.go` resolver, I've separated the helper `func`'s from
the resolver to have them the smallest and more readable way to do this, passing
the resolver to the helper functions will make more sense than just passing the
ORM in the next [Part 4](PART4.md) when we'll need to plug-in authentication and
permissions to the resolvers.

Hey, now we can query and mutate our users from the database!

![gql is showing db data!](images/p3-p1.png?raw=true "Hey there, sexy! (2)")

And a specific user from the database:

![gql is showing db data!](images/p3-p2.png?raw=true "Hey there, sexy! (3)")

We can always create a new one:

![gql is creating db data!](images/p3-p3.png?raw=true "Hey there, sexy! (4)")

And update it:

![gql is updating db data!](images/p3-p4.png?raw=true "Hey there, sexy! (5)")

Okay, and as promised let's enhance our little project a bit more by adding a
couple of packages next.

## Extra: Bonus section

### Logger package

The included `"log"` package could always be enough for us, but I found this
well structured package called [Logrus](http://github.com/sirupsen/logrus) and
created a wrapper for it to be used throughout the project, create the file
`internal/logger/main.go`:

```go
package logger

import (
  "github.com/sirupsen/logrus"
)

var logger *StandardLogger

func init() {
  logger = NewLogger()
}

// Event stores messages to log later, from our standard interface
type Event struct {
  id      int
  message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
  *logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
  var baseLogger = logrus.New()

  var standardLogger = &StandardLogger{baseLogger}

  standardLogger.Formatter = &logrus.TextFormatter{
    FullTimestamp: true,
  }
  // We could transform the errors into a JSON format, for external log SaaS tools such as splunk or logstash
  // standardLogger.Formatter = &logrus.JSONFormatter{
  //   PrettyPrint: true,
  // }

  return standardLogger
}

// Declare variables to store log messages as new Events
var (
  invalidArgMessage      = Event{1, "Invalid arg: %s"}
  invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
  missingArgMessage      = Event{3, "Missing arg: %s"}
)

// Errorfn Log errors with format
func Errorfn(fn string, err error) {
  logger.Errorf("[%s]: %v", fn, err)
}

// InvalidArg is a standard error message
func InvalidArg(argumentName string) {
  logger.Errorf(invalidArgMessage.message, argumentName)
}

// InvalidArgValue is a standard error message
func InvalidArgValue(argumentName string, argumentValue string) {
  logger.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
}

// MissingArg is a standard error message
func MissingArg(argumentName string) {
  logger.Errorf(missingArgMessage.message, argumentName)
}

// Info Log
func Info(args ...interface{}) {
  logger.Infoln(args...)
}

// Infof Log
func Infof(format string, args ...interface{}) {
  logger.Infof(format, args...)
}

// Warn Log
func Warn(args ...interface{}) {
  logger.Warnln(args...)
}

// Warnf Log
func Warnf(format string, args ...interface{}) {
  logger.Warnf(format, args...)
}

// Panic Log
func Panic(args ...interface{}) {
  logger.Panicln(args...)
}

// Panicf Log
func Panicf(format string, args ...interface{}) {
  logger.Panicf(format, args...)
}

// Error Log
func Error(args ...interface{}) {
  logger.Errorln(args...)
}

// Errorf Log
func Errorf(format string, args ...interface{}) {
  logger.Errorf(format, args...)
}

// Fatal Log
func Fatal(args ...interface{}) {
  logger.Fatalln(args...)
}

// Fatalf Log
func Fatalf(format string, args ...interface{}) {
  logger.Fatalf(format, args...)
}
```

Now we have a extendable logger package with nice features

![logger stuff](images/p3-p5.png?raw=true "Logging with Logrus")

### Adding server reloading on change with Realize

If you are annoyed on how every change you make you have to restart you server,
coming from developing quite a while in node.js, I really missed `nodemon`,
luckily we have something even better in Go:

> `$ go get github.com/tockins/realize`

After this you'll have `realize` executable in `$GOPATH/bin/realize` up for use,
so, create a `.realize.yml` file (or run `$GOPATH/bin/realize init` in the root
of our project, to create it interactively, here's my _yml_ file:

```yml
settings:
  files:
    outputs:
      status: false
      path: ""
      name: .r.outputs.log
    logs:
      status: false
      path: ""
      name: .r.logs.log
    errors:
      status: false
      path: ""
      name: .r.errors.log
  legacy:
    force: false
    interval: 0s
schema:
  - name: gql-server
    path: cmd/gql-server
    commands:
      install:
        status: false
      run:
        status: true
    watcher:
      extensions:
        - go
      paths:
        - ../../
      ignore:
        paths:
          - .git
          - .realize
          - .vscode
          - vendor
```

This should be enough for running in a script we're making next, copy the
`scripts/run.sh` to `scripts/dev-run.sh` and modify it like so:

```bash
#!/bin/sh
app="gql-server"
src="$srcPath/$app/$pkgFile"

printf "\nStart running: $app\n"
# Set all ENV vars for the server to run
export $(grep -v '^#' .env | xargs)
time /$GOPATH/bin/realize start run
# This should unset all the ENV vars, just in case.
unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)
printf "\nStopped running: $app\n\n"
```

And also, change `run.sh`to:

```bash
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
```

Well, now you can run our project with `.scripts/dev-run.sh` and have it reload
every time you make a change in your go files! Neat, right?

Again, like the former parts, all the code is available in the
[repository here](https://github.com/cmelgarejo/go-gql-server/tree/tutorial/part-3)! If you have questions, critiques and comments go ahead and let's learn more together!

And...

> Sorry for the long and unwind post, here's a gopher potato: ![gopher potato](https://golang.org/doc/gopher/gophercolor.png)
