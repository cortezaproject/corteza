package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/crusttech/crust/internal/logger"
	"github.com/crusttech/crust/system/internal/repository"
	"github.com/crusttech/crust/system/types"
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
		Find(filter *types.RoleFilter) ([]*types.Role, error)

		Create(role *types.Role) (*types.Role, error)
		Update(role *types.Role) (*types.Role, error)
		Merge(roleID, targetroleID uint64) error
		Move(roleID, organisationID uint64) error

		Archive(ID uint64) error
		Unarchive(ID uint64) error
		Delete(ID uint64) error

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

// log() returns zap's logger with requestID from current context and fields.
func (svc role) log(fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
}

func (svc role) FindByID(id uint64) (*types.Role, error) {
	role, err := svc.role.FindByID(id)
	if err != nil {
		return nil, err
	}

	if !svc.ac.CanReadRole(svc.ctx, role) {
		return nil, errors.New("Not allowed to read role")
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

func (svc role) Create(mod *types.Role) (*types.Role, error) {
	if !svc.ac.CanCreateRole(svc.ctx) {
		return nil, errors.New("Not allowed to create role")
	}
	return svc.role.Create(mod)
}

func (svc role) Update(mod *types.Role) (t *types.Role, err error) {
	if !svc.ac.CanUpdateRole(svc.ctx, mod) {
		return nil, errors.New("Not allowed to update role")
	}

	// @todo: make sure archived & deleted entries can not be edited

	return t, svc.db.Transaction(func() (err error) {
		if t, err = svc.role.FindByID(mod.ID); err != nil {
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

func (svc role) Delete(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that role has been removed (remove from web UI)

	rl := &types.Role{ID: id}
	if !svc.ac.CanDeleteRole(svc.ctx, rl) {
		return errors.New("Not allowed to delete role")
	}
	return svc.role.DeleteByID(id)
}

func (svc role) Archive(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that role has been removed (remove from web UI)
	// @todo: permissions check if current user can remove role
	return svc.role.ArchiveByID(id)
}

func (svc role) Unarchive(id uint64) error {
	// @todo: permissions check if current user can unarchive role
	// @todo: make history accessible
	// @todo: notify users that role has been unarchived
	return svc.role.UnarchiveByID(id)
}

func (svc role) Merge(id, targetroleID uint64) error {
	// @todo: permission check if current user can merge role
	return svc.role.MergeByID(id, targetroleID)
}

func (svc role) Move(id, targetOrganisationID uint64) error {
	// @todo: permission check if current user can move role to another organisation
	return svc.role.MoveByID(id, targetOrganisationID)
}

func (svc role) MemberList(roleID uint64) ([]*types.RoleMember, error) {
	rl := &types.Role{ID: roleID}
	if !svc.ac.CanManageRoleMembers(svc.ctx, rl) {
		return nil, errors.New("Not allowed to manage role members")
	}
	return svc.role.MemberFindByRoleID(roleID)
}

func (svc role) MemberAdd(roleID, userID uint64) error {
	rl := &types.Role{ID: roleID}
	if !svc.ac.CanManageRoleMembers(svc.ctx, rl) {
		return errors.New("Not allowed to manage role members")
	}
	return svc.role.MemberAddByID(roleID, userID)
}

func (svc role) MemberRemove(roleID, userID uint64) error {
	rl := &types.Role{ID: roleID}
	if !svc.ac.CanManageRoleMembers(svc.ctx, rl) {
		return errors.New("Not allowed to manage role members")
	}
	return svc.role.MemberRemoveByID(roleID, userID)
}

var _ RoleService = &role{}
