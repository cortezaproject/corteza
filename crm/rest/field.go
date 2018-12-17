package rest

import (
	"context"

	"github.com/crusttech/crust/crm/rest/request"
	"github.com/crusttech/crust/crm/service"
	"github.com/crusttech/crust/crm/types"
)

type (
	Field struct {
		field service.FieldService
	}

	FieldService interface {
		FindByType(ctx context.Context, t string) (*types.Field, error)
		Find(ctx context.Context) ([]*types.Field, error)
	}
)

func (Field) New() *Field {
	return &Field{
		field: service.DefaultField,
	}
}

func (s *Field) List(ctx context.Context, _ *request.FieldList) (interface{}, error) {
	return s.field.With(ctx).Find()
}

func (s *Field) Type(ctx context.Context, r *request.FieldType) (interface{}, error) {
	return s.field.With(ctx).FindByType(r.TypeID)
}
