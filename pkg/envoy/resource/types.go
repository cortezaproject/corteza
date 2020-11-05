package resource

type (
	Interface interface {
		Identifiers() Identifiers
		ResourceType() string
		Refs() RefSet
	}

	RefSet []*Ref
	Ref    struct {
		// @todo check with Denis regarding strings here (the cdocs comment)
		// @todo should this become node type instead?
		ResourceType string
		Identifiers  Identifiers
	}

	Identifiers map[string]bool
)

const (
	APPLICATION_RESOURCE_TYPE       = "application"
	COMPOSE_CHART_RESOURCE_TYPE     = "composeChart"
	COMPOSE_MODULE_RESOURCE_TYPE    = "composeModule"
	COMPOSE_NAMESPACE_RESOURCE_TYPE = "composeNamespace"
	COMPOSE_PAGE_RESOURCE_TYPE      = "composePage"
	COMPOSE_RECORD_RESOURCE_TYPE    = "ComposeRecord"
	RBAC_RESOURCE_TYPE              = "rbacRule"
	ROLE_RESOURCE_TYPE              = "role"
	USER_RESOURCE_TYPE              = "user"
)

func (ri Identifiers) Add(ii ...string) Identifiers {
	for _, i := range ii {
		if len(i) > 0 {
			ri[i] = true
		}
	}

	return ri
}

func (ri Identifiers) Remove(ii ...string) Identifiers {
	for _, i := range ii {
		delete(ri, i)
	}

	return ri
}

func (ri Identifiers) HasAny(ii ...string) bool {
	for _, i := range ii {
		if ri[i] {
			return true
		}
	}

	return false
}

func (ri Identifiers) StringSlice() []string {
	ss := make([]string, 0, len(ri))
	for k := range ri {
		ss = append(ss, k)
	}
	return ss
}

func (ss RefSet) FilterByResourceType(rt string) RefSet {
	rr := make(RefSet, 0, len(ss))

	for _, s := range ss {
		if s.ResourceType == rt {
			rr = append(rr, s)
		}
	}

	return rr
}
