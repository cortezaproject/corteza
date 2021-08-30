package report

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/spf13/cast"
)

// strategizeLoad uses the given context to determine what join strategy we should
// use to achieve the correct result as optimally as possible
func (d *joinedDataset) strategizeLoad(ctx context.Context, inverted bool, local, foreign *FrameDefinition) (ml *frameBuffer, sl *frameBuffer, err error) {
	if !inverted {
		return d.stratLocalMain(ctx, local, foreign)
	}

	return d.stratForeignMain(ctx, local, foreign)
}

// stratLocalMain uses the local datasource for the main loader
//
// The strategy is used always, except for when the initial sort is controlled
// by the foreign datasource.
//
// In this strategy, the local datasource controlls the initial sort and the
// filter that should be applied when pulling data from the foreign source.
func (d *joinedDataset) stratLocalMain(ctx context.Context, local, foreign *FrameDefinition) (ml *frameBuffer, sl *frameBuffer, err error) {
	var ok bool

	// Prepare the main loader from the local source
	// Local source can be partitioned to support nested joining.
	if d.partitioned {
		ok, err = (d.local.(PartitionableDatasource)).Partition(d.partitionSize, d.partitionCol)
		if err != nil {
			return
		}
		if !ok {
			err = fmt.Errorf("local datasource is not partitionable: %s", d.local.Name())
			return
		}
	}

	ldr, clsr, err := d.local.Load(ctx, local)
	if err != nil {
		return
	}

	ml = &frameBuffer{
		sourceName:  d.Name(),
		keyCol:      d.def.LocalColumn,
		keyColIndex: -1,

		loader: func(_ *Filter, cap uint) ([]*Frame, error) {
			return ldr(int(cap), false)
		},
		closer:    clsr,
		chunkSize: local.Paging.Limit,
		sorting:   local.Sort,

		// Overfetch frames when the last two entries define the same sort.
		// This is required to support paging.
		more: func(ff []*Frame, sc []int) bool {
			// No sorting, we don't care
			if len(local.Sort) == 0 {
				return false
			}

			// With sorting and using primary/unique columns, we don't care
			for _, s := range local.Sort {
				c := local.Columns[local.Columns.Find(s.Column)]
				if c.Primary || c.Unique {
					return false
				}
			}

			// This is the last frame we can pull out, we don't care
			if uint(ff[len(ff)-1].Size()) < local.Paging.Limit {
				return false
			}

			// With sorting and regular'ol columns, we care only if the over-fetched row
			// is the same as the last requested row.
			//
			// If we have multiple frames that means that it is partitioned; check whole frames
			if len(ff) > 1 {
				tmpl := len(ff)
				return ff[tmpl-1].FirstRow().Compare(ff[tmpl-2].FirstRow(), sc...) == 0
			}

			// Else it is a regular'ol regular'ol frame; check rows
			return ff[0].LastRow().Compare(ff[0].LastLastRow(), sc...) == 0
		},
	}

	// Partitioned sources need to be sliced
	if d.partitioned {
		ml.postFetch = func(f []*Frame) ([]*Frame, error) {
			ff, err := d.sliceFrames(f, d.def.LocalColumn, d.def.ForeignColumn)
			for _, f := range ff {
				f.RefValue = ""
				f.RelColumn = ""
			}
			return ff, err
		}
	}

	//
	//

	// Prepare the sub loader from the foreign source
	subDS, ok := d.foreign.(PartitionableDatasource)
	if !ok {
		// @todo allow alternatives also
		err = fmt.Errorf("foreign datasource is not partitionable")
		return
	}

	ffilter := foreign.Filter.Clone()
	sl = &frameBuffer{
		sourceName: d.Name(),
		// meta...
		keyCol:      d.def.ForeignColumn,
		keyColIndex: -1,

		chunkSize: foreign.Paging.Limit,
		loader: func(keyFilter *Filter, cap uint) ([]*Frame, error) {
			foreign.Filter = merger(ffilter.Clone(), keyFilter, "and")

			ok, err := subDS.Partition(cap, d.def.ForeignColumn)
			if err != nil {
				return nil, err
			}
			if !ok {
				return nil, fmt.Errorf("foreign datasource is not partitionable: %s", d.foreign.Name())
			}
			loader, closer, err := subDS.Load(ctx, foreign)
			if closer != nil {
				defer closer()
			}
			if err != nil {
				return nil, err
			}
			return loader(0, false)
		},
		sorting: foreign.Sort,

		more: func(_ []*Frame, _ []int) bool {
			return false
		},

		postFetch: func(ff []*Frame) ([]*Frame, error) {
			ff, err = d.sliceFrames(ff, d.def.ForeignColumn, d.def.LocalColumn)
			if err != nil {
				return nil, err
			}

			for i := range ff {
				if ff[i].Name == "" {
					ff[i].Name = foreign.Name
				}
				if ff[i].Source == "" {
					ff[i].Source = foreign.Source
				}
				if ff[i].Ref == "" {
					ff[i].Ref = foreign.Ref
				}
			}

			return ff, nil
		},
	}

	return
}

