package service

import (
	"context"

	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service/event"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
)

type (
	page struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac pageAccessController

		pageRepo   repository.PageRepository
		moduleRepo repository.ModuleRepository
		nsRepo     repository.NamespaceRepository
	}

	pageAccessController interface {
		CanReadNamespace(context.Context, *types.Namespace) bool
		CanCreatePage(context.Context, *types.Namespace) bool
		CanReadPage(context.Context, *types.Page) bool
		CanUpdatePage(context.Context, *types.Page) bool
		CanDeletePage(context.Context, *types.Page) bool

		FilterReadablePages(ctx context.Context) *permissions.ResourceFilter
	}

	PageService interface {
		With(ctx context.Context) PageService

		FindByID(namespaceID, pageID uint64) (*types.Page, error)
		FindByHandle(namespaceID uint64, handle string) (*types.Page, error)
		FindByModuleID(namespaceID, moduleID uint64) (*types.Page, error)
		FindBySelfID(namespaceID, selfID uint64) (pages types.PageSet, f types.PageFilter, err error)
		Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error)
		Tree(namespaceID uint64) (pages types.PageSet, err error)

		Create(page *types.Page) (*types.Page, error)
		Update(page *types.Page) (*types.Page, error)
		DeleteByID(namespaceID, pageID uint64) error

		Reorder(namespaceID, selfID uint64, pageIDs []uint64) error
	}
)

func Page() PageService {
	return (&page{
		logger: DefaultLogger.Named("page"),
		ac:     DefaultAccessControl,
	}).With(context.Background())
}

func (svc page) With(ctx context.Context) PageService {
	db := repository.DB(ctx)
	return &page{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,

		ac: svc.ac,

		pageRepo:   repository.Page(ctx, db),
		moduleRepo: repository.Module(ctx, db),
		nsRepo:     repository.Namespace(ctx, db),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc page) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc page) FindByID(namespaceID, pageID uint64) (p *types.Page, err error) {
	if _, err = svc.loadNamespace(namespaceID); err != nil {
		return
	} else {
		return svc.checkPermissions(svc.pageRepo.FindByID(namespaceID, pageID))
	}
}

func (svc page) FindByHandle(namespaceID uint64, handle string) (c *types.Page, err error) {
	if _, err = svc.loadNamespace(namespaceID); err != nil {
		return
	} else {
		return svc.checkPermissions(svc.pageRepo.FindByHandle(namespaceID, handle))
	}
}

func (svc page) FindByModuleID(namespaceID, moduleID uint64) (p *types.Page, err error) {
	if _, err = svc.loadNamespace(namespaceID); err != nil {
		return
	} else {
		return svc.checkPermissions(svc.pageRepo.FindByModuleID(namespaceID, moduleID))
	}
}

func (svc page) checkPermissions(p *types.Page, err error) (*types.Page, error) {
	if err != nil {
		return nil, err
	} else if !svc.ac.CanReadPage(svc.ctx, p) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return p, err
}

func (svc page) FindBySelfID(namespaceID, parentID uint64) (pp types.PageSet, f types.PageFilter, err error) {
	if namespaceID == 0 {
		return nil, f, ErrNamespaceRequired.withStack()
	}

	return svc.pageRepo.Find(types.PageFilter{
		NamespaceID: namespaceID,
		ParentID:    parentID,

		// This will enable parentID=0 query
		Root: true,

		IsReadable: svc.ac.FilterReadablePages(svc.ctx),
	})
}

func (svc page) Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	filter.IsReadable = svc.ac.FilterReadablePages(svc.ctx)

	if filter.NamespaceID == 0 {
		return nil, f, ErrNamespaceRequired.withStack()
	}

	return svc.pageRepo.Find(filter)
}

func (svc page) Tree(namespaceID uint64) (pages types.PageSet, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}

	var (
		tree   types.PageSet
		filter = types.PageFilter{
			NamespaceID: namespaceID,
			IsReadable:  svc.ac.FilterReadablePages(svc.ctx),
			Sort:        "weight ASC",
		}
	)

	return tree, svc.db.Transaction(func() (err error) {
		if pages, _, err = svc.pageRepo.Find(filter); err != nil {
			return
		}

		// No preloading - we do not need (or should have) any modules
		// associated with us
		_ = pages.Walk(func(p *types.Page) error {
			if p.SelfID == 0 {
				tree = append(tree, p)
			} else if c := pages.FindByID(p.SelfID); c != nil {
				if c.Children == nil {
					c.Children = types.PageSet{}
				}

				c.Children = append(c.Children, p)
			} else {
				// Move orphans to root
				p.SelfID = 0
				tree = append(tree, p)
			}

			return nil
		})

		return nil
	})
}

