package crm

import (
	"github.com/pkg/errors"
)

var _ = errors.Wrap

func (*Modules) List(r *modulesListRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Modules.list")
}

func (*Modules) Edit(r *modulesEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Modules.edit")
}

func (*Modules) ContentList(r *modulesContentListRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Modules.content/list")
}

func (*Modules) ContentEdit(r *modulesContentEditRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Modules.content/edit")
}

func (*Modules) ContentDelete(r *modulesContentDeleteRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Modules.content/delete")
}
