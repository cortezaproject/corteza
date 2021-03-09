package resource

import (
	ct "github.com/cortezaproject/corteza-server/compose/types"
	st "github.com/cortezaproject/corteza-server/system/types"
)

type (
	Interface interface {
		Identifiers() Identifiers
		ResourceType() string
		Refs() RefSet
	}

	InterfaceSet []Interface

	IdentifiableInterface interface {
		Interface

		SysID() uint64
	}

	RefableInterface interface {
		Interface

		Ref() string
	}

	RefSet []*Ref
	Ref    struct {
		// @todo check with Denis regarding strings here (the cdocs comment)
		// @todo should this become node type instead?
		ResourceType string
		Identifiers  Identifiers
		Constraints  RefSet
	}

	Identifiers map[string]bool
)

var (
	APPLICATION_RESOURCE_TYPE       = st.ApplicationRBACResource.String()
	COMPOSE_CHART_RESOURCE_TYPE     = ct.ChartRBACResource.String()
	COMPOSE_MODULE_RESOURCE_TYPE    = ct.ModuleRBACResource.String()
	COMPOSE_NAMESPACE_RESOURCE_TYPE = ct.NamespaceRBACResource.String()
	COMPOSE_PAGE_RESOURCE_TYPE      = ct.PageRBACResource.String()
	COMPOSE_RECORD_RESOURCE_TYPE    = "compose:record:"
	RBAC_RESOURCE_TYPE              = "rbac:rule:"
	ROLE_RESOURCE_TYPE              = st.RoleRBACResource.String()
	SETTINGS_RESOURCE_TYPE          = "system:setting:"
	USER_RESOURCE_TYPE              = st.UserRBACResource.String()
	DATA_SOURCE_RESOURCE_TYPE       = "data:raw:"
)

func MakeIdentifiers(ss ...string) Identifiers {
	ii := make(Identifiers)
	ii.Add(ss...)
	return ii
}

func (ri Identifiers) Add(ii ...string) Identifiers {
	for _, i := range ii {
		if len(i) > 0 {
			ri[i] = true
		}
	}

	return ri
}

func (ri Identifiers) HasAny(ii Identifiers) bool {
	for i := range ii {
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

func (ri Identifiers) First() string {
	ss := ri.StringSlice()
	if len(ss) <= 0 {
		return ""
	}
	return ss[0]
}

func (rr InterfaceSet) Walk(f func(r Interface) error) (err error) {
	for _, r := range rr {
		err = f(r)
		if err != nil {
			return
		}
	}

	return nil
}

// Constraint returns the current reference with added constraint
func (r *Ref) Constraint(c *Ref) *Ref {
	if r.Constraints == nil {
		r.Constraints = make(RefSet, 0, 1)
	}

	r.Constraints = append(r.Constraints, &Ref{
		ResourceType: c.ResourceType,
		Identifiers:  MakeIdentifiers(c.Identifiers.StringSlice()...),
	})

	return r
}

// IsWildcard checks if this Ref points to all resources of a specific resource type
func (r *Ref) IsWildcard() bool {
	return r.Identifiers["*"]
}
