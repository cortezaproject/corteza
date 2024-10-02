package rest

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	automationTypes "github.com/cortezaproject/corteza/server/automation/types"
	composeEnvoy "github.com/cortezaproject/corteza/server/compose/envoy"
	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/service/event"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/corredor"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	systemEnvoy "github.com/cortezaproject/corteza/server/system/envoy"
	systemService "github.com/cortezaproject/corteza/server/system/service"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
)

type (
	namespacePayload struct {
		*types.Namespace

		CanGrant           bool `json:"canGrant"`
		CanExportNamespace bool `json:"canExportNamespace"`
		CanUpdateNamespace bool `json:"canUpdateNamespace"`
		CanDeleteNamespace bool `json:"canDeleteNamespace"`
		CanManageNamespace bool `json:"canManageNamespace"`
		CanCreateModule    bool `json:"canCreateModule"`
		CanExportModule    bool `json:"canExportModule"`
		CanCreateChart     bool `json:"canCreateChart"`
		CanExportChart     bool `json:"canExportChart"`
		CanCreatePage      bool `json:"canCreatePage"`
		CanExportPage      bool `json:"canExportPage"`
	}

	namespaceSetPayload struct {
		Filter types.NamespaceFilter `json:"filter"`
		Set    []*namespacePayload   `json:"set"`
	}

	pageFinder interface {
		Find(ctx context.Context, filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error)
	}

	pageLayoutFinder interface {
		Find(ctx context.Context, filter types.PageLayoutFilter) (set types.PageLayoutSet, f types.PageLayoutFilter, err error)
	}

	chartFinder interface {
		Find(ctx context.Context, filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error)
	}

	Namespace struct {
		namespace  service.NamespaceService
		module     service.ModuleService
		page       pageFinder
		pageLayout pageLayoutFinder
		chart      chartFinder
		locale     service.ResourceTranslationsManagerService
		attachment service.AttachmentService
		role       systemService.RoleService
		ac         namespaceAccessController
	}

	namespaceAccessController interface {
		CanGrant(context.Context) bool

		CanExportNamespace(context.Context, *types.Namespace) bool
		CanUpdateNamespace(context.Context, *types.Namespace) bool
		CanDeleteNamespace(context.Context, *types.Namespace) bool
		CanManageNamespace(context.Context, *types.Namespace) bool

		CanCreateModuleOnNamespace(context.Context, *types.Namespace) bool
		CanExportModulesOnNamespace(context.Context, *types.Namespace) bool
		CanCreateChartOnNamespace(context.Context, *types.Namespace) bool
		CanExportChartsOnNamespace(context.Context, *types.Namespace) bool
		CanCreatePageOnNamespace(context.Context, *types.Namespace) bool
		CanExportPagesOnNamespace(context.Context, *types.Namespace) bool
	}
)

func (Namespace) New() *Namespace {
	return &Namespace{
		namespace:  service.DefaultNamespace,
		module:     service.DefaultModule,
		page:       service.DefaultPage,
		pageLayout: service.DefaultPageLayout,
		chart:      service.DefaultChart,
		locale:     service.DefaultResourceTranslation,
		role:       systemService.DefaultRole,
		attachment: service.DefaultAttachment,
		ac:         service.DefaultAccessControl,
	}
}

