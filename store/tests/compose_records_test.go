package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type (
	// For
	composeRecordsStore interface {
		SearchComposeRecords(ctx context.Context, m *types.Module, f types.RecordFilter) (types.RecordSet, types.RecordFilter, error)
		LookupComposeRecordByID(ctx context.Context, m *types.Module, id uint64) (*types.Record, error)
		CreateComposeRecord(ctx context.Context, m *types.Module, rr ...*types.Record) error
		//CreateComposeRecordValue(ctx context.Context, m *types.Module, rr ...*types.RecordValue) error
		UpdateComposeRecord(ctx context.Context, m *types.Module, rr ...*types.Record) error
		//UpdateComposeRecordValue(ctx context.Context, m *types.Module, rr ...*types.RecordValue) error
		//PartialUpdateComposeRecord(ctx context.Context, m *types.Module, onlyColumns []string, rr ...*types.Record) error
		//PartialUpdateComposeRecordValue(ctx context.Context, m *types.Module, rr ...*types.RecordValue) error
		RemoveComposeRecord(ctx context.Context, m *types.Module, rr ...*types.Record) error
		//RemoveComposeRecordValue(ctx context.Context, m *types.Module, rr ...*types.RecordValue) error
		//RemoveComposeRecordByID(ctx context.Context, m *types.Module, ID uint64) error
		//RemoveComposeRecordValueByRecordID(ctx context.Context, m *types.Module, recordID uint64, place int, name string) error
		TruncateComposeRecords(ctx context.Context, m *types.Module) error
		//TruncateComposeRecordValues(ctx context.Context, m *types.Module) error
		//ExecUpdateComposeRecords(ctx context.Context, m *types.Module, cnd squirrel.Sqlizer, set store.Payload) error
		//ExecUpdateComposeRecordValues(ctx context.Context, m *types.Module, cnd squirrel.Sqlizer, set store.Payload) error
		//ComposeRecordLookup(ctx context.Context, m *types.Module, cnd squirrel.Sqlizer) (*types.Record, error)
		//ComposeRecordValueLookup(ctx context.Context, m *types.Module, cnd squirrel.Sqlizer) (*types.RecordValue, error)
		//QueryComposeRecords(m *types.Module) squirrel.SelectBuilder
		//QueryComposeRecordValues(m *types.Module) squirrel.SelectBuilder
	}
)

