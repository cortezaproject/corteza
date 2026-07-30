package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD method bodies (FindByID, Search, Create, Update, DeleteByID and
// loadProjectFriaScenario) plus the access-controller interface are generated
// in project_fria_scenario.gen.go from system/project_fria_scenario.cue.
//
// This file owns the constructor, the public service contract and the
// before-create / before-update / before-delete hooks the generated CRUD calls
// into.

type (
	ProjectFriaScenarioService interface {
		FindByID(ctx context.Context, id uint64) (*types.ProjectFriaScenario, error)
		Search(ctx context.Context, f types.ProjectFriaScenarioFilter) (types.ProjectFriaScenarioSet, types.ProjectFriaScenarioFilter, error)
		Create(ctx context.Context, s *types.ProjectFriaScenario) (*types.ProjectFriaScenario, error)
		Update(ctx context.Context, s *types.ProjectFriaScenario) (*types.ProjectFriaScenario, error)
		DeleteByID(ctx context.Context, id uint64) error
	}
)

func ProjectFriaScenario() *projectFriaScenario {
	return &projectFriaScenario{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

// beforeCreate resolves the owning tenant from the project before the generated
// Create assigns the ID / timestamps and persists.
//
// Nothing else is required: a FRIA risk scenario is built up incrementally over
// several editing sessions, so an empty title, severity or meta blob is a legal
// intermediate state, not a validation failure.
func (svc *projectFriaScenario) beforeCreate(ctx context.Context, new *types.ProjectFriaScenario) error {
	if new.ProjectID == 0 {
		return ProjectFriaScenarioErrMissingProject()
	}

	p, err := loadProject(ctx, svc.store, new.ProjectID)
	if err != nil {
		return err
	}

	new.TenantID = p.TenantID
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

// beforeUpdate replaces the meta blob wholesale. Meta is omitSetter, so the
// generated Update does not copy it — and a field-by-field "only if non-empty"
// merge (the pattern project_ai_system uses for its prose-only meta) would make
// it impossible to CLEAR a taxonomy selection: deselecting every trigger type
// sends an empty list, which such a merge would read as "unchanged". The editor
// always submits the whole scenario, so a wholesale replace is both correct and
// consistent with how the plain columns are copied.
//
// The generated Update has already loaded `existing`, run the access check and
// the stale-data guard, and copies AiSystemID / Title / Severity / UpdatedBy
// from `upd` afterwards.
func (svc *projectFriaScenario) beforeUpdate(ctx context.Context, upd, existing *types.ProjectFriaScenario) error {
	existing.Meta = upd.Meta
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectFriaScenario) beforeDelete(ctx context.Context, res *types.ProjectFriaScenario) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
