package report

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/spf13/cast"
)

type (
	stepJoin struct {
		def *JoinStepDefinition
	}

	joinedDataset struct {
		def *JoinStepDefinition

		base    Datasource
		foreign Datasource
	}

	JoinStepDefinition struct {
		Name    string         `json:"name"`
		Local   string         `json:"local"`
		Foreign string         `json:"foreign"`
		Rows    *RowDefinition `json:"rows,omitempty"`
	}
)

const (
	defaultPartitionSize = uint(20)
)

func (j *stepJoin) Run(ctx context.Context, dd ...Datasource) (Datasource, error) {
	if len(dd) == 0 {
		return nil, fmt.Errorf("unknown join datasources")
	}

	if len(dd) < 2 {
		return nil, fmt.Errorf("foreign join datasources not defined: %s", j.def.localDim())
	}

	// @todo multiple joins
	return &joinedDataset{
		def:     j.def,
		base:    dd[0],
		foreign: dd[1],
	}, nil
}

func (j *stepJoin) Validate() error {
	pfx := "invalid join step: "
	switch {
	case j.def.Name == "":
		return errors.New(pfx + "dimension name not defined")

	case j.def.localDim() == "":
		return errors.New(pfx + "local dimension not defined")
	case j.def.localColumn() == "":
		return errors.New(pfx + "local column not defined")
	case j.def.foreignDim() == "":
		return errors.New(pfx + "foreign dimension not defined")
	case j.def.foreignColumn() == "":
		return errors.New(pfx + "foreign column not defined")

	default:
		return nil
	}
}

func (d *stepJoin) Name() string {
	return d.def.Name
}

func (d *stepJoin) Source() []string {
	return []string{d.def.localDim(), d.def.foreignDim()}
}

func (d *stepJoin) Def() *StepDefinition {
	return &StepDefinition{Join: d.def}
}

// // // //

func (d *joinedDataset) Name() string {
	return d.def.Name
}

