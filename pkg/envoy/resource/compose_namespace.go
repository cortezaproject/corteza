package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
)

const (
	COMPOSE_NAMESPACE_RESOURCE_TYPE = "composeNamespace"
)

type (
	composeNamespace struct {
		*base
		Res *types.Namespace
	}
)

func ComposeNamespace(ns *types.Namespace) *composeNamespace {
	r := &composeNamespace{base: &base{}}
	r.SetResourceType(COMPOSE_NAMESPACE_RESOURCE_TYPE)
	r.Res = ns

	r.AddIdentifier(identifiers(ns.Slug, ns.Name, ns.ID)...)

	return r
}

func (m *composeNamespace) SearchQuery() types.NamespaceFilter {
	f := types.NamespaceFilter{
		Slug: m.Res.Slug,
		Name: m.Res.Name,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("namespaceID=%d", m.Res.ID)
	}

	return f
}
