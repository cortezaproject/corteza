package resource

import (
	"fmt"
	"strconv"

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
	r.SetResourceType("compose:namespace")
	r.Res = ns

	if ns.Slug != "" {
		r.AddIdentifier(ns.Slug)
	}
	if ns.Name != "" {
		r.AddIdentifier(ns.Name)
	}
	if ns.ID > 0 {
		r.AddIdentifier(strconv.FormatUint(ns.ID, 10))
	}
	return r
}

func (m *ComposeNamespace) SearchQuery() types.NamespaceFilter {
	f := types.NamespaceFilter{Query: ""}

	f.Slug = m.Res.Slug
	f.Name = m.Res.Name
	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("namespaceID=%d", m.Res.ID)
	}

	return f
}
