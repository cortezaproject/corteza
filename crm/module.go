package crm

import (
	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/rest"
	"github.com/crusttech/crust/crm/types"
)

var _ = errors.Wrap

type Module struct{}

func (Module) New() *Module {
	return &Module{}
}

func (*Module) List(r *rest.ModuleListRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	if r.ID > 0 {
		m := types.Module{}.New()
		return m, db.Get(m, "select * from crm_module id=?", r.ID)
	}

	res := make([]Module, 0)
	err = db.Select(&res, "select * from crm_module order by name asc")
	return res, err
}

func (*Module) Edit(r *rest.ModuleEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	m := types.Module{}.New()
	m.SetID(r.ID).SetName(r.Name)
	if m.GetID() > 0 {
		return m, db.Replace("crm_module", m)
	}
	m.SetID(factory.Sonyflake.NextID())
	return m, db.Insert("crm_module", m)
}

func (*Module) ContentList(r *rest.ModuleContentListRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	if r.ID > 0 {
		m := types.ModuleContentRow{}.New()
		return m, db.Get(m, "select * from crm_module id=?", r.ID)
	}

	res := make([]types.ModuleContentRow, 0)
	err = db.Select(&res, "select * from crm_module order by name asc")
	return res, err
}

func (*Module) ContentEdit(r *rest.ModuleContentEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentDelete(r *rest.ModuleContentDeleteRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/delete")
}
