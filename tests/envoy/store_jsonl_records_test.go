package envoy

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"strconv"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza-server/pkg/envoy/json"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/stretchr/testify/require"
)

// TestStoreJsonl_records takes data from s1, encodes it into jsonl files, decodes
// created jsonl files, encodes into s2 and compares the data from s2.
func TestStoreJsonl_records(t *testing.T) {
	type (
		tc struct {
			name string
			// Before the data gets processed
			pre func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter)
			// After the data gets processed
			postStoreDecode func(req *require.Assertions, err error)
			postJsonlEncode func(req *require.Assertions, err error)
			postStoreEncode func(req *require.Assertions, err error)
			// Data assertions
			check func(ctx context.Context, s store.Storer, req *require.Assertions)
		}
	)

	ctx := context.Background()
	s := initServices(ctx, t)
	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	cases := []*tc{
		{
			name: "base record",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				truncateStore(ctx, s, t)
				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")
				usr := sTestUser(ctx, t, s, "base")
				sTestComposeRecord(ctx, t, s, ns.ID, mod.ID, usr.ID)

				df := su.NewDecodeFilter().
					ComposeRecord(&types.RecordFilter{
						NamespaceID: ns.ID,
						ModuleID:    mod.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "base_module")
				req.NoError(err)
				usr, err := store.LookupUserByHandle(ctx, s, "base_user")
				req.NoError(err)

				rr, _, err := store.SearchComposeRecords(ctx, s, mod, types.RecordFilter{
					ModuleID:    mod.ID,
					NamespaceID: ns.ID,
				})
				req.NoError(err)
				req.Len(rr, 1)
				rec := rr[0]

				req.Equal(ns.ID, rec.NamespaceID)
				req.Equal(mod.ID, rec.ModuleID)

				req.Equal(createdAt.Format(time.RFC3339), rec.CreatedAt.Format(time.RFC3339))
				req.Equal(updatedAt.Format(time.RFC3339), rec.UpdatedAt.Format(time.RFC3339))
				req.Equal(usr.ID, rec.OwnedBy)
				req.Equal(usr.ID, rec.CreatedBy)
				req.Equal(usr.ID, rec.UpdatedBy)

				req.Len(rec.Values, 2)
				vv := rec.Values.FilterByName("module_field_string")
				req.Len(vv, 1)
				req.Equal("string value", vv[0].Value)

				vv = rec.Values.FilterByName("module_field_number")
				req.Len(vv, 1)
				req.Equal("10", vv[0].Value)
			},
		},

		{
			name: "paged",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				truncateStore(ctx, s, t)

				ns := sTestComposeNamespace(ctx, t, s, "base")
				mod := sTestComposeModule(ctx, t, s, ns.ID, "base")
				usr := sTestUser(ctx, t, s, "base")
				sTestComposeRecordRaw(ctx, t, s, ns.ID, mod.ID, usr.ID, &types.RecordValue{Name: "module_field_number", Value: "10"}, &types.RecordValue{Name: "module_field_string", Value: "v"})
				sTestComposeRecordRaw(ctx, t, s, ns.ID, mod.ID, usr.ID, &types.RecordValue{Name: "module_field_number", Value: "11"}, &types.RecordValue{Name: "module_field_string", Value: "v"})
				sTestComposeRecordRaw(ctx, t, s, ns.ID, mod.ID, usr.ID, &types.RecordValue{Name: "module_field_number", Value: "12"}, &types.RecordValue{Name: "module_field_string", Value: "v"})
				sTestComposeRecordRaw(ctx, t, s, ns.ID, mod.ID, usr.ID, &types.RecordValue{Name: "module_field_number", Value: "13"}, &types.RecordValue{Name: "module_field_string", Value: "v"})
				sTestComposeRecordRaw(ctx, t, s, ns.ID, mod.ID, usr.ID, &types.RecordValue{Name: "module_field_number", Value: "14"}, &types.RecordValue{Name: "module_field_string", Value: "v"})

				df := su.NewDecodeFilter().
					ComposeRecord(&types.RecordFilter{
						NamespaceID: ns.ID,
						ModuleID:    mod.ID,
						Paging: filter.Paging{
							Limit: 2,
						},
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "base_module")
				req.NoError(err)
				usr, err := store.LookupUserByHandle(ctx, s, "base_user")
				req.NoError(err)
				_ = usr

				rr, _, err := store.SearchComposeRecords(ctx, s, mod, types.RecordFilter{
					ModuleID:    mod.ID,
					NamespaceID: ns.ID,
				})
				req.NoError(err)
				req.Len(rr, 5)

				for i, rec := range rr {
					req.Equal(ns.ID, rec.NamespaceID)
					req.Equal(mod.ID, rec.ModuleID)

					req.Equal(createdAt.Format(time.RFC3339), rec.CreatedAt.Format(time.RFC3339))
					req.Equal(updatedAt.Format(time.RFC3339), rec.UpdatedAt.Format(time.RFC3339))
					req.Equal(usr.ID, rec.OwnedBy)
					req.Equal(usr.ID, rec.CreatedBy)
					req.Equal(usr.ID, rec.UpdatedBy)

					vv := rec.Values.FilterByName("module_field_number")
					req.Len(vv, 1)
					req.Equal(fmt.Sprintf("1%d", i), vv[0].Value)
				}
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := require.New(t)

			err, df := c.pre(ctx, s)
			if err != nil {
				t.Fatal(err.Error())
			}
			// Decode from store
			sd := su.Decoder()
			nn, err := sd.Decode(ctx, s, df)
			if c.postStoreDecode != nil {
				c.postStoreDecode(req, err)
			} else {
				req.NoError(err)
			}

			// Encode into jsonl
			je := json.NewBulkRecordEncoder(&json.EncoderConfig{})
			bld := envoy.NewBuilder(je)
			g, err := bld.Build(ctx, nn...)
			req.NoError(err)
			err = envoy.Encode(ctx, g, je)
			ss := je.Stream()
			if c.postJsonlEncode != nil {
				c.postJsonlEncode(req, err)
			} else {
				req.NoError(err)
			}

			// Cleanup the store
			truncateStoreRecords(ctx, s, t)

			// Encode back into store
			se := su.NewStoreEncoder(s, dal.Service(), &su.EncoderConfig{})
			jd := json.Decoder()
			nn = make([]resource.Interface, 0, len(nn))
			for _, s := range ss {
				mm, err := jd.Decode(ctx, s.Source, &envoy.DecoderOpts{
					Name: "tmp.jsonl",
					Path: "/tmp.jsonl",
				})
				req.NoError(err)
				nn = append(nn, mm...)
			}

			tpl := resource.NewComposeRecordTemplate(
				"base_module",
				"base_namespace",
				"tmp.jsonl",
				true,
				resource.MappingTplSet{
					{
						Cell:  "id",
						Field: "/",
					},
				},
			)

			nn = append(nn, tpl)
			crs := resource.ComposeRecordShaper()
			nn, err = resource.Shape(nn, crs)
			req.NoError(err)
			bld = envoy.NewBuilder(se)
			g, err = bld.Build(ctx, nn...)
			req.NoError(err)

			err = envoy.Encode(ctx, g, se)
			if c.postStoreEncode != nil {
				c.postStoreEncode(req, err)
			} else {
				req.NoError(err)
			}

			// Assert
			c.check(ctx, s, req)

			// Cleanup the store
			truncateStoreRecords(ctx, s, t)
		})
		ni = 0
	}
}

