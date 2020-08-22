package envoy

import (
	"context"
	"math"
)

type (
	// Graph is the root structure of any graph.
	//
	// The use of a graph layer allows us to tackle relation problems of arbitrary structure;
	// from simple acyclic graphs to complex cyclic graphs.
	// The graph is calculated on-the-fly, meaning that it doesn't addopt the usual approach
	// where all of the nodes are connected via pointers/references.
	// This approach greatly simplifies the entire process of maintaining a chart.
	// In terms of time complexity, in comparison to other layers of the system, this step is free.
	// When it comes to larger setups (custom CRM) the number of nodes and relations is well below 1000.
	// Due to the simple interface, we can define a more optimal graph implementation as a sepperate module.
	Graph struct {
		// nodes defines a set of all available structures conforming to the Node interface
		nodes []Node
		// Since everything is calculated on-the-fly, we need a simple boolean flag to determine if the graph is inverted
		invert bool
		// Allows us to keep track of all the nodes that were determined as conflicting.
		// A node is considered conflicting, when it is part of a cycle.
		// If the data input is consistent, the conflicting nodes will be consistent, but there is no guarantee
		// what node in a cycle will be selected as a conflicting.
		conflicts set
		// Allows us to easily keep track of nodes that were already processed.
		processed set
	}

	// A simple "set" implementation for simpler, quicker checks
	set map[string]bool

	// Node defines an interface that any Graph member must conform to.
	//
	// A Node should define some operation that should be performed when the thing is executed.
	// For example, compose:record resource import, system:user import.
	Node interface {
		// Resource provides the unique resource identifier this Node is designed for, such as compose:module
		Resource() string

		// ID provides the node's identifier, such as the resource's ID
		ID() string

		// Relations provides a list of node's relations
		// This can be calculated on the fly based on the node's state and don't need to be
		// built in to the node struct.
		Relations() map[string][]string

		// Run should implement the actual operation that should be performend when the node is invoked.
		// This can be as simple or as complex as needed
		// It's also provided a set of child and parent nodes so we can easily provide node context
		Run(ctx context.Context, cc []Node, pp []Node) error

		// ResolveConflict should implement any operation that should occur when the node
		// causes a dependency conflict. For example -- partially import the records (without relations)
		// and correct those relations when executing the Run function
		// It's also provided a set of child and parent nodes so we can easily provide node context
		ResolveConflict(ctx context.Context, cc []Node, pp []Node) error
	}
)

// Add registers a given set of nodes nn into the graph g
func (g *Graph) Add(nn ...Node) {
	g.nodes = append(g.nodes, nn...)
}

// Remove removes the set of nodes nn from the graph h
func (g *Graph) Remove(nn ...Node) {
	if len(nn) <= 0 {
		return
	}

	mm := make([]Node, 0, len(g.nodes))
	for _, m := range g.nodes {
		for _, n := range nn {
			if NodeHash(m) == NodeHash(n) && g.canRemove(n) {
				goto skip
			}
		}
		mm = append(mm, m)

	skip:
	}
	g.nodes = mm
}

func (g *Graph) canRemove(n Node) bool {
	if len(g.nodes) <= 1 {
		return true
	}

	// When the first unprocessed child node is found, terminate
	for _, m := range g.Children(n) {
		if !g.processed[NodeHash(m)] {
			return false
		}
	}

	return true
}

func (g *Graph) MarkProcessed(nn ...Node) {
	if g.processed == nil {
		g.processed = make(set)
	}
	for _, n := range nn {
		g.processed[NodeHash(n)] = true
	}
}

// FindNode returns all nodes that match the given resource and identifiers
func (g *Graph) FindNode(res string, IDs ...string) []Node {
	nn := make([]Node, 0, len(IDs))
	for _, n := range g.nodes {
		for _, ID := range IDs {
			if NodeHash(n) == NodeHashRaw(res, ID) {
				nn = append(nn, n)
			}
		}
	}

	return nn
}

// Invert inverts all of the relations in the graph
func (g *Graph) Invert() {
	g.invert = !g.invert
}

// SetNodeConflict is a helper to register this node as a conflictor
func (g *Graph) SetNodeConflict(n Node) {
	if g.conflicts == nil {
		g.conflicts = make(set)
	}

	g.conflicts[NodeHash(n)] = true
}

func (g *Graph) removeProcessed(nn []Node) []Node {
	mm := make([]Node, 0, len(nn))
	for _, n := range nn {
		if !g.processed[NodeHash(n)] {
			mm = append(mm, n)
		}
	}
	return mm
}

// Children provides a list of children of the node n
//
// Child nodes are calculated on the fly based on node Relations()
func (g *Graph) Children(n Node) []Node {
	if !g.invert {
		return g.removeProcessed(g.children(n))
	}
	return g.removeProcessed(g.parents(n))
}

