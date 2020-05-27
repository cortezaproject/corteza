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
	// roleBase
	//
	// This type is auto-generated.
	roleBase struct {
		immutable bool
		role      *types.Role
		oldRole   *types.Role
		invoker   auth.Identifiable
	}

	// roleOnManual
	//
	// This type is auto-generated.
	roleOnManual struct {
		*roleBase
	}

	// roleBeforeCreate
	//
	// This type is auto-generated.
	roleBeforeCreate struct {
		*roleBase
	}

	// roleBeforeUpdate
	//
	// This type is auto-generated.
	roleBeforeUpdate struct {
		*roleBase
	}

	// roleBeforeDelete
	//
	// This type is auto-generated.
	roleBeforeDelete struct {
		*roleBase
	}

	// roleAfterCreate
	//
	// This type is auto-generated.
	roleAfterCreate struct {
		*roleBase
	}

	// roleAfterUpdate
	//
	// This type is auto-generated.
	roleAfterUpdate struct {
		*roleBase
	}

	// roleAfterDelete
	//
	// This type is auto-generated.
	roleAfterDelete struct {
		*roleBase
	}
)

// ResourceType returns "system:role"
//
// This function is auto-generated.
func (roleBase) ResourceType() string {
	return "system:role"
}

// EventType on roleOnManual returns "onManual"
//
// This function is auto-generated.
func (roleOnManual) EventType() string {
	return "onManual"
}

// EventType on roleBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (roleBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on roleBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (roleBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on roleBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (roleBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on roleAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (roleAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on roleAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (roleAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on roleAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (roleAfterDelete) EventType() string {
	return "afterDelete"
}

// RoleOnManual creates onManual for system:role resource
//
// This function is auto-generated.
func RoleOnManual(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleOnManual {
	return &roleOnManual{
		roleBase: &roleBase{
			immutable: false,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleOnManualImmutable creates onManual for system:role resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleOnManualImmutable(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleOnManual {
	return &roleOnManual{
		roleBase: &roleBase{
			immutable: true,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleBeforeCreate creates beforeCreate for system:role resource
//
// This function is auto-generated.
func RoleBeforeCreate(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleBeforeCreate {
	return &roleBeforeCreate{
		roleBase: &roleBase{
			immutable: false,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleBeforeCreateImmutable creates beforeCreate for system:role resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleBeforeCreateImmutable(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleBeforeCreate {
	return &roleBeforeCreate{
		roleBase: &roleBase{
			immutable: true,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleBeforeUpdate creates beforeUpdate for system:role resource
//
// This function is auto-generated.
func RoleBeforeUpdate(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleBeforeUpdate {
	return &roleBeforeUpdate{
		roleBase: &roleBase{
			immutable: false,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleBeforeUpdateImmutable creates beforeUpdate for system:role resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleBeforeUpdateImmutable(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleBeforeUpdate {
	return &roleBeforeUpdate{
		roleBase: &roleBase{
			immutable: true,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleBeforeDelete creates beforeDelete for system:role resource
//
// This function is auto-generated.
func RoleBeforeDelete(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleBeforeDelete {
	return &roleBeforeDelete{
		roleBase: &roleBase{
			immutable: false,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleBeforeDeleteImmutable creates beforeDelete for system:role resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleBeforeDeleteImmutable(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleBeforeDelete {
	return &roleBeforeDelete{
		roleBase: &roleBase{
			immutable: true,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleAfterCreate creates afterCreate for system:role resource
//
// This function is auto-generated.
func RoleAfterCreate(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleAfterCreate {
	return &roleAfterCreate{
		roleBase: &roleBase{
			immutable: false,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleAfterCreateImmutable creates afterCreate for system:role resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleAfterCreateImmutable(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleAfterCreate {
	return &roleAfterCreate{
		roleBase: &roleBase{
			immutable: true,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleAfterUpdate creates afterUpdate for system:role resource
//
// This function is auto-generated.
func RoleAfterUpdate(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleAfterUpdate {
	return &roleAfterUpdate{
		roleBase: &roleBase{
			immutable: false,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleAfterUpdateImmutable creates afterUpdate for system:role resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleAfterUpdateImmutable(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleAfterUpdate {
	return &roleAfterUpdate{
		roleBase: &roleBase{
			immutable: true,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleAfterDelete creates afterDelete for system:role resource
//
// This function is auto-generated.
func RoleAfterDelete(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleAfterDelete {
	return &roleAfterDelete{
		roleBase: &roleBase{
			immutable: false,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// RoleAfterDeleteImmutable creates afterDelete for system:role resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RoleAfterDeleteImmutable(
	argRole *types.Role,
	argOldRole *types.Role,
) *roleAfterDelete {
	return &roleAfterDelete{
		roleBase: &roleBase{
			immutable: true,
			role:      argRole,
			oldRole:   argOldRole,
		},
	}
}

// SetRole sets new role value
//
// This function is auto-generated.
func (res *roleBase) SetRole(argRole *types.Role) {
	res.role = argRole
}

// Role returns role
//
// This function is auto-generated.
func (res roleBase) Role() *types.Role {
	return res.role
}

// OldRole returns oldRole
//
// This function is auto-generated.
func (res roleBase) OldRole() *types.Role {
	return res.oldRole
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *roleBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res roleBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res roleBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["role"], err = json.Marshal(res.role); err != nil {
		return nil, err
	}

	if args["oldRole"], err = json.Marshal(res.oldRole); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Decode return values from Corredor script into struct props
func (res *roleBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.role != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.role); err != nil {
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

	// Do not decode oldRole; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}
