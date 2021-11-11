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
	if i := GetIdentityFromContextWithKey(ctx, identityCtxKey{}); i != nil && i.Valid() {
		return i
	} else {
		return Anonymous()
	}
}

func GetIdentityFromContextWithKey(ctx context.Context, key interface{}) Identifiable {
	if i, ok := ctx.Value(key).(Identifiable); ok {
		return i
	} else {
		return nil
	}
}
