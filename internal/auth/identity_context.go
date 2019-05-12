package auth

import (
	"context"
)

type (
	ctxKey int
)

var (
	identityCtxKey ctxKey
)

func SetIdentityToContext(ctx context.Context, identity Identifiable) context.Context {
	return context.WithValue(ctx, identityCtxKey, identity)
}

func GetIdentityFromContext(ctx context.Context) Identifiable {
	if identity, ok := ctx.Value(identityCtxKey).(Identifiable); ok {
		return identity
	} else {
		return NewIdentity(0)
	}
}
