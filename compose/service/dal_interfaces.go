package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	dalModeler interface {
		ModelIdentFormatter(connectionID uint64) (f *dal.IdentFormatter, err error)

		ReloadModel(ctx context.Context, models ...*dal.Model) (err error)
		CreateModel(ctx context.Context, models ...*dal.Model) (err error)
		UpdateModel(ctx context.Context, old, new *dal.Model) (err error)
		UpdateModelAttribute(ctx context.Context, model *dal.Model, old, new *dal.Attribute, trans ...dal.TransformationFunction) (err error)
		DeleteModel(ctx context.Context, models ...*dal.Model) (err error)

		SearchModelIssues(connectionID, resourceID uint64) (out []error)
	}

	dalDater interface {
		Create(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, vv ...dal.ValueGetter) error
		Update(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, rr ...dal.ValueGetter) (err error)
		Search(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, f filter.Filter) (dal.Iterator, error)
		Lookup(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
		Delete(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, pkv ...dal.ValueGetter) (err error)
		Truncate(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set) (err error)
	}

	dalService interface {
		dalModeler
		dalDater
	}
)
