package resource

import (
	"fmt"
	"strconv"

	automationTypes "github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/compose/types"
)

type (
	// ComposePageLayout represents a ComposePageLayout
	ComposePageLayout struct {
		*base
		Res *types.PageLayout

		RefNs     *Ref
		RefMod    *Ref
		RefParent *Ref

		WfRefs    RefSet
		ModRefs   RefSet
		RefCharts RefSet

		BlockRefs map[int]RefSet
	}
)

func NewComposePageLayout(pg *types.PageLayout, nsRef, modRef, parentRef *Ref) *ComposePageLayout {
	r := &ComposePageLayout{
		base:      &base{},
		WfRefs:    make(RefSet, 0, 10),
		ModRefs:   make(RefSet, 0, 10),
		RefCharts: make(RefSet, 0, 10),
		BlockRefs: make(map[int]RefSet),
	}
	r.SetResourceType(types.PageLayoutResourceType)
	r.Res = pg

	r.AddIdentifier(identifiers(pg.Handle, pg.Meta.Title, pg.ID)...)

	r.RefNs = r.addRef(nsRef)
	if modRef != nil {
		r.RefMod = r.addRef(modRef).Constraint(r.RefNs)
	}

	if parentRef != nil {
		r.RefParent = r.addRef(parentRef).Constraint(r.RefNs)
	}

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDA(&pg.CreatedAt, pg.UpdatedAt, pg.DeletedAt, nil))

	return r
}

func UnpackComposePageLayout(p *types.PageLayout) (*types.PageLayout, *Ref, *Ref) {
	var modRef *Ref
	var parentRef *Ref

	return p, modRef, parentRef
}

func (r *ComposePageLayout) Resource() interface{} {
	return r.Res
}

func (r *ComposePageLayout) ReRef(old RefSet, new RefSet) {
	r.base.ReRef(old, new)

	for _, n := range new {
		switch n.ResourceType {
		case types.NamespaceResourceType:
			r.RefNs = MakeRef(types.NamespaceResourceType, n.Identifiers)
		case types.ModuleResourceType:
			r.RefMod = MakeRef(types.ModuleResourceType, n.Identifiers)
		case types.PageLayoutResourceType:
			r.RefParent = MakeRef(types.PageLayoutResourceType, n.Identifiers)
		}
	}
}

func (r *ComposePageLayout) Prune(ref *Ref) {
}

func (r *ComposePageLayout) SysID() uint64 {
	return r.Res.ID
}

func (r *ComposePageLayout) resourceParts(tpl string) (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref().Constraint(r.RefNs)
	path = []*Ref{r.RefNs}
	resource = fmt.Sprintf(types.PageLayoutRbacResourceTpl(), types.PageLayoutResourceType, r.RefNs.Identifiers.First(), firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Handle))

	return
}

func (r *ComposePageLayout) RBACParts() (resource string, ref *Ref, path []*Ref) {
	return r.resourceParts(types.PageLayoutRbacResourceTpl())
}

func (r *ComposePageLayout) ResourceTranslationParts() (resource string, ref *Ref, path []*Ref) {
	return r.resourceParts(types.PageLayoutResourceTranslationTpl())
}

func (r *ComposePageLayout) encodeTranslations() ([]*ResourceTranslation, error) {
	return nil, nil
}

// FindComposePageLayout looks for the page in the resources
func FindComposePageLayout(rr InterfaceSet, ii Identifiers) (pg *types.PageLayout) {
	var pgRes *ComposePageLayout

	rr.Walk(func(r Interface) error {
		pr, ok := r.(*ComposePageLayout)
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

func ComposePageLayoutErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("compose page unresolved %v", ii.StringSlice())
}

// PageLayout block utilities
func (r *ComposePageLayout) removeBlock(i int) {
	// do the swap remove thing
	r.Res.Blocks[i] = r.Res.Blocks[len(r.Res.Blocks)-1]
	r.Res.Blocks = r.Res.Blocks[:len(r.Res.Blocks)-1]

	// correct block refs.
	// +1 here because it's already cut
	r.BlockRefs[i] = r.BlockRefs[len(r.Res.Blocks)]
	delete(r.BlockRefs, len(r.Res.Blocks))
}

func (r *ComposePageLayout) optString(opt map[string]interface{}, kk ...string) string {
	for _, k := range kk {
		if vr, has := opt[k]; has {
			v, _ := vr.(string)
			return v
		}
	}
	return ""
}

func (r *ComposePageLayout) pbRecordList(opt map[string]interface{}) (out *Ref) {
	id := r.optString(opt, "module", "moduleID")
	if id == "" || id == "0" {
		return
	}

	return MakeRef(types.ModuleResourceType, MakeIdentifiers(id)).Constraint(r.RefNs)
}

func (r *ComposePageLayout) pbComment(opt map[string]interface{}) (out *Ref) {
	id := r.optString(opt, "module", "moduleID")
	if id == "" || id == "0" {
		return
	}

	return MakeRef(types.ModuleResourceType, MakeIdentifiers(id)).Constraint(r.RefNs)
}

func (r *ComposePageLayout) pbAutomation(opt map[string]interface{}) (out *Ref) {
	id := r.optString(opt, "workflow", "workflowID")
	if id == "" || id == "0" {
		return
	}

	return MakeRef(automationTypes.WorkflowResourceType, MakeIdentifiers(id))
}

func (r *ComposePageLayout) pbRecordOrganizer(opt map[string]interface{}) (out *Ref) {
	id := r.optString(opt, "module", "moduleID")
	if id == "" || id == "0" {
		return
	}

	return MakeRef(types.ModuleResourceType, MakeIdentifiers(id)).Constraint(r.RefNs)
}

func (r *ComposePageLayout) pbChart(opt map[string]interface{}) (out *Ref) {
	id := r.optString(opt, "chart", "chartID")
	if id == "" || id == "0" {
		return
	}

	return MakeRef(types.ChartResourceType, MakeIdentifiers(id)).Constraint(r.RefNs)
}

func (r *ComposePageLayout) pbCalendar(opt map[string]interface{}) (out *Ref) {
	id := r.optString(opt, "module", "moduleID")
	if id == "" || id == "0" {
		return
	}

	return MakeRef(types.ModuleResourceType, MakeIdentifiers(id)).Constraint(r.RefNs)
}

func (r *ComposePageLayout) pbMetric(opt map[string]interface{}) (out *Ref) {
	id := r.optString(opt, "module", "moduleID")
	if id == "" || id == "0" {
		return
	}

	return MakeRef(types.ModuleResourceType, MakeIdentifiers(id)).Constraint(r.RefNs)
}
