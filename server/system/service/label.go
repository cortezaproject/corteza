package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/label/types"
	"github.com/cortezaproject/corteza/server/store"  
)

type (
	labelSvc struct {
		store store.Storer
	}
	LabelService interface {
		List(ctx context.Context, filter types.LabelFilter) (types.LabelSet, types.LabelFilter, error)
	}
)

func Label() LabelService {
	return &labelSvc{
		store: DefaultStore,
	}
}

func (svc labelSvc) List(ctx context.Context, filter types.LabelFilter) (set types.LabelSet, f types.LabelFilter, err error) {
	set, f, err = store.SearchLabels(ctx, svc.store, filter)
	if err != nil {
		return nil, f, err
	}
	
	return set, f, nil
}