func (g *Graph) children(n Node) []Node {
	nn := make([]Node, 0)
	for res, IDs := range n.Relations() {
		nn = append(nn, g.FindNode(res, IDs...)...)
	}

	return nn
}

func (g *Graph) ChildrenA(n Node) []Node {
	if !g.invert {
		return g.children(n)
	}
	return g.parents(n)
}

// Parents provides a list of parents of the node n
//
// Parent nodes are calculated on the fly based on node Relations()
func (g *Graph) Parents(n Node) []Node {
	if !g.invert {
		return g.removeProcessed(g.parents(n))
	}
	return g.removeProcessed(g.children(n))
}

func (g *Graph) parents(n Node) []Node {
	nn := make([]Node, 0)

	for _, m := range g.nodes {
		r := m.Relations()
		if r == nil {
			continue
		}
		if IDs, has := r[n.Resource()]; has {
			for _, ID := range IDs {
				if n.ID() == ID {
					nn = append(nn, m)
					break
				}
			}
		}
	}

	return nn
}

// ParentsC returns only the parent nodes that are not registered as conflicting
func (g *Graph) ParentsC(n Node) []Node {
	pp := g.Parents(n)
	mm := make([]Node, 0, int(math.Max(float64(len(pp)-len(g.conflicts)), float64(1))))
	for _, p := range pp {
		if !g.conflicts[NodeHash(p)] {
			mm = append(mm, p)
		}
	}

	return mm
}

func (g *Graph) ParentsA(n Node) []Node {
	if !g.invert {
		return g.parents(n)
	}
	return g.children(n)
}

// Validate performs a basic data validation over all the nodes.
func (g *Graph) Validate() error {
	// @todo...
	return nil
}

func (g *Graph) Nodes() []Node {
	nn := make([]Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		if !g.processed[NodeHash(n)] {
			nn = append(nn, n)
		}
	}
	return nn
}

// Run invokes all operations while respecting relations (dependencies) and solving
// dependency conflicts.
//
// Run does the following:
// 	- Inverts the graph to allow better memory management (@todo docs),
// 	- calls ResolveConflict on any node that causes a dependency conflict (a cycle),
// 	- calls Run on all nodes, respecting dependencies.
//
// The order of above operations is a bit more complex, but the general flow is that.
func (g *Graph) Run(ctx context.Context) error {
	for len(g.nodes) > 0 {
		// Find all root nodes in the current graph state; those nodes are allowed to run.
		nn := make([]Node, 0, len(g.Nodes())/2)
		for _, n := range g.Nodes() {
			// We should not take into account conflicted parent nodes, as they already resolved
			// the conflict.
			if len(g.ParentsC(n)) == 0 {
				nn = append(nn, n)
			}
		}

		if len(nn) > 0 {
			err := g.runRegular(ctx, nn)
			if err != nil {
				return err
			}
		} else {
			err := g.runResolution(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// runRegular doesn't do anything special; it just runs all the nodes that
// are allowed to run.
func (g *Graph) runRegular(ctx context.Context, nn []Node) error {
	for _, n := range nn {
		err := n.Run(ctx, g.Children(n), g.ParentsA(n))
		if err != nil {
			return err
		}

		// When a child node is processed check if it's parent can be removed.
		// A parent node can only be removed if all of it's child nodes have already been
		// processed
		g.MarkProcessed(n)
		g.Remove(g.ParentsA(n)...)
		g.Remove(n)
	}

	return nil
}

// runResolution attempts to resolve dependency conflicts in case there is a cycle (no root node).
//
// Since there are still nodes in the graph and there is no root node (its all just cycles) we can:
// 	- Take any node of any cycle,
// 	- instruct the node to resolve the conflict,
// 	- mark the node as conflicted so it will be properly processed later on,
// 	- keep the node in the graph as it should do another round of processing at the end.
func (g *Graph) runResolution(ctx context.Context) error {
	var n Node
	// @todo taking any node isn't entirely optimal since they might not be in a cycle.
	// For example: A -> B -> A -> c -> D; where A B C is the cycle and C is a branch from the cycle.
	// The code will works just fine, but it won't be that optimal so it should be improved to do
	// actual cycle detection.
	for _, m := range g.Nodes() {
		if !g.conflicts[NodeHash(m)] {
			n = m
			break
		}
	}

	err := n.ResolveConflict(ctx, g.Children(n), g.Parents(n))
	if err != nil {
		return err
	}
	g.SetNodeConflict(n)

	return nil
}

// NodeHash is a helper to calculate a guid for the given node n
func NodeHash(n Node) string {
	return NodeHashRaw(n.Resource(), n.ID())
}

// NodeHashRaw is a helper to calculate a guid for the given resource and ID
func NodeHashRaw(resource, ID string) string {
	return resource + "/" + ID
}
