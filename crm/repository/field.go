package repository

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/crm/types"
)

type (
	FieldRepository interface {
		With(ctx context.Context, db *factory.DB) FieldRepository

		FindByType(t string) (*types.Field, error)
		Find() ([]*types.Field, error)
	}

	field struct {
		*repository
	}
)

func Field(ctx context.Context, db *factory.DB) FieldRepository {
	return (&field{}).With(ctx, db)
}

func (r *field) With(ctx context.Context, db *factory.DB) FieldRepository {
	return &field{
		repository: r.repository.With(ctx, db),
	}
}

// FindByName returns field with a given name
func (f *field) FindByType(t string) (*types.Field, error) {
	res := &types.Field{}
	return res, f.db().Get(res, "SELECT * from crm_field where field_type=?", t)
}

// Find returns all known fields
func (f *field) Find() ([]*types.Field, error) {
	mod := make([]*types.Field, 0)
	return mod, f.db().Select(&mod, "SELECT * FROM crm_field ORDER BY field_name ASC")
}
