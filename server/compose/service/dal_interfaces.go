package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	dalDater interface {
		Create(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, vv ...dal.ValueGetter) error
		Update(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, rr ...dal.ValueGetter) (err error)
		Search(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (dal.Iterator, error)
		Count(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (uint, error)
		Run(ctx context.Context, pp dal.Pipeline) (dal.Iterator, error)
		Lookup(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
		Delete(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, pkv ...dal.ValueGetter) (err error)
		Truncate(ctx context.Context, m dal.ModelRef, operations dal.OperationSet) (err error)
	}
)
