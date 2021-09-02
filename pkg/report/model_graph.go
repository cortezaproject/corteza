package report

import (
	"context"
	"fmt"
)

type (
	modelGraphNode struct {
		step step
		pp   []*modelGraphNode
		cc   []*modelGraphNode
	}
)

// Internally, the model uses a stepGraph to resolve the dependencies between the steps
// allowing us to perform some preprocessing, such as size reduction and shape validation
func (m *model) buildStepGraph(ss stepSet) (map[string]*modelGraphNode, error) {
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
			return nil, fmt.Errorf("step name not unique: %s", s.Name())
		}

		// make sure the child step is in there
		for _, src := range s.Source() {
			c, ok := mp[src]
			if !ok {
				c = &modelGraphNode{
					// step is added later when we get to it
					step: nil,
					pp:   []*modelGraphNode{n},
				}
				mp[src] = c
			}
			n.cc = append(n.cc, c)
		}
	}

	return mp, nil
}

// Datasource attempts to return the smallest possible sub-tree for the given request
//
// Flow outline:
// * Find the start node that corresponds to the provided definition
// * Recursively traverse down the branches
// * When returning, see if the given node can be merged with it's datasource
//
// We try to offload as much work to the datasource to reduce the size of the output frames to reduce
// additional processing.
func (m *model) datasource(ctx context.Context, def *FrameDefinition) (ds Datasource, err error) {
	start, ok := m.nodes[def.Source]
	if !ok {
		return nil, fmt.Errorf("unresolved source: %s", def.Source)
	}

	err = m.validateBranch(start)
	if err != nil {
		return
	}

	return m.reduceBranch(ctx, start)
}

func (m *model) validateBranch(n *modelGraphNode) (err error) {
	err = n.step.Validate()
	if err != nil {
		return err
	}

	for _, c := range n.cc {
		err = m.validateBranch(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *model) reduceBranch(ctx context.Context, n *modelGraphNode) (out Datasource, err error) {
	// leaf node, nothing else to do
	if len(n.cc) == 0 {
		return n.step.Run(ctx)
	}

	// traverse down the branches
	auxO := make([]Datasource, len(n.cc))
	for i, c := range n.cc {
		auxO[i], err = m.reduceBranch(ctx, c)
		if err != nil {
			return
		}
	}

	// for now only the join step expects multiple inputs; it can't be reduced
	if len(auxO) > 1 {
		return n.step.Run(ctx, auxO...)
	}

	// try to reduce the group step
	o := auxO[0]
	if n.step.Def().Group != nil {
		gds, ok := o.(GroupableDatasource)
		if !ok {
			return n.step.Run(ctx, auxO...)
		}

		ok, err = gds.Group(n.step.Def().Group.GroupDefinition, n.step.Name())
		if err != nil {
			return nil, err
		} else if !ok {
			return n.step.Run(ctx, auxO...)
		}

		return gds, nil
	}

	return n.step.Run(ctx, auxO...)
}
