package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache.gen.go.tpl
// Definitions: store/compose_pages.yaml
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

func (c Cache) cacheComposePage(res *types.Page) {
	var (
		ttl  time.Duration = 0
		cost int64         = 1
	)

	if c.composePages.SetWithTTL(res.ID, res, cost, ttl) {
		for _, ikey := range c.composePageIndexes(res) {
			c.composePages.SetWithTTL(ikey, res.ID, cost, ttl)
		}
	}
}

func (c Cache) getCachedComposePageByKey(ikey string) (interface{}, bool) {
	if val, found := c.composePages.Get(ikey); found {
		if id, ok := val.(uint64); ok {
			return c.composePages.Get(id)
		}

		c.composePages.Del(val)
	}

	return nil, false
}

func (c Cache) composePageIndexes(res *types.Page) []string {
	return []string{
		iKey("NamespaceIDHandle",
			store.PreprocessValue(res.NamespaceID, ""),
			store.PreprocessValue(res.Handle, "lower"),
		),
		iKey("NamespaceIDModuleID",
			store.PreprocessValue(res.NamespaceID, ""),
			store.PreprocessValue(res.ModuleID, ""),
		),
	}
}

// LookupComposePageByNamespaceIDHandle searches for page by handle (case-insensitive)
func (c Cache) LookupComposePageByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Page, error) {

	key := iKey(
		"ComposePages",
		"NamespaceIDHandle",
		store.PreprocessValue(namespace_id, ""),
		store.PreprocessValue(handle, "lower"),
	)

	if val, found := c.getCachedComposePageByKey(key); found {
		if res, ok := val.(*types.Page); ok {
			return res, nil
		}

		c.composePages.Del(key)
	}

	if res, err := c.Storer.LookupComposePageByNamespaceIDHandle(ctx, namespace_id, handle); err != nil {
		return nil, err
	} else {
		c.cacheComposePage(res)
		return res, nil
	}
}

// LookupComposePageByNamespaceIDModuleID searches for page by moduleID
func (c Cache) LookupComposePageByNamespaceIDModuleID(ctx context.Context, namespace_id uint64, module_id uint64) (*types.Page, error) {

	key := iKey(
		"ComposePages",
		"NamespaceIDModuleID",
		store.PreprocessValue(namespace_id, ""),
		store.PreprocessValue(module_id, ""),
	)

	if val, found := c.getCachedComposePageByKey(key); found {
		if res, ok := val.(*types.Page); ok {
			return res, nil
		}

		c.composePages.Del(key)
	}

	if res, err := c.Storer.LookupComposePageByNamespaceIDModuleID(ctx, namespace_id, module_id); err != nil {
		return nil, err
	} else {
		c.cacheComposePage(res)
		return res, nil
	}
}

// LookupComposePageByID searches for compose page by ID
//
// It returns compose page even if deleted
func (c Cache) LookupComposePageByID(ctx context.Context, id uint64) (*types.Page, error) {

	if val, found := c.composePages.Get(id); found {
		if res, ok := val.(*types.Page); ok {
			return res, nil
		}

		c.composePages.Del(id)
	}

	if res, err := c.Storer.LookupComposePageByID(ctx, id); err != nil {
		return nil, err
	} else {
		c.cacheComposePage(res)
		return res, nil
	}
}

// CreateComposePage updates cache and forwards call to next configured store
func (c Cache) CreateComposePage(ctx context.Context, rr ...*types.Page) (err error) {
	for _, res := range rr {
		if err = c.Storer.CreateComposePage(ctx, res); err != nil {
			return err
		}

		c.cacheComposePage(res)
	}

	return nil
}

// UpdateComposePage updates cache and forwards call to next configured store
func (c Cache) UpdateComposePage(ctx context.Context, rr ...*types.Page) error {
	for _, res := range rr {
		if err := c.Storer.UpdateComposePage(ctx, res); err != nil {
			return err
		}

		c.cacheComposePage(res)
	}

	return nil
}

// UpsertComposePage updates cache and forwards call to next configured store
func (c Cache) UpsertComposePage(ctx context.Context, rr ...*types.Page) (err error) {
	for _, res := range rr {
		if err = c.Storer.UpsertComposePage(ctx, res); err != nil {
			return err
		}

		c.cacheComposePage(res)
	}

	return nil
}

// DeleteComposePage Deletes one or more rows from compose_page table
func (c Cache) DeleteComposePage(ctx context.Context, rr ...*types.Page) (err error) {
	for _, res := range rr {
		if err = c.Storer.DeleteComposePage(ctx, res); err != nil {
			return
		}

		c.composePages.Del(res.ID)
		for _, key := range c.composePageIndexes(res) {
			c.composePages.Del(key)
		}
	}

	return nil
}

// DeleteComposePageByID Deletes row from the compose_page table
func (c Cache) DeleteComposePageByID(ctx context.Context, ID uint64) error {
	if err := c.Storer.DeleteComposePageByID(ctx, ID); err != nil {
		return err
	}

	c.composePages.Del(ID)
	return nil
}

// TruncateComposePages Deletes all rows from the compose_page table
func (c Cache) TruncateComposePages(ctx context.Context) error {
	if err := c.Storer.TruncateComposePages(ctx); err != nil {
		return err

	}

	c.composePages.Clear()
	return nil
}
