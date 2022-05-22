package dal

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
)

func drain(ctx context.Context, i dal.Iterator) (rr []*types.Record, err error) {
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
