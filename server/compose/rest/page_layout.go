package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/payload"
)

type (
	pageLayoutPayload struct {
		*types.PageLayout

		// CanGrant      bool `json:"canGrant"`
		// CanUpdatePageLayout bool `json:"canUpdatePageLayout"`
		// CanDeletePageLayout bool `json:"canDeletePageLayout"`
	}

	pageLayoutSetPayload struct {
		Filter types.PageLayoutFilter `json:"filter"`
		Set    []*pageLayoutPayload   `json:"set"`
	}

	PageLayout struct {
		pageLayout interface {
			FindByID(ctx context.Context, namespaceID, pageLayoutID uint64) (*types.PageLayout, error)
			FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*types.PageLayout, error)
			FindByPageLayoutID(ctx context.Context, namespaceID, pageLayoutID uint64) (*types.PageLayout, error)
			Find(ctx context.Context, filter types.PageLayoutFilter) (set types.PageLayoutSet, f types.PageLayoutFilter, err error)

			Create(ctx context.Context, pageLayout *types.PageLayout) (*types.PageLayout, error)
			Reorder(ctx context.Context, namespaceID uint64, pageID uint64, pageLayoutIDs []uint64) error
			Update(ctx context.Context, pageLayout *types.PageLayout) (*types.PageLayout, error)
			DeleteByID(ctx context.Context, namespaceID, pageID, pageLayoutID uint64) error
			UndeleteByID(ctx context.Context, namespaceID, pageID, pageLayoutID uint64) error
		}
		locale    service.ResourceTranslationsManagerService
		namespace service.NamespaceService
		ac        pageLayoutAccessController
	}

	pageLayoutAccessController interface {
		// @todo
	}
)

func (PageLayout) New() *PageLayout {
	return &PageLayout{
		pageLayout: service.DefaultPageLayout,
		locale:     service.DefaultResourceTranslation,
		namespace:  service.DefaultNamespace,
		ac:         service.DefaultAccessControl,
	}
}

func (ctrl *PageLayout) List(ctx context.Context, r *request.PageLayoutList) (interface{}, error) {
	var (
		err error
		f   = types.PageLayoutFilter{
			NamespaceID: r.NamespaceID,
			PageID:      r.PageID,
			Labels:      r.Labels,

			Handle: r.Handle,
			Query:  r.Query,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.pageLayout.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *PageLayout) ListNamespace(ctx context.Context, r *request.PageLayoutListNamespace) (interface{}, error) {
	var (
		err error
		f   = types.PageLayoutFilter{
			NamespaceID: r.NamespaceID,
			Labels:      r.Labels,

			Handle: r.Handle,
			Query:  r.Query,
		}
	)

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return nil, err
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return nil, err
	}

	set, filter, err := ctrl.pageLayout.Find(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *PageLayout) Create(ctx context.Context, r *request.PageLayoutCreate) (interface{}, error) {
	var (
		err    error
		layout = &types.PageLayout{
			PageID:      r.PageID,
			ParentID:    r.ParentID,
			Weight:      r.Weight,
			NamespaceID: r.NamespaceID,
			Handle:      r.Handle,
			Meta:        r.Meta,
			Labels:      r.Labels,
			OwnedBy:     r.OwnedBy,
		}
	)

	if len(r.Config) > 2 {
		if err = r.Config.Unmarshal(&layout.Config); err != nil {
			return nil, err
		}
	}

	if len(r.Blocks) > 2 {
		if err = r.Blocks.Unmarshal(&layout.Blocks); err != nil {
			return nil, err
		}
	}

	layout, err = ctrl.pageLayout.Create(ctx, layout)
	return ctrl.makePayload(ctx, layout, err)
}

func (ctrl *PageLayout) Read(ctx context.Context, r *request.PageLayoutRead) (interface{}, error) {
	mod, err := ctrl.pageLayout.FindByID(ctx, r.NamespaceID, r.PageLayoutID)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *PageLayout) ListTranslations(ctx context.Context, r *request.PageLayoutListTranslations) (interface{}, error) {
	return ctrl.locale.PageLayout(ctx, r.NamespaceID, r.PageID, r.PageLayoutID)
}

func (ctrl *PageLayout) UpdateTranslations(ctx context.Context, r *request.PageLayoutUpdateTranslations) (interface{}, error) {
	return api.OK(), ctrl.locale.Upsert(ctx, r.Translations)
}

func (ctrl *PageLayout) Reorder(ctx context.Context, r *request.PageLayoutReorder) (interface{}, error) {
	return api.OK(), ctrl.pageLayout.Reorder(ctx, r.NamespaceID, r.PageID, payload.ParseUint64s(r.PageIDs))
}

func (ctrl *PageLayout) Update(ctx context.Context, r *request.PageLayoutUpdate) (interface{}, error) {
	var (
		err error
		mod = &types.PageLayout{
			ID:          r.PageLayoutID,
			PageID:      r.PageID,
			Weight:      r.Weight,
			ParentID:    r.ParentID,
			NamespaceID: r.NamespaceID,
			Handle:      r.Handle,
			Meta:        r.Meta,
			Labels:      r.Labels,
			OwnedBy:     r.OwnedBy,
			UpdatedAt:   r.UpdatedAt,
		}
	)

	if len(r.Config) > 2 {
		// Process config if it was included in the request
		// if not, do not assume that config has been removed!
		if err = r.Config.Unmarshal(&mod.Config); err != nil {
			return nil, err
		}
	}

	if len(r.Blocks) > 2 {
		// Process blocks if they were included in the request
		// if not, do not assume that blocks were removed!
		if err = r.Blocks.Unmarshal(&mod.Blocks); err != nil {
			return nil, err
		}
	}

	mod, err = ctrl.pageLayout.Update(ctx, mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *PageLayout) Delete(ctx context.Context, r *request.PageLayoutDelete) (interface{}, error) {
	return api.OK(), ctrl.pageLayout.DeleteByID(ctx, r.NamespaceID, r.PageID, r.PageLayoutID)
}

func (ctrl *PageLayout) Undelete(ctx context.Context, r *request.PageLayoutUndelete) (interface{}, error) {
	return api.OK(), ctrl.pageLayout.UndeleteByID(ctx, r.NamespaceID, r.PageID, r.PageLayoutID)
}

func (ctrl PageLayout) makePayload(ctx context.Context, c *types.PageLayout, err error) (*pageLayoutPayload, error) {
	if err != nil || c == nil {
		return nil, err
	}

	return &pageLayoutPayload{
		PageLayout: c,
	}, nil
}

func (ctrl PageLayout) makeFilterPayload(ctx context.Context, nn types.PageLayoutSet, f types.PageLayoutFilter, err error) (*pageLayoutSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &pageLayoutSetPayload{Filter: f, Set: make([]*pageLayoutPayload, len(nn))}

	for i := range nn {
		modp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return modp, nil
}
