package tests

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type (
	usersStoreAdt interface {
		usersStore
		CountUsers(ctx context.Context, f types.UserFilter) (uint, error)
	}
)

func testUsers(t *testing.T, tmp interface{}) {
	var (
		ctx = context.Background()
		req = require.New(t)

		//err  error
		user *types.User

		s = tmp.(usersStoreAdt)
	)

	t.Run("create", func(t *testing.T) {
		user = &types.User{
			ID:        42,
			CreatedAt: time.Now(),
			Email:     "user-crud@crust.test",
			Username:  "UserCRUD",
			Handle:    "usercrud",
		}
		req.NoError(s.CreateUser(ctx, user))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		fetched, err := s.LookupUserByID(ctx, user.ID)
		req.NoError(err)
		req.Equal(user.Email, fetched.Email)
		req.Equal(user.Username, fetched.Username)
		req.Equal(user.Handle, fetched.Handle)
		req.Equal(user.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
		req.Nil(fetched.SuspendedAt)
	})

	t.Run("update", func(t *testing.T) {
		user = &types.User{
			ID:        42,
			CreatedAt: time.Now(),
			Email:     "user-crud+2@crust.test",
			Username:  "UserCRUD+2",
			Handle:    "usercrud+2",
		}
		req.NoError(s.UpdateUser(ctx, user))
	})

	//t.Run("delete/undelete", func(t *testing.T) {
	//	ID := user.ID
	//	user, err = s.LookupUserByID(ctx, ID)
	//	req.NoError(err)
	//
	//	req.NoError(s.DeleteUserByID(ctx, ID))
	//	user, err = s.LookupUserByID(ctx, ID)
	//	req.NoError(err)
	//	req.NotNil(user.DeletedAt)
	//
	//	req.NoError(s.UndeleteUserByID(ctx, ID))
	//	user, err = s.LookupUserByID(ctx, ID)
	//	req.NoError(err)
	//	req.Nil(user.DeletedAt)
	//})
	//
	//t.Run("suspend/suspend", func(t *testing.T) {
	//	ID := user.ID
	//	req.NoError(s.SuspendUserByID(ctx, ID))
	//	user, err = s.LookupUserByID(ctx, ID)
	//	req.NoError(err)
	//	req.NotNil(user.SuspendedAt)
	//
	//	req.NoError(s.UnsuspendUserByID(ctx, ID))
	//	user, err = s.LookupUserByID(ctx, ID)
	//	req.NoError(err)
	//	req.Nil(user.SuspendedAt)
	//})

	t.Run("lookup by email", func(t *testing.T) {
		fetched, err := s.LookupUserByEmail(ctx, user.Email)
		req.NoError(err)
		req.Equal(user.Email, fetched.Email)
	})

	t.Run("lookup by handle", func(t *testing.T) {
		fetched, err := s.LookupUserByHandle(ctx, user.Handle)
		req.NoError(err)
		req.Equal(user.ID, fetched.ID)
	})

	t.Run("lookup by nonexisting handle", func(t *testing.T) {
		fetched, err := s.LookupUserByHandle(ctx, "no such handle")
		req.EqualError(err, "not found")
		req.Nil(fetched)
	})

	t.Run("lookup by username", func(t *testing.T) {
		fetched, err := s.LookupUserByUsername(ctx, user.Username)
		req.NoError(err)
		req.Equal(user.ID, fetched.ID)
	})

	t.Run("search by ID", func(t *testing.T) {
		set, f, err := s.SearchUsers(ctx, types.UserFilter{UserID: []uint64{user.ID}})
		req.NoError(err)
		req.Equal([]uint64{user.ID}, f.UserID)
		req.Len(set, 1)
		req.Equal(uint(1), f.Count)
		//req.Equal(set[0].ID, user.ID)
	})

	t.Run("search by email", func(t *testing.T) {
		set, f, err := s.SearchUsers(ctx, types.UserFilter{Email: user.Email})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal(uint(1), f.Count)
	})

	t.Run("search by username", func(t *testing.T) {
		set, f, err := s.SearchUsers(ctx, types.UserFilter{Username: user.Username})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal(uint(1), f.Count)
	})

	t.Run("search by query", func(t *testing.T) {
		set, f, err := s.SearchUsers(ctx, types.UserFilter{Query: user.Handle})
		req.NoError(err)
		req.Len(set, 1)
		req.Equal(uint(1), f.Count)
	})

	t.Run("search by username", func(t *testing.T) {
		set, f, err := s.SearchUsers(ctx, types.UserFilter{Username: "no such username"})
		req.NoError(err)
		req.Len(set, 0)
		req.Equal(uint(0), f.Count)
	})

	t.Run("search with masked details", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("search by role", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("count", func(t *testing.T) {
		var (
			f      = types.UserFilter{}
			c1, c2 uint
			err    error
			user   = &types.User{ID: id.Next(), CreatedAt: time.Now(), Email: fmt.Sprintf("user-crud+%s@crust.test", time.Now().String())}
		)

		c1, err = s.CountUsers(ctx, f)
		req.NoError(err)

		req.NoError(s.CreateUser(ctx, user))

		c2, err = s.CountUsers(ctx, f)
		req.NoError(err)
		req.Equal(c1+1, c2)

		req.NoError(s.RemoveUserByID(ctx, user.ID))

		c2, err = s.CountUsers(ctx, f)
		req.NoError(err)
		req.Equal(c1, c2)
	})

}
