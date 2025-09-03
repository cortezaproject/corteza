package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/cortezaproject/corteza/server/system/rest/request"
	"github.com/cortezaproject/corteza/server/system/service"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	UserGroup struct {
		userGroup service.UserGroupService
		ac        userGroupAccessController
	}

	userGroupAccessController interface {
		CanGrant(context.Context) bool

		CanCreateUserGroup(context.Context) bool
		CanUpdateUserGroup(context.Context, *types.UserGroup) bool
		CanDeleteUserGroup(context.Context, *types.UserGroup) bool
		CanManageMembersOnUserGroup(context.Context, *types.UserGroup) bool
	}

	userGroupPayload struct {
		*types.UserGroup

		CanGrant                    bool `json:"canGrant"`
		CanUpdateUserGroup          bool `json:"canUpdateUserGroup"`
		CanDeleteUserGroup          bool `json:"canDeleteUserGroup"`
		CanManageMembersOnUserGroup bool `json:"canManageMembersOnUserGroup"`
	}

	userGroupSetPayload struct {
		Filter types.UserGroupFilter `json:"filter"`
		Set    []*userGroupPayload   `json:"set"`
	}
)

func (UserGroup) New() *UserGroup {
	return &UserGroup{
		userGroup: service.DefaultUserGroup,
		ac:        service.DefaultAccessControl,
	}
}

func (ctrl UserGroup) Read(ctx context.Context, r *request.UserGroupRead) (interface{}, error) {
	userGroup, err := ctrl.userGroup.FindByID(ctx, r.UserGroupID.Number())
	return ctrl.makePayload(ctx, userGroup, err)
}

func (ctrl UserGroup) List(ctx context.Context, r *request.UserGroupList) (interface{}, error) {
	var (
		err error
		f   = types.UserGroupFilter{
			Query:       r.Query,
			Labels:      r.Labels,
			MemberID:    r.MemberID,
			UserGroupID: id.StringifySlice(r.UserGroupID...),

			Archived: filter.State(r.Archived),
			Deleted:  filter.State(r.Deleted),
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.userGroup.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl UserGroup) Create(ctx context.Context, r *request.UserGroupCreate) (interface{}, error) {
	var (
		err       error
		userGroup = &types.UserGroup{
			Handle: r.Handle,
			Labels: r.Labels,
			SelfID: r.SelfID,
			Meta:   r.Meta,
		}
	)

	userGroup, err = ctrl.userGroup.Create(ctx, userGroup)
	if err != nil {
		return nil, err
	}

	for _, userID := range payload.ParseUint64s(r.Members) {
		err := ctrl.userGroup.MemberAdd(ctx, userGroup.ID, userID)
		if err != nil {
			return nil, err
		}
	}
	return ctrl.makePayload(ctx, userGroup, err)
}

func (ctrl UserGroup) Update(ctx context.Context, r *request.UserGroupUpdate) (interface{}, error) {
	var (
		err       error
		userGroup = &types.UserGroup{
			ID:        r.UserGroupID.Number(),
			SelfID:    r.SelfID,
			Handle:    r.Handle,
			Labels:    r.Labels,
			UpdatedAt: r.UpdatedAt,
			Meta:      r.Meta,
		}
	)

	userGroup, err = ctrl.userGroup.Update(ctx, userGroup)
	if err != nil {
		return nil, err
	}

	if len(r.Members) > 0 {
		for _, userID := range payload.ParseUint64s(r.Members) {
			err := ctrl.userGroup.MemberAdd(ctx, userGroup.ID, userID)
			if err != nil {
				return nil, err
			}
		}
	}

	return ctrl.makePayload(ctx, userGroup, err)
}

func (ctrl UserGroup) Delete(ctx context.Context, r *request.UserGroupDelete) (interface{}, error) {
	return api.OK(), ctrl.userGroup.Delete(ctx, r.UserGroupID.Number())
}

func (ctrl UserGroup) Undelete(ctx context.Context, r *request.UserGroupUndelete) (interface{}, error) {
	return api.OK(), ctrl.userGroup.Undelete(ctx, r.UserGroupID.Number())
}

func (ctrl UserGroup) MemberList(ctx context.Context, r *request.UserGroupMemberList) (interface{}, error) {
	if mm, err := ctrl.userGroup.MemberList(ctx, r.UserGroupID.Number()); err != nil {
		return nil, err
	} else {
		rval := make([]string, len(mm))
		for i := range mm {
			rval[i] = payload.Uint64toa(mm[i].ID)
		}
		return rval, nil
	}
}

func (ctrl UserGroup) MemberAdd(ctx context.Context, r *request.UserGroupMemberAdd) (interface{}, error) {
	return api.OK(), ctrl.userGroup.MemberAdd(ctx, r.UserGroupID.Number(), r.UserID)
}

func (ctrl UserGroup) makePayload(ctx context.Context, r *types.UserGroup, err error) (*userGroupPayload, error) {
	if err != nil || r == nil {
		return nil, err
	}

	return &userGroupPayload{
		UserGroup: r,

		CanGrant:                    ctrl.ac.CanGrant(ctx),
		CanUpdateUserGroup:          ctrl.ac.CanUpdateUserGroup(ctx, r),
		CanDeleteUserGroup:          ctrl.ac.CanDeleteUserGroup(ctx, r),
		CanManageMembersOnUserGroup: ctrl.ac.CanManageMembersOnUserGroup(ctx, r),
	}, nil
}

func (ctrl UserGroup) makeFilterPayload(ctx context.Context, nn types.UserGroupSet, f types.UserGroupFilter, err error) (*userGroupSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &userGroupSetPayload{Filter: f, Set: make([]*userGroupPayload, len(nn))}

	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
