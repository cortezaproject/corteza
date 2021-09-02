package report

import (
	"context"
	"errors"
	"fmt"
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

		local   Datasource
		foreign Datasource

		partitioned   bool
		partitionSize uint
		partitionCol  string
	}

	JoinStepDefinition struct {
		Name          string  `json:"name"`
		LocalSource   string  `json:"localSource"`
		LocalColumn   string  `json:"localColumn"`
		ForeignSource string  `json:"foreignSource"`
		ForeignColumn string  `json:"foreignColumn"`
		Filter        *Filter `json:"filter,omitempty"`
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
		return nil, fmt.Errorf("foreign join datasources not defined: %s", j.def.LocalSource)
	}

	// @todo temporarily disabled
	for _, d := range dd {
		if _, ok := d.(*joinedDataset); ok {
			return nil, fmt.Errorf("unable to join a joined source: %s", d.Name())
		}
	}

	// @todo multiple joins
	return &joinedDataset{
		def:     j.def,
		local:   dd[0],
		foreign: dd[1],
	}, nil
}

func (j *stepJoin) Validate() error {
	pfx := "invalid join step: "
	switch {
	case j.def.Name == "":
		return errors.New(pfx + "dimension name not defined")

	case j.def.LocalSource == "":
		return errors.New(pfx + "local dimension not defined")
	case j.def.LocalColumn == "":
		return errors.New(pfx + "local column not defined")
	case j.def.ForeignSource == "":
		return errors.New(pfx + "foreign dimension not defined")
	case j.def.ForeignColumn == "":
		return errors.New(pfx + "foreign column not defined")

	default:
		return nil
	}
}

func (d *stepJoin) Name() string {
	return d.def.Name
}

func (d *stepJoin) Source() []string {
	return []string{d.def.LocalSource, d.def.ForeignSource}
}

func (d *stepJoin) Def() *StepDefinition {
	return &StepDefinition{Join: d.def}
}

// // // //

func (d *joinedDataset) Name() string {
	return d.def.Name
}

func (d *joinedDataset) Describe() FrameDescriptionSet {
	out := make(FrameDescriptionSet, 0, 2)

	local := d.local.Describe()
	for _, l := range local {
		l.Source = d.Name()
		if l.Ref == "" {
			l.Ref = l.Source
		}
		out = append(out, l)
	}

	foreign := d.foreign.Describe()
	for _, f := range foreign {
		f.Source = d.Name()
		if f.Ref == "" {
			f.Ref = f.Source
		}
		out = append(out, f)
	}

	return out
}

// Partition marks the DS to partition the response over the given column
func (d *joinedDataset) Partition(partitionSize uint, partitionCol string) (bool, error) {
	if _, ok := d.local.(PartitionableDatasource); !ok {
		return false, fmt.Errorf("local datasource is not partitionable")
	}

	if d.partitioned {
		return true, nil
	}
	if partitionCol == "" {
		return false, errors.New("unable to partition: partition column not defined")
	}

	d.partitioned = true
	d.partitionCol = partitionCol
	d.partitionSize = partitionSize
	return true, nil
}

