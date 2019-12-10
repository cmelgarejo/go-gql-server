package server

import (
	"github.com/cmelgarejo/go-gql-server/pkg/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/google"
)

// InitalizeAuthProviders does just that, with Goth providers
func InitalizeAuthProviders(cfg *utils.ServerConfig) error {
	providers := []goth.Provider{}
	// Initialize Goth providers
	for _, p := range cfg.AuthProviders {
		switch p.Provider {
		case "google":
			providers = append(providers, google.New(p.ClientKey, p.Secret,
				cfg.SchemaVersionedEndpoint("/auth/"+p.Provider+"/callback"),
				p.Scopes...))
		case "auth0":
			providers = append(providers, auth0.New(p.ClientKey, p.Secret,
				cfg.SchemaVersionedEndpoint("/auth/"+p.Provider+"/callback"),
				p.Domain, p.Scopes...))
		}
	}
	goth.UseProviders(providers...)
	return nil
}
