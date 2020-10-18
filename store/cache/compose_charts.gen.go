package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache.gen.go.tpl
// Definitions: store/compose_charts.yaml
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

func (c Cache) cacheComposeChart(res *types.Chart) {
	var (
		ttl  time.Duration = 0
		cost int64         = 1
	)

	if c.composeCharts.SetWithTTL(res.ID, res, cost, ttl) {
		for _, ikey := range c.composeChartIndexes(res) {
			c.composeCharts.SetWithTTL(ikey, res.ID, cost, ttl)
		}
	}
}

func (c Cache) getCachedComposeChartByKey(ikey string) (interface{}, bool) {
	if val, found := c.composeCharts.Get(ikey); found {
		if id, ok := val.(uint64); ok {
			return c.composeCharts.Get(id)
		}

		c.composeCharts.Del(val)
	}

	return nil, false
}

func (c Cache) composeChartIndexes(res *types.Chart) []string {
	return []string{
		iKey("NamespaceIDHandle",
			store.PreprocessValue(res.NamespaceID, ""),
			store.PreprocessValue(res.Handle, "lower"),
		),
	}
}

// LookupComposeChartByID searches for compose chart by ID
//
// It returns compose chart even if deleted
func (c Cache) LookupComposeChartByID(ctx context.Context, id uint64) (*types.Chart, error) {

	if val, found := c.composeCharts.Get(id); found {
		if res, ok := val.(*types.Chart); ok {
			return res, nil
		}

		c.composeCharts.Del(id)
	}

	if res, err := c.Storer.LookupComposeChartByID(ctx, id); err != nil {
		return nil, err
	} else {
		c.cacheComposeChart(res)
		return res, nil
	}
}

// LookupComposeChartByNamespaceIDHandle searches for compose chart by handle (case-insensitive)
func (c Cache) LookupComposeChartByNamespaceIDHandle(ctx context.Context, namespace_id uint64, handle string) (*types.Chart, error) {

	key := iKey(
		"ComposeCharts",
		"NamespaceIDHandle",
		store.PreprocessValue(namespace_id, ""),
		store.PreprocessValue(handle, "lower"),
	)

	if val, found := c.getCachedComposeChartByKey(key); found {
		if res, ok := val.(*types.Chart); ok {
			return res, nil
		}

		c.composeCharts.Del(key)
	}

	if res, err := c.Storer.LookupComposeChartByNamespaceIDHandle(ctx, namespace_id, handle); err != nil {
		return nil, err
	} else {
		c.cacheComposeChart(res)
		return res, nil
	}
}

// CreateComposeChart updates cache and forwards call to next configured store
func (c Cache) CreateComposeChart(ctx context.Context, rr ...*types.Chart) (err error) {
	for _, res := range rr {
		if err = c.Storer.CreateComposeChart(ctx, res); err != nil {
			return err
		}

		c.cacheComposeChart(res)
	}

	return nil
}

// UpdateComposeChart updates cache and forwards call to next configured store
func (c Cache) UpdateComposeChart(ctx context.Context, rr ...*types.Chart) error {
	for _, res := range rr {
		if err := c.Storer.UpdateComposeChart(ctx, res); err != nil {
			return err
		}

		c.cacheComposeChart(res)
	}

	return nil
}

// UpsertComposeChart updates cache and forwards call to next configured store
func (c Cache) UpsertComposeChart(ctx context.Context, rr ...*types.Chart) (err error) {
	for _, res := range rr {
		if err = c.Storer.UpsertComposeChart(ctx, res); err != nil {
			return err
		}

		c.cacheComposeChart(res)
	}

	return nil
}

// DeleteComposeChart Deletes one or more rows from compose_chart table
func (c Cache) DeleteComposeChart(ctx context.Context, rr ...*types.Chart) (err error) {
	for _, res := range rr {
		if err = c.Storer.DeleteComposeChart(ctx, res); err != nil {
			return
		}

		c.composeCharts.Del(res.ID)
		for _, key := range c.composeChartIndexes(res) {
			c.composeCharts.Del(key)
		}
	}

	return nil
}

// DeleteComposeChartByID Deletes row from the compose_chart table
func (c Cache) DeleteComposeChartByID(ctx context.Context, ID uint64) error {
	if err := c.Storer.DeleteComposeChartByID(ctx, ID); err != nil {
		return err
	}

	c.composeCharts.Del(ID)
	return nil
}

// TruncateComposeCharts Deletes all rows from the compose_chart table
func (c Cache) TruncateComposeCharts(ctx context.Context) error {
	if err := c.Storer.TruncateComposeCharts(ctx); err != nil {
		return err

	}

	c.composeCharts.Clear()
	return nil
}
