package tests

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
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

		s = tmp.(usersStoreAdt)

		makeNew = func(nn ...string) *types.User {
			// minimum data set for new user
			name := strings.Join(nn, "")
			return &types.User{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Email:     "user-crud" + name + "@crust.test",
				Username:  "username_" + name,
				Handle:    "handle_" + name,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.User) {
			req := require.New(t)
			req.NoError(s.TruncateUsers(ctx))
			user := makeNew()
			req.NoError(s.CreateUser(ctx, user))
			return req, user
		}

		truncAddFill = func(t *testing.T, l int) (*require.Assertions, types.UserSet) {
			req := require.New(t)
			req.NoError(s.TruncateUsers(ctx))

			set := make([]*types.User, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(string(rand.Bytes(10)))
			}

			req.NoError(s.CreateUser(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateUser(ctx, makeNew()))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, user := truncAndCreate(t)

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
		req, user := truncAndCreate(t)
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
		req, user := truncAndCreate(t)

		fetched, err := s.LookupUserByEmail(ctx, user.Email)
		req.NoError(err)
		req.Equal(user.Email, fetched.Email)
	})

	t.Run("lookup by handle", func(t *testing.T) {
		req, user := truncAndCreate(t)

		fetched, err := s.LookupUserByHandle(ctx, user.Handle)
		req.NoError(err)
		req.Equal(user.ID, fetched.ID)
	})

	t.Run("lookup by nonexisting handle", func(t *testing.T) {
		req, _ := truncAndCreate(t)

		fetched, err := s.LookupUserByHandle(ctx, "no such handle")
		req.EqualError(err, "not found")
		req.Nil(fetched)
	})

	t.Run("lookup by username", func(t *testing.T) {
		req, user := truncAndCreate(t)

		fetched, err := s.LookupUserByUsername(ctx, user.Username)
		req.NoError(err)
		req.Equal(user.ID, fetched.ID)
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by ID", func(t *testing.T) {
			req, prefill := truncAddFill(t, 5)

			set, f, err := s.SearchUsers(ctx, types.UserFilter{UserID: []uint64{prefill[0].ID}})
			req.NoError(err)
			req.Equal([]uint64{prefill[0].ID}, f.UserID)
			req.Len(set, 1)
			req.Equal(uint(1), f.Count)
			//req.Equal(set[0].ID, user.ID)
		})

		t.Run("by email", func(t *testing.T) {
			req, prefill := truncAddFill(t, 5)

			set, f, err := s.SearchUsers(ctx, types.UserFilter{Email: prefill[0].Email})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(uint(1), f.Count)
		})

		t.Run("by username", func(t *testing.T) {
			req, prefill := truncAddFill(t, 5)
			set, f, err := s.SearchUsers(ctx, types.UserFilter{Username: prefill[0].Username})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(uint(1), f.Count)
		})

		t.Run("by query", func(t *testing.T) {
			req, prefill := truncAddFill(t, 5)
			set, f, err := s.SearchUsers(ctx, types.UserFilter{Query: prefill[0].Handle})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(uint(1), f.Count)
		})

		t.Run("by username", func(t *testing.T) {
			req, _ := truncAddFill(t, 5)
			set, f, err := s.SearchUsers(ctx, types.UserFilter{Username: "no such username"})
			req.NoError(err)
			req.Len(set, 0)
			req.Equal(uint(0), f.Count)
		})

		t.Run("with check", func(t *testing.T) {
			req, prefill := truncAddFill(t, 5)
			set, _, err := s.SearchUsers(ctx, types.UserFilter{
				Check: func(user *types.User) (bool, error) {
					// simple check that matches with the first user from prefill
					return user.ID == prefill[0].ID, nil
				},
			})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(prefill[0].ID, set[0].ID)
		})

		t.Run("with check and paging", func(t *testing.T) {
			req, prefill := truncAddFill(t, 5)
			set, _, err := s.SearchUsers(ctx, types.UserFilter{
				// This will cause paging to run multiple queries
				// until it collects all data
				PageFilter: rh.PageFilter{Limit: 2},
				Check: func(user *types.User) (bool, error) {
					// simple check that matches with the 4th user from prefill
					return user.ID == prefill[4].ID, nil
				},
			})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(prefill[4].ID, set[0].ID)
		})

		t.Run("with masked details", func(t *testing.T) {
			t.Skip("not implemented")
		})

		t.Run("by role", func(t *testing.T) {
			t.Skip("not implemented")
		})

		t.Run("search", func(t *testing.T) {
			t.Skip("not implemented")
		})
	})

	t.Run("count", func(t *testing.T) {
		var (
			req = require.New(t)

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
