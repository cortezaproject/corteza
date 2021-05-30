package report

import (
	"context"
	"errors"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/spf13/cast"
)

type (
	model struct {
		steps       []Step
		datasources DatasourceSet
	}

	M interface {
		Add(...Step) M
		Run(context.Context) error
		Load(context.Context, ...*FrameDefinition) ([]*Frame, error)
	}

	StepSet []Step
	Step    interface {
		Name() string
		Source() []string
		Run(context.Context, ...Datasource) (Datasource, error)
		Validate() error
		Def() *StepDefinition
	}

	StepDefinitionSet []*StepDefinition
	StepDefinition    struct {
		Load  *LoadStepDefinition  `json:"load,omitempty"`
		Join  *JoinStepDefinition  `json:"join,omitempty"`
		Group *GroupStepDefinition `json:"group,omitempty"`
		// @todo Transform
	}

	modelGraphNode struct {
		step Step
		ds   Datasource

		pp []*modelGraphNode
		cc []*modelGraphNode
	}
)

func Model(ctx context.Context, sources map[string]DatasourceProvider, dd ...*StepDefinition) (M, error) {
	steps := make([]Step, 0, len(dd))
	ss := make(DatasourceSet, 0, len(steps)*2)

	err := func() error {
		for _, d := range dd {
			switch {
			case d.Load != nil:
				if sources == nil {
					return errors.New("no datasources defined")
				}

				s, ok := sources[d.Load.Source]
				if !ok {
					return fmt.Errorf("unresolved data source: %s", d.Load.Source)
				}
				ds, err := s.Datasource(ctx, d.Load)
				if err != nil {
					return err
				}

				ss = append(ss, ds)

			case d.Join != nil:
				steps = append(steps, &stepJoin{def: d.Join})

			case d.Group != nil:
				steps = append(steps, &stepGroup{def: d.Group})

			// @todo Transform

			default:
				return errors.New("malformed step definition")
			}
		}
		return nil
	}()

	if err != nil {
		return nil, fmt.Errorf("failed to create the model: %s", err.Error())
	}

	return &model{
		steps:       steps,
		datasources: ss,
	}, nil
}

func (m *model) Add(ss ...Step) M {
	m.steps = append(m.steps, ss...)
	return m
}

func (m *model) Run(ctx context.Context) (err error) {
	// initial validation
	err = m.validateModel()
	if err != nil {
		return fmt.Errorf("failed to validate the model: %w", err)
	}

	// nothing left to do
	if len(m.steps) == 0 {
		return nil
	}

	// construct the step graph
	gg, err := m.buildStepGraph(m.steps, m.datasources)
	if err != nil {
		return err
	}

	m.datasources = nil
	for _, n := range gg {
		aux, err := m.reduceGraph(ctx, n)
		if err != nil {
			return err
		}
		m.datasources = append(m.datasources, aux)
	}

	return nil
}

func (m *model) Load(ctx context.Context, dd ...*FrameDefinition) ([]*Frame, error) {
	var err error

	for _, d := range dd {
		err = m.applyPaging(d, d.Paging, d.Sorting)
		if err != nil {
			return nil, err
		}
	}

	// @todo variable root def
	def := dd[0]

	ds := m.datasources.Find(def.Source)
	if ds == nil {
		return nil, fmt.Errorf("unresolved source: %s", def.Source)
	}

	l, c, err := ds.Load(ctx, dd...)
	if err != nil {
		return nil, err
	}
	defer c()

	i := 0
	if def.Paging != nil && def.Paging.Limit > 0 {
		i = int(def.Paging.Limit)
	}

	ff, err := l(i + 1)
	if err != nil {
		return nil, err
	}

	dds := FrameDefinitionSet(dd)
	for i, f := range ff {
		def = dds.FindBySourceRef(f.Source, f.Ref)
		if def == nil {
			return nil, fmt.Errorf("unable to find frame definition for frame: src-%s, ref-%s", f.Source, f.Ref)
		}
		ff[i], err = m.calculatePaging(f, def.Paging, def.Sorting)
		if err != nil {
			return nil, err
		}
	}

	return ff, err
}

func (m *model) calculatePaging(f *Frame, p *filter.Paging, ss filter.SortExprSet) (*Frame, error) {
	if p == nil {
		p = &filter.Paging{}
	}

	var (
		hasPrev = p.PageCursor != nil
		hasNext = f.Size() > int(p.Limit)
		out     = &filter.Paging{}
	)

	out.Limit = p.Limit

	if hasNext {
		f, _ = f.Slice(0, f.Size()-2)
		out.NextPage = m.calculatePageCursor(f.LastRow(), f.Columns, ss)
	}

	if hasPrev {
		out.PrevPage = m.calculatePageCursor(f.FirstRow(), f.Columns, ss)
	}

	f.Paging = out
	f.Sorting = &filter.Sorting{
		Sort: ss,
	}

	return f, nil
}

