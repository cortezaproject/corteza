package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/pkg/errors"
	"github.com/crusttech/human/server/pkg/handle"
	"github.com/crusttech/human/server/pkg/label"
	"github.com/crusttech/human/server/pkg/scope"
	"github.com/crusttech/human/server/store"
	types "github.com/crusttech/human/server/system/types"
	"mime/multipart"
)

type user struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        userAccessController
	services  *userServices
}

func (svc *user) FindByID(ctx context.Context, ID uint64) (res *types.User, err error) {
	var (
		aProps = &userActionProps{user: &types.User{ID: ID}}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		res, err = svc.onLookup(ctx, ID, aProps)
		return err
	}()

	return res, svc.recordAction(ctx, aProps, UserActionLookup, err)
}

func (svc *user) Search(ctx context.Context, filter types.UserFilter) (set types.UserSet, f types.UserFilter, err error) {
	var (
		aProps = &userActionProps{filter: &filter}
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}
		if !svc.ac.CanSearchUsers(ctx) {
			return UserErrNotAllowedToSearch()
		}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
	}()

	return set, f, svc.recordAction(ctx, aProps, UserActionSearch, err)
}

func (svc *user) Create(ctx context.Context, new *types.User) (res *types.User, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &userActionProps{user: new, new: new}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}
		if !svc.ac.CanCreateUser(ctx) {
			return UserErrNotAllowedToCreate()
		}
		res = new
		return svc.onCreate(ctx, new)
	}()

	return res, svc.recordAction(ctx, aProps, UserActionCreate, err)
}

func (svc *user) Update(ctx context.Context, upd *types.User) (res *types.User, err error) {
	var (
		aProps = &userActionProps{update: upd}
		old    *types.User
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadUser(ctx, s, upd.ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setUser(res)
		aProps.setUpdate(res)
		old = res.Clone()

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return UserErrInvalidHandle()
		}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return UserErrStaleData()
		}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
		res.Email = upd.Email
		res.EmailConfirmed = upd.EmailConfirmed
		res.UserGroupID = upd.UserGroupID
		res.Username = upd.Username
		res.Name = upd.Name
		res.Handle = upd.Handle
		res.UpdatedAt = now()

		if err = store.UpdateUser(ctx, s, res); err != nil {
			return err
		}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}

		return nil
	})

	return res, svc.recordAction(ctx, aProps, UserActionUpdate, err, old, res)
}

func (svc *user) DeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &userActionProps{}
		res    *types.User
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadUser(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setUser(res)

		return svc.onDelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, UserActionDelete, err)
}

func (svc *user) UndeleteByID(ctx context.Context, ID uint64) (err error) {
	var (
		aProps = &userActionProps{}
		res    *types.User
	)
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = loadUser(ctx, s, ID); err != nil {
			return
		}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}

		aProps.setUser(res)

		return svc.onUndelete(ctx, s, res, aProps)
	})

	return svc.recordAction(ctx, aProps, UserActionUndelete, err)
}

func loadUser(ctx context.Context, s store.Users, ID uint64) (res *types.User, err error) {
	if ID == 0 {
		return nil, UserErrInvalidID()
	}

	if res, err = store.LookupUserByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, UserErrNotFound()
	}

	return
}

// toLabeledUsers converts to []label.LabeledResource
func toLabeledUsers(set []*types.User) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}

func (svc *user) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

func (svc *user) FindByEmail(ctx context.Context, email string) (u *types.User, err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		u, err = svc.onFindByEmail(ctx, aProps, email)
		return err
	}()

	return u, svc.recordAction(ctx, aProps, UserActionLookup, err)
}

func (svc *user) FindByHandle(ctx context.Context, handle string) (u *types.User, err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

		u, err = svc.onFindByHandle(ctx, aProps, handle)
		return err
	}()

	return u, svc.recordAction(ctx, aProps, UserActionLookup, err)
}

func (svc *user) ToggleEmailConfirmation(ctx context.Context, userID uint64, confirmed bool) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onToggleEmailConfirmation(ctx, aProps, userID, confirmed)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionUpdate, err)
}

func (svc *user) Suspend(ctx context.Context, userID uint64) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onSuspend(ctx, aProps, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionSuspend, err)
}

func (svc *user) Unsuspend(ctx context.Context, userID uint64) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onUnsuspend(ctx, aProps, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionUnsuspend, err)
}

func (svc *user) SetPassword(ctx context.Context, userID uint64, newPassword string) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onSetPassword(ctx, aProps, userID, newPassword)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionSetPassword, err)
}

func (svc *user) DeleteAuthTokensByUserID(ctx context.Context, userID uint64) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onDeleteAuthTokensByUserID(ctx, aProps, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionDeleteAuthTokens, err)
}

func (svc *user) DeleteAuthSessionsByUserID(ctx context.Context, userID uint64) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onDeleteAuthSessionsByUserID(ctx, aProps, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionDeleteAuthSessions, err)
}

func (svc *user) CreateSynthetic(ctx context.Context, src synteticUserDataGen, total uint) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if !svc.ac.CanCreateUser(ctx) {
			return UserErrNotAllowedToCreate()
		}

		err = svc.onCreateSynthetic(ctx, aProps, src, total)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionCreateSynthetic, err)
}

func (svc *user) RemoveSynthetic(ctx context.Context) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if !svc.ac.CanCreateUser(ctx) {
			return UserErrNotAllowedToCreate()
		}

		err = svc.onRemoveSynthetic(ctx, aProps)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionRemoveSynthetic, err)
}

func (svc *user) UploadAvatar(ctx context.Context, userID uint64, upload *multipart.FileHeader) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onUploadAvatar(ctx, aProps, userID, upload)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionUploadAvatar, err)
}

func (svc *user) DeleteAvatar(ctx context.Context, userID uint64) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onDeleteAvatar(ctx, aProps, userID)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionDeleteAvatar, err)
}

func (svc *user) GenerateAvatar(ctx context.Context, userID uint64, bgColor string, initialColor string) (err error) {
	var (
		aProps = &userActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		err = svc.onGenerateAvatar(ctx, aProps, userID, bgColor, initialColor)
		return err
	}()

	return svc.recordAction(ctx, aProps, UserActionGenerateAvatar, err)
}
