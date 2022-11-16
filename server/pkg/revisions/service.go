package revisions

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/filter"
)

type (
	Servicer interface {
		Search(ctx context.Context, mf dal.ModelRef, f filter.Filter) (_ dal.Iterator, err error)
		Create(ctx context.Context, mf dal.ModelRef, revision *Revision) error
	}

	creatorSearcher interface {
		Search(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, f filter.Filter) (dal.Iterator, error)
		Create(ctx context.Context, m dal.ModelRef, operations dal.OperationSet, vv ...dal.ValueGetter) error
	}

	service struct {
		dal creatorSearcher
	}
)

func Service(dal creatorSearcher) *service {
	return &service{dal: dal}
}

func (svc *service) Search(ctx context.Context, mf dal.ModelRef, f filter.Filter) (_ dal.Iterator, err error) {
	return svc.dal.Search(ctx, mf, dal.OperationSet{dal.Search}, f)
}

func (svc *service) Create(ctx context.Context, mf dal.ModelRef, revision *Revision) error {
	return svc.dal.Create(ctx, mf, dal.OperationSet{dal.Create}, revision)

}
