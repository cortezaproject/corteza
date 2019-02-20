package service

import (
	"context"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/repository"
	"github.com/crusttech/crust/system/types"
)

type (
	role struct {
		db  *factory.DB
		ctx context.Context

		role repository.RoleRepository
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

		MemberAdd(roleID, userID uint64) error
		MemberRemove(roleID, userID uint64) error
	}
)

func Role() RoleService {
	return (&role{}).With(context.Background())
}

func (svc *role) With(ctx context.Context) RoleService {
	db := repository.DB(ctx)
	return &role{
		db:   db,
		ctx:  ctx,
		role: repository.Role(ctx, db),
	}
}

func (svc *role) FindByID(id uint64) (*types.Role, error) {
	// @todo: permission check if current user has access to this role
	return svc.role.FindByID(id)
}

func (svc *role) Find(filter *types.RoleFilter) ([]*types.Role, error) {
	// @todo: permission check to return only roles that current user has access to
	return svc.role.Find(filter)
}

func (svc *role) Create(mod *types.Role) (*types.Role, error) {
	// @todo: permission check if current user can add/edit role

	return svc.role.Create(mod)
}

func (svc *role) Update(mod *types.Role) (t *types.Role, err error) {
	// @todo: permission check if current user can add/edit role
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

func (svc *role) Delete(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that role has been removed (remove from web UI)
	// @todo: permissions check if current user can remove role
	return svc.role.DeleteByID(id)
}

func (svc *role) Archive(id uint64) error {
	// @todo: make history unavailable
	// @todo: notify users that role has been removed (remove from web UI)
	// @todo: permissions check if current user can remove role
	return svc.role.ArchiveByID(id)
}

func (svc *role) Unarchive(id uint64) error {
	// @todo: permissions check if current user can unarchive role
	// @todo: make history accessible
	// @todo: notify users that role has been unarchived
	return svc.role.UnarchiveByID(id)
}

func (svc *role) Merge(id, targetroleID uint64) error {
	// @todo: permission check if current user can merge role
	return svc.role.MergeByID(id, targetroleID)
}

func (svc *role) Move(id, targetOrganisationID uint64) error {
	// @todo: permission check if current user can move role to another organisation
	return svc.role.MoveByID(id, targetOrganisationID)
}

func (svc *role) MemberAdd(id, userID uint64) error {
	// @todo: permission check if current user can add user in to a role
	return svc.role.MemberAddByID(id, userID)
}

func (svc *role) MemberRemove(id, userID uint64) error {
	// @todo: permission check if current user can remove user from a role
	return svc.role.MemberRemoveByID(id, userID)
}

var _ RoleService = &role{}
