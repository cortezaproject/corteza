package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache.gen.go.tpl
// Definitions: store/compose_modules.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/store"
	"time"
)

var _ = errors.Is

func (c Cache) cacheComposeModule(res *types.Module) {
	var (
		ttl  time.Duration = 0
		cost int64         = 1
	)

	if c.composeModules.SetWithTTL(res.ID, res, cost, ttl) {
		for _, ikey := range c.composeModuleIndexes(res) {
			c.composeModules.SetWithTTL(ikey, res.ID, cost, ttl)
		}
	}
}

func (c Cache) getCachedComposeModuleByKey(ikey string) (interface{}, bool) {
	if val, found := c.composeModules.Get(ikey); found {
		if id, ok := val.(uint64); ok {
			return c.composeModules.Get(id)
		}

		c.composeModules.Del(val)
	}

	return nil, false
}

func (c Cache) composeModuleIndexes(res *types.Module) []string {
	return []string{
		iKey("NamespaceIDHandle",
			store.PreprocessValue(res.NamespaceID, ""),
			store.PreprocessValue(res.Handle, "lower"),
		),
		iKey("NamespaceIDName",
			store.PreprocessValue(res.NamespaceID, ""),
			store.PreprocessValue(res.Name, "lower"),
		),
	}
}

// LookupComposeModuleByNamespaceIDHandle searches for compose module by handle (case-insensitive)
func (c Cache) LookupComposeModuleByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Module, error) {

	key := iKey(
		"ComposeModules",
		"NamespaceIDHandle",
		store.PreprocessValue(namespace_id, ""),
		store.PreprocessValue(handle, "lower"),
	)

	if val, found := c.getCachedComposeModuleByKey(key); found {
		if res, ok := val.(*types.Module); ok {
			return res, nil
		}

		c.composeModules.Del(key)
	}

	if res, err := c.Storer.LookupComposeModuleByNamespaceIDHandle(ctx, namespace_id, handle); err != nil {
		return nil, err
	} else {
		c.cacheComposeModule(res)
		return res, nil
	}
}

// LookupComposeModuleByNamespaceIDName searches for compose module by name (case-insensitive)
func (c Cache) LookupComposeModuleByNamespaceIDName(ctx context.Context, namespace_id uint64, name string) (*types.Module, error) {

	key := iKey(
		"ComposeModules",
		"NamespaceIDName",
		store.PreprocessValue(namespace_id, ""),
		store.PreprocessValue(name, "lower"),
	)

	if val, found := c.getCachedComposeModuleByKey(key); found {
		if res, ok := val.(*types.Module); ok {
			return res, nil
		}

		c.composeModules.Del(key)
	}

	if res, err := c.Storer.LookupComposeModuleByNamespaceIDName(ctx, namespace_id, name); err != nil {
		return nil, err
	} else {
		c.cacheComposeModule(res)
		return res, nil
	}
}

// LookupComposeModuleByID searches for compose module by ID
//
// It returns compose module even if deleted
func (c Cache) LookupComposeModuleByID(ctx context.Context, id uint64) (*types.Module, error) {

	if val, found := c.composeModules.Get(id); found {
		if res, ok := val.(*types.Module); ok {
			return res, nil
		}

		c.composeModules.Del(id)
	}

	if res, err := c.Storer.LookupComposeModuleByID(ctx, id); err != nil {
		return nil, err
	} else {
		c.cacheComposeModule(res)
		return res, nil
	}
}

// CreateComposeModule updates cache and forwards call to next configured store
func (c Cache) CreateComposeModule(ctx context.Context, rr ...*types.Module) (err error) {
	for _, res := range rr {
		if err = c.Storer.CreateComposeModule(ctx, res); err != nil {
			return err
		}

		c.cacheComposeModule(res)
	}

	return nil
}

// UpdateComposeModule updates cache and forwards call to next configured store
func (c Cache) UpdateComposeModule(ctx context.Context, rr ...*types.Module) error {
	for _, res := range rr {
		if err := c.Storer.UpdateComposeModule(ctx, res); err != nil {
			return err
		}

		c.cacheComposeModule(res)
	}

	return nil
}

// UpsertComposeModule updates cache and forwards call to next configured store
func (c Cache) UpsertComposeModule(ctx context.Context, rr ...*types.Module) (err error) {
	for _, res := range rr {
		if err = c.Storer.UpsertComposeModule(ctx, res); err != nil {
			return err
		}

		c.cacheComposeModule(res)
	}

	return nil
}

// DeleteComposeModule Deletes one or more rows from compose_module table
func (c Cache) DeleteComposeModule(ctx context.Context, rr ...*types.Module) (err error) {
	for _, res := range rr {
		if err = c.Storer.DeleteComposeModule(ctx, res); err != nil {
			return
		}

		c.composeModules.Del(res.ID)
		for _, key := range c.composeModuleIndexes(res) {
			c.composeModules.Del(key)
		}
	}

	return nil
}

// DeleteComposeModuleByID Deletes row from the compose_module table
func (c Cache) DeleteComposeModuleByID(ctx context.Context, ID uint64) error {
	if err := c.Storer.DeleteComposeModuleByID(ctx, ID); err != nil {
		return err
	}

	c.composeModules.Del(ID)
	return nil
}

// TruncateComposeModules Deletes all rows from the compose_module table
func (c Cache) TruncateComposeModules(ctx context.Context) error {
	if err := c.Storer.TruncateComposeModules(ctx); err != nil {
		return err

	}

	c.composeModules.Clear()
	return nil
}
