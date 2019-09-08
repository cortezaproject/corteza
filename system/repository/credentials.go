package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	CredentialsRepository interface {
		With(ctx context.Context, db *factory.DB) CredentialsRepository

		FindByID(ID uint64) (*types.Credentials, error)
		FindByCredentials(kind, credentials string) (cc types.CredentialsSet, err error)
		FindByKind(ownerID uint64, kind string) (cc types.CredentialsSet, err error)
		FindByOwnerID(ownerID uint64) (cc types.CredentialsSet, err error)
		Find() (cc types.CredentialsSet, err error)

		Create(c *types.Credentials) (*types.Credentials, error)
		Update(c *types.Credentials) (*types.Credentials, error)
		DeleteByID(id uint64) error
		DeleteByKind(ownerID uint64, kind string) (err error)
	}

	credentials struct {
		*repository

		// sql table reference
		tblname string
	}
)

const (
	sqlCredentialsColumns = "id, rel_owner, kind, label, credentials, meta, expires_at, " +
		"created_at, deleted_at"
	sqlCredentialsScope  = "deleted_at IS NULL"
	sqlCredentialsSelect = "SELECT " + sqlCredentialsColumns + " FROM %s WHERE " + sqlCredentialsScope

	ErrCredentialsNotFound = repositoryError("CredentialsNotFound")
)

func Credentials(ctx context.Context, db *factory.DB) CredentialsRepository {
	return (&credentials{}).With(ctx, db)
}

func (r *credentials) With(ctx context.Context, db *factory.DB) CredentialsRepository {
	return &credentials{
		repository: r.repository.With(ctx, db),
		tblname:    "sys_credentials",
	}
}

func (r *credentials) FindByID(ID uint64) (*types.Credentials, error) {
	sql := fmt.Sprintf(sqlCredentialsSelect, r.tblname) + " AND id = ?"
	mod := &types.Credentials{}

	return mod, isFound(r.db().Get(mod, sql, ID), mod.ID > 0, ErrCredentialsNotFound)
}

func (r *credentials) FindByCredentials(kind, credentials string) (cc types.CredentialsSet, err error) {
	return r.fetchSet(
		fmt.Sprintf(sqlCredentialsSelect+" AND kind = ? AND credentials = ?", r.tblname),
		kind,
		credentials)
}

func (r *credentials) FindByKind(ownerID uint64, kind string) (cc types.CredentialsSet, err error) {
	return r.fetchSet(
		fmt.Sprintf(sqlCredentialsSelect+" AND rel_owner = ? AND kind = ?", r.tblname),
		ownerID,
		kind)
}

func (r *credentials) FindByOwnerID(ownerID uint64) (cc types.CredentialsSet, err error) {
	return r.fetchSet(
		fmt.Sprintf(sqlCredentialsSelect+" AND rel_owner = ?", r.tblname),
		ownerID)
}

func (r *credentials) Find() (cc types.CredentialsSet, err error) {
	return r.fetchSet(
		fmt.Sprintf(sqlCredentialsSelect, r.tblname))
}

func (r *credentials) fetchSet(sql string, args ...interface{}) (cc types.CredentialsSet, err error) {
	cc = types.CredentialsSet{}
	return cc, r.db().Select(&cc, sql, args...)
}

func (r *credentials) Create(c *types.Credentials) (*types.Credentials, error) {
	c.ID = factory.Sonyflake.NextID()
	c.CreatedAt = time.Now()
	return c, r.db().Insert(r.tblname, c)
}

func (r *credentials) Update(c *types.Credentials) (*types.Credentials, error) {
	updatedAt := time.Now()
	c.UpdatedAt = &updatedAt
	return c, r.db().Replace(r.tblname, c)
}

func (r *credentials) DeleteByID(id uint64) error {
	return r.updateColumnByID(r.tblname, "deleted_at", time.Now(), id)
}

func (r *credentials) DeleteByKind(ownerID uint64, kind string) (err error) {
	return exec(r.db().Exec(
		fmt.Sprintf("UPDATE %s SET deleted_at = NOW() WHERE rel_owner = ? AND kind = ?", r.tblname),
		ownerID,
		kind))
}
