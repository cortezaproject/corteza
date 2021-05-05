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

func TestDataShaping_large(t *testing.T) {
	var (
		ctx = auth.SetSuperUserContext(context.Background())
		s   = initStore(ctx, t)
		err error

		cases = []string{
			"csv_large",
			"jsonl_large",
		}
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("record shaping; data_shaping; large datasets/%s", c), func(t *testing.T) {
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
			req.Len(rr, 2000)

			for i, r := range rr {
				req.Equal(fmt.Sprintf("r%df1", i+1), r.Values.Get("f1", 0).Value)
				req.Equal(fmt.Sprintf("r%df2", i+1), r.Values.Get("f2", 0).Value)
				req.Equal(fmt.Sprintf("r%df3", i+1), r.Values.Get("f3", 0).Value)
				req.Equal(fmt.Sprintf("r%df4", i+1), r.Values.Get("f4", 0).Value)
				req.Equal(fmt.Sprintf("r%df5", i+1), r.Values.Get("f5", 0).Value)
				req.Equal(fmt.Sprintf("r%df6", i+1), r.Values.Get("f6", 0).Value)
			}

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
		t.Run(fmt.Sprintf("record shaping; data_shaping_field_types/%s", c), func(t *testing.T) {
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
		t.Run(fmt.Sprintf("record shaping; data_shaping_refs/%s", c), func(t *testing.T) {
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

func TestDataShaping_xrefsPeer(t *testing.T) {
	var (
		ctx = auth.SetSuperUserContext(context.Background())
		s   = initStore(ctx, t)
		err error

		cases = []string{
			"csv_xrefs",
		}
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("record shaping; data_shaping_xrefs_peer/%s", c), func(t *testing.T) {
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

			m, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod1")
			req.NotNil(m)
			refM, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod2")
			req.NotNil(refM)

			rr, _, err := store.SearchComposeRecords(ctx, s, m, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr, 4)
			refRR, _, err := store.SearchComposeRecords(ctx, s, refM, types.RecordFilter{})
			req.NoError(err)
			req.Len(refRR, 4)

			r1 := rr[0]
			r2 := rr[1]
			r3 := rr[2]
			r4 := rr[3]

			refR1 := refRR[0]
			refR2 := refRR[1]
			refR3 := refRR[2]
			refR4 := refRR[3]

			req.Len(r1.Values, 1)
			req.Equal(strconv.FormatUint(refR1.ID, 10), r1.Values.Get("f_record", 0).Value)
			req.Equal(refR1.ID, r1.Values.Get("f_record", 0).Ref)

			req.Len(r2.Values, 1)
			req.Equal(strconv.FormatUint(refR2.ID, 10), r2.Values.Get("f_record", 0).Value)
			req.Equal(refR2.ID, r2.Values.Get("f_record", 0).Ref)

			req.Len(r3.Values, 1)
			req.Equal(strconv.FormatUint(refR3.ID, 10), r3.Values.Get("f_record", 0).Value)
			req.Equal(refR3.ID, r3.Values.Get("f_record", 0).Ref)

			req.Len(r4.Values, 1)
			req.Equal(strconv.FormatUint(refR4.ID, 10), r4.Values.Get("f_record", 0).Value)
			req.Equal(refR4.ID, r4.Values.Get("f_record", 0).Ref)

			s.TruncateComposeRecords(ctx, m)
		})
	}
}

func TestDataShaping_xrefsStore(t *testing.T) {
	var (
		ctx = auth.SetSuperUserContext(context.Background())
		s   = initStore(ctx, t)
		err error

		cases = []string{
			"csv_xrefs_store",
		}
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("record shaping; data_shaping_xrefs_store/%s", c), func(t *testing.T) {
			var (
				req = require.New(t)
			)

			truncateStore(ctx, s, t)
			err = collect(
				err,
				storeRole(ctx, s, 1, "everyone"),
				storeRole(ctx, s, 2, "admins"),

				storeComposeNamespace(ctx, s, 1001, "ns1"),
				storeComposeModule(ctx, s, 1001, 2001, "mod_ref"),
				storeComposeModuleField(ctx, s, 2001, 2101, "label"),

				storeComposeRecord(ctx, s, 1001, 2001, 3001, "label"),
				storeComposeRecord(ctx, s, 1001, 2001, 3002, "label"),
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

			m, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod1")
			req.NotNil(m)
			refM, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod_ref")
			req.NotNil(refM)

			rr, _, err := store.SearchComposeRecords(ctx, s, m, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr, 4)
			refRR, _, err := store.SearchComposeRecords(ctx, s, refM, types.RecordFilter{})
			req.NoError(err)
			req.Len(refRR, 2)

			r1 := rr[0]
			r2 := rr[1]
			r3 := rr[2]
			r4 := rr[3]

			refR1 := refRR[0]
			refR2 := refRR[1]

			req.Len(r1.Values, 1)
			req.Equal(strconv.FormatUint(refR1.ID, 10), r1.Values.Get("f_record", 0).Value)
			req.Equal(refR1.ID, r1.Values.Get("f_record", 0).Ref)

			req.Len(r2.Values, 1)
			req.Equal(strconv.FormatUint(refR2.ID, 10), r2.Values.Get("f_record", 0).Value)
			req.Equal(refR2.ID, r2.Values.Get("f_record", 0).Ref)

			req.Len(r3.Values, 1)
			req.Equal(strconv.FormatUint(refR1.ID, 10), r3.Values.Get("f_record", 0).Value)
			req.Equal(refR1.ID, r3.Values.Get("f_record", 0).Ref)

			req.Len(r4.Values, 1)
			req.Equal(strconv.FormatUint(refR2.ID, 10), r4.Values.Get("f_record", 0).Value)
			req.Equal(refR2.ID, r4.Values.Get("f_record", 0).Ref)

			s.TruncateComposeRecords(ctx, m)
		})
	}
}

func TestDataShaping_xrefsMix(t *testing.T) {
	var (
		ctx = auth.SetSuperUserContext(context.Background())
		s   = initStore(ctx, t)
		err error

		cases = []string{
			"csv_xrefs_mix",
		}
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("record shaping; data_shaping_xrefs_mix/%s", c), func(t *testing.T) {
			var (
				req = require.New(t)
			)

			truncateStore(ctx, s, t)
			err = collect(
				err,
				storeRole(ctx, s, 1, "everyone"),
				storeRole(ctx, s, 2, "admins"),

				storeComposeNamespace(ctx, s, 1001, "ns1"),

				storeComposeModule(ctx, s, 1001, 2001, "mod1"),
				storeComposeModuleField(ctx, s, 2001, 2101, "f_label"),
				storeComposeRecord(ctx, s, 1001, 2001, 3001, "f_label"),

				storeComposeModule(ctx, s, 1001, 2002, "mod2"),
				storeComposeModuleField(ctx, s, 2002, 2201, "f_label"),
				storeComposeRecord(ctx, s, 1001, 2002, 3101, "f_label"),
				storeComposeRecord(ctx, s, 1001, 2002, 3102, "f_label"),
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

			mod1, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod1")
			req.NotNil(mod1)
			mod2, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod2")
			req.NotNil(mod2)

			rr1, _, err := store.SearchComposeRecords(ctx, s, mod1, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr1, 5)
			rr2, _, err := store.SearchComposeRecords(ctx, s, mod2, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr2, 4)

			r1 := rr1[0]
			r2 := rr1[1]
			r3 := rr1[2]
			r4 := rr1[3]
			refStoreSelf := rr1[4]

			refStoreR1 := rr2[2]
			refStoreR2 := rr2[3]
			refCSVR1 := rr2[0]

			req.Len(r1.Values, 2)
			req.Equal(strconv.FormatUint(refCSVR1.ID, 10), r1.Values.Get("f_record", 0).Value)
			req.Equal(refCSVR1.ID, r1.Values.Get("f_record", 0).Ref)

			req.Len(r2.Values, 2)
			req.Equal(strconv.FormatUint(refStoreSelf.ID, 10), r2.Values.Get("f_record_self", 0).Value)
			req.Equal(refStoreSelf.ID, r2.Values.Get("f_record_self", 0).Ref)

			req.Len(r3.Values, 2)
			req.Equal(strconv.FormatUint(refStoreR1.ID, 10), r3.Values.Get("f_record", 0).Value)
			req.Equal(refStoreR1.ID, r3.Values.Get("f_record", 0).Ref)

			req.Len(r4.Values, 2)
			req.Equal(strconv.FormatUint(refStoreR2.ID, 10), r4.Values.Get("f_record", 0).Value)
			req.Equal(refStoreR2.ID, r4.Values.Get("f_record", 0).Ref)

			s.TruncateComposeRecords(ctx, mod1)
		})
	}
}

func TestDataShaping_update(t *testing.T) {
	var (
		ctx = auth.SetSuperUserContext(context.Background())
		s   = initStore(ctx, t)
		err error

		cases = []string{
			"csv_update",
		}
	)

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("record shaping; data_shaping_update/%s", c), func(t *testing.T) {
			var (
				req = require.New(t)
			)

			truncateStore(ctx, s, t)
			err = collect(
				err,
				storeRole(ctx, s, 1, "everyone"),
				storeRole(ctx, s, 2, "admins"),

				storeComposeNamespace(ctx, s, 1001, "ns1"),
				storeComposeModule(ctx, s, 1001, 2001, "mod1"),
				storeComposeModuleField(ctx, s, 2001, 2101, "f_label"),

				storeComposeRecord(ctx, s, 1001, 2001, 3001, "f_label"),
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

			mod1, err := loadComposeModuleFull(ctx, s, req, ns.ID, "mod1")
			req.NotNil(mod1)

			rr, _, err := store.SearchComposeRecords(ctx, s, mod1, types.RecordFilter{})
			req.NoError(err)
			req.Len(rr, 2)

			r1 := rr[0]
			r2 := rr[1]

			req.Len(r1.Values, 1)
			req.Equal("created", r1.Values.Get("f_label", 0).Value)

			req.Len(r2.Values, 1)
			req.Equal("updated", r2.Values.Get("f_label", 0).Value)

			s.TruncateComposeRecords(ctx, mod1)
		})
	}
}
