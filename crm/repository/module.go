package repository

import (
	"context"
	"github.com/crusttech/crust/crm/types"
	"github.com/titpetric/factory"
)

type (
	module struct{}
)

func Module() module {
	return module{}
}

func (r module) FindByID(ctx context.Context, id uint64) (*types.Module, error) {
	db := factory.Database.MustGet()
	mod := &types.Module{}
	if err := db.Get(mod, "SELECT * FROM crm_module WHERE id = ?", id); err != nil {
		println(err.Error())
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r module) Find(ctx context.Context) ([]*types.Module, error) {
	db := factory.Database.MustGet()
	mod := make([]*types.Module, 0)
	if err := db.Select(&mod, "SELECT * FROM crm_module ORDER BY name ASC"); err != nil {
		println(err.Error())
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r module) Create(ctx context.Context, mod *types.Module) (*types.Module, error) {
	db := factory.Database.MustGet()
	mod.SetID(factory.Sonyflake.NextID())
	if err := db.Insert("crm_module", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r module) Update(ctx context.Context, mod *types.Module) (*types.Module, error) {
	db := factory.Database.MustGet()
	if err := db.Replace("crm_module", mod); err != nil {
		return nil, ErrDatabaseError
	} else {
		return mod, nil
	}
}

func (r module) DeleteByID(ctx context.Context, id uint64) error {
	db := factory.Database.MustGet()
	if _, err := db.Exec("DELETE FROM crm_module WHERE ID = ?", id); err != nil {
		return ErrDatabaseError
	} else {
		return nil
	}
}

//func (r module) Edit(r *moduleEditRequest) (interface{}, error) {
//	db := factory.Database.MustGet()
//	m := module{}.New()
//	m.SetID(r.id).SetName(r.name)
//	if m.GetID() > 0 {
//		return m, db.Replace("crm_module", m)
//	}
//	m.SetID(factory.Sonyflake.NextID())
//	return m, db.Insert("crm_module", m)
//}
//
//func (r module) ContentList(r *moduleContentListRequest) (interface{}, error) {
//	db := factory.Database.MustGet()
//	if r.id > 0 {
//		m := ModuleContentRow{}.New()
//		return m, db.Get(m, "select * from crm_module id=?", r.id)
//	}
//
//	res := make([]ModuleContentRow, 0)
//	err = db.Select(&res, "select * from crm_module order by name asc")
//	return res, err
//}
//
//func (r module) ContentEdit(r *moduleContentEditRequest) (interface{}, error) {
//	return nil, errors.New("Not implemented: module.content/edit")
//}
//
//func (r module) ContentDelete(r *moduleContentDeleteRequest) (interface{}, error) {
//	return nil, errors.New("Not implemented: module.content/delete")
//}
