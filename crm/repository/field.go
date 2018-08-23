package repository

import (
	"context"
	"github.com/crusttech/crust/crm/types"
)

type (
	Field interface {
		With(ctx context.Context) Field

		FindByType(t string) (*types.Field, error)
		Find() ([]*types.Field, error)
	}

	field struct {
		*repository
	}
)

func NewField(ctx context.Context) Field {
	return (&field{}).With(ctx)
}

func (r *field) With(ctx context.Context) Field {
	return &field{
		repository: r.repository.With(ctx),
	}
}

// FindByName returns field with a given name
func (f *field) FindByType(t string) (*types.Field, error) {
	res := &types.Field{}
	return res, f.db().Get(res, "SELECT * from crm_fields where field_type=?", t)
}

// Find returns all known fields
func (f *field) Find() ([]*types.Field, error) {
	mod := make([]*types.Field, 0)
	return mod, f.db().Select(&mod, "SELECT * FROM crm_fields ORDER BY field_name ASC")
}
