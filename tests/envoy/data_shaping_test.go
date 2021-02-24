package envoy

import (
	"context"
	"fmt"
	"path"
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
