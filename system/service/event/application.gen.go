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
)

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
