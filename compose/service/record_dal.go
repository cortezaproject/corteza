package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	dalDML interface {
		Create(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, vv ...dal.ValueGetter) error
		Update(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, r dal.ValueGetter) (err error)
		Search(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, f filter.Filter) (dal.Iterator, error)
		Lookup(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
		Delete(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, pkv dal.ValueGetter) (err error)
		Truncate(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set) (err error)
	}
)

func (svc *record) drainIterator(ctx context.Context, iter dal.Iterator, f types.RecordFilter, module *types.Module) (set types.RecordSet, outFilter types.RecordFilter, err error) {
	set = make(types.RecordSet, 0, f.Limit)

	i := 0
	for iter.Next(ctx) {
		auxr := svc.prepareRecordTarget(module)
		if err = iter.Scan(auxr); err != nil {
			return
		}

		set = append(set, auxr)

		i++
	}
	err = iter.Err()

	outFilter = f
	pp := f.Paging.Clone()

	if len(set) > 0 && f.PrevPage != nil {
		pp.PrevPage, err = iter.BackCursor(set[0])
		if err != nil {
			return
		}
	}

	if len(set) > 0 {
		pp.NextPage, err = iter.ForwardCursor(set[len(set)-1])
		if err != nil {
			return
		}
	}

	outFilter.Paging = *pp

	return
}

func (svc *record) prepareRecordTarget(module *types.Module) *types.Record {
	// so we can avoid some code later involving (non)partitioned modules :seenoevil:
	return &types.Record{
		ModuleID:    module.ID,
		NamespaceID: module.NamespaceID,
		Values:      make(types.RecordValueSet, 0, len(module.Fields)),
	}
}

func (svc *record) recToGetters(rr ...*types.Record) (out []dal.ValueGetter) {
	out = make([]dal.ValueGetter, len(rr))

	for i := range rr {
		out[i] = rr[i]
	}

	return
}

func (svc *record) recToGetter(rr ...*types.Record) (out dal.ValueGetter) {
	if len(rr) == 0 {
		return
	}

	return svc.recToGetters(rr...)[0]
}

func (svc *record) recCreateCapabilities(m *types.Module) (out capabilities.Set) {
	return capabilities.CreateCapabilities(m.ModelConfig.Capabilities...)
}

func (svc *record) recUpdateCapabilities(m *types.Module) (out capabilities.Set) {
	return capabilities.UpdateCapabilities(m.ModelConfig.Capabilities...)
}

func (svc *record) recDeleteCapabilities(m *types.Module) (out capabilities.Set) {
	return capabilities.DeleteCapabilities(m.ModelConfig.Capabilities...)
}

func (svc *record) recFilterCapabilities(f types.RecordFilter) (out capabilities.Set) {
	if f.PageCursor != nil {
		out = append(out, capabilities.Paging)
	}

	if f.IncPageNavigation {
		out = append(out, capabilities.Paging)
	}

	if f.IncTotal {
		out = append(out, capabilities.Stats)
	}

	if f.Sort != nil {
		out = append(out, capabilities.Sorting)
	}

	return
}

func (svc *record) recSearchCapabilities(m *types.Module, f types.RecordFilter) (out capabilities.Set) {
	return capabilities.SearchCapabilities(m.ModelConfig.Capabilities...).
		Union(svc.recFilterCapabilities(f))
}
