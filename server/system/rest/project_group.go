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
	ProjectGroup struct {
		svc projectGroupService
		ac  projectGroupAccessController
	}

	projectGroupPayload struct {
		*types.ProjectGroup

		Entries          types.ProjectGroupEntrySet `json:"entries"`
		CanUpdate        bool                       `json:"canUpdate"`
		CanDelete        bool                       `json:"canDelete"`
		CanManageMembers bool                       `json:"canManageMembers"`
	}

	projectGroupSetPayload struct {
		Filter types.ProjectGroupFilter `json:"filter"`
		Set    []*projectGroupPayload   `json:"set"`
	}

	projectGroupService interface {
		FindByID(ctx context.Context, id uint64) (*types.ProjectGroup, error)
		Search(ctx context.Context, f types.ProjectGroupFilter) (types.ProjectGroupSet, types.ProjectGroupFilter, error)
		Create(ctx context.Context, g *types.ProjectGroup) (*types.ProjectGroup, error)
		Update(ctx context.Context, g *types.ProjectGroup) (*types.ProjectGroup, error)
		DeleteByID(ctx context.Context, id uint64) error
		MemberList(ctx context.Context, projectGroupID uint64) (types.ProjectGroupEntrySet, error)
		MemberAdd(ctx context.Context, projectGroupID uint64, resourceRef string) error
		MemberRemove(ctx context.Context, projectGroupID uint64, resourceRef string) error
	}

	projectGroupAccessController interface {
		CanGrant(context.Context) bool
		CanUpdateProjectGroup(context.Context, *types.ProjectGroup) bool
		CanDeleteProjectGroup(context.Context, *types.ProjectGroup) bool
		CanManageMembersOnProjectGroup(context.Context, *types.ProjectGroup) bool
	}
)

func (ProjectGroup) New() *ProjectGroup {
	return &ProjectGroup{
		svc: service.DefaultProjectGroup,
		ac:  service.DefaultAccessControl,
	}
}

func (ctrl *ProjectGroup) List(ctx context.Context, r *request.ProjectGroupList) (interface{}, error) {
	var (
		err error
		f   = types.ProjectGroupFilter{
			ProjectID: r.ProjectID,
			Handle:    r.Handle,
			Deleted:   filter.State(r.Deleted),
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

func (ctrl *ProjectGroup) Create(ctx context.Context, r *request.ProjectGroupCreate) (interface{}, error) {
	g := &types.ProjectGroup{
		ProjectID: r.ProjectID,
		Handle:    r.Handle,
		Meta: types.ProjectGroupMeta{
			Short:       r.Name,
			Description: r.Description,
		},
	}

	g, err := ctrl.svc.Create(ctx, g)
	return ctrl.makePayload(ctx, g, err)
}

func (ctrl *ProjectGroup) Read(ctx context.Context, r *request.ProjectGroupRead) (interface{}, error) {
	res, err := ctrl.svc.FindByID(ctx, r.ProjectGroupID)
	return ctrl.makePayload(ctx, res, err)
}

func (ctrl *ProjectGroup) Update(ctx context.Context, r *request.ProjectGroupUpdate) (interface{}, error) {
	g := &types.ProjectGroup{
		ID:     r.ProjectGroupID,
		Handle: r.Handle,
		Meta: types.ProjectGroupMeta{
			Short:       r.Name,
			Description: r.Description,
		},
		UpdatedAt: r.UpdatedAt,
	}

	g, err := ctrl.svc.Update(ctx, g)
	return ctrl.makePayload(ctx, g, err)
}

func (ctrl *ProjectGroup) Delete(ctx context.Context, r *request.ProjectGroupDelete) (interface{}, error) {
	return api.OK(), ctrl.svc.DeleteByID(ctx, r.ProjectGroupID)
}

func (ctrl *ProjectGroup) EntryAdd(ctx context.Context, r *request.ProjectGroupEntryAdd) (interface{}, error) {
	return api.OK(), ctrl.svc.MemberAdd(ctx, r.ProjectGroupID, r.ResourceRef)
}

func (ctrl *ProjectGroup) EntryRemove(ctx context.Context, r *request.ProjectGroupEntryRemove) (interface{}, error) {
	return api.OK(), ctrl.svc.MemberRemove(ctx, r.ProjectGroupID, r.ResourceRef)
}

func (ctrl *ProjectGroup) makePayload(ctx context.Context, g *types.ProjectGroup, err error) (*projectGroupPayload, error) {
	if err != nil || g == nil {
		return nil, err
	}

	entries, _ := ctrl.svc.MemberList(ctx, g.ID)

	return &projectGroupPayload{
		ProjectGroup:     g,
		Entries:          entries,
		CanUpdate:        ctrl.ac.CanUpdateProjectGroup(ctx, g),
		CanDelete:        ctrl.ac.CanDeleteProjectGroup(ctx, g),
		CanManageMembers: ctrl.ac.CanManageMembersOnProjectGroup(ctx, g),
	}, nil
}

func (ctrl *ProjectGroup) makeFilterPayload(ctx context.Context, nn types.ProjectGroupSet, f types.ProjectGroupFilter, err error) (*projectGroupSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectGroupSetPayload{Filter: f, Set: make([]*projectGroupPayload, len(nn))}
	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
