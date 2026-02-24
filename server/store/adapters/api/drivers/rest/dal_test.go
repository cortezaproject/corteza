package rest

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestDing(t *testing.T) {
	ctx := context.Background()
	dc, err := dalConnector(ctx, "rest://localhost:3399")
	require.NoError(t, err)

	m := &dal.Model{
		Resource: "r",
		ResourceID: 1,
		ConnectionID: 0,
		Ident: "test",
		Attributes: dal.AttributeSet{{
			Ident: "id",
			Label: "id",
			Type: &dal.TypeNumber{},
			Store: &dal.CodecPlain{},
			PrimaryKey: true,
		}, {
			Ident: "name",
			Label: "name",
			Type: &dal.TypeText{},
			Store: &dal.CodecPlain{},
			PrimaryKey: false,
		}, {
			Ident: "status",
			Label: "status",
			Type: &dal.TypeText{},
			Store: &dal.CodecPlain{},
			PrimaryKey: false,
		}},
	}

	err = dc.CreateModel(ctx, m)
	require.NoError(t, err)

	// dc.Create(ctx, m, rawRecord{
	// 	vals: map[string][]any{
	// 		"f1": {"v1"},
	// 	},
	// })
	// require.NoError(t, err)


	// out := &rawRecord{}
	// err = dc.Lookup(ctx, m, rawRecord{
	// 	vals: map[string][]any{
	// 		"id": {"25"},
	// 	},
	// }, out)
	// require.NoError(t, err)

	// err = dc.Update(ctx, m, rawRecord{
	// 	vals: map[string][]any{
	// 		"id": {"20"},
	// 		"name": {"edited"},
	// 		"status": {"back active"},
	// 	},
	// })
	// require.NoError(t, err)

	// err = dc.Delete(ctx, m, rawRecord{
	// 	vals: map[string][]any{
	// 		"id": {"20"},
	// 		"name": {"edited"},
	// 		"status": {"back active"},
	// 	},
	// })
	// require.NoError(t, err)

	iter, err := dc.Search(ctx, m, filter.Generic())
	require.NoError(t, err)

	ok := iter.Next(ctx)
	_ = ok
	_ = iter

	ding := &rawRecord{}
	err = iter.Scan(ding)
	require.NoError(t, err)

	spew.Dump("ding", ding)

	_ = dc
	panic("AAAAAAA")
}


type (
	rawRecord struct {
		vals map[string][]any
	}
)


func (rr rawRecord) CountValues() map[string]uint {
	out := make(map[string]uint)
	for k, v := range rr.vals {
		out[k] = uint(len(v))
	}

	return out
}

func (rr rawRecord) GetValue(k string,i uint) (any, error) {
	return rr.vals[k][i], nil
}

func (rr *rawRecord) SetValue(k string, i uint, v any) error {
	if rr.vals == nil {
		rr.vals = map[string][]any{}
	}

	rr.vals[k] = []any{v}

	return nil
}


