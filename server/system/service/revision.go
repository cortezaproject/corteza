package service

import (
	"context"

	cauth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/revisions"
	"github.com/cortezaproject/corteza/server/system/model"
)

type (
	revision struct {
		dal interface {
			Search(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (dal.Iterator, error)
			Create(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, vv ...dal.ValueGetter) error
			Update(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, rr ...dal.ValueGetter) error
		}
	}
)

func Revision(d dal.FullService) *revision {
	return &revision{dal: d}
}

func revisionModelRef() dal.ModelRef {
	return dal.ModelRef{
		ResourceType: model.Revision.ResourceType,
	}
}

func (svc *revision) Search(ctx context.Context, f revisions.Filter) (dal.Iterator, error) {
	return svc.dal.Search(ctx, revisionModelRef(), dal.OperationSet{dal.Search}, f)
}

func (svc *revision) Create(ctx context.Context, rev *revisions.Revision) error {
	return svc.dal.Create(ctx, revisionModelRef(), dal.OperationSet{dal.Create}, rev)
}

func (svc *revision) Delete(ctx context.Context, revisionID uint64) error {
	invokerID := cauth.GetIdentityFromContext(ctx).Identity()

	rev := &revisions.Revision{
		ID:        revisionID,
		DeletedAt: now(),
		DeletedBy: invokerID,
	}

	return svc.dal.Update(ctx, revisionModelRef(), dal.OperationSet{dal.Update}, rev)
}

