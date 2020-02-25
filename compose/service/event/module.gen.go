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
	"encoding/json"

	"github.com/cortezaproject/corteza-server/compose/types"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

type (
	// moduleBase
	//
	// This type is auto-generated.
	moduleBase struct {
		immutable bool
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
			immutable: false,
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleOnManualImmutable creates onManual for compose:module resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ModuleOnManualImmutable(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleOnManual {
	return &moduleOnManual{
		moduleBase: &moduleBase{
			immutable: true,
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
			immutable: false,
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleBeforeCreateImmutable creates beforeCreate for compose:module resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ModuleBeforeCreateImmutable(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleBeforeCreate {
	return &moduleBeforeCreate{
		moduleBase: &moduleBase{
			immutable: true,
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
			immutable: false,
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleBeforeUpdateImmutable creates beforeUpdate for compose:module resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ModuleBeforeUpdateImmutable(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleBeforeUpdate {
	return &moduleBeforeUpdate{
		moduleBase: &moduleBase{
			immutable: true,
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
			immutable: false,
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleBeforeDeleteImmutable creates beforeDelete for compose:module resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ModuleBeforeDeleteImmutable(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleBeforeDelete {
	return &moduleBeforeDelete{
		moduleBase: &moduleBase{
			immutable: true,
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
			immutable: false,
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleAfterCreateImmutable creates afterCreate for compose:module resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ModuleAfterCreateImmutable(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleAfterCreate {
	return &moduleAfterCreate{
		moduleBase: &moduleBase{
			immutable: true,
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
			immutable: false,
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleAfterUpdateImmutable creates afterUpdate for compose:module resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ModuleAfterUpdateImmutable(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleAfterUpdate {
	return &moduleAfterUpdate{
		moduleBase: &moduleBase{
			immutable: true,
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
			immutable: false,
			module:    argModule,
			oldModule: argOldModule,
			namespace: argNamespace,
		},
	}
}

// ModuleAfterDeleteImmutable creates afterDelete for compose:module resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ModuleAfterDeleteImmutable(
	argModule *types.Module,
	argOldModule *types.Module,
	argNamespace *types.Namespace,
) *moduleAfterDelete {
	return &moduleAfterDelete{
		moduleBase: &moduleBase{
			immutable: true,
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

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res moduleBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["module"], err = json.Marshal(res.module); err != nil {
		return nil, err
	}

	if args["oldModule"], err = json.Marshal(res.oldModule); err != nil {
		return nil, err
	}

	if args["namespace"], err = json.Marshal(res.namespace); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Decode return values from Corredor script into struct props
func (res *moduleBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.module != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.module); err != nil {
				return
			}
		}
	}

	if res.module != nil {
		if r, ok := results["module"]; ok {
			if err = json.Unmarshal(r, res.module); err != nil {
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
