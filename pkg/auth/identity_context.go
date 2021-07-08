package auth

import (
	"context"
)

type (
	identityCtxKey struct{}
)

func SetIdentityToContext(ctx context.Context, identity Identifiable) context.Context {
	return context.WithValue(ctx, identityCtxKey{}, identity)
}

// GetIdentityFromContext always returns identity, either valid or anonymous
//
// For anonymous user, it auto appends all anonymous defined on the system
func GetIdentityFromContext(ctx context.Context) Identifiable {
	if identity, ok := ctx.Value(identityCtxKey{}).(Identifiable); ok && identity != nil && identity.Valid() {
		return identity
	} else {
		return Anonymous()
	}
}
