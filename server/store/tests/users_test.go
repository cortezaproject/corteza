package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testUsers(t *testing.T, s store.Users) {
	var (
		dummyID uint64 = 0
		ctx            = context.Background()

		makeNew = func(nn ...string) *types.User {
			// minimum data set for new user
			name := strings.Join(nn, "")
			dummyID = dummyID + 1
			return &types.User{
				ID:        dummyID,
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

	t.Run("create and update deleted with existing email", func(t *testing.T) {
		req, user := truncAndCreate(t)
		deletedUser := makeNew("copy")
		deletedUser.DeletedAt = now()
		deletedUser.Email = user.Email
		req.NoError(store.CreateUser(ctx, s, deletedUser))
		req.NoError(store.UpdateUser(ctx, s, deletedUser))
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
	})

	t.Run("paging and sorting", func(t *testing.T) {
		require.NoError(t, s.TruncateUsers(ctx))
		require.NoError(t, s.CreateUser(ctx,
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
		))

		{
			tcc := []struct {
				sort string
				rval string
			}{
				{"email", "01..10"},
				{"email DESC", "10..01"},
				{"id, email", "01..10"},
				{"id DESC, email DESC", "10..01"},
				{"id DESC, email", "10..01"},
				{"id, email DESC", "01..10"},
			}

			for _, tc := range tcc {
				t.Run(tc.sort, func(t *testing.T) {
					var (
						req = require.New(t)

						f   = types.UserFilter{}
						set types.UserSet
						err error
					)

					f.Sort.Set(tc.sort)
					set, f, err = store.SearchUsers(ctx, s, f)
					req.NoError(err)
					req.Equal(tc.rval, stringifySetRange(set))
					req.Nil(f.PrevPage)
					req.Nil(f.NextPage)
				})
			}
		}

		{
			prevCur := 1
			nextCur := 2
			bothCur := prevCur | nextCur

			// tests if cursors are properly set/unset by inspecting req. bits
			testCursors := func(req *require.Assertions, b int, f types.UserFilter) {
				if b&prevCur == 0 {
					req.Nil(f.PrevPage)
				} else {
					req.NotNil(f.PrevPage)
				}

				if b&nextCur == 0 {
					req.Nil(f.NextPage)
				} else {
					req.NotNil(f.NextPage)
				}
			}

			tcc := []struct {
				// how data is sorted
				sort string

				// expected data
				rval []string

				// how cursors should be set when moving forward/backward
				curr []int
			}{
				{
					"id",
					[]string{"01..03", "04..06", "07..09", "10"},
					[]int{nextCur, bothCur, bothCur, prevCur},
				},
				{
					"id DESC",
					[]string{"10..08", "07..05", "04..02", "01"},
					[]int{nextCur, bothCur, bothCur, prevCur},
				},
				{
					"email",
					[]string{"01..03", "04..06", "07..09", "10"},
					[]int{nextCur, bothCur, bothCur, prevCur},
				},
				{
					"email DESC",
					[]string{"10..08", "07..05", "04..02", "01"},
					[]int{nextCur, bothCur, bothCur, prevCur},
				},
				{
					"id DESC, email",
					[]string{"10..08", "07..05", "04..02", "01"},
					[]int{nextCur, bothCur, bothCur, prevCur},
				},
			}

			for _, tc := range tcc {
				t.Run("crawling: "+tc.sort+"x", func(t *testing.T) {

					var (
						req = require.New(t)

						f   = types.UserFilter{}
						set types.UserSet
						err error
					)

					f.Sort.Set(tc.sort)
					f.Limit = 3 // 3, 3, 3, 1

					t.Log("going from page 1 to 4")
					for p := 0; p < 4; p++ {
						t.Logf("0123 next page cursor: %35s", f.PageCursor)
						set, f, err = store.SearchUsers(ctx, s, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifySetRange(set))

						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}

					t.Log("and back from 3 to 1")
					f.PageCursor = f.PrevPage
					for p := 2; p >= 0; p-- {
						f.Sort = nil
						t.Logf("_210 next page cursor: %35s", f.PageCursor)
						set, f, err = store.SearchUsers(ctx, s, f)
						t.Log(stringifySetRange(set))
						req.NoError(err)
						//req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","), "search should not return altered sort on filter")

						req.Equal(tc.rval[p], stringifySetRange(set))
						testCursors(req, tc.curr[p], f)

						// reverse to previous page
						f.PageCursor = f.PrevPage
					}

					t.Log("and again all the way to 4th page")
					f.PageCursor = f.NextPage
					for p := 1; p < 4; p++ {
						set, f, err = store.SearchUsers(ctx, s, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifySetRange(set))
						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}
				})
			}
		}

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

			f.Limit = uint(len(set))
			set, f, err := store.SearchUsers(ctx, s, f)
			req.Len(set, int(f.Limit))
			req.NoError(err)

			f.Limit = 1
			set, f, err = store.SearchUsers(ctx, s, f)
			req.NoError(err)

			// go to next page with different sorting
			f.PageCursor = f.NextPage
			f.Sort = filter.SortExprSet{&filter.SortExpr{Column: "username", Descending: false}}
			set, f, err = store.SearchUsers(ctx, s, f)
			req.EqualError(err, filter.ErrIncompatibleSort.Error())
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
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: "user-metrics-1@crust.test", UpdatedAt: &oct}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: "user-metrics-2@crust.test", UpdatedAt: &oct}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: "user-metrics-3@crust.test", SuspendedAt: &oct}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: "user-metrics-4@crust.test", SuspendedAt: &oct}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: "user-metrics-5@crust.test", SuspendedAt: &nov}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: "user-metrics-6@crust.test", DeletedAt: &nov}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: oct, Email: "user-metrics-7@crust.test", DeletedAt: &nov}))
		req.NoError(s.CreateUser(ctx, &types.User{ID: id.Next(), CreatedAt: nov, Email: "user-metrics-8@crust.test", DeletedAt: &nov}))

		m, err := store.UserMetrics(ctx, s)
		req.NoError(err)
		req.Equal(e, m)
	})
}
