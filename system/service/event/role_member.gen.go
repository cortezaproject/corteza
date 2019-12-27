package event

// This file is auto-generated.
//
// YAML event definitions:
//   system/service/event/events.yaml
//
// Regenerate with:
//   go run codegen/v2/events.go --service system
//

import (
	"github.com/cortezaproject/corteza-server/system/types"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// roleMemberBase
	//
	// This type is auto-generated.
	roleMemberBase struct {
		user    *types.User
		role    *types.Role
		invoker auth.Identifiable
	}

	// roleMemberBeforeAdd
	//
	// This type is auto-generated.
	roleMemberBeforeAdd struct {
		*roleMemberBase
	}

	// roleMemberBeforeRemove
	//
	// This type is auto-generated.
	roleMemberBeforeRemove struct {
		*roleMemberBase
	}

	// roleMemberAfterAdd
	//
	// This type is auto-generated.
	roleMemberAfterAdd struct {
		*roleMemberBase
	}

	// roleMemberAfterRemove
	//
	// This type is auto-generated.
	roleMemberAfterRemove struct {
		*roleMemberBase
	}
)

// ResourceType returns "system:role:member"
//
// This function is auto-generated.
func (roleMemberBase) ResourceType() string {
	return "system:role:member"
}

// EventType on roleMemberBeforeAdd returns "beforeAdd"
//
// This function is auto-generated.
func (roleMemberBeforeAdd) EventType() string {
	return "beforeAdd"
}

// EventType on roleMemberBeforeRemove returns "beforeRemove"
//
// This function is auto-generated.
func (roleMemberBeforeRemove) EventType() string {
	return "beforeRemove"
}

// EventType on roleMemberAfterAdd returns "afterAdd"
//
// This function is auto-generated.
func (roleMemberAfterAdd) EventType() string {
	return "afterAdd"
}

// EventType on roleMemberAfterRemove returns "afterRemove"
//
// This function is auto-generated.
func (roleMemberAfterRemove) EventType() string {
	return "afterRemove"
}

// RoleMemberBeforeAdd creates beforeAdd for system:role:member resource
//
// This function is auto-generated.
func RoleMemberBeforeAdd(
	argUser *types.User,
	argRole *types.Role,
) *roleMemberBeforeAdd {
	return &roleMemberBeforeAdd{
		roleMemberBase: &roleMemberBase{
			user: argUser,
			role: argRole,
		},
	}
}

// RoleMemberBeforeRemove creates beforeRemove for system:role:member resource
//
// This function is auto-generated.
func RoleMemberBeforeRemove(
	argUser *types.User,
	argRole *types.Role,
) *roleMemberBeforeRemove {
	return &roleMemberBeforeRemove{
		roleMemberBase: &roleMemberBase{
			user: argUser,
			role: argRole,
		},
	}
}

// RoleMemberAfterAdd creates afterAdd for system:role:member resource
//
// This function is auto-generated.
func RoleMemberAfterAdd(
	argUser *types.User,
	argRole *types.Role,
) *roleMemberAfterAdd {
	return &roleMemberAfterAdd{
		roleMemberBase: &roleMemberBase{
			user: argUser,
			role: argRole,
		},
	}
}

// RoleMemberAfterRemove creates afterRemove for system:role:member resource
//
// This function is auto-generated.
func RoleMemberAfterRemove(
	argUser *types.User,
	argRole *types.Role,
) *roleMemberAfterRemove {
	return &roleMemberAfterRemove{
		roleMemberBase: &roleMemberBase{
			user: argUser,
			role: argRole,
		},
	}
}

// SetUser sets new user value
//
// This function is auto-generated.
func (res *roleMemberBase) SetUser(argUser *types.User) {
	res.user = argUser
}

// User returns user
//
// This function is auto-generated.
func (res roleMemberBase) User() *types.User {
	return res.user
}

// SetRole sets new role value
//
// This function is auto-generated.
func (res *roleMemberBase) SetRole(argRole *types.Role) {
	res.role = argRole
}

// Role returns role
//
// This function is auto-generated.
func (res roleMemberBase) Role() *types.Role {
	return res.role
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *roleMemberBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res roleMemberBase) Invoker() auth.Identifiable {
	return res.invoker
}
