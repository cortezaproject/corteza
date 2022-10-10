package dalutils

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	creator interface {
		Create(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, vv ...dal.ValueGetter) error
	}

	updater interface {
		Update(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, rr ...dal.ValueGetter) (err error)
	}

	searcher interface {
		Search(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (dal.Iterator, error)
	}

	lookuper interface {
		Lookup(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, lookup dal.ValueGetter, dst dal.ValueSetter) (err error)
	}

	deleter interface {
		Delete(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, pkv ...dal.ValueGetter) (err error)
	}
)

// ComposeRecordsList iterates over results and collects all available records
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

	err = l.Lookup(ctx, mod.ModelRef(), recLookupOperations(mod), dal.PKValues{"id": recordID}, out)
	if err != nil {
		return
	}

	return
}

func ComposeRecordCreate(ctx context.Context, c creator, mod *types.Module, records ...*types.Record) (err error) {
	return c.Create(ctx, mod.ModelRef(), recCreateOperations(mod), recToGetters(records...)...)
}

func ComposeRecordUpdate(ctx context.Context, u updater, mod *types.Module, records ...*types.Record) (err error) {
	return u.Update(ctx, mod.ModelRef(), recUpdateOperations(mod), recToGetters(records...)...)
}

func ComposeRecordSoftDelete(ctx context.Context, u updater, mod *types.Module, records ...*types.Record) (err error) {
	return u.Update(ctx, mod.ModelRef(), recUpdateOperations(mod), recToGetters(records...)...)
}

func ComposeRecordDelete(ctx context.Context, d deleter, mod *types.Module, records ...*types.Record) (err error) {
	return d.Delete(ctx, mod.ModelRef(), recDeleteOperations(mod), recToGetters(records...)...)
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

func prepFilter(filter types.RecordFilter, mod *types.Module) filter.Filter {
	return filter.ToConstraintedFilter(mod.Config.DAL.Constraints)
}

func prepIterator(ctx context.Context, dal searcher, mod *types.Module, filter types.RecordFilter) (iter dal.Iterator, err error) {
	dalFilter := prepFilter(filter, mod)

	iter, err = dal.Search(ctx, mod.ModelRef(), recSearchOperations(mod, filter), dalFilter)
	return
}

// drains iterator and collects all records
//
// Collection of records is done with respect to check function and limit constraint on record filter
// For any other filter constraint we assume that underlying DAL took care of it
func drainIterator(ctx context.Context, iter dal.Iterator, mod *types.Module, f types.RecordFilter) (set types.RecordSet, outFilter types.RecordFilter, err error) {
	// close iterator after we've drained it
	defer iter.Close()

	var (
		r *types.Record
	)

	// Get the requested number of record
	if f.Limit > 0 {
		set = make(types.RecordSet, 0, f.Limit)
	} else {
		set = make(types.RecordSet, 0, 1000)
	}

	// Make out filter
	outFilter = f
	pp := f.Paging.Clone()

	err = dal.IteratorPaging(ctx, iter, pp, f.Sorting, func(i dal.Iterator) (out dal.ValueGetter, ok bool) {
		r = prepareRecordTarget(mod)
		if err = i.Scan(r); err != nil {
			return
		}

		// check fetched record
		if f.Check != nil {
			if ok, err = f.Check(r); err != nil {
				return nil, false
			} else if !ok {
				return nil, ok
			}
		}

		if f.Limit == 0 || uint(len(set)) < f.Limit {
			set = append(set, r)
		}

		return r, true
	})
	if err != nil {
		return
	}

	outFilter.Paging = *pp

	return
}

func prepareRecordTarget(module *types.Module) *types.Record {
	// so we can avoid some code later involving (non)partitioned modules :seenoevil:
	r := &types.Record{
		ModuleID:    module.ID,
		NamespaceID: module.NamespaceID,
		Values:      make(types.RecordValueSet, 0, len(module.Fields)),
	}
	r.SetModule(module)

	return r
}

func recToGetters(rr ...*types.Record) (out []dal.ValueGetter) {
	out = make([]dal.ValueGetter, len(rr))

	for i := range rr {
		out[i] = rr[i]
	}

	return
}

func recCreateOperations(m *types.Module) (out dal.OperationSet) {
	return dal.CreateOperations()
}

func recUpdateOperations(m *types.Module) (out dal.OperationSet) {
	return dal.UpdateOperations()
}

func recDeleteOperations(m *types.Module) (out dal.OperationSet) {
	return dal.DeleteOperations()
}

func recFilterOperations(f types.RecordFilter) (out dal.OperationSet) {
	if f.PageCursor != nil {
		out = append(out, dal.Paging)
	}

	if f.IncPageNavigation {
		out = append(out, dal.Paging)
	}

	if f.Sort != nil {
		out = append(out, dal.Sorting)
	}

	return
}

func recSearchOperations(m *types.Module, f types.RecordFilter) (out dal.OperationSet) {
	return dal.SearchOperations().
		Union(recFilterOperations(f))
}

func recLookupOperations(m *types.Module) (out dal.OperationSet) {
	return dal.LookupOperations()
}
