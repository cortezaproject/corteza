package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	dalModeler interface {
		SearchModels(ctx context.Context) (out dal.ModelSet, err error)
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
		RemoveModel(ctx context.Context, connectionID, ID uint64) (err error)
		ReplaceModelAttribute(ctx context.Context, model *dal.Model, old, new *dal.Attribute, trans ...dal.TransformationFunction) (err error)

		GetConnectionMeta(ctx context.Context, ID uint64) (cm dal.ConnectionConfig, err error)

		SearchModelIssues(connectionID, resourceID uint64) (out []error)
	}

	dalDater interface {
		Create(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, vv ...dal.ValueGetter) error
		Update(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, rr ...dal.ValueGetter) (err error)
		Search(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (dal.Iterator, error)
		Lookup(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
		Delete(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, pkv ...dal.ValueGetter) (err error)
		Truncate(ctx context.Context, m dal.ModelRef, operations dal.OperationSet) (err error)
	}

	dalService interface {
		dalModeler
		dalDater
	}
)
