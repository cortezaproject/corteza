package envoyx

type (
	// DepGraph provides a collection of optionally connected subgraphs
	//
	// Each subgraph is dedicated for a specific node scope.
	// A subgraph with scoped nodes may be connected to the subgraph with no
	// defined scope.
	DepGraph struct {
		graphs []*depSubgraph
	}

	// depSubgraph provides a bi-directional dependency graph
	depSubgraph struct {
		scope Scope

		nodes   []*depNode
		nodeMap map[*Node]*depNode

		parentSubgraphs map[*depSubgraph]bool
		childSubgraphs  map[*depSubgraph]bool
	}

	// depNode provides a wrapper around the resource for eas of use within the dep. graph
	depNode struct {
		Node *Node

		parents  map[*depNode]bool
		children map[*depNode]bool

		// Index to keep track of missing references at build time
		missingReferences map[string]Ref
	}
)

// BuildDepGraph constructs a dependency graph from the provided nodes
//
// We firstly group the nodes by scope, then build a subgraph for each scope,
// and lastly merge the subgraphs into a single graph.
func BuildDepGraph(nn ...*Node) (out *DepGraph) {
	scopes := scopeNodes(nn...)

	aux := make([]*depSubgraph, 0, len(scopes))
	for _, ss := range scopes {
		aux = append(aux, buildDepSubgraph(ss))
	}

	return buildDepGraph(aux...)
}

// scopeNodes groups the nodes by the scope
//
// The function returns a slice of NodeSet where each NodeSet only contains
// nodes of the same scope (or no scope if none defined).
func scopeNodes(nn ...*Node) (out []NodeSet) {
	type (
		// Defining a little wrapper around the NodeSet so we can use pointers
		scopeWrap struct {
			nn NodeSet
		}
	)

	var (
		// resource type -> identifier -> wrap
		scopes = make(map[string]map[string]*scopeWrap)

		// Holds the nodes with no defined scope
		empty NodeSet
	)

	// Bucket nodes into scopes
	//
	// Use maps to make the process as efficient as possible.
	for _, n := range nn {
		// Empty scopes are a special case
		if n.Scope.IsEmpty() {
			empty = append(empty, n)
			continue
		}

		// New resource type
		if _, ok := scopes[n.Scope.ResourceType]; !ok {
			scopes[n.Scope.ResourceType] = make(map[string]*scopeWrap, 4)
		}

		// Check if the identifiers of the current node scope exist in the index.
		// If they do, update the wrap struct by adding this node, then register
		// the same wrap struct to all other identifiers -- this allows for some
		// recovery in case some resource uses a subset of one but not the other.
		hasIdent := false
		firstIdent := ""
		for _, i := range n.Scope.Identifiers.Slice {
			hasIdent = hasIdent || scopes[n.Scope.ResourceType][i] != nil
			if hasIdent {
				firstIdent = i
				break
			}
		}

		if !hasIdent {
			// Not registered yet, create a new wrap for all of the identifiers
			w := &scopeWrap{nn: append(make(NodeSet, 0, 10), n)}
			for _, i := range n.Scope.Identifiers.Slice {
				scopes[n.Scope.ResourceType][i] = w
			}
		} else {
			// Already registered; update it and add missing identifiers
			w := scopes[n.Scope.ResourceType][firstIdent]
			w.nn = append(w.nn, n)
			for _, i := range n.Scope.Identifiers.Slice {
				scopes[n.Scope.ResourceType][i] = w
			}
		}
	}

	// Unpack the map into a slice of NodeSet the return expects
	//
	// @note since all identifiers are registered, the wraps would appear
	//       duplicated so we need to filter them a bit
	out = make([]NodeSet, 0, 10)
	if len(empty) > 0 {
		out = append(out, empty)
	}
	seen := make(map[*scopeWrap]bool)
	for _, s := range scopes {
		for _, ss := range s {
			if seen[ss] {
				continue
			}
			seen[ss] = true
			out = append(out, ss.nn)
		}
	}

	return out
}

