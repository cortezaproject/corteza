package crm

import (
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

var _ = errors.Wrap

func (*Module) List(r *moduleListRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	if r.id > 0 {
		m := Module{}.new()
		return m, db.Get(m, "select * from crm_module id=?", r.id)
	}

	res := make([]Module, 0)
	err = db.Select(&res, "select * from crm_module order by name asc")
	return res, err
}

func (*Module) Edit(r *moduleEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	m := Module{}.new()
	m.SetID(r.id).SetName(r.name)
	if m.GetID() > 0 {
		return m, db.Replace("crm_module", m)
	}
	m.SetID(factory.Sonyflake.NextID())
	return m, db.Insert("crm_module", m)
}

func (*Module) ContentList(r *moduleContentListRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	if r.id > 0 {
		m := ModuleContentRow{}.new()
		return m, db.Get(m, "select * from crm_module id=?", r.id)
	}

	res := make([]ModuleContentRow, 0)
	err = db.Select(&res, "select * from crm_module order by name asc")
	return res, err
}

func (*Module) ContentEdit(r *moduleContentEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentDelete(r *moduleContentDeleteRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/delete")
}