func (d *joinedDataset) Load(ctx context.Context, dd ...*FrameDefinition) (l Loader, c Closer, err error) {
	dscr := d.Describe()

	// Preparation
	// - Assure sort columns for paging purposes
	// - Assure local/foreign definitions
	//   Keep a cloned original version so we don't overwrite the initial definition.
	oLocalDef, oForeignDef, err := d.prepareDefinitions(FrameDefinitionSet(dd), dscr)
	if err != nil {
		return
	}

	// - Validate the sort of the local frame for paging purposes
	err = d.validateSort(oLocalDef, dscr)
	if err != nil {
		return
	}

	localDef := oLocalDef.Clone()
	foreignDef := oForeignDef.Clone()

	// - Preprocess sorting definitions for additional context.
	//   The join is inverted when the foreign DS governs the initial sort.
	inverted, err := d.prepareSorting(localDef, foreignDef)
	if err != nil {
		return
	}

	// - Preprocess additional paging filtering
	var pp []partialPagingCnd
	if oLocalDef.Paging.PageCursor != nil {
		pp, err = d.strategizePaging(localDef, foreignDef, inverted)
		if err != nil {
			return
		}
	}

	// - Determine the join strategy to use.
	//   For now the strategy is the combination of the following parameters, in the
	//   future we can have a wrapper struct as well.
	mainLoader, subLoader, err := d.strategizeLoad(ctx, inverted, localDef, foreignDef)
	if err != nil {
		return
	}

	// Load and join data
	//
	// The loader function will iterate indefinitely until the requested frame
	// definition is satisfied.
	//
	// Outline:
	//   . prepare additional (partial) filters based on the paging cursor
	//   . pull data from the main source
	//   . prepare an additional key-based filter for sub source
	//   . pull data from the sub source
	//   . additional processing
	//   .. apply additional filtering based on page cursors
	//   . prepare response
	isEmpty := false
	return func(cap int, processed bool) (oo []*Frame, err error) {
			var keys []string
			if isEmpty {
				return
			}
			if processed {
				cap++
			}

			// The modified flag will help us determine if we need another iteration or not.
			// The flag is only set to true if we do any additional modifications in here.
			//
			// Underlying loaders must be able to provide all of the requested data, so if
			// we don't do any modifications here, the iteration should not repeat.
			modified := false
			// m is defined here so we don't re initialize it every time later on
			m := false

			pagingSatisfied := oLocalDef.Paging.PageCursor == nil
			more := false
			for {
				modified = false

				// . Pull data from the main source
				more, err = mainLoader.load(uint(cap), nil)
				if err != nil {
					return
				}
				if !more {
					isEmpty = true
					break
				}

				// . Prepare an additional key-based filter for sub sources
				keys, err = mainLoader.keys()
				if err != nil {
					return
				}
				keyFilter := subLoader.keyFilter(keys)

				// . Pull data from the sub source
				more, err = subLoader.load(uint(cap), keyFilter)
				if err != nil {
					return nil, err
				}
				if !more {
					break
				}

				// Cut out any local rows that do not have any foreign rows
				modified = d.validateRows(mainLoader, subLoader, inverted)
				if modified {
					continue
				}

				// . Sort the collected data
				//
				// @todo Most of the sorting has already been done by the datasource,
				//       so we can get away with a lot of partial processing.
				// if d.shouldSort(oLocalDef.Sort) {
				// }
				err = d.strategizeSort(mainLoader, subLoader, inverted, oLocalDef, d.def.LocalColumn)
				if err != nil {
					return nil, err
				}

				// . Additional processing
				//
				// Any additional filtering and goes here.
				// Transformations and other bits should reside in the buffers when loading data.

				// .. Additional filters based on page cursors
				if !pagingSatisfied {
					m, pagingSatisfied = d.pagingFilter(mainLoader, subLoader, pp)
					modified = modified || m
				}

				// . Prepare response
				//
				// If there were no modifications to what we pulled, we can safely
				// assume that we can produce a response.
				if !modified {
					break
				}
			}

			return prepareResponse(mainLoader, subLoader, inverted, processed, oLocalDef, d.def.LocalColumn, dscr)
		}, func() {
			if mainLoader.closer != nil {
				mainLoader.closer()
			}
			if subLoader.closer != nil {
				subLoader.closer()
			}
		}, nil
}

func (d *joinedDataset) prepareDefinitions(dd FrameDefinitionSet, dscr FrameDescriptionSet) (localDef *FrameDefinition, foreignDef *FrameDefinition, err error) {
	if len(dd) == 0 {
		err = errors.New("joining requires at least one frame definition")
		return
	}

	localDef = FrameDefinitionSet(dd).FindBySourceRef(d.Name(), d.def.LocalSource)
	foreignDef = FrameDefinitionSet(dd).FindBySourceRef(d.Name(), d.def.ForeignSource)

	if localDef == nil {
		localDef = &FrameDefinition{
			Name:   dd[0].Name,
			Source: d.Name(),
			Ref:    d.def.LocalSource,
			Paging: dd[0].Paging,
			Sort:   dd[0].Sort,
			Filter: dd[0].Filter,
		}
	}

	if foreignDef == nil {
		foreignDef = &FrameDefinition{
			Name:   dd[0].Name,
			Source: d.Name(),
			Ref:    d.def.ForeignSource,
			Paging: &filter.Paging{
				Limit: localDef.Paging.Limit,
			},
			Sort: filter.SortExprSet{},
		}
	}

	if len(localDef.Columns) == 0 {
		dscr = d.local.Describe()
		sc := dscr.FilterBySource(localDef.Ref)[0]
		localDef.Columns = sc.Columns
	}

	if len(foreignDef.Columns) == 0 {
		dscr = d.foreign.Describe()
		sc := dscr.FilterBySource(foreignDef.Ref)[0]
		foreignDef.Columns = sc.Columns
	}

	return
}