func (m *model) calculatePageCursor(r FrameRow, cc FrameColumnSet, ss filter.SortExprSet) *filter.PagingCursor {
	out := &filter.PagingCursor{LThen: ss.Reversed()}

	for _, s := range ss {
		ci := cc.Find(s.Column)
		out.Set(s.Column, r[ci].Get(), s.Descending)
	}

	return out
}

func (m *model) applyPaging(def *FrameDefinition, p *filter.Paging, ss filter.SortExprSet) (err error) {
	if p == nil {
		return nil
	}

	ss, err = p.PageCursor.Sort(ss)
	if err != nil {
		return err
	}

	// @todo somesort of a primary key to avoid edgecases
	sort := ss.Clone()
	if p.PageCursor != nil && p.PageCursor.ROrder {
		sort.Reverse()
	}
	def.Sorting = sort

	// convert cursor to rows def
	if p.PageCursor == nil {
		return nil
	}

	rd := &RowDefinition{
		Cells: make(map[string]*CellDefinition),
	}
	kk := p.PageCursor.Keys()
	vv := p.PageCursor.Values()
	for i, k := range kk {
		v, err := cast.ToStringE(vv[i])
		if err != nil {
			return err
		}

		lt := p.PageCursor.Desc()[i]
		if p.PageCursor.IsROrder() {
			lt = !lt
		}
		op := ""
		if lt {
			op = "lt"
		} else {
			op = "gt"
		}

		rd.Cells[k] = &CellDefinition{
			Op:    op,
			Value: fmt.Sprintf("'%s'", v),
		}
	}
	def.Rows = rd.MergeAnd(def.Rows)

	return nil
}

func (m *model) validateModel() error {
	if len(m.steps)+len(m.datasources) == 0 {
		return errors.New("no model steps defined")
	}

	var err error
	for _, s := range m.steps {
		err = s.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *model) buildStepGraph(ss StepSet, dd DatasourceSet) ([]*modelGraphNode, error) {
	mp := make(map[string]*modelGraphNode)

	for _, s := range ss {
		s := s

		// make sure that the step is in the graph
		n, ok := mp[s.Name()]
		if !ok {
			n = &modelGraphNode{
				step: s,
			}
			mp[s.Name()] = n
		} else {
			n.step = s
		}

		// make sure the child step is in there
		for _, src := range s.Source() {
			c, ok := mp[src]
			if !ok {
				c = &modelGraphNode{
					// will be added later
					step: nil,
					pp:   []*modelGraphNode{n},

					ds: dd.Find(src),
				}
				mp[src] = c
			}
			n.cc = append(n.cc, c)
		}
	}

	// return all of the root nodes
	out := make([]*modelGraphNode, 0, len(ss))
	for _, n := range mp {
		if len(n.pp) == 0 {
			out = append(out, n)
		}
	}

	return out, nil
}

func (m *model) reduceGraph(ctx context.Context, n *modelGraphNode) (out Datasource, err error) {
	auxO := make([]Datasource, len(n.cc))
	if len(n.cc) > 0 {
		for i, c := range n.cc {
			out, err = m.reduceGraph(ctx, c)
			if err != nil {
				return nil, err
			}
			auxO[i] = out
		}
	}

	bail := func() (out Datasource, err error) {
		if n.step == nil {
			if n.ds != nil {
				return n.ds, nil
			}

			return out, nil
		}

		aux, err := n.step.Run(ctx, auxO...)
		if err != nil {
			return nil, err
		}

		return aux, nil
	}

	if n.step == nil {
		return bail()
	}

	// check if this one can reduce the existing datasources
	//
	// for now, only "simple branches are supported"
	var o Datasource
	if len(auxO) > 1 {
		return bail()
	} else if len(auxO) > 0 {
		// use the only available output
		o = auxO[0]
	} else {
		// use own datasource (in case of leaves nodes)
		o = n.ds
	}

	if n.step.Def().Group != nil {
		gds, ok := o.(GroupableDatasource)
		if !ok {
			return bail()
		}

		ok, err = gds.Group(n.step.Def().Group.GroupDefinition, n.step.Name())
		if err != nil {
			return nil, err
		} else if !ok {
			return bail()
		}

		out = gds
		// we've covered this step with the child step; ignore it
		return out, nil
	}

	return bail()
}

// @todo cleanup the bellow two?

func (sd *StepDefinition) source() string {
	switch {
	case sd.Load != nil:
		return sd.Load.Source
	case sd.Group != nil:
		return sd.Group.Source
	// @todo Transform
	default:
		return ""
	}
}

func (sd *StepDefinition) name() string {
	switch {
	case sd.Load != nil:
		return sd.Load.Name
	case sd.Group != nil:
		return sd.Group.Name
	// @todo Transform
	default:
		return ""
	}
}
