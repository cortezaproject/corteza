package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	Tenant struct {
		svc tenantService
		ac  tenantAccessController
	}

	tenantPayload struct {
		*types.Tenant

		CanGrant         bool `json:"canGrant"`
		CanUpdateTenant  bool `json:"canUpdateTenant"`
		CanDeleteTenant  bool `json:"canDeleteTenant"`
		CanSuspendTenant bool `json:"canSuspendTenant"`
		CanManageMembers bool `json:"canManageMembers"`
	}

	tenantSetPayload struct {
		Filter types.TenantFilter `json:"filter"`
		Set    []*tenantPayload   `json:"set"`
	}

	tenantMemberPayload struct {
		*types.TenantMembership
	}

	tenantMemberSetPayload struct {
		Filter types.TenantMembershipFilter `json:"filter"`
		Set    []*tenantMemberPayload       `json:"set"`
	}

	tenantService interface {
		FindByID(ctx context.Context, ID uint64) (*types.Tenant, error)
		Create(ctx context.Context, new *types.Tenant) (*types.Tenant, error)
		Update(ctx context.Context, upd *types.Tenant) (*types.Tenant, error)
		DeleteByID(ctx context.Context, ID uint64) error
		UndeleteByID(ctx context.Context, ID uint64) error
		Search(ctx context.Context, filter types.TenantFilter) (types.TenantSet, types.TenantFilter, error)

		Suspend(ctx context.Context, ID uint64) error
		Activate(ctx context.Context, ID uint64) error
		Archive(ctx context.Context, ID uint64) error

		SearchMembers(ctx context.Context, filter types.TenantMembershipFilter) (types.TenantMembershipSet, types.TenantMembershipFilter, error)
		Invite(ctx context.Context, m *types.TenantMembership) (*types.TenantMembership, error)
		UpdateMember(ctx context.Context, m *types.TenantMembership) (*types.TenantMembership, error)
		RemoveMember(ctx context.Context, tenantID, userID uint64) error
		SuspendMember(ctx context.Context, tenantID, userID uint64) error
		ActivateMember(ctx context.Context, tenantID, userID uint64) error
	}

	tenantAccessController interface {
		CanGrant(context.Context) bool

		CanCreateTenant(context.Context) bool
		CanUpdateTenant(context.Context, *types.Tenant) bool
		CanDeleteTenant(context.Context, *types.Tenant) bool
		CanSuspendTenant(context.Context, *types.Tenant) bool
		CanManageMembersOnTenant(context.Context, *types.Tenant) bool
	}
)

func (Tenant) New() *Tenant {
	return &Tenant{
		svc: service.DefaultTenant,
		ac:  service.DefaultAccessControl,
	}
}

func (ctrl *Tenant) List(ctx context.Context, r *request.TenantList) (interface{}, error) {
	var (
		err error
		f   = types.TenantFilter{
			Query:   r.Query,
			Handle:  r.Handle,
			Status:  types.TenantStatus(r.Status),
			Labels:  r.Labels,
			Deleted: filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, f, err := ctrl.svc.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, f, err)
}

func (ctrl *Tenant) Create(ctx context.Context, r *request.TenantCreate) (interface{}, error) {
	t := &types.Tenant{
		Handle: r.Handle,
		Status: types.TenantStatus(r.Status),
		Config: r.Config,
		Meta:   r.Meta,
		Labels: r.Labels,
	}

	t, err := ctrl.svc.Create(ctx, t)
	return ctrl.makePayload(ctx, t, err)
}

func (ctrl *Tenant) Read(ctx context.Context, r *request.TenantRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.TenantID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *Tenant) Update(ctx context.Context, r *request.TenantUpdate) (interface{}, error) {
	t := &types.Tenant{
		ID:        r.TenantID,
		Handle:    r.Handle,
		Status:    types.TenantStatus(r.Status),
		Config:    r.Config,
		Meta:      r.Meta,
		UpdatedAt: r.UpdatedAt,
		Labels:    r.Labels,
	}

	t, err := ctrl.svc.Update(ctx, t)
	return ctrl.makePayload(ctx, t, err)
}

func (ctrl *Tenant) Delete(ctx context.Context, r *request.TenantDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.TenantID)
}