func (svc page) Reorder(namespaceID, selfID uint64, pageIDs []uint64) error {
	if ns, err := svc.loadNamespace(namespaceID); err != nil {
		return err
	} else if selfID == 0 {
		// Reordering on root mode
		if !svc.ac.CanCreatePage(svc.ctx, ns) {
			return ErrNoUpdatePermissions.withStack()
		}
	} else if p, err := svc.checkPermissions(svc.pageRepo.FindByID(ns.ID, selfID)); err != nil {
		return err
	} else if !svc.ac.CanUpdatePage(svc.ctx, p) {
		return ErrNoUpdatePermissions.withStack()
	}

	return svc.pageRepo.Reorder(namespaceID, selfID, pageIDs)
}

func (svc page) Create(new *types.Page) (p *types.Page, err error) {
	var (
		ns *types.Namespace
	)

	new.ID = 0

	if !handle.IsValid(new.Handle) {
		return nil, ErrInvalidHandle
	}

	if ns, err = svc.loadNamespace(new.NamespaceID); err != nil {
		return nil, err
	} else if !svc.ac.CanCreatePage(svc.ctx, ns) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	if err = eventbus.WaitFor(svc.ctx, event.PageBeforeCreate(new, nil, ns)); err != nil {
		return
	}

	if err = svc.UniqueCheck(new); err != nil {
		return
	}

	if p, err = svc.pageRepo.Create(new); err != nil {
		return
	}

	defer eventbus.Dispatch(svc.ctx, event.PageAfterCreate(new, nil, ns))
	return
}

func (svc page) Update(upd *types.Page) (p *types.Page, err error) {
	var (
		ns *types.Namespace
	)

	if upd.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if !handle.IsValid(upd.Handle) {
		return nil, ErrInvalidHandle
	}

	if ns, err = svc.loadNamespace(upd.NamespaceID); err != nil {
		return
	}

	if p, err = svc.pageRepo.FindByID(upd.NamespaceID, upd.ID); err != nil {
		return
	}

	if isStale(upd.UpdatedAt, p.UpdatedAt, p.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if !svc.ac.CanUpdatePage(svc.ctx, p) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	if err = eventbus.WaitFor(svc.ctx, event.PageBeforeUpdate(upd, p, ns)); err != nil {
		return
	}

	if err = svc.UniqueCheck(upd); err != nil {
		return
	}

	p.ModuleID = upd.ModuleID
	p.SelfID = upd.SelfID
	p.Blocks = upd.Blocks
	p.Title = upd.Title
	p.Handle = upd.Handle
	p.Description = upd.Description
	p.Visible = upd.Visible
	p.Weight = upd.Weight

	if p, err = svc.pageRepo.Update(p); err != nil {
		return
	}

	defer eventbus.Dispatch(svc.ctx, event.PageAfterUpdate(upd, p, ns))
	return
}

func (svc page) UniqueCheck(p *types.Page) (err error) {
	if p.Handle != "" {
		if e, _ := svc.pageRepo.FindByHandle(p.NamespaceID, p.Handle); e != nil && e.ID != p.ID {
			return repository.ErrPageHandleNotUnique
		}
	}

	if p.ModuleID > 0 {
		if e, _ := svc.pageRepo.FindByModuleID(p.NamespaceID, p.ModuleID); e != nil && e.ID != e.ID {
			return ErrModulePageExists
		}
	}

	return nil
}

func (svc page) DeleteByID(namespaceID, pageID uint64) (err error) {
	var (
		del *types.Page
		ns  *types.Namespace
	)

	if pageID == 0 {
		return ErrInvalidID.withStack()
	}

	if ns, err = svc.loadNamespace(namespaceID); err != nil {
		return err
	}

	if del, err = svc.pageRepo.FindByID(namespaceID, pageID); err != nil {
		return err
	} else if !svc.ac.CanDeletePage(svc.ctx, del) {
		return ErrNoDeletePermissions.withStack()
	}

	if err = eventbus.WaitFor(svc.ctx, event.PageBeforeDelete(nil, del, ns)); err != nil {
		return
	}

	if err = svc.pageRepo.DeleteByID(namespaceID, pageID); err != nil {
		return
	}

	defer eventbus.Dispatch(svc.ctx, event.PageAfterDelete(nil, del, ns))
	return
}

func (svc page) loadNamespace(namespaceID uint64) (ns *types.Namespace, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}

	if ns, err = svc.nsRepo.FindByID(namespaceID); err != nil {
		return
	}

	if !svc.ac.CanReadNamespace(svc.ctx, ns) {
		return nil, ErrNoReadPermissions.withStack()
	}

	return
}
