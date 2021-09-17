package resource

import (
	"fmt"
	"strconv"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
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

		WfRefs    RefSet
		ModRefs   RefSet
		RefCharts RefSet

		BlockRefs map[int]RefSet
	}
)

func NewComposePage(pg *types.Page, nsRef, modRef, parentRef string) *ComposePage {
	r := &ComposePage{
		base:      &base{},
		WfRefs:    make(RefSet, 0, 10),
		ModRefs:   make(RefSet, 0, 10),
		RefCharts: make(RefSet, 0, 10),
		BlockRefs: make(map[int]RefSet),
	}
	r.SetResourceType(types.PageResourceType)
	r.Res = pg

	r.AddIdentifier(identifiers(pg.Handle, pg.Title, pg.ID)...)

	r.RefNs = r.AddRef(types.NamespaceResourceType, nsRef)
	if modRef != "" {
		r.RefMod = r.AddRef(types.ModuleResourceType, modRef).Constraint(r.RefNs)
	}

	if parentRef != "" {
		r.RefParent = r.AddRef(types.PageResourceType, parentRef).Constraint(r.RefNs)
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
				ref := r.AddRef(types.ModuleResourceType, id).Constraint(r.RefNs)
				r.BlockRefs[i] = add(r.BlockRefs[i], ref)
				r.ModRefs = append(r.ModRefs, ref)
			}

		case "Automation":
			bb, _ := b.Options["buttons"].([]interface{})
			for _, b := range bb {
				button, _ := b.(map[string]interface{})
				id := ss(button, "workflow", "workflowID")
				if id != "" {
					ref := r.AddRef(automationTypes.WorkflowResourceType, id)
					r.BlockRefs[i] = add(r.BlockRefs[i], ref)
					r.WfRefs = append(r.WfRefs, ref)
				}
			}

		case "RecordOrganizer":
			id := ss(b.Options, "module", "moduleID")
			if id != "" {
				ref := r.AddRef(types.ModuleResourceType, id).Constraint(r.RefNs)
				r.BlockRefs[i] = add(r.BlockRefs[i], ref)
				r.ModRefs = append(r.ModRefs, ref)
			}

		case "Chart":
			id := ss(b.Options, "chart", "chartID")
			if id != "" {
				ref := r.AddRef(types.ChartResourceType, id).Constraint(r.RefNs)
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
					ref := r.AddRef(types.ModuleResourceType, id).Constraint(r.RefNs)
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
					ref := r.AddRef(types.ModuleResourceType, id).Constraint(r.RefNs)
					r.BlockRefs[i] = add(r.BlockRefs[i], ref)
					r.ModRefs = append(r.ModRefs, ref)
				}
			}

		case "Comment":
			id := ss(b.Options, "module", "moduleID")
			if id != "" {
				ref := r.AddRef(types.ModuleResourceType, id).Constraint(r.RefNs)
				r.BlockRefs[i] = add(r.BlockRefs[i], ref)
				r.ModRefs = append(r.ModRefs, ref)
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

func (r *ComposePage) RBACPath() []*Ref {
	return []*Ref{r.RefNs}
}

func (r *ComposePage) ResourceTranslationParts() (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref()
	path = []*Ref{r.RefNs}
	resource = fmt.Sprintf(types.PageResourceTranslationTpl(), types.PageResourceTranslationType, r.RefNs.Identifiers.First(), firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Handle))

	return
}

func (r *ComposePage) encodeTranslations() ([]*ResourceTranslation, error) {
	return nil, nil
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
