package rest

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/rest/request"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/service/event"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/corredor"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/payload"
)

type (
	pagePayload struct {
		*types.Page

		Children []*pagePayload `json:"children,omitempty"`

		CanGrant      bool `json:"canGrant"`
		CanExportPage bool `json:"canExportPage"`
		CanUpdatePage bool `json:"canUpdatePage"`
		CanDeletePage bool `json:"canDeletePage"`
	}

	pageSetPayload struct {
		Filter types.PageFilter `json:"filter"`
		Set    []*pagePayload   `json:"set"`
	}

	pageIconPayload struct {
		*types.PageConfigIcon
	}

	Page struct {
		page interface {
			FindByID(ctx context.Context, namespaceID, pageID uint64) (*types.Page, error)
			FindByHandle(ctx context.Context, namespaceID uint64, handle string) (*types.Page, error)
			FindByPageID(ctx context.Context, namespaceID, pageID uint64) (*types.Page, error)
			FindBySelfID(ctx context.Context, namespaceID, selfID uint64) (pages types.PageSet, f types.PageFilter, err error)
			Find(ctx context.Context, filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error)
			Tree(ctx context.Context, namespaceID uint64) (pages types.PageSet, err error)

			Create(ctx context.Context, page *types.Page) (*types.Page, error)
			Update(ctx context.Context, page *types.Page) (*types.Page, error)
			DeleteByID(ctx context.Context, namespaceID, pageID uint64, pds types.PageChildrenDeleteStrategy) error

			UpdateIcon(ctx context.Context, namespaceID, pageID uint64, icon *types.PageConfigIcon) (out *types.PageConfigIcon, err error)

			Reorder(ctx context.Context, namespaceID, selfID uint64, pageIDs []uint64) error
		}
		locale     service.ResourceTranslationsManagerService
		namespace  service.NamespaceService
		attachment service.AttachmentService
		ac         pageAccessController
	}

	pageAccessController interface {
		CanGrant(context.Context) bool

		CanUpdatePage(context.Context, *types.Page) bool
		CanExportPage(context.Context, *types.Page) bool
		CanDeletePage(context.Context, *types.Page) bool
	}
)

func (Page) New() *Page {
	return &Page{
		page:       service.DefaultPage,
		locale:     service.DefaultResourceTranslation,
		namespace:  service.DefaultNamespace,
		attachment: service.DefaultAttachment,
		ac:         service.DefaultAccessControl,
	}
}

