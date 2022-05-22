package dal

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func RecordSearch(t *testing.T, ctx context.Context, d dal.Connection) {
	const (
		totalRecords = 10
	)

	var (
		req = require.New(t)

		m = &dal.Model{
			Ident: "crs_test_search",
			Attributes: dal.AttributeSet{
				&dal.Attribute{Ident: "ID", Type: &dal.TypeID{}, Store: &dal.CodecAlias{Ident: "id"}, PrimaryKey: true},
				&dal.Attribute{Ident: "createdAt", Type: &dal.TypeTimestamp{}, Store: &dal.CodecAlias{Ident: "created_at"}},
				&dal.Attribute{Ident: "updatedAt", Type: &dal.TypeTimestamp{}, Store: &dal.CodecAlias{Ident: "updated_at"}},

				&dal.Attribute{Ident: "v_string", Filterable: true, Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "v_number", Filterable: true, Type: &dal.TypeNumber{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "v_is_odd", Filterable: true, Type: &dal.TypeBoolean{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "p_string", Filterable: true, Type: &dal.TypeText{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "p_number", Filterable: true, Type: &dal.TypeNumber{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "p_is_odd", Filterable: true, Type: &dal.TypeBoolean{}, Store: &dal.CodecPlain{}},
			},
		}
	)

	for ID := uint64(1); ID <= totalRecords; ID++ {
		r := &types.Record{ID: ID, CreatedAt: time.Now()}

		i := int(ID)
		r.Values = r.Values.Set(&types.RecordValue{Name: "v_string", Value: "tens_" + strconv.Itoa(i%10)})
		r.Values = r.Values.Set(&types.RecordValue{Name: "v_number", Value: strconv.Itoa(i)})
		r.Values = r.Values.Set(&types.RecordValue{Name: "v_is_odd", Value: strconv.FormatBool(i%2 == 1)})
		r.Values = r.Values.Set(&types.RecordValue{Name: "p_string", Value: "tens_" + strconv.Itoa(i%10)})
		r.Values = r.Values.Set(&types.RecordValue{Name: "p_number", Value: strconv.Itoa(i)})
		r.Values = r.Values.Set(&types.RecordValue{Name: "p_is_odd", Value: strconv.FormatBool(i%2 == 1)})

		req.NoError(d.Create(ctx, m, r))
	}

	cases := []struct {
		f     types.RecordFilter
		total int
	}{
		{
			total: totalRecords,
		},
		{
			f:     types.RecordFilter{Query: "v_string == p_string"},
			total: totalRecords,
		},
		{
			f:     types.RecordFilter{Query: "v_number == p_number"},
			total: totalRecords,
		},
		{
			f:     types.RecordFilter{Query: "p_is_odd"},
			total: totalRecords / 2,
		},
		{
			f:     types.RecordFilter{Query: "true = p_is_odd"},
			total: totalRecords / 2,
		},
		{
			f:     types.RecordFilter{Query: "p_is_odd = true"},
			total: totalRecords / 2,
		},
		{
			f:     types.RecordFilter{Query: "!p_is_odd"},
			total: totalRecords / 2,
		},
		{
			f:     types.RecordFilter{Query: "p_number = 1"},
			total: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.f.Query, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			i, err := d.Search(ctx, m, c.f.ToFilter())
			req.NoError(err)

			rr, err := drain(ctx, i)
			req.NoError(err)
			req.Len(rr, c.total)
		})
	}

	t.Run("paging", func(t *testing.T) {
		var (
			req      = require.New(t)
			ids      string
			fwd, bck *filter.PagingCursor

			search = func(where, orderBy string, lim uint, cur *filter.PagingCursor) (ids string, fwd, bck *filter.PagingCursor) {
				f := types.RecordFilter{Query: where}
				f.PageCursor = cur
				f.Limit = lim
				req.NoError(f.Sort.Set(orderBy))
				i, err := d.Search(ctx, m, f.ToFilter())
				req.NoError(err)
				req.NoError(i.Err())

				if !i.Next(ctx) {
					req.NoError(i.Err())
					return
				}

				r := new(types.Record)
				req.NoError(i.Scan(r))

				bck, err = i.BackCursor(r)
				req.NoError(err)
				t.Logf("bck-cursor (from the 1st fetched record): %v", bck)

				rr, err := drain(ctx, i)
				req.NoError(err)

				if len(rr) > 0 {
					fwd, err = i.ForwardCursor(rr[len(rr)-1])
					req.NoError(err)
					t.Logf("fwd-cursor (from the lst fetched record): %v", fwd)
				}

				ids = fmt.Sprintf("%d", r.ID)
				for _, r = range rr {
					ids += fmt.Sprintf(",%d", r.ID)
				}

				return
			}
		)

		ids, fwd, _ = search("", "", 3, nil)
		req.Equal("1,2,3", ids)
		ids, fwd, _ = search("", "", 3, fwd)
		req.Equal("4,5,6", ids)
		ids, fwd, _ = search("", "", 3, fwd)
		req.Equal("7,8,9", ids)
		ids, _, bck = search("", "", 3, fwd)
		req.Equal("10", ids)

		ids, _, bck = search("", "", 3, bck)
		req.Equal("7,8,9", ids)
		ids, _, bck = search("", "", 3, bck)
		req.Equal("4,5,6", ids)
		ids, _, bck = search("", "", 3, bck)
		req.Equal("1,2,3", ids)

		ids, fwd, _ = search("", "p_is_odd", 3, nil)
		req.Equal("2,4,6", ids)
		ids, _, bck = search("", "", 3, fwd)
		req.Equal("8,10,1", ids)
		ids, _, _ = search("", "", 3, bck)
		req.Equal("2,4,6", ids)

		ids, fwd, _ = search("", "v_is_odd", 3, nil)
		req.Equal("2,4,6", ids)
		ids, _, bck = search("", "", 3, fwd)
		req.Equal("8,10,1", ids)
		ids, _, _ = search("", "", 3, bck)
		req.Equal("2,4,6", ids)

		_, _ = fwd, bck // avoiding unused var. error
	})
}
