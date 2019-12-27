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
	// userBase
	//
	// This type is auto-generated.
	userBase struct {
		user    *types.User
		oldUser *types.User
		invoker auth.Identifiable
	}

	// userOnManual
	//
	// This type is auto-generated.
	userOnManual struct {
		*userBase
	}

	// userBeforeCreate
	//
	// This type is auto-generated.
	userBeforeCreate struct {
		*userBase
	}

	// userBeforeUpdate
	//
	// This type is auto-generated.
	userBeforeUpdate struct {
		*userBase
	}

	// userBeforeDelete
	//
	// This type is auto-generated.
	userBeforeDelete struct {
		*userBase
	}

	// userAfterCreate
	//
	// This type is auto-generated.
	userAfterCreate struct {
		*userBase
	}

	// userAfterUpdate
	//
	// This type is auto-generated.
	userAfterUpdate struct {
		*userBase
	}

	// userAfterDelete
	//
	// This type is auto-generated.
	userAfterDelete struct {
		*userBase
	}
)

// ResourceType returns "system:user"
//
// This function is auto-generated.
func (userBase) ResourceType() string {
	return "system:user"
}

// EventType on userOnManual returns "onManual"
//
// This function is auto-generated.
func (userOnManual) EventType() string {
	return "onManual"
}

// EventType on userBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (userBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on userBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (userBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on userBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (userBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on userAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (userAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on userAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (userAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on userAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (userAfterDelete) EventType() string {
	return "afterDelete"
}

// UserOnManual creates onManual for system:user resource
//
// This function is auto-generated.
func UserOnManual(
	argUser *types.User,
	argOldUser *types.User,
) *userOnManual {
	return &userOnManual{
		userBase: &userBase{
			user:    argUser,
			oldUser: argOldUser,
		},
	}
}

// UserBeforeCreate creates beforeCreate for system:user resource
//
// This function is auto-generated.
func UserBeforeCreate(
	argUser *types.User,
	argOldUser *types.User,
) *userBeforeCreate {
	return &userBeforeCreate{
		userBase: &userBase{
			user:    argUser,
			oldUser: argOldUser,
		},
	}
}

// UserBeforeUpdate creates beforeUpdate for system:user resource
//
// This function is auto-generated.
func UserBeforeUpdate(
	argUser *types.User,
	argOldUser *types.User,
) *userBeforeUpdate {
	return &userBeforeUpdate{
		userBase: &userBase{
			user:    argUser,
			oldUser: argOldUser,
		},
	}
}

// UserBeforeDelete creates beforeDelete for system:user resource
//
// This function is auto-generated.
func UserBeforeDelete(
	argUser *types.User,
	argOldUser *types.User,
) *userBeforeDelete {
	return &userBeforeDelete{
		userBase: &userBase{
			user:    argUser,
			oldUser: argOldUser,
		},
	}
}

// UserAfterCreate creates afterCreate for system:user resource
//
// This function is auto-generated.
func UserAfterCreate(
	argUser *types.User,
	argOldUser *types.User,
) *userAfterCreate {
	return &userAfterCreate{
		userBase: &userBase{
			user:    argUser,
			oldUser: argOldUser,
		},
	}
}

// UserAfterUpdate creates afterUpdate for system:user resource
//
// This function is auto-generated.
func UserAfterUpdate(
	argUser *types.User,
	argOldUser *types.User,
) *userAfterUpdate {
	return &userAfterUpdate{
		userBase: &userBase{
			user:    argUser,
			oldUser: argOldUser,
		},
	}
}

// UserAfterDelete creates afterDelete for system:user resource
//
// This function is auto-generated.
func UserAfterDelete(
	argUser *types.User,
	argOldUser *types.User,
) *userAfterDelete {
	return &userAfterDelete{
		userBase: &userBase{
			user:    argUser,
			oldUser: argOldUser,
		},
	}
}

// SetUser sets new user value
//
// This function is auto-generated.
func (res *userBase) SetUser(argUser *types.User) {
	res.user = argUser
}

// User returns user
//
// This function is auto-generated.
func (res userBase) User() *types.User {
	return res.user
}

// OldUser returns oldUser
//
// This function is auto-generated.
func (res userBase) OldUser() *types.User {
	return res.oldUser
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *userBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res userBase) Invoker() auth.Identifiable {
	return res.invoker
}