func (ctrl *Page) List(ctx context.Context, r *request.PageList) (interface{}, error) {
	var (
		err error
		f   = types.PageFilter{
			NamespaceID: r.NamespaceID,
			ModuleID:    r.ModuleID,
			ParentID:    r.SelfID,
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

	set, filter, err := ctrl.page.Find(ctx, f)

	return ctrl.makeFilterPayload(ctx, set, filter, err)
}

func (ctrl *Page) Tree(ctx context.Context, r *request.PageTree) (interface{}, error) {
	tree, err := ctrl.page.Tree(ctx, r.NamespaceID)
	return ctrl.makeTreePayload(ctx, tree, err)
}

func (ctrl *Page) Create(ctx context.Context, r *request.PageCreate) (interface{}, error) {
	var (
		err error
		mod = &types.Page{
			NamespaceID: r.NamespaceID,
			SelfID:      r.SelfID,
			ModuleID:    r.ModuleID,
			Title:       r.Title,
			Handle:      r.Handle,
			Description: r.Description,
			Visible:     r.Visible,
			Weight:      r.Weight,
			Labels:      r.Labels,
			Meta:        r.Meta,
		}
	)

	if len(r.Config) > 2 {
		if err = r.Config.Unmarshal(&mod.Config); err != nil {
			return nil, err
		}
	}

	if len(r.Blocks) > 2 {
		if err = r.Blocks.Unmarshal(&mod.Blocks); err != nil {
			return nil, err
		}
	}

	mod, err = ctrl.page.Create(ctx, mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Page) Read(ctx context.Context, r *request.PageRead) (interface{}, error) {
	mod, err := ctrl.page.FindByID(ctx, r.NamespaceID, r.PageID)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Page) ListTranslations(ctx context.Context, r *request.PageListTranslations) (interface{}, error) {
	return ctrl.locale.Page(ctx, r.NamespaceID, r.PageID)
}

func (ctrl *Page) UpdateTranslations(ctx context.Context, r *request.PageUpdateTranslations) (interface{}, error) {
	return api.OK(), ctrl.locale.Upsert(ctx, r.Translations)
}

func (ctrl *Page) Reorder(ctx context.Context, r *request.PageReorder) (interface{}, error) {
	return api.OK(), ctrl.page.Reorder(ctx, r.NamespaceID, r.SelfID, payload.ParseUint64s(r.PageIDs))
}

func (ctrl *Page) Update(ctx context.Context, r *request.PageUpdate) (interface{}, error) {
	var (
		err error
		mod = &types.Page{
			NamespaceID: r.NamespaceID,
			ID:          r.PageID,
			SelfID:      r.SelfID,
			ModuleID:    r.ModuleID,
			Title:       r.Title,
			Handle:      r.Handle,
			Description: r.Description,
			Visible:     r.Visible,
			Weight:      r.Weight,
			Labels:      r.Labels,
			Meta:        r.Meta,
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

	mod, err = ctrl.page.Update(ctx, mod)
	return ctrl.makePayload(ctx, mod, err)
}

func (ctrl *Page) Delete(ctx context.Context, r *request.PageDelete) (interface{}, error) {
	var strategy types.PageChildrenDeleteStrategy

	switch aux := types.PageChildrenDeleteStrategy(r.Strategy); aux {
	case types.PageChildrenOnDeleteForce,
		types.PageChildrenOnDeleteRebase,
		types.PageChildrenOnDeleteCascade:
		strategy = aux
	default:
		strategy = types.PageChildrenOnDeleteAbort
	}

	return api.OK(), ctrl.page.DeleteByID(ctx, r.NamespaceID, r.PageID, strategy)
}

func (ctrl *Page) Upload(ctx context.Context, r *request.PageUpload) (interface{}, error) {
	file, err := r.Upload.Open()
	if err != nil {
		return nil, err
	}

	defer file.Close()

	a, err := ctrl.attachment.CreatePageAttachment(
		ctx,
		r.NamespaceID,
		r.Upload.Filename,
		r.Upload.Size,
		file,
		r.PageID,
	)

	return makeAttachmentPayload(ctx, a, err)
}

func (ctrl *Page) TriggerScript(ctx context.Context, r *request.PageTriggerScript) (rsp interface{}, err error) {
	var (
		page      *types.Page
		namespace *types.Namespace
	)

	if page, err = ctrl.page.FindByID(ctx, r.NamespaceID, r.PageID); err != nil {
		return
	}
	if namespace, err = ctrl.namespace.FindByID(ctx, r.NamespaceID); err != nil {
		return
	}

	// @todo implement same behaviour as we have on record - page+oldPage
	err = corredor.Service().Exec(ctx, r.Script, corredor.ExtendScriptArgs(event.PageOnManual(page, page, namespace, nil), r.Args))
	return ctrl.makePayload(ctx, page, err)
}

func (ctrl *Page) UpdateIcon(ctx context.Context, r *request.PageUpdateIcon) (interface{}, error) {
	var (
		err  error
		icon = &types.PageConfigIcon{
			Type:  r.Type,
			Src:   r.Source,
			Style: r.Style,
		}
	)

	icon, err = ctrl.page.UpdateIcon(ctx, r.NamespaceID, r.PageID, icon)
	return ctrl.makeIconPayload(ctx, icon, err)
}

func (ctrl Page) makePayload(ctx context.Context, c *types.Page, err error) (*pagePayload, error) {
	if err != nil || c == nil {
		return nil, err
	}

	return &pagePayload{
		Page: c,

		CanGrant: ctrl.ac.CanGrant(ctx),

		CanUpdatePage: ctrl.ac.CanUpdatePage(ctx, c),
		CanExportPage: ctrl.ac.CanExportPage(ctx, c),
		CanDeletePage: ctrl.ac.CanDeletePage(ctx, c),
	}, nil
}

func (ctrl Page) makeIconPayload(_ context.Context, i *types.PageConfigIcon, err error) (*pageIconPayload, error) {
	if err != nil || i == nil {
		return nil, err
	}

	return &pageIconPayload{
		PageConfigIcon: i,
	}, nil
}

func (ctrl Page) makeTreePayload(ctx context.Context, pp types.PageSet, err error) ([]*pagePayload, error) {
	if err != nil {
		return nil, err
	}

	set := make([]*pagePayload, len(pp))

	for i := range pp {
		set[i], err = ctrl.makePayload(ctx, pp[i], nil)
		if err != nil {
			return nil, err
		}

		if len(pp[i].Children) > 0 {
			set[i].Children, err = ctrl.makeTreePayload(ctx, pp[i].Children, nil)
			if err != nil {
				return nil, err
			}
		}
	}

	return set, nil
}

func (ctrl Page) makeFilterPayload(ctx context.Context, nn types.PageSet, f types.PageFilter, err error) (*pageSetPayload, error) {
	if err != nil {
		return nil, err
	}

	modp := &pageSetPayload{Filter: f, Set: make([]*pagePayload, len(nn))}

	for i := range nn {
		modp.Set[i], _ = ctrl.makePayload(ctx, nn[i], nil)
	}

	return modp, nil
}
