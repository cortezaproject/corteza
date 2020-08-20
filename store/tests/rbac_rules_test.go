package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"testing"
)

func testRbacRules(t *testing.T, s store.RbacRules) {
	var (
		ctx = context.Background()
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateRbacRules(ctx))
		req.NoError(s.CreateRbacRule(ctx, permissions.AllowRule(42, "res1", "op1")))
	})

	t.Run("update", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateRbacRules(ctx))
		req.NoError(s.UpdateRbacRule(ctx, permissions.AllowRule(42, "res1", "op1")))
	})

	t.Run("upsert", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateRbacRules(ctx))
		req.NoError(s.UpsertRbacRule(ctx, permissions.AllowRule(42, "res1", "op1")))
		set, _, err := s.SearchRbacRules(ctx, permissions.RuleFilter{})
		req.NoError(err)
		req.Len(set, 1)
		req.True(set[0].Access == permissions.Allow)

		req.NoError(s.UpsertRbacRule(ctx, permissions.DenyRule(42, "res1", "op1")))
		set, _, err = s.SearchRbacRules(ctx, permissions.RuleFilter{})
		req.NoError(err)
		req.Len(set, 1)
		req.True(set[0].Access == permissions.Deny)
	})

	t.Run("search", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateRbacRules(ctx))
		req.NoError(s.CreateRbacRule(ctx,
			permissions.AllowRule(42, "res1", "op1"),
			permissions.AllowRule(42, "res2", "op2"),
		))

		set, _, err := s.SearchRbacRules(ctx, permissions.RuleFilter{})
		req.NoError(err)
		req.Len(set, 2)
	})
}
