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

		// localFrames contains the pulled local frames
		localFrames []*Frame
		// foreignFrames contains the pulled foreign frames grouped by ds
		foreignFrames [][]*Frame
		// foreignSourceIndex maps the ds to the ds index of foreignFrames slice
		foreignSourceIndex map[string]int
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
	err = d.validatePagingSort(oLocalDef, dscr)
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
		pp, err = d.strategizePaging(localDef, foreignDef, inverted, dscr)
		if err != nil {
			return
		}
	}

	// - Do some preparation tasks before the loading occurs
	localLoader, foreignLoader, err := d.strategizeLoading(ctx, inverted, localDef, foreignDef)
	if err != nil {
		return
	}

	// Loading and joining data
	//
	// The loader function will iterate indefinitely until the requested frame
	// definition is satisfied.
	//
	// Outline:
	//   . load data from the two sources (local, foreign)
	//   . remove rows that don't have any related data
	//   . sort the pulled chunk
	//   . apply additional processing
	//   .. apply additional paging cursor based filtering
	//   . update buffers
	//
	// When preparing the response, the paging cursor is calculated.
	isEmpty := false
	return func(cap int, processed bool) (oo []*Frame, err error) {
			if isEmpty {
				return
			}

			// We need to go 1 over to see if we can calculate paging cursors
			loadCap := cap
			if processed {
				loadCap++
			}

			// The modified flag will help us determine if we need another iteration or not.
			// The flag is only set to true if we do any additional modifications in here.
			//
			// Underlying loaders must be able to provide all of the requested data, so if
			// we don't do any modifications here, the iteration should not repeat.
			modified := false
			m := false

			pagingSatisfied := oLocalDef.Paging.PageCursor == nil
			var buffLocal []*Frame
			var buffForeign []*Frame
			for {
				modified = false

				// . load data from the two sources (local, foreign)
				//
				// The loading order differs a bit based on the provided sort.
				if !inverted {
					buffLocal, buffForeign, err = d.localBasedLoad(ctx, loadCap, localLoader, foreignLoader, localDef, foreignDef)
				} else {
					buffLocal, buffForeign, err = d.foreignBasedLoad(ctx, loadCap, localLoader, foreignLoader, localDef, foreignDef)
				}
				if err != nil {
					return
				}

				// .. if both pulled buffersare empty, there is nothing left to do
				if d.sizeOfBuffer(buffLocal, d.partitioned)+d.sizeOfBuffer(buffForeign, true) == 0 {
					break
				}

				// . remove rows that don't have any related data
				buffLocal, buffForeign, modified = d.removeOrphanRows(buffLocal, buffForeign)

				// . sort the pulled chunk
				//
				// The way we are fetching data, most of the sorting is already done.
				// We need to assure the sort based on both local and foreign data sources.
				buffLocal, buffForeign = d.sortFrameBuffers(localLoader, foreignLoader, buffLocal, buffForeign, localDef.Sort, inverted)

				// . apply additional processing
				//
				// Any additional filtering and goes here.
				// Any processing should only occur based on the provided buffer to simplify
				// the entire algorithm.

				// .. apply additional paging cursor based filtering
				//
				// As we're combining data from N sources, the paging cursors are also
				// based on the combination of sources and we can't do the "classical"
				// "cut the data and call it a day" as we do with records.
				//
				// The "primary" DS (the one that defines the initial order) may strictly
				// cut off unneeded rows on the DS level, but the rest need to be calculated
				// manually here.
				//
				// This evaluation must only occur while the paging filter is not yet
				// satisfied.
				if !pagingSatisfied {
					buffLocal, buffForeign, m, pagingSatisfied, err = d.applyPagingFilter(buffLocal, buffForeign, pp)
					modified = modified || m
					if err != nil {
						return
					}
				}

				// . Update buffers
				//
				// As we're fetching and providing sorted data, the buffer
				// merging algorithm may be a simple concatenation.
				d.updateBuffers(buffLocal, buffForeign)

				// If there were no modifications to what we pulled, we can safely
				// assume that we can produce a response.
				if !modified {
					break
				}
			}

			if d.bufferSize() == 0 {
				return d.prepareResponseGeneric(oLocalDef, dscr), nil
			}

			return d.prepareResponse(oLocalDef, oForeignDef, dscr, cap, processed)
		}, func() {
			if localLoader.closer != nil {
				localLoader.closer()
			}
			if foreignLoader.closer != nil {
				foreignLoader.closer()
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

// validatePagingSort assures that we can apply paging to the provided sort
//
// The final datasource that is defined in the given sort must be sorted by
// a primary key or a unique value.
func (d *joinedDataset) validatePagingSort(def *FrameDefinition, dd FrameDescriptionSet) (err error) {
	// The first one is always local so this is ok
	localDscr := dd[0]
	var dscr *FrameDescription

	// - determine the last ds of the given sort definition
	sortDS := ""
	auxSS := make(filter.SortExprSet, 0, len(def.Sort))
	for i := len(def.Sort) - 2; i >= 0; i-- {
		aa := strings.Split(def.Sort[i+1].Column, ".")
		bb := strings.Split(def.Sort[i].Column, ".")

		if len(aa) != len(bb) || (len(aa) > 1 && aa[0] != aa[1]) {
			auxSS = append(auxSS, def.Sort[i+1])
			if len(aa) > 1 {
				sortDS = aa[0]
			}
			break
		}

		auxSS = append(auxSS, def.Sort[i+1])
		if len(aa) > 1 {
			sortDS = aa[0]
		}
	}

	// when local, ref is omitted
	if sortDS == "" {
		// Do this to avoid extra work afterwords
		auxSS = def.Sort

		dscr = localDscr
		sortDS = dscr.Ref
	} else {
		dscr = dd.FilterByRef(sortDS)[0]
	}

	// - check if we're sorting by a unique value
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

// prepareSorting parses the provided local sort definition and applies it to
// the local and foreign definition
//
// This allows the store to provide the ordered data, simplifying our lives.
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

// localBasedLoad primary loads from local source then foreign one
func (d *joinedDataset) localBasedLoad(ctx context.Context, cap int, local, foreign *frameLoadCtx, localDef, foreignDef *FrameDefinition) (buffLocal, buffForeign []*Frame, err error) {
	// Fetch data from local ds
	var localPost func(ff []*Frame) ([]*Frame, error)
	if d.partitioned {
		localPost = func(ff []*Frame) ([]*Frame, error) {
			var err error

			ff, err = d.sliceFrames(ff, d.def.LocalColumn, d.def.ForeignColumn)
			for _, f := range ff {
				f.RefValue = ""
				f.RelColumn = ""
			}
			return ff, err
		}
	}

	buffLocal, err = d.frameLoader(local.loader, local, cap, d.partitioned, localPost)
	if err != nil {
		return
	}
	if len(buffLocal) == 0 {
		return
	}

	// Init local frame loader metadata
	if !local.metaInitialized {
		local.metaInitialized = true
		local.sortColumns = make([]int, len(local.sorting))
		for i, s := range local.sorting {
			local.sortColumns[i] = buffLocal[0].Columns.Find(s.Column)
		}

		local.keyColIndex = buffLocal[0].Columns.Find(local.keyCol)
	}

	// Prepare filtering for foreign DS
	keys, err := local.keys(buffLocal)
	if err != nil {
		return
	}
	if len(keys) == 0 {
		return
	}
	keyFilter := foreign.keyFilter(keys)

	// Fetch data from foreign ds
	foreignPost := func(ff []*Frame) ([]*Frame, error) {
		ff, err = d.sliceFrames(ff, d.def.ForeignColumn, d.def.LocalColumn)
		if err != nil {
			return nil, err
		}

		for i := range ff {
			if ff[i].Name == "" {
				ff[i].Name = foreignDef.Name
			}
			if ff[i].Source == "" {
				ff[i].Source = foreignDef.Source
			}
			if ff[i].Ref == "" {
				ff[i].Ref = foreignDef.Ref
			}
		}

		return ff, nil
	}

	fLdr, fClsr, err := foreign.initLoader(cap, keyFilter)
	if err != nil {
		return
	}
	defer fClsr()
	buffForeign, err = d.frameLoader(fLdr, foreign, 0, true, foreignPost)

	if len(buffForeign) == 0 {
		return
	}

	// Init foreign frame loader metadata
	if !foreign.metaInitialized {
		foreign.metaInitialized = true
		foreign.sortColumns = make([]int, len(foreign.sorting))
		for i, s := range foreign.sorting {
			foreign.sortColumns[i] = buffLocal[0].Columns.Find(s.Column)
		}

		foreign.keyColIndex = buffLocal[0].Columns.Find(foreign.keyCol)
	}

	return
}

// foreignBasedLoad primary loads from foreign source and then the local one
func (d *joinedDataset) foreignBasedLoad(ctx context.Context, cap int, local, foreign *frameLoadCtx, localDef, foreignDef *FrameDefinition) (buffLocal, buffForeign []*Frame, err error) {
	// Fetch data from foreign ds
	foreignPost := func(ff []*Frame) ([]*Frame, error) {
		ff, err = d.sliceFrames(ff, d.def.ForeignColumn, d.def.LocalColumn)
		if err != nil {
			return nil, err
		}

		for i := range ff {
			if ff[i].Name == "" {
				ff[i].Name = foreignDef.Name
			}
			if ff[i].Source == "" {
				ff[i].Source = foreignDef.Source
			}
			if ff[i].Ref == "" {
				ff[i].Ref = foreignDef.Ref
			}
		}

		return ff, nil
	}

	buffForeign, err = d.frameLoader(foreign.loader, foreign, 0, true, foreignPost)
	if err != nil {
		return
	}

	// Init foreign frame loader metadata
	if !foreign.metaInitialized {
		foreign.metaInitialized = true
		foreign.sortColumns = make([]int, len(foreign.sorting))
		for i, s := range foreign.sorting {
			foreign.sortColumns[i] = buffForeign[0].Columns.Find(s.Column)
		}

		foreign.keyColIndex = buffForeign[0].Columns.Find(foreign.keyCol)
	}

	// Prepare filtering for foreign DS
	keys, err := foreign.keys(buffForeign)
	if err != nil {
		return
	}
	if len(keys) == 0 {
		return
	}
	keyFilter := local.keyFilter(keys)

	// Fetch data from local ds
	var localPost func(ff []*Frame) ([]*Frame, error)
	if d.partitioned {
		localPost = func(ff []*Frame) ([]*Frame, error) {
			var err error

			ff, err = d.sliceFrames(ff, d.def.LocalColumn, d.def.ForeignColumn)
			for _, f := range ff {
				f.RefValue = ""
				f.RelColumn = ""
			}
			return ff, err
		}
	}

	lLdr, lClsr, err := local.initLoader(cap, keyFilter)
	if err != nil {
		return
	}
	defer lClsr()

	buffLocal, err = d.frameLoader(lLdr, local, 0, true, localPost)
	if err != nil {
		return
	}

	// Init local frame loader metadata
	if !local.metaInitialized {
		local.metaInitialized = true
		local.sortColumns = make([]int, len(local.sorting))
		for i, s := range local.sorting {
			local.sortColumns[i] = buffLocal[0].Columns.Find(s.Column)
		}

		local.keyColIndex = buffLocal[0].Columns.Find(local.keyCol)
	}

	return
}

// updateBuffers updates existing frame buffers with newly provided ones
func (d *joinedDataset) updateBuffers(local, foreign []*Frame) {
	// - local
	if len(local) != 0 {
		if len(d.localFrames) == 0 {
			d.localFrames = local
		} else {
			if d.partitioned {
				d.localFrames = append(d.localFrames, local...)
			} else {
				d.localFrames[0].Rows = append(d.localFrames[0].Rows, local[0].Rows...)
			}
		}
	}

	// - foreign
	if len(foreign) != 0 {
		if d.foreignSourceIndex == nil {
			d.foreignSourceIndex = make(map[string]int)
		}
		if ix, ok := d.foreignSourceIndex[foreign[0].Ref]; !ok {
			d.foreignSourceIndex[foreign[0].Ref] = len(d.foreignFrames)
			d.foreignFrames = append(d.foreignFrames, foreign)
		} else {
			d.foreignFrames[ix] = append(d.foreignFrames[ix], foreign...)
		}
	}
}

// removeOrphanRows removes all local rows with no related foreign frame
func (d *joinedDataset) removeOrphanRows(local, foreign []*Frame) ([]*Frame, []*Frame, bool) {
	// special edgecase where either of the buffers are empty
	if len(local) == 0 || len(foreign) == 0 {
		return nil, nil, true
	}

	// - index foreign frames based on ref value so that we have an easier time validating
	fIndex := make(map[string]bool)
	for _, ff := range foreign {
		fIndex[ff.RefValue] = true
	}

	// - go over local frame rows and see which ones need to be removed
	var k string
	var modified bool

	cols := local[0].Columns
	keyColIx := cols.Find(foreign[0].RelColumn)
	for _, lf := range local {
		aux := make([]FrameRow, 0, len(lf.Rows))

		lf.WalkRows(func(i int, r FrameRow) error {
			if r[keyColIx] == nil {
				return nil
			} else {
				k = cast.ToString(r[keyColIx].Get())
				if k == "" {
					return nil
				}
			}

			if fIndex[k] {
				aux = append(aux, r)
			}
			return nil
		})

		modified = modified || len(aux) != len(lf.Rows)
		lf.Rows = aux
	}

	return local, foreign, modified
}

func (d *joinedDataset) prepareResponse(localDef, foreignDef *FrameDefinition, dd FrameDescriptionSet, cap int, applyPaging bool) (oo []*Frame, err error) {
	// - Use the requested size of the two buffers
	// -- first come local frames
	oo = append(oo, d.cutLocalFrameBuffer(cap)...)
	localDelimiter := len(oo)
	//
	// -- followed by foreign
	aux := d.cutForeignFrameBuffer(cap)
	for _, a := range aux {
		a.RelSource = oo[0].Ref
	}
	oo = append(oo, aux...)

	// - Default for empty response
	if len(oo) == 0 {
		return d.prepareResponseGeneric(localDef, dd), nil
	}

	// - Optionally calculate paging
	if applyPaging {
		// when the buffer is empty, the isn't anything else to pull -- no next page
		if d.bufferSize() == 0 {
			return
		}

		// Paging can only be applied to non-partitioned local frames
		// @todo generalize?
		if localDelimiter != 1 {
			return
		}

		// @todo calculate prev page cursor; for now only next page is provided

		if oo[0].Paging == nil {
			oo[0].Paging = &filter.Paging{}
		}
		oo[0].Paging.NextPage = d.calculatePagingCursor(oo[0], oo[1:], true, localDef.Sort)
		oo[0].Paging.NextPage.LThen = localDef.Sort.Reversed()
	}

	return
}

// prepareResponseGeneric constructs an empty response in case where no data would be provided
func (d *joinedDataset) prepareResponseGeneric(lfd *FrameDefinition, dd FrameDescriptionSet) (oo []*Frame) {
	dscr := dd.FilterByRef(lfd.Ref)[0]

	f := &Frame{
		Name:    lfd.Name,
		Ref:     lfd.Ref,
		Source:  lfd.Source,
		Columns: make(FrameColumnSet, 0, len(lfd.Columns)),
	}

	if len(lfd.Columns) == 0 {
		f.Columns = dscr.Columns
	} else {
		for _, c := range lfd.Columns {
			f.Columns = append(f.Columns, dscr.Columns[dscr.Columns.Find(c.Name)])
		}
	}
	return []*Frame{f}
}

// calculatePagingCursor calculates the paging cursor for the given frame
//
// The check for existence should be performed way in advanced so we won't bother here.
// A unique value is also assured at way before.
func (d *joinedDataset) calculatePagingCursor(local *Frame, foreign []*Frame, desc bool, ss filter.SortExprSet) *filter.PagingCursor {
	cursor := &filter.PagingCursor{LThen: ss.Reversed()}
	var foreignFrames []*Frame
	var v interface{}
	var localRow FrameRow

	if desc {
		localRow = local.LastRow()
	} else {
		localRow = local.FirstRow()
	}

	// Index foreign frames by ref for nicer lookups
	foreignIndex := make(map[string][]*Frame)
	for _, f := range foreign {
		foreignIndex[f.Ref] = append(foreignIndex[f.Ref], f)
	}

	// Go over sorting and construct the cursor
	for _, s := range ss {
		foreignFrames = nil

		if strings.Contains(s.Column, ".") {
			pts := strings.Split(s.Column, ".")
			foreignFrames = foreignIndex[pts[0]]
			var r FrameRow
			var f *Frame
			if desc {
				f = foreignFrames[len(foreignFrames)-1]
				r = f.FirstRow()
			} else {
				f = foreignFrames[0]
				r = f.FirstRow()
			}

			if r[f.Columns.Find(pts[1])] != nil {
				v = r[f.Columns.Find(pts[1])].Get()
			}
			cursor.Set(s.Column, v, s.Descending)
		} else {
			if localRow[local.Columns.Find(s.Column)] != nil {
				v = localRow[local.Columns.Find(s.Column)].Get()
			}
			cursor.Set(s.Column, v, s.Descending)
		}
	}

	return cursor
}

func (d *joinedDataset) cutLocalFrameBuffer(cap int) (oo []*Frame) {
	// When partitioned, extract frames
	if d.partitioned {
		// We can use everything
		if len(d.localFrames) <= cap {
			oo = d.localFrames
			d.localFrames = nil
			return
		}

		oo = d.localFrames[0:cap]
		d.localFrames = d.localFrames[cap:]
		return
	}

	// When not partitioned, extract rows
	aux := d.localFrames[0].CloneMeta()
	if d.localFrames[0].Size() <= cap {
		aux.Rows = d.localFrames[0].Rows
		d.localFrames = nil
		return []*Frame{aux}
	}

	aux.Rows = d.localFrames[0].Rows[0:cap]
	d.localFrames[0].Rows = d.localFrames[0].Rows[cap:]
	return []*Frame{aux}
}

func (d *joinedDataset) cutForeignFrameBuffer(cap int) (oo []*Frame) {
	for i, dsFrames := range d.foreignFrames {
		// We can use everything
		if len(dsFrames) <= cap {
			oo = append(oo, dsFrames...)
			d.foreignFrames[i] = nil
			continue
		}

		oo = append(oo, dsFrames[0:cap]...)
		d.foreignFrames[i] = dsFrames[cap:]
	}

	return
}

func (d *joinedDataset) bufferSize() int {
	if len(d.localFrames) == 0 {
		return 0
	}

	if d.partitioned {
		return len(d.localFrames)
	}

	return d.localFrames[0].Size()
}

func (d *joinedDataset) sizeOfBuffer(b []*Frame, partitioned bool) int {
	if partitioned {
		return len(b)
	}

	if len(b) == 0 {
		return 0
	}
	return b[0].Size()
}
