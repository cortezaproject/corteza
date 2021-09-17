package store

import (
	"context"
	"strconv"

	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
)

func newComposePageFromResource(res *resource.ComposePage, cfg *EncoderConfig) resourceState {
	return &composePage{
		cfg: mergeConfig(cfg, res.Config()),

		res: res,

		relWfs:    make(map[string]*atypes.Workflow),
		relMods:   make(map[string]*types.Module),
		relCharts: make(map[string]*types.Chart),
	}
}

func (n *composePage) Prepare(ctx context.Context, pl *payload) (err error) {
	// Get related namespace
	if !n.cfg.IgnoreStore {
		n.relNS, err = findComposeNamespace(ctx, pl.s, pl.state.ParentResources, n.res.RefNs.Identifiers)
		if err != nil {
			return err
		}
	} else {
		n.res.Res.ID = 0
		n.relNS = resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefNs.Identifiers)
	}
	if n.relNS == nil {
		return resource.ComposeNamespaceErrUnresolved(n.res.RefNs.Identifiers)
	}

	// Get related module
	// If this isn't a record page, there is no related module
	if n.res.RefMod != nil {
		if !n.cfg.IgnoreStore {
			n.relMod, err = findComposeModule(ctx, pl.s, n.relNS.ID, pl.state.ParentResources, n.res.RefMod.Identifiers)
			if err != nil {
				return err
			}
		} else {
			n.relMod = resource.FindComposeModule(pl.state.ParentResources, n.res.RefMod.Identifiers)
		}
		if n.relMod == nil {
			return resource.ComposeModuleErrUnresolved(n.res.RefMod.Identifiers)
		}
	}

	// Get parent page
	if n.res.RefParent != nil {
		if !n.cfg.IgnoreStore {
			n.relParent, err = findComposePage(ctx, pl.s, n.relNS.ID, pl.state.ParentResources, n.res.RefParent.Identifiers)
			if err != nil {
				return err
			}
		} else {
			n.relParent = resource.FindComposePage(pl.state.ParentResources, n.res.RefParent.Identifiers)
		}
		if n.relParent == nil {
			return resource.ComposePageErrUnresolved(n.res.RefParent.Identifiers)
		}
	}

	// Get related workflows
	var wf *atypes.Workflow
	for _, wfr := range n.res.WfRefs {
		if !n.cfg.IgnoreStore {
			wf, err = findAutomationWorkflow(ctx, pl.s, pl.state.ParentResources, wfr.Identifiers)
			if err != nil {
				return err
			}
		} else {
			wf = resource.FindAutomationWorkflow(pl.state.ParentResources, wfr.Identifiers)
		}
		if wf == nil {
			return resource.AutomationWorkflowErrUnresolved(wfr.Identifiers)
		}
		for id := range wfr.Identifiers {
			n.relWfs[id] = wf
		}
	}

	// Get other related modules
	var mod *types.Module
	for _, mr := range n.res.ModRefs {
		if !n.cfg.IgnoreStore {
			mod, err = findComposeModule(ctx, pl.s, n.relNS.ID, pl.state.ParentResources, mr.Identifiers)
			if err != nil {
				return err
			}
		} else {
			mod = resource.FindComposeModule(pl.state.ParentResources, mr.Identifiers)
		}
		if mod == nil {
			return resource.ComposeModuleErrUnresolved(mr.Identifiers)
		}
		for id := range mr.Identifiers {
			n.relMods[id] = mod
		}
	}

	// Get related charts
	var chr *types.Chart
	for _, refChart := range n.res.RefCharts {
		if !n.cfg.IgnoreStore {
			chr, err = findComposeChart(ctx, pl.s, n.relNS.ID, pl.state.ParentResources, refChart.Identifiers)
			if err != nil {
				return err
			}
		} else {
			chr = resource.FindComposeChart(pl.state.ParentResources, refChart.Identifiers)
		}
		if chr == nil {
			return resource.ComposeChartErrUnresolved(refChart.Identifiers)
		}
		for id := range refChart.Identifiers {
			n.relCharts[id] = chr
		}
	}

	if n.cfg.IgnoreStore {
		return nil
	}

	// Try to get the original page
	n.pg, err = findComposePageStore(ctx, pl.s, n.relNS.ID, makeGenericFilter(n.res.Identifiers()))
	if err != nil {
		return err
	}

	if n.pg != nil {
		n.res.Res.ID = n.pg.ID
		n.res.Res.NamespaceID = n.pg.NamespaceID
	}
	return nil
}

