package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

const (
	ErrRoleNameNotUnique   = serviceError("RoleNameNotUnique")
	ErrRoleHandleNotUnique = serviceError("RoleHandleNotUnique")
)

type (
	role struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac   roleAccessController
		user UserService

		role repository.RoleRepository
	}

	roleAccessController interface {
		CanAccess(context.Context) bool

		CanCreateRole(context.Context) bool
		CanReadRole(context.Context, *types.Role) bool
		CanUpdateRole(context.Context, *types.Role) bool
		CanDeleteRole(context.Context, *types.Role) bool
		CanManageRoleMembers(context.Context, *types.Role) bool

		FilterReadableRoles(ctx context.Context) *permissions.ResourceFilter
	}

	RoleService interface {
		With(ctx context.Context) RoleService

		FindByID(roleID uint64) (*types.Role, error)
		FindByName(name string) (*types.Role, error)
		FindByHandle(handle string) (*types.Role, error)
		Find(types.RoleFilter) (types.RoleSet, types.RoleFilter, error)

		Create(role *types.Role) (*types.Role, error)
		Update(role *types.Role) (*types.Role, error)
		Merge(roleID, targetroleID uint64) error
		Move(roleID, organisationID uint64) error

		Archive(ID uint64) error
		Unarchive(ID uint64) error
		Delete(ID uint64) error
		Undelete(ID uint64) error

		Membership(userID uint64) ([]*types.RoleMember, error)
		MemberList(roleID uint64) ([]*types.RoleMember, error)
		MemberAdd(roleID, userID uint64) error
		MemberRemove(roleID, userID uint64) error
	}
)

func Role(ctx context.Context) RoleService {
	return (&role{
		ac:     DefaultAccessControl,
		logger: DefaultLogger.Named("role"),
	}).With(ctx)
}

func (svc role) With(ctx context.Context) RoleService {
	db := repository.DB(ctx)
	return &role{
		db:     db,
		ctx:    ctx,
		logger: svc.logger,
		ac:     svc.ac,

		role: repository.Role(ctx, db),
		user: DefaultUser.With(ctx),
	}
}

func (svc role) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc role) FindByID(roleID uint64) (*types.Role, error) {
	return svc.findByID(roleID)
}

func (svc role) findByID(roleID uint64) (*types.Role, error) {
	if roleID == 0 {
		return nil, ErrInvalidID
	}

	role, err := svc.role.FindByID(roleID)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadRole(svc.ctx, role) {
		return nil, ErrNoPermissions.withStack()
	}
	return role, nil
}

func (svc role) Find(f types.RoleFilter) (types.RoleSet, types.RoleFilter, error) {
	f.IsReadable = svc.ac.FilterReadableRoles(svc.ctx)

	if f.Deleted > 0 {
		// If list with deleted or suspended users is requested
		// user must have access permissions to system (ie: is admin)
		//
		// not the best solution but ATM it allows us to have at least
		// some kind of control over who can see deleted or archived roles
		if !svc.ac.CanAccess(svc.ctx) {
			return nil, f, ErrNoPermissions.withStack()
		}
	}

	return svc.role.Find(f)
}

func (svc role) FindByName(rolename string) (*types.Role, error) {
	return svc.role.FindByName(rolename)
}

func (svc role) FindByHandle(handle string) (*types.Role, error) {
	return svc.role.FindByHandle(handle)
}

func (svc role) Create(new *types.Role) (r *types.Role, err error) {

	if !handle.IsValid(new.Handle) {
		return nil, ErrInvalidHandle
	}

	if !svc.ac.CanCreateRole(svc.ctx) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	return r, svc.db.Transaction(func() (err error) {
		if err = eventbus.WaitFor(svc.ctx, event.RoleBeforeCreate(new, r)); err != nil {
			return
		}

		if err = svc.UniqueCheck(new); err != nil {
			return
		}

		if r, err = svc.role.Create(new); err != nil {
			return
		}

		defer eventbus.Dispatch(svc.ctx, event.RoleAfterCreate(new, r))
		return
	})
}

func (svc role) Update(upd *types.Role) (r *types.Role, err error) {
	if upd.ID == 0 {
		return nil, ErrInvalidID
	}

	if !handle.IsValid(upd.Handle) {
		return nil, ErrInvalidHandle
	}

	if !svc.ac.CanUpdateRole(svc.ctx, upd) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	return r, svc.db.Transaction(func() (err error) {
		if r, err = svc.role.FindByID(upd.ID); err != nil {
			return
		}

		if err = eventbus.WaitFor(svc.ctx, event.RoleBeforeUpdate(upd, r)); err != nil {
			return
		}

		if err = svc.UniqueCheck(upd); err != nil {
			return
		}

		// Assign changed values
		r.Name = upd.Name
		r.Handle = upd.Handle

		if r, err = svc.role.Update(r); err != nil {
			return err
		}

		defer eventbus.Dispatch(svc.ctx, event.RoleAfterUpdate(upd, r))

		return nil
	})
}

