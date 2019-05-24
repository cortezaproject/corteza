package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/cortezaproject/corteza-server/compose/internal/service"
	"github.com/cortezaproject/corteza-server/compose/rest/request"
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	namespacePayload struct {
		*types.Namespace

		CanGrant           bool `json:"canGrant"`
		CanUpdateNamespace bool `json:"canUpdateNamespace"`
		CanDeleteNamespace bool `json:"canDeleteNamespace"`
		CanCreateModule    bool `json:"canCreateModule"`
		CanCreateChart     bool `json:"canCreateChart"`
		CanCreateTrigger   bool `json:"canCreateTrigger"`
		CanCreatePage      bool `json:"canCreatePage"`
	}

	namespaceSetPayload struct {
		Filter types.NamespaceFilter `json:"filter"`
		Set    []*namespacePayload   `json:"set"`
	}

	Namespace struct {
		namespace service.NamespaceService
		ac        namespaceAccessController
	}

	namespaceAccessController interface {
		CanGrant(context.Context) bool

		CanUpdateNamespace(context.Context, *types.Namespace) bool
		CanDeleteNamespace(context.Context, *types.Namespace) bool

		CanCreateModule(context.Context, *types.Namespace) bool
		CanCreateChart(context.Context, *types.Namespace) bool
		CanCreateTrigger(context.Context, *types.Namespace) bool
		CanCreatePage(context.Context, *types.Namespace) bool
	}
)

func (Namespace) New() *Namespace {
	return &Namespace{
		namespace: service.DefaultNamespace,
		ac:        service.DefaultAccessControl,
	}
}

func (ctrl Namespace) List(ctx context.Context, r *request.NamespaceList) (interface{}, error) {
	f := types.NamespaceFilter{
		Query:   r.Query,
		PerPage: r.PerPage,
		Page:    r.Page,
	}

	set, filter, err := ctrl.namespace.With(ctx).Find(f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl Namespace) Create(ctx context.Context, r *request.NamespaceCreate) (interface{}, error) {
	var err error
	ns := &types.Namespace{
		Name:    r.Name,
		Slug:    r.Slug,
		Meta:    r.Meta,
		Enabled: r.Enabled,
	}

	ns, err = ctrl.namespace.With(ctx).Create(ns)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Namespace) Read(ctx context.Context, r *request.NamespaceRead) (interface{}, error) {
	ns, err := ctrl.namespace.With(ctx).FindByID(r.NamespaceID)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Namespace) Update(ctx context.Context, r *request.NamespaceUpdate) (interface{}, error) {
	var (
		ns  = &types.Namespace{}
		err error
	)

	ns.ID = r.NamespaceID
	ns.Name = r.Name
	ns.Slug = r.Slug
	ns.Meta = r.Meta
	ns.Enabled = r.Enabled
	ns.UpdatedAt = r.UpdatedAt

	ns, err = ctrl.namespace.With(ctx).Update(ns)
	return ctrl.makePayload(ctx, ns, err)
}

func (ctrl Namespace) Delete(ctx context.Context, r *request.NamespaceDelete) (interface{}, error) {
	_, err := ctrl.namespace.With(ctx).FindByID(r.NamespaceID)
	if err != nil {
		return nil, err
	}

	return resputil.OK(), ctrl.namespace.With(ctx).DeleteByID(r.NamespaceID)
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

		CanCreateModule:  ctrl.ac.CanCreateModule(ctx, ns),
		CanCreateChart:   ctrl.ac.CanCreateChart(ctx, ns),
		CanCreateTrigger: ctrl.ac.CanCreateTrigger(ctx, ns),
		CanCreatePage:    ctrl.ac.CanCreatePage(ctx, ns),
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
