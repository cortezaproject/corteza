package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

func ProjectIncident() *projectIncident {
	return &projectIncident{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

// beforeSearch normalises the project scope to the chain root, mirroring what
// beforeCreate already does on the way in. See rootProjectID (project.go).
func (svc *projectIncident) beforeSearch(ctx context.Context, f *types.ProjectIncidentFilter) error {
	f.ProjectID = rootProjectID(ctx, svc.store, f.ProjectID)
	return nil
}

func (svc *projectIncident) beforeCreate(ctx context.Context, new *types.ProjectIncident) error {
	if new.ProjectID == 0 {
		return ProjectIncidentErrMissingProject()
	}

	// Work items always file against the chain root, never a revision, so a
	// revision reads like a milestone over one shared item pool. New items are
	// deliberately left unassigned (RevisionID stays 0) here.
	project, err := loadProject(ctx, svc.store, new.ProjectID)
	if err != nil {
		return err
	}
	new.ProjectID = project.RootProjectID()

	if ok, err := revisionInChain(ctx, svc.store, new.ProjectID, new.RevisionID); err != nil {
		return err
	} else if !ok {
		return ProjectIncidentErrInvalidRevision()
	}

	if new.Title == "" {
		new.Title = firstNonEmpty(new.Description, "(untitled)")
	}
	if new.Status == "" {
		new.Status = "Open"
	}
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectIncident) beforeUpdate(ctx context.Context, upd, res *types.ProjectIncident) error {
	if ok, err := revisionInChain(ctx, svc.store, res.ProjectID, upd.RevisionID); err != nil {
		return err
	} else if !ok {
		return ProjectIncidentErrInvalidRevision()
	}
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectIncident) beforeDelete(ctx context.Context, res *types.ProjectIncident) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
