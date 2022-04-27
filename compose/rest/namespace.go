package rest

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	automationTypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	envoyStore "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	systemService "github.com/cortezaproject/corteza-server/system/service"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	namespacePayload struct {
		*types.Namespace

		CanGrant           bool `json:"canGrant"`
		CanUpdateNamespace bool `json:"canUpdateNamespace"`
		CanDeleteNamespace bool `json:"canDeleteNamespace"`
		CanManageNamespace bool `json:"canManageNamespace"`
		CanCreateModule    bool `json:"canCreateModule"`
		CanCreateChart     bool `json:"canCreateChart"`
		CanCreatePage      bool `json:"canCreatePage"`
	}

	namespaceSetPayload struct {
		Filter types.NamespaceFilter `json:"filter"`
		Set    []*namespacePayload   `json:"set"`
	}

	pageFinder interface {
		Find(ctx context.Context, filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error)
	}

	chartFinder interface {
		Find(ctx context.Context, filter types.ChartFilter) (set types.ChartSet, f types.ChartFilter, err error)
	}

	Namespace struct {
		namespace  service.NamespaceService
		module     service.ModuleService
		page       pageFinder
		chart      chartFinder
		locale     service.ResourceTranslationsManagerService
		attachment service.AttachmentService
		role       systemService.RoleService
		ac         namespaceAccessController
	}

	namespaceAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateNamespace(context.Context, *types.Namespace) bool
		CanDeleteNamespace(context.Context, *types.Namespace) bool
		CanManageNamespace(context.Context, *types.Namespace) bool

		CanCreateModuleOnNamespace(context.Context, *types.Namespace) bool
		CanCreateChartOnNamespace(context.Context, *types.Namespace) bool
		CanCreatePageOnNamespace(context.Context, *types.Namespace) bool
	}
)

