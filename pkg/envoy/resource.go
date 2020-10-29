package envoy

type (
	Resource interface {
		Identifiers() ResourceIdentifiers
		ResourceType() string
		Refs() NodeRefSet
	}

	NodeRefSet []*NodeRef
	NodeRef    struct {
		// @todo check with Denis regarding strings here (the cdocs comment)
		ResourceType string
		Identifiers  ResourceIdentifiers
	}

	ResourceIdentifiers map[string]bool

	nodeIndex map[string]map[string]*node
)

func (ri ResourceIdentifiers) Add(ii ...string) ResourceIdentifiers {
	for _, i := range ii {
		if len(i) > 0 {
			ri[i] = true
		}
	}

	return ri
}

func (ri ResourceIdentifiers) Remove(ii ...string) ResourceIdentifiers {
	for _, i := range ii {
		delete(ri, i)
	}

	return ri
}

func (ri ResourceIdentifiers) HasAny(ii ...string) bool {
	for _, i := range ii {
		if ri[i] {
			return true
		}
	}

	return false
}

func (ri ResourceIdentifiers) StringSlice() []string {
	ss := make([]string, 0, len(ri))
	for k := range ri {
		ss = append(ss, k)
	}
	return ss
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

func (ri nodeIndex) GetRef(ref *NodeRef) *node {
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
