package report

import (
	"context"
	"errors"
	"fmt"

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

// @todo allow x-join sorting
//       - determine the lead sort datasource (first sort expr. definition); use as base
//       - sort the output based on the lead datasource
//
// @todo allow x-join filtering
//
// @todo improve join datasource loading
//       use SQL partitioning & grouping to determine chunks that fall into the same group.
func (d *joinedDataset) Load(ctx context.Context, dd ...*FrameDefinition) (Loader, Closer, error) {
	// to hold closer references for all underlying datasources
	closers := make([]Closer, 0, 10)

	return func(cap int) ([]*Frame, error) {
			out := make([]*Frame, 0, 10)

			// fetch base frame
			baseDef := FrameDefinitionSet(dd).FindBySourceRef(d.Name(), d.def.localDim())
			if baseDef == nil {
				return nil, fmt.Errorf("could not find base definition: %s, %s", d.Name(), d.def.localDim())
			}
			baseL, baseC, err := d.base.Load(ctx, baseDef)
			if err != nil {
				return nil, err
			}
			closers = append(closers, baseC)
			baseFrame, err := baseL(cap)
			if err != nil {
				return nil, err
			}
			for i := range baseFrame {
				baseFrame[i].Name = baseDef.Name
				baseFrame[i].Source = baseDef.Source
				baseFrame[i].Ref = baseDef.Ref
			}
			out = append(out, baseFrame...)

			// fetch foreign frames
			// - foreign ref
			foreignDef := FrameDefinitionSet(dd).FindBySourceRef(d.Name(), d.def.foreignDim())
			if foreignDef == nil {
				return nil, fmt.Errorf("could not find foreign definition: %s", d.foreign.Name())
			}

			// - extract keys
			kk := make([]string, 0, baseFrame[0].Size())
			kx := make(map[string]bool)
			keyCol := baseFrame[0].Columns.Find(d.def.localColumn())
			if keyCol < 0 {
				return nil, fmt.Errorf("could not find local key column: %s", d.def.localColumn())
			}
			var ok bool
			err = baseFrame[0].WalkRows(func(i int, r FrameRow) (err error) {
				c := r[keyCol].Get()
				k, err := cast.ToStringE(c)
				if err != nil {
					return err
				}
				if ok = kx[k]; !ok {
					kk = append(kk, k)
					kx[k] = true
				}
				return nil
			})

			// @todo partitioning
			// @todo parallel
			fdr := foreignDef.Rows
			for _, k := range kk {
				f := (&RowDefinition{
					Cells: map[string]*CellDefinition{
						d.def.foreignColumn(): {Op: "eq", Value: "'" + k + "'"},
					},
				}).MergeAnd(fdr)

				foreignDef.Rows = f
				foreignL, foreignC, err := d.foreign.Load(ctx, foreignDef)
				if err != nil {
					return nil, err
				}
				closers = append(closers, foreignC)
				foreignFrame, err := foreignL(cap)
				if err != nil {
					return nil, err
				}
				for i := range foreignFrame {
					foreignFrame[i].Name = foreignDef.Name
					foreignFrame[i].Source = foreignDef.Source
					foreignFrame[i].Ref = foreignDef.Ref
				}

				out = append(out, foreignFrame...)
			}

			return out, nil
		}, func() {
			for _, c := range closers {
				c()
			}
		}, nil
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
