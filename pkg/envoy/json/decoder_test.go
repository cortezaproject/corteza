package json

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/stretchr/testify/require"
)

func parseDocument(ctx context.Context, name string) (*resource.ResourceDataset, error) {
	path := fmt.Sprintf("testdata/%s.jsonl", name)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cd := Decoder()
	ii, err := cd.Decode(ctx, f, &envoy.DecoderOpts{
		Name: name,
		Path: path,
	})
	if err != nil {
		return nil, err
	}
	rd, _ := ii[0].(*resource.ResourceDataset)
	return rd, nil
}

func TestDecoder(t *testing.T) {
	ctx := context.Background()

	t.Run("doc 1", func(t *testing.T) {
		req := require.New(t)

		ds, err := parseDocument(ctx, "records_1")
		req.NoError(err)
		req.NotNil(ds)
		req.Equal(resource.MakeIdentifiers("records_1"), ds.Identifiers())
		req.Subset([]string{"id", "c1", "c2", "c3"}, ds.P.Fields())
		req.Equal(uint64(3), ds.P.Count())

		for i := 0; i < 3; i++ {
			n, err := ds.P.Next()
			req.NoError(err)
			req.NotNil(n)

			req.Len(n, 4)
		}
		n, err := ds.P.Next()
		req.Nil(n)
		req.Nil(err)
	})
}
