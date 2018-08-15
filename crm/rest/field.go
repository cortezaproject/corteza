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
		field service.FieldService
	}

	FieldService interface {
		FindByName(ctx context.Context, name string) (*types.Field, error)
		Find(ctx context.Context) ([]*types.Field, error)
	}
)

func (Field) New() server.FieldAPI {
	return &Field{
		field: service.Field(),
	}
}

func (s *Field) List(ctx context.Context, _ *server.FieldListRequest) (interface{}, error) {
	return s.field.With(ctx).Find()
}

func (s *Field) Type(ctx context.Context, r *server.FieldTypeRequest) (interface{}, error) {
	return s.field.With(ctx).FindByName(r.ID)
}
