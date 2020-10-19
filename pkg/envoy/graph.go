package envoy

import (
	"context"
	"errors"
	"math"

	"github.com/cortezaproject/corteza-server/pkg/envoy/types"
)

type (
	// Graph struct handless and prcesses all of the dependency related operations
	//
	// This is a cyclic graph where node relationships are determined on-the-fly
	// based on the node properties.
	// Refer to the documentation for additional details.
	Graph struct {
		nodes []types.Node

		// Since it's calculated on the fly, this is all we need
		invert bool

		// A cycle is interpreted as a dependency conflict (deadlock).
		// It's up to the graph's discression to determine what node in the cycle will be used.
		// There is no guarantee that this list will be consistent across multiple runs.
		conflicts types.NodeSet

		// Nodes aren't immediately removed from the graph, so they are firstly marked as processed
		processed types.NodeSet
	}
)

var (
	ErrorDependencyConflict = errors.New("graph: dependency conflict")
)

// NewGraph returns a new Graph struct
//
// It handles some initialization bits, so it's better to use it instead of
// making one yourself.
func NewGraph() *Graph {
	return &Graph{
		conflicts: make(types.NodeSet, 0, 100),
		processed: make(types.NodeSet, 0, 100),
		nodes:     make([]types.Node, 0, 100),
		invert:    false,
	}
}

// Add adds the provided set of nodes into the given graph g
//
// The method doesn't do any existence checks for duplicates.
// It simply pushes the provided nodes.
func (g *Graph) Add(nn ...types.Node) {
	g.nodes = append(g.nodes, nn...)
}

// Remove removes the set of nodes nn from the graph h
//
// The nodes canonly be removed if it doesn't have any unprocessed dependencies (child nodes)
func (g *Graph) Remove(nn ...types.Node) {
	if len(nn) <= 0 {
		return
	}

	mm := make([]types.Node, 0, len(g.nodes))
	for _, m := range g.nodes {
		for _, n := range nn {
			if g.nodesMatch(m, n) && g.canRemove(n) {
				goto skip
			}
		}
		mm = append(mm, m)

	skip:
	}
	g.nodes = mm
}

// FindNode returns all nodes that match the given resource and identifiers
func (g *Graph) FindNode(res string, identifiers ...string) []types.Node {
	nn := make([]types.Node, 0, len(identifiers))
	for _, n := range g.nodes {
		if n.Matches(res, identifiers...) {
			nn = append(nn, n)
		}
	}

	return nn
}

// Invert inverts the graph
//
// Since relationships are calculated on-the-fly, this is a simple bit switch
func (g *Graph) Invert() {
	g.invert = !g.invert
}

// Children provides node n child nodes **excluding** processed nodes
func (g *Graph) Children(n types.Node) []types.Node {
	if !g.invert {
		return g.removeProcessedNodes(g.children(n))
	}
	return g.removeProcessedNodes(g.parents(n))
}

// ChildrenA provides node n child nodes **including** processed nodes
func (g *Graph) ChildrenA(n types.Node) []types.Node {
	if !g.invert {
		return g.children(n)
	}
	return g.parents(n)
}

// Parents provides node n parent nodes **excluding** processed nodes
func (g *Graph) Parents(n types.Node) []types.Node {
	if !g.invert {
		return g.removeProcessedNodes(g.parents(n))
	}
	return g.removeProcessedNodes(g.children(n))
}

// ParentsA provides node n parent nodes **incliding** processed nodes
func (g *Graph) ParentsA(n types.Node) []types.Node {
	if !g.invert {
		return g.parents(n)
	}
	return g.children(n)
}

// ParentsAC provides node n parent nodes **including** processed nodes, **excluding** conflicting nodes
func (g *Graph) ParentsAC(n types.Node) []types.Node {
	pp := g.Parents(n)
	mm := make([]types.Node, 0, int(math.Max(float64(len(pp)-len(g.conflicts)), 1.0)))
	for _, p := range pp {
		if !g.conflicts.Has(p) {
			mm = append(mm, p)
		}
	}

	return mm
}