func (ctrl *Tenant) Undelete(ctx context.Context, r *request.TenantUndelete) (interface{}, error) {
	return api.OK(), ctrl.svc.UndeleteByID(ctx, r.TenantID)
}

func (ctrl *Tenant) Suspend(ctx context.Context, r *request.TenantSuspend) (interface{}, error) {
	return api.OK(), ctrl.svc.Suspend(ctx, r.TenantID)
}

func (ctrl *Tenant) Activate(ctx context.Context, r *request.TenantActivate) (interface{}, error) {
	return api.OK(), ctrl.svc.Activate(ctx, r.TenantID)
}

func (ctrl *Tenant) Archive(ctx context.Context, r *request.TenantArchive) (interface{}, error) {
	return api.OK(), ctrl.svc.Archive(ctx, r.TenantID)
}

func (ctrl *Tenant) ListMembers(ctx context.Context, r *request.TenantListMembers) (interface{}, error) {
	set, f, err := ctrl.svc.SearchMembers(ctx, types.TenantMembershipFilter{
		TenantID: r.TenantID,
	})
	return ctrl.makeMemberFilterPayload(set, f, err)
}

func (ctrl *Tenant) AddMember(ctx context.Context, r *request.TenantAddMember) (interface{}, error) {
	m, err := ctrl.svc.Invite(ctx, &types.TenantMembership{
		TenantID: r.TenantID,
		UserID:   r.UserID,
		Role:     types.TenantMemberRole(r.Role),
	})
	return ctrl.makeMemberPayload(m, err)
}

func (ctrl *Tenant) UpdateMember(ctx context.Context, r *request.TenantUpdateMember) (interface{}, error) {
	m, err := ctrl.svc.UpdateMember(ctx, &types.TenantMembership{
		TenantID: r.TenantID,
		UserID:   r.UserID,
		Role:     types.TenantMemberRole(r.Role),
		Status:   types.TenantMemberStatus(r.Status),
	})
	return ctrl.makeMemberPayload(m, err)
}

func (ctrl *Tenant) RemoveMember(ctx context.Context, r *request.TenantRemoveMember) (interface{}, error) {
	return api.OK(), ctrl.svc.RemoveMember(ctx, r.TenantID, r.UserID)
}

func (ctrl *Tenant) SuspendMember(ctx context.Context, r *request.TenantSuspendMember) (interface{}, error) {
	return api.OK(), ctrl.svc.SuspendMember(ctx, r.TenantID, r.UserID)
}

func (ctrl *Tenant) ActivateMember(ctx context.Context, r *request.TenantActivateMember) (interface{}, error) {
	return api.OK(), ctrl.svc.ActivateMember(ctx, r.TenantID, r.UserID)
}

func (ctrl *Tenant) makePayload(ctx context.Context, t *types.Tenant, err error) (*tenantPayload, error) {
	if err != nil || t == nil {
		return nil, err
	}

	return &tenantPayload{
		Tenant: t,

		CanGrant:         ctrl.ac.CanGrant(ctx),
		CanUpdateTenant:  ctrl.ac.CanUpdateTenant(ctx, t),
		CanDeleteTenant:  ctrl.ac.CanDeleteTenant(ctx, t),
		CanSuspendTenant: ctrl.ac.CanSuspendTenant(ctx, t),
		CanManageMembers: ctrl.ac.CanManageMembersOnTenant(ctx, t),
	}, nil
}

func (ctrl *Tenant) makeFilterPayload(ctx context.Context, nn types.TenantSet, f types.TenantFilter, err error) (*tenantSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &tenantSetPayload{Filter: f, Set: make([]*tenantPayload, len(nn))}
	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}

func (ctrl *Tenant) makeMemberPayload(m *types.TenantMembership, err error) (*tenantMemberPayload, error) {
	if err != nil || m == nil {
		return nil, err
	}

	return &tenantMemberPayload{TenantMembership: m}, nil
}

func (ctrl *Tenant) makeMemberFilterPayload(nn types.TenantMembershipSet, f types.TenantMembershipFilter, err error) (*tenantMemberSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &tenantMemberSetPayload{Filter: f, Set: make([]*tenantMemberPayload, len(nn))}
	for i := range nn {
		msp.Set[i], _ = ctrl.makeMemberPayload(nn[i], nil)
	}

	return msp, nil
}
