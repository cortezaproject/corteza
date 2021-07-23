package report

import "context"

type (
	modelGraphNode struct {
		step step
		ds   Datasource

		pp []*modelGraphNode
		cc []*modelGraphNode
	}
)

// Internally, the model uses a stepGraph to resolve the dependencies between the steps
// allowing us to perform some preprocessing, such as size reduction and shape validation
//
// @todo most of this might need to be done at runtime, not buildtime
func (m *model) buildStepGraph(ss stepSet, dd DatasourceSet) ([]*modelGraphNode, error) {
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