// buildDepSubgraph constructs a dep. subgraph from the provided nodes
//
// The function returns a bidirectional graph of the nodes where the parent
// represents the dependency of the current resource (a namespace would be
// a parent of a module).
//
// The build process indexes some data for optimal operations down the line.
func buildDepSubgraph(nn NodeSet) (out *depSubgraph) {
	out = &depSubgraph{
		scope: nn[0].Scope,

		nodes:   make([]*depNode, len(nn)),
		nodeMap: make(map[*Node]*depNode, len(nn)),

		parentSubgraphs: make(map[*depSubgraph]bool),
		childSubgraphs:  make(map[*depSubgraph]bool),
	}

	byIdentifier := make(map[string]map[string]*depNode, 8)
	byNode := make(map[*Node]*depNode, len(nn))

	// 1. index all of the nodes in a map so we can trivially connect them later
	for i, _n := range nn {
		n := _n

		// Function blindly trusts it will be called with the correct data.
		// @todo consider implementing this check but make sure all of the decoders
		//       correctly set the scopes.
		// if !n.Scope.Equals(nn[0].Scope) {
		// 	panic("invalid state: subgraphs can only be constructed with nodes from the same scope")
		// }

		aux := &depNode{
			Node: n,

			// Keep track of missing references so we can figure them out optimally
			missingReferences: make(map[string]Ref),

			parents:  make(map[*depNode]bool),
			children: make(map[*depNode]bool),
		}

		for field, ref := range n.References {
			aux.missingReferences[field] = ref
		}

		byNode[n] = aux
		out.nodes[i] = aux
		out.nodeMap[n] = aux

		if _, ok := byIdentifier[n.ResourceType]; !ok {
			byIdentifier[n.ResourceType] = make(map[string]*depNode, 8)
		}

		for _, i := range n.Identifiers.Slice {
			byIdentifier[n.ResourceType][i] = byNode[n]
		}
	}

	// 2. link up the node with it's dependencies
	for _, _n := range out.nodes {
		n := _n

		for field, ref := range n.Node.References {
			resource := ref.ResourceType
			found := false

			for _, i := range ref.Identifiers.Slice {
				ch, ok := byIdentifier[resource][i]
				found = found || ok

				if ok {
					delete(n.missingReferences, field)
					n.parents[ch] = true
					ch.children[n] = true
					break
				}
			}
		}
	}

	return
}

// buildDepGraph constructs a dependency graph from the given set of subgraphs
//
// The subgraphs can be connected in case a scoped node would reference a
// unscoped node (unscoped nodes can not reference scoped nodes, nor can nodes
// from different scopes -- unneeded and removes some complexity).
func buildDepGraph(gg ...*depSubgraph) (out *DepGraph) {
	out = &DepGraph{
		graphs: make([]*depSubgraph, 0, len(gg)),
	}

	// Get the unscoped graph
	//
	// For now, we can only xref to unscoped graphs.
	// This might need to be generalized.
	unscopedG := unscopedSubgraph(gg...)

	// Iterate all graphs and try to resolve unresolved deps
	for _, g := range gg {
		out.graphs = append(out.graphs, g)

		if unscopedG == nil {
			continue
		}

		// Get nodes with missing references
		missingNodes := g.nodesWithMissingRefs()
		if len(missingNodes) == 0 {
			continue
		}

		// Try to resolve missing refs using the unscoped graph
		for _, _mn := range missingNodes {
			mn := _mn

			for field, ref := range mn.missingReferences {
				n := depNodeForRef(ref, unscopedG.nodes...)
				if n == nil {
					continue
				}

				mn.parents[n] = true
				n.children[mn] = true

				delete(mn.missingReferences, field)
				if len(mn.missingReferences) == 0 {
					mn.missingReferences = nil
				}
			}
		}

		// xref the subgraphs
		g.parentSubgraphs[unscopedG] = true
		unscopedG.childSubgraphs[g] = true
	}

	return
}

// Graph traversal functions

// Roots returns all nodes which are considered as root resources based on the current state
//
// For the most part, these are all resources with no parent resources.
// If all resources define parents, then some home brew logic is ran
func (g DepGraph) Roots() (out NodeSet) {
	for _, sg := range g.graphs {
		out = append(out, sg.Roots()...)
	}
	return
}

// ParentForRef returns a parent node of n which matches ref (nil if none)
func (g DepGraph) ParentForRef(n *Node, ref Ref) (out *Node) {
	for _, sg := range g.graphs {
		out = sg.ParentForRef(n, ref)
		if out != nil {
			return
		}
	}
	return
}

// ParentForRT returns a set of parent nodes matching the resource type
func (g DepGraph) ParentForRT(n *Node, rt string) (out NodeSet) {
	for _, sg := range g.graphs {
		out = sg.ParentForRT(n, rt)
		if out != nil {
			return
		}
	}
	return
}

// ChildrenForResourceType returns child nodes of n which match the resource type
func (g DepGraph) ChildrenForResourceType(n *Node, rt string) (out NodeSet) {
	for _, sg := range g.graphs {
		out = sg.ChildrenForResourceType(n, rt)
		if out != nil {
			return
		}
	}
	return
}

