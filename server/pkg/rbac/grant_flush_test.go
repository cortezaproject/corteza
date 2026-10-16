package rbac

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// failingRulesStore fails the first upsert and records the rest
type failingRulesStore struct {
	nopRulesStore
	failed   bool
	upserted RuleSet
}

func (s *failingRulesStore) UpsertRbacRule(_ context.Context, rr ...*Rule) error {
	if !s.failed {
		s.failed = true
		return errors.New("store unavailable")
	}
	s.upserted = append(s.upserted, rr...)
	return nil
}

// A rule whose flush failed is flushed with the next grant, even when that
// grant gives it the same access again
func TestGrantRetriesUnflushedRule(t *testing.T) {
	var (
		req   = require.New(t)
		ctx   = context.Background()
		store = &failingRulesStore{}
		svc   = NewService(zap.NewNop(), store)
		rule  = func() *Rule { return AllowRule(1, "corteza::compose:record/1/2/3", "read") }
	)

	svc.setRules(RuleSet{AllowRule(9, "corteza::compose:record/1/2/9", "read")})

	req.Error(svc.Grant(ctx, rule()))
	req.NoError(svc.Grant(ctx, rule()))

	req.Len(store.upserted, 1)
	req.Equal("corteza::compose:record/1/2/3", store.upserted[0].Resource)
}