func (svc role) UniqueCheck(r *types.Role) (err error) {
	if r.Handle != "" {
		if ex, _ := svc.role.FindByHandle(r.Handle); ex != nil && ex.ID > 0 && ex.ID != r.ID {
			return ErrRoleHandleNotUnique
		}
	}

	if r.Name != "" {
		if ex, _ := svc.role.FindByName(r.Name); ex != nil && ex.ID > 0 && ex.ID != r.ID {
			return ErrRoleNameNotUnique
		}
	}

	return nil
}

func (svc role) Delete(roleID uint64) (err error) {
	var (
		role *types.Role
	)

	if role, err = svc.findByID(roleID); err != nil {
		return err
	}

	if !svc.ac.CanDeleteRole(svc.ctx, role) {
		return ErrNoPermissions.withStack()
	}

	if err = eventbus.WaitFor(svc.ctx, event.RoleBeforeDelete(nil, role)); err != nil {
		return
	}

	if err = svc.role.DeleteByID(roleID); err != nil {
		return
	}

	defer eventbus.Dispatch(svc.ctx, event.RoleAfterDelete(nil, role))
	return
}

func (svc role) Undelete(roleID uint64) error {
	role, err := svc.findByID(roleID)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteRole(svc.ctx, role) {
		return ErrNoPermissions.withStack()
	}

	return svc.role.UndeleteByID(roleID)
}

func (svc role) Archive(roleID uint64) error {
	role, err := svc.findByID(roleID)
	if err != nil {
		return err
	}

	if !svc.ac.CanUpdateRole(svc.ctx, role) {
		return ErrNoPermissions.withStack()
	}

	return svc.role.ArchiveByID(roleID)
}

func (svc role) Unarchive(roleID uint64) error {
	role, err := svc.findByID(roleID)
	if err != nil {
		return err
	}

	if !svc.ac.CanUpdateRole(svc.ctx, role) {
		return ErrNoPermissions.withStack()
	}

	return svc.role.UnarchiveByID(roleID)
}

func (svc role) Merge(roleID, targetRoleID uint64) error {
	role, err := svc.findByID(roleID)
	if err != nil {
		return err
	}

	if targetRoleID == 0 {
		return ErrInvalidID
	}

	if !svc.ac.CanUpdateRole(svc.ctx, role) {
		return ErrNoPermissions.withStack()
	}

	return svc.role.MergeByID(roleID, targetRoleID)
}

func (svc role) Move(roleID, targetOrganisationID uint64) error {
	role, err := svc.findByID(roleID)
	if err != nil {
		return err
	}

	if targetOrganisationID == 0 {
		return ErrInvalidID
	}

	if !svc.ac.CanUpdateRole(svc.ctx, role) {
		return ErrNoPermissions.withStack()
	}

	return svc.role.MoveByID(roleID, targetOrganisationID)
}

func (svc role) Membership(userID uint64) ([]*types.RoleMember, error) {
	return svc.role.MembershipsFindByUserID(userID)
}

func (svc role) MemberList(roleID uint64) ([]*types.RoleMember, error) {
	_, err := svc.findByID(roleID)
	if err != nil {
		return nil, err
	}

	return svc.role.MemberFindByRoleID(roleID)
}

func (svc role) MemberAdd(roleID, userID uint64) (err error) {
	var (
		role *types.Role
		user *types.User
	)

	if role, err = svc.findByID(roleID); err != nil {
		return
	}

	if user, err = svc.user.FindByID(userID); err != nil {
		return
	}

	if err = eventbus.WaitFor(svc.ctx, event.RoleMemberBeforeAdd(user, role)); err != nil {
		return
	}

	if !svc.ac.CanManageRoleMembers(svc.ctx, role) {
		return errors.New("Not allowed to manage role members")
	}

	if err = svc.role.MemberAddByID(role.ID, user.ID); err != nil {
		return
	}

	defer eventbus.Dispatch(svc.ctx, event.RoleMemberAfterAdd(user, role))
	return nil
}

func (svc role) MemberRemove(roleID, userID uint64) (err error) {
	var (
		role *types.Role
		user *types.User
	)

	if role, err = svc.findByID(roleID); err != nil {
		return
	}

	if user, err = svc.user.FindByID(userID); err != nil {
		return
	}

	if err = eventbus.WaitFor(svc.ctx, event.RoleMemberBeforeRemove(user, role)); err != nil {
		return
	}

	if !svc.ac.CanManageRoleMembers(svc.ctx, role) {
		return errors.New("Not allowed to manage role members")
	}

	if err = svc.role.MemberRemoveByID(role.ID, user.ID); err != nil {
		return
	}

	defer eventbus.Dispatch(svc.ctx, event.RoleMemberAfterRemove(user, role))
	return nil
}

var _ RoleService = &role{}
