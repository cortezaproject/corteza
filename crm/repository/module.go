package repository

import (
	"context"
	"github.com/crusttech/crust/crm/types"
	"github.com/titpetric/factory"
)

type (
	Module interface {
		With(ctx context.Context) Module

		FindByID(id uint64) (*types.Module, error)
		Find() ([]*types.Module, error)
		Create(mod *types.Module) (*types.Module, error)
		Update(mod *types.Module) (*types.Module, error)
		DeleteByID(id uint64) error
	}

	module struct {
		*repository
	}
)

func NewModule(ctx context.Context) Module {
	return &module{
		repository: &repository{
			ctx: ctx,
		},
	}
}

func (r *module) With(ctx context.Context) Module {
	return &module{
		repository: r.repository.With(ctx),
	}
}


func (r *module) FindByID(id uint64) (*types.Module, error) {
	mod := &types.Module{}
	if err := r.db().Get(mod, "SELECT * FROM crm_module WHERE id = ?", id); err != nil {
		println(err.Error())
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r *module) Find() ([]*types.Module, error) {
	mod := make([]*types.Module, 0)
	if err := r.db().Select(&mod, "SELECT * FROM crm_module ORDER BY name ASC"); err != nil {
		println(err.Error())
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r *module) Create(mod *types.Module) (*types.Module, error) {
	mod.ID = factory.Sonyflake.NextID()
	return mod, r.db().Insert("crm_module", mod)
}

func (r *module) Update(mod *types.Module) (*types.Module, error) {
	return mod, r.db().Replace("crm_module", mod)

}

func (r *module) DeleteByID(id uint64) error {
	if _, err := r.db().Exec("DELETE FROM crm_module WHERE ID = ?", id); err != nil {
		return ErrDatabaseError
	} else {
		return nil
	}
}

//func (r *module) Edit(r *moduleEditRequest) (interface{}, error) {
//	db := factory.Database.MustGet()
//	m := module{}.New()
//	m.SetID(r.id).SetName(r.name)
//	if m.GetID() > 0 {
//		return m, db.With(ctx).Replace("crm_module", m)
//	}
//	m.SetID(factory.Sonyflake.NextID())
//	return m, db.With(ctx).Insert("crm_module", m)
//}
//
//func (r *module) ContentList(r *moduleContentListRequest) (interface{}, error) {
//	db := factory.Database.MustGet()
//	if r.id > 0 {
//		m := ModuleContentRow{}.New()
//		return m, db.Get(m, "select * from crm_module id=?", r.id)
//	}
//
//	res := make([]ModuleContentRow, 0)
//	err = db.With(ctx).Select(&res, "select * from crm_module order by name asc")
//	return res, err
//}
//
//func (r *module) ContentEdit(r *moduleContentEditRequest) (interface{}, error) {
//	return nil, errors.New("Not implemented: module.content/edit")
//}
//
//func (r *module) ContentDelete(r *moduleContentDeleteRequest) (interface{}, error) {
//	return nil, errors.New("Not implemented: module.content/delete")
//}
