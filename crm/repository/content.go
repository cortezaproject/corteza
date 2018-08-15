package repository

import (
	"context"
	"github.com/crusttech/crust/crm/types"
	"github.com/titpetric/factory"
)

type (
	Content interface {
		With(ctx context.Context) Content

		FindByID(id uint64) (*types.Content, error)
		Find() ([]*types.Content, error)
		Create(mod *types.Content) (*types.Content, error)
		Update(mod *types.Content) (*types.Content, error)
		DeleteByID(id uint64) error
	}

	content struct {
		*repository
	}
)

func NewContent(ctx context.Context) Content {
	return &content{
		repository: &repository{
			ctx: ctx,
		},
	}
}

func (r *content) With(ctx context.Context) Content {
	return &content{
		repository: r.repository.With(ctx),
	}
}

// @todo: update to accepted DeletedAt column semantics from SAM

func (r *content) FindByID(id uint64) (*types.Content, error) {
	mod := &types.Content{}
	return mod, r.db().Get(mod, "SELECT * FROM crm_module_content WHERE id=?", id)
}

func (r *content) Find() ([]*types.Content, error) {
	mod := make([]*types.Content, 0)
	return mod, r.db().Select(&mod, "SELECT * FROM crm_module_content ORDER BY id DESC")
}

func (r *content) Create(mod *types.Content) (*types.Content, error) {
	mod.ID = factory.Sonyflake.NextID()
	return mod, r.db().Insert("crm_module_content", mod)
}

func (r *content) Update(mod *types.Content) (*types.Content, error) {
	return mod, r.db().Replace("crm_module_content", mod)

}

func (r *content) DeleteByID(id uint64) error {
	_, err := r.db().Exec("DELETE FROM crm_module_content WHERE id=?", id)
	return err
}
