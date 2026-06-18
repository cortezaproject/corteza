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
		projectGroup projectGroupService
		ac           projectGroupAccessController
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
		projectGroup: service.DefaultProjectGroup,
		ac:           service.DefaultAccessControl,
	}
}

func (ctrl *ProjectGroup) makeFilter(ctx context.Context, r *request.ProjectGroupList) (types.ProjectGroupFilter, error) {
	var (
		err error
		f   = types.ProjectGroupFilter{
			ProjectID: r.ProjectID,
			Handle:    r.Handle,
			Deleted:   filter.State(r.Deleted),
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

// beforeCreate fills the request params the generated Create cannot assign
// directly onto the resource struct: name/description map onto res.Meta
// (genHook: true) rather than plain fields, and projectID is a scope id the
// endpoint drops under compound:false, so create restores it here. The
// generated Create has already assigned the plain Handle field.
func (ctrl *ProjectGroup) beforeCreate(ctx context.Context, res *types.ProjectGroup, r *request.ProjectGroupCreate) error {
	res.ProjectID = r.ProjectID
	res.Meta.Short = r.Name
	res.Meta.Description = r.Description
	return nil
}

// beforeUpdate fills the name/description request params onto res.Meta;
// they are genHook: true because they map onto res.Meta rather than plain
// struct fields. The generated Update has already assigned ID, the plain
// Handle field and UpdatedAt.
func (ctrl *ProjectGroup) beforeUpdate(ctx context.Context, res *types.ProjectGroup, r *request.ProjectGroupUpdate) error {
	res.Meta.Short = r.Name
	res.Meta.Description = r.Description
	return nil
}

func (ctrl *ProjectGroup) EntryAdd(ctx context.Context, r *request.ProjectGroupEntryAdd) (interface{}, error) {
	return api.OK(), ctrl.projectGroup.MemberAdd(ctx, r.ProjectGroupID, r.ResourceRef)
}

func (ctrl *ProjectGroup) EntryRemove(ctx context.Context, r *request.ProjectGroupEntryRemove) (interface{}, error) {
	return api.OK(), ctrl.projectGroup.MemberRemove(ctx, r.ProjectGroupID, r.ResourceRef)
}

func (ctrl *ProjectGroup) makePayload(ctx context.Context, g *types.ProjectGroup, err error) (*projectGroupPayload, error) {
	if err != nil || g == nil {
		return nil, err
	}

	entries, _ := ctrl.projectGroup.MemberList(ctx, g.ID)

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
