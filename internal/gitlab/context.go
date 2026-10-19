package gitlab

import "context"

type contextKey int

const credsContextKey contextKey = iota

// Creds holds GitLab instance URL and personal access token for a request.
type Creds struct {
	BaseURL string
	Token   string
}

// ContextWithCreds attaches GitLab credentials to ctx.
func ContextWithCreds(ctx context.Context, creds Creds) context.Context {
	return context.WithValue(ctx, credsContextKey, creds)
}

// CredsFromContext returns credentials previously stored on ctx.
func CredsFromContext(ctx context.Context) (Creds, bool) {
	v, ok := ctx.Value(credsContextKey).(Creds)
	return v, ok
}
