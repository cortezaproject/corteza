package rest

import (
	"github.com/pkg/errors"

	"context"
	"github.com/crusttech/crust/crm/rest/server"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/crm/types"
)

var _ = errors.Wrap

type (
	Field struct {
		service fieldService
	}

	fieldService interface {
		FindByName(context.Context, string) (*types.Field, error)
		Find(context.Context) ([]*types.Field, error)
	}
)

func (Field) New() *Field {
	return &Field{service: service.Field()}
}

func (self *Field) List(_ *server.FieldListRequest) (interface{}, error) {
	return self.service.Find(context.TODO())
}

func (self *Field) Type(r *server.FieldTypeRequest) (interface{}, error) {
	return self.service.FindByName(context.TODO(), r.ID)
}
