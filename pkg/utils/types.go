package utils

// ContextKey defines a type for context keys shared in the app
type ContextKey string

// ServerConfig defines the configuration for the server
type ServerConfig struct {
	Version        string
	Env            string
	Host           string
	Port           string
	URISchema      string
	ServiceVersion string
	SessionSecret  string
	Frontend       FrontendConfig
	JWT            JWTConfig
	GraphQL        GQLConfig
	Database       DBConfig
	AuthProviders  []AuthProvider
}

//JWTConfig defines the options for JWT tokens
type JWTConfig struct {
	Secret    string
	Algorithm string
}

// GQLConfig defines the configuration for the GQL Server
type GQLConfig struct {
	ComplexityLimit        int
	Path                   string
	PlaygroundPath         string
	IsPlaygroundEnabled    bool
	IsIntrospectionEnabled bool
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

//FrontendConfig defines the options for the Frontend
type FrontendConfig struct {
	LoginCallbackURL string
}

func getValidHost(host string) string {
	if host == ":" {
		return "localhost"
	}
	return host
}

// ListenEndpoint builds the endpoint string (host + port)
func (s *ServerConfig) ListenEndpoint() string {
	if s.Port == "80" {
		return s.Host
	}
	if s.Host == ":" {
		return s.Host + s.Port

	}
	return s.Host + ":" + s.Port
}

// VersionedEndpoint builds the endpoint `string (host + port + version)
func (s *ServerConfig) VersionedEndpoint(path string) string {
	return "/" + s.ServiceVersion + path
}

// SchemaVersionedEndpoint builds the schema endpoint string (schema + host + port + version)
func (s *ServerConfig) SchemaVersionedEndpoint(path string) string {
	if s.Port == "80" {
		return s.URISchema + getValidHost(s.Host) + "/" + s.ServiceVersion + path
	}
	return s.URISchema + getValidHost(s.Host) + ":" + s.Port + "/" + s.ServiceVersion + path
}
