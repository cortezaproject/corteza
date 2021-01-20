package wfexec

import (
	"context"
)

type (
	Steps []Step
	Step  interface {
		ID() uint64
		SetID(uint64)
		Exec(context.Context, *ExecRequest) (ExecResponse, error)
	}

	// list of Graph steps with relations
	Graph struct {
		steps    []Step
		children map[Step][]Step
		parents  map[Step][]Step
		index    map[uint64]Step
	}

	stepIdentifier struct{ id uint64 }
)

func NewGraph() *Graph {
	wf := &Graph{
		steps:    make([]Step, 0, 1024),
		children: make(map[Step][]Step),
		parents:  make(map[Step][]Step),
		index:    make(map[uint64]Step),
	}

	return wf
}

func (i *stepIdentifier) ID() uint64      { return i.id }
func (i *stepIdentifier) SetID(id uint64) { i.id = id }

func (g *Graph) AddStep(s Step, cc ...Step) {
	g.steps = append(g.steps, s)

	if id := s.ID(); id != 0 {
		g.index[id] = s
	}

	if len(cc) > 0 {
		for _, c := range cc {
			g.AddParent(c, s)
		}
	}
}

func (g *Graph) Len() int {
	return len(g.steps)
}

func (g *Graph) StepByID(ID uint64) Step {
	return g.index[ID]
}

func (g *Graph) AddParent(c, p Step) {
	g.parents[c] = append(g.parents[c], p)
	g.children[p] = append(g.children[p], c)
}

func (g *Graph) Children(s Step) Steps {
	return g.children[s]
}

func (g *Graph) Parents(s Step) Steps {
	return g.parents[s]
}

func (g *Graph) Exec(context.Context, *ExecRequest) (ExecResponse, error) {
	// @todo
	return nil, nil
}

func (g *Graph) Orphans() (oo Steps) {
	for _, step := range g.steps {
		if len(g.Parents(step)) > 0 {
			continue
		}

		oo = append(oo, step)
	}

	return
}
