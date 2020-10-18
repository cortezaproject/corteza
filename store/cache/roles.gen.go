package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache.gen.go.tpl
// Definitions: store/roles.yaml
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

func (c Cache) cacheRole(res *types.Role) {
	var (
		ttl  time.Duration = 0
		cost int64         = 1
	)

	if c.roles.SetWithTTL(res.ID, res, cost, ttl) {
		for _, ikey := range c.roleIndexes(res) {
			c.roles.SetWithTTL(ikey, res.ID, cost, ttl)
		}
	}
}

func (c Cache) getCachedRoleByKey(ikey string) (interface{}, bool) {
	if val, found := c.roles.Get(ikey); found {
		if id, ok := val.(uint64); ok {
			return c.roles.Get(id)
		}

		c.roles.Del(val)
	}

	return nil, false
}

func (c Cache) roleIndexes(res *types.Role) []string {
	return []string{
		iKey("Handle",
			store.PreprocessValue(res.Handle, "lower"),
		),
		iKey("Name",
			store.PreprocessValue(res.Name, ""),
		),
	}
}

// LookupRoleByID searches for role by ID
//
// It returns role even if deleted or suspended
func (c Cache) LookupRoleByID(ctx context.Context, id uint64) (*types.Role, error) {

	if val, found := c.roles.Get(id); found {
		if res, ok := val.(*types.Role); ok {
			return res, nil
		}

		c.roles.Del(id)
	}

	if res, err := c.Storer.LookupRoleByID(ctx, id); err != nil {
		return nil, err
	} else {
		c.cacheRole(res)
		return res, nil
	}
}

// LookupRoleByHandle searches for role by its handle
//
// It returns only valid roles (not deleted, not archived)
func (c Cache) LookupRoleByHandle(ctx context.Context, handle string) (*types.Role, error) {

	key := iKey(
		"Roles",
		"Handle",
		store.PreprocessValue(handle, "lower"),
	)

	if val, found := c.getCachedRoleByKey(key); found {
		if res, ok := val.(*types.Role); ok {
			return res, nil
		}

		c.roles.Del(key)
	}

	if res, err := c.Storer.LookupRoleByHandle(ctx, handle); err != nil {
		return nil, err
	} else {
		c.cacheRole(res)
		return res, nil
	}
}

// LookupRoleByName searches for role by its name
//
// It returns only valid roles (not deleted, not archived)
func (c Cache) LookupRoleByName(ctx context.Context, name string) (*types.Role, error) {

	key := iKey(
		"Roles",
		"Name",
		store.PreprocessValue(name, ""),
	)

	if val, found := c.getCachedRoleByKey(key); found {
		if res, ok := val.(*types.Role); ok {
			return res, nil
		}

		c.roles.Del(key)
	}

	if res, err := c.Storer.LookupRoleByName(ctx, name); err != nil {
		return nil, err
	} else {
		c.cacheRole(res)
		return res, nil
	}
}

// CreateRole updates cache and forwards call to next configured store
func (c Cache) CreateRole(ctx context.Context, rr ...*types.Role) (err error) {
	for _, res := range rr {
		if err = c.Storer.CreateRole(ctx, res); err != nil {
			return err
		}

		c.cacheRole(res)
	}

	return nil
}

// UpdateRole updates cache and forwards call to next configured store
func (c Cache) UpdateRole(ctx context.Context, rr ...*types.Role) error {
	for _, res := range rr {
		if err := c.Storer.UpdateRole(ctx, res); err != nil {
			return err
		}

		c.cacheRole(res)
	}

	return nil
}

// UpsertRole updates cache and forwards call to next configured store
func (c Cache) UpsertRole(ctx context.Context, rr ...*types.Role) (err error) {
	for _, res := range rr {
		if err = c.Storer.UpsertRole(ctx, res); err != nil {
			return err
		}

		c.cacheRole(res)
	}

	return nil
}

// DeleteRole Deletes one or more rows from roles table
func (c Cache) DeleteRole(ctx context.Context, rr ...*types.Role) (err error) {
	for _, res := range rr {
		if err = c.Storer.DeleteRole(ctx, res); err != nil {
			return
		}

		c.roles.Del(res.ID)
		for _, key := range c.roleIndexes(res) {
			c.roles.Del(key)
		}
	}

	return nil
}

// DeleteRoleByID Deletes row from the roles table
func (c Cache) DeleteRoleByID(ctx context.Context, ID uint64) error {
	if err := c.Storer.DeleteRoleByID(ctx, ID); err != nil {
		return err
	}

	c.roles.Del(ID)
	return nil
}

// TruncateRoles Deletes all rows from the roles table
func (c Cache) TruncateRoles(ctx context.Context) error {
	if err := c.Storer.TruncateRoles(ctx); err != nil {
		return err

	}

	c.roles.Clear()
	return nil
}