// @todo allow x-join filtering
func (d *joinedDataset) Load(ctx context.Context, dd ...*FrameDefinition) (Loader, Closer, error) {
	// to hold closer references for all underlying datasources
	closers := make([]Closer, 0, 10)

	return func(cap int) ([]*Frame, error) {
			// determine local and foreign frame definitions
			localDef := FrameDefinitionSet(dd).FindBySourceRef(d.Name(), d.def.localDim())
			foreignDef := FrameDefinitionSet(dd).FindBySourceRef(d.Name(), d.def.foreignDim())

			// basic pre-run validation
			// - definitions
			if localDef == nil {
				return nil, fmt.Errorf("definition for local datasource not found: %s", d.def.localDim())
			}
			if foreignDef == nil {
				return nil, fmt.Errorf("definition for foreign datasource not found: %s", d.def.foreignDim())
			}
			if localDef.Paging == nil {
				localDef.Paging = &filter.Paging{}
			}
			if foreignDef.Paging == nil {
				foreignDef.Paging = &filter.Paging{}
			}

			// - key columns
			if localDef.Columns.Find(d.def.localColumn()) < 0 {
				return nil, fmt.Errorf("local frame definition must include the key column: %s", d.def.localColumn())
			}
			if foreignDef.Columns.Find(d.def.foreignColumn()) < 0 {
				return nil, fmt.Errorf("foreign frame definition must include the key column: %s", d.def.foreignColumn())
			}

			// based on the passed sort, determine main/sub datasources
			//
			// if the local datasource is being initially sorted by the foreign
			// datasource, we need to resolve the foreign datasource first and
			// adjust the local sorting based on that
			localSort, foreignSort, foreignDS, err := d.splitSort(localDef.Sorting)
			localDef.Sorting = localSort
			useSubSort := foreignDS != ""
			if err != nil {
				return nil, err
			}

			if foreignDS != "" {
				if foreignDS != d.foreign.Name() {
					return nil, fmt.Errorf("foreign sort datasource not part of the join: %s", foreignDS)
				}

				foreignDef.Sorting = append(foreignSort, foreignDef.Sorting...)
			}

			// pull frames from the datasource that defines the initial sort
			var mainLoader Loader
			var mainCloser Closer
			var mainPageCap uint
			// - when using foreign for base sort, firstly pull frames from the foreign datasource
			if useSubSort {
				prtDS, ok := d.foreign.(PartitionableDatasource)
				if !ok {
					// @todo allow alternatives also
					return nil, fmt.Errorf("foreign datasource is not partitionable")
				}

				// - determine partition size
				partitionSize := defaultPartitionSize
				if foreignDef.Paging != nil && foreignDef.Paging.Limit > 0 {
					partitionSize = foreignDef.Paging.Limit
				}

				// - determine local limit
				localLimit := defaultPageSize
				if localDef.Paging != nil && localDef.Paging.Limit > 0 {
					localLimit = localDef.Paging.Limit
				}

				// - determine maximum foreign page size after partitioning
				//   +1 because we are going one over for paging stuff
				mainPageCap = (partitionSize + 1) * (localLimit + 1)
				foreignDef.Paging.Limit = mainPageCap

				// - prepare loader, closer
				mainLoader, mainCloser, err = prtDS.Partition(ctx, partitionSize, d.def.foreignColumn(), foreignDef)
			} else {
				mainPageCap = defaultPageSize
				if localDef.Paging != nil && localDef.Paging.Limit > 0 {
					mainPageCap = localDef.Paging.Limit
				}
				// nothing special needed
				mainLoader, mainCloser, err = d.base.Load(ctx, localDef)
			}
			if mainCloser != nil {
				defer mainCloser()
			}
			if err != nil {
				return nil, err
			}

			// pull rows from the main datasource
			mainFrames, err := mainLoader(int(mainPageCap))
			if err != nil {
				return nil, err
			}

			if useSubSort {
				// here we need to slice the partitioned datasource
				// @todo should this be layed off to the lowe level?
				mainFrames, err = d.sliceFramesFurther(mainFrames, d.def.foreignColumn())
				if err != nil {
					return nil, err
				}

				for i := range mainFrames {
					mainFrames[i].Name = foreignDef.Name
					mainFrames[i].Source = foreignDef.Source
					mainFrames[i].Ref = foreignDef.Ref
				}
			} else {
				for i := range mainFrames {
					mainFrames[i].Name = localDef.Name
					mainFrames[i].Source = localDef.Source
					mainFrames[i].Ref = localDef.Ref
				}
			}

			// determine keys to filter over sub datasource
			var mainKeyColIndex int
			keys := make([]string, 0, defaultPageSize)
			keySet := make(map[string]bool)
			if useSubSort {
				mainKeyColIndex = mainFrames[0].Columns.Find(d.def.foreignColumn())
				if mainKeyColIndex < 0 {
					return nil, fmt.Errorf("key column on foreign datasource does not exist: %s", d.def.foreignColumn())
				}
			} else {
				mainKeyColIndex = mainFrames[0].Columns.Find(d.def.localColumn())
				if mainKeyColIndex < 0 {
					return nil, fmt.Errorf("key column on local datasource does not exist: %s", d.def.localColumn())
				}
			}
			var k string
			if useSubSort {
				// since these frames are grouped by the ref value, we can use any one of the rows
				// and we don't have to look at all of the rows

				for _, mf := range mainFrames {
					r := mf.FirstRow()
					k, err = cast.ToStringE(r[mainKeyColIndex].Get())
					if ok := keySet[k]; !ok {
						keys = append(keys, k)
						keySet[k] = true
					}
					if err != nil {
						return nil, err
					}
				}
			} else {
				for _, mf := range mainFrames {
					err = mf.WalkRows(func(i int, r FrameRow) error {
						k, err = cast.ToStringE(r[mainKeyColIndex].Get())
						if ok := keySet[k]; !ok {
							keys = append(keys, k)
							keySet[k] = true
						}
						return err
					})
					if err != nil {
						return nil, err
					}
				}
			}

			// filter over sub datasource
			var subFrames []*Frame
			if useSubSort {
				// here we use the LOCAL datasource, because it's flipped
				// - prepare key pre-filter
				localDef.Rows = d.keySliceToFilter(d.def.localColumn(), keys).MergeAnd(localDef.Rows)

				// - go!
				loader, closer, err := d.base.Load(ctx, localDef)
				if closer != nil {
					defer closer()
				}
				if err != nil {
					return nil, err
				}
				subFrames, err = loader(0)

				// -- slice keys based on mainFrames and their sort key matching
				//
				// Each key from mainFrames points to a bucket where all of the
				// sort columns match.
				// Later used for easier sorting.
				// @todo join with above mainFrames traversal?
				var (
					prevSortBucketKey = ""
					sortBucketKey     = ""
					buckets           = make(map[string]int)
					bucketIndex       = 0
				)

				sortKeyColIndexes := make([]int, 0, len(foreignDef.Sorting))
				for _, s := range foreignDef.Sorting {
					sortKeyColIndexes = append(sortKeyColIndexes, mainFrames[0].Columns.Find(s.Column))
				}

				for _, mf := range mainFrames {
					// since these frames are grouped by the ref value, we can use any one of the rows
					// and we don't have to look at all of the rows
					r := mf.FirstRow()

					// get key value
					k = mf.RefKey

					// get sort value
					sortBucketKey, err = d.cellsToString(r, sortKeyColIndexes)
					if err != nil {
						return nil, err
					}

					// initial bucket
					if prevSortBucketKey == "" {
						buckets[k] = bucketIndex
						prevSortBucketKey = sortBucketKey
						continue
					}

					// falls into the same bucket
					if prevSortBucketKey == sortBucketKey {
						buckets[k] = bucketIndex
						continue
					}

					// falls into another bucket
					bucketIndex++
					buckets[k] = bucketIndex
					prevSortBucketKey = sortBucketKey
				}

				localColIndex := subFrames[0].Columns.Find(d.def.localColumn())
				for i := range subFrames {
					subFrames[i], err = d.bucketSort(subFrames[i], buckets, localColIndex, localDef.Sorting...)
					if err != nil {
						return nil, err
					}

					subFrames[i].Name = localDef.Name
					subFrames[i].Source = localDef.Source
					subFrames[i].Ref = localDef.Ref
				}
			} else {
				prtDS, ok := d.foreign.(PartitionableDatasource)
				if !ok {
					// @todo allow alternatives also
					return nil, fmt.Errorf("foreign datasource is not partitionable")
				}

				// - determine partition size
				partitionSize := defaultPartitionSize
				if foreignDef.Paging != nil && foreignDef.Paging.Limit > 0 {
					partitionSize = foreignDef.Paging.Limit
				}

				// - prepare key pre-filter
				foreignDef.Rows = d.keySliceToFilter(d.def.foreignColumn(), keys).MergeAnd(foreignDef.Rows)

				// - go!
				loader, closer, err := prtDS.Partition(ctx, partitionSize, d.def.foreignColumn(), foreignDef)
				if closer != nil {
					defer closer()
				}
				if err != nil {
					return nil, err
				}
				subFrames, err = loader(0)
				subFrames, err = d.sliceFramesFurther(subFrames, d.def.foreignColumn())
				if err != nil {
					return nil, err
				}

				for i := range subFrames {
					subFrames[i].Name = foreignDef.Name
					subFrames[i].Source = foreignDef.Source
					subFrames[i].Ref = foreignDef.Ref
				}
			}
			if err != nil {
				return nil, err
			}

			// just to make sure the local frames are always before foreign frames
			if useSubSort {
				return append(subFrames, mainFrames...), nil
			} else {
				return append(mainFrames, subFrames...), nil
			}
		}, func() {
			for _, c := range closers {
				c()
			}
		}, nil
}