func (d *joinedDataset) stratForeignMain(ctx context.Context, local, foreign *FrameDefinition) (ml *frameBuffer, sl *frameBuffer, err error) {
	// - main
	mainDS, ok := d.foreign.(PartitionableDatasource)
	if !ok {
		// @todo allow alternatives also
		err = fmt.Errorf("foreign datasource is not partitionable")
		return
	}

	partitionSize := foreign.Paging.Limit
	ok, err = mainDS.Partition(partitionSize, d.def.ForeignColumn)
	if err != nil {
		return
	}
	if !ok {
		err = fmt.Errorf("foreign datasource is not partitionable")
		return
	}
	ldr, clsr, err := mainDS.Load(ctx, foreign)
	ml = &frameBuffer{
		// meta...
		keyCol:      d.def.ForeignColumn,
		keyColIndex: -1,

		loader: func(_ *Filter, cap uint) ([]*Frame, error) {
			return ldr(int(cap), false)
		},
		closer: clsr,
		// This is the maximum if all of the pages have full partitions
		chunkSize: local.Paging.Limit * partitionSize,
		sorting:   foreign.Sort,

		// In case where there is no unique column present in the sort, we need
		// to overfetch rows until we reach a row that has a different sort.
		more: func(ff []*Frame, sc []int) bool {
			f := ff[0]
			// No sorting, we don't care
			if len(foreign.Sort) == 0 {
				return false
			}

			// This is the last frame we can pull out, we don't care
			if uint(f.Size()) < foreign.Paging.Limit {
				return false
			}

			// With sorting and using primary/unique columns, we don't care
			for _, s := range foreign.Sort {
				c := foreign.Columns[foreign.Columns.Find(s.Column)]
				if c.Primary || c.Unique {
					return false
				}
			}

			// Only the last frame of the buffer is passed in here, so when it is too small
			if f.Size() <= 1 {
				return true
			}

			// With sorting and regular'ol columns, we care only if the over-fetched row
			// is the same as the last requested row.
			return f.LastRow().Compare(f.LastLastRow(), sc...) == 0
		},

		postFetch: func(f []*Frame) ([]*Frame, error) {
			return d.sliceFrames(f, d.def.ForeignColumn, d.def.LocalColumn)
		},
	}

	// - sub
	if d.partitioned {
		ok, err = (d.local.(PartitionableDatasource)).Partition(d.partitionSize, d.partitionCol)
		if err != nil {
			return
		}
		if !ok {
			err = fmt.Errorf("foreign datasource is not partitionable")
			return
		}
	}

	lfilter := local.Filter.Clone()
	sl = &frameBuffer{
		// meta...
		keyCol:      d.def.LocalColumn,
		keyColIndex: -1,

		loader: func(keyFilter *Filter, _ uint) ([]*Frame, error) {
			// @todo not ok for multiple pulls!
			local.Filter = merger(lfilter.Clone(), keyFilter, "and")

			loader, closer, err := d.local.Load(ctx, local)
			if closer != nil {
				defer closer()
			}
			if err != nil {
				return nil, err
			}
			return loader(0, false)
		},
		sorting: local.Sort,

		more: func(ff []*Frame, sc []int) bool {
			return false
		},
	}

	if d.partitioned {
		sl.postFetch = func(f []*Frame) ([]*Frame, error) {
			ff, err := d.sliceFrames(f, d.def.LocalColumn, d.def.ForeignColumn)
			for _, f := range ff {
				f.RefValue = ""
				f.RelColumn = ""
			}
			return ff, err
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Sorting

// shouldSort determines if we need to perform additional sorting based on the
// provided frame definition sorting.
//
// * If no sort was requested; consider it sorted.
// * If all sort expressions use only the local source; consider it sorted.
// * If at least one sort expression references a foreign source; consider it NOT sorted.
func (d *joinedDataset) shouldSort(ss filter.SortExprSet) bool {
	if len(ss) == 0 {
		return false
	}

	for _, s := range ss {
		if strings.Contains(s.Column, ".") {
			return true
		}
	}

	return false
}

func (d *joinedDataset) strategizeSort(main, sub *frameBuffer, inverted bool, lfd *FrameDefinition, keyColumn string) (err error) {
	var local, foreign *frameBuffer

	// Determine which one was local/foreign
	if inverted {
		local = sub
		foreign = main
	} else {
		local = main
		foreign = sub
	}

	// - index foreign frames into buckets; here the foreign sort must be respected
	buckets := make(map[string]int)
	for i, mf := range foreign.localFrames {
		buckets[mf.RefValue] = i
	}

	localSort := make(filter.SortExprSet, 0, len(lfd.Sort))
	for _, s := range lfd.Sort {
		if !strings.Contains(s.Column, ".") {
			localSort = append(localSort, s)
		}
	}

	if len(local.localFrames) == 1 {
		err = d.sortFrameRows(local, localSort, buckets, inverted)
	} else {
		err = d.sortBufferFrames(local, localSort, buckets, inverted)
	}

	if err != nil {
		return
	}

	// - assure bucket's order
	// @todo can we do this in one go alongside the original sort?
	aux := make([]*Frame, len(foreign.localFrames))
	for _, f := range foreign.localFrames {
		aux[buckets[f.RefValue]] = f
	}
	foreign.localFrames = aux

	return
}

func (d *joinedDataset) sortFrameRows(local *frameBuffer, ss filter.SortExprSet, buckets map[string]int, inverted bool) error {
	f := local.localFrames[0]

	// - sort based on bucket index
	var ri FrameRow
	var rj FrameRow

	var bi int
	var bj int

	var ci expr.Comparable
	var ok bool

	sort.SliceStable(f.Rows, func(i, j int) bool {
		ri = f.Rows[i]
		rj = f.Rows[j]
		bi = buckets[cast.ToString(ri[local.keyColIndex].Get())]
		bj = buckets[cast.ToString(rj[local.keyColIndex].Get())]

		if inverted {
			// if bi is before bj, the row bellongs to a complitely different bucket
			// so we shouldn't do anything extra
			if bi < bj {
				return !ss.Reversed()
			}
		}

		// go through the sort definitions and sort based on that
		for si, s := range ss {
			ci, ok = ri[local.sortColumns[si]].(expr.Comparable)
			if !ok {
				return !s.Descending
			}

			r, err := ci.Compare(rj[local.sortColumns[si]])
			if err != nil {
				return s.Descending
			}

			if r != 0 {
				if s.Descending {
					return r > 0
				}
				return r < 0
			}
		}

		return bi < bj
	})

	f.WalkRows(func(i int, r FrameRow) error {
		buckets[cast.ToString(r[local.keyColIndex].Get())] = i
		return nil
	})

	return nil
}

func (d *joinedDataset) sortBufferFrames(local *frameBuffer, ss filter.SortExprSet, buckets map[string]int, inverted bool) error {
	f := local.localFrames[0]

	var ri *Frame
	var rj *Frame

	var bi int
	var bj int

	var ci expr.Comparable
	var ok bool

	sort.SliceStable(local.localFrames, func(i, j int) bool {
		ri = local.localFrames[i]
		rj = local.localFrames[j]
		bi = buckets[cast.ToString(ri.Rows[0][local.keyColIndex].Get())]
		bj = buckets[cast.ToString(rj.Rows[0][local.keyColIndex].Get())]

		if inverted {
			// if bi is before bj, the row bellongs to a complitely different bucket
			// so we shouldn't do anything extra
			if bi < bj {
				return !ss.Reversed()
			}
		}

		// go through the sort definitions and sort based on that
		for si, s := range ss {
			ci, ok = ri.Rows[0][local.sortColumns[si]].(expr.Comparable)
			if !ok {
				return !s.Descending
			}

			r, err := ci.Compare(rj.Rows[0][local.sortColumns[si]])
			if err != nil {
				return s.Descending
			}

			if r != 0 {
				if s.Descending {
					return r > 0
				}
				return r < 0
			}
		}

		return bi < bj
	})

	f.WalkRows(func(i int, r FrameRow) error {
		buckets[cast.ToString(r[local.keyColIndex].Get())] = i
		return nil
	})

	return nil
}
