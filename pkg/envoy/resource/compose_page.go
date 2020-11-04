package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
)

const (
	COMPOSE_PAGE_RESOURCE_TYPE = "composePage"
)

type (
	// ComposePage represents a ComposePage
	ComposePage struct {
		*base
		Res *types.Page

		// Might keep track of related namespace, page

	}
)

func NewComposePage(pg *types.Page, nsRef, modRef string) *ComposePage {
	r := &ComposePage{base: &base{}}
	r.SetResourceType(COMPOSE_PAGE_RESOURCE_TYPE)
	r.Res = pg

	r.AddIdentifier(identifiers(pg.Handle, pg.Title, pg.ID)...)

	r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)
	r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, modRef)

	return r
}

func (m *ComposePage) SearchQuery() types.PageFilter {
	f := types.PageFilter{
		Handle: m.Res.Handle,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("pageID=%d", m.Res.ID)
	}

	return f
}
