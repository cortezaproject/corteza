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

func GetIdentityFromContext(ctx context.Context) Identifiable {
	if identity, ok := ctx.Value(identityCtxKey{}).(Identifiable); ok && identity != nil {
		return identity
	} else {
		return NewIdentity(0)
	}
}

// SetSuperUserContext stores system user as identity
// and accompanying JWT for it to the context
func SetSuperUserContext(ctx context.Context) context.Context {
	su := NewSuperUserIdentity()

	return SetIdentityToContext(ctx, su)
}
