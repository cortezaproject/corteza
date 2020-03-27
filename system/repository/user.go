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
		UndeleteByID(id uint64) error

		Metrics() (*types.UserMetrics, error)
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
		From(r.table() + " AS u")
}

func (r *user) With(ctx context.Context, db *factory.DB) UserRepository {
	return &user{
		repository: r.repository.With(ctx, db),
	}
}

// FindByUsername searches for valid (not deleted, not suspended) users by username
func (r user) FindByUsername(username string) (*types.User, error) {
	return r.findOneBy("username", username, true)
}

// FindByHandle searches for valid (not deleted, not suspended) users by handle
func (r user) FindByHandle(handle string) (*types.User, error) {
	return r.findOneBy("handle", handle, true)
}

// FindByEmail searches for valid (not deleted, not suspended) users by email
func (r user) FindByEmail(email string) (*types.User, error) {
	return r.findOneBy("email", email, true)
}

// FindByID searches for users by their ID
func (r user) FindByID(id uint64) (*types.User, error) {
	return r.findOneBy("id", id, false)
}

// findOneBy searches for users by given field & value
func (r user) findOneBy(field string, value interface{}, validOnly bool) (*types.User, error) {
	var (
		u = &types.User{}

		q = r.query().
			Where(squirrel.Eq{field: value})
	)

	if validOnly {
		q = q.Where(squirrel.Eq{
			"deleted_at":   nil,
			"suspended_at": nil,
		})
	}

	if err := rh.FetchOne(r.db(), q, u); err != nil {
		return nil, err
	} else if u.ID == 0 {
		return nil, ErrUserNotFound
	}

	return u, nil
}

// Find searches for users
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

	query = rh.FilterNullByState(query, "u.deleted_at", f.Deleted)
	query = rh.FilterNullByState(query, "u.suspended_at", f.Suspended)

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
		query = query.Where(whereMasked(f.IsEmailUnmaskable, squirrel.Eq{"u.email": f.Email}))
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

	return set, f, rh.FetchPaged(r.db(), query, f.PageFilter, &set)
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
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"suspended_at": time.Now()}, squirrel.Eq{"id": id})
}

func (r *user) UnsuspendByID(id uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"suspended_at": nil}, squirrel.Eq{"id": id})
}

func (r *user) DeleteByID(id uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"deleted_at": time.Now()}, squirrel.Eq{"id": id})
}

func (r *user) UndeleteByID(id uint64) error {
	return rh.UpdateColumns(r.db(), r.table(), rh.Set{"deleted_at": nil}, squirrel.Eq{"id": id})
}

// Metrics collects and returns user metrics
func (r user) Metrics() (rval *types.UserMetrics, err error) {
	var (
		counters = squirrel.
			Select(
				"COUNT(*) as total",
				"SUM(IF(deleted_at IS NULL, 0, 1)) as deleted",
				"SUM(IF(suspended_at IS NULL, 0, 1)) as suspended",
				"SUM(IF(deleted_at IS NULL AND suspended_at IS NULL, 1, 0)) as valid",
			).
			From(r.table() + " AS u")
	)

	rval = &types.UserMetrics{}

	if err = rh.FetchOne(r.db(), counters, rval); err != nil {
		return
	}

	// Fetch daily metrics for created, updated, deleted and suspended users
	err = rh.MultiDailyMetrics(
		r.db(),
		squirrel.Select().From(r.table()+" AS u"),
		[]string{
			"created_at",
			"updated_at",
			"deleted_at",
			"suspended_at",
		},
		&rval.DailyCreated,
		&rval.DailyUpdated,
		&rval.DailyDeleted,
		&rval.DailySuspended,
	)

	if err != nil {
		return
	}

	return
}
