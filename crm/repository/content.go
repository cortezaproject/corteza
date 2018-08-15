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


func (r *content) FindByID(id uint64) (*types.Content, error) {
	mod := &types.Content{}
	if err := r.db().Get(mod, "SELECT * FROM crm_content WHERE id = ?", id); err != nil {
		println(err.Error())
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r *content) Find() ([]*types.Content, error) {
	mod := make([]*types.Content, 0)
	if err := r.db().Select(&mod, "SELECT * FROM crm_content ORDER BY name ASC"); err != nil {
		println(err.Error())
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r *content) Create(mod *types.Content) (*types.Content, error) {
	mod.ID = factory.Sonyflake.NextID()
	return mod, r.db().Insert("crm_content", mod)
}

func (r *content) Update(mod *types.Content) (*types.Content, error) {
	return mod, r.db().Replace("crm_content", mod)

}

func (r *content) DeleteByID(id uint64) error {
	if _, err := r.db().Exec("DELETE FROM crm_content WHERE ID = ?", id); err != nil {
		return ErrDatabaseError
	} else {
		return nil
	}
}
