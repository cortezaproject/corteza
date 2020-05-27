package event

// This file is auto-generated.
//
// YAML event definitions:
//   system/service/event/events.yaml
//
// Regenerate with:
//   go run ./codegen/v2/events --service system
//

import (
	"encoding/json"

	"github.com/cortezaproject/corteza-server/system/types"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// roleMemberBase
	//
	// This type is auto-generated.
	roleMemberBase struct {
		immutable bool
		user      *types.User
		role      *types.Role
		invoker   auth.Identifiable
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
			immutable: false,
			user:      argUser,
			role:      argRole,
		},
	}
}

// RoleMemberBeforeAddImmutable creates beforeAdd for system:role:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleMemberBeforeAddImmutable(
	argUser *types.User,
	argRole *types.Role,
) *roleMemberBeforeAdd {
	return &roleMemberBeforeAdd{
		roleMemberBase: &roleMemberBase{
			immutable: true,
			user:      argUser,
			role:      argRole,
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
			immutable: false,
			user:      argUser,
			role:      argRole,
		},
	}
}

// RoleMemberBeforeRemoveImmutable creates beforeRemove for system:role:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleMemberBeforeRemoveImmutable(
	argUser *types.User,
	argRole *types.Role,
) *roleMemberBeforeRemove {
	return &roleMemberBeforeRemove{
		roleMemberBase: &roleMemberBase{
			immutable: true,
			user:      argUser,
			role:      argRole,
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
			immutable: false,
			user:      argUser,
			role:      argRole,
		},
	}
}

// RoleMemberAfterAddImmutable creates afterAdd for system:role:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleMemberAfterAddImmutable(
	argUser *types.User,
	argRole *types.Role,
) *roleMemberAfterAdd {
	return &roleMemberAfterAdd{
		roleMemberBase: &roleMemberBase{
			immutable: true,
			user:      argUser,
			role:      argRole,
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
			immutable: false,
			user:      argUser,
			role:      argRole,
		},
	}
}

// RoleMemberAfterRemoveImmutable creates afterRemove for system:role:member resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleMemberAfterRemoveImmutable(
	argUser *types.User,
	argRole *types.Role,
) *roleMemberAfterRemove {
	return &roleMemberAfterRemove{
		roleMemberBase: &roleMemberBase{
			immutable: true,
			user:      argUser,
			role:      argRole,
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res roleMemberBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["user"], err = json.Marshal(res.user); err != nil {
		return nil, err
	}

	if args["role"], err = json.Marshal(res.role); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Decode return values from Corredor script into struct props
func (res *roleMemberBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.user != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.user); err != nil {
				return
			}
		}
	}

	if res.user != nil {
		if r, ok := results["user"]; ok {
			if err = json.Unmarshal(r, res.user); err != nil {
				return
			}
		}
	}

	if res.role != nil {
		if r, ok := results["role"]; ok {
			if err = json.Unmarshal(r, res.role); err != nil {
				return
			}
		}
	}

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}
