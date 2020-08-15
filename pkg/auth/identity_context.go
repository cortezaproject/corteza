package auth

import (
	"context"
)

type (
	identityCtxKey struct{}
	jwtCtxKey      struct{}
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

func SetJwtToContext(ctx context.Context, jwt string) context.Context {
	return context.WithValue(ctx, jwtCtxKey{}, jwt)

}

func GetJwtFromContext(ctx context.Context) string {
	if jwt, ok := ctx.Value(jwtCtxKey{}).(string); ok {
		return jwt
	} else {
		return ""
	}
}

// SetSuperUserContext stores system user as identity
// and accompanying JWT for it to the context
func SetSuperUserContext(ctx context.Context) context.Context {
	su := NewSuperUserIdentity()

	ctx = SetIdentityToContext(ctx, su)

	if DefaultJwtHandler != nil {
		ctx = SetJwtToContext(ctx, DefaultJwtHandler.Encode(su))
	}

	return ctx
}
