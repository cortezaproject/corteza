package cache

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_cache.gen.go.tpl
// Definitions: store/users.yaml
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

func (c Cache) cacheUser(res *types.User) {
	var (
		ttl  time.Duration = 0
		cost int64         = 1
	)

	if c.users.SetWithTTL(res.ID, res, cost, ttl) {
		for _, ikey := range c.userIndexes(res) {
			c.users.SetWithTTL(ikey, res.ID, cost, ttl)
		}
	}
}

func (c Cache) getCachedUserByKey(ikey string) (interface{}, bool) {
	if val, found := c.users.Get(ikey); found {
		if id, ok := val.(uint64); ok {
			return c.users.Get(id)
		}

		c.users.Del(val)
	}

	return nil, false
}

func (c Cache) userIndexes(res *types.User) []string {
	return []string{
		iKey("Email",
			store.PreprocessValue(res.Email, "lower"),
		),
		iKey("Handle",
			store.PreprocessValue(res.Handle, "lower"),
		),
		iKey("Username",
			store.PreprocessValue(res.Username, "lower"),
		),
	}
}

// LookupUserByID searches for user by ID
//
// It returns user even if deleted or suspended
func (c Cache) LookupUserByID(ctx context.Context, id uint64) (*types.User, error) {

	if val, found := c.users.Get(id); found {
		if res, ok := val.(*types.User); ok {
			return res, nil
		}

		c.users.Del(id)
	}

	if res, err := c.Storer.LookupUserByID(ctx, id); err != nil {
		return nil, err
	} else {
		c.cacheUser(res)
		return res, nil
	}
}

// LookupUserByEmail searches for user by their email
//
// It returns only valid users (not deleted, not suspended)
func (c Cache) LookupUserByEmail(ctx context.Context, email string) (*types.User, error) {

	key := iKey(
		"Users",
		"Email",
		store.PreprocessValue(email, "lower"),
	)

	if val, found := c.getCachedUserByKey(key); found {
		if res, ok := val.(*types.User); ok {
			return res, nil
		}

		c.users.Del(key)
	}

	if res, err := c.Storer.LookupUserByEmail(ctx, email); err != nil {
		return nil, err
	} else {
		c.cacheUser(res)
		return res, nil
	}
}

// LookupUserByHandle searches for user by their email
//
// It returns only valid users (not deleted, not suspended)
func (c Cache) LookupUserByHandle(ctx context.Context, handle string) (*types.User, error) {

	key := iKey(
		"Users",
		"Handle",
		store.PreprocessValue(handle, "lower"),
	)

	if val, found := c.getCachedUserByKey(key); found {
		if res, ok := val.(*types.User); ok {
			return res, nil
		}

		c.users.Del(key)
	}

	if res, err := c.Storer.LookupUserByHandle(ctx, handle); err != nil {
		return nil, err
	} else {
		c.cacheUser(res)
		return res, nil
	}
}

// LookupUserByUsername searches for user by their username
//
// It returns only valid users (not deleted, not suspended)
func (c Cache) LookupUserByUsername(ctx context.Context, username string) (*types.User, error) {

	key := iKey(
		"Users",
		"Username",
		store.PreprocessValue(username, "lower"),
	)

	if val, found := c.getCachedUserByKey(key); found {
		if res, ok := val.(*types.User); ok {
			return res, nil
		}

		c.users.Del(key)
	}

	if res, err := c.Storer.LookupUserByUsername(ctx, username); err != nil {
		return nil, err
	} else {
		c.cacheUser(res)
		return res, nil
	}
}

// CreateUser updates cache and forwards call to next configured store
func (c Cache) CreateUser(ctx context.Context, rr ...*types.User) (err error) {
	for _, res := range rr {
		if err = c.Storer.CreateUser(ctx, res); err != nil {
			return err
		}

		c.cacheUser(res)
	}

	return nil
}

// UpdateUser updates cache and forwards call to next configured store
func (c Cache) UpdateUser(ctx context.Context, rr ...*types.User) error {
	for _, res := range rr {
		if err := c.Storer.UpdateUser(ctx, res); err != nil {
			return err
		}

		c.cacheUser(res)
	}

	return nil
}

// UpsertUser updates cache and forwards call to next configured store
func (c Cache) UpsertUser(ctx context.Context, rr ...*types.User) (err error) {
	for _, res := range rr {
		if err = c.Storer.UpsertUser(ctx, res); err != nil {
			return err
		}

		c.cacheUser(res)
	}

	return nil
}

// DeleteUser Deletes one or more rows from users table
func (c Cache) DeleteUser(ctx context.Context, rr ...*types.User) (err error) {
	for _, res := range rr {
		if err = c.Storer.DeleteUser(ctx, res); err != nil {
			return
		}

		c.users.Del(res.ID)
		for _, key := range c.userIndexes(res) {
			c.users.Del(key)
		}
	}

	return nil
}

// DeleteUserByID Deletes row from the users table
func (c Cache) DeleteUserByID(ctx context.Context, ID uint64) error {
	if err := c.Storer.DeleteUserByID(ctx, ID); err != nil {
		return err
	}

	c.users.Del(ID)
	return nil
}

// TruncateUsers Deletes all rows from the users table
func (c Cache) TruncateUsers(ctx context.Context) error {
	if err := c.Storer.TruncateUsers(ctx); err != nil {
		return err

	}

	c.users.Clear()
	return nil
}
