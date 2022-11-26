package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	dalModeler interface {
		SearchModels(ctx context.Context) (out dal.ModelSet, err error)
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
		RemoveModel(ctx context.Context, connectionID, ID uint64) (err error)
		ReplaceModelAttribute(ctx context.Context, model *dal.Model, diff *dal.ModelDiff, hasRecords bool, trans ...dal.TransformationFunction) (err error)

		GetConnectionByID(uint64) *dal.ConnectionWrap

		SearchModelIssues(resourceID uint64) (out []error)
	}

	dalDater interface {
		Create(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, vv ...dal.ValueGetter) error
		Update(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, rr ...dal.ValueGetter) (err error)
		Search(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (dal.Iterator, error)
		Run(ctx context.Context, pp dal.Pipeline) (dal.Iterator, error)
		Lookup(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
		Delete(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, pkv ...dal.ValueGetter) (err error)
		Truncate(ctx context.Context, m dal.ModelRef, operations dal.OperationSet) (err error)
	}

	dalService interface {
		dalModeler
		dalDater
	}
)