func testComposeRecords(t *testing.T, s composeRecordsStore) {
	var (
		ctx = context.Background()
		req = require.New(t)

		mod = &types.Module{
			ID:          id.Next(),
			NamespaceID: id.Next(),
			Handle:      "",
			Name:        "testComposeRecords",
			CreatedAt:   time.Now(),
		}

		makeNew = func() *types.Record {
			// minimum data set for new composeRecord
			return &types.Record{
				ID:          id.Next(),
				NamespaceID: mod.NamespaceID,
				ModuleID:    mod.ID,
				CreatedAt:   time.Now(),
			}
		}
	)

	t.Run("create", func(t *testing.T) {
		composeRecord := makeNew()
		req.NoError(s.CreateComposeRecord(ctx, mod, composeRecord))
	})

	t.Run("create with duplicate handle", func(t *testing.T) {
		t.Skip("not implemented")
	})

	t.Run("lookup by ID", func(t *testing.T) {
		composeRecord := makeNew()
		req.NoError(s.CreateComposeRecord(ctx, mod, composeRecord))
		fetched, err := s.LookupComposeRecordByID(ctx, mod, composeRecord.ID)
		req.NoError(err)
		req.Equal(composeRecord.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("remove", func(t *testing.T) {
		composeRecord := makeNew()
		req.NoError(s.CreateComposeRecord(ctx, mod, composeRecord))
		req.NoError(s.RemoveComposeRecord(ctx, mod))
	})

	t.Run("remove by ID", func(t *testing.T) {
		composeRecord := makeNew()
		req.NoError(s.CreateComposeRecord(ctx, mod, composeRecord))
		req.NoError(s.RemoveComposeRecord(ctx, mod))
	})

	t.Run("update", func(t *testing.T) {
		composeRecord := makeNew()
		req.NoError(s.CreateComposeRecord(ctx, mod, composeRecord))

		composeRecord = &types.Record{
			ID:        composeRecord.ID,
			CreatedAt: composeRecord.CreatedAt,
			OwnedBy:   1,
		}
		req.NoError(s.UpdateComposeRecord(ctx, mod, composeRecord))

		updated, err := s.LookupComposeRecordByID(ctx, mod, composeRecord.ID)
		req.NoError(err)
		req.Equal(composeRecord.OwnedBy, updated.OwnedBy)
	})

	t.Run("search", func(t *testing.T) {
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

		req.NoError(s.TruncateComposeRecords(ctx, mod))
		req.NoError(s.CreateComposeRecord(ctx, mod, prefill...))

		// search for all valid
		set, f, err := s.SearchComposeRecords(ctx, mod, types.RecordFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one
		req.Equal(valid, int(f.Count))

		// search for ALL
		set, f, err = s.SearchComposeRecords(ctx, mod, types.RecordFilter{Deleted: rh.FilterStateInclusive})
		req.NoError(err)
		req.Len(set, count) // we've deleted one

		// search for deleted only
		set, f, err = s.SearchComposeRecords(ctx, mod, types.RecordFilter{Deleted: rh.FilterStateExclusive})
		req.NoError(err)
		req.Len(set, 1) // we've deleted one
	})

	t.Run("filtered search", func(t *testing.T) {
		// testComposeRecordsSearch(t, s)
	})
}

//func testComposeRecordsSearch(t *testing.T, s store.ComposeRecords) {
//	m := &types.Module{
//		ID:          123,
//		NamespaceID: 456,
//		Fields: types.ModuleFieldSet{
//			&types.ModuleField{Name: "foo"},
//			&types.ModuleField{Name: "bar"},
//			&types.ModuleField{Name: "booly", Kind: "Bool"},
//		},
//	}
//
//	ttc := []struct {
//		name    string
//		f       types.RecordFilter
//		match   []string
//		noMatch []string
//		args    []interface{}
//		err     string
//	}{
//		{
//			name: "default filter",
//			match: []string{
//				"SELECT r.id, r.module_id, r.rel_namespace, r.owned_by, r.created_at, " +
//					"r.created_by, r.updated_at, r.updated_by, r.deleted_at, r.deleted_by " +
//					"FROM compose_record AS r " +
//					"WHERE r.module_id = ? AND r.rel_namespace = ? AND r.deleted_at IS NULL",
//			},
//		},
//		{
//			name: "simple query",
//			f:    types.RecordFilter{Query: "id = 5 AND foo = 7"},
//			match: []string{
//				"r.id  = 5",
//				"rv_foo.value  = 7"},
//			args: []interface{}{"foo"},
//		},
//		{
//			name: "sorting",
//			f:    types.RecordFilter{Sort: "id ASC, bar DESC"},
//			match: []string{
//				" r.id ASC",
//				" rv_bar.value DESC",
//			},
//			args: []interface{}{"bar"},
//		},
//		{
//			name:  "exclude deleted records (def. behaviour)",
//			f:     types.RecordFilter{Deleted: rh.FilterStateExcluded},
//			match: []string{" r.deleted_at IS "},
//		},
//		{
//			name:    "include deleted records",
//			f:       types.RecordFilter{Deleted: rh.FilterStateInclusive},
//			noMatch: []string{" r.deleted_at IS NULL "},
//		},
//		{
//			name:  "only deleted record",
//			f:     types.RecordFilter{Deleted: rh.FilterStateExclusive},
//			match: []string{" r.deleted_at IS NOT NULL"},
//		},
//		{
//			name:  "boolean",
//			f:     types.RecordFilter{Query: "booly"},
//			match: []string{"(rv_booly.value NOT IN ("},
//			args:  []interface{}{"booly"},
//		},
//	}
//
//	for _, tc := range ttc {
//		t.Run(tc.name, func(t *testing.T) {
//			req := require.New(t)
//			sb, err := s.query(m, tc.f)
//
//			if tc.err != nil {
//				req.Error(err, tc.err, "buildQuery(%+v) did not return an expected error %q but %q", tc.f, tc.err, err)
//			} else {
//				req.NoError(err,"buildQuery(%+v) returned an unexpected error: %v", tc.f, err)
//			}
//
//			sql, args, err := sb.ToSql()
//
//			for _, m := range tc.match {
//				require.True(t, strings.Contains(sql, m),
//					"assertion failed; query %q \n  "+
//						"             did not contain  %q", sql, m)
//			}
//
//			for _, m := range tc.noMatch {
//				require.False(t, strings.Contains(sql, m),
//					"assertion failed; query %q \n  "+
//						"             must not contain  %q", sql, m)
//			}
//
//			tc.args = append(tc.args, m.ID, m.NamespaceID)
//			require.True(t, fmt.Sprintf("%+v", args) == fmt.Sprintf("%+v", tc.args),
//				"assertion failed; args %+v \n  "+
//					"     do not match expected %+v", args, tc.args)
//		})
//	}
//}
