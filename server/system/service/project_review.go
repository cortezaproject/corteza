package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

func ProjectReview() *projectReview {
	return &projectReview{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *projectReview) beforeCreate(ctx context.Context, new *types.ProjectReview) error {
	if new.ProjectID == 0 {
		return ProjectReviewErrMissingProject()
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
		return ProjectReviewErrInvalidRevision()
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

func (svc *projectReview) beforeUpdate(ctx context.Context, upd, res *types.ProjectReview) error {
	if ok, err := revisionInChain(ctx, svc.store, res.ProjectID, upd.RevisionID); err != nil {
		return err
	} else if !ok {
		return ProjectReviewErrInvalidRevision()
	}
	// Stamp res (the persisted record); updated_by is not in the generated
	// field-copy, so stamping upd would be dropped.
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectReview) beforeDelete(ctx context.Context, res *types.ProjectReview) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
