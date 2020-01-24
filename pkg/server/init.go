package server

import (
	"github.com/cmelgarejo/go-gql-server/pkg/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitter"
)

// InitalizeAuthProviders does just that, with Goth providers
func InitalizeAuthProviders(cfg *utils.ServerConfig) error {
	providers := []goth.Provider{}
	// Initialize Goth providers
	for _, p := range cfg.AuthProviders {
		switch p.Provider {
		case "auth0":
			providers = append(providers, auth0.New(p.ClientKey, p.Secret,
				cfg.SchemaVersionedEndpoint("/auth/"+p.Provider+"/callback"),
				p.Domain, p.Scopes...))
		case "facebook":
			providers = append(providers, facebook.New(p.ClientKey, p.Secret,
				cfg.SchemaVersionedEndpoint("/auth/"+p.Provider+"/callback"),
				p.Scopes...))
		case "google":
			providers = append(providers, google.New(p.ClientKey, p.Secret,
				cfg.SchemaVersionedEndpoint("/auth/"+p.Provider+"/callback"),
				p.Scopes...))
		case "twitter":
			providers = append(providers, twitter.New(p.ClientKey, p.Secret,
				cfg.SchemaVersionedEndpoint("/auth/"+p.Provider+"/callback")))
		}
	}
	goth.UseProviders(providers...)
	return nil
}
