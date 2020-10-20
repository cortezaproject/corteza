package cache

// Custom cache implementation for role membership
//
// Caching role members and user memberships

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/system/types"
	"time"
)

var _ = errors.Is

// cacheRoleMembership handles both - membership and members!
func (c Cache) cacheRoleRelationships(i uint64, mm []uint64) {
	var (
		ttl  time.Duration = 0
		cost int64         = 1
	)

	c.roleMembers.SetWithTTL(i, mm, cost, ttl)
}

func (c Cache) getRoleRelationships(i uint64) ([]uint64, bool) {
	if val, found := c.roleMembers.Get(i); found {
		if mm, ok := val.([]uint64); ok {
			return mm, true
		}
	}

	c.roleMembers.Del(i)
	return nil, false
}

func (c Cache) SearchRoleMembers(ctx context.Context, roleID uint64) ([]uint64, error) {
	if mm, found := c.getRoleRelationships(roleID); found {
		return mm, nil
	}

	if mm, err := c.Storer.SearchRoleMembers(ctx, roleID); err != nil {
		return nil, err
	} else {
		c.cacheRoleRelationships(roleID, mm)
		return mm, nil
	}
}

func (c Cache) SearchUserMemberships(ctx context.Context, userID uint64) ([]uint64, error) {
	if mm, found := c.getRoleRelationships(userID); found {
		return mm, nil
	}

	if mm, err := c.Storer.SearchUserMemberships(ctx, userID); err != nil {
		return nil, err
	} else {
		c.cacheRoleRelationships(userID, mm)
		return mm, nil
	}
}

// CreateRoleMember updates cache and forwards call to next configured store
func (c Cache) CreateRoleMember(ctx context.Context, rr ...*types.RoleMember) (err error) {
	for _, res := range rr {
		if err = c.Storer.CreateRoleMember(ctx, res); err != nil {
			return err
		}

		// remove all cache entries to this role/user
		c.roleMembers.Del(res.RoleID)
		c.roleMembers.Del(res.UserID)
	}

	return nil
}

// DeleteRoleMember Deletes one or more rows from role_members table
func (c Cache) DeleteRoleMember(ctx context.Context, rr ...*types.RoleMember) (err error) {
	for _, res := range rr {
		if err = c.Storer.DeleteRoleMember(ctx, res); err != nil {
			return
		}

		// remove all cache entries to this role/user
		c.roleMembers.Del(res.RoleID)
		c.roleMembers.Del(res.UserID)
	}

	return nil
}

// DeleteRoleMemberByUserIDRoleID Deletes row from the role_members table
func (c Cache) DeleteRoleMemberByUserIDRoleID(ctx context.Context, userID uint64, roleID uint64) error {
	if err := c.Storer.DeleteRoleMemberByUserIDRoleID(ctx, userID, roleID); err != nil {
		return err
	}

	// remove all cache entries to this role/user
	c.roleMembers.Del(userID)
	c.roleMembers.Del(roleID)
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
