package rest

import (
	"github.com/pkg/errors"

	"context"
	"github.com/crusttech/crust/crm/rest/server"
	"github.com/crusttech/crust/crm/types"
)

var _ = errors.Wrap

type (
	Field struct {
		svc fieldService
	}

	fieldService interface {
		FindByName(ctx context.Context, name string) (*types.Field, error)
		Find(ctx context.Context) ([]*types.Field, error)
	}
)

func (Field) New(fieldSvc fieldService) *Field {
	var ctrl = &Field{}
	ctrl.svc = fieldSvc
	return ctrl
}

func (self *Field) List(ctx context.Context, _ *server.FieldListRequest) (interface{}, error) {
	return self.svc.Find(ctx)
}

func (self *Field) Type(ctx context.Context, r *server.FieldTypeRequest) (interface{}, error) {
	return self.svc.FindByName(ctx, r.ID)
}
