package wfexec

import (
	"context"
)

type (
	Steps []Step
	Step  interface {
		Exec(context.Context, *ExecRequest) (ExecResponse, error)
	}

	// list of Graph steps with relations
	Graph struct {
		steps    []Step
		children map[Step][]Step
		parents  map[Step][]Step
		index    map[uint64]Step
	}
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

func (g *Graph) AddStep(s Step, cc ...Step) {
	g.steps = append(g.steps, s)

	if len(cc) > 0 {
		g.children[s] = cc
		for _, c := range cc {
			g.AddParent(c, s)
		}
	}
}

func (g *Graph) Len() int {
	return len(g.steps)
}

func (g *Graph) SetStepIdentifier(s Step, ID uint64) {
	g.index[ID] = s
}

func (g *Graph) GetStepByIdentifier(ID uint64) Step {
	return g.index[ID]
}

func (g *Graph) AddParent(c, p Step) {
	g.parents[c] = append(g.parents[c], p)
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
