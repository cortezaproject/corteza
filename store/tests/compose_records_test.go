package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

func testComposeRecords(t *testing.T, s store.ComposeRecords) {
	var (
		ctx = context.Background()

		mod = &types.Module{
			ID:          id.Next(),
			NamespaceID: id.Next(),
			Handle:      "",
			Name:        "testComposeRecords",
			CreatedAt:   time.Now(),
			Fields: types.ModuleFieldSet{
				&types.ModuleField{Kind: "String", Name: "str1"},
				&types.ModuleField{Kind: "String", Name: "str2"},
				&types.ModuleField{Kind: "String", Name: "str3"},
				&types.ModuleField{Kind: "Number", Name: "num1"},
				&types.ModuleField{Kind: "Number", Name: "num2"},
				&types.ModuleField{Kind: "Number", Name: "num3"},
				&types.ModuleField{Kind: "DateTime", Name: "dt1"},
				&types.ModuleField{Kind: "String", Name: "strMulti", Multi: true},

				// These are used bellow to test (almost) all permutations
				&types.ModuleField{Kind: "Bool", Name: "bool1"},
				&types.ModuleField{Kind: "Bool", Name: "bool2"},
				&types.ModuleField{Kind: "Bool", Name: "bool3"},

				&types.ModuleField{Kind: "DateTime", Name: "datetime1"},
				&types.ModuleField{Kind: "DateTime", Name: "datetime2"},
				&types.ModuleField{Kind: "DateTime", Name: "datetime3"},

				&types.ModuleField{Kind: "Email", Name: "email1"},
				&types.ModuleField{Kind: "Email", Name: "email2"},
				&types.ModuleField{Kind: "Email", Name: "email3"},

				// &types.ModuleField{Kind: "File", Name: "file1"},
				// &types.ModuleField{Kind: "File", Name: "file2"},
				// &types.ModuleField{Kind: "File", Name: "file3"},

				&types.ModuleField{Kind: "Select", Name: "select1"},
				&types.ModuleField{Kind: "Select", Name: "select2"},
				&types.ModuleField{Kind: "Select", Name: "select3"},

				&types.ModuleField{Kind: "Number", Name: "number1"},
				&types.ModuleField{Kind: "Number", Name: "number2"},
				&types.ModuleField{Kind: "Number", Name: "number3"},

				// &types.ModuleField{Kind: "Record", Name: "record1"},
				// &types.ModuleField{Kind: "Record", Name: "record2"},
				// &types.ModuleField{Kind: "Record", Name: "record3"},

				&types.ModuleField{Kind: "String", Name: "string1"},
				&types.ModuleField{Kind: "String", Name: "string2"},
				&types.ModuleField{Kind: "String", Name: "string3"},

				&types.ModuleField{Kind: "Url", Name: "url1"},
				&types.ModuleField{Kind: "Url", Name: "url2"},
				&types.ModuleField{Kind: "Url", Name: "url3"},

				// &types.ModuleField{Kind: "User", Name: "user1"},
				// &types.ModuleField{Kind: "User", Name: "user2"},
				// &types.ModuleField{Kind: "User", Name: "user3"},
			},
		}

		makeNew = func(vv ...*types.RecordValue) *types.Record {
			// minimum data set for new composeRecord
			var recordID = id.Next()

			for _, v := range vv {
				v.RecordID = recordID
			}

			return &types.Record{
				ID:          recordID,
				NamespaceID: mod.NamespaceID,
				ModuleID:    mod.ID,
				CreatedAt:   time.Now(),
				Values:      vv,
			}
		}

		makeNewUpd = func(t *testing.T, dt string, vv ...*types.RecordValue) *types.Record {
			r := makeNew(vv...)

			if dt != "" {
				n, err := time.Parse(time.RFC3339, dt)
				if err != nil {
					t.Error(err)
				}
				r.UpdatedAt = &n
			}
			return r
		}

		truncAndCreate = func(t *testing.T, rr ...*types.Record) (*require.Assertions, types.RecordSet) {
			req := require.New(t)
			req.NoError(s.TruncateComposeRecords(ctx, mod))

			if len(rr) == 0 {
				rr = []*types.Record{makeNew()}
			}

			for _, rec := range rr {
				req.NoError(s.CreateComposeRecord(ctx, mod, rec))
			}

			return req, rr
		}

		stringifyValues = func(set types.RecordSet, fields ...string) string {
			var out string
			for r := range set {
				if r > 0 {
					out += ";"
				}

				for f := range fields {
					if f > 0 {
						out += ","
					}

					v := set[r].Values.Get(fields[f], 0)
					if v != nil {
						out += v.Value
					} else {
						out += "<NIL>"
					}

				}

			}

			return out
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		composeRecord := makeNew()
		req.NoError(s.CreateComposeRecord(ctx, mod, composeRecord))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, rr := truncAndCreate(t, makeNew(
			&types.RecordValue{Name: "str1", Value: "v1", Ref: 1},
			&types.RecordValue{Name: "str2", Value: "v2", Ref: 2},
			&types.RecordValue{Name: "str3", Value: "v3", Ref: 3},
		))
		rec := rr[0]

		fetched, err := s.LookupComposeRecordByID(ctx, mod, rec.ID)
		req.NoError(err)
		req.Equal(rec.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
		req.Len(fetched.Values, len(rec.Values))
		req.Equal("str2", fetched.Values[1].Name)
		req.Equal("v2", fetched.Values[1].Value)
		req.Equal(uint64(2), fetched.Values[1].Ref)
	})

	t.Run("update", func(t *testing.T) {
		req, rr := truncAndCreate(t)
		rec := rr[0]

		rec = &types.Record{
			ID:          rec.ID,
			CreatedAt:   rec.CreatedAt,
			ModuleID:    mod.ID,
			NamespaceID: mod.NamespaceID,
			OwnedBy:     id.Next(),
		}

		req.NoError(s.UpdateComposeRecord(ctx, mod, rec))

		updated, err := s.LookupComposeRecordByID(ctx, mod, rec.ID)
		req.NoError(err)
		req.Equal(rec.OwnedBy, updated.OwnedBy)
	})

	t.Run("update values", func(t *testing.T) {
		req, rr := truncAndCreate(t, makeNew(
			&types.RecordValue{Name: "str1", Value: "v1", Ref: 1},
			&types.RecordValue{Name: "str2", Value: "v2", Ref: 2},
		))
		rec := rr[0]

		rec = &types.Record{
			ID:          rec.ID,
			CreatedAt:   rec.CreatedAt,
			OwnedBy:     id.Next(),
			Values:      rec.Values,
			ModuleID:    mod.ID,
			NamespaceID: mod.NamespaceID,
		}

		rec.Values[0].Value = "vv10"
		rec.Values[1].Value = "vv20"
		rec.Values = append(rec.Values, &types.RecordValue{Name: "str3", Value: "vv30", Ref: 3})
		rec.Values.SetRecordID(rec.ID)

		req.NoError(s.UpdateComposeRecord(ctx, mod, rec))

		updated, err := s.LookupComposeRecordByID(ctx, mod, rec.ID)
		req.NoError(err)
		req.Equal(rec.OwnedBy, updated.OwnedBy)
		req.Len(updated.Values, len(rec.Values))
		req.Equal("str2", updated.Values[1].Name)
		req.Equal("vv20", updated.Values[1].Value)
	})

	t.Run("soft delete values", func(t *testing.T) {
		req, rr := truncAndCreate(t, makeNew(
			&types.RecordValue{Name: "str1", Value: "v1", Ref: 1},
			&types.RecordValue{Name: "str2", Value: "v2", Ref: 2},
		))
		rec := rr[0]
		rec.DeletedAt = &rec.CreatedAt

		req.NoError(s.UpdateComposeRecord(ctx, mod, rec))

		updated, err := s.LookupComposeRecordByID(ctx, mod, rec.ID)

		req.NoError(err)
		req.NotNil(rec)
		req.NotNil(rec.DeletedAt)
		req.Len(updated.Values, len(rec.Values))
		req.NotNil(updated.Values[0].DeletedAt)
		req.NotNil(updated.Values[1].DeletedAt)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Record", func(t *testing.T) {
			req, rr := truncAndCreate(t)
			rec := rr[0]

			req.NoError(s.DeleteComposeRecord(ctx, mod, rec))
			_, err := s.LookupComposeRecordByID(ctx, mod, rec.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, rr := truncAndCreate(t)
			rec := rr[0]

			req.NoError(s.DeleteComposeRecordByID(ctx, mod, rec.ID))
			_, err := s.LookupComposeRecordByID(ctx, mod, rec.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by record attributes", func(t *testing.T) {
			prefill := []*types.Record{
				makeNew(),
				makeNew(),
				makeNew(),
				makeNew(),
				makeNew(),
			}

			count := len(prefill)

			prefill[4].DeletedAt = &prefill[4].CreatedAt
			valid := count - 1

			req, _ := truncAndCreate(t, prefill...)

			// search for all valid
			set, _, err := s.SearchComposeRecords(ctx, mod, types.RecordFilter{})
			req.NoError(err)
			req.Len(set, valid) // we've deleted one

			// search for ALL
			set, _, err = s.SearchComposeRecords(ctx, mod, types.RecordFilter{Deleted: filter.StateInclusive})
			req.NoError(err)
			req.Len(set, count) // we've deleted one

			// search for deleted only
			set, _, err = s.SearchComposeRecords(ctx, mod, types.RecordFilter{Deleted: filter.StateExclusive})
			req.NoError(err)
			req.Len(set, 1) // we've deleted one
		})

		t.Run("by values", func(t *testing.T) {
			var (
				err error
				set types.RecordSet

				req, _ = truncAndCreate(t,
					makeNew(&types.RecordValue{Name: "str1", Value: "v1"}, &types.RecordValue{Name: "str2", Value: "same"}, &types.RecordValue{Name: "str3", Value: "three"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "v2"}, &types.RecordValue{Name: "str2", Value: "same"}, &types.RecordValue{Name: "str3", Value: "three"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "v3"}, &types.RecordValue{Name: "str2", Value: "same"}, &types.RecordValue{Name: "str3", Value: "three"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "v4"}, &types.RecordValue{Name: "str2", Value: "same"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "v5"}, &types.RecordValue{Name: "str2", Value: "same"}),

					// Add one additional record with deleted values
					makeNew(&types.RecordValue{Name: "str1", Value: "v6", DeletedAt: now()}, &types.RecordValue{Name: "str2", Value: "deleted", DeletedAt: now()}),
				)

				f = types.RecordFilter{
					ModuleID:    mod.ID,
					NamespaceID: mod.NamespaceID,
				}
			)

			f.Query = `str1 = 'v1'`
			set, _, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			req.Len(set, 1)

			f.Query = `str2 = 'same'`
			set, _, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			req.Len(set, 5)

			f.Query = `str2 = 'different'`
			set, _, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			req.Len(set, 0)

			f.Query = `str3 = 'three' AND str1 = 'v1'`
			set, _, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			req.Len(set, 1)
		})
	})

	t.Run("paging and sorting", func(t *testing.T) {
		type (
			tc struct {
				// how data is sorted
				sort string

				// expected data
				rval []string

				// how cursors should be set when moving forward/backward
				curr []int
			}
		)
		var (
			_, _ = truncAndCreate(t,
				makeNew(&types.RecordValue{Name: "str1", Value: "v1"}, &types.RecordValue{Name: "str3", Value: "a"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v2"}, &types.RecordValue{Name: "str3", Value: "b"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v3"}, &types.RecordValue{Name: "str3", Value: "b"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v4"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v5"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v6"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v7"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v8"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v9"}, &types.RecordValue{Name: "str3", Value: "c"}),
			)
		)

		prevCur := 1
		nextCur := 2
		bothCur := prevCur | nextCur

		// tests if cursors are properly set/unset by inspecting req. bits
		testCursors := func(req *require.Assertions, b int, f types.RecordFilter) {
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

		tcc := []tc{
			{
				"id",
				[]string{"v1,a;v2,b;v3,b", "v4,<NIL>;v5,<NIL>;v6,<NIL>", "v7,<NIL>;v8,<NIL>;v9,c"},
				[]int{nextCur, bothCur, prevCur},
			},
			{
				"id DESC",
				[]string{"v9,c;v8,<NIL>;v7,<NIL>", "v6,<NIL>;v5,<NIL>;v4,<NIL>", "v3,b;v2,b;v1,a"},
				[]int{nextCur, bothCur, prevCur},
			},
			{
				"str1",
				[]string{"v1,a;v2,b;v3,b", "v4,<NIL>;v5,<NIL>;v6,<NIL>", "v7,<NIL>;v8,<NIL>;v9,c"},
				[]int{nextCur, bothCur, prevCur},
			},
			{
				"str1 DESC",
				[]string{"v9,c;v8,<NIL>;v7,<NIL>", "v6,<NIL>;v5,<NIL>;v4,<NIL>", "v3,b;v2,b;v1,a"},
				[]int{nextCur, bothCur, prevCur},
			},
			{
				"str3",
				[]string{"v4,<NIL>;v5,<NIL>;v6,<NIL>", "v7,<NIL>;v8,<NIL>;v1,a", "v2,b;v3,b;v9,c"},
				[]int{nextCur, bothCur, prevCur},
			},
			{
				"str3 DESC",
				[]string{"v9,c;v3,b;v2,b", "v1,a;v8,<NIL>;v7,<NIL>", "v6,<NIL>;v5,<NIL>;v4,<NIL>"},
				[]int{nextCur, bothCur, prevCur},
			},
		}

		for _, tc := range tcc {
			t.Run("crawling: "+tc.sort, func(t *testing.T) {
				var (
					req = require.New(t)

					f   = types.RecordFilter{}
					set types.RecordSet
					err error
				)

				f.Sort.Set(tc.sort)
				f.Limit = 3

				for p := 0; p < 3; p++ {
					set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
					req.NoError(err)
					req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))
					req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))

					testCursors(req, tc.curr[p], f)

					// advance to next page
					f.PageCursor = f.NextPage
				}

				f.PageCursor = f.PrevPage
				for p := 1; p >= 0; p-- {
					f.Sort = nil
					set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
					req.NoError(err)
					req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

					req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))
					testCursors(req, tc.curr[p], f)

					// reverse to previous page
					f.PageCursor = f.PrevPage
				}

				f.PageCursor = f.NextPage

				for p := 1; p < 3; p++ {
					set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
					req.NoError(err)
					req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

					req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))
					testCursors(req, tc.curr[p], f)

					// advance to next page
					f.PageCursor = f.NextPage
				}
			})
		}

		t.Run("advanced sorting", func(t *testing.T) {
			var (
				_, _ = truncAndCreate(t,
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0010"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0009"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0008"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0007"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0006"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0005"}),
				)
			)

			tcc := []tc{
				{
					"str1",
					[]string{"0001,0010;0001,0009;0001,0008", "0001,0007;0001,0006;0001,0005"},
					[]int{nextCur, prevCur},
				},
				{
					"str1, str3 DESC",
					[]string{"0001,0010;0001,0009;0001,0008", "0001,0007;0001,0006;0001,0005"},
					[]int{nextCur, prevCur},
				},
				{
					"str1 DESC, str3",
					[]string{"0001,0005;0001,0006;0001,0007", "0001,0008;0001,0009;0001,0010"},
					[]int{nextCur, prevCur},
				},
			}

			for _, tc := range tcc {
				t.Run("crawling: "+tc.sort, func(t *testing.T) {
					var (
						req = require.New(t)

						f   = types.RecordFilter{}
						set types.RecordSet
						err error
					)

					f.Sort.Set(tc.sort)
					f.Limit = 3

					for p := 0; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))
						req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))

						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}

					f.PageCursor = f.PrevPage
					for p := 0; p >= 0; p-- {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))
						testCursors(req, tc.curr[p], f)

						// reverse to previous page
						f.PageCursor = f.PrevPage
					}

					f.PageCursor = f.NextPage

					for p := 1; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))
						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}
				})
			}
		})

		t.Run("advanced sorting; NULL; #1", func(t *testing.T) {
			var (
				_, _ = truncAndCreate(t,
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0010"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0009"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0008"}),
					makeNew( /* padding                                   */ &types.RecordValue{Name: "str3", Value: "0007"}),
					makeNew( /* padding                                   */ &types.RecordValue{Name: "str3", Value: "0006"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0005"}),
				)
			)

			tcc := []tc{
				{
					"str1",
					[]string{"<NIL>,0007;<NIL>,0006;0001,0010", "0001,0009;0001,0008;0001,0005"},
					[]int{nextCur, prevCur},
				},
				{
					"str1 DESC",
					[]string{"0001,0005;0001,0008;0001,0009", "0001,0010;<NIL>,0006;<NIL>,0007"},
					[]int{nextCur, prevCur},
				},
				{
					"str1, str3 DESC",
					[]string{"<NIL>,0007;<NIL>,0006;0001,0010", "0001,0009;0001,0008;0001,0005"},
					[]int{nextCur, prevCur},
				},
				{
					"str1 DESC, str3 DESC",
					[]string{"0001,0010;0001,0009;0001,0008", "0001,0005;<NIL>,0007;<NIL>,0006"},
					[]int{nextCur, prevCur},
				},
			}

			for _, tc := range tcc {
				t.Run("crawling: "+tc.sort, func(t *testing.T) {
					var (
						req = require.New(t)

						f   = types.RecordFilter{}
						set types.RecordSet
						err error
					)

					f.Sort.Set(tc.sort)
					f.Limit = 3

					for p := 0; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))
						req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))

						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}

					f.PageCursor = f.PrevPage
					for p := 0; p >= 0; p-- {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))
						testCursors(req, tc.curr[p], f)

						// reverse to previous page
						f.PageCursor = f.PrevPage
					}

					f.PageCursor = f.NextPage

					for p := 1; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))
						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}
				})
			}
		})

		t.Run("advanced sorting; NULL; #2", func(t *testing.T) {
			var (
				_, _ = truncAndCreate(t,
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"}, &types.RecordValue{Name: "str3", Value: "0010"}),
					makeNew( /* padding                                   */ &types.RecordValue{Name: "str3", Value: "0009"}),
					makeNew(&types.RecordValue{Name: "str1", Value: "0001"} /* padding                                    */),
					makeNew( /* padding                                  */ /* padding                                   */ ),
				)
			)

			tcc := []tc{
				{
					"str1, str3",
					[]string{"<NIL>,<NIL>;<NIL>,0009", "0001,<NIL>;0001,0010"},
					[]int{nextCur, prevCur},
				},
				{
					"str1, str3 DESC",
					[]string{"<NIL>,0009;<NIL>,<NIL>", "0001,0010;0001,<NIL>"},
					[]int{nextCur, prevCur},
				},
				{
					"str1 DESC, str3",
					[]string{"0001,<NIL>;0001,0010", "<NIL>,<NIL>;<NIL>,0009"},
					[]int{nextCur, prevCur},
				},
				{
					"str1 DESC, str3 DESC",
					[]string{"0001,0010;0001,<NIL>", "<NIL>,0009;<NIL>,<NIL>"},
					[]int{nextCur, prevCur},
				},
			}

			for _, tc := range tcc {
				t.Run("crawling: "+tc.sort, func(t *testing.T) {
					var (
						req = require.New(t)

						f   = types.RecordFilter{}
						set types.RecordSet
						err error
					)

					f.Sort.Set(tc.sort)
					f.Limit = 2

					for p := 0; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))
						req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))

						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}

					f.PageCursor = f.PrevPage
					for p := 0; p >= 0; p-- {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))
						testCursors(req, tc.curr[p], f)

						// reverse to previous page
						f.PageCursor = f.PrevPage
					}

					f.PageCursor = f.NextPage

					for p := 1; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifyValues(set, "str1", "str3"))
						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}
				})
			}
		})

		t.Run("advanced sorting; NULL; #4", func(t *testing.T) {
			// Build all record value permutations for the given field based on a pattern in a 3x8 matrix
			//
			// The matrix defines all NULL / NOT NULL permutations where NOT NULL values define a mini permutation
			// of increasing values -- a good enough permutation to prove correctness.
			//
			// Boolean and pure Timestamp values are tested elsewhere
			permute := func(ptrn string, f1, f2, f3 string) {
				spf := func(ptr, val string) string {
					return fmt.Sprintf(ptrn, val)
				}

				_, _ = truncAndCreate(t,
					makeNew(&types.RecordValue{Name: f1, Value: spf(ptrn, "1")}, &types.RecordValue{Name: f2, Value: spf(ptrn, "1")}, &types.RecordValue{Name: f3, Value: spf(ptrn, "1")}),
					makeNew( /* padding                                       */ &types.RecordValue{Name: f2, Value: spf(ptrn, "1")}, &types.RecordValue{Name: f3, Value: spf(ptrn, "2")}),
					makeNew(&types.RecordValue{Name: f1, Value: spf(ptrn, "1")} /* padding                                        */, &types.RecordValue{Name: f3, Value: spf(ptrn, "3")}),
					makeNew( /* padding                                       */ /* padding                                        */ &types.RecordValue{Name: f3, Value: spf(ptrn, "4")}),
					makeNew(&types.RecordValue{Name: f1, Value: spf(ptrn, "2")}, &types.RecordValue{Name: f2, Value: spf(ptrn, "3")} /* padding                                        */),
					makeNew( /* padding                                       */ &types.RecordValue{Name: f2, Value: spf(ptrn, "3")} /* padding                                        */),
					makeNew(&types.RecordValue{Name: f1, Value: spf(ptrn, "2")} /* padding                                        */ /* padding                                        */),
					makeNew( /* padding                                       */ /* padding                                       */ /* padding                                       */ ),
				)
			}

			kinds := []struct {
				kind string
				ptrn string
			}{
				{kind: "DateTime", ptrn: "2021-02-0%sT01:00:00.000Z"},
				{kind: "Email", ptrn: "test+e%s@test.tld"},
				{kind: "Select", ptrn: "opt%s"},
				{kind: "Number", ptrn: "%s"},
				{kind: "String", ptrn: "string%s"},
				{kind: "Url", ptrn: "https://www.testko-%s.tld"},
				// { kind: "Bool", ptrn: ""},
				// { kind: "File", ptrn: ""},
				// { kind: "Record", ptrn: ""},
				// { kind: "User", ptrn: ""},
			}

			// These test cases check all 3 field sort permutations.
			// Each field kind must return items in the same order.
			//
			// 0 == <NIL>
			tcc := []struct {
				sort string
				curr []int
				exp  [][]string
			}{
				{
					sort: "%s, %s, %s",
					curr: []int{nextCur, bothCur, bothCur, prevCur},
					exp: [][]string{
						{"0", "0", "0"},
						{"0", "0", "4"},
						{"0", "1", "2"},
						{"0", "3", "0"},
						{"1", "0", "3"},
						{"1", "1", "1"},
						{"2", "0", "0"},
						{"2", "3", "0"},
					},
				},
				{
					sort: "%s DESC, %s, %s",
					curr: []int{nextCur, bothCur, bothCur, prevCur},
					exp: [][]string{
						{"2", "0", "0"},
						{"2", "3", "0"},
						{"1", "0", "3"},
						{"1", "1", "1"},
						{"0", "0", "0"},
						{"0", "0", "4"},
						{"0", "1", "2"},
						{"0", "3", "0"},
					},
				},
				{
					sort: "%s, %s DESC, %s",
					curr: []int{nextCur, bothCur, bothCur, prevCur},
					exp: [][]string{
						{"0", "3", "0"},
						{"0", "1", "2"},
						{"0", "0", "0"},
						{"0", "0", "4"},
						{"1", "1", "1"},
						{"1", "0", "3"},
						{"2", "3", "0"},
						{"2", "0", "0"},
					},
				},
				{
					sort: "%s DESC, %s DESC, %s",
					curr: []int{nextCur, bothCur, bothCur, prevCur},
					exp: [][]string{
						{"2", "3", "0"},
						{"2", "0", "0"},
						{"1", "1", "1"},
						{"1", "0", "3"},
						{"0", "3", "0"},
						{"0", "1", "2"},
						{"0", "0", "0"},
						{"0", "0", "4"},
					},
				},
				{
					sort: "%s, %s, %s DESC",
					curr: []int{nextCur, bothCur, bothCur, prevCur},
					exp: [][]string{
						{"0", "0", "4"},
						{"0", "0", "0"},
						{"0", "1", "2"},
						{"0", "3", "0"},
						{"1", "0", "3"},
						{"1", "1", "1"},
						{"2", "0", "0"},
						{"2", "3", "0"},
					},
				},
				{
					sort: "%s DESC, %s, %s DESC",
					curr: []int{nextCur, bothCur, bothCur, prevCur},
					exp: [][]string{
						{"2", "0", "0"},
						{"2", "3", "0"},
						{"1", "0", "3"},
						{"1", "1", "1"},
						{"0", "0", "4"},
						{"0", "0", "0"},
						{"0", "1", "2"},
						{"0", "3", "0"},
					},
				},
				{
					sort: "%s, %s DESC, %s DESC",
					curr: []int{nextCur, bothCur, bothCur, prevCur},
					exp: [][]string{
						{"0", "3", "0"},
						{"0", "1", "2"},
						{"0", "0", "4"},
						{"0", "0", "0"},
						{"1", "1", "1"},
						{"1", "0", "3"},
						{"2", "3", "0"},
						{"2", "0", "0"},
					},
				},
				{
					sort: "%s DESC, %s DESC, %s DESC",
					curr: []int{nextCur, bothCur, bothCur, prevCur},
					exp: [][]string{
						{"2", "3", "0"},
						{"2", "0", "0"},
						{"1", "1", "1"},
						{"1", "0", "3"},
						{"0", "3", "0"},
						{"0", "1", "2"},
						{"0", "0", "4"},
						{"0", "0", "0"},
					},
				},
			}

			// Helper to comvert expected items from above-defined cases into a string.
			expToString := func(ptrn string, exp [][]string) string {
				rr := make([]string, 0, len(exp))

				for _, page := range exp {
					rec := make([]string, 0)
					for _, val := range page {
						if val == "0" {
							rec = append(rec, "<NIL>")
						} else {
							rec = append(rec, fmt.Sprintf(ptrn, val))
						}
					}
					rr = append(rr, strings.Join(rec, ","))
				}
				return strings.Join(rr, ";")
			}

			for _, k := range kinds {
				for _, tc := range tcc {

					f1 := strings.ToLower(k.kind) + "1"
					f2 := strings.ToLower(k.kind) + "2"
					f3 := strings.ToLower(k.kind) + "3"

					tc.sort = fmt.Sprintf(tc.sort, f1, f2, f3)

					t.Run("crawling: "+k.kind+"; "+tc.sort, func(t *testing.T) {
						var (
							req = require.New(t)

							f   = types.RecordFilter{}
							set types.RecordSet
							err error
						)

						f.Sort.Set(tc.sort)
						f.Limit = 2

						permute(k.ptrn, f1, f2, f3)

						for p := 0; p < 4; p++ {
							set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
							req.NoError(err)
							req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

							ff := p * 2
							req.Equal(expToString(k.ptrn, tc.exp[ff:ff+int(f.Limit)]), stringifyValues(set, strings.ToLower(k.kind)+"1", strings.ToLower(k.kind)+"2", strings.ToLower(k.kind)+"3"))

							testCursors(req, tc.curr[p], f)

							// advance to next page
							f.PageCursor = f.NextPage
						}

						f.PageCursor = f.PrevPage
						for p := 2; p >= 0; p-- {
							set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
							req.NoError(err)
							req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

							ff := p * 2
							req.Equal(expToString(k.ptrn, tc.exp[ff:ff+int(f.Limit)]), stringifyValues(set, strings.ToLower(k.kind)+"1", strings.ToLower(k.kind)+"2", strings.ToLower(k.kind)+"3"))
							testCursors(req, tc.curr[p], f)

							// reverse to previous page
							f.PageCursor = f.PrevPage
						}

						f.PageCursor = f.NextPage

						for p := 1; p < 4; p++ {
							set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
							req.NoError(err)
							req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

							ff := p * 2
							req.Equal(expToString(k.ptrn, tc.exp[ff:ff+int(f.Limit)]), stringifyValues(set, strings.ToLower(k.kind)+"1", strings.ToLower(k.kind)+"2", strings.ToLower(k.kind)+"3"))
							testCursors(req, tc.curr[p], f)

							// advance to next page
							f.PageCursor = f.NextPage
						}
					})
				}
			}
		})

		t.Run("advanced sorting; NULL; Booleans", func(t *testing.T) {
			var (
				_, _ = truncAndCreate(t,
					makeNew(&types.RecordValue{Name: "bool1", Value: ""}, &types.RecordValue{Name: "bool3", Value: ""}),
					makeNew( /* padding                                */ &types.RecordValue{Name: "bool3", Value: "1"}),
					makeNew(&types.RecordValue{Name: "bool1", Value: "1"} /* padding                                 */),
					makeNew( /* padding                                */ /* padding                                */ ),
				)
			)

			tcc := []tc{
				{
					"bool1, bool3",
					[]string{"<NIL>,<NIL>;<NIL>,1", ",;1,<NIL>"},
					[]int{nextCur, prevCur},
				},
				{
					"bool1, bool3 DESC",
					[]string{"<NIL>,1;<NIL>,<NIL>", ",;1,<NIL>"},
					[]int{nextCur, prevCur},
				},
				{
					"bool1 DESC, bool3",
					[]string{"1,<NIL>;,", "<NIL>,<NIL>;<NIL>,1"},
					[]int{nextCur, prevCur},
				},
				{
					"bool1 DESC, bool3 DESC",
					[]string{"1,<NIL>;,", "<NIL>,1;<NIL>,<NIL>"},
					[]int{nextCur, prevCur},
				},
			}

			for _, tc := range tcc {
				t.Run("crawling: "+tc.sort, func(t *testing.T) {
					var (
						req = require.New(t)

						f   = types.RecordFilter{}
						set types.RecordSet
						err error
					)

					f.Sort.Set(tc.sort)
					f.Limit = 2

					for p := 0; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))
						req.Equal(tc.rval[p], stringifyValues(set, "bool1", "bool3"))

						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}

					f.PageCursor = f.PrevPage
					for p := 0; p >= 0; p-- {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifyValues(set, "bool1", "bool3"))
						testCursors(req, tc.curr[p], f)

						// reverse to previous page
						f.PageCursor = f.PrevPage
					}

					f.PageCursor = f.NextPage

					for p := 1; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifyValues(set, "bool1", "bool3"))
						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}
				})
			}
		})

		t.Run("advanced sorting; disable multi-value sorting", func(t *testing.T) {
			var (
				_, _ = truncAndCreate(t,
					makeNew(&types.RecordValue{Name: "strMulti", Place: 0, Value: "a"}, &types.RecordValue{Name: "strMulti", Place: 1, Value: "b"}, &types.RecordValue{Name: "strMulti", Place: 2, Value: "c"}),
				)
			)

			var (
				req = require.New(t)

				f   = types.RecordFilter{}
				err error
			)

			f.Sort.Set("strMulti DESC")
			f.Limit = 100

			_, f, err = store.SearchComposeRecords(ctx, s, mod, f)
			req.Error(err, "not allowed to sort by multi-value fields: strMulti")
		})

		t.Run("advanced sorting; NULL; record value + sys fields", func(t *testing.T) {
			var (
				_, _ = truncAndCreate(t,
					makeNewUpd(t, "2021-01-01T01:00:00Z", &types.RecordValue{Name: "str1", Value: "a"}),
					makeNewUpd(t, "" /* padding       */, &types.RecordValue{Name: "str1", Value: "b"}),
					makeNewUpd(t, "2021-01-01T02:00:00Z"),
					makeNewUpd(t, "" /* padding       */),
				)
			)

			tcc := []tc{
				{
					"updatedAt, str1",
					[]string{"<NIL>,<NIL>;<NIL>,b", "2021-01-01T01:00:00Z,a;2021-01-01T02:00:00Z,<NIL>"},
					[]int{nextCur, prevCur},
				},
				{
					"updatedAt, str1 DESC",
					[]string{"<NIL>,b;<NIL>,<NIL>", "2021-01-01T01:00:00Z,a;2021-01-01T02:00:00Z,<NIL>"},
					[]int{nextCur, prevCur},
				},
				{
					"updatedAt DESC, str1",
					[]string{"2021-01-01T02:00:00Z,<NIL>;2021-01-01T01:00:00Z,a", "<NIL>,<NIL>;<NIL>,b"},
					[]int{nextCur, prevCur},
				},
				{
					"updatedAt DESC, str1 DESC",
					[]string{"2021-01-01T02:00:00Z,<NIL>;2021-01-01T01:00:00Z,a", "<NIL>,b;<NIL>,<NIL>"},
					[]int{nextCur, prevCur},
				},
			}

			stringifyMix := func(set types.RecordSet, fields ...string) string {
				var out string
				for r := range set {
					if r > 0 {
						out += ";"
					}

					for f := range fields {
						if f > 0 {
							out += ","
						}

						if fields[f] == "updatedAt" {
							v := set[r].UpdatedAt
							if v != nil && !v.IsZero() {
								out += v.Format(time.RFC3339)
							} else {
								out += "<NIL>"
							}
						} else {
							v := set[r].Values.Get(fields[f], 0)
							if v != nil {
								out += v.Value
							} else {
								out += "<NIL>"
							}
						}
					}
				}
				return out
			}

			for _, tc := range tcc {
				t.Run("crawling: "+tc.sort, func(t *testing.T) {
					var (
						req = require.New(t)

						f   = types.RecordFilter{}
						set types.RecordSet
						err error
					)

					f.Sort.Set(tc.sort)
					f.Limit = 2

					for p := 0; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))
						req.Equal(tc.rval[p], stringifyMix(set, "updatedAt", "str1"))

						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}

					f.PageCursor = f.PrevPage
					for p := 0; p >= 0; p-- {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))

						req.Equal(tc.rval[p], stringifyMix(set, "updatedAt", "str1"))
						testCursors(req, tc.curr[p], f)

						// reverse to previous page
						f.PageCursor = f.PrevPage
					}

					f.PageCursor = f.NextPage

					for p := 1; p < 2; p++ {
						set, f, err = store.SearchComposeRecords(ctx, s, mod, f)
						req.NoError(err)
						req.True(tc.sort == f.Sort.String() || strings.HasPrefix(f.Sort.String(), tc.sort+","))
						req.Equal(tc.rval[p], stringifyMix(set, "updatedAt", "str1"))

						testCursors(req, tc.curr[p], f)

						// advance to next page
						f.PageCursor = f.NextPage
					}
				})
			}
		})

	})

	t.Run("sort by system field, paged", func(t *testing.T) {
		var (
			req, _ = truncAndCreate(t,
				makeNew(&types.RecordValue{Name: "str1", Value: "v1"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v2"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v3"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v4"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v5"}),
			)

			err error
			set types.RecordSet
			f   = types.RecordFilter{
				ModuleID:    mod.ID,
				NamespaceID: mod.NamespaceID,
			}
		)

		// sorting by system field

		f.Limit = 2

		req.NoError(f.Sort.Set("str1, createdAt"))
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v1;v2", stringifyValues(set, "str1"))
		req.NotNil(f.NextPage)
		req.Nil(f.PrevPage)

		f.PageCursor = f.NextPage
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v3;v4", stringifyValues(set, "str1"))
		req.NotNil(f.PrevPage)
		req.NotNil(f.NextPage)

		f.PageCursor = f.NextPage
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v5", stringifyValues(set, "str1"))
		req.NotNil(f.PrevPage)
		req.Nil(f.NextPage)

		f.PageCursor = f.PrevPage
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v3;v4", stringifyValues(set, "str1"))
		req.NotNil(f.PrevPage)
		req.NotNil(f.NextPage)

		f.Limit = 1
		f.PageCursor = f.PrevPage
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v2", stringifyValues(set, "str1"))
		req.NotNil(f.PrevPage)
		req.NotNil(f.NextPage)

		f.PageCursor = f.NextPage
		f.Limit = 3
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v3;v4;v5", stringifyValues(set, "str1"))
		req.Nil(f.NextPage)
		req.NotNil(f.PrevPage) // we can't actually know if we're on the last page or not..

		// sorting by module field

		f.PageCursor = nil
		f.Limit = 2
		req.NoError(f.Sort.Set("str1 DESC"))
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v5;v4", stringifyValues(set, "str1"))
		req.NotNil(f.NextPage)
		req.Nil(f.PrevPage)

		f.PageCursor = f.NextPage
		f.Sort = nil
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v3;v2", stringifyValues(set, "str1"))
		req.NotNil(f.NextPage)
		req.NotNil(f.PrevPage)

		f.PageCursor = f.PrevPage
		f.Sort = nil
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v5;v4", stringifyValues(set, "str1"))
		req.NotNil(f.NextPage)
		req.Nil(f.PrevPage)
	})

	t.Run("paged", func(t *testing.T) {
		var (
			err error
			set types.RecordSet

			req, _ = truncAndCreate(t,
				makeNew(&types.RecordValue{Name: "str1", Value: "v1"}, &types.RecordValue{Name: "str3", Value: "w3.1"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v2"}, &types.RecordValue{Name: "str3", Value: "w3.2"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v3"}, &types.RecordValue{Name: "str3", Value: "w3.3"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v4"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v5"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v6"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v7"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v8"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v9"}),
			)

			f = types.RecordFilter{
				ModuleID:    mod.ID,
				NamespaceID: mod.NamespaceID,
			}
		)

		req.NoError(f.Sort.Set("str1"))
		f.Limit = 3
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.NotNil(f.NextPage)
		req.Nil(f.PrevPage)
		req.Equal("v1,w3.1;v2,w3.2;v3,w3.3", stringifyValues(set, "str1", "str3"))

		f.PageCursor = f.NextPage
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("v4,<NIL>;v5,<NIL>;v6,<NIL>", stringifyValues(set, "str1", "str3"))

		req.NoError(f.Sort.Set("str3 DESC, str1 ASC"))
		f.PageCursor = nil
		f.Limit = 1
		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.Equal("w3.3,v3", stringifyValues(set, "str3", "str1"))

		f.Limit = 1
		req.NoError(f.Sort.Set("str1"))

		validOrderOfStr1 := []int{0, 1, 2, 3, 4, 5, 6, 7}

		f.PageCursor = nil
		for p := 1; p <= 5; p++ {
			if p > 1 {
				f.Sort = nil
			}
			set, f, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			req.Equal(fmt.Sprintf("v%d", validOrderOfStr1[p]), stringifyValues(set, "str1"))
			//fmt.Printf("%-30s\t%v\t%30s\t%30s\n", stringifyValues(set, "str1", "str3"), f.Sort, f.NextPage, f.PrevPage)
			f.PageCursor = f.NextPage
		}

		f.PageCursor = f.PrevPage
		for p := 4; p >= 1; p-- {
			f.Sort = nil
			set, f, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			//fmt.Printf("%-30s\t%v\t%30s\t%30s\n", stringifyValues(set, "str1", "str3"), f.Sort, f.NextPage, f.PrevPage)
			req.Equal(fmt.Sprintf("v%d", validOrderOfStr1[p]), stringifyValues(set, "str1"))
			f.PageCursor = f.PrevPage
		}

		f.Limit = 1

		req.NoError(f.Sort.Set("str3 DESC, id DESC"))
		f.PageCursor = nil

		validOrderOfStr1 = []int{0, 3, 2, 1, 9, 8, 7}

		for p := 1; p <= 6; p++ {
			if p > 1 {
				f.Sort = nil
			}
			set, f, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			//fmt.Printf("%-30s\t%v\t%30s\t%30s\n", stringifyValues(set, "str3", "str1"), f.Sort, f.NextPage, f.PrevPage)
			req.Equal(fmt.Sprintf("v%d", validOrderOfStr1[p]), stringifyValues(set, "str1"))
			f.PageCursor = f.NextPage
		}

		f.PageCursor = f.PrevPage
		for p := 5; p >= 1; p-- {
			f.Sort = nil
			set, f, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			//fmt.Printf("%-30s\t%v\t%30s\t%30s\n", stringifyValues(set, "str3", "str1"), f.Sort, f.NextPage, f.PrevPage)
			req.Equal(fmt.Sprintf("v%d", validOrderOfStr1[p]), stringifyValues(set, "str1"))
			f.PageCursor = f.PrevPage
		}

		req.NoError(f.Sort.Set("str3 ASC, id ASC"))
		f.PageCursor = nil

		validOrderOfStr1 = []int{0, 4, 5, 6, 7, 8, 9, 1, 2}

		for p := 1; p <= len(validOrderOfStr1)-1; p++ {
			if p > 1 {
				f.Sort = nil
			}
			set, f, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			//fmt.Printf("%-30s\t%v\t%30s\t%30s\n", stringifyValues(set, "str3", "str1"), f.Sort, f.NextPage, f.PrevPage)
			req.Equal(fmt.Sprintf("v%d", validOrderOfStr1[p]), stringifyValues(set, "str1"))
			f.PageCursor = f.NextPage
		}

		f.PageCursor = f.PrevPage
		for p := len(validOrderOfStr1) - 2; p >= 1; p-- {
			f.Sort = nil
			set, f, err = s.SearchComposeRecords(ctx, mod, f)
			req.NoError(err)
			//fmt.Printf("%-30s\t%v\t%30s\t%30s\n", stringifyValues(set, "str3", "str1"), f.Sort, f.NextPage, f.PrevPage)
			req.Equal(fmt.Sprintf("v%d", validOrderOfStr1[p]), stringifyValues(set, "str1"))
			f.PageCursor = f.PrevPage
		}
	})

	t.Run("page navigation", func(t *testing.T) {
		var (
			err error
			set types.RecordSet

			req, _ = truncAndCreate(t,
				makeNew(&types.RecordValue{Name: "str1", Value: "v1"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v2"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v3"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v4"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v5"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v6"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v7"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v8"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "v9"}),
			)

			f = types.RecordFilter{
				ModuleID:    mod.ID,
				NamespaceID: mod.NamespaceID,
			}
		)

		// ascending

		req.NoError(f.Sort.Set("str1"))
		f.Limit = 2
		f.IncPageNavigation = true
		f.IncTotal = true

		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.NotNil(set)
		req.NotNil(f.PageNavigation)
		req.Equal(uint(9), f.Total)
		req.Len(f.PageNavigation, 5)

		f.PageCursor = f.PageNavigation[1].Cursor
		f.IncPageNavigation = false
		f.IncTotal = false
		set, _, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.NotNil(set)
		req.Equal("v3;v4", stringifyValues(set, "str1"))

		f.PageCursor = f.PageNavigation[4].Cursor
		f.IncPageNavigation = false
		f.IncTotal = false
		set, _, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.NotNil(set)
		req.Equal("v9", stringifyValues(set, "str1"))

		// descending

		req.NoError(f.Sort.Set("str1 DESC"))
		f.Limit = 2
		f.PageCursor = nil
		f.IncPageNavigation = true
		f.IncTotal = true

		set, f, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.NotNil(set)
		req.NotNil(f.PageNavigation)
		req.Equal(uint(9), f.Total)
		req.Len(f.PageNavigation, 5)

		f.PageCursor = f.PageNavigation[1].Cursor
		f.IncPageNavigation = false
		f.IncTotal = false
		set, _, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.NotNil(set)
		req.Equal("v7;v6", stringifyValues(set, "str1"))

		f.PageCursor = f.PageNavigation[4].Cursor
		f.IncPageNavigation = false
		f.IncTotal = false
		set, _, err = s.SearchComposeRecords(ctx, mod, f)
		req.NoError(err)
		req.NotNil(set)
		req.Equal("v1", stringifyValues(set, "str1"))
	})

	t.Run("report", func(t *testing.T) {
		var (
			err error

			req, _ = truncAndCreate(t,
				makeNew(&types.RecordValue{Name: "dt1", Value: "2020-01-01"}, &types.RecordValue{Name: "num1", Value: "1"}, &types.RecordValue{Name: "str3", Value: "three"}),
				makeNew(&types.RecordValue{Name: "dt1", Value: "2020-01-01"}, &types.RecordValue{Name: "num1", Value: "2"}, &types.RecordValue{Name: "str3", Value: "three"}),
				makeNew(&types.RecordValue{Name: "dt1", Value: "2020-01-01"}, &types.RecordValue{Name: "num1", Value: "3"}, &types.RecordValue{Name: "str3", Value: "three"}),
				makeNew(&types.RecordValue{Name: "dt1", Value: "2020-05-01"}, &types.RecordValue{Name: "num1", Value: "4"}),
				makeNew(&types.RecordValue{Name: "dt1", Value: "2020-05-01"}, &types.RecordValue{Name: "num1", Value: "5"}),

				// Add one additional record with deleted values
				makeNew(&types.RecordValue{Name: "dt1", Value: "2020-05-01", DeletedAt: now()}, &types.RecordValue{Name: "num1", Value: "6", DeletedAt: now()}, &types.RecordValue{Name: "str2", Value: "deleted", DeletedAt: now()}),
			)

			report []map[string]interface{}
		)

		report, err = s.ComposeRecordReport(ctx, mod, "MAX(num1)", "QUARTER(dt1)", "")
		req.NoError(err)
		req.Len(report, 3)

		// @todo find a way to compare the results

		//expected := []map[string]interface{}{
		//	{"count": 3, "dimension_0": 1, "metric_0": 3},
		//	{"count": 2, "dimension_0": 2, "metric_0": 5},
		//	{"count": 1, "dimension_0": nil, "metric_0": nil},
		//}
		//
		//req.True(
		//	reflect.DeepEqual(report, expected),
		//	"report does not match expected results:\n%#v\n%#v", report, expected)

		report, err = s.ComposeRecordReport(ctx, mod, "COUNT(num1)", "YEAR(dt1)", "")
		req.NoError(err)

		report, err = s.ComposeRecordReport(ctx, mod, "SUM(num1)", "DATE(dt1)", "")
		req.NoError(err)

		report, err = s.ComposeRecordReport(ctx, mod, "MIN(num1)", "DATE(NOW())", "")
		req.NoError(err)

		report, err = s.ComposeRecordReport(ctx, mod, "AVG(num1)", "DATE(NOW())", "")
		req.NoError(err)

		// Note that not all functions are compatible across all backends
	})

	t.Run("partial value update", func(t *testing.T) {
		var (
			err error
			set types.RecordSet

			req, _ = truncAndCreate(t,
				makeNew(&types.RecordValue{Name: "str1", Value: "1st"}, &types.RecordValue{Name: "num1", Value: "1"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "2nd"}, &types.RecordValue{Name: "num1", Value: "2"}),
				makeNew(&types.RecordValue{Name: "str1", Value: "3rd"}, &types.RecordValue{Name: "num1", Value: "3"}),
			)
		)

		set, _, err = s.SearchComposeRecords(ctx, mod, types.RecordFilter{})
		req.NoError(err)
		req.Equal("1st,1;2nd,2;3rd,3", stringifyValues(set, "str1", "num1"))

		ups := &types.RecordValue{RecordID: set[1].ID, Name: "num1", Value: "22"}
		req.NoError(s.PartialComposeRecordValueUpdate(ctx, mod, ups))

		set, _, err = s.SearchComposeRecords(ctx, mod, types.RecordFilter{})
		req.NoError(err)
		req.Equal("1st,1;2nd,22;3rd,3", stringifyValues(set, "str1", "num1"))

	})
}
