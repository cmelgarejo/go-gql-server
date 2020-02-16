package utils

// ContextKeys holds the context keys throughout the project
type ContextKeys struct {
	GothicProviderCtxKey ContextKey // Provider for Gothic library
	ProviderCtxKey       ContextKey // Provider in Auth
	UserCtxKey           ContextKey // User db object in Auth
}

var (
	// ProjectContextKeys the project's context keys
	ProjectContextKeys = ContextKeys{
		GothicProviderCtxKey: "provider",
		ProviderCtxKey:       "gg-provider",
		UserCtxKey:           "gg-auth-user",
	}
)
