package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

// The CRUD skeleton (FindByID, Search, Create, DeleteByID, loadKnowledgeBase),
// the struct, the access-controller interface and the constructor are generated
// in knowledge_base.gen.go from system/knowledge_base.cue.
//
// This file owns the before-create / before-delete hooks the generated Create /
// DeleteByID call into (audit-author bookkeeping), and the bespoke Update /
// Undelete bodies (onUpdate / onUndelete) the generated wrappers delegate to.

// beforeCreate runs after the access check and before the generated Create
// assigns the ID / timestamps and persists. It stamps the creating identity.
func (svc *knowledgeBase) beforeCreate(ctx context.Context, new *types.KnowledgeBase) error {
	new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

// beforeDelete runs after the access check and before the generated DeleteByID
// sets DeletedAt and persists. It stamps the deleting identity.
func (svc *knowledgeBase) beforeDelete(ctx context.Context, res *types.KnowledgeBase) error {
	res.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
	return nil
}

// onUpdate is the bespoke body for the generated Update. Unlike the standard
// scaffold (load existing, copy mutable fields onto it, persist the loaded
// record), the original service checks access on the incoming `upd`, copies the
// existing audit fields ONTO `upd`, and persists `upd` itself.
func (svc *knowledgeBase) onUpdate(ctx context.Context, upd *types.KnowledgeBase, aProps *knowledgeBaseActionProps) (kb *types.KnowledgeBase, err error) {
	if !svc.ac.CanUpdateKnowledgeBase(ctx, upd) {
		return nil, KnowledgeBaseErrNotAllowedToUpdate()
	}

	var existing *types.KnowledgeBase
	if existing, err = store.LookupKnowledgeBaseByID(ctx, svc.store, upd.ID); err != nil {
		return nil, KnowledgeBaseErrNotFound()
	}

	if isStale(upd.UpdatedAt, existing.UpdatedAt, existing.CreatedAt) {
		return nil, KnowledgeBaseErrStaleData()
	}

	upd.UpdatedAt = now()
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	upd.CreatedAt = existing.CreatedAt
	upd.CreatedBy = existing.CreatedBy
	upd.DeletedAt = existing.DeletedAt

	if err = store.UpdateKnowledgeBase(ctx, svc.store, upd); err != nil {
		return nil, err
	}

	kb = upd
	return kb, nil
}

// onUndelete is the bespoke body for the generated UndeleteByID. It clears the
// soft-delete marker while also stamping UpdatedAt / UpdatedBy, and reuses the
// delete access check (there is no dedicated undelete permission/error).
func (svc *knowledgeBase) onUndelete(ctx context.Context, ID uint64, aProps *knowledgeBaseActionProps) (err error) {
	var kb *types.KnowledgeBase
	if kb, err = loadKnowledgeBase(ctx, svc.store, ID); err != nil {
		return
	}

	if !svc.ac.CanDeleteKnowledgeBase(ctx, kb) {
		return KnowledgeBaseErrNotAllowedToDelete()
	}

	kb.DeletedAt = nil
	kb.UpdatedAt = now()
	kb.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	if err = store.UpdateKnowledgeBase(ctx, svc.store, kb); err != nil {
		return
	}

	return nil
}
