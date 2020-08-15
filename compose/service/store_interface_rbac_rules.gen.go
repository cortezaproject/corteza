package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/rbac_rules.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	rbacRulesStore interface {
		SearchRbacRules(ctx context.Context, f permissions.RuleFilter) (permissions.RuleSet, permissions.RuleFilter, error)
		CreateRbacRule(ctx context.Context, rr ...*permissions.Rule) error
		UpdateRbacRule(ctx context.Context, rr ...*permissions.Rule) error
		PartialUpdateRbacRule(ctx context.Context, onlyColumns []string, rr ...*permissions.Rule) error
		RemoveRbacRule(ctx context.Context, rr ...*permissions.Rule) error
		RemoveRbacRuleByRoleIDResourceOperation(ctx context.Context, roleID uint64, resource string, operation string) error

		TruncateRbacRules(ctx context.Context) error
	}
)
