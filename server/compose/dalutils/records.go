package dalutils

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
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

	counter interface {
		Count(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (uint, error)
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

func ComposeRecordsCount(ctx context.Context, c counter, mod *types.Module, filter types.RecordFilter) (cnt uint, err error) {
	constraints := map[string][]interface{}{
		"namespaceID": {mod.NamespaceID},
	}

	if filter.ModuleID != 0 {
		constraints["moduleID"] = []interface{}{filter.ModuleID}
	}

	dalFilter := filter.ToConstraintedFilter(constraints)

	return c.Count(ctx, mod.ModelRef(), recLookupOperations(mod), dalFilter)
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

func ComposeRecordUndelete(ctx context.Context, u updater, mod *types.Module, records ...*types.Record) (err error) {
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

	if f.PageCursor != nil {
		if f.IncPageNavigation || f.IncTotal {
			err = fmt.Errorf("not allowed to fetch page navigation or total item count with page cursor")
			return
		}
	}

	var (
		ok       bool
		fetched  uint
		filtered uint
		r        *types.Record
	)

	// Get the requested number of record
	if f.Limit > 0 {
		set = make(types.RecordSet, 0, f.Limit)
	} else {
		set = make(types.RecordSet, 0, 1000)
	}

	for f.Limit == 0 || uint(len(set)) < f.Limit {
		// reset counters every drain
		fetched = 0
		filtered = 0

		err = WalkIterator(ctx, iter, mod, func(r *types.Record) error {
			// check fetched record
			if f.Check != nil {
				if ok, err = f.Check(r); err != nil {
					return err
				} else if !ok {
					filtered++
					return nil
				}
			}

			fetched++
			set = append(set, r)
			return err
		})

		// if an error occurred inside Next()/WalkIterator,
		// we need to stop draining
		if err != nil {
			return
		}

		total := fetched + filtered
		if total == 0 || f.Limit == 0 || (0 < f.Limit && total < f.Limit) {
			// do not re-fetch if:
			// 1) nothing was fetch in the previous run
			// 2) there was no limit (everything was fetched)
			// 3) there are less total (fetched and filtered) items then value of limit
			break
		}

		// Fetch more records
		setLen := uint(len(set))
		if setLen > 0 && setLen < f.Limit {
			fetchMore := f.Limit - setLen

			// request more items
			if err = iter.More(fetchMore, r); err != nil {
				return
			}
		}
	}

	// Get the page nav/total/next-prev cursors
	nav, err := generatePageNavigation(ctx, iter, mod, f, set)
	if err != nil {
		return
	}

	// Make out filter
	outFilter = f
	outFilter.Paging = nav.Paging
	outFilter.Sorting = nav.Sorting

	return
}

// generatePageNavigation generates page navigation for a given record set using an iterator and filter limit.
// If the limit is not defined, the page navigation will consist of only one page without a cursor.
// If the limit is defined and is greater than the total number of records in the set,
// the page navigation will consist of only one page without a cursor.
// If the limit is defined and is less than the total number of records in the set,
// the page navigation will have multiple pages with cursor(s) based on the total number of records and the provided limit.
// @todo revisit and clean up this function properly
func generatePageNavigation(ctx context.Context, iter dal.Iterator, mod *types.Module, p types.RecordFilter, set types.RecordSet) (out types.RecordFilter, err error) {
	const (
		howMuchMore = 1000
	)
	var (
		ok      bool
		total   = uint(len(set))
		setLen  = len(set)
		counter = 0

		first *types.Record
		last  *types.Record
		page  filter.Page

		pageNavigation = []*filter.Page{
			{
				Page:   1,
				Count:  0,
				Cursor: nil,
			},
		}

		// generatePage generates pageNavigation for given record set
		generatePage = func(last *types.Record) (err error) {
			if !p.IncPageNavigation || p.Limit == 0 || len(pageNavigation) == 0 {
				return
			}

			lastNavPageNo := len(pageNavigation) - 1
			nextPage, err := iter.ForwardCursor(last)
			if err != nil {
				return
			}

			if total < p.Limit {
				pageNavigation[lastNavPageNo].Count = total
			}

			// prepare page
			if total != 0 && (total%p.Limit) == 0 {
				pageNavigation[lastNavPageNo].Count = p.Limit
				page = filter.Page{
					Page:   uint(len(pageNavigation) + 1),
					Count:  p.Limit,
					Cursor: nextPage,
				}
			}

			expectedItemCountUpToPage := uint(lastNavPageNo+1) * p.Limit
			if p.Limit == 1 {
				expectedItemCountUpToPage = uint(lastNavPageNo) * p.Limit
			}

			if expectedItemCountUpToPage < total {
				// push page when limit is matched with the previous page item size
				if pageNavigation[lastNavPageNo].Count == p.Limit {
					pageNavigation = append(pageNavigation, &filter.Page{
						Page:   page.Page,
						Count:  total % p.Limit,
						Cursor: page.Cursor, // prev cursor
					})
				}
			}

			return
		}
	)

	if setLen == 0 {
		return
	}

	first = set[0]
	last = set[setLen-1]

	// Limit
	out.Limit = p.Limit

	// Sorting
	out.Sort = dal.IteratorSorting(iter)

	// No need to generate prev/next cursor
	// if limit is not defined and set is empty
	if p.Limit > 0 && len(set) > 0 {
		// PrevPage
		out.PrevPage, err = dal.PreLoadCursor(ctx, iter, 1, true, first)
		if err != nil {
			return
		}

		// NextPage
		out.NextPage, err = dal.PreLoadCursor(ctx, iter, 1, false, last)
		if err != nil {
			return
		}
	}

	if p.IncTotal || p.IncPageNavigation {
		// For the first page nav
		err = generatePage(last)
		if err != nil {
			return
		}

		for counter == 0 || counter < howMuchMore {
			counter++

			interLoop := 0

			if err = iter.More(howMuchMore, last); err != nil {
				return
			}

			err = WalkIterator(ctx, iter, mod, func(rec *types.Record) error {
				// check fetched record
				if p.Check != nil {
					if ok, err = p.Check(rec); err != nil {
						return err
					} else if !ok {
						return nil
					}
				}

				interLoop++
				total++
				last = rec
				return generatePage(rec)
			})
			if err != nil {
				return
			}

			if interLoop < howMuchMore {
				break
			}
		}
	}

	// Total
	if p.IncTotal {
		out.Total = total
	}

	// Page navigation
	if p.IncPageNavigation {
		// Ensure that the last page count is correct if it's not equal to the limit.
		lastPageCount := pageNavigation[len(pageNavigation)-1].Count
		if lastPageCount > 0 && lastPageCount != p.Limit && lastPageCount != total%p.Limit {
			pageNavigation[len(pageNavigation)-1].Count = total % p.Limit
		}

		if p.Limit == 1 {
			pageNavigation = pageNavigation[:len(pageNavigation)-1]
		}

		out.PageNavigation = pageNavigation
	}

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
