package repository

import (
	"context"
	"io"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	UserRepository interface {
		With(ctx context.Context, db *factory.DB) UserRepository

		FindByEmail(email string) (*types.User, error)
		FindByUsername(username string) (*types.User, error)
		FindByHandle(handle string) (*types.User, error)
		FindByID(id uint64) (*types.User, error)
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
	return squirrel.
		Select(r.columns()...).
		From(r.table() + " AS u").
		Where(squirrel.And{
			squirrel.Eq{"deleted_at": nil},
			squirrel.Eq{"suspended_at": nil},
		})
}

func (r *user) With(ctx context.Context, db *factory.DB) UserRepository {
	return &user{
		repository: r.repository.With(ctx, db),
	}
}

func (r user) FindByUsername(username string) (*types.User, error) {
	return r.findOneBy("username", username)
}

func (r user) FindByHandle(handle string) (*types.User, error) {
	return r.findOneBy("handle", handle)
}

func (r user) FindByEmail(email string) (*types.User, error) {
	return r.findOneBy("email", email)
}

func (r user) FindByID(id uint64) (*types.User, error) {
	return r.findOneBy("id", id)
}

func (r user) findOneBy(field string, value interface{}) (*types.User, error) {
	var (
		u = &types.User{}

		q = r.query().
			Where(squirrel.Eq{field: value})

		err = rh.FetchOne(r.db(), q, u)
	)

	if err != nil {
		return nil, err
	} else if u.ID == 0 {
		return nil, ErrUserNotFound
	}

	return u, nil
}

func (r user) Find(filter types.UserFilter) (set types.UserSet, f types.UserFilter, err error) {
	f = filter

	if f.Sort == "" {
		f.Sort = "id"
	}

	query := r.query()

	// Returns user filter (flt) wrapped in IF() function with cnd as condition (when cnd != nil)
	whereMasked := func(cnd *permissions.ResourceFilter, flt squirrel.Sqlizer) squirrel.Sqlizer {
		if cnd != nil {
			return rh.SquirrelFunction("IF", cnd, flt, squirrel.Expr("false"))
		} else {
			return flt
		}
	}

	if !f.IncDeleted {
		query = query.Where(squirrel.Eq{"u.deleted_at": nil})
	}

	if !f.IncSuspended {
		query = query.Where(squirrel.Eq{"u.suspended_at": nil})
	}

	if len(f.UserID) > 0 {
		query = query.Where(squirrel.Eq{"u.ID": f.UserID})
	}

	if len(f.RoleID) > 0 {
		or := squirrel.Or{}
		// Due to lack of support for more exotic expressions (slice of values inside subquery)
		// we'll use set of OR expressions as a workaround
		for _, roleID := range f.RoleID {
			or = append(or, squirrel.Expr("u.ID IN (SELECT rel_user FROM sys_role_member WHERE rel_role = ?)", roleID))
		}

		query = query.Where(or)
	}

	if f.Query != "" {
		qs := f.Query + "%"
		query = query.Where(squirrel.Or{
			squirrel.Like{"u.username": qs},
			squirrel.Like{"u.handle": qs},
			whereMasked(f.IsEmailUnmaskable, squirrel.Like{"u.email": qs}),
			whereMasked(f.IsNameUnmaskable, squirrel.Like{"u.name": qs}),
		})
	}

	if f.Email != "" {
		query = query.Where(whereMasked(f.IsNameUnmaskable, squirrel.Eq{"u.name": f.Email}))
	}

	if f.Username != "" {
		query = query.Where(squirrel.Eq{"u.username": f.Username})
	}

	if f.Handle != "" {
		query = query.Where(squirrel.Eq{"u.handle": f.Handle})
	}

	if f.Kind != "" {
		query = query.Where(squirrel.Eq{"u.kind": f.Kind})
	}

	if f.IsReadable != nil {
		query = query.Where(f.IsReadable)
	}

	var orderBy []string
	if orderBy, err = rh.ParseOrder(f.Sort, r.columns()...); err != nil {
		return
	} else {
		query = query.OrderBy(orderBy...)
	}

	if f.Count, err = rh.Count(r.db(), query); err != nil || f.Count == 0 {
		return
	}

	return set, f, rh.FetchPaged(r.db(), query, f.Page, f.PerPage, &set)
}

func (r user) Total() (count uint) {
	count, _ = rh.Count(r.db(), squirrel.Select().From(r.table()))
	return
}

func (r *user) Create(mod *types.User) (*types.User, error) {
	mod.ID = factory.Sonyflake.NextID()
	mod.CreatedAt = time.Now()
	return mod, r.db().Insert(r.table(), mod)
}

func (r *user) Update(mod *types.User) (*types.User, error) {
	rh.SetCurrentTimeRounded(&mod.UpdatedAt)
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
