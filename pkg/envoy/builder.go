package envoy

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	graphBuilder struct {
		pp []Processor
	}

	// Some special bits
	Processor interface {
		Process(ctx context.Context, state *ExecState) error
	}
)

func NewGraphBuilder(pp ...Processor) *graphBuilder {
	return &graphBuilder{
		pp: pp,
	}
}

func (b *graphBuilder) Build(ctx context.Context, rr ...resource.Interface) (*graph, error) {
	g := newGraph()

	// Prepare nodes for all resources
	nn := make(nodeSet, 0, len(rr))
	for _, r := range rr {
		nn = nn.add(newNode(r))
	}

	// Let's keep track of the missing deps.
	mMap := make(map[resource.Interface]resource.RefSet)

	// Index all resources for nicer lookups
	nIndex := make(nodeIndex)
	nIndex.Add(nn...)

	// Build the graph
	for _, cNode := range nn {
		refs := cNode.res.Refs()
		missingRefs := make(resource.RefSet, 0, len(refs))

		// Attempt to connect all available nodes
		for _, ref := range refs {
			rNode := nIndex.GetRef(ref)

			if rNode == nil {
				missingRefs = append(missingRefs, ref)
				continue
			}

			// The resource is available; we can connect the two
			g.addChild(cNode, rNode)
			g.addParent(rNode, cNode)
		}

		g.addNode(cNode)
		if len(missingRefs) > 0 {
			mMap[cNode.res] = missingRefs
		}
	}

	// @todo fixme..
	g.invert()
	g.DryRun()

	// Do any dep. related preprocessing
	var state *ExecState
	var err error
	err = func() error {
		for {
			state, err = g.Next(ctx)
			if err != nil {
				return err
			} else if state == nil {
				return nil
			}

			// Copy state so we don't alter the original one
			nState := state
			nState.MissingDeps = mMap[state.Res]

			for _, p := range b.pp {
				err = p.Process(ctx, nState)
				if err != nil {
					return err
				}
			}
		}
	}()

	if err != nil {
		return nil, err
	}

	// @todo fixme..
	g.invert()
	g.ProdRun()

	return g, nil
}
