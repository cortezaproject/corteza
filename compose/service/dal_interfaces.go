package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	dalModeler interface {
		SearchModels(ctx context.Context) (out dal.ModelSet, err error)
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
		RemoveModel(ctx context.Context, connectionID, ID uint64) (err error)
		ReplaceModelAttribute(ctx context.Context, model *dal.Model, old, new *dal.Attribute, trans ...dal.TransformationFunction) (err error)

		GetConnectionMeta(ctx context.Context, ID uint64) (cm dal.ConnectionMeta, err error)

		SearchModelIssues(connectionID, resourceID uint64) (out []error)
	}

	dalDater interface {
		Create(ctx context.Context, m dal.ModelRef, capabilities capabilities.Set, vv ...dal.ValueGetter) error
		Update(ctx context.Context, m dal.ModelRef, capabilities capabilities.Set, rr ...dal.ValueGetter) (err error)
		Search(ctx context.Context, m dal.ModelRef, capabilities capabilities.Set, f filter.Filter) (dal.Iterator, error)
		Lookup(ctx context.Context, m dal.ModelRef, capabilities capabilities.Set, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
		Delete(ctx context.Context, m dal.ModelRef, capabilities capabilities.Set, pkv ...dal.ValueGetter) (err error)
		Truncate(ctx context.Context, m dal.ModelRef, capabilities capabilities.Set) (err error)
	}

	dalService interface {
		dalModeler
		dalDater
	}
)
