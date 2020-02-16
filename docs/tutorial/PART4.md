# Creating an opinionated Go GQL Server - Part 4

This is part of a series, not suited for beginners, but welcome too! You can
check it out from the start

- [Part 1](PART1.md)

And review how to get up to this part, as always all the code is at github for
[part 4](https://github.com/cmelgarejo/go-gql-server/tree/tutorial/part-4)

---

Today we'll add Goth and user Auth0 as our authentication provider to connect
users via google and/or facebook - or any Identity provider using OAuth2

## Adding Goth into the project

Let's add Goth to our project:

> `\$ go get -u github.com/markbates/goth

We can start referencing goth into our server code, there are a few steps to
integrate Goth completely into the flow of the app.

First, I want to keep the configuration as clean as possible and do it from the
go-gql-server definition in `cmd/go-gql-server` for this I create a configuration type
in `pkg/utils/types.go`:

```go
// AuthProvider defines the configuration for the Goth config
type AuthProvider struct {
    Provider  string
    ClientKey string
    Secret    string
    Domain    string // If needed, like with auth0
    Scopes    []string
}
```

This way we can use configuration properties needed to instantiate a Goth in
the server, also, we can create a way to type in all configurations in a ordered
fashion in `utils/types.go`:

```go
package utils

// ContextKey defines a type for context keys shared in the app
type ContextKey string

// ServerConfig defines the configuration for the server
type ServerConfig struct {
    Host          string
    Port          string
    URISchema     string
    Version       string
    SessionSecret string
    JWT           JWTConfig
    GraphQL       GQLConfig
    Database      DBConfig
    AuthProviders []AuthProvider
}

//JWTConfig defines the options for JWT tokens
type JWTConfig struct {
    Secret    string
    Algorithm string
}

// GQLConfig defines the configuration for the GQL Server
type GQLConfig struct {
    Path                string
    PlaygroundPath      string
    IsPlaygroundEnabled bool
}

// DBConfig defines the configuration for the DB config
type DBConfig struct {
    Dialect     string
    DSN         string
    SeedDB      bool
    LogMode     bool
    AutoMigrate bool
}

// AuthProvider defines the configuration for the Goth config
type AuthProvider struct {
    Provider  string
    ClientKey string
    Secret    string
    Domain    string // If needed, like with auth0
    Scopes    []string
}

// ListenEndpoint builds the endpoint string (host + port)
func (s *ServerConfig) ListenEndpoint() string {
    if s.Port == "80" {
        return s.Host
    }
    return s.Host + ":" + s.Port
}

// VersionedEndpoint builds the endpoint string (host + port + version)
func (s *ServerConfig) VersionedEndpoint(path string) string {
    return "/" + s.ServiceVersion + path
}

// SchemaVersionedEndpoint builds the schema endpoint string (schema + host + port + version)
func (s *ServerConfig) SchemaVersionedEndpoint(path string) string {
    if s.Port == "80" {
        return s.URISchema + s.Host + "/" + s.ServiceVersion + path
    }
    return s.URISchema + s.Host + ":" + s.Port + "/" + s.ServiceVersion + path
}

```

You can notice that now we can even versionate the endpoints here, all from the
configuration file `.env` (there's the `.env/example` file in the repo too)

So, then we can add into the `serverconf` variable a list of providers, and all
the different configurations and pass the server configuration down to the
`pkg/server/main.go` from `cmd/go-gql-server/main.go`:

```go
var serverconf = &utils.ServerConfig{
        Host:          utils.MustGet("SERVER_HOST"),
        Port:          utils.MustGet("SERVER_PORT"),
        URISchema:     utils.MustGet("SERVER_URI_SCHEMA"),
        Version:       utils.MustGet("SERVER_PATH_VERSION"),
        SessionSecret: utils.MustGet("SESSION_SECRET"),
        JWT: utils.JWTConfig{
            Secret:    utils.MustGet("AUTH_JWT_SECRET"),
            Algorithm: utils.MustGet("AUTH_JWT_SIGNING_ALGORITHM"),
        },
        GraphQL: utils.GQLConfig{
            Path:                utils.MustGet("GQL_SERVER_GRAPHQL_PATH"),
            PlaygroundPath:      utils.MustGet("GQL_SERVER_GRAPHQL_PLAYGROUND_PATH"),
            IsPlaygroundEnabled: utils.MustGetBool("GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED"),
        },
        Database: utils.DBConfig{
            Dialect:     utils.MustGet("GORM_DIALECT"),
            DSN:         utils.MustGet("GORM_CONNECTION_DSN"),
            SeedDB:      utils.MustGetBool("GORM_SEED_DB"),
            LogMode:     utils.MustGetBool("GORM_LOGMODE"),
            AutoMigrate: utils.MustGetBool("GORM_AUTOMIGRATE"),
        },
        AuthProviders: []utils.AuthProvider{
            utils.AuthProvider{
                Provider:  "google",
                ClientKey: utils.MustGet("PROVIDER_GOOGLE_KEY"),
                Secret:    utils.MustGet("PROVIDER_GOOGLE_SECRET"),
            },
            utils.AuthProvider{
                Provider:  "auth0",
                ClientKey: utils.MustGet("PROVIDER_AUTH0_KEY"),
                Secret:    utils.MustGet("PROVIDER_AUTH0_SECRET"),
                Domain:    utils.MustGet("PROVIDER_AUTH0_DOMAIN"),
                Scopes:    strings.Split(utils.MustGet("PROVIDER_AUTH0_SCOPES"), ","),
            },
        },
    }
    orm, err := orm.Factory(serverconf)
    defer orm.DB.Close()
    if err != nil {
        logger.Panic(err)
    }
    server.Run(serverconf, orm)
```

Now, getting back into the flow, we have to create handlers for our server that
will manage the authenticated paths, this is done creating the proper middleware
for the auth handling, let's create a route in our server package called _auth_
`pkg/server/routes/auth.go`:

```go
package routes

import (
    "github.com/cmelgarejo/go-gql-server/internal/handlers/auth"
    "github.com/cmelgarejo/go-gql-server/internal/orm"
    "github.com/cmelgarejo/go-gql-server/pkg/utils"
    "github.com/gin-gonic/gin"
)

// Auth routes
func Auth(cfg *utils.ServerConfig, r *gin.Engine, orm *orm.ORM) error {
    // OAuth handlers
    g := r.Group(cfg.VersionedEndpoint("/auth"))
    g.GET("/:provider", auth.Begin())
    g.GET("/:provider/callback", auth.Callback(cfg, orm))
    // g.GET(:provider/refresh", auth.Refresh(cfg, orm))
    return nil
}
```

As you might have noticed, now the routes receive the `ServerConfig` struct and
that way we can use the `VersionedEndpoint` method to build the versioned path
for our server, and the `/:provider` path refers to the provider you want to
authenticate with, as configured now, I just added google and auth0, but since
Goth takes many providers, you can add you own as you like with in the config
array struct of _providers_

Now, to create the `internal/handlers/auth` middleware and handler for this
to work: `internal/handlers/auth/main.go`:

```go
//
package auth

import (
    "context"
    "net/http"

    "github.com/cmelgarejo/go-gql-server/pkg/utils"
    "github.com/gin-gonic/gin"
)

func addProviderToContext(c *gin.Context, value interface{}) *http.Request {
    return c.Request.WithContext(context.WithValue(c.Request.Context(),
        string(utils.ProjectContextKeys.ProviderCtxKey), value))
}
```

Here we see that we add the provider handler to the current context of the
request, in order for this handle the auth downstream of it, how? well, binding
the auth middleware `internal/handlers/auth/middleware/main.go`:

```go
package middleware

import (
    "context"
    "errors"
    "net/http"
    "strings"

    "github.com/cmelgarejo/go-gql-server/pkg/utils"
    "github.com/dgrijalva/jwt-go"

    "github.com/gin-gonic/gin"
)

var (
    // TokenHeadName is a string in the header. Default value is "Bearer"
    TokenHeadName = "Bearer"

    // APIKeyLookup is a string in the form of "<source>:<name>" that is used
    // to extract token from the request.
    // Optional. Default value "header:Authorization".
    // Possible values:
    // - "header:<name>"
    // - "query:<name>"
    // - "cookie:<name>"
    APIKeyLookup = "query:api_key,cookie:api_key,header:X-API-KEY"

    // TokenLookup is a string in the form of "<source>:<name>" that is used
    // to extract token from the request.
    // Optional. Default value "header:Authorization".
    // Possible values:
    // - "header:<name>"
    // - "query:<name>"
    // - "cookie:<name>"
    TokenLookup = "query:token,cookie:jwt,header:Authorization"

    // ErrNoClaims when HTTP status 403 is given
    ErrNoClaims = errors.New("invalid token")

    // ErrForbidden when HTTP status 403 is given
    ErrForbidden = errors.New("you don't have permission to access this resource")

    // ErrExpiredToken indicates JWT token has expired. Can't refresh.
    ErrExpiredToken = errors.New("token is expired")

    // ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
    ErrEmptyAuthHeader = errors.New("auth header is empty")

    // ErrEmptyAPIKeyHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
    ErrEmptyAPIKeyHeader = errors.New("api key header is empty")

    // ErrMissingExpField missing exp field in token
    ErrMissingExpField = errors.New("missing exp field")

    // ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
    ErrInvalidAuthHeader = errors.New("auth header is invalid")

    // ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
    ErrEmptyQueryToken = errors.New("query token is empty")

    // ErrEmptyCookieToken can be thrown if authing with a cookie, the token cokie is empty
    ErrEmptyCookieToken = errors.New("cookie token is empty")

    // ErrEmptyParamToken can be thrown if authing with parameter in path, the parameter in path is empty
    ErrEmptyParamToken = errors.New("parameter token is empty")

    // ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
    ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")
)

func jwtFromHeader(c *gin.Context, key string) (string, error) {
    authHeader := c.Request.Header.Get(key)

    if authHeader == "" {
        return "", ErrEmptyAuthHeader
    }

    parts := strings.SplitN(authHeader, " ", 2)
    if !(len(parts) == 2 && parts[0] == TokenHeadName) {
        return "", ErrInvalidAuthHeader
    }

    return parts[1], nil
}

func apiKeyFromHeader(c *gin.Context, key string) (string, error) {
    apiKey := c.Request.Header.Get(key)
    if apiKey == "" {
        return "", ErrEmptyAPIKeyHeader
    }
    return apiKey, nil
}

func tokenFromQuery(c *gin.Context, key string) (string, error) {
    token := c.Query(key)
    if token == "" {
        return "", ErrEmptyQueryToken
    }
    return token, nil
}

func tokenFromCookie(c *gin.Context, key string) (string, error) {
    cookie, _ := c.Cookie(key)
    if cookie == "" {
        return "", ErrEmptyCookieToken
    }
    return cookie, nil
}

func tokenFromParam(c *gin.Context, key string) (string, error) {
    token := c.Param(key)
    if token == "" {
        return "", ErrEmptyParamToken
    }
    return token, nil
}

// ParseToken parse jwt token from gin context
func ParseToken(c *gin.Context, cfg *utils.ServerConfig) (t *jwt.Token, err error) {
    var token string
    methods := strings.Split(TokenLookup, ",")
    for _, method := range methods {
        if len(token) > 0 {
            break
        }
        parts := strings.Split(strings.TrimSpace(method), ":")
        k := strings.TrimSpace(parts[0])
        v := strings.TrimSpace(parts[1])
        switch k {
        case "header":
            token, err = jwtFromHeader(c, v)
        case "query":
            token, err = tokenFromQuery(c, v)
        case "cookie":
            token, err = tokenFromCookie(c, v)
        case "param":
            token, err = tokenFromParam(c, v)
        }
    }
    if err != nil {
        return nil, err
    }
    SigningAlgorithm := cfg.JWT.Algorithm
    Key := []byte(cfg.JWT.Secret)
    return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        if jwt.GetSigningMethod(SigningAlgorithm) != t.Method {
            return nil, ErrInvalidSigningAlgorithm
        }
        // save token string if vaild
        // c.Set("AUTH_JWT_TOKEN", token)
        return Key, nil
    })
}

