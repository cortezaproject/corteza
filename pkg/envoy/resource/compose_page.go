package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	// ComposePage represents a ComposePage
	ComposePage struct {
		*base
		Res *types.Page

		RefNs     *Ref
		RefMod    *Ref
		RefParent *Ref

		ModRefs   RefSet
		RefCharts RefSet

		BlockRefs map[int]RefSet
	}
)

func NewComposePage(pg *types.Page, nsRef, modRef, parentRef string) *ComposePage {
	r := &ComposePage{
		base:      &base{},
		ModRefs:   make(RefSet, 0, 10),
		RefCharts: make(RefSet, 0, 10),
		BlockRefs: make(map[int]RefSet),
	}
	r.SetResourceType(COMPOSE_PAGE_RESOURCE_TYPE)
	r.Res = pg

	r.AddIdentifier(identifiers(pg.Handle, pg.Title, pg.ID)...)

	r.RefNs = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)
	if modRef != "" {
		r.RefMod = r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, modRef).Constraint(r.RefNs)
	}

	if parentRef != "" {
		r.RefParent = r.AddRef(COMPOSE_PAGE_RESOURCE_TYPE, parentRef).Constraint(r.RefNs)
	}

	// Quick utility to extract references from options
	ss := func(m map[string]interface{}, kk ...string) string {
		for _, k := range kk {
			if vr, has := m[k]; has {
				v, _ := vr.(string)
				return v
			}
		}
		return ""
	}

	add := func(rr RefSet, r *Ref) RefSet {
		if rr == nil {
			rr = make(RefSet, 0, 2)
		}

		return append(rr, r)
	}

	for i, b := range pg.Blocks {
		switch b.Kind {
		case "RecordList":
			id := ss(b.Options, "module", "moduleID")
			if id != "" {
				ref := r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, id).Constraint(r.RefNs)
				r.BlockRefs[i] = add(r.BlockRefs[i], ref)
				r.ModRefs = append(r.ModRefs, ref)
			}

		case "RecordOrganizer":
			id := ss(b.Options, "module", "moduleID")
			if id != "" {
				ref := r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, id).Constraint(r.RefNs)
				r.BlockRefs[i] = add(r.BlockRefs[i], ref)
				r.ModRefs = append(r.ModRefs, ref)
			}

		case "Chart":
			id := ss(b.Options, "chart", "chartID")
			if id != "" {
				ref := r.AddRef(COMPOSE_CHART_RESOURCE_TYPE, id).Constraint(r.RefNs)
				r.BlockRefs[i] = add(r.BlockRefs[i], ref)
				r.RefCharts = append(r.RefCharts, ref)
			}

		case "Calendar":
			ff, _ := b.Options["feeds"].([]interface{})
			for _, f := range ff {
				feed, _ := f.(map[string]interface{})
				fOpts, _ := (feed["options"]).(map[string]interface{})
				id := ss(fOpts, "module", "moduleID")
				if id != "" {
					ref := r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, id).Constraint(r.RefNs)
					r.BlockRefs[i] = add(r.BlockRefs[i], ref)
					r.ModRefs = append(r.ModRefs, ref)
				}
			}

		case "Metric":
			mm, _ := b.Options["metrics"].([]interface{})
			for _, m := range mm {
				mops, _ := m.(map[string]interface{})
				id := ss(mops, "module", "moduleID")
				if id != "" {
					ref := r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, id).Constraint(r.RefNs)
					r.BlockRefs[i] = add(r.BlockRefs[i], ref)
					r.ModRefs = append(r.ModRefs, ref)
				}
			}
		}
	}

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDA(&pg.CreatedAt, pg.UpdatedAt, pg.DeletedAt, nil))

	return r
}

func (r *ComposePage) SysID() uint64 {
	return r.Res.ID
}

func (r *ComposePage) Ref() string {
	return firstOkString(r.Res.Handle, r.Res.Title, strconv.FormatUint(r.Res.ID, 10))
}

// FindComposePage looks for the page in the resources
func FindComposePage(rr InterfaceSet, ii Identifiers) (pg *types.Page) {
	var pgRes *ComposePage

	rr.Walk(func(r Interface) error {
		pr, ok := r.(*ComposePage)
		if !ok {
			return nil
		}

		if pr.Identifiers().HasAny(ii) {
			pgRes = pr
		}
		return nil
	})

	// Found it
	if pgRes != nil {
		return pgRes.Res
	}
	return nil
}

func ComposePageErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("compose page unresolved %v", ii.StringSlice())
}
