package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/id"
	"github.com/crusttech/human/server/pkg/payload"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
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

// makeFilter builds the userGroup search filter from the list request.
//
// Companion hook for the generated List controller.
func (ctrl UserGroup) makeFilter(ctx context.Context, r *request.UserGroupList) (types.UserGroupFilter, error) {
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
		return f, err
	}

	f.IncTotal = r.IncTotal

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return f, err
	}

	return f, nil
}

// beforeCreate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Create controller
// already mapped the plain-value params (handle, labels).
func (ctrl UserGroup) beforeCreate(ctx context.Context, res *types.UserGroup, r *request.UserGroupCreate) error {
	res.Config = r.Config
	res.Meta = r.Meta

	return nil
}

// afterCreate adds the requested members to the freshly created user group.
// types.UserGroup has no Members field -- the ids are added one-by-one via
// MemberAdd, which needs the assigned resource id, so it runs after Create.
func (ctrl UserGroup) afterCreate(ctx context.Context, res *types.UserGroup, r *request.UserGroupCreate) error {
	for _, userID := range payload.ParseUint64s(r.Members) {
		if err := ctrl.userGroup.MemberAdd(ctx, res.ID, userID); err != nil {
			return err
		}
	}

	return nil
}

// beforeUpdate fills the complex/hook-managed fields onto the resource
// before it is handed to the service. The generated Update controller
// already mapped the plain-value params (handle, labels) and ID/UpdatedAt.
func (ctrl UserGroup) beforeUpdate(ctx context.Context, res *types.UserGroup, r *request.UserGroupUpdate) error {
	res.Config = r.Config
	res.Meta = r.Meta

	return nil
}

// afterUpdate adds the requested members to the updated user group. As with
// create, MemberAdd needs the resource id and so runs after Update.
func (ctrl UserGroup) afterUpdate(ctx context.Context, res *types.UserGroup, r *request.UserGroupUpdate) error {
	if len(r.Members) > 0 {
		for _, userID := range payload.ParseUint64s(r.Members) {
			if err := ctrl.userGroup.MemberAdd(ctx, res.ID, userID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ctrl UserGroup) MemberList(ctx context.Context, r *request.UserGroupMemberList) (interface{}, error) {
	if mm, err := ctrl.userGroup.MemberList(ctx, r.UserGroupID.Num()); err != nil {
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
	return api.OK(), ctrl.userGroup.MemberAdd(ctx, r.UserGroupID.Num(), r.UserID)
}

func (ctrl UserGroup) makePayload(ctx context.Context, r *types.UserGroup, err error) (*userGroupPayload, error) {
	if err != nil || r == nil {
		return nil, err
	}

	if r.Config == nil {
		r.Config = &types.UserGroupConfig{}
	}

	if len(r.Config.Paths) == 0 {
		r.Config.Paths = []types.UserGroupPath{}
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