// ParseAPIKey parse api key from gin context
func ParseAPIKey(c *gin.Context, cfg *utils.ServerConfig) (apiKey string, err error) {
    methods := strings.Split(APIKeyLookup, ",")
    for _, method := range methods {
        if len(apiKey) > 0 {
            break
        }
        parts := strings.Split(strings.TrimSpace(method), ":")
        k := strings.TrimSpace(parts[0])
        v := strings.TrimSpace(parts[1])
        switch k {
        case "header":
            apiKey, err = apiKeyFromHeader(c, v)
        case "query":
            apiKey, err = tokenFromQuery(c, v)
        case "cookie":
            apiKey, err = tokenFromCookie(c, v)
        case "param":
            apiKey, err = tokenFromParam(c, v)
        }
    }
    if err != nil {
        return "", err
    }
    return apiKey, nil
}

func addToContext(c *gin.Context, key utils.ContextKey, value interface{}) *http.Request {
    return c.Request.WithContext(context.WithValue(c.Request.Context(), key, value))
}
```

You might noticed again that here we also want to use API keys to auth requests,
these will be handled on our side, generating secure enough API keys for machine
to machine communications and integrations to our graphql API.

Again, all that this middleware does is trying to get the JWT tokens generated
by our providers and try to autheticate against our `users` records, that might
come from a cookie, a header or event a querystring param, as needed or
implemented by our client(s)

And the actual auth middleware handling method that will enforce all this is in
`internal/handlers/auth/middleware/auth.go`:

```go
package middleware

