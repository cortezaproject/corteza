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
	composePage struct {
		*base
		Res *types.Page

		// Might keep track of related namespace, page
	}
)

func ComposePage(pg *types.Page) *composePage {
	r := &composePage{base: &base{}}
	r.SetResourceType(COMPOSE_PAGE_RESOURCE_TYPE)
	r.Res = pg

	r.AddIdentifier(identifiers(pg.Handle, pg.Title, pg.ID)...)

	return r
}

func (m *composePage) SearchQuery() types.PageFilter {
	f := types.PageFilter{
		Handle: m.Res.Handle,
	}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("pageID=%d", m.Res.ID)
	}

	return f
}