func (d *joinedDataset) sliceFrames(ff []*Frame, selfCol, relCol string) (out []*Frame, err error) {
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
				RefValue:  k,
				RelColumn: relCol,
				Columns:   ff[0].Columns,
				Paging:    ff[0].Paging,
				Sort:      ff[0].Sort,
				Filter:    ff[0].Filter,
				Ref:       ff[0].Ref,
			})
		}

		out[i].Rows = append(out[i].Rows, r)
	}

	for _, f := range ff {
		// slice the output; one frame per key
		var k string
		fColI := f.Columns.Find(selfCol)
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

func (d *joinedDataset) validateSort(def *FrameDefinition, dd FrameDescriptionSet) (err error) {
	sortDS := ""
	auxSS := make(filter.SortExprSet, 0, len(def.Sort))

	for i := len(def.Sort) - 2; i >= 0; i-- {
		aa := strings.Split(def.Sort[i].Column, ".")
		bb := strings.Split(def.Sort[i+1].Column, ".")

		if len(aa) != len(bb) || (len(aa) > 1 && aa[0] != aa[1]) {
			auxSS = append(auxSS, def.Sort[i+1])
			break
		}

		auxSS = append(auxSS, def.Sort[i+1])
		if len(aa) > 1 {
			sortDS = aa[0]
		}
	}

	// @todo temporarily disabled
	if sortDS != "" {
		return fmt.Errorf("[temporary] initial sort can not by by a foreign column: %s", sortDS)
	}

	// The first one is always local so this is ok
	localDscr := dd[0]
	var dscr *FrameDescription

	// When local, ref is omitted
	if sortDS == "" {
		// Do this to avoid extra work afterwords
		auxSS = def.Sort

		dscr = localDscr
		sortDS = dscr.Ref
	} else {
		dscr = dd.FilterByRef(sortDS)[0]
	}

	// Check if we're sorting by a unique value
	def.Sort = func() filter.SortExprSet {
		unique := ""
		for _, c := range dscr.Columns {
			if c.Primary || c.Unique {
				if unique == "" {
					unique = c.Name
				}
				if auxSS.Get(c.Name) != nil {
					return def.Sort
				}
			}
		}

		if sortDS == localDscr.Ref {
			return append(def.Sort, &filter.SortExpr{Column: unique, Descending: auxSS.LastDescending()})
		} else {
			return append(def.Sort, &filter.SortExpr{Column: fmt.Sprintf("%s.%s", sortDS, unique), Descending: auxSS.LastDescending()})
		}

	}()
	return nil
}

func (d *joinedDataset) prepareSorting(local, foreign *FrameDefinition) (inverted bool, err error) {
	var (
		localSS   filter.SortExprSet
		foreignSS filter.SortExprSet
	)

	foreignDS := ""

	for i, s := range local.Sort {
		spts := strings.Split(s.Column, ".")
		if len(spts) > 1 {
			inverted = i == 0
			if foreignDS != "" {
				if foreignDS != spts[0] {
					// @todo allow this also
					err = fmt.Errorf("cannot sort local datasource by multiple foreign datasources: %s, %s", foreignDS, spts[0])
					return
				}
			} else {
				foreignDS = spts[0]
			}

			foreignSS = append(foreignSS, &filter.SortExpr{Column: spts[1], Descending: s.Descending})
		} else {
			localSS = append(localSS, s)
		}
	}

	local.Sort = localSS
	foreign.Sort = append(foreignSS, foreign.Sort...)

	return
}

func (d *joinedDataset) validateRows(local, foreign *frameBuffer, inverted bool) (modified bool) {
	keyCol := foreign.localFrames[0].RelColumn
	cols := local.localFrames[0].Columns
	keyColIx := cols.Find(keyCol)

	// - index foreign frames into buckets; here the foreign sort must be respected
	buckets := make(map[string]int)
	for i, mf := range foreign.localFrames {
		buckets[mf.RefValue] = i
	}

	var k string
	aux := make([]FrameRow, 0, 100)
	local.walkRowsLocal(func(i int, r FrameRow) error {
		v := r[keyColIx]
		if v == nil {
			k = ""
		} else {
			k = cast.ToString(v.Get())
		}

		if _, ok := buckets[k]; ok {
			aux = append(aux, r)
		}

		return nil
	})

	modified = len(local.localFrames[0].Rows) != len(aux)
	local.localFrames[0].Rows = aux

	return

}
