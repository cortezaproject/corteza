package rest

import (
	"context"

	"github.com/crusttech/human/server/compose/rest/request"
	"github.com/crusttech/human/server/compose/service"
	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/api"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/payload"
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
			Search(ctx context.Context, filter types.PageLayoutFilter) (set types.PageLayoutSet, f types.PageLayoutFilter, err error)

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

	set, filter, err := ctrl.pageLayout.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
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

// makeFilter builds the search filter for the generated List controller.
func (ctrl PageLayout) makeFilter(ctx context.Context, r *request.PageLayoutList) (f types.PageLayoutFilter, err error) {
	f = types.PageLayoutFilter{
		NamespaceID: r.NamespaceID,
		PageID:      r.PageID,
		Labels:      r.Labels,

		Handle: r.Handle,
		Query:  r.Query,
	}

	if f.Paging, err = filter.NewPaging(r.Limit, r.PageCursor); err != nil {
		return
	}

	if f.Sorting, err = filter.NewSorting(r.Sort); err != nil {
		return
	}

	return
}

// beforeCreate assigns the non-plain create params onto the resource for the
// generated Create controller: the per-op page parent, meta, and the JSON
// config/blocks payloads.
func (ctrl PageLayout) beforeCreate(ctx context.Context, res *types.PageLayout, r *request.PageLayoutCreate) (err error) {
	res.PageID = r.PageID
	res.Meta = r.Meta

	if len(r.Config) > 2 {
		if err = r.Config.Unmarshal(&res.Config); err != nil {
			return err
		}
	}

	if len(r.Blocks) > 2 {
		if err = r.Blocks.Unmarshal(&res.Blocks); err != nil {
			return err
		}
	}

	return nil
}

// beforeUpdate assigns the non-plain update params onto the resource for the
// generated Update controller: the per-op page parent, meta, and the JSON
// config/blocks payloads.
func (ctrl PageLayout) beforeUpdate(ctx context.Context, res *types.PageLayout, r *request.PageLayoutUpdate) (err error) {
	res.PageID = r.PageID
	res.Meta = r.Meta

	if len(r.Config) > 2 {
		// Process config if it was included in the request
		// if not, do not assume that config has been removed!
		if err = r.Config.Unmarshal(&res.Config); err != nil {
			return err
		}
	}

	if len(r.Blocks) > 2 {
		// Process blocks if they were included in the request
		// if not, do not assume that blocks were removed!
		if err = r.Blocks.Unmarshal(&res.Blocks); err != nil {
			return err
		}
	}

	return nil
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
