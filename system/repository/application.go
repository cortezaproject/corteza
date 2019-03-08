package repository

import (
	"context"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/types"
)

type (
	ApplicationRepository interface {
		With(ctx context.Context, db *factory.DB) ApplicationRepository

		FindByID(id uint64) (*types.Application, error)
		Find() (types.ApplicationSet, error)

		Create(mod *types.Application) (*types.Application, error)
		Update(mod *types.Application) (*types.Application, error)

		DeleteByID(id uint64) error
	}

	application struct {
		*repository

		// sql table reference
		applications string
		members      string
	}
)

const (
	sqlApplicationColumns = "id, rel_owner, name, enabled, unify, created_at, updated_at, deleted_at"
	sqlApplicationScope   = "deleted_at IS NULL"

	ErrApplicationNotFound = repositoryError("ApplicationNotFound")
)

func Application(ctx context.Context, db *factory.DB) ApplicationRepository {
	return (&application{}).With(ctx, db)
}

func (r *application) With(ctx context.Context, db *factory.DB) ApplicationRepository {
	return &application{
		repository:   r.repository.With(ctx, db),
		applications: "sys_application",
	}
}

func (r *application) FindByID(id uint64) (*types.Application, error) {
	sql := "SELECT " + sqlApplicationColumns + " FROM " + r.applications + " WHERE id = ? AND " + sqlApplicationScope
	mod := &types.Application{}

	return mod, isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrApplicationNotFound)
}

func (r *application) Find() (types.ApplicationSet, error) {
	rval := make([]*types.Application, 0)
	params := make([]interface{}, 0)

	sql := "SELECT " + sqlApplicationColumns + " FROM " + r.applications + " WHERE " + sqlApplicationScope

	sql += " ORDER BY id ASC"

	return rval, r.db().Select(&rval, sql, params...)
}

func (r *application) Create(mod *types.Application) (*types.Application, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()

	return mod, r.db().Insert(r.applications, mod)
}

func (r *application) Update(mod *types.Application) (*types.Application, error) {
	mod.UpdatedAt = timeNowPtr()

	return mod, r.db().Replace(r.applications, mod)
}

func (r *application) DeleteByID(id uint64) error {
	return r.updateColumnByID(r.applications, "deleted_at", time.Now(), id)
}
