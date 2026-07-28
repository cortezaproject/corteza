package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

func ProjectTask() *projectTask {
	return &projectTask{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *projectTask) beforeCreate(ctx context.Context, new *types.ProjectTask) error {
	if new.ProjectID == 0 {
		return ProjectTaskErrMissingProject()
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
		return ProjectTaskErrInvalidRevision()
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

func (svc *projectTask) beforeUpdate(ctx context.Context, upd, res *types.ProjectTask) error {
	if ok, err := revisionInChain(ctx, svc.store, res.ProjectID, upd.RevisionID); err != nil {
		return err
	} else if !ok {
		return ProjectTaskErrInvalidRevision()
	}
	// Stamp res (the persisted record); updated_by is not in the generated
	// field-copy, so stamping upd would be dropped.
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectTask) beforeDelete(ctx context.Context, res *types.ProjectTask) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
