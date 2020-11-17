package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeNamespace struct {
		*base
		Res *types.Namespace
	}
)

func NewComposeNamespace(ns *types.Namespace) *ComposeNamespace {
	r := &ComposeNamespace{base: &base{}}
	r.SetResourceType(COMPOSE_NAMESPACE_RESOURCE_TYPE)
	r.Res = ns

	r.AddIdentifier(identifiers(ns.Slug, ns.Name, ns.ID)...)

	return r
}

func (r *ComposeNamespace) SysID() uint64 {
	return r.Res.ID
}
