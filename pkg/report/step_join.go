package report

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/qlng"
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
		out = append(out,
			&FrameDescription{
				Source:  d.Name(),
				Ref:     l.Source,
				Columns: l.Columns,
			},
		)
	}

	foreign := d.foreign.Describe()
	for _, f := range foreign {
		out = append(out,
			&FrameDescription{
				Source:  d.Name(),
				Ref:     f.Source,
				Columns: f.Columns,
			},
		)
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
	// Preparation
	// - Assure sort columns for paging purposes
	// - Assure local/foreign definitions
	//   Keep a cloned original version so we don't overwrite the initial definition.
	oLocalDef, oForeignDef, err := d.prepareDefinitions(FrameDefinitionSet(dd))
	if err != nil {
		return
	}

	// - Validate the sort of the local frame for paging purposes
	dscr := d.Describe()
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
	var cndMain, apxx, cndSub *qlng.ASTNode
	if oLocalDef.Paging.PageCursor != nil {
		cndMain, apxx, cndSub, err = d.calculatePagingFilters(localDef, inverted)

		localDef.Filter = merger(&Filter{cndMain}, localDef.Filter, "and")
		localDef.Paging.PageCursor = nil
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
	return func(cap int, paged bool) (oo []*Frame, err error) {
			var keys []string
			if isEmpty {
				return
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
				more, err = mainLoader.load(nil, paged)
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
				more, err = subLoader.load(keyFilter, paged)
				if err != nil {
					return nil, err
				}
				if !more {
					break
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
					m, pagingSatisfied = d.pagingFilter(mainLoader, subLoader, apxx, cndSub)
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

			return prepareResponse(mainLoader, subLoader, inverted, oLocalDef, d.def.LocalColumn, dscr)
		}, func() {
			return
		}, nil
}

func (d *joinedDataset) prepareDefinitions(dd FrameDefinitionSet) (localDef *FrameDefinition, foreignDef *FrameDefinition, err error) {
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

// pagingFilter applies additional filtering based on the given page cursor
func (d *joinedDataset) pagingFilter(main, sub *frameBuffer, cndMain, cndSub *qlng.ASTNode) (modified, satisfied bool) {
	cutSize := 0
	done := false
	if cndMain == nil {
		return false, true
	}
	main.walkRowsLocal(func(i int, r FrameRow) error {
		if done {
			return nil
		}

		// Firstly we evaluate if the local row falls in the "danger zone"
		// (if the row was right on the edge of where the paging cursor filter applied)
		if d.eval(cndMain, r, main.localFrames[0].Columns) {

			// If we are in the "danger zone", we check what foreign frames don't pass
			// the cursor filter.
			//
			// If the foreign frame does not pass it, we should remove it along with the local row.
			if cndSub != nil && !d.eval(cndSub, sub.getByRefValue(r[main.keyColIndex]).FirstRow(), sub.localFrames[0].Columns) {
				cutSize++
			} else {
				done = true
				return nil
			}
		} else {
			done = true
			return nil
		}

		return nil
	})

	if cutSize > 0 {
		main.removeLocal(cutSize)
		sub.removeLocal(cutSize)

		if main.sizeLocal() <= cutSize {
			// We removed all of the local buffer so the paging is not yet satisfied
			return true, false
		}

		// We removed the portion of the local buffer, so the paging is satisfied
		return true, true
	}

	return false, true
}

// calculatePagingFilters produces additional filtering that should be done
// on the datasource level and/or in the join logic.
//
// The core logic is extracted from store/rdbms/builders/cursor.go
func (d *joinedDataset) calculatePagingFilters(local *FrameDefinition, inverted bool) (localCondition, localAppendix, foreignCondition *qlng.ASTNode, err error) {
	if len(local.Paging.PageCursor.Keys()) == 0 {
		return
	}

	var (
		cur = local.Paging.PageCursor

		// baseCndAppx is the initial AST for finding rows that match the sort column
		// It's basically the second part of the wrap condition (if the value equals)
		//
		// The correlated string version is: (%s OR ((%s IS NULL AND %s) OR %s = %s))
		baseCndAppx = func(field string, checkNull bool, value interface{}) *qlng.ASTNode {
			return &qlng.ASTNode{
				Ref: "or",
				Args: qlng.ASTNodeSet{
					&qlng.ASTNode{
						Ref: "and",
						Args: qlng.ASTNodeSet{
							&qlng.ASTNode{
								Ref: "is",
								Args: qlng.ASTNodeSet{
									&qlng.ASTNode{
										Symbol: field,
									},
									&qlng.ASTNode{
										Ref: "null",
									},
								},
							}, &qlng.ASTNode{
								Value: qlng.MakeValueOf("Boolean", checkNull),
							},
						},
					},
					&qlng.ASTNode{
						Ref: "eq",
						Args: qlng.ASTNodeSet{{
							Symbol: field,
						}, {
							// @todo type
							Value: qlng.MakeValueOf("String", value),
						}},
					},
				},
			}
		}

		// baseCnd is the initial AST for filtering over the given sort column
		//
		// The correlated string version is: ((%s IS %s AND %s) OR (%s %s %s))
		baseCnd = func(field string, nullVal *qlng.ASTNode, checkNull bool, compOp string, value interface{}, appendix bool) *qlng.ASTNode {
			pp := strings.Split(field, ".")
			field = pp[len(pp)-1]

			out := &qlng.ASTNode{
				Ref: "or",
				Args: qlng.ASTNodeSet{&qlng.ASTNode{
					Ref: "and",
					Args: qlng.ASTNodeSet{&qlng.ASTNode{
						Ref: "is",
						Args: qlng.ASTNodeSet{{
							Symbol: field,
						}, nullVal},
					}, &qlng.ASTNode{
						Value: qlng.MakeValueOf("Boolean", checkNull),
					}},
				}, &qlng.ASTNode{
					Ref: compOp,
					Args: qlng.ASTNodeSet{{
						Symbol: field,
					}, {
						// @todo type
						Value: qlng.MakeValueOf("String", value),
					}},
				},
				},
			}

			if appendix {
				localAppendix = baseCndAppx(field, checkNull, value)

				return &qlng.ASTNode{
					Ref: "or",
					Args: qlng.ASTNodeSet{
						out,
						localAppendix,
					},
				}
			}

			return out
		}

		// wrapCnd is the conjunction between two paging cursor columns
		//
		// The correlated string version is: (%s OR (((%s IS NULL AND %s) OR %s = %s) AND %s))
		wrapCnd = func(base *qlng.ASTNode, field string, value interface{}, checkNull bool, condition *qlng.ASTNode) *qlng.ASTNode {
			pp := strings.Split(field, ".")
			field = pp[len(pp)-1]

			return &qlng.ASTNode{
				Ref: "or",
				Args: qlng.ASTNodeSet{
					base,
					&qlng.ASTNode{
						Ref: "and",
						Args: qlng.ASTNodeSet{
							&qlng.ASTNode{
								Ref: "or",
								Args: qlng.ASTNodeSet{
									&qlng.ASTNode{
										Ref: "and",
										Args: qlng.ASTNodeSet{
											&qlng.ASTNode{
												Ref: "is",
												Args: qlng.ASTNodeSet{
													&qlng.ASTNode{
														Symbol: field,
													},
													&qlng.ASTNode{
														Ref: "null",
													},
												},
											}, &qlng.ASTNode{
												Value: qlng.MakeValueOf("Boolean", checkNull),
											},
										},
									},
									&qlng.ASTNode{
										Ref: "eq",
										Args: qlng.ASTNodeSet{
											&qlng.ASTNode{
												Symbol: field,
											},
											&qlng.ASTNode{
												// @todo type
												Value: qlng.MakeValueOf("String", value),
											},
										},
									},
								},
							},
							condition,
						},
					},
				},
			}
		}
	)

	var (
		cc = cur.Keys()
		vv = cur.Values()

		ltOp = map[bool]string{
			true:  "lt",
			false: "gt",
		}

		notOp = map[bool]*qlng.ASTNode{
			true:  {Ref: "nnull"},
			false: {Ref: "null"},
		}

		isNull = func(i int, neg bool) bool {
			if (isNil(vv[i]) && !neg) || (!isNil(vv[i]) && neg) {
				return true
			}

			return false
		}
	)

	// Determine the point at which we switch local sorts and foreign sorts
	sourceDelimiter := len(cc) - 1
	for j := range cc {
		if j > 0 {
			if strings.Contains(cc[j-1], ".") != strings.Contains(cc[j], ".") {
				sourceDelimiter = j
				break
			}
		}
	}

	// Some temporary variables to avoid initialization
	var tmp []string
	var field string

	calculateAST := func(cc []string, vv []interface{}, dsc []bool, cut bool) (cnd *qlng.ASTNode) {
		// going from the last key/column to the 1st one
		for i := len(cc) - 1; i >= 0; i-- {
			// We need to cut off the values that are before the cursor (when ascending)
			// and vice-versa for descending.
			lt := dsc[i]
			if cut && cur.IsROrder() {
				lt = !lt
			}
			op := ltOp[lt]

			tmp = strings.Split(cc[i], ".")
			field = tmp[len(tmp)-1]

			base := baseCnd(field, notOp[!lt], isNull(i, lt), op, vv[i], cut && i == len(cc)-1)

			if cnd == nil {
				cnd = base
			} else {
				cnd = wrapCnd(base, field, vv[i], isNull(i, false), cnd)
			}
		}

		return
	}

	// when there is no delimiter we can fully filter the ds
	if sourceDelimiter == len(cc)-1 {
		localCondition = calculateAST(cc, vv, cur.Desc(), false)
	} else {
		localCondition = calculateAST(cc[0:sourceDelimiter], vv[0:sourceDelimiter], cur.Desc()[0:sourceDelimiter], true)
		foreignCondition = calculateAST(cc[sourceDelimiter:], vv[sourceDelimiter:], cur.Desc()[sourceDelimiter:], false)
	}

	return
}

func (d *joinedDataset) validateSort(def *FrameDefinition, dd FrameDescriptionSet) (err error) {
	sortDS := ""
	auxSS := make(filter.SortExprSet, 0, len(def.Sort))

	// Get the last sorting delimiter
	for i := len(def.Sort) - 1; i >= 0; i-- {
		s := def.Sort[i]

		spts := strings.Split(s.Column, ".")
		if len(spts) == 1 && sortDS != "" {
			break
		}
		if len(spts) > 1 {
			if sortDS == "" {
				sortDS = spts[0]
			} else if sortDS != spts[0] {
				break
			}
		}
		auxSS = append(auxSS, s)
	}

	// Check if we're sorting by a unique value
	if sortDS == "" {
		sortDS = def.Ref
	}

	dscr := dd.FilterByRef(sortDS)[0]
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
		if sortDS == def.Ref {
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

	if foreignDS != "" && foreignDS != foreign.Ref {
		return false, fmt.Errorf("foreign datasource in sort expression not found: %s", foreignDS)
	}

	local.Sort = localSS
	foreign.Sort = append(foreignSS, foreign.Sort...)

	return
}
