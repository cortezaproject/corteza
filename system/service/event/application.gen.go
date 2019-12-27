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
	// applicationBase
	//
	// This type is auto-generated.
	applicationBase struct {
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
