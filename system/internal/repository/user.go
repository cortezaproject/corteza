package repository

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/types"
)

type (
	UserRepository interface {
		With(ctx context.Context, db *factory.DB) UserRepository

		FindByEmail(email string) (*types.User, error)
		FindByUsername(username string) (*types.User, error)
		FindByID(id uint64) (*types.User, error)
		FindByIDs(id ...uint64) (types.UserSet, error)
		FindBySatosaID(id string) (*types.User, error)
		Find(filter *types.UserFilter) ([]*types.User, error)

		Create(mod *types.User) (*types.User, error)
		Update(mod *types.User) (*types.User, error)

		BindAvatar(user *types.User, avatar io.Reader) (*types.User, error)

		SuspendByID(id uint64) error
		UnsuspendByID(id uint64) error
		DeleteByID(id uint64) error
	}

	user struct {
		*repository

		// sql table reference
		users string
	}
)

const (
	sqlUserColumns = "id, email, username, password, name, handle, " +
		"meta, satosa_id, rel_organisation, " +
		"created_at, updated_at, suspended_at, deleted_at"
	sqlUserScope  = "suspended_at IS NULL AND deleted_at IS NULL"
	sqlUserSelect = "SELECT " + sqlUserColumns + " FROM %s WHERE " + sqlUserScope

	ErrUserNotFound = repositoryError("UserNotFound")
)

func User(ctx context.Context, db *factory.DB) UserRepository {
	return (&user{}).With(ctx, db)
}

func (r *user) With(ctx context.Context, db *factory.DB) UserRepository {
	return &user{
		repository: r.repository.With(ctx, db),
		users:      "sys_user",
	}
}

func (r *user) FindByUsername(username string) (*types.User, error) {
	sql := fmt.Sprintf(sqlUserSelect, r.users) + " AND username = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, username), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindBySatosaID(satosaID string) (*types.User, error) {
	sql := fmt.Sprintf(sqlUserSelect, r.users) + " AND satosa_id = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, satosaID), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindByEmail(email string) (*types.User, error) {
	sql := fmt.Sprintf(sqlUserSelect, r.users) + " AND email = ?"
	mod := &types.User{}

	return mod, isFound(r.db().Get(mod, sql, email), mod.ID > 0, ErrUserNotFound)
}

func (r *user) FindByID(id uint64) (*types.User, error) {
	sql := fmt.Sprintf(sqlUserSelect, r.users) + " AND id = ?"
	mod := &types.User{}
	if err := isFound(r.db().Get(mod, sql, id), mod.ID > 0, ErrUserNotFound); err != nil {
		return nil, err
	}

	err := r.prepare(mod, "roles")
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func (r *user) FindByIDs(IDs ...uint64) (uu types.UserSet, err error) {
	if len(IDs) == 0 {
		return
	}

	sql := fmt.Sprintf(sqlUserSelect, r.users) + " AND id IN (?)"

	if sql, args, err := sqlx.In(sql, IDs); err != nil {
		return nil, err
	} else {
		return uu, r.db().Select(&uu, sql, args...)
	}

}

func (r *user) Find(filter *types.UserFilter) ([]*types.User, error) {
	if filter == nil {
		filter = &types.UserFilter{}
	}

	rval := make([]*types.User, 0)
	params := make([]interface{}, 0)
	sql := fmt.Sprintf(sqlUserSelect, r.users)

	if filter.Query != "" {
		sql += " AND (username LIKE ?"
		params = append(params, filter.Query+"%")
		sql += " OR email LIKE ?"
		params = append(params, filter.Query+"%")
		sql += " OR name LIKE ?)"
		params = append(params, filter.Query+"%")
	}

	if filter.Email != "" {
		sql += " AND (email = ?)"
		params = append(params, filter.Email)
	}

	if filter.Username != "" {
		sql += " AND (username = ?)"
		params = append(params, filter.Username)
	}

	switch filter.OrderBy {
	case "updated_at", "createdAt":
		sql += " ORDER BY updated_at DESC"
	default:
		sql += " ORDER BY username ASC"
	}

	if err := r.db().Select(&rval, sql, params...); err != nil {
		return nil, err
	}
	if err := r.prepareAll(rval, "roles"); err != nil {
		return nil, err
	}

	return rval, nil
}

func (r *user) Create(mod *types.User) (*types.User, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	return mod, r.db().Insert(r.users, mod)
}

func (r *user) Update(mod *types.User) (*types.User, error) {
	mod.UpdatedAt = timeNowPtr()
	return mod, r.db().Replace(r.users, mod)
}

func (r *user) BindAvatar(user *types.User, avatar io.Reader) (*types.User, error) {
	if user.Meta == nil {
		user.Meta = new(types.UserMeta)
	}
	// @todo: IMPORTANT: implement avatar uploading
	user.Meta.Avatar = ""
	return user, nil
}

func (r *user) SuspendByID(id uint64) error {
	return r.updateColumnByID(r.users, "suspend_at", time.Now(), id)
}

func (r *user) UnsuspendByID(id uint64) error {
	return r.updateColumnByID(r.users, "suspend_at", nil, id)
}

func (r *user) DeleteByID(id uint64) error {
	return r.updateColumnByID(r.users, "deleted_at", time.Now(), id)
}

func (r *user) prepareAll(users []*types.User, fields ...string) error {
	for _, user := range users {
		if err := r.prepare(user, fields...); err != nil {
			return err
		}
	}
	return nil
}

func (r *user) prepare(user *types.User, fields ...string) (err error) {
	api := Role(r.Context(), r.db())
	for _, field := range fields {
		switch field {
		case "roles":
			if user.ID > 0 {
				roles, err := api.FindByMemberID(user.ID)
				if err != nil {
					return err
				}
				user.Roles = roles
			}
		default:
		}
	}
	return
}