func (d *joinedDataset) splitSort(ss filter.SortExprSet) (local filter.SortExprSet, foreign filter.SortExprSet, foreignDS string, err error) {
	for _, s := range ss {
		spts := strings.Split(s.Column, ".")
		if len(spts) > 1 {
			if foreignDS != "" {
				if foreignDS != spts[0] {
					// @todo allow this also
					err = fmt.Errorf("cannot sort local datasource by multiple foreign datasources: %s, %s", foreignDS, spts[0])
					return
				}
			} else {
				foreignDS = spts[0]
			}

			foreign = append(foreign, &filter.SortExpr{Column: spts[1], Descending: s.Descending})
		} else {
			local = append(local, s)
		}
	}

	return
}

func (d *joinedDataset) keySliceToFilter(col string, keys []string) *RowDefinition {
	cf := &RowDefinition{Or: make([]*RowDefinition, 0, len(keys))}
	for _, k := range keys {
		cf.Or = append(cf.Or, &RowDefinition{
			Cells: map[string]*CellDefinition{
				col: {Op: "eq", Value: "'" + k + "'"},
			},
		})
	}

	return cf
}

func (d *joinedDataset) sliceFramesFurther(ff []*Frame, col string) (out []*Frame, err error) {
	outMap := make(map[string]int)

	cellToString := func(t expr.TypedValue) (string, error) {
		return cast.ToStringE(t.Get())
	}

	push := func(k string, r FrameRow) {
		var i int
		var ok bool
		if i, ok = outMap[k]; !ok {
			i = len(out)
			outMap[k] = i
			out = append(out, &Frame{
				RefKey:  k,
				Columns: ff[0].Columns,
				Paging:  ff[0].Paging,
				Sorting: ff[0].Sorting,
			})
		}

		out[i].Rows = append(out[i].Rows, r)
	}

	for _, f := range ff {
		// slice the output; one frame per key
		var k string
		fColI := f.Columns.Find(col)
		err = f.WalkRows(func(i int, r FrameRow) error {
			k, err = cellToString(r[fColI])
			push(k, r)
			return err
		})
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// @todo make in place
func (d *joinedDataset) bucketSort(f *Frame, buckets map[string]int, bucketColIndex int, ss ...*filter.SortExpr) (*Frame, error) {
	// - prepare sort col indexes
	sortColIndexes := make([]int, len(ss))
	for i, s := range ss {
		sortColIndexes[i] = f.Columns.Find(s.Column)
	}

	// - sort #1 based on key sub-slice order to place common rows together
	//   ordered determined by corresponding bucket index
	var ri FrameRow
	var rj FrameRow

	var bi int
	var bj int

	var ci expr.Comparable
	var ok bool
	sort.SliceStable(f.Rows, func(i, j int) bool {
		ri = f.Rows[i]
		rj = f.Rows[j]
		bi = buckets[cast.ToString(ri[bucketColIndex].Get())]
		bj = buckets[cast.ToString(rj[bucketColIndex].Get())]

		// if bi is before bj, the row bellongs to a complitely different bucket
		// so we shouldn't do anything extra
		if bi < bj {
			return true
		}

		// go through the sort definitions and sort based on that
		for si, s := range ss {
			ci, ok = ri[sortColIndexes[si]].(expr.Comparable)
			if !ok {
				return true
			}

			r, err := ci.Compare(rj[sortColIndexes[si]])
			if err != nil {
				return false
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
	return f, nil
}

func (d *joinedDataset) cellsToString(row FrameRow, indexes []int) (out string, err error) {
	var k string
	for _, i := range indexes {
		k, err = cast.ToStringE(row[i].Get())
		if err != nil {
			return
		}
		out += k
	}

	return
}

func (def *JoinStepDefinition) localDim() string {
	return dimensionOf(def.Local)
}
func (def *JoinStepDefinition) localColumn() string {
	return columnOf(def.Local)
}
func (def *JoinStepDefinition) foreignDim() string {
	return dimensionOf(def.Foreign)
}
func (def *JoinStepDefinition) foreignColumn() string {
	return columnOf(def.Foreign)
}