import (
    "net/http"

    "github.com/cmelgarejo/go-gql-server/internal/logger"
    "github.com/cmelgarejo/go-gql-server/internal/orm"
    "github.com/cmelgarejo/go-gql-server/pkg/utils"
    "github.com/dgrijalva/jwt-go"

    "github.com/gin-gonic/gin"
)

func authError(c *gin.Context, err error) {
    errKey := "message"
    errMsgHeader := "[Auth] error: "
    e := gin.H{errKey: errMsgHeader + err.Error()}
    c.AbortWithStatusJSON(http.StatusUnauthorized, e)
}

// Middleware wraps the request with auth middleware
func Middleware(path string, cfg *utils.ServerConfig, orm *orm.ORM) gin.HandlerFunc {
    logger.Info("[Auth.Middleware] Applied to path: ", path)
    return gin.HandlerFunc(func(c *gin.Context) {
        if a, err := ParseAPIKey(c, cfg); err == nil {
            user, err := orm.FindUserByAPIKey(a)
            if err != nil {
                authError(c, ErrForbidden)
            }
            c.Next()
        } else {
            if err != ErrEmptyAPIKeyHeader {
                authError(c, err)
            } else {
                t, err := ParseToken(c, cfg)
                if err != nil {
                    authError(c, err)
                } else {
                    if claims, ok := t.Claims.(jwt.MapClaims); ok {
                        if claims["exp"] != nil {
                            issuer := claims["iss"].(string)
                            userid := claims["jti"].(string)
                            email := claims["email"].(string)
                            if user, err := orm.FindUserByJWT(email, issuer, userid); err != nil {
                                authError(c, ErrForbidden)
                            } else {
                                c.Request = addToContext(c, utils.ProjectContextKeys.UserCtxKey, user)
                                c.Next()
                            }
                        } else {
                            authError(c, ErrMissingExpField)
                        }
                    } else {
                        authError(c, err)
                    }
                }
            }
        }
    })
}
```

Here the real search for the user is made in `FindUserByJWT` for JWT tokens and
`FindUserByAPIKey` for the API keys if sent, those functions live in
`internal/orm/main.go`.

```go
//FindUserByAPIKey finds the user that is related to the API key
func (o *ORM) FindUserByAPIKey(apiKey string) (*models.User, error) {
    db := o.DB
    uak := &models.UserAPIKey{}
    if apiKey == "" {
        return nil, errors.New("API key is empty")
    }
    if err := db.Preload("User").Where("api_key = ?", apiKey).Find(uak).Error; err != nil {
        return nil, err
    }
    return &uak.User, nil
}

