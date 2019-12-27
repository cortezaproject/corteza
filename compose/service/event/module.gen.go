package event

// This file is auto-generated.
//
// YAML event definitions:
//   compose/service/event/events.yaml
//
// Regenerate with:
//   go run codegen/v2/events.go --service compose
//

import (
	"github.com/cortezaproject/corteza-server/compose/types"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// moduleBase
	//
	// This type is auto-generated.
	moduleBase struct {
		module    *types.Module
		oldModule *types.Module
		namespace *types.Namespace
		invoker   auth.Identifiable
	}

	// moduleOnManual
	//
	// This type is auto-generated.
	moduleOnManual struct {
		*moduleBase
	}

	// moduleBeforeCreate
	//
	// This type is auto-generated.
	moduleBeforeCreate struct {
		*moduleBase
	}

	// moduleBeforeUpdate
	//
	// This type is auto-generated.
	moduleBeforeUpdate struct {
		*moduleBase
	}

	// moduleBeforeDelete
	//
	// This type is auto-generated.
	moduleBeforeDelete struct {
		*moduleBase
	}

	// moduleAfterCreate
	//
	// This type is auto-generated.
	moduleAfterCreate struct {
		*moduleBase
	}

	// moduleAfterUpdate
	//
	// This type is auto-generated.
	moduleAfterUpdate struct {
		*moduleBase
	}

	// moduleAfterDelete
	//
	// This type is auto-generated.
	moduleAfterDelete struct {
		*moduleBase
	}
)

// ResourceType returns "compose:module"
//
// This function is auto-generated.
func (moduleBase) ResourceType() string {
	return "compose:module"
}

// EventType on moduleOnManual returns "onManual"
//
// This function is auto-generated.
func (moduleOnManual) EventType() string {
	return "onManual"
}

// EventType on moduleBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (moduleBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on moduleBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (moduleBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on moduleBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (moduleBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on moduleAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (moduleAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on moduleAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (moduleAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on moduleAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (moduleAfterDelete) EventType() string {
	return "afterDelete"
}

// ModuleOnManual creates onManual for compose:module resource
//
// This function is auto-generated.
func ModuleOnManual(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleOnManual {
	return &moduleOnManual{
		moduleBase: &moduleBase{
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleBeforeCreate creates beforeCreate for compose:module resource
//
// This function is auto-generated.
func ModuleBeforeCreate(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleBeforeCreate {
	return &moduleBeforeCreate{
		moduleBase: &moduleBase{
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleBeforeUpdate creates beforeUpdate for compose:module resource
//
// This function is auto-generated.
func ModuleBeforeUpdate(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleBeforeUpdate {
	return &moduleBeforeUpdate{
		moduleBase: &moduleBase{
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleBeforeDelete creates beforeDelete for compose:module resource
//
// This function is auto-generated.
func ModuleBeforeDelete(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleBeforeDelete {
	return &moduleBeforeDelete{
		moduleBase: &moduleBase{
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleAfterCreate creates afterCreate for compose:module resource
//
// This function is auto-generated.
func ModuleAfterCreate(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleAfterCreate {
	return &moduleAfterCreate{
		moduleBase: &moduleBase{
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleAfterUpdate creates afterUpdate for compose:module resource
//
// This function is auto-generated.
func ModuleAfterUpdate(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleAfterUpdate {
	return &moduleAfterUpdate{
		moduleBase: &moduleBase{
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleAfterDelete creates afterDelete for compose:module resource
//
// This function is auto-generated.
func ModuleAfterDelete(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleAfterDelete {
	return &moduleAfterDelete{
		moduleBase: &moduleBase{
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// SetModule sets new module value
//
// This function is auto-generated.
func (res *moduleBase) SetModule(argModule *types.Module) {
	res.module = argModule
}

// Module returns module
//
// This function is auto-generated.
func (res moduleBase) Module() *types.Module {
	return res.module
}

// OldModule returns oldModule
//
// This function is auto-generated.
func (res moduleBase) OldModule() *types.Module {
	return res.oldModule
}

// Namespace returns namespace
//
// This function is auto-generated.
func (res moduleBase) Namespace() *types.Namespace {
	return res.namespace
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *moduleBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res moduleBase) Invoker() auth.Identifiable {
	return res.invoker
}