func (Namespace) New() *Namespace {
	return &Namespace{
		namespace:  service.DefaultNamespace,
		module:     service.DefaultModule,
		page:       service.DefaultPage,
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

	// prepare filters
	df := envoyStore.NewDecodeFilter()

	// - compose resources
	df = df.ComposeNamespace(&types.NamespaceFilter{
		NamespaceID: []uint64{r.NamespaceID},
	}).
		ComposeModule(&types.ModuleFilter{}).
		ComposePage(&types.PageFilter{}).
		ComposeChart(&types.ChartFilter{})

	// - workflow
	// @todo how do we want to handle these ones?
	//       do we handle these ones?

	decoder := func() (resource.InterfaceSet, error) {
		// get from store
		return envoyStore.Decoder().Decode(ctx, service.DefaultStore, df)
	}

	encoder := func(nn resource.InterfaceSet) error {
		// prepare for encoding
		se := envoyStore.NewStoreEncoder(service.DefaultStore, &envoyStore.EncoderConfig{})
		bld := envoy.NewBuilder(se)
		g, err := bld.Build(ctx, nn...)
		if err != nil {
			return err
		}

		return envoy.Encode(ctx, g, se)
	}

	ns, err := ctrl.namespace.Clone(ctx, r.NamespaceID, dup, decoder, encoder)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Namespace) Export(ctx context.Context, r *request.NamespaceExport) (out interface{}, err error) {
	// Get resources
	resources, err := ctrl.gatherResources(ctx, r.NamespaceID)
	if err != nil {
		return
	}

	// Encode
	ye := yaml.NewYamlEncoder(&yaml.EncoderConfig{})
	bld := envoy.NewBuilder(ye)
	g, err := bld.Build(ctx, resources...)
	if err != nil {
		return nil, err
	}

	err = envoy.Encode(ctx, g, ye)
	if err != nil {
		return
	}

	// Archive encoded resources
	buf := bytes.NewBuffer(nil)
	w := zip.NewWriter(buf)

	var (
		f  io.Writer
		bb []byte
	)
	for _, s := range ye.Stream() {
		// @todo generalize when needed
		f, err = w.Create(fmt.Sprintf("%s.yaml", s.Resource))
		if err != nil {
			return
		}

		bb, err = ioutil.ReadAll(s.Source)
		if err != nil {
			return
		}

		_, err = f.Write(bb)
		if err != nil {
			return
		}
	}

	err = w.Close()
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
}

func (ctrl Namespace) ImportRun(ctx context.Context, r *request.NamespaceImportRun) (interface{}, error) {
	var (
		dup = &types.Namespace{
			Name: r.Name,
			Slug: r.Slug,
		}

		encoder = func(nn resource.InterfaceSet) error {
			se := envoyStore.NewStoreEncoder(service.DefaultStore, &envoyStore.EncoderConfig{})

			bld := envoy.NewBuilder(se)
			g, err := bld.Build(ctx, nn...)
			if err != nil {
				return err
			}

			err = envoy.Encode(ctx, g, se)
			if err != nil {
				return err
			}

			return nil
		}
	)

	ns, err := ctrl.namespace.ImportRun(ctx, r.SessionID, dup, encoder)
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
		CanUpdateNamespace: ctrl.ac.CanUpdateNamespace(ctx, ns),
		CanDeleteNamespace: ctrl.ac.CanDeleteNamespace(ctx, ns),
		CanManageNamespace: ctrl.ac.CanManageNamespace(ctx, ns),

		CanCreateModule: ctrl.ac.CanCreateModuleOnNamespace(ctx, ns),
		CanCreateChart:  ctrl.ac.CanCreateChartOnNamespace(ctx, ns),
		CanCreatePage:   ctrl.ac.CanCreatePageOnNamespace(ctx, ns),
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

func (ctrl Namespace) gatherResources(ctx context.Context, namespaceID uint64) (resources resource.InterfaceSet, err error) {
	var (
		nsII resource.Identifiers
	)

	// Prepare resources
	resources, nsII, err = ctrl.exportCompose(ctx, namespaceID)
	if err != nil {
		return
	}

	// Tweak exported resources
	resources = ctrl.tweakExport(ctx, resources, nsII)

	// Role placeholders for RBAC
	var roleIndex map[uint64]*systemTypes.Role
	resources, roleIndex, err = ctrl.preparePlaceholders(ctx, resources)
	if err != nil {
		return
	}

	// RBAC
	auxRBAC, err := ctrl.exportRBAC(ctx, roleIndex, resources)
	if err != nil {
		return
	}

	// Translations
	auxResTrans, err := ctrl.exportResourceTranslations(ctx, resources)
	if err != nil {
		return
	}
	resources = append(resources, auxRBAC...)
	resources = append(resources, auxResTrans...)

	return
}

func (ctrl Namespace) exportCompose(ctx context.Context, namespaceID uint64) (resources resource.InterfaceSet, nsII resource.Identifiers, err error) {
	// - namespace
	n, err := ctrl.namespace.FindByID(ctx, namespaceID)
	if err != nil {
		return
	}
	nsRes := resource.NewComposeNamespace(n)
	nsII = nsRes.Identifiers()
	resources = append(resources, nsRes)
	// - modules
	mm, _, err := ctrl.module.Find(ctx, types.ModuleFilter{NamespaceID: n.ID})
	if err != nil {
		return
	}
	for _, m := range mm {
		km := resource.NewComposeModule(m, resource.MakeNamespaceRef(n.ID, n.Slug, n.Name))
		for _, f := range m.Fields {
			km.AddField(resource.NewComposeModuleField(f, km.RefNs, km.Ref()))
		}
		resources = append(resources, km)
	}
	// - pages
	pp, _, err := ctrl.page.Find(ctx, types.PageFilter{NamespaceID: n.ID})
	if err != nil {
		return
	}
	for _, p := range pp {
		p, modRef, parentRef := resource.UnpackComposePage(p)
		resources = append(resources, resource.NewComposePage(
			p,
			resource.MakeNamespaceRef(n.ID, n.Slug, n.Name),
			modRef,
			parentRef,
		))
	}
	// - charts
	cc, _, err := ctrl.chart.Find(ctx, types.ChartFilter{NamespaceID: n.ID})
	if err != nil {
		return
	}
	for _, c := range cc {
		refMods := make(resource.RefSet, 0, 2)
		for _, r := range c.Config.Reports {
			refMods = append(refMods, resource.MakeModuleRef(r.ModuleID, "", ""))
		}
		resources = append(resources, resource.NewComposeChart(
			c,
			resource.MakeNamespaceRef(n.ID, n.Slug, n.Name),
			refMods,
		))
	}

	return
}

func (ctrl Namespace) exportRBAC(ctx context.Context, roleIndex map[uint64]*systemTypes.Role, base resource.InterfaceSet) (resources resource.InterfaceSet, err error) {
	// Prepare RBAC Rules
	rawRules := rbac.Global().Rules()
	rules := make([]*resource.RbacRule, 0, len(rawRules))
	for _, rule := range rawRules {
		_, ref, pp, err := resource.ParseRule(rule.Resource)
		if err != nil {
			return nil, err
		}

		role, ok := roleIndex[rule.RoleID]
		if !ok {
			continue
		}

		rules = append(rules, resource.NewRbacRule(
			rule,
			resource.MakeRoleRef(role.ID, role.Handle, role.Name),
			ref,
			rule.Resource,
			pp...,
		))
	}

	for _, r := range envoy.FilterRequestedRBACRules(base, rules) {
		resources = append(resources, r)
	}

	return
}

func (ctrl Namespace) exportResourceTranslations(ctx context.Context, base resource.InterfaceSet) (resources resource.InterfaceSet, err error) {
	var (
		lsvc         = locale.Global()
		tags         = lsvc.Tags()
		translations = make([]*resource.ResourceTranslation, 0, 124)

		resKeyTrans map[string]map[string]*locale.ResourceTranslation
	)

	for _, t := range tags {
		resKeyTrans, err = lsvc.LoadResourceTranslations(ctx, t)
		if err != nil {
			return
		}

		for transRes, keyTrans := range resKeyTrans {
			rawTranslations := make(locale.ResourceTranslationSet, 0, len(resKeyTrans))
			for _, trans := range keyTrans {
				rawTranslations = append(rawTranslations, trans)
			}

			_, ref, pp, err := resource.ParseResourceTranslation(transRes)
			if err != nil {
				return nil, err
			}

			translations = append(translations, resource.NewResourceTranslation(
				systemTypes.FromLocale(rawTranslations),
				ref.Identifiers.First(),
				ref,
				pp...,
			))
		}
	}

	for _, t := range envoy.FilterRequiredResourceTranslations(base, translations) {
		resources = append(resources, t)
	}

	return
}

func (ctrl Namespace) tweakExport(ctx context.Context, resources resource.InterfaceSet, nsII resource.Identifiers) resource.InterfaceSet {
	oldNsRef := resource.MakeRef(types.NamespaceResourceType, nsII)
	prune := resource.RefSet{resource.MakeWildRef(automationTypes.WorkflowResourceType)}
	ns := resource.FindComposeNamespace(resources, nsII)

	// - remove logo and icon references as attachments are not exported by default
	// @todo code in attachment exporting, most likely when we do attachment handling rework
	ns.Meta.Icon = ""
	ns.Meta.IconID = 0
	ns.Meta.Logo = ""
	ns.Meta.LogoID = 0
	ns.Meta.LogoEnabled = false

	// - prune resources we won't preserve
	resources.SearchForReferences(oldNsRef).Walk(func(r resource.Interface) error {
		pp, ok := r.(resource.PrunableInterface)
		if !ok {
			return nil
		}

		for _, p := range prune {
			pp.Prune(p)
		}
		return nil
	})

	return resources
}

func (ctrl Namespace) preparePlaceholders(ctx context.Context, base resource.InterfaceSet) (resources resource.InterfaceSet, roleIndex map[uint64]*systemTypes.Role, err error) {
	resources = base

	// Get roles as we'll need them later for some resources
	roleIndex = make(map[uint64]*systemTypes.Role)
	rr, _, err := ctrl.role.Find(ctx, systemTypes.RoleFilter{})
	if err != nil {
		return
	}
	for _, role := range rr {
		roleIndex[role.ID] = role

		// Add them as placeholders since we don't want to export them
		r := resource.NewRole(role)
		r.MarkPlaceholder()
		resources = append(resources, r)
	}

	return
}