// Validate performs a basic data validation over all the nodes.
//
// @todo Do we need this on the graph layer?
func (g *Graph) Validate() error {
	return nil
}

// Nodes returns all unprocessed nodes in the given graph g
func (g *Graph) Nodes() []types.Node {
	nn := make([]types.Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		if !g.processed.Has(n) {
			nn = append(nn, n)
		}
	}
	return nn
}

// Next provides the next node that should be processed (inclide the nodes context)
//
// Flow outline:
//  * If there are no more nodes, return nothing (all nil values).
//  * If there is a node with no parent nodes; select that as the next node.
//  * If there is no node with no parent nodes; determine a conflicting node.
//    This returns the conflicting node n, it's parents, it's children and an ErrorDependencyConflict.
func (g *Graph) Next(ctx context.Context) (n types.Node, pp []types.Node, cc []types.Node, err error) {
	if len(g.Nodes()) <= 0 {
		return nil, nil, nil, nil
	}

	for _, m := range g.Nodes() {
		// We should not take into account conflicted parent nodes,
		// as they already resolved the conflict.
		if len(g.ParentsAC(m)) == 0 {
			n = m
			break
		}
	}

	if n != nil {
		// Get the node's child and parent nodes.
		// Attempt parent cleanup
		cc = g.Children(n)
		pp = g.ParentsA(n)
		g.markProcessed(n)
		g.Remove(g.ParentsA(n)...)
		g.Remove(n)

		return n, pp, cc, nil
	}

	// Determine a conflicting node if we stumbled on a conflict
	for _, m := range g.Nodes() {
		if !g.conflicts.Has(m) {
			n = m
			break
		}
	}

	// Get the node's child and parent nodes.
	// No cleanup should be done here since the node isn't yet fully processed.
	cc = g.Children(n)
	pp = g.ParentsA(n)
	g.markConflicting(n)
	return n, pp, cc, ErrorDependencyConflict
}

// Helper methods
// ------------------------------------------------------------------------

func (g *Graph) nodesMatch(n, m types.Node) bool {
	mRes := m.Resource()
	mIdd := m.Identifiers()

	return n.Matches(mRes, mIdd...)
}

func (g *Graph) canRemove(n types.Node) bool {
	if len(g.nodes) <= 1 {
		return true
	}

	// A node can only be removed if all of it's child nodes are processed
	for _, m := range g.Children(n) {
		if !g.processed.Has(m) {
			return false
		}
	}

	return true
}

func (g *Graph) markProcessed(nn ...types.Node) {
	if g.processed == nil {
		g.processed = make(types.NodeSet, 0, len(nn))
	}

	for _, n := range nn {
		if !g.processed.Has(n) {
			g.processed = append(g.processed, n)
		}
	}
}

// helper to mark the node as a conflicting node
func (g *Graph) markConflicting(n types.Node) {
	if g.conflicts == nil {
		g.conflicts = make(types.NodeSet, 0, 1)
	}

	if !g.conflicts.Has(n) {
		g.conflicts = append(g.conflicts, n)
	}
}

func (g *Graph) removeProcessedNodes(nn []types.Node) []types.Node {
	mm := make([]types.Node, 0, len(nn))
	for _, n := range nn {
		if !g.processed.Has(n) {
			mm = append(mm, n)
		}
	}
	return mm
}

func (g *Graph) children(n types.Node) []types.Node {
	nn := make([]types.Node, 0)
	// A simple find all nodes that n is in a relationship with will do the trick
	for res, IDs := range n.Relations() {
		nn = append(nn, g.FindNode(res, IDs...)...)
	}

	return nn
}

func (g *Graph) parents(n types.Node) []types.Node {
	nn := make([]types.Node, 0)

	// A more complex, find all nodes that have n in their relationship.
	// @note can we make this nicer?
	for _, m := range g.nodes {
		r := m.Relations()
		if r == nil {
			continue
		}

		if IDs, has := r[n.Resource()]; has {
			if n.Matches(n.Resource(), IDs...) {
				nn = append(nn, m)
			}
		}
	}

	return nn
}
