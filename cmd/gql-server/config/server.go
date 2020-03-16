package config

import (
	"strings"

	"github.com/cmelgarejo/go-gql-server/pkg/utils"
)

// Server meh
func Server() *utils.ServerConfig {
	return &utils.ServerConfig{
		Version:        utils.MustGet("APP_VERSION"),
		Env:            utils.MustGet("APP_ENV"),
		Host:           utils.MustGet("SERVER_HOST"),
		Port:           utils.MustGet("SERVER_PORT"),
		URISchema:      utils.MustGet("SERVER_URI_SCHEMA"),
		ServiceVersion: utils.MustGet("SERVER_PATH_VERSION"),
		SessionSecret:  utils.MustGet("SESSION_SECRET"),
		Frontend: utils.FrontendConfig{
			LoginCallbackURL: utils.MustGet("FRONTEND_LOGIN_CALLBACK_URL"),
		},
		JWT: utils.JWTConfig{
			Secret:    utils.MustGet("AUTH_JWT_SECRET"),
			Algorithm: utils.MustGet("AUTH_JWT_SIGNING_ALGORITHM"),
		},
		GraphQL: utils.GQLConfig{
			ComplexityLimit:        utils.MustGetInt32("GQL_SERVER_GRAPHQL_COMPLEXITY_LIMIT"),
			Path:                   utils.MustGet("GQL_SERVER_GRAPHQL_PATH"),
			PlaygroundPath:         utils.MustGet("GQL_SERVER_GRAPHQL_PLAYGROUND_PATH"),
			IsPlaygroundEnabled:    utils.MustGetBool("GQL_SERVER_GRAPHQL_PLAYGROUND_ENABLED"),
			IsIntrospectionEnabled: utils.MustGetBool("GQL_SERVER_GRAPHQL_INTROSPECTION_ENABLED"),
		},
		Database: utils.DBConfig{
			Dialect:     utils.MustGet("GORM_DIALECT"),
			DSN:         utils.MustGet("GORM_CONNECTION_DSN"),
			SeedDB:      utils.MustGetBool("GORM_SEED_DB"),
			LogMode:     utils.MustGetBool("GORM_LOGMODE"),
			AutoMigrate: utils.MustGetBool("GORM_AUTOMIGRATE"),
		},
		AuthProviders: []utils.AuthProvider{
			{
				Provider:  "auth0",
				ClientKey: utils.MustGet("PROVIDER_AUTH0_KEY"),
				Secret:    utils.MustGet("PROVIDER_AUTH0_SECRET"),
				Domain:    utils.MustGet("PROVIDER_AUTH0_DOMAIN"),
				Scopes:    strings.Split(utils.MustGet("PROVIDER_AUTH0_SCOPES"), ","),
			},
			{
				Provider:  "facebook",
				ClientKey: utils.MustGet("PROVIDER_FACEBOOK_KEY"),
				Secret:    utils.MustGet("PROVIDER_FACEBOOK_SECRET"),
			},
			{
				Provider:  "google",
				ClientKey: utils.MustGet("PROVIDER_GOOGLE_KEY"),
				Secret:    utils.MustGet("PROVIDER_GOOGLE_SECRET"),
			},
			{
				Provider:  "twitter",
				ClientKey: utils.MustGet("PROVIDER_TWITTER_KEY"),
				Secret:    utils.MustGet("PROVIDER_TWITTER_SECRET"),
			},
		},
	}
}
