package envoy

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func TestStoreYaml_APIGateway(t *testing.T) {
	type (
		tc struct {
			name string
			// Before the data gets processed
			pre func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter)
			// After the data gets processed
			postStoreDecode func(req *require.Assertions, err error)
			postYamlEncode  func(req *require.Assertions, err error)
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
			name: "base",
			pre: func(ctx context.Context, s store.Storer) (error, *su.DecodeFilter) {
				gwr := sTestAPIGatewayRoute(ctx, t, s, "test")
				_ = sTestAPIGatewayFilter(ctx, t, s, gwr.ID, "test")

				df := su.NewDecodeFilter().
					APIGWRoutes(&types.ApigwRouteFilter{})

				return nil, df
			},
			check: func(ctx context.Context, s store.Storer, req *require.Assertions) {
				rr, _, err := store.SearchApigwRoutes(ctx, s, types.ApigwRouteFilter{})
				req.NoError(err)
				req.Len(rr, 1)

				r := rr[0]
				req.Equal("/testing/test", r.Endpoint)
				req.Equal("POST", r.Method)
				req.True(r.Enabled)
				req.True(r.Meta.Debug)
				req.True(r.Meta.Async)

				ff, _, err := store.SearchApigwFilters(ctx, s, types.ApigwFilterFilter{RouteID: rr[0].ID})
				req.NoError(err)
				req.Len(ff, 1)

				f := ff[0]
				req.Equal(r.ID, f.Route)
				req.Equal("test_ref", f.Ref)
				req.Equal("test_kind", f.Kind)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := require.New(t)

			truncateStore(ctx, s, t)
			err, df := c.pre(ctx, s)
			if err != nil {
				t.Fatal(err.Error())
			}
			// Decode from store
			sd := su.Decoder()
			nn, err := sd.Decode(ctx, s, dal.Service(), df)
			if c.postStoreDecode != nil {
				c.postStoreDecode(req, err)
			} else {
				req.NoError(err)
			}

			// Encode into YAML
			ye := yaml.NewYamlEncoder(&yaml.EncoderConfig{})
			bld := envoy.NewBuilder(ye)
			g, err := bld.Build(ctx, nn...)
			req.NoError(err)
			err = envoy.Encode(ctx, g, ye)
			ss := ye.Stream()
			if c.postYamlEncode != nil {
				c.postYamlEncode(req, err)
			} else {
				req.NoError(err)
			}

			// Cleanup the store
			truncateStore(ctx, s, t)

			// Encode back into store
			se := su.NewStoreEncoder(s, dal.Service(), &su.EncoderConfig{})
			yd := yaml.Decoder()
			nn = make([]resource.Interface, 0, len(nn))
			for _, s := range ss {
				mm, err := yd.Decode(ctx, s.Source, nil)
				req.NoError(err)
				nn = append(nn, mm...)
			}
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
			truncateStore(ctx, s, t)
		})
		ni = 0
		truncateStore(ctx, s, t)
	}
}
