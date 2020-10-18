package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache.gen.go.tpl
// Definitions: store/compose_namespaces.yaml
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

func (c Cache) cacheComposeNamespace(res *types.Namespace) {
	var (
		ttl  time.Duration = 0
		cost int64         = 1
	)

	if c.composeNamespaces.SetWithTTL(res.ID, res, cost, ttl) {
		for _, ikey := range c.composeNamespaceIndexes(res) {
			c.composeNamespaces.SetWithTTL(ikey, res.ID, cost, ttl)
		}
	}
}

func (c Cache) getCachedComposeNamespaceByKey(ikey string) (interface{}, bool) {
	if val, found := c.composeNamespaces.Get(ikey); found {
		if id, ok := val.(uint64); ok {
			return c.composeNamespaces.Get(id)
		}

		c.composeNamespaces.Del(val)
	}

	return nil, false
}

func (c Cache) composeNamespaceIndexes(res *types.Namespace) []string {
	return []string{
		iKey("Slug",
			store.PreprocessValue(res.Slug, "lower"),
		),
	}
}

// LookupComposeNamespaceBySlug searches for namespace by slug (case-insensitive)
func (c Cache) LookupComposeNamespaceBySlug(ctx context.Context, slug string) (*types.Namespace, error) {

	key := iKey(
		"ComposeNamespaces",
		"Slug",
		store.PreprocessValue(slug, "lower"),
	)

	if val, found := c.getCachedComposeNamespaceByKey(key); found {
		if res, ok := val.(*types.Namespace); ok {
			return res, nil
		}

		c.composeNamespaces.Del(key)
	}

	if res, err := c.Storer.LookupComposeNamespaceBySlug(ctx, slug); err != nil {
		return nil, err
	} else {
		c.cacheComposeNamespace(res)
		return res, nil
	}
}

// LookupComposeNamespaceByID searches for compose namespace by ID
//
// It returns compose namespace even if deleted
func (c Cache) LookupComposeNamespaceByID(ctx context.Context, id uint64) (*types.Namespace, error) {

	if val, found := c.composeNamespaces.Get(id); found {
		if res, ok := val.(*types.Namespace); ok {
			return res, nil
		}

		c.composeNamespaces.Del(id)
	}

	if res, err := c.Storer.LookupComposeNamespaceByID(ctx, id); err != nil {
		return nil, err
	} else {
		c.cacheComposeNamespace(res)
		return res, nil
	}
}

// CreateComposeNamespace updates cache and forwards call to next configured store
func (c Cache) CreateComposeNamespace(ctx context.Context, rr ...*types.Namespace) (err error) {
	for _, res := range rr {
		if err = c.Storer.CreateComposeNamespace(ctx, res); err != nil {
			return err
		}

		c.cacheComposeNamespace(res)
	}

	return nil
}

// UpdateComposeNamespace updates cache and forwards call to next configured store
func (c Cache) UpdateComposeNamespace(ctx context.Context, rr ...*types.Namespace) error {
	for _, res := range rr {
		if err := c.Storer.UpdateComposeNamespace(ctx, res); err != nil {
			return err
		}

		c.cacheComposeNamespace(res)
	}

	return nil
}

// UpsertComposeNamespace updates cache and forwards call to next configured store
func (c Cache) UpsertComposeNamespace(ctx context.Context, rr ...*types.Namespace) (err error) {
	for _, res := range rr {
		if err = c.Storer.UpsertComposeNamespace(ctx, res); err != nil {
			return err
		}

		c.cacheComposeNamespace(res)
	}

	return nil
}

// DeleteComposeNamespace Deletes one or more rows from compose_namespace table
func (c Cache) DeleteComposeNamespace(ctx context.Context, rr ...*types.Namespace) (err error) {
	for _, res := range rr {
		if err = c.Storer.DeleteComposeNamespace(ctx, res); err != nil {
			return
		}

		c.composeNamespaces.Del(res.ID)
		for _, key := range c.composeNamespaceIndexes(res) {
			c.composeNamespaces.Del(key)
		}
	}

	return nil
}

// DeleteComposeNamespaceByID Deletes row from the compose_namespace table
func (c Cache) DeleteComposeNamespaceByID(ctx context.Context, ID uint64) error {
	if err := c.Storer.DeleteComposeNamespaceByID(ctx, ID); err != nil {
		return err
	}

	c.composeNamespaces.Del(ID)
	return nil
}

// TruncateComposeNamespaces Deletes all rows from the compose_namespace table
func (c Cache) TruncateComposeNamespaces(ctx context.Context) error {
	if err := c.Storer.TruncateComposeNamespaces(ctx); err != nil {
		return err

	}

	c.composeNamespaces.Clear()
	return nil
}
