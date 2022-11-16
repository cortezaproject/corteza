package service

import (
	"context"
	"github.com/cortezaproject/corteza/server/pkg/discovery/types"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	resourceActivity struct {
		store store.Storer
	}

	ResourceActivityService interface {
		Find(context.Context, types.ResourceActivityFilter) (types.ResourceActivitySet, types.ResourceActivityFilter, error)
	}
)

func ResourceActivity() *resourceActivity {
	return &resourceActivity{
		store: DefaultStore,
	}
}

func (svc resourceActivity) Find(ctx context.Context, filter types.ResourceActivityFilter) (aa types.ResourceActivitySet, f types.ResourceActivityFilter, err error) {
	err = func() error {
		aa, f, err = store.SearchResourceActivityLogs(ctx, svc.store, filter)
		if err != nil {
			return err
		}

		return nil
	}()

	return aa, f, err
}
