# Creating an opinionated Go GQL Server - Part 1

Let's make an opinionated GraphQL server using:

- [Gin-gonic](https://gin-gonic.com) web framework
- [Goth](https://github.com/markbates/goth) for OAuth2 connections
- [GORM](http://gorm.io) as DB ORM
  - [Gomigrate](https://gopkg.in/gormigrate.v1)
- [GQLGen](https://gqlgen.com/) for building GraphQL servers without any fuss

This assumes you have:

- At least, basic Go knowledge, and have go 1.12+ installed.
- VSCode (preferred) or similar IDE

## Project Setup

We'll follow the Go standard project layout, for this service, take a look at
the specification, it is opinionated, but serves as a good base, and we might
strand from it a little bit from it's guidelines.

Start by creating a directory anywhere we want to:

```bash
> mkdir go-gql-server
> cd go-gql-server
/go-gql-server $
```

Let's create the whole project layout with:

```bash
$ mkdir -p {build,cmd/gql-server,internal/gql-server,pkg,scripts}
# directories are created
```

- `internal/gql-server` will hold all the related files for the gql server
- `cmd/gql-server` will house the `main.go` file for the server, the entrypoint
  that will glue all together.

Since we're using go 1.12+ that will allow to use any directory you want outside
the `$GOPATH/src` path, we want to use `go modules` to initialize our project
with it we have to run:

```bash
$ go mod init github.com/cmelgarejo/go-gql-server # Replace w/your user handle
go: creating new go.mod: module github.com/cmelgarejo/go-gql-server
```

## Coding our web server

### 1. Web framework: Gin

So, now we can start adding packages to our project! Let's start by getting our
web framework: `gin-gonic`

From gin-gonic.com:

> **What is Gin?** Gin is a web framework written in Golang. It features a
> martini-like API with much better performance, up to 40 times faster. If you
> need performance and good productivity, you will love Gin.

```bash
$ go get -u github.com/gin-gonic/gin
go: creating new go.mod: module cmd/gql-server/main.go
```

### 2. Code the web server

So lets start creating the web server, to keep going, create a `main.go` file
in `cmd/gql-server`

```bash
$ vi cmd/gql-server/main.go
# vi ensues, I hope you know how to exit
```

And paste this placeholder code:

```go
package main

import (
    "github.com/cmelgarejo/go-gql-server/internal/logger"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    host := "localhost"
    port := "7777"
    pathGQL := "/graphql"
    r := gin.Default()
    // Setup a route
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, "OK")
    })
    // Inform the user where the server is listening
    logger.Println("Running @ http://" + host + ":" + port + pathGQL)
    // Print out and exit(1) to the OS if the server cannot run
    logger.Fatalln(r.Run(host + ":" + port))
}
```

If you `go run cmd/gql-server/main.go` this code, it will bring up a GIN server
listening to [locahost:7777](http://localhost:7777/ping) and get an _OK_ printed
out in the browser. It works! Now, lets refactor this code using our already
present directory structure, `script`, `internal` and `pkg` folders

- `internal/handlers/ping.go`

```go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// Ping is simple keep-alive/ping handler
func Ping() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.String(http.StatusOK, "OK")
    }
}
```

- `pkg/srv/main.go`

```go
package server

import (
    "github.com/cmelgarejo/go-gql-server/internal/logger"

    "github.com/gin-gonic/gin"
    "github.com/cmelgarejo/go-gql-server/internal/handlers"
)

var HOST, PORT string

func init() {
    HOST = "localhost"
    PORT = "7777"
}

// Run web server
func Run() {
    r := gin.Default()
    // Setup routes
    r.GET("/ping", handlers.Ping())
    logger.Println("Running @ http://" + HOST + ":" + PORT )
    logger.Fatalln(r.Run(HOST + ":" + PORT))
}
```

- `cmd/gql-server/main.go` file can be changed to just this:

```go
package main

import (
    "github.com/cmelgarejo/go-gql-server/pkg/server"
)

func main() {
    server.Run()
}

```

How about it? One line and we have our server running out of the `pkg` folder

Now we can build up the server with a script:

- `scripts/build.sh` (`/bin/sh` because it is more obiquitus)

```bash
#!/bin/sh
srcPath="cmd"
pkgFile="main.go"
outputPath="build"
app="gql-server"
output="$outputPath/$app"
src="$srcPath/$app/$pkgFile"

printf "\nBuilding: $app\n"
time go build -o $output $src
printf "\nBuilt: $app size:"
ls -lah $output | awk '{print $5}'
printf "\nDone building: $app\n\n"
```

And before running, make sure you `chmod +x` it

```bash
$ chmod +x scripts/build.sh
# sets execution permission on file
```

Now we can start building our server, like so:

```bash
> .scripts/build.sh

Building: gql-server

real    0m0.317s
user    0m0.531s
sys     0m0.529s

Built: gql-server size:16M

Done building: gql-server
```

16M standalone server, not bad I think, and this could be the size of it's
docker image! Ok, now onto trying out what has been built:

```bash
$ ./build/gql-server
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> github.com/cmelgarejo/go-gql-server/internal/handlers.Ping.func1 (3 handlers)
2019/07/13 00:00:00 Running @ http://localhost:7777/gql
[GIN-debug] Listening and serving HTTP on :7777
[GIN] 2019/07/13 - 00:00:00 | 200 |      38.936µs |             ::1 | GET      /ping
```

> That's some serious speed, 39 µ(_micro_) seconds.

And we can further improve our server code, and make configurations load from a
`.env` file, let's create some `utils` for our server:

- `pkg/utils/env.go`

```go
package utils

import (
    "github.com/cmelgarejo/go-gql-server/internal/logger"
    "os"
    "strconv"
)

// MustGet will return the env or panic if it is not present
func MustGet(k string) string {
    v := os.Getenv(k)
    if v == "" {
        logger.Panicln("ENV missing, key: " + k)
    }
    return v
}

// MustGetBool will return the env as boolean or panic if it is not present
func MustGetBool(k string) bool {
    v := os.Getenv(k)
    if v == "" {
        logger.Panicln("ENV missing, key: " + k)
    }
    b, err := strconv.ParseBool(v)
    if err != nil {
        logger.Panicln("ENV err: [" + k + "]\n" + err.Error())
    }
    return b
}
```

The code is pretty much self-explainatory, if a ENV var does not exists, the
program will panic, we need these to run. Now, changing:

- `pkg/server/main.go` to this:

```go
package server

import (
    "github.com/cmelgarejo/go-gql-server/internal/logger"

    "github.com/gin-gonic/gin"
    "github.com/cmelgarejo/go-gql-server/internal/handlers"
    "github.com/cmelgarejo/go-gql-server/pkg/utils"
)

var host, port string

func init() {
    host = utils.MustGet("SERVER_HOST")
    port = utils.MustGet("SERVER_PORT")
}

// Run spins up the server
func Run() {
    r := gin.Default()
    // Simple keep-alive/ping handler
    r.GET("/ping", handlers.Ping())
    // Inform the user where the server is listening
    logger.Println("Running @ http://" + host + ":" + port)
    // Print out and exit(1) to the OS if the server cannot run
    logger.Fatalln(r.Run(host + ":" + port))

}
```

And we see how's it is starting to take form as a well laid out project!

We can still make a couple of things, like running server locally by using this
script:

- `scripts/run.sh` (don't forget to `chmod +x` this one too)

```bash
#!/bin/sh
srcPath="cmd"
pkgFile="main.go"
app="gql-server"
src="$srcPath/$app/$pkgFile"

printf "\nStart running: $app\n"
# Set all ENV vars for the server to run
export $(grep -v '^#' .env | xargs) && time go run $src
# This should unset all the ENV vars, just in case.
unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)
printf "\nStopped running: $app\n\n"
```

This will set up ENV vars, and `go run` the server.

You can also add a `.gitignore` file, run `git init` set up the origin to your
own repository and `git push` it :)

This will continue in [Part 2](PART2.md) where we'll add the GQLGen portion of
the server!

All the code is available [here](https://github.com/cmelgarejo/go-gql-server/tree/tutorial/part-1)