func TestStoreJsonl_records_fieldTypes(t *testing.T) {
	type (
		tc struct {
			name string
			// Before the data gets processed
			pre func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter)
			// After the data gets processed
			postStoreDecode func(req *require.Assertions, err error)
			postJsonlEncode func(req *require.Assertions, err error)
			postStoreEncode func(req *require.Assertions, err error)
			// Data assertions
			check func(ctx context.Context, s store.Storer, req *require.Assertions)
		}
	)

	ctx := context.Background()
	s := initServices(ctx, t)
	ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

	ni := uint64(10)
	su.NextID = func() uint64 {
		ni++
		return ni
	}

	cases := []*tc{
		{
			name: "base field types",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				truncateStore(ctx, s, t)
				ns := sTestComposeNamespace(ctx, t, s, "base")
				usr := sTestUser(ctx, t, s, "base")
				mod := sTestComposeModuleFull(ctx, s, t, ns.ID, "base")

				recID := su.NextID()
				rec := &types.Record{
					ID:          recID,
					NamespaceID: ns.ID,
					ModuleID:    mod.ID,

					Values: types.RecordValueSet{
						{
							RecordID: recID,
							Name:     "BoolTrue",
							Value:    "1",
						},
						{
							RecordID: recID,
							Name:     "BoolFalse",
							Value:    "0",
						},
						{
							RecordID: recID,
							Name:     "DateTime",
							Value:    "2021-01-01T11:10:09Z",
						},
						{
							RecordID: recID,
							Name:     "Email",
							Value:    "test@mail.tld",
						},
						{
							RecordID: recID,
							Name:     "Select",
							Value:    "v1",
						},
						{
							RecordID: recID,
							Name:     "Number",
							Value:    "10.01",
						},
						{
							RecordID: recID,
							Name:     "String",
							Value:    "testing",
						},
						{
							RecordID: recID,
							Name:     "Url",
							Value:    "htts://www.testing.tld",
						},
						{
							RecordID: recID,
							Name:     "User",
							Value:    strconv.FormatUint(usr.ID, 10),
							Ref:      usr.ID,
						},
					},
				}
				err := store.CreateComposeRecord(ctx, s, mod, rec)
				if err != nil {
					t.Fatal(err)
				}

				df := su.NewDecodeFilter().
					ComposeRecord(&types.RecordFilter{
						NamespaceID: ns.ID,
						ModuleID:    mod.ID,
					})
				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				ns, err := store.LookupComposeNamespaceBySlug(ctx, s, "base_namespace")
				req.NoError(err)
				mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, ns.ID, "base_module")
				req.NoError(err)
				usr, err := store.LookupUserByHandle(ctx, s, "base_user")
				req.NoError(err)

				rr, _, err := store.SearchComposeRecords(ctx, s, mod, types.RecordFilter{
					ModuleID:    mod.ID,
					NamespaceID: ns.ID,
				})
				req.NoError(err)
				req.Len(rr, 1)
				rec := rr[0]

				req.Equal("1", rec.Values.FilterByName("BoolTrue")[0].Value)
				req.Equal("", rec.Values.FilterByName("BoolFalse")[0].Value)
				req.Equal("2021-01-01T11:10:09Z", rec.Values.FilterByName("DateTime")[0].Value)
				req.Equal("test@mail.tld", rec.Values.FilterByName("Email")[0].Value)
				req.Equal("v1", rec.Values.FilterByName("Select")[0].Value)
				req.Equal("10.01", rec.Values.FilterByName("Number")[0].Value)
				req.Equal("testing", rec.Values.FilterByName("String")[0].Value)
				req.Equal("htts://www.testing.tld", rec.Values.FilterByName("Url")[0].Value)
				req.Equal(strconv.FormatUint(usr.ID, 10), rec.Values.FilterByName("User")[0].Value)
				req.Equal(usr.ID, rec.Values.FilterByName("User")[0].Ref)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := require.New(t)

			err, df := c.pre(ctx, s)
			if err != nil {
				t.Fatal(err.Error())
			}
			// Decode from store
			sd := su.Decoder()
			nn, err := sd.Decode(ctx, s, df)
			if c.postStoreDecode != nil {
				c.postStoreDecode(req, err)
			} else {
				req.NoError(err)
			}

			// Encode into CSV
			ce := csv.NewBulkRecordEncoder(&csv.EncoderConfig{})
			bld := envoy.NewBuilder(ce)
			g, err := bld.Build(ctx, nn...)
			req.NoError(err)
			err = envoy.Encode(ctx, g, ce)
			ss := ce.Stream()
			if c.postJsonlEncode != nil {
				c.postJsonlEncode(req, err)
			} else {
				req.NoError(err)
			}

			// Cleanup the store
			truncateStoreRecords(ctx, s, t)

			// Encode back into store
			se := su.NewStoreEncoder(s, dal.Service(), &su.EncoderConfig{})
			yd := csv.Decoder()
			nn = make([]resource.Interface, 0, len(nn))
			for _, s := range ss {
				mm, err := yd.Decode(ctx, s.Source, &envoy.DecoderOpts{
					Name: "tmp.csv",
					Path: "/tmp.csv",
				})
				req.NoError(err)
				nn = append(nn, mm...)
			}

			tpl := resource.NewComposeRecordTemplate(
				"base_module",
				"base_namespace",
				"tmp.csv",
				true,
				resource.MappingTplSet{
					{
						Cell:  "id",
						Field: "/",
					},
				},
			)

			nn = append(nn, tpl)
			crs := resource.ComposeRecordShaper()
			nn, err = resource.Shape(nn, crs)
			req.NoError(err)
			bld = envoy.NewBuilder(se)
			g, err = bld.Build(ctx, nn...)
			req.NoError(err)

			err = envoy.Encode(ctx, g, se)
			if c.postStoreEncode != nil {
				c.postStoreEncode(req, err)
			} else {
				req.NoError(err)
			}

			// Assert
			c.check(ctx, s, req)

			// Cleanup the store
			truncateStoreRecords(ctx, s, t)
		})
		ni = 0
	}
}
