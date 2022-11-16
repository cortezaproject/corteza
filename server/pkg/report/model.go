package report

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	model struct {
		ran   bool
		steps []step
		nodes map[string]*modelGraphNode
	}

	// M is the model interface that should be used when trying to model the datasource
	M interface {
		Add(...step) M
		Run(context.Context) error
		Load(context.Context, ...*FrameDefinition) ([]*Frame, error)
		Describe(ctx context.Context, source string) (FrameDescriptionSet, error)
		GetStep(name string) step
	}

	stepSet []step
	step    interface {
		Name() string
		Source() []string
		Run(context.Context, ...Datasource) (Datasource, error)
		Validate() error
		Def() *StepDefinition
	}

	StepDefinitionSet []*StepDefinition
	StepDefinition    struct {
		Kind string `json:"kind,omitempty"`

		Load  *LoadStepDefinition  `json:"load,omitempty"`
		Join  *JoinStepDefinition  `json:"join,omitempty"`
		Group *GroupStepDefinition `json:"group,omitempty"`
		// @todo Transform
	}
)

// Model initializes the model based on the provided sources and step definitions.
//
// Additional steps may be added after the model is constructed.
// Call `M.Run(context.Context)` to allow the model to be used for requesting data.
// Additional steps may not be added after the `M.Run(context.Context)` was called
func Model(ctx context.Context, sources map[string]DatasourceProvider, dd ...*StepDefinition) (M, error) {
	steps := make([]step, 0, len(dd))

	err := func() error {
		for _, d := range dd {
			switch {
			case d.Load != nil:
				if sources == nil {
					return fmt.Errorf("no datasource providers defined")
				}
				steps = append(steps, &stepLoad{def: d.Load, dsp: sources[d.Load.Source]})

			case d.Join != nil:
				steps = append(steps, &stepJoin{def: d.Join})

			case d.Group != nil:
				steps = append(steps, &stepGroup{def: d.Group})

			// @todo Transform

			default:
				return fmt.Errorf("malformed step definition: unsupported step kind")
			}
		}
		return nil
	}()

	if err != nil {
		return nil, fmt.Errorf("failed to create the model: %s", err.Error())
	}

	return &model{
		steps: steps,
	}, nil
}

func (ss StepDefinitionSet) Validate() error {
	var f *Filter
	for _, s := range ss {
		switch {
		case s.Load != nil:
			f = s.Load.Filter
		case s.Join != nil:
			f = s.Join.Filter
		case s.Group != nil:
			f = s.Group.Filter
		}
		if f != nil && f.Error != "" {
			return errors.InvalidData(f.Error)
		}
	}
	return nil
}

// Add adds additional steps to the model
func (m *model) Add(ss ...step) M {
	m.steps = append(m.steps, ss...)
	return m
}

// Run bakes the model configuration and makes the requested data available
func (m *model) Run(ctx context.Context) (err error) {
	const errPfx = "failed to run the model"
	defer func() {
		m.ran = true
	}()

	// initial validation
	err = func() (err error) {
		if m.ran {
			return fmt.Errorf("model already ran")
		}

		if len(m.steps) == 0 {
			return fmt.Errorf("no model steps defined")
		}

		return nil
	}()
	if err != nil {
		return fmt.Errorf("%s: failed to validate the model: %w", errPfx, err)
	}

	// construct a model graph for future optimizations
	m.nodes, err = m.buildStepGraph(m.steps)
	if err != nil {
		return fmt.Errorf("%s: %w", errPfx, err)
	}

	return nil
}

// Describe returns the descriptions for the requested model datasources
//
// The Run method must be called before the description can be provided.
func (m *model) Describe(ctx context.Context, source string) (out FrameDescriptionSet, err error) {
	var ds Datasource

	err = func() error {
		if !m.ran {
			return fmt.Errorf("model was not yet ran")
		}

		ds, err = m.datasource(ctx, &FrameDefinition{Source: source})
		if err != nil {
			return fmt.Errorf("model does not contain the datasource %s", source)
		}

		return nil
	}()
	if err != nil {
		return nil, fmt.Errorf("unable to describe the model source: %w", err)
	}

	return ds.Describe(), nil
}

// GetStep returns the details of the requested step
func (m *model) GetStep(name string) step {
	for _, s := range m.steps {
		if s.Name() == name {
			return s
		}
	}

	return nil
}

// Load returns the Frames based on the provided FrameDefinitions
//
// The Run method must be called before the frames can be provided.
func (m *model) Load(ctx context.Context, dd ...*FrameDefinition) (ff []*Frame, err error) {
	var (
		def *FrameDefinition
		ds  Datasource
	)

	// request validation
	err = func() error {
		// - all frame definitions must define the same datasource; call Load multiple times if
		//   you need to access multiple datasources
		for i, d := range dd {
			if i == 0 {
				continue
			}
			if d.Source != dd[i-1].Source {
				return fmt.Errorf("frame definition source missmatch: expected %s, got %s", dd[i-1].Source, d.Source)
			}
		}

		def = dd[0]
		ds, err = m.datasource(ctx, def)
		if err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		return nil, fmt.Errorf("unable to load frames: invalid request: %w", err)
	}

	// apply any frame definition defaults
	aux := make([]*FrameDefinition, len(dd))
	for i, d := range dd {
		aux[i] = d.Clone()

		// assure paging is always provided so we can ignore nil checks
		if aux[i].Paging == nil {
			aux[i].Paging = &filter.Paging{
				Limit: defaultPageSize,
			}
		}

		// assure sorting is always provided so we can ignore nil checks
		if aux[i].Sort == nil {
			aux[i].Sort = filter.SortExprSet{}
		}
	}
	dd = aux

	// assure paging is always provided so we can ignore nil checks
	if def.Paging == nil {
		def.Paging = &filter.Paging{
			Limit: defaultPageSize,
		}
	}

	// assure sorting is always provided so we can ignore nil checks
	if def.Sort == nil {
		def.Sort = filter.SortExprSet{}
	}

	// load the data
	err = func() error {
		l, c, err := ds.Load(ctx, dd...)
		if err != nil {
			return err
		}
		defer c()

		ff, err = l(int(def.Paging.Limit), true)
		if err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		return nil, fmt.Errorf("unable to load frames: %w", err)
	}

	return ff, nil
}
