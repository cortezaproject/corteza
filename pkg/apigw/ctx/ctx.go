package ctx

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/apigw/types"
)

type ContextKey string

const ContextKeyScope ContextKey = "scope"
const ContextKeyProfiler ContextKey = "profiler"

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

func ProfilerToContext(ctx context.Context, h interface{}) context.Context {
	return context.WithValue(ctx, ContextKeyProfiler, h)
}

func ProfilerFromContext(ctx context.Context) (h interface{}) {
	hh := ctx.Value(ContextKeyProfiler)

	if hh == nil {
		return nil
	}

	return hh
}
