package rbac

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap"
)

type nopRulesStore struct{}

func (nopRulesStore) SearchRbacRules(context.Context, RuleFilter) (RuleSet, RuleFilter, error) {
	return nil, RuleFilter{}, nil
}
func (nopRulesStore) UpsertRbacRule(context.Context, ...*Rule) error { return nil }
func (nopRulesStore) DeleteRbacRule(context.Context, ...*Rule) error { return nil }
func (nopRulesStore) TruncateRbacRules(context.Context) error        { return nil }

func benchGrant(b *testing.B, existing int) {
	ctx := context.Background()
	svc := NewService(zap.NewNop(), nopRulesStore{})
	rr := make(RuleSet, 0, existing)
	for i := 0; i < existing; i++ {
		rr = append(rr, AllowRule(uint64(1+i%3), fmt.Sprintf("corteza::compose:record/1/2/%d", 1000000+i), "read"))
	}
	svc.setRules(rr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = svc.Grant(ctx, AllowRule(7, fmt.Sprintf("corteza::compose:record/1/2/%d", 9000000+i), "read"))
	}
}

func BenchmarkGrant1k(b *testing.B)  { benchGrant(b, 1000) }
func BenchmarkGrant20k(b *testing.B) { benchGrant(b, 20000) }
