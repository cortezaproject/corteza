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
