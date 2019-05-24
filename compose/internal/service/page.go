package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/compose/internal/repository"
	"github.com/cortezaproject/corteza-server/compose/types"
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
	}

	PageService interface {
		With(ctx context.Context) PageService

		FindByID(namespaceID, pageID uint64) (*types.Page, error)
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
// func (svc page) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

func (svc page) FindByID(namespaceID, pageID uint64) (p *types.Page, err error) {
	return svc.checkPermissions(svc.pageRepo.FindByID(namespaceID, pageID))
}

func (svc page) FindByModuleID(namespaceID, moduleID uint64) (p *types.Page, err error) {
	return svc.checkPermissions(svc.pageRepo.FindByModuleID(namespaceID, moduleID))
}

func (svc page) checkPermissions(p *types.Page, err error) (*types.Page, error) {
	if err != nil {
		return nil, err
	} else if !svc.ac.CanReadPage(svc.ctx, p) {
		return nil, errors.New("not allowed to access this page")
	}

	return p, err
}

func (svc page) FindBySelfID(namespaceID, parentID uint64) (pp types.PageSet, f types.PageFilter, err error) {
	if namespaceID == 0 {
		return nil, f, ErrNamespaceRequired.withStack()
	}

	return svc.filterPageSetByPermission(svc.pageRepo.Find(types.PageFilter{
		NamespaceID: namespaceID,
		ParentID:    parentID,

		// This will enable parentID=0 query
		Root: true,
	}))
}

func (svc page) Find(filter types.PageFilter) (set types.PageSet, f types.PageFilter, err error) {
	if filter.NamespaceID == 0 {
		return nil, f, ErrNamespaceRequired.withStack()
	}

	return svc.filterPageSetByPermission(svc.pageRepo.Find(filter))
}

func (svc page) Tree(namespaceID uint64) (pages types.PageSet, err error) {
	if namespaceID == 0 {
		return nil, ErrNamespaceRequired.withStack()
	}

	var (
		tree   types.PageSet
		filter = types.PageFilter{
			NamespaceID: namespaceID,
		}
	)

	return tree, svc.db.Transaction(func() (err error) {
		if pages, _, err = svc.filterPageSetByPermission(svc.pageRepo.Find(filter)); err != nil {
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

func (svc page) filterPageSetByPermission(pp types.PageSet, f types.PageFilter, err error) (types.PageSet, types.PageFilter, error) {
	if err != nil {
		return nil, f, err
	}

	// @todo Filter-by-permission can/will mess up filter's count & paging...
	pp, err = pp.Filter(func(m *types.Page) (bool, error) {
		return svc.ac.CanReadPage(svc.ctx, m), nil
	})

	return pp, f, err
}

func (svc page) Reorder(namespaceID, selfID uint64, pageIDs []uint64) error {
	if ns, err := svc.loadNamespace(namespaceID); err != nil {
		return err
	} else if p, err := svc.checkPermissions(svc.pageRepo.FindByID(ns.ID, selfID)); err != nil {
		return err
	} else if !svc.ac.CanUpdatePage(svc.ctx, p) {
		return ErrNoUpdatePermissions.withStack()
	}

	return svc.pageRepo.Reorder(namespaceID, selfID, pageIDs)
}

func (svc page) Create(mod *types.Page) (p *types.Page, err error) {
	mod.ID = 0

	if ns, err := svc.loadNamespace(mod.NamespaceID); err != nil {
		return nil, err
	} else if !svc.ac.CanCreatePage(svc.ctx, ns) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	if err = svc.checkModulePage(mod); err != nil {
		return
	}

	p, err = svc.pageRepo.Create(mod)
	return
}

func (svc page) Update(mod *types.Page) (p *types.Page, err error) {
	if mod.ID == 0 {
		return nil, ErrInvalidID.withStack()
	}

	if _, err = svc.loadNamespace(mod.NamespaceID); err != nil {
		return
	}

	if p, err = svc.pageRepo.FindByID(mod.NamespaceID, mod.ID); err != nil {
		return
	}

	if isStale(mod.UpdatedAt, p.UpdatedAt, p.CreatedAt) {
		return nil, ErrStaleData.withStack()
	}

	if !svc.ac.CanUpdatePage(svc.ctx, p) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	if err = svc.checkModulePage(mod); err != nil {
		return
	}

	p.ModuleID = mod.ModuleID
	p.SelfID = mod.SelfID
	p.Blocks = mod.Blocks
	p.Title = mod.Title
	p.Description = mod.Description
	p.Visible = mod.Visible
	p.Weight = mod.Weight

	p, err = svc.pageRepo.Update(p)
	return
}

func (svc page) checkModulePage(mod *types.Page) error {
	if mod.ModuleID > 0 {
		if p, err := svc.pageRepo.FindByModuleID(mod.NamespaceID, mod.ModuleID); err != nil {
			if err.Error() != repository.ErrPageNotFound.Error() {
				return err
			}
		} else if p.ID > 0 && mod.ID != p.ID {
			return ErrModulePageExists
		}
	}

	return nil
}

func (svc page) DeleteByID(namespaceID, pageID uint64) error {
	if pageID == 0 {
		return ErrInvalidID.withStack()
	}

	if _, err := svc.loadNamespace(namespaceID); err != nil {
		return err
	}

	if p, err := svc.pageRepo.FindByID(namespaceID, pageID); err != nil {
		return err
	} else if !svc.ac.CanDeletePage(svc.ctx, p) {
		return ErrNoDeletePermissions.withStack()
	}

	return svc.pageRepo.DeleteByID(namespaceID, pageID)
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
