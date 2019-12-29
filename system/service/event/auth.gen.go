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
	// authBase
	//
	// This type is auto-generated.
	authBase struct {
		user     *types.User
		provider *types.AuthProvider
		invoker  auth.Identifiable
	}

	// authBeforeLogin
	//
	// This type is auto-generated.
	authBeforeLogin struct {
		*authBase
	}

	// authBeforeSignup
	//
	// This type is auto-generated.
	authBeforeSignup struct {
		*authBase
	}

	// authAfterLogin
	//
	// This type is auto-generated.
	authAfterLogin struct {
		*authBase
	}

	// authAfterSignup
	//
	// This type is auto-generated.
	authAfterSignup struct {
		*authBase
	}
)

// ResourceType returns "system:auth"
//
// This function is auto-generated.
func (authBase) ResourceType() string {
	return "system:auth"
}

// EventType on authBeforeLogin returns "beforeLogin"
//
// This function is auto-generated.
func (authBeforeLogin) EventType() string {
	return "beforeLogin"
}

// EventType on authBeforeSignup returns "beforeSignup"
//
// This function is auto-generated.
func (authBeforeSignup) EventType() string {
	return "beforeSignup"
}

// EventType on authAfterLogin returns "afterLogin"
//
// This function is auto-generated.
func (authAfterLogin) EventType() string {
	return "afterLogin"
}

// EventType on authAfterSignup returns "afterSignup"
//
// This function is auto-generated.
func (authAfterSignup) EventType() string {
	return "afterSignup"
}

// AuthBeforeLogin creates beforeLogin for system:auth resource
//
// This function is auto-generated.
func AuthBeforeLogin(
	argUser *types.User,
	argProvider *types.AuthProvider,
) *authBeforeLogin {
	return &authBeforeLogin{
		authBase: &authBase{
			user:     argUser,
			provider: argProvider,
		},
	}
}

// AuthBeforeSignup creates beforeSignup for system:auth resource
//
// This function is auto-generated.
func AuthBeforeSignup(
	argUser *types.User,
	argProvider *types.AuthProvider,
) *authBeforeSignup {
	return &authBeforeSignup{
		authBase: &authBase{
			user:     argUser,
			provider: argProvider,
		},
	}
}

// AuthAfterLogin creates afterLogin for system:auth resource
//
// This function is auto-generated.
func AuthAfterLogin(
	argUser *types.User,
	argProvider *types.AuthProvider,
) *authAfterLogin {
	return &authAfterLogin{
		authBase: &authBase{
			user:     argUser,
			provider: argProvider,
		},
	}
}

// AuthAfterSignup creates afterSignup for system:auth resource
//
// This function is auto-generated.
func AuthAfterSignup(
	argUser *types.User,
	argProvider *types.AuthProvider,
) *authAfterSignup {
	return &authAfterSignup{
		authBase: &authBase{
			user:     argUser,
			provider: argProvider,
		},
	}
}

// SetUser sets new user value
//
// This function is auto-generated.
func (res *authBase) SetUser(argUser *types.User) {
	res.user = argUser
}

// User returns user
//
// This function is auto-generated.
func (res authBase) User() *types.User {
	return res.user
}

// SetProvider sets new provider value
//
// This function is auto-generated.
func (res *authBase) SetProvider(argProvider *types.AuthProvider) {
	res.provider = argProvider
}

// Provider returns provider
//
// This function is auto-generated.
func (res authBase) Provider() *types.AuthProvider {
	return res.provider
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *authBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res authBase) Invoker() auth.Identifiable {
	return res.invoker
}
