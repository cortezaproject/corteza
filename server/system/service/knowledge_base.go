package service

import (
	"context"

	a "github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/errors"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
)

type (
	knowledgeBase struct {
		actionlog actionlog.Recorder
		store     store.Storer
		ac        knowledgeBaseAccessController
	}

	knowledgeBaseAccessController interface {
		CanCreateKnowledgeBase(ctx context.Context) bool
		CanSearchKnowledgeBases(ctx context.Context) bool
		CanReadKnowledgeBase(ctx context.Context, kb *types.KnowledgeBase) bool
		CanUpdateKnowledgeBase(ctx context.Context, kb *types.KnowledgeBase) bool
		CanDeleteKnowledgeBase(ctx context.Context, kb *types.KnowledgeBase) bool
	}
)

func KnowledgeBase() *knowledgeBase {
	return &knowledgeBase{
		ac:        DefaultAccessControl,
		actionlog: DefaultActionlog,
		store:     DefaultStore,
	}
}

func (svc *knowledgeBase) FindByID(ctx context.Context, ID uint64) (kb *types.KnowledgeBase, err error) {
	err = func() error {
		if kb, err = loadKnowledgeBase(ctx, svc.store, ID); err != nil {
			return err
		}

		if !svc.ac.CanReadKnowledgeBase(ctx, kb) {
			return KnowledgeBaseErrNotAllowedToRead()
		}

		return nil
	}()

	return kb, err
}

func (svc *knowledgeBase) Create(ctx context.Context, new *types.KnowledgeBase) (kb *types.KnowledgeBase, err error) {
	err = func() (err error) {
		if !svc.ac.CanCreateKnowledgeBase(ctx) {
			return KnowledgeBaseErrNotAllowedToCreate()
		}

		new.ID = nextID()
		new.CreatedAt = *now()
		new.CreatedBy = a.GetIdentityFromContext(ctx).Identity()

		if err = store.CreateKnowledgeBase(ctx, svc.store, new); err != nil {
			return
		}

		kb = new
		return nil
	}()

	return kb, err
}

func (svc *knowledgeBase) Update(ctx context.Context, upd *types.KnowledgeBase) (kb *types.KnowledgeBase, err error) {
	err = func() (err error) {
		if !svc.ac.CanUpdateKnowledgeBase(ctx, upd) {
			return KnowledgeBaseErrNotAllowedToUpdate()
		}

		var existing *types.KnowledgeBase
		if existing, err = store.LookupKnowledgeBaseByID(ctx, svc.store, upd.ID); err != nil {
			return KnowledgeBaseErrNotFound()
		}

		if isStale(upd.UpdatedAt, existing.UpdatedAt, existing.CreatedAt) {
			return KnowledgeBaseErrStaleData()
		}

		upd.UpdatedAt = now()
		upd.UpdatedBy = a.GetIdentityFromContext(ctx).Identity()
		upd.CreatedAt = existing.CreatedAt
		upd.CreatedBy = existing.CreatedBy
		upd.DeletedAt = existing.DeletedAt

		if err = store.UpdateKnowledgeBase(ctx, svc.store, upd); err != nil {
			return
		}

		kb = upd
		return nil
	}()

	return kb, err
}

func (svc *knowledgeBase) DeleteByID(ctx context.Context, ID uint64) (err error) {
	err = func() (err error) {
		var kb *types.KnowledgeBase
		if kb, err = loadKnowledgeBase(ctx, svc.store, ID); err != nil {
			return
		}

		if !svc.ac.CanDeleteKnowledgeBase(ctx, kb) {
			return KnowledgeBaseErrNotAllowedToDelete()
		}

		kb.DeletedAt = now()
		kb.DeletedBy = a.GetIdentityFromContext(ctx).Identity()
		if err = store.UpdateKnowledgeBase(ctx, svc.store, kb); err != nil {
			return
		}

		return nil
	}()

	return err
}

func (svc *knowledgeBase) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	err = func() (err error) {
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
	}()

	return err
}

func (svc *knowledgeBase) Search(ctx context.Context, filter types.KnowledgeBaseFilter) (set types.KnowledgeBaseSet, f types.KnowledgeBaseFilter, err error) {
	filter.Check = func(res *types.KnowledgeBase) (bool, error) {
		if !svc.ac.CanReadKnowledgeBase(ctx, res) {
			return false, nil
		}

		return true, nil
	}

	err = func() error {
		if !svc.ac.CanSearchKnowledgeBases(ctx) {
			return KnowledgeBaseErrNotAllowedToSearch()
		}

		if set, f, err = store.SearchKnowledgeBases(ctx, svc.store, filter); err != nil {
			return err
		}

		return nil
	}()

	return set, f, err
}

func loadKnowledgeBase(ctx context.Context, s store.KnowledgeBases, ID uint64) (res *types.KnowledgeBase, err error) {
	if ID == 0 {
		return nil, KnowledgeBaseErrInvalidID()
	}

	if res, err = store.LookupKnowledgeBaseByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, KnowledgeBaseErrNotFound()
	}

	return
}
