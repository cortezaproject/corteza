package store

// This file is auto-generated.
//
// Template:    pkg/codegen/assets/store_base.gen.go.tpl
// Definitions: store/rbac_rules.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	RbacRules interface {
		SearchRbacRules(ctx context.Context, f permissions.RuleFilter) (permissions.RuleSet, permissions.RuleFilter, error)

		CreateRbacRule(ctx context.Context, rr ...*permissions.Rule) error

		UpdateRbacRule(ctx context.Context, rr ...*permissions.Rule) error

		UpsertRbacRule(ctx context.Context, rr ...*permissions.Rule) error

		DeleteRbacRule(ctx context.Context, rr ...*permissions.Rule) error
		DeleteRbacRuleByRoleIDResourceOperation(ctx context.Context, roleID uint64, resource string, operation string) error

		TruncateRbacRules(ctx context.Context) error
	}
)

var _ *permissions.Rule
var _ context.Context

// SearchRbacRules returns all matching RbacRules from store
func SearchRbacRules(ctx context.Context, s RbacRules, f permissions.RuleFilter) (permissions.RuleSet, permissions.RuleFilter, error) {
	return s.SearchRbacRules(ctx, f)
}

// CreateRbacRule creates one or more RbacRules in store
func CreateRbacRule(ctx context.Context, s RbacRules, rr ...*permissions.Rule) error {
	return s.CreateRbacRule(ctx, rr...)
}

// UpdateRbacRule updates one or more (existing) RbacRules in store
func UpdateRbacRule(ctx context.Context, s RbacRules, rr ...*permissions.Rule) error {
	return s.UpdateRbacRule(ctx, rr...)
}

// UpsertRbacRule creates new or updates existing one or more RbacRules in store
func UpsertRbacRule(ctx context.Context, s RbacRules, rr ...*permissions.Rule) error {
	return s.UpsertRbacRule(ctx, rr...)
}

// DeleteRbacRule Deletes one or more RbacRules from store
func DeleteRbacRule(ctx context.Context, s RbacRules, rr ...*permissions.Rule) error {
	return s.DeleteRbacRule(ctx, rr...)
}

// DeleteRbacRuleByRoleIDResourceOperation Deletes RbacRule from store
func DeleteRbacRuleByRoleIDResourceOperation(ctx context.Context, s RbacRules, roleID uint64, resource string, operation string) error {
	return s.DeleteRbacRuleByRoleIDResourceOperation(ctx, roleID, resource, operation)
}

// TruncateRbacRules Deletes all RbacRules from store
func TruncateRbacRules(ctx context.Context, s RbacRules) error {
	return s.TruncateRbacRules(ctx)
}
