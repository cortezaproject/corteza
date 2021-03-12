package event

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/event/events.yaml

import (
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/system/automation"
	"github.com/cortezaproject/corteza-server/system/types"
)

// dummy placing to simplify import generation logic
var _ = json.NewEncoder

type (

	// systemBase
	//
	// This type is auto-generated.
	systemBase struct {
		immutable bool
		invoker   auth.Identifiable
	}

	// systemOnManual
	//
	// This type is auto-generated.
	systemOnManual struct {
		*systemBase
	}

	// systemOnInterval
	//
	// This type is auto-generated.
	systemOnInterval struct {
		*systemBase
	}

	// systemOnTimestamp
	//
	// This type is auto-generated.
	systemOnTimestamp struct {
		*systemBase
	}

	// applicationBase
	//
	// This type is auto-generated.
	applicationBase struct {
		immutable      bool
		application    *types.Application
		oldApplication *types.Application
		invoker        auth.Identifiable
	}

	// applicationOnManual
	//
	// This type is auto-generated.
	applicationOnManual struct {
		*applicationBase
	}

	// applicationBeforeCreate
	//
	// This type is auto-generated.
	applicationBeforeCreate struct {
		*applicationBase
	}

	// applicationBeforeUpdate
	//
	// This type is auto-generated.
	applicationBeforeUpdate struct {
		*applicationBase
	}

	// applicationBeforeDelete
	//
	// This type is auto-generated.
	applicationBeforeDelete struct {
		*applicationBase
	}

	// applicationAfterCreate
	//
	// This type is auto-generated.
	applicationAfterCreate struct {
		*applicationBase
	}

	// applicationAfterUpdate
	//
	// This type is auto-generated.
	applicationAfterUpdate struct {
		*applicationBase
	}

	// applicationAfterDelete
	//
	// This type is auto-generated.
	applicationAfterDelete struct {
		*applicationBase
	}

	// authBase
	//
	// This type is auto-generated.
	authBase struct {
		immutable bool
		user      *types.User
		provider  *types.AuthProvider
		invoker   auth.Identifiable
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

	// authClientBase
	//
	// This type is auto-generated.
	authClientBase struct {
		immutable     bool
		authClient    *types.AuthClient
		oldAuthClient *types.AuthClient
		invoker       auth.Identifiable
	}

	// authClientOnManual
	//
	// This type is auto-generated.
	authClientOnManual struct {
		*authClientBase
	}

	// authClientBeforeCreate
	//
	// This type is auto-generated.
	authClientBeforeCreate struct {
		*authClientBase
	}

	// authClientBeforeUpdate
	//
	// This type is auto-generated.
	authClientBeforeUpdate struct {
		*authClientBase
	}

	// authClientBeforeDelete
	//
	// This type is auto-generated.
	authClientBeforeDelete struct {
		*authClientBase
	}

	// authClientAfterCreate
	//
	// This type is auto-generated.
	authClientAfterCreate struct {
		*authClientBase
	}

	// authClientAfterUpdate
	//
	// This type is auto-generated.
	authClientAfterUpdate struct {
		*authClientBase
	}

	// authClientAfterDelete
	//
	// This type is auto-generated.
	authClientAfterDelete struct {
		*authClientBase
	}

	// mailBase
	//
	// This type is auto-generated.
	mailBase struct {
		immutable bool
		message   *types.MailMessage
		invoker   auth.Identifiable
	}

	// mailOnManual
	//
	// This type is auto-generated.
	mailOnManual struct {
		*mailBase
	}

	// mailOnReceive
	//
	// This type is auto-generated.
	mailOnReceive struct {
		*mailBase
	}

	// mailOnSend
	//
	// This type is auto-generated.
	mailOnSend struct {
		*mailBase
	}

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

	// sinkBase
	//
	// This type is auto-generated.
	sinkBase struct {
		immutable bool
		response  *types.SinkResponse
		request   *types.SinkRequest
		invoker   auth.Identifiable
	}

	// sinkOnRequest
	//
	// This type is auto-generated.
	sinkOnRequest struct {
		*sinkBase
	}

	// userBase
	//
	// This type is auto-generated.
	userBase struct {
		immutable bool
		user      *types.User
		oldUser   *types.User
		invoker   auth.Identifiable
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

// ResourceType returns "system"
//
// This function is auto-generated.
func (systemBase) ResourceType() string {
	return "system"
}

// EventType on systemOnManual returns "onManual"
//
// This function is auto-generated.
func (systemOnManual) EventType() string {
	return "onManual"
}

// EventType on systemOnInterval returns "onInterval"
//
// This function is auto-generated.
func (systemOnInterval) EventType() string {
	return "onInterval"
}

// EventType on systemOnTimestamp returns "onTimestamp"
//
// This function is auto-generated.
func (systemOnTimestamp) EventType() string {
	return "onTimestamp"
}

// SystemOnManual creates onManual for system resource
//
// This function is auto-generated.
func SystemOnManual() *systemOnManual {
	return &systemOnManual{
		systemBase: &systemBase{
			immutable: false,
		},
	}
}

// SystemOnManualImmutable creates onManual for system resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func SystemOnManualImmutable() *systemOnManual {
	return &systemOnManual{
		systemBase: &systemBase{
			immutable: true,
		},
	}
}

// SystemOnInterval creates onInterval for system resource
//
// This function is auto-generated.
func SystemOnInterval() *systemOnInterval {
	return &systemOnInterval{
		systemBase: &systemBase{
			immutable: false,
		},
	}
}

// SystemOnIntervalImmutable creates onInterval for system resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func SystemOnIntervalImmutable() *systemOnInterval {
	return &systemOnInterval{
		systemBase: &systemBase{
			immutable: true,
		},
	}
}

// SystemOnTimestamp creates onTimestamp for system resource
//
// This function is auto-generated.
func SystemOnTimestamp() *systemOnTimestamp {
	return &systemOnTimestamp{
		systemBase: &systemBase{
			immutable: false,
		},
	}
}

// SystemOnTimestampImmutable creates onTimestamp for system resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func SystemOnTimestampImmutable() *systemOnTimestamp {
	return &systemOnTimestamp{
		systemBase: &systemBase{
			immutable: true,
		},
	}
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *systemBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res systemBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res systemBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res systemBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *systemBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
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

func (res *systemBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

// ResourceType returns "system:application"
//
// This function is auto-generated.
func (applicationBase) ResourceType() string {
	return "system:application"
}

// EventType on applicationOnManual returns "onManual"
//
// This function is auto-generated.
func (applicationOnManual) EventType() string {
	return "onManual"
}

// EventType on applicationBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (applicationBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on applicationBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (applicationBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on applicationBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (applicationBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on applicationAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (applicationAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on applicationAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (applicationAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on applicationAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (applicationAfterDelete) EventType() string {
	return "afterDelete"
}

// ApplicationOnManual creates onManual for system:application resource
//
// This function is auto-generated.
func ApplicationOnManual(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationOnManual {
	return &applicationOnManual{
		applicationBase: &applicationBase{
			immutable:      false,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationOnManualImmutable creates onManual for system:application resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ApplicationOnManualImmutable(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationOnManual {
	return &applicationOnManual{
		applicationBase: &applicationBase{
			immutable:      true,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationBeforeCreate creates beforeCreate for system:application resource
//
// This function is auto-generated.
func ApplicationBeforeCreate(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationBeforeCreate {
	return &applicationBeforeCreate{
		applicationBase: &applicationBase{
			immutable:      false,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationBeforeCreateImmutable creates beforeCreate for system:application resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ApplicationBeforeCreateImmutable(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationBeforeCreate {
	return &applicationBeforeCreate{
		applicationBase: &applicationBase{
			immutable:      true,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationBeforeUpdate creates beforeUpdate for system:application resource
//
// This function is auto-generated.
func ApplicationBeforeUpdate(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationBeforeUpdate {
	return &applicationBeforeUpdate{
		applicationBase: &applicationBase{
			immutable:      false,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationBeforeUpdateImmutable creates beforeUpdate for system:application resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ApplicationBeforeUpdateImmutable(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationBeforeUpdate {
	return &applicationBeforeUpdate{
		applicationBase: &applicationBase{
			immutable:      true,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationBeforeDelete creates beforeDelete for system:application resource
//
// This function is auto-generated.
func ApplicationBeforeDelete(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationBeforeDelete {
	return &applicationBeforeDelete{
		applicationBase: &applicationBase{
			immutable:      false,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationBeforeDeleteImmutable creates beforeDelete for system:application resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ApplicationBeforeDeleteImmutable(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationBeforeDelete {
	return &applicationBeforeDelete{
		applicationBase: &applicationBase{
			immutable:      true,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationAfterCreate creates afterCreate for system:application resource
//
// This function is auto-generated.
func ApplicationAfterCreate(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationAfterCreate {
	return &applicationAfterCreate{
		applicationBase: &applicationBase{
			immutable:      false,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationAfterCreateImmutable creates afterCreate for system:application resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ApplicationAfterCreateImmutable(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationAfterCreate {
	return &applicationAfterCreate{
		applicationBase: &applicationBase{
			immutable:      true,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationAfterUpdate creates afterUpdate for system:application resource
//
// This function is auto-generated.
func ApplicationAfterUpdate(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationAfterUpdate {
	return &applicationAfterUpdate{
		applicationBase: &applicationBase{
			immutable:      false,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationAfterUpdateImmutable creates afterUpdate for system:application resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ApplicationAfterUpdateImmutable(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationAfterUpdate {
	return &applicationAfterUpdate{
		applicationBase: &applicationBase{
			immutable:      true,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationAfterDelete creates afterDelete for system:application resource
//
// This function is auto-generated.
func ApplicationAfterDelete(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationAfterDelete {
	return &applicationAfterDelete{
		applicationBase: &applicationBase{
			immutable:      false,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// ApplicationAfterDeleteImmutable creates afterDelete for system:application resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ApplicationAfterDeleteImmutable(
	argApplication *types.Application,
	argOldApplication *types.Application,
) *applicationAfterDelete {
	return &applicationAfterDelete{
		applicationBase: &applicationBase{
			immutable:      true,
			application:    argApplication,
			oldApplication: argOldApplication,
		},
	}
}

// SetApplication sets new application value
//
// This function is auto-generated.
func (res *applicationBase) SetApplication(argApplication *types.Application) {
	res.application = argApplication
}

// Application returns application
//
// This function is auto-generated.
func (res applicationBase) Application() *types.Application {
	return res.application
}

// OldApplication returns oldApplication
//
// This function is auto-generated.
func (res applicationBase) OldApplication() *types.Application {
	return res.oldApplication
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *applicationBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res applicationBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res applicationBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["application"], err = json.Marshal(res.application); err != nil {
		return nil, err
	}

	if args["oldApplication"], err = json.Marshal(res.oldApplication); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res applicationBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for *types.Application

	// Could not found expression-type counterpart for *types.Application

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *applicationBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.application != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.application); err != nil {
				return
			}
		}
	}

	if res.application != nil {
		if r, ok := results["application"]; ok {
			if err = json.Unmarshal(r, res.application); err != nil {
				return
			}
		}
	}

	// Do not decode oldApplication; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}

func (res *applicationBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for *types.Application
	// oldApplication marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

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
			immutable: false,
			user:      argUser,
			provider:  argProvider,
		},
	}
}

// AuthBeforeLoginImmutable creates beforeLogin for system:auth resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthBeforeLoginImmutable(
	argUser *types.User,
	argProvider *types.AuthProvider,
) *authBeforeLogin {
	return &authBeforeLogin{
		authBase: &authBase{
			immutable: true,
			user:      argUser,
			provider:  argProvider,
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
			immutable: false,
			user:      argUser,
			provider:  argProvider,
		},
	}
}

// AuthBeforeSignupImmutable creates beforeSignup for system:auth resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthBeforeSignupImmutable(
	argUser *types.User,
	argProvider *types.AuthProvider,
) *authBeforeSignup {
	return &authBeforeSignup{
		authBase: &authBase{
			immutable: true,
			user:      argUser,
			provider:  argProvider,
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
			immutable: false,
			user:      argUser,
			provider:  argProvider,
		},
	}
}

// AuthAfterLoginImmutable creates afterLogin for system:auth resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthAfterLoginImmutable(
	argUser *types.User,
	argProvider *types.AuthProvider,
) *authAfterLogin {
	return &authAfterLogin{
		authBase: &authBase{
			immutable: true,
			user:      argUser,
			provider:  argProvider,
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
			immutable: false,
			user:      argUser,
			provider:  argProvider,
		},
	}
}

// AuthAfterSignupImmutable creates afterSignup for system:auth resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthAfterSignupImmutable(
	argUser *types.User,
	argProvider *types.AuthProvider,
) *authAfterSignup {
	return &authAfterSignup{
		authBase: &authBase{
			immutable: true,
			user:      argUser,
			provider:  argProvider,
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res authBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["user"], err = json.Marshal(res.user); err != nil {
		return nil, err
	}

	if args["provider"], err = json.Marshal(res.provider); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res authBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	if rvars["user"], err = automation.NewUser(res.user); err != nil {
		return nil, err
	}

	// Could not found expression-type counterpart for *types.AuthProvider

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *authBase) Decode(results map[string][]byte) (err error) {
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

	if res.provider != nil {
		if r, ok := results["provider"]; ok {
			if err = json.Unmarshal(r, res.provider); err != nil {
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

func (res *authBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.user != nil && vars.Has("user") {
		var aux *automation.User
		aux, err = automation.NewUser(expr.Must(vars.Select("user")))
		if err != nil {
			return
		}

		res.user = aux.GetValue()
	}
	// Could not find expression-type counterpart for *types.AuthProvider
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

// ResourceType returns "system:auth-client"
//
// This function is auto-generated.
func (authClientBase) ResourceType() string {
	return "system:auth-client"
}

// EventType on authClientOnManual returns "onManual"
//
// This function is auto-generated.
func (authClientOnManual) EventType() string {
	return "onManual"
}

// EventType on authClientBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (authClientBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on authClientBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (authClientBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on authClientBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (authClientBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on authClientAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (authClientAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on authClientAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (authClientAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on authClientAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (authClientAfterDelete) EventType() string {
	return "afterDelete"
}

// AuthClientOnManual creates onManual for system:auth-client resource
//
// This function is auto-generated.
func AuthClientOnManual(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientOnManual {
	return &authClientOnManual{
		authClientBase: &authClientBase{
			immutable:     false,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientOnManualImmutable creates onManual for system:auth-client resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthClientOnManualImmutable(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientOnManual {
	return &authClientOnManual{
		authClientBase: &authClientBase{
			immutable:     true,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientBeforeCreate creates beforeCreate for system:auth-client resource
//
// This function is auto-generated.
func AuthClientBeforeCreate(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientBeforeCreate {
	return &authClientBeforeCreate{
		authClientBase: &authClientBase{
			immutable:     false,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientBeforeCreateImmutable creates beforeCreate for system:auth-client resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthClientBeforeCreateImmutable(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientBeforeCreate {
	return &authClientBeforeCreate{
		authClientBase: &authClientBase{
			immutable:     true,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientBeforeUpdate creates beforeUpdate for system:auth-client resource
//
// This function is auto-generated.
func AuthClientBeforeUpdate(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientBeforeUpdate {
	return &authClientBeforeUpdate{
		authClientBase: &authClientBase{
			immutable:     false,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientBeforeUpdateImmutable creates beforeUpdate for system:auth-client resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthClientBeforeUpdateImmutable(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientBeforeUpdate {
	return &authClientBeforeUpdate{
		authClientBase: &authClientBase{
			immutable:     true,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientBeforeDelete creates beforeDelete for system:auth-client resource
//
// This function is auto-generated.
func AuthClientBeforeDelete(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientBeforeDelete {
	return &authClientBeforeDelete{
		authClientBase: &authClientBase{
			immutable:     false,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientBeforeDeleteImmutable creates beforeDelete for system:auth-client resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthClientBeforeDeleteImmutable(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientBeforeDelete {
	return &authClientBeforeDelete{
		authClientBase: &authClientBase{
			immutable:     true,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientAfterCreate creates afterCreate for system:auth-client resource
//
// This function is auto-generated.
func AuthClientAfterCreate(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientAfterCreate {
	return &authClientAfterCreate{
		authClientBase: &authClientBase{
			immutable:     false,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientAfterCreateImmutable creates afterCreate for system:auth-client resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthClientAfterCreateImmutable(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientAfterCreate {
	return &authClientAfterCreate{
		authClientBase: &authClientBase{
			immutable:     true,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientAfterUpdate creates afterUpdate for system:auth-client resource
//
// This function is auto-generated.
func AuthClientAfterUpdate(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientAfterUpdate {
	return &authClientAfterUpdate{
		authClientBase: &authClientBase{
			immutable:     false,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientAfterUpdateImmutable creates afterUpdate for system:auth-client resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthClientAfterUpdateImmutable(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientAfterUpdate {
	return &authClientAfterUpdate{
		authClientBase: &authClientBase{
			immutable:     true,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientAfterDelete creates afterDelete for system:auth-client resource
//
// This function is auto-generated.
func AuthClientAfterDelete(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientAfterDelete {
	return &authClientAfterDelete{
		authClientBase: &authClientBase{
			immutable:     false,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// AuthClientAfterDeleteImmutable creates afterDelete for system:auth-client resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func AuthClientAfterDeleteImmutable(
	argAuthClient *types.AuthClient,
	argOldAuthClient *types.AuthClient,
) *authClientAfterDelete {
	return &authClientAfterDelete{
		authClientBase: &authClientBase{
			immutable:     true,
			authClient:    argAuthClient,
			oldAuthClient: argOldAuthClient,
		},
	}
}

// SetAuthClient sets new authClient value
//
// This function is auto-generated.
func (res *authClientBase) SetAuthClient(argAuthClient *types.AuthClient) {
	res.authClient = argAuthClient
}

// AuthClient returns authClient
//
// This function is auto-generated.
func (res authClientBase) AuthClient() *types.AuthClient {
	return res.authClient
}

// OldAuthClient returns oldAuthClient
//
// This function is auto-generated.
func (res authClientBase) OldAuthClient() *types.AuthClient {
	return res.oldAuthClient
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *authClientBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res authClientBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res authClientBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["authClient"], err = json.Marshal(res.authClient); err != nil {
		return nil, err
	}

	if args["oldAuthClient"], err = json.Marshal(res.oldAuthClient); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res authClientBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for *types.AuthClient

	// Could not found expression-type counterpart for *types.AuthClient

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *authClientBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.authClient != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.authClient); err != nil {
				return
			}
		}
	}

	if res.authClient != nil {
		if r, ok := results["authClient"]; ok {
			if err = json.Unmarshal(r, res.authClient); err != nil {
				return
			}
		}
	}

	// Do not decode oldAuthClient; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}

func (res *authClientBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for *types.AuthClient
	// oldAuthClient marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

// ResourceType returns "system:mail"
//
// This function is auto-generated.
func (mailBase) ResourceType() string {
	return "system:mail"
}

// EventType on mailOnManual returns "onManual"
//
// This function is auto-generated.
func (mailOnManual) EventType() string {
	return "onManual"
}

// EventType on mailOnReceive returns "onReceive"
//
// This function is auto-generated.
func (mailOnReceive) EventType() string {
	return "onReceive"
}

// EventType on mailOnSend returns "onSend"
//
// This function is auto-generated.
func (mailOnSend) EventType() string {
	return "onSend"
}

// MailOnManual creates onManual for system:mail resource
//
// This function is auto-generated.
func MailOnManual(
	argMessage *types.MailMessage,
) *mailOnManual {
	return &mailOnManual{
		mailBase: &mailBase{
			immutable: false,
			message:   argMessage,
		},
	}
}

// MailOnManualImmutable creates onManual for system:mail resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MailOnManualImmutable(
	argMessage *types.MailMessage,
) *mailOnManual {
	return &mailOnManual{
		mailBase: &mailBase{
			immutable: true,
			message:   argMessage,
		},
	}
}

// MailOnReceive creates onReceive for system:mail resource
//
// This function is auto-generated.
func MailOnReceive(
	argMessage *types.MailMessage,
) *mailOnReceive {
	return &mailOnReceive{
		mailBase: &mailBase{
			immutable: false,
			message:   argMessage,
		},
	}
}

// MailOnReceiveImmutable creates onReceive for system:mail resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MailOnReceiveImmutable(
	argMessage *types.MailMessage,
) *mailOnReceive {
	return &mailOnReceive{
		mailBase: &mailBase{
			immutable: true,
			message:   argMessage,
		},
	}
}

// MailOnSend creates onSend for system:mail resource
//
// This function is auto-generated.
func MailOnSend(
	argMessage *types.MailMessage,
) *mailOnSend {
	return &mailOnSend{
		mailBase: &mailBase{
			immutable: false,
			message:   argMessage,
		},
	}
}

// MailOnSendImmutable creates onSend for system:mail resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func MailOnSendImmutable(
	argMessage *types.MailMessage,
) *mailOnSend {
	return &mailOnSend{
		mailBase: &mailBase{
			immutable: true,
			message:   argMessage,
		},
	}
}

// SetMessage sets new message value
//
// This function is auto-generated.
func (res *mailBase) SetMessage(argMessage *types.MailMessage) {
	res.message = argMessage
}

// Message returns message
//
// This function is auto-generated.
func (res mailBase) Message() *types.MailMessage {
	return res.message
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *mailBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res mailBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res mailBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["message"], err = json.Marshal(res.message); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res mailBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for *types.MailMessage

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *mailBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.message != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.message); err != nil {
				return
			}
		}
	}

	if res.message != nil {
		if r, ok := results["message"]; ok {
			if err = json.Unmarshal(r, res.message); err != nil {
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

func (res *mailBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for *types.MailMessage
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

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

// Encode internal data to be passed as event params & arguments to workflow
func (res roleBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	if rvars["role"], err = automation.NewRole(res.role); err != nil {
		return nil, err
	}

	if rvars["oldRole"], err = automation.NewRole(res.oldRole); err != nil {
		return nil, err
	}

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
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

func (res *roleBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.role != nil && vars.Has("role") {
		var aux *automation.Role
		aux, err = automation.NewRole(expr.Must(vars.Select("role")))
		if err != nil {
			return
		}

		res.role = aux.GetValue()
	}
	// oldRole marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

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

// Encode internal data to be passed as event params & arguments to workflow
func (res roleMemberBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	if rvars["user"], err = automation.NewUser(res.user); err != nil {
		return nil, err
	}

	if rvars["role"], err = automation.NewRole(res.role); err != nil {
		return nil, err
	}

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
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

func (res *roleMemberBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.user != nil && vars.Has("user") {
		var aux *automation.User
		aux, err = automation.NewUser(expr.Must(vars.Select("user")))
		if err != nil {
			return
		}

		res.user = aux.GetValue()
	}
	if res.role != nil && vars.Has("role") {
		var aux *automation.Role
		aux, err = automation.NewRole(expr.Must(vars.Select("role")))
		if err != nil {
			return
		}

		res.role = aux.GetValue()
	}
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

// ResourceType returns "system:sink"
//
// This function is auto-generated.
func (sinkBase) ResourceType() string {
	return "system:sink"
}

// EventType on sinkOnRequest returns "onRequest"
//
// This function is auto-generated.
func (sinkOnRequest) EventType() string {
	return "onRequest"
}

// SinkOnRequest creates onRequest for system:sink resource
//
// This function is auto-generated.
func SinkOnRequest(
	argResponse *types.SinkResponse,
	argRequest *types.SinkRequest,
) *sinkOnRequest {
	return &sinkOnRequest{
		sinkBase: &sinkBase{
			immutable: false,
			response:  argResponse,
			request:   argRequest,
		},
	}
}

// SinkOnRequestImmutable creates onRequest for system:sink resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func SinkOnRequestImmutable(
	argResponse *types.SinkResponse,
	argRequest *types.SinkRequest,
) *sinkOnRequest {
	return &sinkOnRequest{
		sinkBase: &sinkBase{
			immutable: true,
			response:  argResponse,
			request:   argRequest,
		},
	}
}

// SetResponse sets new response value
//
// This function is auto-generated.
func (res *sinkBase) SetResponse(argResponse *types.SinkResponse) {
	res.response = argResponse
}

// Response returns response
//
// This function is auto-generated.
func (res sinkBase) Response() *types.SinkResponse {
	return res.response
}

// Request returns request
//
// This function is auto-generated.
func (res sinkBase) Request() *types.SinkRequest {
	return res.request
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *sinkBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res sinkBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res sinkBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["response"], err = json.Marshal(res.response); err != nil {
		return nil, err
	}

	if args["request"], err = json.Marshal(res.request); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res sinkBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for *types.SinkResponse

	// Could not found expression-type counterpart for *types.SinkRequest

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *sinkBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.response != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.response); err != nil {
				return
			}
		}
	}

	if res.response != nil {
		if r, ok := results["response"]; ok {
			if err = json.Unmarshal(r, res.response); err != nil {
				return
			}
		}
	}

	// Do not decode request; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}

func (res *sinkBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for *types.SinkResponse
	// request marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

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
			immutable: false,
			user:      argUser,
			oldUser:   argOldUser,
		},
	}
}

// UserOnManualImmutable creates onManual for system:user resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func UserOnManualImmutable(
	argUser *types.User,
	argOldUser *types.User,
) *userOnManual {
	return &userOnManual{
		userBase: &userBase{
			immutable: true,
			user:      argUser,
			oldUser:   argOldUser,
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
			immutable: false,
			user:      argUser,
			oldUser:   argOldUser,
		},
	}
}

// UserBeforeCreateImmutable creates beforeCreate for system:user resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func UserBeforeCreateImmutable(
	argUser *types.User,
	argOldUser *types.User,
) *userBeforeCreate {
	return &userBeforeCreate{
		userBase: &userBase{
			immutable: true,
			user:      argUser,
			oldUser:   argOldUser,
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
			immutable: false,
			user:      argUser,
			oldUser:   argOldUser,
		},
	}
}

// UserBeforeUpdateImmutable creates beforeUpdate for system:user resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func UserBeforeUpdateImmutable(
	argUser *types.User,
	argOldUser *types.User,
) *userBeforeUpdate {
	return &userBeforeUpdate{
		userBase: &userBase{
			immutable: true,
			user:      argUser,
			oldUser:   argOldUser,
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
			immutable: false,
			user:      argUser,
			oldUser:   argOldUser,
		},
	}
}

// UserBeforeDeleteImmutable creates beforeDelete for system:user resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func UserBeforeDeleteImmutable(
	argUser *types.User,
	argOldUser *types.User,
) *userBeforeDelete {
	return &userBeforeDelete{
		userBase: &userBase{
			immutable: true,
			user:      argUser,
			oldUser:   argOldUser,
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
			immutable: false,
			user:      argUser,
			oldUser:   argOldUser,
		},
	}
}

// UserAfterCreateImmutable creates afterCreate for system:user resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func UserAfterCreateImmutable(
	argUser *types.User,
	argOldUser *types.User,
) *userAfterCreate {
	return &userAfterCreate{
		userBase: &userBase{
			immutable: true,
			user:      argUser,
			oldUser:   argOldUser,
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
			immutable: false,
			user:      argUser,
			oldUser:   argOldUser,
		},
	}
}

// UserAfterUpdateImmutable creates afterUpdate for system:user resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func UserAfterUpdateImmutable(
	argUser *types.User,
	argOldUser *types.User,
) *userAfterUpdate {
	return &userAfterUpdate{
		userBase: &userBase{
			immutable: true,
			user:      argUser,
			oldUser:   argOldUser,
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
			immutable: false,
			user:      argUser,
			oldUser:   argOldUser,
		},
	}
}

// UserAfterDeleteImmutable creates afterDelete for system:user resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func UserAfterDeleteImmutable(
	argUser *types.User,
	argOldUser *types.User,
) *userAfterDelete {
	return &userAfterDelete{
		userBase: &userBase{
			immutable: true,
			user:      argUser,
			oldUser:   argOldUser,
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res userBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["user"], err = json.Marshal(res.user); err != nil {
		return nil, err
	}

	if args["oldUser"], err = json.Marshal(res.oldUser); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res userBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	if rvars["user"], err = automation.NewUser(res.user); err != nil {
		return nil, err
	}

	if rvars["oldUser"], err = automation.NewUser(res.oldUser); err != nil {
		return nil, err
	}

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *userBase) Decode(results map[string][]byte) (err error) {
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

	// Do not decode oldUser; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}

func (res *userBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.user != nil && vars.Has("user") {
		var aux *automation.User
		aux, err = automation.NewUser(expr.Must(vars.Select("user")))
		if err != nil {
			return
		}

		res.user = aux.GetValue()
	}
	// oldUser marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}
