package envoy

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	graphBuilder struct {
		pp []Preparer
	}

	// Rc is an alias for the ResourceState channel
	Rc chan *ResourceState

	Preparer interface {
		Prepare(ctx context.Context, ee ...*ResourceState) error
	}

	// Encoder encodes all resources provided by the Rc until nil is passed
	//
	// Encoding errors are passed via Ec.
	Encoder interface {
		Encode(ctx context.Context, p Provider) error
	}

	PrepareEncoder interface {
		Preparer
		Encoder
	}
)

func NewGraphBuilder(pp ...Preparer) *graphBuilder {
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

	// Do any dep. related preprocessing
	var state *ResourceState
	var err error

	err = func() error {
		for {
			state, err = g.NextInverted(ctx)
			if err != nil {
				return err
			} else if state == nil {
				return nil
			}

			// Copy state so we don't alter the original one
			nState := state
			nState.MissingDeps = mMap[state.Res]

			for _, p := range b.pp {
				err = p.Prepare(ctx, nState)
				if err != nil {
					return err
				}
			}
		}
	}()

	if err != nil {
		return nil, err
	}

	g.Relink()
	g.reset()

	return g, nil
}
