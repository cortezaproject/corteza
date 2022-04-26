package store

import (
	"context"
	"fmt"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"golang.org/x/text/language"
)

var (
	gResourceTranslations map[language.Tag]map[string]*types.ResourceTranslation

	resourceTranslationIndex = func(l *types.ResourceTranslation) string {
		return fmt.Sprintf("%s/%s", l.Resource, l.K)
	}
)

func newResourceTranslationFromResource(res *resource.ResourceTranslation, cfg *EncoderConfig) resourceState {
	return &resourceTranslation{
		cfg: mergeConfig(cfg, res.Config()),

		locales: res.Res,

		refResourceTranslation: res.Res[0].Resource,
		refLocaleRes:           res.RefRes,

		refPathRes: res.RefPath,
	}
}

func (n *resourceTranslation) Prepare(ctx context.Context, pl *payload) (err error) {
	// Init global state
	if gResourceTranslations == nil {
		err = n.initGlobalIndex(ctx, pl.s)
		if err != nil {
			return err
		}
	}

	// For now we will only allow resource specific locale resources if that resource is
	// also present.
	// Here we check if we can find it in case we're handling a resource specific rule.
	refRes := n.refLocaleRes
	if refRes.ResourceType == composeTypes.ModuleFieldResourceType {
		for _, r := range pl.state.ParentResources {
			if r.ResourceType() == composeTypes.ModuleResourceType && r.Identifiers().HasAny(n.refPathRes[1].Identifiers) {
				return
			}
		}
	} else {
		for _, r := range pl.state.ParentResources {
			if refRes.ResourceType == r.ResourceType() && r.Identifiers().HasAny(refRes.Identifiers) {
				return
			}
		}
		// We couldn't find it...
		return resource.ResourceTranslationErrNotFound(refRes.Identifiers)
	}

	return
}

func (n *resourceTranslation) Encode(ctx context.Context, pl *payload) (err error) {

	// @todo move out of the encoding logic
	if n == nil {
		return
	}

	localeRes, err := n.makeResourceTranslation(pl)
	if err != nil {
		return err
	}

	// Assure correct resources
	for _, l := range n.locales {
		l.Resource = localeRes
	}

	ll := make(types.ResourceTranslationSet, 0, len(n.locales))
	for _, l := range n.locales {
		old, ok := gResourceTranslations[l.Lang.Tag][resourceTranslationIndex(l)]

		// not exist
		if !ok {
			l.ID = NextID()
			l.CreatedAt = *now()
			l.CreatedBy = pl.invokerID

			ll = append(ll, l)
			continue
		}

		// exists
		switch n.cfg.OnExisting {
		case resource.Skip,
			resource.MergeLeft:
			continue

		case resource.MergeRight,
			resource.Default:
			old.Message = l.Message

			ll = append(ll, old)
		}
	}

	// Update global index for following runs
	n.updateGlobalIndex(ll)

	return store.UpsertResourceTranslation(ctx, pl.s, ll...)
}

