package scope

import "context"

type (
	Scope struct {
		TenantID      uint64
		ProjectID     uint64 // 0 = tenant-level operation
		RootProjectID uint64 // root of the revision chain; equals ProjectID when not a revision
	}

	scopeCtxKey struct{}
)

// TenantOnly returns true when no project is set.
func (s Scope) TenantOnly() bool { return s.ProjectID == 0 }

// Valid returns true when at minimum a tenant is set.
func (s Scope) Valid() bool { return s.TenantID != 0 }

func SetScopeToContext(ctx context.Context, s Scope) context.Context {
	return context.WithValue(ctx, scopeCtxKey{}, s)
}

// GetScopeFromContext returns the scope stored in ctx.
// Returns zero-value Scope (invalid) when not set.
func GetScopeFromContext(ctx context.Context) Scope {
	if s, ok := ctx.Value(scopeCtxKey{}).(Scope); ok {
		return s
	}
	return Scope{}
}
