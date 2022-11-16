package store

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestUser_Merger(t *testing.T) {
	req := require.New(t)

	now := time.Time{}
	nowP := &time.Time{}

	empty := &types.User{}
	full := &types.User{
		Username:    "username",
		Email:       "email",
		Name:        "name",
		Handle:      "handle",
		Kind:        types.SystemUser,
		Meta:        &types.UserMeta{},
		CreatedAt:   now,
		UpdatedAt:   nowP,
		SuspendedAt: nowP,
		DeletedAt:   nowP,
	}

	t.Run("merge on empty", func(t *testing.T) {
		c := mergeUsers(empty, full)
		req.Equal("username", c.Username)
		req.Equal("email", c.Email)
		req.Equal("name", c.Name)
		req.Equal("handle", c.Handle)
		req.Equal(types.SystemUser, c.Kind)
		req.Equal(&types.UserMeta{}, c.Meta)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.SuspendedAt)
		req.Equal(nowP, c.DeletedAt)
	})

	t.Run("merge with empty", func(t *testing.T) {
		c := mergeUsers(full, empty)
		req.Equal("username", c.Username)
		req.Equal("email", c.Email)
		req.Equal("name", c.Name)
		req.Equal("handle", c.Handle)
		req.Equal(types.SystemUser, c.Kind)
		req.Equal(&types.UserMeta{}, c.Meta)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.SuspendedAt)
		req.Equal(nowP, c.DeletedAt)
	})
}
