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

func ScopeFromContext(ctx context.Context) *types.Scp {
	return ctx.Value(ContextKeyScope).(*types.Scp)
}
