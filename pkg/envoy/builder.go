package graph

import (
	"context"
)

type (
	graphBuilder struct {
		s resourceSanitizer
		v resourceValidator
	}

	// Some special bits
	resourceSanitizer interface {
		Sanitize(ctx context.Context, r Resource, missing NodeRefSet) error
	}
	resourceValidator interface {
		Validate(ctx context.Context, r Resource) error
	}
)

func NewGraphBuilder(s resourceSanitizer, v resourceValidator) *graphBuilder {
	return &graphBuilder{
		s: s,
		v: v,
	}
}

func (b *graphBuilder) Build(ctx context.Context, rr ...Resource) (*graph, error) {
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
		missingRefs := make(NodeRefSet, 0, len(refs))

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

		// Go for a quick sanitization
		if b.s != nil {
			err := b.s.Sanitize(ctx, cNode.res, missingRefs)
			if err != nil {
				return nil, err
			}
		}

		// Go for a quick validation
		if b.v != nil {
			err := b.v.Validate(ctx, cNode.res)
			if err != nil {
				return nil, err
			}
		}

		g.addNode(cNode)
	}

	return g, nil
}
