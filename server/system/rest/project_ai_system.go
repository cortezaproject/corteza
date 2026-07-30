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
	ProjectAiSystem struct {
		projectAiSystem projectAiSystemService
		ac              projectAiSystemAccessController
	}

	projectAiSystemPayload struct {
		*types.ProjectAiSystem

		Entries            types.ProjectAiSystemEntrySet `json:"entries"`
		CanUpdate          bool                          `json:"canUpdate"`
		CanDelete          bool                          `json:"canDelete"`
		CanManageResources bool                          `json:"canManageResources"`
	}

	projectAiSystemSetPayload struct {
		Filter types.ProjectAiSystemFilter `json:"filter"`
		Set    []*projectAiSystemPayload   `json:"set"`
	}

	projectAiSystemService interface {
		FindByID(ctx context.Context, id uint64) (*types.ProjectAiSystem, error)
		Search(ctx context.Context, f types.ProjectAiSystemFilter) (types.ProjectAiSystemSet, types.ProjectAiSystemFilter, error)
		Create(ctx context.Context, g *types.ProjectAiSystem) (*types.ProjectAiSystem, error)
		Update(ctx context.Context, g *types.ProjectAiSystem) (*types.ProjectAiSystem, error)
		DeleteByID(ctx context.Context, id uint64) error
		MemberList(ctx context.Context, projectAiSystemID uint64) (types.ProjectAiSystemEntrySet, error)
		MemberAdd(ctx context.Context, projectAiSystemID uint64, resourceRef string) error
		MemberRemove(ctx context.Context, projectAiSystemID uint64, resourceRef string) error
	}

	projectAiSystemAccessController interface {
		CanGrant(context.Context) bool
		CanUpdateProjectAiSystem(context.Context, *types.ProjectAiSystem) bool
		CanDeleteProjectAiSystem(context.Context, *types.ProjectAiSystem) bool
		CanManageResourcesOnProjectAiSystem(context.Context, *types.ProjectAiSystem) bool
	}
)

func (ProjectAiSystem) New() *ProjectAiSystem {
	return &ProjectAiSystem{
		projectAiSystem: service.DefaultProjectAiSystem,
		ac:              service.DefaultAccessControl,
	}
}

func (ctrl *ProjectAiSystem) makeFilter(ctx context.Context, r *request.ProjectAiSystemList) (types.ProjectAiSystemFilter, error) {
	var (
		err error
		f   = types.ProjectAiSystemFilter{
			ProjectID: r.ProjectID,
			Handle:    r.Handle,
			RiskClass: r.RiskClass,
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
// directly onto the resource struct: name/description/intendedPurpose map onto
// res.Meta (genHook: true) rather than plain fields, and projectID is a scope
// id the endpoint drops under compound:false, so create restores it here. The
// generated Create has already assigned the plain Handle and RiskClass fields.
func (ctrl *ProjectAiSystem) beforeCreate(ctx context.Context, res *types.ProjectAiSystem, r *request.ProjectAiSystemCreate) error {
	res.ProjectID = r.ProjectID
	res.Meta.Short = r.Name
	res.Meta.Description = r.Description
	res.Meta.IntendedPurpose = r.IntendedPurpose
	return nil
}

// beforeUpdate fills the name/description/intendedPurpose request params onto
// res.Meta; they are genHook: true because they map onto res.Meta rather than
// plain struct fields. The generated Update has already assigned ID, the plain
// Handle and RiskClass fields and UpdatedAt.
func (ctrl *ProjectAiSystem) beforeUpdate(ctx context.Context, res *types.ProjectAiSystem, r *request.ProjectAiSystemUpdate) error {
	res.Meta.Short = r.Name
	res.Meta.Description = r.Description
	res.Meta.IntendedPurpose = r.IntendedPurpose
	return nil
}

func (ctrl *ProjectAiSystem) EntryAdd(ctx context.Context, r *request.ProjectAiSystemEntryAdd) (interface{}, error) {
	return api.OK(), ctrl.projectAiSystem.MemberAdd(ctx, r.ProjectAiSystemID, r.ResourceRef)
}

func (ctrl *ProjectAiSystem) EntryRemove(ctx context.Context, r *request.ProjectAiSystemEntryRemove) (interface{}, error) {
	return api.OK(), ctrl.projectAiSystem.MemberRemove(ctx, r.ProjectAiSystemID, r.ResourceRef)
}

func (ctrl *ProjectAiSystem) makePayload(ctx context.Context, g *types.ProjectAiSystem, err error) (*projectAiSystemPayload, error) {
	if err != nil || g == nil {
		return nil, err
	}

	entries, _ := ctrl.projectAiSystem.MemberList(ctx, g.ID)

	return &projectAiSystemPayload{
		ProjectAiSystem:    g,
		Entries:            entries,
		CanUpdate:          ctrl.ac.CanUpdateProjectAiSystem(ctx, g),
		CanDelete:          ctrl.ac.CanDeleteProjectAiSystem(ctx, g),
		CanManageResources: ctrl.ac.CanManageResourcesOnProjectAiSystem(ctx, g),
	}, nil
}

func (ctrl *ProjectAiSystem) makeFilterPayload(ctx context.Context, nn types.ProjectAiSystemSet, f types.ProjectAiSystemFilter, err error) (*projectAiSystemSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectAiSystemSetPayload{Filter: f, Set: make([]*projectAiSystemPayload, len(nn))}
	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
