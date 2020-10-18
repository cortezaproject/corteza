package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache.gen.go.tpl
// Definitions: store/role_members.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"time"
)

var _ = errors.Is

func (c Cache) cacheRoleMember(res *types.RoleMember) {
	var (
		ttl  time.Duration = 0
		cost int64         = 1
	)

	if c.roleMembers.SetWithTTL(res.ID, res, cost, ttl) {
		for _, ikey := range c.roleMemberIndexes(res) {
			c.roleMembers.SetWithTTL(ikey, res.ID, cost, ttl)
		}
	}
}

func (c Cache) getCachedRoleMemberByKey(ikey string) (interface{}, bool) {
	if val, found := c.roleMembers.Get(ikey); found {
		if id, ok := val.(uint64); ok {
			return c.roleMembers.Get(id)
		}

		c.roleMembers.Del(val)
	}

	return nil, false
}

func (c Cache) roleMemberIndexes(res *types.RoleMember) []string {
	return []string{}
}

// CreateRoleMember updates cache and forwards call to next configured store
func (c Cache) CreateRoleMember(ctx context.Context, rr ...*types.RoleMember) (err error) {
	for _, res := range rr {
		if err = c.Storer.CreateRoleMember(ctx, res); err != nil {
			return err
		}

		c.cacheRoleMember(res)
	}

	return nil
}

// UpdateRoleMember updates cache and forwards call to next configured store
func (c Cache) UpdateRoleMember(ctx context.Context, rr ...*types.RoleMember) error {
	for _, res := range rr {
		if err := c.Storer.UpdateRoleMember(ctx, res); err != nil {
			return err
		}

		c.cacheRoleMember(res)
	}

	return nil
}

// UpsertRoleMember updates cache and forwards call to next configured store
func (c Cache) UpsertRoleMember(ctx context.Context, rr ...*types.RoleMember) (err error) {
	for _, res := range rr {
		if err = c.Storer.UpsertRoleMember(ctx, res); err != nil {
			return err
		}

		c.cacheRoleMember(res)
	}

	return nil
}

// DeleteRoleMember Deletes one or more rows from role_members table
func (c Cache) DeleteRoleMember(ctx context.Context, rr ...*types.RoleMember) (err error) {
	for _, res := range rr {
		if err = c.Storer.DeleteRoleMember(ctx, res); err != nil {
			return
		}

		c.roleMembers.Del(res.ID)
		for _, key := range c.roleMemberIndexes(res) {
			c.roleMembers.Del(key)
		}
	}

	return nil
}

// DeleteRoleMemberByUserIDRoleID Deletes row from the role_members table
func (c Cache) DeleteRoleMemberByUserIDRoleID(ctx context.Context, userID uint64, roleID uint64) error {
	if err := c.Storer.DeleteRoleMemberByUserIDRoleID(ctx, userID, roleID); err != nil {
		return err
	}

	c.roleMembers.Del(ID)
	return nil
}

// TruncateRoleMembers Deletes all rows from the role_members table
func (c Cache) TruncateRoleMembers(ctx context.Context) error {
	if err := c.Storer.TruncateRoleMembers(ctx); err != nil {
		return err

	}

	c.roleMembers.Clear()
	return nil
}
