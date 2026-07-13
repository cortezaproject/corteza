package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

func KnowledgeBase() *knowledgeBase {
	return &knowledgeBase{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
	}
}

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
func (svc *knowledgeBase) onUpdate(ctx context.Context, s store.Storer, upd, res *types.KnowledgeBase, aProps *knowledgeBaseActionProps, _ func() error, _ func() error) error {
	if !svc.ac.CanUpdateKnowledgeBase(ctx, upd) {
		return KnowledgeBaseErrNotAllowedToUpdate()
	}

	upd.UpdatedAt = now()
	upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	upd.CreatedAt = res.CreatedAt
	upd.CreatedBy = res.CreatedBy
	upd.DeletedAt = res.DeletedAt

	if err := store.UpdateKnowledgeBase(ctx, s, upd); err != nil {
		return err
	}

	*res = *upd
	return nil
}

// onUndelete is the bespoke body for the generated UndeleteByID. It clears the
// soft-delete marker while also stamping UpdatedAt / UpdatedBy, and reuses the
// delete access check (there is no dedicated undelete permission/error).
func (svc *knowledgeBase) onUndelete(ctx context.Context, s store.Storer, res *types.KnowledgeBase, aProps *knowledgeBaseActionProps) error {
	if !svc.ac.CanDeleteKnowledgeBase(ctx, res) {
		return KnowledgeBaseErrNotAllowedToDelete()
	}

	res.DeletedAt = nil
	res.UpdatedAt = now()
	res.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
	return store.UpdateKnowledgeBase(ctx, s, res)
}
