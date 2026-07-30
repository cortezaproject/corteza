package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/system/types"
)

// The struct, access-controller interface, and CRUD method bodies (FindByID,
// Search, Create, Update, DeleteByID, loadProjectFeature) are generated in
// project_feature.gen.go from the CUE model.
//
// This file owns the constructor and the before-create / before-update hooks
// the generated Create / Update call into.

func ProjectFeature() *projectFeature {
	return &projectFeature{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

// beforeSearch normalises the project scope to the chain root, mirroring what
// beforeCreate already does on the way in. See rootProjectID (project.go).
func (svc *projectFeature) beforeSearch(ctx context.Context, f *types.ProjectFeatureFilter) error {
	f.ProjectID = rootProjectID(ctx, svc.store, f.ProjectID)
	return nil
}

func (svc *projectFeature) beforeCreate(ctx context.Context, new *types.ProjectFeature) error {
	if new.ProjectID == 0 {
		return ProjectFeatureErrMissingProject()
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
		return ProjectFeatureErrInvalidRevision()
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

func (svc *projectFeature) beforeUpdate(ctx context.Context, upd, res *types.ProjectFeature) error {
	if ok, err := revisionInChain(ctx, svc.store, res.ProjectID, upd.RevisionID); err != nil {
		return err
	} else if !ok {
		return ProjectFeatureErrInvalidRevision()
	}
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

func (svc *projectFeature) beforeDelete(ctx context.Context, res *types.ProjectFeature) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}
