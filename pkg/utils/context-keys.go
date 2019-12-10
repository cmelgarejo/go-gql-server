package utils

// ContextKeys holds the context keys throught the project
type ContextKeys struct {
	ProviderCtxKey ContextKey // Provider in Auth
	UserCtxKey     ContextKey // User db object in Auth
}

var (
	// ProjectContextKeys the project's context keys
	ProjectContextKeys = ContextKeys{
		ProviderCtxKey: "provider",
		UserCtxKey:     "auth-user",
	}
)
