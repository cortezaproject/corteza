package tests

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func testUsers(t *testing.T, s store.Users) {
	var (
		ctx = context.Background()

		makeNew = func(nn ...string) *types.User {
			// minimum data set for new user
			name := strings.Join(nn, "")
			return &types.User{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				Email:     "user-crud+" + name + "@crust.test",
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

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, types.UserSet) {
			req := require.New(t)
			req.NoError(s.TruncateUsers(ctx))

			set := make([]*types.User, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(string(rand.Bytes(10)))
			}

			req.NoError(s.CreateUser(ctx, set...))
			return req, set
		}

		stringifySetRange = func(set types.UserSet) string {
			if len(set) == 0 {
				return ""
			}

			var out = set[0].Handle[7:]

			if len(set) > 1 {
				out += ".." + set[len(set)-1].Handle[7:]
			}

			return out
		}

		// in case we need some quick old-school debugging
		//dbg := func(uu ...*types.User) {
		//	for i, u := range uu {
		//		fmt.Printf(" => #%2d %s\n", i+1, u.Handle)
		//	}
		//}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateUsers(ctx))
		req.NoError(store.CreateUser(ctx, s, makeNew()))
	})

	t.Run("create with duplicate email", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateUsers(ctx))

		req.NoError(s.CreateUser(ctx, makeNew("a")))
		req.EqualError(s.CreateUser(ctx, makeNew("a")), store.ErrNotUnique.Error())
	})

	t.Run("create with duplicate ID", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateUsers(ctx))
		user := makeNew("a")
		req.NoError(s.CreateUser(ctx, user))
		req.EqualError(s.CreateUser(ctx, user), store.ErrNotUnique.Error())
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, user := truncAndCreate(t)
		fetched, err := store.LookupUserByID(ctx, s, user.ID)
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

	t.Run("lookup by email", func(t *testing.T) {
		req, user := truncAndCreate(t)
		fetched, err := store.LookupUserByEmail(ctx, s, user.Email)
		req.NoError(err)
		req.Equal(user.Email, fetched.Email)
	})

	t.Run("lookup by handle", func(t *testing.T) {
		req, user := truncAndCreate(t)
		fetched, err := store.LookupUserByHandle(ctx, s, user.Handle)
		req.NoError(err)
		req.Equal(user.ID, fetched.ID)
	})

	t.Run("lookup by nonexisting handle", func(t *testing.T) {
		req, _ := truncAndCreate(t)
		fetched, err := store.LookupUserByHandle(ctx, s, "no such handle")
		req.EqualError(err, "not found")
		req.Nil(fetched)
	})

	t.Run("lookup by username", func(t *testing.T) {
		req, user := truncAndCreate(t)
		fetched, err := store.LookupUserByUsername(ctx, s, user.Username)
		req.NoError(err)
		req.Equal(user.ID, fetched.ID)
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by ID", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)
			set, f, err := store.SearchUsers(ctx, s, types.UserFilter{UserID: []uint64{prefill[0].ID}})
			req.NoError(err)
			req.Equal([]uint64{prefill[0].ID}, f.UserID)
			req.Len(set, 1)
			//req.Equal(set[0].ID, user.ID)
		})

		t.Run("by email", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)
			set, _, err := store.SearchUsers(ctx, s, types.UserFilter{Email: prefill[0].Email})
			req.NoError(err)
			req.Len(set, 1)
		})

		t.Run("by username", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)
			set, _, err := store.SearchUsers(ctx, s, types.UserFilter{Username: prefill[0].Username})
			req.NoError(err)
			req.Len(set, 1)
		})

		t.Run("by query", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)
			set, _, err := store.SearchUsers(ctx, s, types.UserFilter{Query: prefill[0].Handle})
			req.NoError(err)
			req.Len(set, 1)
		})

		t.Run("by username", func(t *testing.T) {
			req, _ := truncAndFill(t, 5)
			set, _, err := store.SearchUsers(ctx, s, types.UserFilter{Username: "no such username"})
			req.NoError(err)
			req.Len(set, 0)
		})

		t.Run("with check", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)
			set, _, err := store.SearchUsers(ctx, s, types.UserFilter{
				Check: func(user *types.User) (bool, error) {
					// simple check that matches with the first user from prefill
					return user.ID == prefill[0].ID, nil
				},
			})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(prefill[0].ID, set[0].ID)
		})

		t.Run("with paging", func(t *testing.T) {
			req := require.New(t)
			req.NoError(s.TruncateUsers(ctx))

			set := []*types.User{
				makeNew("01"),
				makeNew("02"),
				makeNew("03"),
				makeNew("04"),
				makeNew("05"),
				makeNew("06"),
				makeNew("07"),
				makeNew("08"),
				makeNew("09"),
				makeNew("10"),
			}

			req.NoError(s.CreateUser(ctx, set...))
			f := types.UserFilter{}
			f.Sort.Set("email")

			// Fetch first page
			f.Limit = 3
			set, f, err := store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Len(set, 3)
			req.NotNil(f.NextPage)
			req.Nil(f.PrevPage)
			req.Equal("01..03", stringifySetRange(set))

			// 2nd page
			f.Limit = 6
			f.PageCursor = f.NextPage
			set, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Len(set, 6)
			req.NotNil(f.NextPage)
			req.NotNil(f.PrevPage)
			req.Equal("04..09", stringifySetRange(set))

			// 3rd, last page (1 item left)
			f.Limit = 2
			f.PageCursor = f.NextPage
			set, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Len(set, 1)
			req.Nil(f.NextPage)
			req.NotNil(f.PrevPage)
			req.Equal("10", stringifySetRange(set))

			// now, in reverse, last 3 items
			f.Limit = 3
			f.PageCursor = f.PrevPage
			set, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Len(set, 3)
			req.NotNil(f.NextPage)
			req.NotNil(f.PrevPage)
			req.Equal("07..09", stringifySetRange(set))

			// still in reverse, next 5 items
			f.Limit = 5
			f.PageCursor = f.PrevPage
			set, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Len(set, 5)
			req.NotNil(f.NextPage)
			req.NotNil(f.PrevPage)
			req.Equal("02..06", stringifySetRange(set))

			// still in reverse, last 5 items (actually, we'll only get 1)
			f.Limit = 5
			f.PageCursor = f.PrevPage
			set, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Len(set, 1)
			req.Nil(f.PrevPage)
			req.NotNil(f.NextPage)
			req.Equal("01", stringifySetRange(set))
		})

		t.Run("with keyed paging and multi-key sorting", func(t *testing.T) {
			req := require.New(t)
			req.NoError(s.TruncateUsers(ctx))

			set := []*types.User{
				makeNew("01"),
				makeNew("02"),
				makeNew("03"),
				makeNew("04"),
				makeNew("05"),
			}

			req.NoError(s.CreateUser(ctx, set...))
			f := types.UserFilter{}
			f.Sort = filter.SortExprSet{&filter.SortExpr{Column: "email", Descending: true}, &filter.SortExpr{Column: "handle", Descending: true}}

			// Fetch first page
			f.Limit = 3
			set, f, err := store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Len(set, 3)
			req.NotNil(f.NextPage)
			req.Nil(f.PrevPage)
			req.Equal("handle_05", set[0].Handle)
			req.Equal("handle_03", set[2].Handle)

		})

		t.Run("with incompatible sort", func(t *testing.T) {
			req := require.New(t)
			req.NoError(s.TruncateUsers(ctx))

			set := []*types.User{
				makeNew("01"),
				makeNew("02"),
				makeNew("03"),
				makeNew("04"),
				makeNew("05"),
			}

			req.NoError(s.CreateUser(ctx, set...))
			f := types.UserFilter{}
			f.Sort = filter.SortExprSet{&filter.SortExpr{Column: "email", Descending: true}, &filter.SortExpr{Column: "handle", Descending: true}}

			f.Limit = 1
			set, f, err := store.SearchUsers(ctx, s, f)
			req.NoError(err)

			// go to next page with different sorting
			f.PageCursor = f.NextPage
			f.Sort = filter.SortExprSet{&filter.SortExpr{Column: "email", Descending: true}}
			set, f, err = store.SearchUsers(ctx, s, f)
			req.EqualError(err, "incompatible sort")
		})

		t.Run("with empty cursors", func(t *testing.T) {
			req := require.New(t)
			req.NoError(s.TruncateUsers(ctx))

			set := []*types.User{makeNew("01"), makeNew("02")}

			req.NoError(s.CreateUser(ctx, set...))

			f := types.UserFilter{}
			f.Limit = 2
			_, f, err := store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Nil(f.NextPage)
			req.Nil(f.PrevPage)

			f.Limit = 1
			_, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.NotNil(f.NextPage)
			req.Nil(f.PrevPage)

			f.PageCursor = f.NextPage
			_, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Nil(f.NextPage)
			req.NotNil(f.PrevPage)

			f = types.UserFilter{}
			f.Sort = filter.SortExprSet{&filter.SortExpr{Column: "id", Descending: true}}
			f.Limit = 2
			_, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Nil(f.NextPage)
			req.Nil(f.PrevPage)

			f.Limit = 1
			_, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.NotNil(f.NextPage)
			req.Nil(f.PrevPage)

			f.PageCursor = f.NextPage
			_, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Nil(f.NextPage)
			req.NotNil(f.PrevPage)

			// expecting empty next if check fn filters out resources
			f = types.UserFilter{Check: func(user *types.User) (bool, error) {
				return user.Handle != "handle_02", nil
			}}
			f.Sort = filter.SortExprSet{&filter.SortExpr{Column: "id", Descending: true}}
			f.Limit = 1
			_, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)
			req.Nil(f.NextPage)
			req.Nil(f.PrevPage)
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

		c1, err = store.CountUsers(ctx, s, f)
		req.NoError(err)

		req.NoError(s.CreateUser(ctx, user))

		c2, err = store.CountUsers(ctx, s, f)
		req.NoError(err)
		req.Equal(c1+1, c2)

		req.NoError(s.DeleteUserByID(ctx, user.ID))

		c2, err = store.CountUsers(ctx, s, f)
		req.NoError(err)
		req.Equal(c1, c2)
	})

	t.Run("metrics", func(t *testing.T) {
		var (
			req = require.New(t)

			oct, _ = time.Parse(time.RFC3339, "2020-10-02T10:01:10Z")
			nov, _ = time.Parse(time.RFC3339, "2020-11-02T20:02:10Z")

			octu = uint(oct.Truncate(time.Hour * 24).Unix())
			novu = uint(nov.Truncate(time.Hour * 24).Unix())

			e = &types.UserMetrics{
				Total:          8,
				Valid:          2,
				Deleted:        3,
				Suspended:      3,
				DailyCreated:   []uint{octu, 7, novu, 1},
				DailyUpdated:   []uint{octu, 2},
				DailySuspended: []uint{octu, 2, novu, 1},
				DailyDeleted:   []uint{novu, 3},
			}
		)

		req.NoError(s.TruncateUsers(ctx))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: fmt.Sprintf("user-crud+%s@crust.test", rand.Bytes(10)), UpdatedAt: &oct}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: fmt.Sprintf("user-crud+%s@crust.test", rand.Bytes(10)), UpdatedAt: &oct}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: fmt.Sprintf("user-crud+%s@crust.test", rand.Bytes(10)), SuspendedAt: &oct}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: fmt.Sprintf("user-crud+%s@crust.test", rand.Bytes(10)), SuspendedAt: &oct}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: fmt.Sprintf("user-crud+%s@crust.test", rand.Bytes(10)), SuspendedAt: &nov}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: fmt.Sprintf("user-crud+%s@crust.test", rand.Bytes(10)), DeletedAt: &nov}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: fmt.Sprintf("user-crud+%s@crust.test", rand.Bytes(10)), DeletedAt: &nov}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: nov, Email: fmt.Sprintf("user-crud+%s@crust.test", rand.Bytes(10)), DeletedAt: &nov}))

		m, err := store.UserMetrics(ctx, s)
		req.NoError(err)
		req.Equal(e, m)
	})
}
