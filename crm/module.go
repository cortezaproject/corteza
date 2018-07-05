package crm

import (
	"github.com/pkg/errors"
)

var _ = errors.Wrap

func (*Module) List(r *moduleListRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.list")
}

func (*Module) Edit(r *moduleEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.edit")
}

func (*Module) ContentList(r *moduleContentListRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/list")
}

func (*Module) ContentEdit(r *moduleContentEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/edit")
}

func (*Module) ContentDelete(r *moduleContentDeleteRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Module.content/delete")
}