// FindUserByJWT finds the user that is related to the APIKey token
func (o *ORM) FindUserByJWT(email string, provider string, userID string) (*models.User, error) {
    db := o.DB
    up := &models.UserProfile{}
    if provider == "" || userID == "" {
        return nil, errors.New("provider or userId empty")
    }
    if err := db.Preload("User").Where("email  = ? AND provider = ? AND external_user_id = ?", email, provider, userID).First(up).Error; err != nil {
        return nil, err
    }
    return &up.User, nil
}

// UpsertUserProfile saves the user if doesn't exists and adds the OAuth profile
func (o *ORM) UpsertUserProfile(input *goth.User) (*models.User, error) {
    db := o.DB
    u := &models.User{}
    up := &models.UserProfile{}
    u, err := transformations.GothUserToDBUser(input, false)
    if err != nil {
        return nil, err
    }
    if tx := db.Where("email = ?", input.Email).First(u); !tx.RecordNotFound() && tx.Error != nil {
        return nil, tx.Error
    }
    if tx := db.Model(u).Save(u); tx.Error != nil {
        return nil, err
    }
    if tx := db.Where("email = ? AND provider = ? AND external_user_id = ?",
        input.Email, input.Provider, input.UserID).First(up); !tx.RecordNotFound() && tx.Error != nil {
        return nil, err
    }
    up, err = transformations.GothUserToDBUserProfile(input, false)
    if err != nil {
        return nil, err
    }
    up.User = *u
    if tx := db.Model(up).Save(up); tx.Error != nil {
        return nil, tx.Error
    }
    logger.Debug(u.ID)
    return u, nil
}

```

There's also the `UpsertUserProfile` method that saves the profile into our
database so it can be retrieved later as we want at `user_profiles`

Continuing the auth flow, here's the handlers for the routing of the login/out
with oAuth providers `internal/handlers/auth/handlers.go`:

```go
package auth

