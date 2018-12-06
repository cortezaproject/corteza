package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/repository"
	"github.com/crusttech/crust/crm/types"

	systemService "github.com/crusttech/crust/system/service"
)

type (
	content struct {
		db  *factory.DB
		ctx context.Context

		repository repository.ContentRepository
		pageRepo   repository.PageRepository

		userSvc systemService.UserService
	}

	ContentService interface {
		With(ctx context.Context) ContentService

		FindByID(contentID uint64) (*types.Content, error)

		Report(moduleID uint64, params *types.ContentReport) (interface{}, error)
		Find(moduleID uint64, query string, page int, perPage int, sort string) (*repository.FindResponse, error)

		Create(content *types.Content) (*types.Content, error)
		Update(content *types.Content) (*types.Content, error)
		DeleteByID(contentID uint64) error

		Fields(mod *types.Content) ([]*types.ContentColumn, error)
	}
)

func Content() ContentService {
	return (&content{
		userSvc: systemService.DefaultUser,
	}).With(context.Background())
}

func (s *content) With(ctx context.Context) ContentService {
	db := repository.DB(ctx)
	return &content{
		db:         db,
		ctx:        ctx,
		repository: repository.Content(ctx, db),
		pageRepo:   repository.Page(ctx, db),
		userSvc:    s.userSvc.With(ctx),
	}
}

func (s *content) FindByID(id uint64) (*types.Content, error) {
	response, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return response, s.preload(response, "page", "user", "fields")
}

func (s *content) Report(moduleID uint64, params *types.ContentReport) (interface{}, error) {
	return s.repository.Report(moduleID, params)
}

func (s *content) Find(moduleID uint64, query string, page int, perPage int, sort string) (*repository.FindResponse, error) {
	response, err := s.repository.Find(moduleID, query, page, perPage, sort)
	if err != nil {
		return nil, err
	}
	if err := s.preloadAll(response.Contents, "user", "fields"); err != nil {
		return nil, err
	}
	return response, nil
}

func (s *content) Create(mod *types.Content) (*types.Content, error) {
	response, err := s.repository.Create(mod)
	if err != nil {
		return nil, err
	}
	return response, s.preload(response, "user", "fields")
}

func (s *content) Update(content *types.Content) (c *types.Content, err error) {
	validate := func() error {
		if content.ID == 0 {
			return errors.New("Error updating content: invalid ID")
		} else if c, err = s.repository.FindByID(content.ID); err != nil {
			return errors.Wrap(err, "Error while loading content for update")
		} else {
			content.CreatedAt = c.CreatedAt
		}

		return nil
	}

	if err = validate(); err != nil {
		return nil, err
	}

	return c, s.db.Transaction(func() (err error) {
		c, err = s.repository.Update(content)
		return
	})
}

func (s *content) Fields(mod *types.Content) ([]*types.ContentColumn, error) {
	return s.repository.Fields(mod)
}

func (s *content) DeleteByID(id uint64) error {
	return s.repository.DeleteByID(id)
}
