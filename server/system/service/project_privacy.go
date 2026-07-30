package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

func ProjectPrivacy() *projectPrivacy {
	return &projectPrivacy{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

// beforeSearch normalises the project scope to the chain root, mirroring what
// beforeCreate already does on the way in. See rootProjectID (project.go).
func (svc *projectPrivacy) beforeSearch(ctx context.Context, f *types.ProjectPrivacyFilter) error {
	f.ProjectID = rootProjectID(ctx, svc.store, f.ProjectID)
	return nil
}

func (svc *projectPrivacy) beforeCreate(ctx context.Context, new *types.ProjectPrivacy) error {
	if new.ProjectID == 0 {
		return ProjectPrivacyErrMissingProject()
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
		return ProjectPrivacyErrInvalidRevision()
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

func (svc *projectPrivacy) beforeUpdate(ctx context.Context, upd, res *types.ProjectPrivacy) error {
	if ok, err := revisionInChain(ctx, svc.store, res.ProjectID, upd.RevisionID); err != nil {
		return err
	} else if !ok {
		return ProjectPrivacyErrInvalidRevision()
	}
	// Stamp res (the persisted record); updated_by is not in the generated
	// field-copy, so stamping upd would be dropped.
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectPrivacy) beforeDelete(ctx context.Context, res *types.ProjectPrivacy) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
