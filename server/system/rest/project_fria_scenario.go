package rest

import (
	"context"

	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/system/rest/request"
	"github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

type (
	ProjectFriaScenario struct {
		projectFriaScenario projectFriaScenarioService
		ac                  projectFriaScenarioAccessController
	}

	projectFriaScenarioPayload struct {
		*types.ProjectFriaScenario

		CanUpdate bool `json:"canUpdate"`
		CanDelete bool `json:"canDelete"`
	}

	projectFriaScenarioSetPayload struct {
		Filter types.ProjectFriaScenarioFilter `json:"filter"`
		Set    []*projectFriaScenarioPayload   `json:"set"`
	}

	projectFriaScenarioService interface {
		FindByID(ctx context.Context, id uint64) (*types.ProjectFriaScenario, error)
		Search(ctx context.Context, f types.ProjectFriaScenarioFilter) (types.ProjectFriaScenarioSet, types.ProjectFriaScenarioFilter, error)
		Create(ctx context.Context, s *types.ProjectFriaScenario) (*types.ProjectFriaScenario, error)
		Update(ctx context.Context, s *types.ProjectFriaScenario) (*types.ProjectFriaScenario, error)
		DeleteByID(ctx context.Context, id uint64) error
	}

	projectFriaScenarioAccessController interface {
		CanGrant(context.Context) bool
		CanUpdateProjectFriaScenario(context.Context, *types.ProjectFriaScenario) bool
		CanDeleteProjectFriaScenario(context.Context, *types.ProjectFriaScenario) bool
	}
)

func (ProjectFriaScenario) New() *ProjectFriaScenario {
	return &ProjectFriaScenario{
		projectFriaScenario: service.DefaultProjectFriaScenario,
		ac:                  service.DefaultAccessControl,
	}
}

func (ctrl *ProjectFriaScenario) makeFilter(ctx context.Context, r *request.ProjectFriaScenarioList) (types.ProjectFriaScenarioFilter, error) {
	var (
		err error
		f   = types.ProjectFriaScenarioFilter{
			Query:      r.Query,
			ProjectID:  r.ProjectID,
			AiSystemID: r.AiSystemID,
			Severity:   r.Severity,
			Deleted:    filter.State(r.Deleted),
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
// directly onto the resource struct: the taxonomy key lists and the free prose
// map onto res.Meta (genHook: true) rather than plain fields, and projectID is
// a scope id the endpoint drops under compound:false, so create restores it
// here. The generated Create has already assigned the plain AiSystemID, Title
// and Severity fields.
func (ctrl *ProjectFriaScenario) beforeCreate(ctx context.Context, res *types.ProjectFriaScenario, r *request.ProjectFriaScenarioCreate) error {
	res.ProjectID = r.ProjectID
	res.Meta = types.ProjectFriaScenarioMeta{
		Description:            r.Description,
		TriggerTypes:           r.TriggerTypes,
		TriggerDescription:     r.TriggerDescription,
		ImpactedParties:        r.ImpactedParties,
		VulnerableGroups:       r.VulnerableGroups,
		VulnerableGroupsNotes:  r.VulnerableGroupsNotes,
		Rights:                 r.Rights,
		HarmVectors:            r.HarmVectors,
		HarmVectorsDescription: r.HarmVectorsDescription,
	}
	return nil
}

// beforeUpdate rebuilds res.Meta from the request params; they are genHook:true
// because they map onto res.Meta rather than plain struct fields. The whole
// blob is replaced (not merged field-by-field) so a taxonomy selection can be
// cleared — see the service's beforeUpdate. The generated Update has already
// assigned ID, the plain AiSystemID / Title / Severity fields and UpdatedAt.
func (ctrl *ProjectFriaScenario) beforeUpdate(ctx context.Context, res *types.ProjectFriaScenario, r *request.ProjectFriaScenarioUpdate) error {
	res.Meta = types.ProjectFriaScenarioMeta{
		Description:            r.Description,
		TriggerTypes:           r.TriggerTypes,
		TriggerDescription:     r.TriggerDescription,
		ImpactedParties:        r.ImpactedParties,
		VulnerableGroups:       r.VulnerableGroups,
		VulnerableGroupsNotes:  r.VulnerableGroupsNotes,
		Rights:                 r.Rights,
		HarmVectors:            r.HarmVectors,
		HarmVectorsDescription: r.HarmVectorsDescription,
	}
	return nil
}

func (ctrl *ProjectFriaScenario) makePayload(ctx context.Context, s *types.ProjectFriaScenario, err error) (*projectFriaScenarioPayload, error) {
	if err != nil || s == nil {
		return nil, err
	}

	return &projectFriaScenarioPayload{
		ProjectFriaScenario: s,
		CanUpdate:           ctrl.ac.CanUpdateProjectFriaScenario(ctx, s),
		CanDelete:           ctrl.ac.CanDeleteProjectFriaScenario(ctx, s),
	}, nil
}

func (ctrl *ProjectFriaScenario) makeFilterPayload(ctx context.Context, nn types.ProjectFriaScenarioSet, f types.ProjectFriaScenarioFilter, err error) (*projectFriaScenarioSetPayload, error) {
	if err != nil {
		return nil, err
	}

	msp := &projectFriaScenarioSetPayload{Filter: f, Set: make([]*projectFriaScenarioPayload, len(nn))}
	for i := range nn {
		msp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return msp, nil
}
