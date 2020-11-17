package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	// ComposePage represents a ComposePage
	ComposePage struct {
		*base
		Res *types.Page

		NsRef  *Ref
		ModRef *Ref
	}
)

func NewComposePage(pg *types.Page, nsRef, modRef string) *ComposePage {
	r := &ComposePage{base: &base{}}
	r.SetResourceType(COMPOSE_PAGE_RESOURCE_TYPE)
	r.Res = pg

	r.AddIdentifier(identifiers(pg.Handle, pg.Title, pg.ID)...)

	r.NsRef = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)
	if modRef != "" {
		r.ModRef = r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, modRef)
	}

	return r
}

func (r *ComposePage) SysID() uint64 {
	return r.Res.ID
}
