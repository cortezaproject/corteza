package rest

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/corredor"
	"github.com/cortezaproject/corteza-server/pkg/filter"
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

	Namespace struct {
		namespace  service.NamespaceService
		locale     service.ResourceTranslationsManagerService
		attachment service.AttachmentService
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
		locale:     service.DefaultResourceTranslation,
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

	a, err := ctrl.attachment.With(ctx).CreateNamespaceAttachment(
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

func (ctrl *Namespace) TriggerScript(ctx context.Context, r *request.NamespaceTriggerScript) (rsp interface{}, err error) {
	var (
		namespace *types.Namespace
	)

	if namespace, err = ctrl.namespace.FindByID(ctx, r.NamespaceID); err != nil {
		return
	}

	err = corredor.Service().Exec(ctx, r.Script, event.NamespaceOnManual(namespace, nil))
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