func (n *composePage) Encode(ctx context.Context, pl *payload) (err error) {
	res := n.res.Res
	exists := n.pg != nil && n.pg.ID > 0

	// Determine the ID
	if res.ID <= 0 && exists {
		res.ID = n.pg.ID
	}
	if res.ID <= 0 {
		res.ID = NextID()
	}

	ts := n.res.Timestamps()
	if ts != nil {
		if ts.CreatedAt != nil {
			res.CreatedAt = *ts.CreatedAt.T
		} else {
			res.CreatedAt = *now()
		}
		if ts.UpdatedAt != nil {
			res.UpdatedAt = ts.UpdatedAt.T
		}
		if ts.DeletedAt != nil {
			res.DeletedAt = ts.DeletedAt.T
		}
	}

	// Namespace
	res.NamespaceID = n.relNS.ID
	if res.NamespaceID <= 0 {
		ns := resource.FindComposeNamespace(pl.state.ParentResources, n.res.RefNs.Identifiers)
		res.NamespaceID = ns.ID
	}

	if res.NamespaceID <= 0 {
		return resource.ComposeNamespaceErrUnresolved(n.res.RefNs.Identifiers)
	}

	// Module?
	if n.res.RefMod != nil {
		res.ModuleID = n.relMod.ID
		if res.ModuleID <= 0 {
			mod := resource.FindComposeModule(pl.state.ParentResources, n.res.RefMod.Identifiers)
			res.ModuleID = mod.ID
		}
	}

	// Parent?
	if n.res.RefParent != nil {
		res.SelfID = n.relParent.ID
		if res.SelfID <= 0 {
			mod := resource.FindComposePage(pl.state.ParentResources, n.res.RefParent.Identifiers)
			res.SelfID = mod.ID
		}
	}

	// Blocks
	getModID := func(id string) uint64 {
		mod := n.relMods[id]
		if mod == nil || mod.ID <= 0 {
			mod = resource.FindComposeModule(pl.state.ParentResources, resource.MakeIdentifiers(id))
			if mod == nil || mod.ID <= 0 {
				return 0
			}
		}
		return mod.ID
	}
	getWfID := func(id string) uint64 {
		wf := n.relWfs[id]
		if wf == nil || wf.ID <= 0 {
			wf = resource.FindAutomationWorkflow(pl.state.ParentResources, resource.MakeIdentifiers(id))
			if wf == nil || wf.ID <= 0 {
				return 0
			}
		}
		return wf.ID
	}

	getChartID := func(id string) uint64 {
		chr := n.relCharts[id]
		if chr == nil || chr.ID <= 0 {
			chr = resource.FindComposeChart(pl.state.ParentResources, resource.MakeIdentifiers(id))
			if chr == nil || chr.ID <= 0 {
				return 0
			}
		}
		return chr.ID
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

	for _, b := range res.Blocks {
		switch b.Kind {
		case "RecordList":
			id := ss(b.Options, "module", "moduleID")
			if id == "" {
				continue
			}
			mID := getModID(id)
			if mID <= 0 {
				return resource.ComposeModuleErrUnresolved(resource.MakeIdentifiers(id))
			}
			b.Options["moduleID"] = strconv.FormatUint(mID, 10)
			delete(b.Options, "module")

		case "Automation":
			bb, _ := b.Options["buttons"].([]interface{})
			for _, b := range bb {
				button, _ := b.(map[string]interface{})
				id := ss(button, "workflow", "workflowID")
				if id == "" {
					continue
				}
				wfID := getWfID(id)
				if wfID <= 0 {
					return resource.AutomationWorkflowErrUnresolved(resource.MakeIdentifiers(id))
				}
				button["workflowID"] = strconv.FormatUint(wfID, 10)
				delete(button, "workflow")
			}

		case "RecordOrganizer":
			id := ss(b.Options, "module", "moduleID")
			if id == "" {
				continue
			}
			mID := getModID(id)
			if mID <= 0 {
				return resource.ComposeModuleErrUnresolved(resource.MakeIdentifiers(id))
			}
			b.Options["moduleID"] = strconv.FormatUint(mID, 10)
			delete(b.Options, "module")

		case "Calendar":
			ff, _ := b.Options["feeds"].([]interface{})
			for _, f := range ff {
				feed, _ := f.(map[string]interface{})
				fOpts, _ := (feed["options"]).(map[string]interface{})
				id := ss(fOpts, "module", "moduleID")
				if id == "" {
					continue
				}
				mID := getModID(id)
				if mID <= 0 {
					return resource.ComposeModuleErrUnresolved(resource.MakeIdentifiers(id))
				}
				fOpts["moduleID"] = strconv.FormatUint(mID, 10)
				delete(fOpts, "module")
			}

		case "Chart":
			id := ss(b.Options, "chart", "chartID")
			if id == "" {
				continue
			}
			chrID := getChartID(id)
			if chrID == 0 {
				return resource.ComposeChartErrUnresolved(resource.MakeIdentifiers(id))
			}
			b.Options["chartID"] = strconv.FormatUint(chrID, 10)
			delete(b.Options, "chart")

		case "Metric":
			mm, _ := b.Options["metrics"].([]interface{})
			for _, m := range mm {
				mops, _ := m.(map[string]interface{})
				id := ss(mops, "module", "moduleID")
				if id == "" {
					continue
				}
				mID := getModID(id)
				if mID <= 0 {
					return resource.ComposeModuleErrUnresolved(resource.MakeIdentifiers(id))
				}
				mops["moduleID"] = strconv.FormatUint(mID, 10)
				delete(mops, "module")

			}

		case "Comment":
			id := ss(b.Options, "module", "moduleID")
			if id == "" {
				continue
			}
			mID := getModID(id)
			if mID <= 0 {
				return resource.ComposeModuleErrUnresolved(resource.MakeIdentifiers(id))
			}
			b.Options["moduleID"] = strconv.FormatUint(mID, 10)
			delete(b.Options, "module")
		}
	}

	// Evaluate the resource skip expression
	// @todo expand available parameters; similar implementation to compose/types/record@Dict
	if skip, err := basicSkipEval(ctx, n.cfg, !exists); err != nil {
		return err
	} else if skip {
		return nil
	}

	// Create a fresh page
	if !exists {
		return store.CreateComposePage(ctx, pl.s, res)
	}

	// Update existing page
	switch n.cfg.OnExisting {
	case resource.Skip:
		return nil

	case resource.MergeLeft:
		res = mergeComposePage(n.pg, res)

	case resource.MergeRight:
		res = mergeComposePage(res, n.pg)
	}

	err = store.UpdateComposePage(ctx, pl.s, res)
	if err != nil {
		return err
	}

	n.res.Res = res
	return nil
}
