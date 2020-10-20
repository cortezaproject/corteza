package envoy

type (
	// Node defines the signature of any valid graph node
	Node interface {
		// Matches checks if the node matches the given resources and **any** of the identifiers.
		Matches(resource string, identifiers ...string) bool

		// Identifiers returns a set of values that identify the node.
		//
		// The identifiers may **not** be unique across all resources, but
		// they **must** be unique inside a given resource.
		Identifiers() NodeIdentifiers

		// Resource returns the Corteza resource identifier that this node handles
		Resource() string

		// Relations returns a set of NodeRelationships regarding this node
		//
		// The graph layer **must** be able to handle dynamic relationships (changed in runtime).
		Relations() NodeRelationships
	}

	// NodeUpdater defines a node that can update its state based on the given set of Nodes
	//
	// For example, a ComposeRecordNode should know how to update the referenced ComposeModule resource.
	NodeUpdater interface {
		Node

		// Update receives a set of nodes that should be used when updating the given node n
		//
		// The caller **must** only provide nodes that the given node n is dependent of (it's parent nodes).
		Update(...Node)
	}

	// NodeSet is a set of Nodes
	NodeSet []Node

	// NodeRelationships holds relationships for a specific node
	NodeRelationships map[string]NodeIdentifiers

	// NodeIdentifiers represents a set of node identifiers
	NodeIdentifiers []string
)

// Add adds a new identifier for the given resource
func (n NodeRelationships) Add(resource string, identifier ...string) {
	if _, has := n[resource]; !has {
		n[resource] = make(NodeIdentifiers, 0, 1)
	}

	n[resource] = n[resource].Add(identifier...)
}

// Add adds a new identifiers
func (ii NodeIdentifiers) Add(identifier ...string) NodeIdentifiers {
	exists := false
	for _, i := range ii {
		for _, j := range identifier {
			exists = exists || i == j
		}
	}

	if !exists {
		ii = append(ii, identifier...)
	}

	return ii
}

// HasAny checks if any of the provided identifiers appear in the given set of identifiers
func (ii NodeIdentifiers) HasAny(jj ...string) bool {
	for _, i := range ii {
		for _, j := range jj {
			if i == j {
				return true
			}
		}
	}

	return false
}

// Has checks if the given NodeSet contains a specific Node
func (ss NodeSet) Has(n Node) bool {
	has := false
	for _, s := range ss {
		mRes := n.Resource()
		mIdd := n.Identifiers()

		has = has || s.Matches(mRes, mIdd...)
	}
	return has
}
