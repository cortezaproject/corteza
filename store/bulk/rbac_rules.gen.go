package bulk

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
// Definitions file that controls how this file is generated:
// store/rbac_rules.yaml

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	rbacRuleCreate struct {
		Done chan struct{}
		res  *permissions.Rule
		err  error
	}

	rbacRuleUpdate struct {
		Done chan struct{}
		res  *permissions.Rule
		err  error
	}

	rbacRuleRemove struct {
		Done chan struct{}
		res  *permissions.Rule
		err  error
	}
)

// CreateRbacRule creates a new RbacRule
// create job that can be pushed to store's transaction handler
func CreateRbacRule(res *permissions.Rule) *rbacRuleCreate {
	return &rbacRuleCreate{res: res}
}

// Do Executes rbacRuleCreate job
func (j *rbacRuleCreate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.CreateRbacRule(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// UpdateRbacRule creates a new RbacRule
// update job that can be pushed to store's transaction handler
func UpdateRbacRule(res *permissions.Rule) *rbacRuleUpdate {
	return &rbacRuleUpdate{res: res}
}

// Do Executes rbacRuleUpdate job
func (j *rbacRuleUpdate) Do(ctx context.Context, s storeInterface) error {
	j.err = s.UpdateRbacRule(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}

// RemoveRbacRule creates a new RbacRule
// remove job that can be pushed to store's transaction handler
func RemoveRbacRule(res *permissions.Rule) *rbacRuleRemove {
	return &rbacRuleRemove{res: res}
}

// Do Executes rbacRuleRemove job
func (j *rbacRuleRemove) Do(ctx context.Context, s storeInterface) error {
	j.err = s.RemoveRbacRule(ctx, j.res)
	j.Done <- struct{}{}
	return j.err
}