func (n *resourceTranslation) makeResourceTranslation(pl *payload) (string, error) {
	p0ID := uint64(0)
	p1ID := uint64(0)
	p2ID := uint64(0)

	localeRes := ""
	_ = localeRes

	switch n.refLocaleRes.ResourceType {
	// case automationTypes.WorkflowResourceType:
	// 	p1 := resource.FindAutomationWorkflow(pl.state.ParentResources, n.refLocaleRes.Identifiers)
	// 	if p1 == nil {
	// 		return "", resource.AutomationWorkflowErrUnresolved(n.refLocaleRes.Identifiers)
	// 	}
	// 	p1ID = p1.ID

	// 	return automationTypes.WorkflowResourceTranslation(p1ID), nil

	case composeTypes.NamespaceResourceType:
		p1 := resource.FindComposeNamespace(pl.state.ParentResources, n.refLocaleRes.Identifiers)
		if p1 == nil {
			return "", resource.ComposeNamespaceErrUnresolved(n.refLocaleRes.Identifiers)
		}
		p1ID = p1.ID

		return composeTypes.NamespaceResourceTranslation(p1ID), nil

	case composeTypes.ModuleResourceType:
		p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.refPathRes[0].Identifiers)
		if p0 == nil {
			return "", resource.ComposeNamespaceErrUnresolved(n.refPathRes[0].Identifiers)
		}
		p0ID = p0.ID

		p1 := resource.FindComposeModule(pl.state.ParentResources, n.refLocaleRes.Identifiers)
		if p1 == nil {
			return "", resource.ComposeModuleErrUnresolved(n.refLocaleRes.Identifiers)
		}
		p1ID = p1.ID

		return composeTypes.ModuleResourceTranslation(p0ID, p1ID), nil

	case composeTypes.PageResourceType:
		p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.refPathRes[0].Identifiers)
		if p0 == nil {
			return "", resource.ComposeNamespaceErrUnresolved(n.refPathRes[0].Identifiers)
		}
		p0ID = p0.ID

		p1 := resource.FindComposePage(pl.state.ParentResources, n.refLocaleRes.Identifiers)
		if p1 == nil {
			return "", resource.ComposePageErrUnresolved(n.refLocaleRes.Identifiers)
		}
		p1ID = p1.ID

		return composeTypes.PageResourceTranslation(p0ID, p1ID), nil

	case composeTypes.ChartResourceType:
		p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.refPathRes[0].Identifiers)
		if p0 == nil {
			return "", resource.ComposeNamespaceErrUnresolved(n.refPathRes[0].Identifiers)
		}
		p0ID = p0.ID

		p1 := resource.FindComposeChart(pl.state.ParentResources, n.refLocaleRes.Identifiers)
		if p1 == nil {
			return "", resource.ComposeChartErrUnresolved(n.refLocaleRes.Identifiers)
		}
		p1ID = p1.ID

		return composeTypes.ChartResourceTranslation(p0ID, p1ID), nil

	case composeTypes.ModuleFieldResourceType:
		p0 := resource.FindComposeNamespace(pl.state.ParentResources, n.refPathRes[0].Identifiers)
		if p0 == nil {
			return "", resource.ComposeNamespaceErrUnresolved(n.refPathRes[0].Identifiers)
		}
		p0ID = p0.ID

		p1 := resource.FindComposeModule(pl.state.ParentResources, n.refPathRes[1].Identifiers)
		if p1 == nil {
			return "", resource.ComposeModuleErrUnresolved(n.refPathRes[1].Identifiers)
		}
		p1ID = p1.ID

		// field
		f := resource.FindComposeModuleField(pl.state.ParentResources, n.refPathRes[1].Identifiers, n.refLocaleRes.Identifiers)
		if f == nil {
			return "", resource.ComposeModuleFieldErrUnresolved(n.refLocaleRes.Identifiers)
		}
		p2ID = f.ID

		return composeTypes.ModuleFieldResourceTranslation(p0ID, p1ID, p2ID), nil

	// case types.ReportResourceType:
	// 	p1 := resource.FindReport(pl.state.ParentResources, n.refLocaleRes.Identifiers)
	// 	if p1 == nil {
	// 		return "", resource.ReportErrUnresolved(n.refLocaleRes.Identifiers)
	// 	}
	// 	p1ID = p1.ID

	// 	return types.ReportResourceTranslation(p1ID), nil

	default:
		// @todo if we wish to support res. trans. for external stuff, this needs to pass through.
		//       this also requires some tweaks in the path ID thing.
		return "", fmt.Errorf("unsupported resource type '%s' for resource translation store encode", n.refLocaleRes.ResourceType)
	}
}

func (n *resourceTranslation) initGlobalIndex(ctx context.Context, s store.Storer) (err error) {
	gResourceTranslations = make(map[language.Tag]map[string]*types.ResourceTranslation)
	ll, _, err := store.SearchResourceTranslations(ctx, s, types.ResourceTranslationFilter{})
	if err != store.ErrNotFound && err != nil {
		return err
	}

	n.updateGlobalIndex(ll)
	return nil
}

func (n *resourceTranslation) updateGlobalIndex(ll types.ResourceTranslationSet) {
	for _, l := range ll {
		if _, ok := gResourceTranslations[l.Lang.Tag]; !ok {
			gResourceTranslations[l.Lang.Tag] = make(map[string]*types.ResourceTranslation)
		}
		gResourceTranslations[l.Lang.Tag][resourceTranslationIndex(l)] = l
	}
}