func (g DepGraph) NodeForRef(ref Ref) (out *Node) {
	aux := make([]*Node, 0, 10)
	for _, sg := range g.graphs {
		for _, n := range sg.nodes {
			aux = append(aux, n.Node)
		}
	}

	return NodeForRef(ref, aux...)
}

// Children returns all child nodes of n
func (g DepGraph) Children(n *Node) (out NodeSet) {
	for _, sg := range g.graphs {
		out = sg.Children(n)
		if out != nil {
			return
		}
	}
	return
}

// MissingRegs returns a slice of all refs that are requested but not found in the graph
func (g DepGraph) MissingRefs() (out []map[string]Ref) {
	for _, sg := range g.graphs {
		for _, n := range sg.nodes {
			if len(n.missingReferences) == 0 {
				continue
			}
			out = append(out, n.missingReferences)
		}
	}

	return
}

func (g DepGraph) allNodes() (out []*depNode) {
	for _, sg := range g.graphs {
		out = append(out, sg.nodes...)
	}

	return
}

// Roots returns all nodes which are considered as root resources based on the current state
//
// For the most part, these are all resources with no parent resources.
// If all resources define parents, then some home brew logic is ran
//
// @todo when we add more resources, we might need to expand this; for now
//       it should work just fine.
func (g depSubgraph) Roots() (out NodeSet) {
	for _, n := range g.nodes {
		if needyResources[n.Node.ResourceType] {
			continue
		}
		out = append(out, n.Node)
	}

	if len(out) != 0 {
		return
	}

	for _, n := range g.nodes {
		if superNeedyResources[n.Node.ResourceType] {
			continue
		}
		out = append(out, n.Node)
	}

	return
}

// ParentForRef returns a parent node of n which matches ref (nil if none)
func (g depSubgraph) ParentForRef(n *Node, ref Ref) *Node {
	return NodeForRef(ref, g.parent(n)...)
}

// ParentForRT returns a set of parent nodes matching the resource type
func (g depSubgraph) ParentForRT(n *Node, rt string) NodeSet {
	return NodesForResourceType(rt, g.parent(n)...)
}

// ChildrenForResourceType returns child nodes of n which match the resource type
func (g depSubgraph) ChildrenForResourceType(n *Node, rt string) (out NodeSet) {
	out = g.Children(n)
	out = NodesForResourceType(rt, out...)
	return
}

// Children returns all child nodes of n
func (g depSubgraph) Children(n *Node) (out NodeSet) {
	aux, ok := g.nodeMap[n]
	if !ok {
		return
	}

	for n := range aux.children {
		out = append(out, n.Node)
	}
	return
}

func (g depSubgraph) parent(n *Node) (out NodeSet) {
	aux, ok := g.nodeMap[n]
	if !ok {
		return
	}

	for n := range aux.parents {
		out = append(out, n.Node)
	}
	return
}

// Utility functions

func (g *depSubgraph) nodesWithMissingRefs() (out []*depNode) {
	out = make([]*depNode, 0, 3)
	for _, n := range g.nodes {
		if len(n.missingReferences) > 0 {
			out = append(out, n)
		}
	}

	return
}

func unscopedSubgraph(gg ...*depSubgraph) *depSubgraph {
	for _, g := range gg {
		if g.scope.IsEmpty() {
			return g
		}
	}

	return nil
}

// depNodeForRef returns a node which matches the ref (nil if none)
func depNodeForRef(ref Ref, nn ...*depNode) (out *depNode) {
	for _, n := range nn {
		if n.Node.ResourceType != ref.ResourceType {
			continue
		}

		if n.Node.Identifiers.HasIntersection(ref.Identifiers) {
			return n
		}
	}

	return
}

// depNodesByResourceType returns a map where key is resource type, value slice of nodes
func depNodesByResourceType(nn ...*depNode) (out map[string][]*depNode) {
	out = make(map[string][]*depNode, 4)
	for _, n := range nn {
		out[n.Node.ResourceType] = append(out[n.Node.ResourceType], n)
	}

	return
}

// unpackDepNodes extracts envoy Nodes from the dep. nodes
//
// The function does no validation nor filtering for nil values.
func unpackDepNodes(nn ...*depNode) (out NodeSet) {
	out = make(NodeSet, 0, len(nn))
	for _, n := range nn {
		out = append(out, n.Node)
	}
	return out
}
