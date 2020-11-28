package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	// ComposePage represents a ComposePage
	ComposePage struct {
		*base
		Res *types.Page

		NsRef     *Ref
		ModRef    *Ref
		ParentRef *Ref

		ModRefs   RefSet
		ChartRefs RefSet
	}
)

func NewComposePage(pg *types.Page, nsRef, modRef, parentRef string) *ComposePage {
	r := &ComposePage{
		base:      &base{},
		ModRefs:   make(RefSet, 0, 10),
		ChartRefs: make(RefSet, 0, 10),
	}
	r.SetResourceType(COMPOSE_PAGE_RESOURCE_TYPE)
	r.Res = pg

	r.AddIdentifier(identifiers(pg.Handle, pg.Title, pg.ID)...)

	r.NsRef = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)
	if modRef != "" {
		r.ModRef = r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, modRef)
	}

	if parentRef != "" {
		r.ParentRef = r.AddRef(COMPOSE_PAGE_RESOURCE_TYPE, parentRef)
	}

	for _, b := range pg.Blocks {
		switch b.Kind {
		case "RecordList":
			id, _ := b.Options["module"].(string)
			if id != "" {
				r.ModRefs = append(r.ModRefs, r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, id))
			}

		case "Chart":
			id, _ := b.Options["chart"].(string)
			if id != "" {
				r.ChartRefs = append(r.ChartRefs, r.AddRef(COMPOSE_CHART_RESOURCE_TYPE, id))
			}

		case "Calendar":
			ff, _ := b.Options["feeds"].([]interface{})
			for _, f := range ff {
				feed, _ := f.(map[string]interface{})
				fOpts, _ := (feed["options"]).(map[string]interface{})
				id, _ := fOpts["module"].(string)
				if id != "" {
					r.ModRefs = append(r.ModRefs, r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, id))
				}
			}

		case "Metric":
			mm, _ := b.Options["metrics"].([]interface{})
			for _, m := range mm {
				mops, _ := m.(map[string]interface{})
				id, _ := mops["module"].(string)
				if id != "" {
					r.ModRefs = append(r.ModRefs, r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, id))
				}
			}
		}
	}

	return r
}

func (r *ComposePage) SysID() uint64 {
	return r.Res.ID
}
