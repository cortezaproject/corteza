package crm

import (
	"github.com/pkg/errors"

	"github.com/crusttech/crust/crm/rest"
	"github.com/crusttech/crust/crm/service"
)

var _ = errors.Wrap

type (
	Field struct {
		service FieldInterface
	}

	FieldInterface interface {
		List() (interface{}, error)
		Type(id string) (interface{}, error)
	}
)

func (Field) New() *Field {
	return &Field{service.Field{}.New()}
}

func (self *Field) List(_ *rest.FieldListRequest) (interface{}, error) {
	return self.service.List()
}

func (self *Field) Type(r *rest.FieldTypeRequest) (interface{}, error) {
	return self.service.Type(r.ID)
}
