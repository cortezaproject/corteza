package auth

import (
	"context"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/pkg/errors"
)

type (
	ctxKey int
)

var (
	identityCtxKey ctxKey
)

func getIdentityClaimFromContext(ctx context.Context) (uint64, error) {
	if jwt, claims, err := jwtauth.FromContext(ctx); err != nil && err.Error() == "jwtauth: no token found" {
		// Token is not required. We're handling anonymous request lower down the line
		return 0, nil
	} else if err != nil {
		return 0, errors.Wrap(err, "failed to authorize request")
	} else if !jwt.Valid {
		return 0, errors.New("JWT not valid")
	} else if idInterface, has := claims.Get("sub"); !has {
		return 0, errors.New("Malformed JWT claims")
	} else if idString, ok := idInterface.(string); !ok {
		return 0, errors.New("Malformed JWT claims (sub not a string)")
	} else if identityId, err := strconv.ParseUint(idString, 10, 64); err != nil {
		return 0, errors.Wrap(err, "Malformed JWT claims")
	} else {
		return identityId, nil
	}
}

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
