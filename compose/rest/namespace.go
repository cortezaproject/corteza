package rest

import (
	"context"

	"github.com/titpetric/factory/resputil"

	"github.com/crusttech/crust/compose/internal/service"
	"github.com/crusttech/crust/compose/rest/request"
	"github.com/crusttech/crust/compose/types"
)

type (
	namespacePayload struct {
		*types.Namespace

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
)

type Namespace struct {
	namespace   service.NamespaceService
	permissions service.PermissionsService
}

func (Namespace) New() *Namespace {
	return &Namespace{
		namespace:   service.DefaultNamespace,
		permissions: service.DefaultPermissions,
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

	perm := ctrl.permissions.With(ctx)

	return &namespacePayload{
		Namespace: ns,

		CanUpdateNamespace: perm.CanUpdateNamespace(ns),
		CanDeleteNamespace: perm.CanDeleteNamespace(ns),
		CanCreateModule:    perm.CanCreateModule(ns),
		CanCreateChart:     perm.CanCreateChart(ns),
		CanCreateTrigger:   perm.CanCreateTrigger(ns),
		CanCreatePage:      perm.CanCreatePage(ns),
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
