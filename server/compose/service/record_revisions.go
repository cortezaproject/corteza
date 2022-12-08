package service

import (
	"context"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/revisions"
)

type (
	// wrapper around revision.service
	//
	// handles all module-to-model conversion
	recordRevisions struct {
		r interface {
			Search(ctx context.Context, mf dal.ModelRef, f filter.Filter) (_ dal.Iterator, err error)
			Create(ctx context.Context, mf dal.ModelRef, revision *revisions.Revision) error
		}
	}
)

// constructs DAL model-ref from module
//
// Since we're using revisions model resource type needs to be altered.
func (svc *recordRevisions) modelRef(mod *types.Module) (mf dal.ModelRef) {
	mf = mod.ModelRef()

	// removing resource ID from model filter
	// this will force DAL service to use FindModelByResourceIdent
	// to find model with same Resource value and RevisionResourceType
	// for ResourceType
	mf.ResourceID = 0
	mf.ResourceType = revisions.RevisionResourceType
	return
}

func (svc *recordRevisions) search(ctx context.Context, rec *types.Record) (_ dal.Iterator, err error) {
	var (
		revModRef = svc.modelRef(rec.GetModule())
		revFilter = filter.Generic(filter.WithConstraint("rel_resource", rec.ID))
	)

	return svc.r.Search(ctx, revModRef, revFilter)
}

func (svc *recordRevisions) created(ctx context.Context, new *types.Record) (err error) {
	var (
		skipList  = svc.skippedField(new.GetModule())
		invokerID = auth.GetIdentityFromContext(ctx).Identity()
		rev       *revisions.Revision
	)

	rev = revisions.Make(revisions.Created, new.Revision, new.ID, invokerID)
	err = rev.CollectChanges(new, nil, skipList...)
	if err != nil {
		return
	}

	return svc.r.Create(ctx, svc.modelRef(new.GetModule()), rev)
}

func (svc *recordRevisions) updated(ctx context.Context, upd, old *types.Record) (err error) {
	var (
		skipList  = svc.skippedField(old.GetModule())
		invokerID = auth.GetIdentityFromContext(ctx).Identity()
		rev       *revisions.Revision
	)

	rev = revisions.Make(revisions.Updated, upd.Revision, upd.ID, invokerID)
	err = rev.CollectChanges(upd, old, skipList...)
	if err != nil {
		return
	}

	return svc.r.Create(ctx, svc.modelRef(upd.GetModule()), rev)
}

func (svc *recordRevisions) softDeleted(ctx context.Context, del *types.Record) (err error) {
	var (
		invokerID = auth.GetIdentityFromContext(ctx).Identity()
		rev       *revisions.Revision
	)

	rev = revisions.Make(revisions.SoftDeleted, del.Revision, del.ID, invokerID)
	if err != nil {
		return
	}

	return svc.r.Create(ctx, svc.modelRef(del.GetModule()), rev)
}

func (svc *recordRevisions) undeleted(ctx context.Context, undel *types.Record) (err error) {
	var (
		invokerID = auth.GetIdentityFromContext(ctx).Identity()
		rev       *revisions.Revision
	)

	rev = revisions.Make(revisions.Undeleted, undel.Revision, undel.ID, invokerID)
	if err != nil {
		return
	}

	return svc.r.Create(ctx, svc.modelRef(undel.GetModule()), rev)
}

func (svc *recordRevisions) skippedField(mod *types.Module) []string {
	list := []string{
		"ID",
		"namespaceID",
		"moduleID",
		"revision",
		"meta",
		"createdBy",
		"createdAt",
		"updatedBy",
		"updatedAt",
		"deletedBy",
		"deletedAt",
	}

	for _, f := range mod.Fields {
		if f.Config.RecordRevisions.Skip {
			list = append(list, f.Name)
		}
	}

	return list
}

// @todo uncomment when supported
//func (svc *recordRevisions) hardDeleted(ctx context.Context, del *types.Record) (err error) {
//	var (
//		invokerID = auth.GetIdentityFromContext(ctx).Identity()
//		rev       *revisions.Revision
//	)
//
//	rev, err = revisions.Make(revisions.HardDeleted, del.Revision, del.ID, invokerID, nil, nil)
//	if err != nil {
//		return
//	}
//
//	return svc.r.Create(ctx, svc.modelRef(del.GetModule()), rev)
//}
