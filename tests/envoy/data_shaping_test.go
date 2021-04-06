package envoy

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

func TestDataShaping(t *testing.T) {
	var (
		ctx = auth.SetSuperUserContext(context.Background())
		s   = initStore(ctx, t)
		err error

		cases = []string{
			"csv_simple",
			"jsonl_simple",

			"csv_selfref",
			"jsonl_selfref",
		}
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("record shaping; data_shaping/%s", c), func(t *testing.T) {
			var (
				req = require.New(t)
			)

			truncateStore(ctx, s, t)
			err = collect(
				err,
				storeRole(ctx, s, 1, "everyone"),
				storeRole(ctx, s, 2, "admins"),
			)
			if err != nil {
				t.Fatal(err.Error())
			}

			nn, err := decodeDirectory(ctx, path.Join("data_shaping", c))
			req.NoError(err)

			crs := resource.ComposeRecordShaper()
			nn, err = resource.Shape(nn, crs)
			req.NoError(err)

			req.NoError(encode(ctx, s, nn))

			ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
			req.NotNil(ns)
			ms, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod1")
			req.NotNil(ms)

			rr, _, err := store.SearchComposeRecords(ctx, s, ms, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr, 2)

			r1 := rr[0]
			r2 := rr[1]

			req.Len(r1.Values, 2)
			req.Equal("c1.v1", r1.Values.FilterByName("f1")[0].Value)
			req.Equal("c2.v1", r1.Values.FilterByName("f2")[0].Value)

			req.Len(r2.Values, 2)
			req.Equal("c1.v2", r2.Values.FilterByName("f1")[0].Value)
			req.Equal("c2.v2", r2.Values.FilterByName("f2")[0].Value)

			s.TruncateComposeRecords(ctx, ms)
		})
	}
}

func TestDataShaping_fieldTypes(t *testing.T) {
	var (
		ctx = auth.SetSuperUserContext(context.Background())
		s   = initStore(ctx, t)
		err error

		cases = []string{
			"csv_fieldtypes",
			"jsonl_fieldtypes",
		}
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("record shaping; data_shaping_field types/%s", c), func(t *testing.T) {
			var (
				req = require.New(t)
			)

			truncateStore(ctx, s, t)
			err = collect(
				err,
				storeRole(ctx, s, 1, "everyone"),
				storeRole(ctx, s, 2, "admins"),
			)
			if err != nil {
				t.Fatal(err.Error())
			}

			nn, err := decodeDirectory(ctx, path.Join("data_shaping", c))
			req.NoError(err)

			crs := resource.ComposeRecordShaper()
			nn, err = resource.Shape(nn, crs)
			req.NoError(err)

			req.NoError(encode(ctx, s, nn))

			ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
			req.NotNil(ns)
			ms, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod1")
			req.NotNil(ms)

			rr, _, err := store.SearchComposeRecords(ctx, s, ms, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr, 2)

			r1 := rr[0]
			r2 := rr[1]

			req.Len(r1.Values, 7)
			req.Equal("1", r1.Values.Get("f_bool", 0).Value)
			req.Equal("2021-04-01T10:00:00Z", r1.Values.Get("f_datetime", 0).Value)
			req.Equal("mail@test.tld", r1.Values.Get("f_email", 0).Value)
			req.Equal("1.23", r1.Values.Get("f_number", 0).Value)
			req.Equal("opt_2", r1.Values.Get("f_select", 0).Value)
			req.Equal("test here", r1.Values.Get("f_string", 0).Value)
			req.Equal("https://test.tld", r1.Values.Get("f_url", 0).Value)

			req.Len(r2.Values, 7)
			req.Equal("", r2.Values.Get("f_bool", 0).Value)
			req.Equal("2021-04-02T10:00:00Z", r2.Values.Get("f_datetime", 0).Value)
			req.Equal("mail@test.tld", r2.Values.Get("f_email", 0).Value)
			req.Equal("20", r2.Values.Get("f_number", 0).Value)
			req.Equal("opt_3", r2.Values.Get("f_select", 0).Value)
			req.Equal("test here", r2.Values.Get("f_string", 0).Value)
			req.Equal("https://test.tld", r2.Values.Get("f_url", 0).Value)

			s.TruncateComposeRecords(ctx, ms)
		})
	}
}

func TestDataShaping_refs(t *testing.T) {
	var (
		ctx = auth.SetSuperUserContext(context.Background())
		s   = initStore(ctx, t)
		err error

		cases = []string{
			"csv_refs",
		}
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("record shaping; data_shaping_field types/%s", c), func(t *testing.T) {
			var (
				req = require.New(t)
			)

			truncateStore(ctx, s, t)
			err = collect(
				err,
				storeRole(ctx, s, 1, "everyone"),
				storeRole(ctx, s, 2, "admins"),

				storeUser(ctx, s, 201, "u1"),
				storeUser(ctx, s, 202, "u2"),
			)
			if err != nil {
				t.Fatal(err.Error())
			}

			nn, err := decodeDirectory(ctx, path.Join("data_shaping", c))
			req.NoError(err)

			crs := resource.ComposeRecordShaper()
			nn, err = resource.Shape(nn, crs)
			req.NoError(err)

			req.NoError(encode(ctx, s, nn))

			ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "ns1")
			req.NotNil(ns)
			ms, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod1")
			req.NotNil(ms)

			rr, _, err := store.SearchComposeRecords(ctx, s, ms, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr, 4)

			r1 := rr[0]
			r2 := rr[1]
			r3 := rr[2]
			r4 := rr[3]

			req.Len(r1.Values, 2)
			req.Equal(strconv.FormatUint(r1.ID, 10), r1.Values.Get("f_record", 0).Value)
			req.Equal(r1.ID, r1.Values.Get("f_record", 0).Ref)
			req.Equal("201", r1.Values.Get("f_user", 0).Value)
			req.Equal(uint64(201), r1.Values.Get("f_user", 0).Ref)

			req.Len(r2.Values, 2)
			req.Equal(strconv.FormatUint(r2.ID, 10), r2.Values.Get("f_record", 0).Value)
			req.Equal(r2.ID, r2.Values.Get("f_record", 0).Ref)
			req.Equal("202", r2.Values.Get("f_user", 0).Value)
			req.Equal(uint64(202), r2.Values.Get("f_user", 0).Ref)

			req.Len(r3.Values, 2)
			req.Equal(strconv.FormatUint(r4.ID, 10), r3.Values.Get("f_record", 0).Value)
			req.Equal(r4.ID, r3.Values.Get("f_record", 0).Ref)
			req.Equal("201", r3.Values.Get("f_user", 0).Value)
			req.Equal(uint64(201), r3.Values.Get("f_user", 0).Ref)

			req.Len(r4.Values, 2)
			req.Equal(strconv.FormatUint(r3.ID, 10), r4.Values.Get("f_record", 0).Value)
			req.Equal(r3.ID, r4.Values.Get("f_record", 0).Ref)
			req.Equal("202", r4.Values.Get("f_user", 0).Value)
			req.Equal(uint64(202), r4.Values.Get("f_user", 0).Ref)

			s.TruncateComposeRecords(ctx, ms)
		})
	}
}
