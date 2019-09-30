package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/repository"
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

		ac roleAccessController

		role repository.RoleRepository
	}

	roleAccessController interface {
		CanCreateRole(context.Context) bool
		CanReadRole(context.Context, *types.Role) bool
		CanUpdateRole(context.Context, *types.Role) bool
		CanDeleteRole(context.Context, *types.Role) bool
		CanManageRoleMembers(context.Context, *types.Role) bool
	}

	RoleService interface {
		With(ctx context.Context) RoleService

		FindByID(roleID uint64) (*types.Role, error)
		FindByName(name string) (*types.Role, error)
		FindByHandle(handle string) (*types.Role, error)
		Find(filter *types.RoleFilter) ([]*types.Role, error)

		Create(role *types.Role) (*types.Role, error)
		Update(role *types.Role) (*types.Role, error)
		Merge(roleID, targetroleID uint64) error
		Move(roleID, organisationID uint64) error

		Archive(ID uint64) error
		Unarchive(ID uint64) error
		Delete(ID uint64) error

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

func (svc role) Find(filter *types.RoleFilter) ([]*types.Role, error) {
	roles, err := svc.role.Find(filter)
	if err != nil {
		return nil, err
	}

	ret := []*types.Role{}
	for _, role := range roles {
		if svc.ac.CanReadRole(svc.ctx, role) {
			ret = append(ret, role)
		}
	}
	return ret, nil
}

func (svc role) FindByName(rolename string) (*types.Role, error) {
	return svc.role.FindByName(rolename)
}

func (svc role) FindByHandle(handle string) (*types.Role, error) {
	return svc.role.FindByHandle(handle)
}

func (svc role) Create(mod *types.Role) (t *types.Role, err error) {

	if !handle.IsValid(mod.Handle) {
		return nil, ErrInvalidHandle
	}

	if !svc.ac.CanCreateRole(svc.ctx) {
		return nil, ErrNoCreatePermissions.withStack()
	}

	return t, svc.db.Transaction(func() (err error) {
		if err = svc.UniqueCheck(mod); err != nil {
			return
		}

		t, err = svc.role.Create(mod)
		return
	})
}

func (svc role) Update(mod *types.Role) (t *types.Role, err error) {
	if mod.ID == 0 {
		return nil, ErrInvalidID
	}

	if !handle.IsValid(mod.Handle) {
		return nil, ErrInvalidHandle
	}

	if !svc.ac.CanUpdateRole(svc.ctx, mod) {
		return nil, ErrNoUpdatePermissions.withStack()
	}

	// @todo: make sure archived & deleted entries can not be edited

	return t, svc.db.Transaction(func() (err error) {
		if t, err = svc.role.FindByID(mod.ID); err != nil {
			return
		}

		if err = svc.UniqueCheck(mod); err != nil {
			return
		}

		// Assign changed values
		t.Name = mod.Name
		t.Handle = mod.Handle

		if t, err = svc.role.Update(t); err != nil {
			return err
		}

		return nil
	})
}

func (svc role) UniqueCheck(r *types.Role) (err error) {
	if r.Handle != "" {
		if ex, _ := svc.role.FindByHandle(r.Handle); ex.ID > 0 && ex.ID != r.ID {
			return ErrRoleHandleNotUnique
		}
	}

	if r.Name != "" {
		if ex, _ := svc.role.FindByName(r.Name); ex.ID > 0 && ex.ID != r.ID {
			return ErrRoleNameNotUnique
		}
	}

	return nil
}

func (svc role) Delete(roleID uint64) error {
	role, err := svc.findByID(roleID)
	if err != nil {
		return err
	}

	if !svc.ac.CanDeleteRole(svc.ctx, role) {
		return ErrNoPermissions.withStack()
	}

	// @todo: make history unavailable
	// @todo: notify users that role has been removed (remove from web UI)

	return svc.role.DeleteByID(roleID)
}

func (svc role) Archive(roleID uint64) error {
	role, err := svc.findByID(roleID)
	if err != nil {
		return err
	}

	if !svc.ac.CanUpdateRole(svc.ctx, role) {
		return ErrNoPermissions.withStack()
	}

	// @todo: make history unavailable
	// @todo: notify users that role has been removed (remove from web UI)
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

	// @todo: make history accessible
	// @todo: notify users that role has been unarchived
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

func (svc role) MemberAdd(roleID, userID uint64) error {
	role, err := svc.findByID(roleID)
	if err != nil {
		return err
	}

	if userID == 0 {
		return ErrInvalidID
	}

	if !svc.ac.CanManageRoleMembers(svc.ctx, role) {
		return errors.New("Not allowed to manage role members")
	}
	return svc.role.MemberAddByID(roleID, userID)
}

func (svc role) MemberRemove(roleID, userID uint64) error {
	role, err := svc.findByID(roleID)
	if err != nil {
		return err
	}

	if userID == 0 {
		return ErrInvalidID
	}

	if !svc.ac.CanManageRoleMembers(svc.ctx, role) {
		return errors.New("Not allowed to manage role members")
	}
	return svc.role.MemberRemoveByID(roleID, userID)
}

var _ RoleService = &role{}
