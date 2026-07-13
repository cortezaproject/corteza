package service

import (
	"context"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/store"
)

func loadChartScoped(ctx context.Context, s store.ComposeCharts, namespaceID, chartID uint64) (res *types.Chart, err error) {
	if chartID == 0 || namespaceID == 0 {
		return nil, ChartErrInvalidID()
	}
	if res, err = store.LookupComposeChartByID(ctx, s, chartID); errors.IsNotFound(err) {
		return nil, ChartErrNotFound()
	} else if err != nil {
		return nil, err
	}
	if res.NamespaceID != namespaceID {
		return nil, ChartErrNotFound()
	}
	return
}

func loadModuleScoped(ctx context.Context, s store.ComposeModules, namespaceID, moduleID uint64) (res *types.Module, err error) {
	if moduleID == 0 || namespaceID == 0 {
		return nil, ModuleErrInvalidID()
	}
	if res, err = store.LookupComposeModuleByID(ctx, s, moduleID); errors.IsNotFound(err) {
		return nil, ModuleErrNotFound()
	} else if err != nil {
		return nil, err
	}
	if res.NamespaceID != namespaceID {
		return nil, ModuleErrNotFound()
	}
	return
}

func loadModuleWithFields(ctx context.Context, s store.Storer, namespaceID, moduleID uint64) (res *types.Module, err error) {
	if res, err = loadModuleScoped(ctx, s, namespaceID, moduleID); err != nil {
		return
	}
	err = loadModuleFields(ctx, s, res)
	return
}

func loadModuleFieldScoped(ctx context.Context, s store.Storer, namespaceID, moduleID, fieldID uint64) (res *types.ModuleField, err error) {
	return loadModuleField(ctx, s, namespaceID, moduleID, fieldID)
}

func loadPageScoped(ctx context.Context, s store.ComposePages, namespaceID, pageID uint64) (res *types.Page, err error) {
	if pageID == 0 || namespaceID == 0 {
		return nil, PageErrInvalidID()
	}
	if res, err = store.LookupComposePageByID(ctx, s, pageID); errors.IsNotFound(err) {
		return nil, PageErrNotFound()
	} else if err != nil {
		return nil, err
	}
	if res.NamespaceID != namespaceID {
		return nil, PageErrNotFound()
	}
	return
}

func loadPageLayoutScoped(ctx context.Context, s store.ComposePageLayouts, namespaceID, pageID, pageLayoutID uint64) (res *types.PageLayout, err error) {
	if pageLayoutID == 0 || namespaceID == 0 {
		return nil, PageLayoutErrInvalidID()
	}
	if res, err = store.LookupComposePageLayoutByID(ctx, s, pageLayoutID); errors.IsNotFound(err) {
		return nil, PageLayoutErrNotFound()
	} else if err != nil {
		return nil, err
	}
	if res.NamespaceID != namespaceID {
		return nil, PageLayoutErrNotFound()
	}
	return
}