func (ctrl Namespace) List(ctx context.Context, r *request.NamespaceList) (interface{}, error) {
	var (
		err error
		f   = types.NamespaceFilter{
			Query:  r.Query,
			Slug:   r.Slug,
			Labels: r.Labels,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.namespace.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Namespace) Create(ctx context.Context, r *request.NamespaceCreate) (interface{}, error) {
	var (
		err error
		ns  = &types.Namespace{
			Name:    r.Name,
			Slug:    r.Slug,
			Enabled: r.Enabled,
			Labels:  r.Labels,
		}
	)

	if err = r.Meta.Unmarshal(&ns.Meta); err != nil {
		return nil, err
	}

	ns, err = ctrl.namespace.Create(ctx, ns)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Namespace) Read(ctx context.Context, r *request.NamespaceRead) (interface{}, error) {
	ns, err := ctrl.namespace.FindByID(ctx, r.NamespaceID)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Namespace) ListTranslations(ctx context.Context, r *request.NamespaceListTranslations) (interface{}, error) {
	return ctrl.locale.Namespace(ctx, r.NamespaceID)
}

func (ctrl Namespace) UpdateTranslations(ctx context.Context, r *request.NamespaceUpdateTranslations) (interface{}, error) {
	return api.OK(), ctrl.locale.Upsert(ctx, r.Translations)
}

func (ctrl Namespace) Update(ctx context.Context, r *request.NamespaceUpdate) (interface{}, error) {
	var (
		err error
		ns  = &types.Namespace{
			ID:        r.NamespaceID,
			Name:      r.Name,
			Slug:      r.Slug,
			Enabled:   r.Enabled,
			Labels:    r.Labels,
			UpdatedAt: r.UpdatedAt,
		}
	)

	if err = r.Meta.Unmarshal(&ns.Meta); err != nil {
		return nil, err
	}

	ns, err = ctrl.namespace.Update(ctx, ns)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Namespace) Delete(ctx context.Context, r *request.NamespaceDelete) (interface{}, error) {
	_, err := ctrl.namespace.FindByID(ctx, r.NamespaceID)
	if err != nil {
		return nil, err
	}

	return api.OK(), ctrl.namespace.DeleteByID(ctx, r.NamespaceID)
}

func (ctrl Namespace) Upload(ctx context.Context, r *request.NamespaceUpload) (interface{}, error) {
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	a, err := ctrl.attachment.CreateNamespaceAttachment(
		ctx,
		r.Upload.Filename,
		r.Upload.Size,
		file,
	)
	if err != nil {
		return nil, err
	}

	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl Namespace) Clone(ctx context.Context, r *request.NamespaceClone) (interface{}, error) {
	dup := &types.Namespace{
		Name: r.Name,
		Slug: r.Slug,
	}

	// @todo temporary workaround cause Envoy requires some identifiable thing
	if dup.Slug == "" {
		dup.Slug = fmt.Sprintf("cl_%d", r.NamespaceID)
	}

	nodes, err := ctrl.gatherNodes(ctx, r.NamespaceID)
	if err != nil {
		return nil, err
	}

	decoder := func() (envoyx.NodeSet, error) {
		return nodes, nil
	}

	ns, err := ctrl.namespace.Clone(ctx, r.NamespaceID, dup, decoder)
	if err != nil {
		return nil, err
	}

	// @todo temporary workaround cause Envoy requires some identifiable thing
	if r.Slug == "" {
		ns.Slug = ""
		ns, err = ctrl.namespace.Update(ctx, ns)
		if err != nil {
			return nil, err
		}
	}
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Namespace) Export(ctx context.Context, r *request.NamespaceExport) (out interface{}, err error) {
	nodes, err := ctrl.gatherNodes(ctx, r.NamespaceID)
	if err != nil {
		return
	}

	p := envoyx.EncodeParams{
		Type:   envoyx.EncodeTypeIo,
		Params: map[string]any{},
	}

	evsvc := envoyx.Global()
	gg, err := evsvc.Bake(ctx, p, nil, nodes...)
	if err != nil {
		return
	}

	// Archive encoded resources
	buf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(buf)

	f, err := zw.Create(fmt.Sprintf("%s.yaml", r.Filename))
	if err != nil {
		return
	}

	p.Params["writer"] = f
	err = evsvc.Encode(ctx, p, gg)
	if err != nil {
		return
	}

	err = zw.Close()
	if err != nil {
		return
	}

	return ctrl.serveExport(ctx, fmt.Sprintf("%s.zip", r.Filename), bytes.NewReader(buf.Bytes()), nil)
}

func (ctrl Namespace) ImportInit(ctx context.Context, r *request.NamespaceImportInit) (interface{}, error) {
	f, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ctrl.namespace.ImportInit(ctx, f, r.Upload.Size)
	// return ctrl.namespace.ImportInit(ctx, f, r.Upload.Header.Get("content-type"), r.Upload.Size)
}

func (ctrl Namespace) ImportRun(ctx context.Context, r *request.NamespaceImportRun) (interface{}, error) {
	var (
		dup = &types.Namespace{
			Name: r.Name,
			Slug: r.Slug,
		}
	)

	// @todo temporary workaround cause Envoy requires some identifiable thing
	if dup.Slug == "" {
		dup.Slug = fmt.Sprintf("cl_%d", r.SessionID)
	}

	ns, err := ctrl.namespace.ImportRun(ctx, r.SessionID, dup)
	if err != nil {
		return nil, err
	}

	// @todo temporary workaround cause Envoy requires some identifiable thing
	if r.Slug == "" {
		ns.Slug = ""
		ns, err = ctrl.namespace.Update(ctx, ns)
		if err != nil {
			return nil, err
		}
	}
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Namespace) serveExport(ctx context.Context, fn string, archive io.ReadSeeker, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Disposition", "attachment; filename="+fn)

		http.ServeContent(w, req, fn, time.Now(), archive)
	}, nil
}

func (ctrl *Namespace) TriggerScript(ctx context.Context, r *request.NamespaceTriggerScript) (rsp interface{}, err error) {
	var (
		namespace *types.Namespace
	)

	if namespace, err = ctrl.namespace.FindByID(ctx, r.NamespaceID); err != nil {
		return
	}

	err = corredor.Service().Exec(ctx, r.Script, corredor.ExtendScriptArgs(event.NamespaceOnManual(namespace, nil), r.Args))
	return ctrl.makePayload(ctx, namespace, err)
}

func (ctrl Namespace) makePayload(ctx context.Context, ns *types.Namespace, err error) (*namespacePayload, error) {
	if err != nil || ns == nil {
		return nil, err
	}

	return &namespacePayload{
		Namespace: ns,

		CanGrant:           ctrl.ac.CanGrant(ctx),
		CanExportNamespace: ctrl.ac.CanExportNamespace(ctx, ns),
		CanUpdateNamespace: ctrl.ac.CanUpdateNamespace(ctx, ns),
		CanDeleteNamespace: ctrl.ac.CanDeleteNamespace(ctx, ns),
		CanManageNamespace: ctrl.ac.CanManageNamespace(ctx, ns),

		CanCreateModule: ctrl.ac.CanCreateModuleOnNamespace(ctx, ns),
		CanExportModule: ctrl.ac.CanExportModulesOnNamespace(ctx, ns),
		CanCreateChart:  ctrl.ac.CanCreateChartOnNamespace(ctx, ns),
		CanExportChart:  ctrl.ac.CanExportChartsOnNamespace(ctx, ns),
		CanCreatePage:   ctrl.ac.CanCreatePageOnNamespace(ctx, ns),
		CanExportPage:   ctrl.ac.CanExportPagesOnNamespace(ctx, ns),
	}, nil
}

func (ctrl Namespace) makeFilterPayload(ctx context.Context, nn types.NamespaceSet, f types.NamespaceFilter, err error) (*namespaceSetPayload, error) {
	if err != nil {
		return nil, err
	}

	nsp := &namespaceSetPayload{Filter: f, Set: make([]*namespacePayload, len(nn))}

	for i := range nn {
		nsp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return nsp, nil
}

func (ctrl Namespace) gatherNodes(ctx context.Context, namespaceID uint64) (resources envoyx.NodeSet, err error) {
	var (
		nsII envoyx.Identifiers
		aux  envoyx.NodeSet
	)

	// Prepare resources
	aux, nsII, err = ctrl.exportCompose(ctx, namespaceID)
	if err != nil {
		return
	}
	resources = append(resources, aux...)

	// Tweak exported resources
	resources = ctrl.tweakExport(ctx, resources, nsII)

	// Role placeholders for RBAC
	aux, err = ctrl.preparePlaceholders(ctx)
	if err != nil {
		return
	}
	resources = append(resources, aux...)

	// RBAC
	aux, err = ctrl.exportRBAC(ctx, resources)
	if err != nil {
		return
	}
	resources = append(resources, aux...)

	// Translations
	aux, err = ctrl.exportResourceTranslations(ctx, resources)
	if err != nil {
		return
	}
	resources = append(resources, aux...)

	return
}

func (ctrl Namespace) exportCompose(ctx context.Context, namespaceID uint64) (resources envoyx.NodeSet, nsII envoyx.Identifiers, err error) {
	// - namespace
	n, err := ctrl.namespace.FindByID(ctx, namespaceID)
	if err != nil {
		return
	}

	// @todo this isn't ok, will do for now
	if !ctrl.ac.CanExportNamespace(ctx, n) {
		err = fmt.Errorf("not allowed to export namespace %s", n.Name)
		return
	}

	nsNode, err := composeEnvoy.NamespaceToEnvoyNode(n)
	if err != nil {
		return
	}
	nsII = nsNode.Identifiers
	resources = append(resources, nsNode)

	// - modules
	mm, _, err := ctrl.module.Find(ctx, types.ModuleFilter{NamespaceID: n.ID})
	if err != nil {
		return
	}
	for _, m := range mm {
		var aux *envoyx.Node
		aux, err = composeEnvoy.ModuleToEnvoyNode(m)
		if err != nil {
			return
		}
		resources = append(resources, aux)

		for _, f := range m.Fields {
			aux, err = composeEnvoy.ModuleFieldToEnvoyNode(f)
			if err != nil {
				return
			}
			resources = append(resources, aux)
		}
	}

	// - pages
	pp, _, err := ctrl.page.Find(ctx, types.PageFilter{NamespaceID: n.ID})
	if err != nil {
		return
	}
	for _, p := range pp {
		var aux *envoyx.Node
		aux, err = composeEnvoy.PageToEnvoyNode(p)
		if err != nil {
			return
		}
		resources = append(resources, aux)
	}

	// - page layouts
	ll, _, err := ctrl.pageLayout.Find(ctx, types.PageLayoutFilter{NamespaceID: n.ID})
	if err != nil {
		return
	}
	for _, l := range ll {
		var aux *envoyx.Node
		aux, err = composeEnvoy.PageLayoutToEnvoyNode(l)
		if err != nil {
			return
		}
		resources = append(resources, aux)
	}

	// - charts
	cc, _, err := ctrl.chart.Find(ctx, types.ChartFilter{NamespaceID: n.ID})
	if err != nil {
		return
	}
	for _, c := range cc {
		var aux *envoyx.Node
		aux, err = composeEnvoy.ChartToEnvoyNode(c)
		if err != nil {
			return
		}
		resources = append(resources, aux)
	}

	return
}

func (ctrl Namespace) exportRBAC(ctx context.Context, base envoyx.NodeSet) (resources envoyx.NodeSet, err error) {
	// Prepare RBAC Rules
	rawRules := rbac.Global().Rules()

	resources, err = envoyx.RBACRulesForNodes(rawRules, base...)
	if err != nil {
		return
	}

	return
}

func (ctrl Namespace) exportResourceTranslations(ctx context.Context, base envoyx.NodeSet) (resources envoyx.NodeSet, err error) {
	var (
		lsvc         = locale.Global()
		tags         = lsvc.Tags()
		translations = make([]*locale.ResourceTranslation, 0, 128)

		resKeyTrans map[string]map[string]*locale.ResourceTranslation
	)

	for _, t := range tags {
		resKeyTrans, err = lsvc.LoadResourceTranslations(ctx, t)
		if err != nil {
			return
		}

		for _, keyTrans := range resKeyTrans {
			for _, trans := range keyTrans {
				translations = append(translations, trans)
			}
		}
	}

	resources, err = envoyx.ResourceTranslationsForNodes(systemTypes.FromLocale(translations), base...)
	return
}

func (ctrl Namespace) tweakExport(ctx context.Context, nodes envoyx.NodeSet, nsII envoyx.Identifiers) envoyx.NodeSet {
	nsRef := envoyx.Ref{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  nsII,
		Scope: envoyx.Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  nsII,
		},
	}
	nsNode := envoyx.NodeForRef(nsRef, nodes...)

	// - remove logo and icon references as attachments are not exported by default
	// @todo code in attachment exporting, most likely when we do attachment handling rework
	ns := nsNode.Resource.(*types.Namespace)
	ns.Meta.Icon = ""
	ns.Meta.IconID = 0
	ns.Meta.Logo = ""
	ns.Meta.LogoID = 0
	ns.Meta.LogoEnabled = false

	// - prune resources we won't preserve
	pref := envoyx.Ref{
		ResourceType: automationTypes.WorkflowResourceType,
	}
	for _, n := range nodes {
		n.Prune(pref)
	}

	return nodes
}

func (ctrl Namespace) preparePlaceholders(ctx context.Context) (resources envoyx.NodeSet, err error) {
	rr, _, err := ctrl.role.Find(ctx, systemTypes.RoleFilter{})
	if err != nil {
		return
	}
	var aux *envoyx.Node
	for _, role := range rr {
		aux, err = systemEnvoy.RoleToEnvoyNode(role)
		if err != nil {
			return
		}

		aux.Placeholder = true
		resources = append(resources, aux)
	}

	return
}
