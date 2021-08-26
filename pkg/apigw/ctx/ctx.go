package ctx

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
)

type ContextKey string

const ContextKeyScope ContextKey = "scope"

func ScopeToContext(ctx context.Context, s *types.Scp) context.Context {
	return context.WithValue(ctx, ContextKeyScope, s)
}

func ScopeFromContext(ctx context.Context) (ss *types.Scp) {
	s := ctx.Value(ContextKeyScope)

	if s == nil {
		return &types.Scp{}
	}

	return s.(*types.Scp)
}
