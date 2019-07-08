package repository

import (
	"context"
	"io"
	"time"

	"github.com/titpetric/factory"
	"gopkg.in/Masterminds/squirrel.v1"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	UserRepository interface {
		With(ctx context.Context, db *factory.DB) UserRepository

		FindByEmail(email string) (*types.User, error)
		FindByUsername(username string) (*types.User, error)
		FindByID(id uint64) (*types.User, error)
		FindByIDs(id ...uint64) (types.UserSet, error)
		Find(filter types.UserFilter) (set types.UserSet, f types.UserFilter, err error)
		Total() uint

		Create(mod *types.User) (*types.User, error)
		Update(mod *types.User) (*types.User, error)

		BindAvatar(user *types.User, avatar io.Reader) (*types.User, error)

		SuspendByID(id uint64) error
		UnsuspendByID(id uint64) error
		DeleteByID(id uint64) error
	}

	user struct {
		*repository
	}
)

const (
	ErrUserNotFound = repositoryError("UserNotFound")
)

func User(ctx context.Context, db *factory.DB) UserRepository {
	return (&user{}).With(ctx, db)
}

func (r user) table() string {
	return "sys_user"
}

func (r user) columns() []string {
	return []string{
		"u.id",
		"u.email",
		"u.username",
		"u.name",
		"u.handle",
		"u.meta",
		"u.kind",
		"u.rel_organisation",
		"u.email_confirmed",
		"u.created_at",
		"u.updated_at",
		"u.suspended_at",
		"u.deleted_at",
	}
}

func (r user) query() squirrel.SelectBuilder {
	return r.queryNoFilter().Where("u.deleted_at IS NULL AND u.suspended_at IS NULL")
}

func (r user) queryNoFilter() squirrel.SelectBuilder {
	return squirrel.
		Select().
		From(r.table() + " AS u").
		Columns(r.columns()...)
}

func (r *user) With(ctx context.Context, db *factory.DB) UserRepository {
	return &user{
		repository: r.repository.With(ctx, db),
	}
}

func (r user) findBy(field string, value interface{}) (*types.User, error) {
	var (
		query = r.query().Where("u."+field+" = ?", value)
		u     = &types.User{}
	)

	return u, isFound(r.fetchOne(u, query), u.ID > 0, ErrUserNotFound)
}

func (r user) FindByUsername(username string) (*types.User, error) {
	return r.findBy("username", username)
}

func (r user) FindByEmail(email string) (*types.User, error) {
	return r.findBy("email", email)
}

func (r user) FindByID(id uint64) (*types.User, error) {
	return r.findBy("id", id)
}

func (r user) FindByIDs(IDs ...uint64) (types.UserSet, error) {
	if len(IDs) == 0 {
		return nil, nil
	}

	var (
		query = r.query().Where("u.id IN (?)", IDs)
		uu    = types.UserSet{}
	)

	return uu, r.fetchSet(&uu, query)
}

func (r user) Find(filter types.UserFilter) (set types.UserSet, f types.UserFilter, err error) {
	f = filter
	q := r.queryNoFilter()

	if !f.IncDeleted {
		q = q.Where("u.deleted_at IS NULL")
	}

	if !f.IncSuspended {
		q = q.Where("u.suspended_at IS NULL")
	}

	if f.Query != "" {
		qs := f.Query + "%"
		q = q.Where("u.username LIKE ? OR u.email LIKE ? OR u.name LIKE ?", qs, qs, qs)
	}

	if f.Email != "" {
		q = q.Where("u.email = ?", f.Email)
	}

	if f.Username != "" {
		q = q.Where("u.username = ?", f.Username)
	}

	if f.Kind != "" {
		q = q.Where("u.kind = ?", f.Kind)
	}

	if f.Email != "" {
		q = q.Where("u.email = ?", f.Email)
	}

	// @todo add support for more sophisticated sorting through ql
	//       refactor github.com/cortezaproject/corteza-server/compose/internal/repository/ql
	//       for common use (out of compose pkg)
	switch f.Sort {
	case "createdAt":
		q = q.OrderBy("created_at")
	case "updatedAt":
		q = q.OrderBy("updated_at")
	case "deletedAt":
		q = q.OrderBy("deleted_at")
	case "suspendedAt":
		q = q.OrderBy("suspended_at")
	case "email", "username":
		q = q.OrderBy(f.Sort)
	case "userID":
		q = q.OrderBy("id")
	default:
		q = q.OrderBy("id")
	}

	return set, f, r.fetchPaged(&set, q, f.Page, f.PerPage)
}

func (r user) Total() (count uint) {

	count, _ = r.count(r.query())
	return
}

func (r *user) Create(mod *types.User) (*types.User, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	return mod, r.db().Insert(r.table(), mod)
}

func (r *user) Update(mod *types.User) (*types.User, error) {
	mod.UpdatedAt = timeNowPtr()
	return mod, r.db().Replace(r.table(), mod)
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
	return r.updateColumnByID(r.table(), "suspended_at", time.Now(), id)
}

func (r *user) UnsuspendByID(id uint64) error {
	return r.updateColumnByID(r.table(), "suspended_at", nil, id)
}

func (r *user) DeleteByID(id uint64) error {
	return r.updateColumnByID(r.table(), "deleted_at", time.Now(), id)
}
