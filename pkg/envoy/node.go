package envoy

import "github.com/cortezaproject/corteza-server/pkg/envoy/resource"

type (
	// the node struct is used for nicer graph state management
	nodeSet []*node
	node    struct {
		res     resource.Interface
		missing resource.RefSet

		pp nodeSet
		cc nodeSet
	}

	// resource type -> identifier -> []nodes
	// There can be multiple resources with same identifier; for example
	// two modules under different namespaces.
	nodeIndex map[string]map[string]nodeSet
)

func newNode(res resource.Interface) *node {
	return &node{
		res: res,
		cc:  make(nodeSet, 0, 10),
		pp:  make(nodeSet, 0, 10),
	}
}

func (nn nodeSet) add(mm ...*node) nodeSet {
	return append(nn, mm...)
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

func (ri nodeIndex) Add(nn ...*node) {
	for _, n := range nn {
		rt := n.res.ResourceType()
		if _, has := ri[rt]; !has {
			ri[rt] = make(map[string]nodeSet)
		}

		for i := range n.res.Identifiers() {
			if ri[rt][i] == nil {
				ri[rt][i] = make(nodeSet, 0, 5)
			}
			ri[rt][i] = append(ri[rt][i], n)
		}
	}
}

func (ri nodeIndex) GetRef(ref *resource.Ref) *node {
	refIi, has := ri[ref.ResourceType]
	if !has {
		return nil
	}

	for i := range ref.Identifiers {
		rr, has := refIi[i]
		if !has || len(rr) == 0 {
			continue
		}

		// No constraints? no worries
		if ref.Constraints == nil || len(ref.Constraints) == 0 {
			return rr[0]
		}

		// Constraints? check if ok
		// If this loop makes you sick, don't worry; this will be at most 3x3x1
		for _, r := range rr {
			for _, c := range ref.Constraints {
				for _, ref := range r.res.Refs() {
					if ref.ResourceType == c.ResourceType && ref.Identifiers.HasAny(c.Identifiers) {
						return r
					}
				}
			}
		}

	}

	return nil
}

func (ri nodeIndex) GetResourceType(rt string) nodeSet {
	ix, has := ri[rt]
	if !has {
		return nil
	}
	nm := make(map[*node]bool)
	for _, nn := range ix {
		for _, n := range nn {
			nm[n] = true
		}
	}

	nn := make(nodeSet, 0, len(nm))
	for n := range nm {
		nn = append(nn, n)
	}
	return nn
}