import (
    "net/http"
    "time"

    "github.com/cmelgarejo/go-gql-server/internal/orm"

    "github.com/dgrijalva/jwt-go"

    "github.com/cmelgarejo/go-gql-server/internal/logger"
    "github.com/cmelgarejo/go-gql-server/pkg/utils"
    "github.com/gin-gonic/gin"
    "github.com/markbates/goth/gothic"
)

// Claims JWT claims
type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

// Begin login with the auth provider
func Begin() gin.HandlerFunc {
    return func(c *gin.Context) {
        // You have to add value context with provider name to get provider name in GetProviderName method
        c.Request = addProviderToContext(c, c.Param("provider"))
        // try to get the user without re-authenticating
        if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err != nil {
            gothic.BeginAuthHandler(c.Writer, c.Request)
        } else {
            logger.Debugf("user: %#v", gothUser)
        }
    }
}

// Callback callback to complete auth provider flow
func Callback(cfg *utils.ServerConfig, orm *orm.ORM) gin.HandlerFunc {
    return func(c *gin.Context) {
        // You have to add value context with provider name to get provider name in GetProviderName method
        c.Request = addProviderToContext(c, c.Param("provider"))
        user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
        if err != nil {
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }
        u, err := orm.FindUserByJWT(user.Email, user.Provider, user.UserID)
        // logger.Debugf("gothUser: %#v", user)
        if err != nil {
            if u, err = orm.UpsertUserProfile(&user); err != nil {
                logger.Errorf("[Auth.CallBack.UserLoggedIn.UpsertUserProfile.Error]: %q", err)
                c.AbortWithError(http.StatusInternalServerError, err)
            }
        }
        // logger.Debug("[Auth.CallBack.UserLoggedIn.USER]: ", u)
        logger.Debug("[Auth.CallBack.UserLoggedIn]: ", u.ID)
        jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod(cfg.JWT.Algorithm), Claims{
            Email: user.Email,
            StandardClaims: jwt.StandardClaims{
                Id:        user.UserID,
                Issuer:    user.Provider,
                IssuedAt:  time.Now().UTC().Unix(),
                NotBefore: time.Now().UTC().Unix(),
                ExpiresAt: user.ExpiresAt.UTC().Unix(),
            },
        })
        token, err := jwtToken.SignedString([]byte(cfg.JWT.Secret))
        if err != nil {
            logger.Error("[Auth.Callback.JWT] error: ", err)
            c.AbortWithError(http.StatusInternalServerError, err)
            return
        }
        logger.Debug("token: ", token)
        json := gin.H{
            "type":          "Bearer",
            "token":         token,
            "refresh_token": user.RefreshToken,
        }
        c.JSON(http.StatusOK, json)
    }
}

// Logout logs out of the auth provider
func Logout() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Request = addProviderToContext(c, c.Param("provider"))
        gothic.Logout(c.Writer, c.Request)
        c.Writer.Header().Set("Location", "/")
        c.Writer.WriteHeader(http.StatusTemporaryRedirect)
    }
}
```

This will allow us to go to `http://localhost:7777/v1/auth/:provider` such as we
have configured `google` and `auth0` and we'll have an user created in our API!

Now, you have to configure an app properly at [GCP](https://console.cloud.google.com/)
and then, fill in the data needed for the provider as indicated in
`cmd/go-gql-server/main.go`, this will allow you to use the provider and
navigate to `http://localhost:7777/v1/auth/google` and to test it out:

![oauth consent](images/p4-p3.png?raw=true "Google App consent page")

Sorry for the spanish but it's the usual OAuth2 consent screen for apps from
google :D

Should you have it all correct, this will happen:

![oauth consent](images/p4-p2.png?raw=true "OAuth consent page")

And then we have a user Signup/in already doing all the work :D

Now, use the token in the playground `http://localhost:7777/v1/graphql/playground`

![oauth bearer](images/p4-p1.png?raw=true "token usage")

And run a query or mutation on our server

![gql is showing db data!](images/p4-p4.png?raw=true "Our API responding with data")

So with this, we've integrted all!

- gqlgen
- goth
- gorm

How about this: Next part it's setting up a RBAC layer, so we can determine
which roles/users can access information in our API :)

Be on the lookup for that! Again, like the former parts, all the code is
available in the [repository here](https://github.com/cmelgarejo/go-gql-server/tree/tutorial/part-4)!
If you have questions, critiques and comments go ahead and let's learn more
together!
