package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"testing"
)

func testRbacRules(t *testing.T, s rbacRulesStore) {
	var (
		ctx = context.Background()
		req = require.New(t)

		rule *permissions.Rule
	)

	t.Run("create", func(t *testing.T) {
		rule = &permissions.Rule{
			RoleID:    42,
			Resource:  permissions.Resource("res"),
			Operation: permissions.Operation("op"),
			Access:    permissions.Allow,
		}
		req.NoError(s.CreateRbacRule(ctx, rule))
	})

	t.Run("update", func(t *testing.T) {
		rule = &permissions.Rule{
			RoleID:    42,
			Resource:  permissions.Resource("res"),
			Operation: permissions.Operation("op"),
			Access:    permissions.Allow,
		}
		req.NoError(s.UpdateRbacRule(ctx, rule))
	})

	t.Run("search", func(t *testing.T) {
		set, _, err := s.SearchRbacRules(ctx, permissions.RuleFilter{})
		req.NoError(err)
		req.Len(set, 1)
	})
}
