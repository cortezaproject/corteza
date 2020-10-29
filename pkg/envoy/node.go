package envoy

type (
	// the node struct is used for nicer graph state management
	nodeSet []*node
	node    struct {
		res Resource

		pp nodeSet
		cc nodeSet

		state NodeState
	}

	// @tbd
	// stateLegacyID: stateSysID
	nodeStateEntry map[string]string
	// stateResourceType: nodeStateEntry
	NodeState map[string]nodeStateEntry
)

func newNode(res Resource) *node {
	return &node{
		res:   res,
		cc:    make(nodeSet, 0, 10),
		pp:    make(nodeSet, 0, 10),
		state: make(NodeState),
	}
}

func (nn nodeSet) add(mm ...*node) nodeSet {
	return append(nn, mm...)
}

func (nn nodeSet) filter(f func(n *node) bool) nodeSet {
	mm := make(nodeSet, 0, len(nn))
	for _, n := range nn {
		if f(n) {
			mm = append(mm, n)
		}
	}
	return mm
}

func (nn nodeSet) has(m *node) bool {
	for _, n := range nn {
		if n == m {
			return true
		}
	}

	return false
}

func (nn nodeSet) remove(mm ...*node) nodeSet {
	if len(mm) <= 0 {
		return nn
	}

	nClean := make(nodeSet, 0, len(nn))
	mmSet := make(nodeSet, 0, len(mm))
	mmSet = append(mmSet, mm...)

	for _, n := range nn {
		if !mmSet.has(n) {
			nClean = append(nClean, n)
		}
	}

	return nClean
}

func (s NodeState) addEntry(res string, e nodeStateEntry) {
	if s[res] == nil {
		s[res] = make(nodeStateEntry)
	}

	s[res] = s[res].merge(e)
}

func (s NodeState) merge(e NodeState) NodeState {
	for k, v := range e {
		if s[k] == nil {
			s[k] = v
		} else {
			s[k].merge(v)
		}
	}

	return s
}

func (se nodeStateEntry) merge(e nodeStateEntry) nodeStateEntry {
	for k, v := range e {
		se[k] = v
	}

	return se
}
