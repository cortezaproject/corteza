package test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/crs"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/stretchr/testify/require"
)

func All(t *testing.T, d crs.StoreConnection) {
	t.Run("RecordCodec", func(t *testing.T) { RecordCodec(t, d) })
	t.Run("RecordSearch", func(t *testing.T) { RecordSearch(t, d) })
}

func RecordCodec(t *testing.T, d crs.StoreConnection) {
	var (
		req = require.New(t)

		// enable query logging when +debug is used on DSN schema
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

		m = &data.Model{
			Ident: "crs_test_codec",
			Attributes: data.AttributeSet{
				&data.Attribute{Ident: data.SysID, Type: &data.TypeID{}, Store: &data.StoreCodecAlias{Ident: "id"}, PrimaryKey: true},
				&data.Attribute{Ident: data.SysCreatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "created_at"}},
				&data.Attribute{Ident: data.SysUpdatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "updated_at"}},

				&data.Attribute{Ident: "vID", Type: &data.TypeID{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vRef", Type: &data.TypeRef{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vTimestamp", Type: &data.TypeTimestamp{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vTime", Type: &data.TypeTime{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vDate", Type: &data.TypeDate{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vNumber", Type: &data.TypeNumber{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vText", Type: &data.TypeText{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vBoolean_T", Type: &data.TypeBoolean{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vBoolean_F", Type: &data.TypeBoolean{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vEnum", Type: &data.TypeEnum{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vGeometry", Type: &data.TypeGeometry{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vJSON", Type: &data.TypeJSON{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vBlob", Type: &data.TypeBlob{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "vUUID", Type: &data.TypeUUID{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "pID", Type: &data.TypeID{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pRef", Type: &data.TypeRef{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pTimestamp_TZT", Type: &data.TypeTimestamp{Timezone: true}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pTimestamp_TZF", Type: &data.TypeTimestamp{Timezone: false}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pTime", Type: &data.TypeTime{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pDate", Type: &data.TypeDate{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pNumber", Type: &data.TypeNumber{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pText", Type: &data.TypeText{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pBoolean_T", Type: &data.TypeBoolean{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pBoolean_F", Type: &data.TypeBoolean{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pEnum", Type: &data.TypeEnum{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pGeometry", Type: &data.TypeGeometry{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pJSON", Type: &data.TypeJSON{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pBlob", Type: &data.TypeBlob{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "pUUID", Type: &data.TypeUUID{}, Store: &data.StoreCodecPlain{}},
			},
		}

		rIn  = types.Record{ID: 42}
		err  error
		rOut *types.Record

		piTime time.Time
	)

	piTime, err = time.Parse("2006-01-02T15:04:05", "2006-01-02T15:04:05")
	req.NoError(err)
	piTime = piTime.UTC()

	rIn.CreatedAt = piTime
	rIn.UpdatedAt = &piTime

	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vID", Value: "34324"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vRef", Value: "32423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vTimestamp", Value: "2022-01-01T10:10:10+02:00"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vTime", Value: "04:10:10+04:00"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vDate", Value: "2022-01-01"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vNumber", Value: "2423423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vText", Value: "lorm ipsum "})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vBoolean_T", Value: "true"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vBoolean_F", Value: "false"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vEnum", Value: "abc"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vGeometry", Value: `{"lat":1,"lng":1}`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vJSON", Value: `[{"bool":true"}]`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vBlob", Value: "0110101"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vUUID", Value: "ba485865-54f9-44de-bde8-6965556c022a"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pID", Value: "34324"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pRef", Value: "32423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pTimestamp_TZF", Value: "2022-02-01T10:10:10"})

	// @todo how (if at all) should we know if underlying DB supports timezone?
	//rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pTimestamp_TZT", Value: "2022-02-01T10:10:10"})

	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pTime", Value: "06:06:06"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pDate", Value: "2022-01-01"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pNumber", Value: "2423423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pText", Value: "lorm ipsum "})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pBoolean_T", Value: "true"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pBoolean_F", Value: "false"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pEnum", Value: "abc"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pGeometry", Value: `{"lat":1,"lng":1}`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pJSON", Value: `[{"bool":true"}]`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pBlob", Value: "0110101"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pUUID", Value: "ba485865-54f9-44de-bde8-6965556c022a"})
	rIn.Values = rIn.Values.GetClean()

	req.NoError(d.CreateRecords(ctx, m, &rIn))

	rOut = new(types.Record)
	req.NoError(d.LookupRecord(ctx, m, crs.PKValues{"id": rIn.ID}, rOut))

	{
		// normalize timezone on timestamps
		rOut.CreatedAt = rOut.CreatedAt.UTC()
		aux := rOut.UpdatedAt.UTC()
		rOut.UpdatedAt = &aux
	}

	for _, attr := range m.Attributes {
		vIn, err := rIn.GetValue(attr.Ident, 0)
		req.NoError(err)
		vOut, err := rOut.GetValue(attr.Ident, 0)
		req.NoError(err)
		req.Equal(vIn, vOut, "values for attribute %q are not equal", attr.Ident)
	}
}

func RecordSearch(t *testing.T, d crs.StoreConnection) {
	const (
		totalRecords = 10
	)

	var (
		req = require.New(t)

		// enable query logging when +debug is used on DSN schema
		ctx = logger.ContextWithValue(context.Background(), logger.MakeDebugLogger())

		m = &data.Model{
			Ident: "crs_test_search",
			Attributes: data.AttributeSet{
				&data.Attribute{Ident: data.SysID, Type: &data.TypeID{}, Store: &data.StoreCodecAlias{Ident: "id"}, PrimaryKey: true},
				&data.Attribute{Ident: data.SysCreatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "created_at"}},
				&data.Attribute{Ident: data.SysUpdatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "updated_at"}},

				&data.Attribute{Ident: "v_string", Type: &data.TypeText{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "v_number", Type: &data.TypeNumber{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "v_is_odd", Type: &data.TypeBoolean{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
				&data.Attribute{Ident: "p_string", Type: &data.TypeText{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "p_number", Type: &data.TypeNumber{}, Store: &data.StoreCodecPlain{}},
				&data.Attribute{Ident: "p_is_odd", Type: &data.TypeBoolean{}, Store: &data.StoreCodecPlain{}},
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

		req.NoError(d.CreateRecords(ctx, m, r))
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

			i, err := d.SearchRecords(ctx, m, c.f)
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
				i, err := d.SearchRecords(ctx, m, f)
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

func drain(ctx context.Context, i crs.Iterator) (rr []*types.Record, err error) {
	var r *types.Record
	rr = make([]*types.Record, 0, 100)
	for i.Next(ctx) {
		if i.Err() != nil {
			return nil, i.Err()
		}

		r = new(types.Record)
		if err = i.Scan(r); err != nil {
			return
		}

		rr = append(rr, r)
	}

	return rr, i.Err()
}
