package envoy

import (
	"context"
	"io"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	builder struct {
		pp []Preparer
	}

	Stream struct {
		Source     io.Reader
		Resource   string
		Identifier string
	}

	Preparer interface {
		Prepare(ctx context.Context, ee ...*ResourceState) error
	}

	// Encoder encodes all resources provided by the Rc until nil is passed
	//
	// Encoding errors are passed via Ec.
	Encoder interface {
		Encode(ctx context.Context, p Provider) error
	}

	// Streammer provides a set of streams of encoded documents
	Streammer interface {
		Stream() []*Stream
	}

	PrepareEncoder interface {
		Preparer
		Encoder
	}
	PrepareEncodeStreammer interface {
		Preparer
		Encoder
		Streammer
	}
)

func NewBuilder(pp ...Preparer) *builder {
	return &builder{
		pp: pp,
	}
}

// Build builds the graph that is used for structured data processing
//
// Outline:
// 1. Build an initial graph so that we can do some structured preprocessing.
// 2. Preprocess the resources based on the initial graph. The initial graph
//    should remain unchanged. Preprocessing can request additional references and
//    constraints.
// 3. Build a final graph based on the preprocessing modifications.
func (b *builder) Build(ctx context.Context, rr ...resource.Interface) (*graph, error) {
	var err error

	g := b.buildGraph(rr)

	// Do any dep. related preprocessing
	var state *ResourceState
	err = func() error {
		for {
			state, err = g.NextInverted(ctx)
			if err != nil {
				return err
			} else if state == nil {
				return nil
			}

			for _, p := range b.pp {
				err = p.Prepare(ctx, state)
				if err != nil {
					return err
				}
			}
		}
	}()

	if err != nil {
		return nil, err
	}

	g = b.buildGraph(rr)
	return g, nil
}

func (b *builder) buildGraph(rr []resource.Interface) *graph {
	g := newGraph()

	// Prepare nodes for all resources
	nn := make(nodeSet, 0, len(rr))
	for _, r := range rr {
		nn = nn.add(newNode(r))
	}

	// Index all resources for nicer lookups
	nIndex := make(nodeIndex)
	nIndex.Add(nn...)

	// Build the graph
	for _, cNode := range nn {
		refs := cNode.res.Refs()
		missingRefs := make(resource.RefSet, 0, len(refs))

		// Attempt to connect all available nodes
		for _, ref := range refs {
			// Handle wildcard references
			if ref.IsWildcard() {
				nn := nIndex.GetResourceType(ref.ResourceType)
				// Connect the nodes
				for _, n := range nn {
					g.addChild(cNode, n)
					g.addParent(n, cNode)
				}
			} else {
				// Handle regular references
				rNode := nIndex.GetRef(ref)

				if rNode == nil {
					missingRefs = append(missingRefs, ref)
					continue
				}

				// The resource is available; we can connect the two
				g.addChild(cNode, rNode)
				g.addParent(rNode, cNode)
			}
		}

		g.addNode(cNode)
	}

	return g
}
