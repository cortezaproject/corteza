package dalutils

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/dal/capabilities"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	creator interface {
		Create(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, vv ...dal.ValueGetter) error
	}

	updater interface {
		Update(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, rr ...dal.ValueGetter) (err error)
	}

	searcher interface {
		Search(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, f filter.Filter) (dal.Iterator, error)
	}

	lookuper interface {
		Lookup(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
	}

	deleter interface {
		Delete(ctx context.Context, m dal.ModelFilter, capabilities capabilities.Set, pkv ...dal.ValueGetter) (err error)
	}
)

func ComposeRecordsList(ctx context.Context, s searcher, mod *types.Module, filter types.RecordFilter) (set types.RecordSet, outFilter types.RecordFilter, err error) {
	iter, err := prepIterator(ctx, s, mod, filter)
	if err != nil {
		return
	}

	set, outFilter, err = drainIterator(ctx, iter, mod, filter)
	return
}

func ComposeRecordsIterator(ctx context.Context, s searcher, mod *types.Module, filter types.RecordFilter) (iter dal.Iterator, outFilter types.RecordFilter, err error) {
	iter, err = prepIterator(ctx, s, mod, filter)
	if err != nil {
		return
	}

	outFilter = filter
	outFilter.Paging = *filter.Paging.Clone()

	return
}

func ComposeRecordsFind(ctx context.Context, l lookuper, mod *types.Module, recordID uint64) (out *types.Record, err error) {
	out = prepareRecordTarget(mod)

	err = l.Lookup(ctx, mod.ModelFilter(), recLookupCapabilities(mod), dal.PKValues{"id": recordID}, out)
	if err != nil {
		return
	}

	return
}

func ComposeRecordCreate(ctx context.Context, c creator, mod *types.Module, records ...*types.Record) (err error) {
	return c.Create(ctx, mod.ModelFilter(), recCreateCapabilities(mod), recToGetters(records...)...)
}

func ComposeRecordUpdate(ctx context.Context, u updater, mod *types.Module, records ...*types.Record) (err error) {
	return u.Update(ctx, mod.ModelFilter(), recUpdateCapabilities(mod), recToGetters(records...)...)
}

func ComposeRecordSoftDelete(ctx context.Context, u updater, invoker uint64, mod *types.Module, records ...*types.Record) (err error) {
	n := time.Now()
	for _, r := range records {
		r.DeletedAt = &n
		r.DeletedBy = invoker
	}

	return u.Update(ctx, mod.ModelFilter(), recUpdateCapabilities(mod), recToGetters(records...)...)
}

func ComposeRecordDelete(ctx context.Context, d deleter, mod *types.Module, records ...*types.Record) (err error) {
	return d.Delete(ctx, mod.ModelFilter(), recDeleteCapabilities(mod), recToGetters(records...)...)
}

func WalkIterator(ctx context.Context, iter dal.Iterator, mod *types.Module, f func(r *types.Record) error) (err error) {
	for iter.Next(ctx) {
		r := prepareRecordTarget(mod)
		if err = iter.Scan(r); err != nil {
			return
		}

		if err = f(r); err != nil {
			return
		}
	}

	return iter.Err()
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utils

func prepFilter(filter types.RecordFilter, mod *types.Module) (dalFilter filter.Filter) {
	dalFilter = filter.ToFilter()
	if mod.ModelConfig.Partitioned {
		dalFilter = filter.ToConstraintedFilter(mod.ModelConfig.Constraints)
	}

	return
}

func prepIterator(ctx context.Context, dal searcher, mod *types.Module, filter types.RecordFilter) (iter dal.Iterator, err error) {
	dalFilter := prepFilter(filter, mod)

	iter, err = dal.Search(ctx, mod.ModelFilter(), recSearchCapabilities(mod, filter), dalFilter)
	return
}

func drainIterator(ctx context.Context, iter dal.Iterator, mod *types.Module, filter types.RecordFilter) (set types.RecordSet, outFilter types.RecordFilter, err error) {
	defer iter.Close()

	// Get the requested number of recrds
	set = make(types.RecordSet, 0, filter.Limit)
	i := 0
	err = WalkIterator(ctx, iter, mod, func(r *types.Record) error {
		set = append(set, r)
		i++
		return nil
	})
	if err != nil {
		return
	}

	// Make out filter
	outFilter = filter
	pp := filter.Paging.Clone()

	if len(set) > 0 && filter.PrevPage != nil {
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

func prepareRecordTarget(module *types.Module) *types.Record {
	// so we can avoid some code later involving (non)partitioned modules :seenoevil:
	return &types.Record{
		ModuleID:    module.ID,
		NamespaceID: module.NamespaceID,
		Values:      make(types.RecordValueSet, 0, len(module.Fields)),
	}
}

func recToGetters(rr ...*types.Record) (out []dal.ValueGetter) {
	out = make([]dal.ValueGetter, len(rr))

	for i := range rr {
		out[i] = rr[i]
	}

	return
}

func recToGetter(rr ...*types.Record) (out dal.ValueGetter) {
	if len(rr) == 0 {
		return
	}

	return recToGetters(rr...)[0]
}
