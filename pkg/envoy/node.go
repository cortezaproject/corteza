package envoy

import "github.com/cortezaproject/corteza-server/pkg/envoy/resource"

type (
	// the node struct is used for nicer graph state management
	nodeSet []*node
	node    struct {
		res resource.Interface

		pp nodeSet
		cc nodeSet
	}

	nodeIndex map[string]map[string]*node
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

func (ri nodeIndex) Add(nn ...*node) {
	for _, n := range nn {
		rt := n.res.ResourceType()
		if _, has := ri[rt]; !has {
			ri[rt] = make(map[string]*node)
		}

		for i := range n.res.Identifiers() {
			ri[rt][i] = n
		}
	}
}

func (ri nodeIndex) GetRef(ref *resource.Ref) *node {
	res, has := ri[ref.ResourceType]
	if !has {
		return nil
	}

	for i := range ref.Identifiers {
		r, has := res[i]
		if has {
			return r
		}
	}

	return nil
}
