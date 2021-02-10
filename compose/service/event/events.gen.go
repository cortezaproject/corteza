package event

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/service/event/events.yaml

import (
	"encoding/json"
	"github.com/cortezaproject/corteza-server/compose/automation"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/expr"
)

// dummy placing to simplify import generation logic
var _ = json.NewEncoder

type (

	// composeBase
	//
	// This type is auto-generated.
	composeBase struct {
		immutable bool
		invoker   auth.Identifiable
	}

	// composeOnManual
	//
	// This type is auto-generated.
	composeOnManual struct {
		*composeBase
	}

	// composeOnInterval
	//
	// This type is auto-generated.
	composeOnInterval struct {
		*composeBase
	}

	// composeOnTimestamp
	//
	// This type is auto-generated.
	composeOnTimestamp struct {
		*composeBase
	}

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

	// namespaceBase
	//
	// This type is auto-generated.
	namespaceBase struct {
		immutable    bool
		namespace    *types.Namespace
		oldNamespace *types.Namespace
		invoker      auth.Identifiable
	}

	// namespaceOnManual
	//
	// This type is auto-generated.
	namespaceOnManual struct {
		*namespaceBase
	}

	// namespaceBeforeCreate
	//
	// This type is auto-generated.
	namespaceBeforeCreate struct {
		*namespaceBase
	}

	// namespaceBeforeUpdate
	//
	// This type is auto-generated.
	namespaceBeforeUpdate struct {
		*namespaceBase
	}

	// namespaceBeforeDelete
	//
	// This type is auto-generated.
	namespaceBeforeDelete struct {
		*namespaceBase
	}

	// namespaceAfterCreate
	//
	// This type is auto-generated.
	namespaceAfterCreate struct {
		*namespaceBase
	}

	// namespaceAfterUpdate
	//
	// This type is auto-generated.
	namespaceAfterUpdate struct {
		*namespaceBase
	}

	// namespaceAfterDelete
	//
	// This type is auto-generated.
	namespaceAfterDelete struct {
		*namespaceBase
	}

	// pageBase
	//
	// This type is auto-generated.
	pageBase struct {
		immutable bool
		page      *types.Page
		oldPage   *types.Page
		namespace *types.Namespace
		invoker   auth.Identifiable
	}

	// pageOnManual
	//
	// This type is auto-generated.
	pageOnManual struct {
		*pageBase
	}

	// pageBeforeCreate
	//
	// This type is auto-generated.
	pageBeforeCreate struct {
		*pageBase
	}

	// pageBeforeUpdate
	//
	// This type is auto-generated.
	pageBeforeUpdate struct {
		*pageBase
	}

	// pageBeforeDelete
	//
	// This type is auto-generated.
	pageBeforeDelete struct {
		*pageBase
	}

	// pageAfterCreate
	//
	// This type is auto-generated.
	pageAfterCreate struct {
		*pageBase
	}

	// pageAfterUpdate
	//
	// This type is auto-generated.
	pageAfterUpdate struct {
		*pageBase
	}

	// pageAfterDelete
	//
	// This type is auto-generated.
	pageAfterDelete struct {
		*pageBase
	}

	// recordBase
	//
	// This type is auto-generated.
	recordBase struct {
		immutable         bool
		record            *types.Record
		oldRecord         *types.Record
		module            *types.Module
		namespace         *types.Namespace
		recordValueErrors *types.RecordValueErrorSet
		invoker           auth.Identifiable
	}

	// recordOnManual
	//
	// This type is auto-generated.
	recordOnManual struct {
		*recordBase
	}

	// recordOnIteration
	//
	// This type is auto-generated.
	recordOnIteration struct {
		*recordBase
	}

	// recordBeforeCreate
	//
	// This type is auto-generated.
	recordBeforeCreate struct {
		*recordBase
	}

	// recordBeforeUpdate
	//
	// This type is auto-generated.
	recordBeforeUpdate struct {
		*recordBase
	}

	// recordBeforeDelete
	//
	// This type is auto-generated.
	recordBeforeDelete struct {
		*recordBase
	}

	// recordAfterCreate
	//
	// This type is auto-generated.
	recordAfterCreate struct {
		*recordBase
	}

	// recordAfterUpdate
	//
	// This type is auto-generated.
	recordAfterUpdate struct {
		*recordBase
	}

	// recordAfterDelete
	//
	// This type is auto-generated.
	recordAfterDelete struct {
		*recordBase
	}
)

// ResourceType returns "compose"
//
// This function is auto-generated.
func (composeBase) ResourceType() string {
	return "compose"
}

// EventType on composeOnManual returns "onManual"
//
// This function is auto-generated.
func (composeOnManual) EventType() string {
	return "onManual"
}

// EventType on composeOnInterval returns "onInterval"
//
// This function is auto-generated.
func (composeOnInterval) EventType() string {
	return "onInterval"
}

// EventType on composeOnTimestamp returns "onTimestamp"
//
// This function is auto-generated.
func (composeOnTimestamp) EventType() string {
	return "onTimestamp"
}

// ComposeOnManual creates onManual for compose resource
//
// This function is auto-generated.
func ComposeOnManual() *composeOnManual {
	return &composeOnManual{
		composeBase: &composeBase{
			immutable: false,
		},
	}
}

// ComposeOnManualImmutable creates onManual for compose resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ComposeOnManualImmutable() *composeOnManual {
	return &composeOnManual{
		composeBase: &composeBase{
			immutable: true,
		},
	}
}

// ComposeOnInterval creates onInterval for compose resource
//
// This function is auto-generated.
func ComposeOnInterval() *composeOnInterval {
	return &composeOnInterval{
		composeBase: &composeBase{
			immutable: false,
		},
	}
}

// ComposeOnIntervalImmutable creates onInterval for compose resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ComposeOnIntervalImmutable() *composeOnInterval {
	return &composeOnInterval{
		composeBase: &composeBase{
			immutable: true,
		},
	}
}

// ComposeOnTimestamp creates onTimestamp for compose resource
//
// This function is auto-generated.
func ComposeOnTimestamp() *composeOnTimestamp {
	return &composeOnTimestamp{
		composeBase: &composeBase{
			immutable: false,
		},
	}
}

// ComposeOnTimestampImmutable creates onTimestamp for compose resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func ComposeOnTimestampImmutable() *composeOnTimestamp {
	return &composeOnTimestamp{
		composeBase: &composeBase{
			immutable: true,
		},
	}
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *composeBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res composeBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res composeBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res composeBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *composeBase) Decode(results map[string][]byte) (err error) {
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

func (res *composeBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

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

// Encode internal data to be passed as event params & arguments to workflow
func (res moduleBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	if rvars["module"], err = automation.NewComposeModule(res.module); err != nil {
		return nil, err
	}

	if rvars["oldModule"], err = automation.NewComposeModule(res.oldModule); err != nil {
		return nil, err
	}

	if rvars["namespace"], err = automation.NewComposeNamespace(res.namespace); err != nil {
		return nil, err
	}

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
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

	// Do not decode oldModule; marked as immutable

	// Do not decode namespace; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}

func (res *moduleBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.module != nil && vars.Has("module") {
		var aux *automation.ComposeModule
		aux, err = automation.NewComposeModule(expr.Must(vars.Select("module")))
		if err != nil {
			return
		}

		res.module = aux.GetValue()
	}
	// oldModule marked as immutable
	// namespace marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

// ResourceType returns "compose:namespace"
//
// This function is auto-generated.
func (namespaceBase) ResourceType() string {
	return "compose:namespace"
}

// EventType on namespaceOnManual returns "onManual"
//
// This function is auto-generated.
func (namespaceOnManual) EventType() string {
	return "onManual"
}

// EventType on namespaceBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (namespaceBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on namespaceBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (namespaceBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on namespaceBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (namespaceBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on namespaceAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (namespaceAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on namespaceAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (namespaceAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on namespaceAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (namespaceAfterDelete) EventType() string {
	return "afterDelete"
}

// NamespaceOnManual creates onManual for compose:namespace resource
//
// This function is auto-generated.
func NamespaceOnManual(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceOnManual {
	return &namespaceOnManual{
		namespaceBase: &namespaceBase{
			immutable:    false,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceOnManualImmutable creates onManual for compose:namespace resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func NamespaceOnManualImmutable(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceOnManual {
	return &namespaceOnManual{
		namespaceBase: &namespaceBase{
			immutable:    true,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceBeforeCreate creates beforeCreate for compose:namespace resource
//
// This function is auto-generated.
func NamespaceBeforeCreate(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceBeforeCreate {
	return &namespaceBeforeCreate{
		namespaceBase: &namespaceBase{
			immutable:    false,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceBeforeCreateImmutable creates beforeCreate for compose:namespace resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func NamespaceBeforeCreateImmutable(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceBeforeCreate {
	return &namespaceBeforeCreate{
		namespaceBase: &namespaceBase{
			immutable:    true,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceBeforeUpdate creates beforeUpdate for compose:namespace resource
//
// This function is auto-generated.
func NamespaceBeforeUpdate(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceBeforeUpdate {
	return &namespaceBeforeUpdate{
		namespaceBase: &namespaceBase{
			immutable:    false,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceBeforeUpdateImmutable creates beforeUpdate for compose:namespace resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func NamespaceBeforeUpdateImmutable(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceBeforeUpdate {
	return &namespaceBeforeUpdate{
		namespaceBase: &namespaceBase{
			immutable:    true,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceBeforeDelete creates beforeDelete for compose:namespace resource
//
// This function is auto-generated.
func NamespaceBeforeDelete(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceBeforeDelete {
	return &namespaceBeforeDelete{
		namespaceBase: &namespaceBase{
			immutable:    false,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceBeforeDeleteImmutable creates beforeDelete for compose:namespace resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func NamespaceBeforeDeleteImmutable(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceBeforeDelete {
	return &namespaceBeforeDelete{
		namespaceBase: &namespaceBase{
			immutable:    true,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceAfterCreate creates afterCreate for compose:namespace resource
//
// This function is auto-generated.
func NamespaceAfterCreate(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceAfterCreate {
	return &namespaceAfterCreate{
		namespaceBase: &namespaceBase{
			immutable:    false,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceAfterCreateImmutable creates afterCreate for compose:namespace resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func NamespaceAfterCreateImmutable(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceAfterCreate {
	return &namespaceAfterCreate{
		namespaceBase: &namespaceBase{
			immutable:    true,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceAfterUpdate creates afterUpdate for compose:namespace resource
//
// This function is auto-generated.
func NamespaceAfterUpdate(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceAfterUpdate {
	return &namespaceAfterUpdate{
		namespaceBase: &namespaceBase{
			immutable:    false,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceAfterUpdateImmutable creates afterUpdate for compose:namespace resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func NamespaceAfterUpdateImmutable(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceAfterUpdate {
	return &namespaceAfterUpdate{
		namespaceBase: &namespaceBase{
			immutable:    true,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceAfterDelete creates afterDelete for compose:namespace resource
//
// This function is auto-generated.
func NamespaceAfterDelete(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceAfterDelete {
	return &namespaceAfterDelete{
		namespaceBase: &namespaceBase{
			immutable:    false,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// NamespaceAfterDeleteImmutable creates afterDelete for compose:namespace resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func NamespaceAfterDeleteImmutable(
	argNamespace *types.Namespace,
	argOldNamespace *types.Namespace,
) *namespaceAfterDelete {
	return &namespaceAfterDelete{
		namespaceBase: &namespaceBase{
			immutable:    true,
			namespace:    argNamespace,
			oldNamespace: argOldNamespace,
		},
	}
}

// SetNamespace sets new namespace value
//
// This function is auto-generated.
func (res *namespaceBase) SetNamespace(argNamespace *types.Namespace) {
	res.namespace = argNamespace
}

// Namespace returns namespace
//
// This function is auto-generated.
func (res namespaceBase) Namespace() *types.Namespace {
	return res.namespace
}

// OldNamespace returns oldNamespace
//
// This function is auto-generated.
func (res namespaceBase) OldNamespace() *types.Namespace {
	return res.oldNamespace
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *namespaceBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res namespaceBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res namespaceBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["namespace"], err = json.Marshal(res.namespace); err != nil {
		return nil, err
	}

	if args["oldNamespace"], err = json.Marshal(res.oldNamespace); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res namespaceBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	if rvars["namespace"], err = automation.NewComposeNamespace(res.namespace); err != nil {
		return nil, err
	}

	if rvars["oldNamespace"], err = automation.NewComposeNamespace(res.oldNamespace); err != nil {
		return nil, err
	}

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *namespaceBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.namespace != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.namespace); err != nil {
				return
			}
		}
	}

	if res.namespace != nil {
		if r, ok := results["namespace"]; ok {
			if err = json.Unmarshal(r, res.namespace); err != nil {
				return
			}
		}
	}

	// Do not decode oldNamespace; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}

func (res *namespaceBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.namespace != nil && vars.Has("namespace") {
		var aux *automation.ComposeNamespace
		aux, err = automation.NewComposeNamespace(expr.Must(vars.Select("namespace")))
		if err != nil {
			return
		}

		res.namespace = aux.GetValue()
	}
	// oldNamespace marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

// ResourceType returns "compose:page"
//
// This function is auto-generated.
func (pageBase) ResourceType() string {
	return "compose:page"
}

// EventType on pageOnManual returns "onManual"
//
// This function is auto-generated.
func (pageOnManual) EventType() string {
	return "onManual"
}

// EventType on pageBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (pageBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on pageBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (pageBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on pageBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (pageBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on pageAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (pageAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on pageAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (pageAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on pageAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (pageAfterDelete) EventType() string {
	return "afterDelete"
}

// PageOnManual creates onManual for compose:page resource
//
// This function is auto-generated.
func PageOnManual(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageOnManual {
	return &pageOnManual{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageOnManualImmutable creates onManual for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageOnManualImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageOnManual {
	return &pageOnManual{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeCreate creates beforeCreate for compose:page resource
//
// This function is auto-generated.
func PageBeforeCreate(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeCreate {
	return &pageBeforeCreate{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeCreateImmutable creates beforeCreate for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageBeforeCreateImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeCreate {
	return &pageBeforeCreate{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeUpdate creates beforeUpdate for compose:page resource
//
// This function is auto-generated.
func PageBeforeUpdate(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeUpdate {
	return &pageBeforeUpdate{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeUpdateImmutable creates beforeUpdate for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageBeforeUpdateImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeUpdate {
	return &pageBeforeUpdate{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeDelete creates beforeDelete for compose:page resource
//
// This function is auto-generated.
func PageBeforeDelete(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeDelete {
	return &pageBeforeDelete{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageBeforeDeleteImmutable creates beforeDelete for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageBeforeDeleteImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageBeforeDelete {
	return &pageBeforeDelete{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterCreate creates afterCreate for compose:page resource
//
// This function is auto-generated.
func PageAfterCreate(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterCreate {
	return &pageAfterCreate{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterCreateImmutable creates afterCreate for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageAfterCreateImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterCreate {
	return &pageAfterCreate{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterUpdate creates afterUpdate for compose:page resource
//
// This function is auto-generated.
func PageAfterUpdate(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterUpdate {
	return &pageAfterUpdate{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterUpdateImmutable creates afterUpdate for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageAfterUpdateImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterUpdate {
	return &pageAfterUpdate{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterDelete creates afterDelete for compose:page resource
//
// This function is auto-generated.
func PageAfterDelete(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterDelete {
	return &pageAfterDelete{
		pageBase: &pageBase{
			immutable: false,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// PageAfterDeleteImmutable creates afterDelete for compose:page resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func PageAfterDeleteImmutable(
	argPage *types.Page,
	argOldPage *types.Page,
	argNamespace *types.Namespace,
) *pageAfterDelete {
	return &pageAfterDelete{
		pageBase: &pageBase{
			immutable: true,
			page:      argPage,
			oldPage:   argOldPage,
			namespace: argNamespace,
		},
	}
}

// SetPage sets new page value
//
// This function is auto-generated.
func (res *pageBase) SetPage(argPage *types.Page) {
	res.page = argPage
}

// Page returns page
//
// This function is auto-generated.
func (res pageBase) Page() *types.Page {
	return res.page
}

// OldPage returns oldPage
//
// This function is auto-generated.
func (res pageBase) OldPage() *types.Page {
	return res.oldPage
}

// Namespace returns namespace
//
// This function is auto-generated.
func (res pageBase) Namespace() *types.Namespace {
	return res.namespace
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *pageBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res pageBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res pageBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["page"], err = json.Marshal(res.page); err != nil {
		return nil, err
	}

	if args["oldPage"], err = json.Marshal(res.oldPage); err != nil {
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

// Encode internal data to be passed as event params & arguments to workflow
func (res pageBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	// Could not found expression-type counterpart for *types.Page

	// Could not found expression-type counterpart for *types.Page

	if rvars["namespace"], err = automation.NewComposeNamespace(res.namespace); err != nil {
		return nil, err
	}

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *pageBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.page != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.page); err != nil {
				return
			}
		}
	}

	if res.page != nil {
		if r, ok := results["page"]; ok {
			if err = json.Unmarshal(r, res.page); err != nil {
				return
			}
		}
	}

	// Do not decode oldPage; marked as immutable

	// Do not decode namespace; marked as immutable

	if res.invoker != nil {
		if r, ok := results["invoker"]; ok {
			if err = json.Unmarshal(r, res.invoker); err != nil {
				return
			}
		}
	}
	return
}

func (res *pageBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	// Could not find expression-type counterpart for *types.Page
	// oldPage marked as immutable
	// namespace marked as immutable
	// Could not find expression-type counterpart for auth.Identifiable

	return
}

// ResourceType returns "compose:record"
//
// This function is auto-generated.
func (recordBase) ResourceType() string {
	return "compose:record"
}

// EventType on recordOnManual returns "onManual"
//
// This function is auto-generated.
func (recordOnManual) EventType() string {
	return "onManual"
}

// EventType on recordOnIteration returns "onIteration"
//
// This function is auto-generated.
func (recordOnIteration) EventType() string {
	return "onIteration"
}

// EventType on recordBeforeCreate returns "beforeCreate"
//
// This function is auto-generated.
func (recordBeforeCreate) EventType() string {
	return "beforeCreate"
}

// EventType on recordBeforeUpdate returns "beforeUpdate"
//
// This function is auto-generated.
func (recordBeforeUpdate) EventType() string {
	return "beforeUpdate"
}

// EventType on recordBeforeDelete returns "beforeDelete"
//
// This function is auto-generated.
func (recordBeforeDelete) EventType() string {
	return "beforeDelete"
}

// EventType on recordAfterCreate returns "afterCreate"
//
// This function is auto-generated.
func (recordAfterCreate) EventType() string {
	return "afterCreate"
}

// EventType on recordAfterUpdate returns "afterUpdate"
//
// This function is auto-generated.
func (recordAfterUpdate) EventType() string {
	return "afterUpdate"
}

// EventType on recordAfterDelete returns "afterDelete"
//
// This function is auto-generated.
func (recordAfterDelete) EventType() string {
	return "afterDelete"
}

// RecordOnManual creates onManual for compose:record resource
//
// This function is auto-generated.
func RecordOnManual(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordOnManual {
	return &recordOnManual{
		recordBase: &recordBase{
			immutable:         false,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordOnManualImmutable creates onManual for compose:record resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RecordOnManualImmutable(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordOnManual {
	return &recordOnManual{
		recordBase: &recordBase{
			immutable:         true,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordOnIteration creates onIteration for compose:record resource
//
// This function is auto-generated.
func RecordOnIteration(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordOnIteration {
	return &recordOnIteration{
		recordBase: &recordBase{
			immutable:         false,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordOnIterationImmutable creates onIteration for compose:record resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RecordOnIterationImmutable(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordOnIteration {
	return &recordOnIteration{
		recordBase: &recordBase{
			immutable:         true,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordBeforeCreate creates beforeCreate for compose:record resource
//
// This function is auto-generated.
func RecordBeforeCreate(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordBeforeCreate {
	return &recordBeforeCreate{
		recordBase: &recordBase{
			immutable:         false,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordBeforeCreateImmutable creates beforeCreate for compose:record resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RecordBeforeCreateImmutable(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordBeforeCreate {
	return &recordBeforeCreate{
		recordBase: &recordBase{
			immutable:         true,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordBeforeUpdate creates beforeUpdate for compose:record resource
//
// This function is auto-generated.
func RecordBeforeUpdate(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordBeforeUpdate {
	return &recordBeforeUpdate{
		recordBase: &recordBase{
			immutable:         false,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordBeforeUpdateImmutable creates beforeUpdate for compose:record resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RecordBeforeUpdateImmutable(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordBeforeUpdate {
	return &recordBeforeUpdate{
		recordBase: &recordBase{
			immutable:         true,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordBeforeDelete creates beforeDelete for compose:record resource
//
// This function is auto-generated.
func RecordBeforeDelete(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordBeforeDelete {
	return &recordBeforeDelete{
		recordBase: &recordBase{
			immutable:         false,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordBeforeDeleteImmutable creates beforeDelete for compose:record resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RecordBeforeDeleteImmutable(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordBeforeDelete {
	return &recordBeforeDelete{
		recordBase: &recordBase{
			immutable:         true,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordAfterCreate creates afterCreate for compose:record resource
//
// This function is auto-generated.
func RecordAfterCreate(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordAfterCreate {
	return &recordAfterCreate{
		recordBase: &recordBase{
			immutable:         false,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordAfterCreateImmutable creates afterCreate for compose:record resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RecordAfterCreateImmutable(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordAfterCreate {
	return &recordAfterCreate{
		recordBase: &recordBase{
			immutable:         true,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordAfterUpdate creates afterUpdate for compose:record resource
//
// This function is auto-generated.
func RecordAfterUpdate(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordAfterUpdate {
	return &recordAfterUpdate{
		recordBase: &recordBase{
			immutable:         false,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordAfterUpdateImmutable creates afterUpdate for compose:record resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RecordAfterUpdateImmutable(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordAfterUpdate {
	return &recordAfterUpdate{
		recordBase: &recordBase{
			immutable:         true,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordAfterDelete creates afterDelete for compose:record resource
//
// This function is auto-generated.
func RecordAfterDelete(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordAfterDelete {
	return &recordAfterDelete{
		recordBase: &recordBase{
			immutable:         false,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// RecordAfterDeleteImmutable creates afterDelete for compose:record resource
//
// None of the arguments will be mutable!
//
// This function is auto-generated.
func RecordAfterDeleteImmutable(
	argRecord *types.Record,
	argOldRecord *types.Record,
	argModule *types.Module,
	argNamespace *types.Namespace,
	argRecordValueErrors *types.RecordValueErrorSet,
) *recordAfterDelete {
	return &recordAfterDelete{
		recordBase: &recordBase{
			immutable:         true,
			record:            argRecord,
			oldRecord:         argOldRecord,
			module:            argModule,
			namespace:         argNamespace,
			recordValueErrors: argRecordValueErrors,
		},
	}
}

// SetRecord sets new record value
//
// This function is auto-generated.
func (res *recordBase) SetRecord(argRecord *types.Record) {
	res.record = argRecord
}

// Record returns record
//
// This function is auto-generated.
func (res recordBase) Record() *types.Record {
	return res.record
}

// OldRecord returns oldRecord
//
// This function is auto-generated.
func (res recordBase) OldRecord() *types.Record {
	return res.oldRecord
}

// Module returns module
//
// This function is auto-generated.
func (res recordBase) Module() *types.Module {
	return res.module
}

// Namespace returns namespace
//
// This function is auto-generated.
func (res recordBase) Namespace() *types.Namespace {
	return res.namespace
}

// SetRecordValueErrors sets new recordValueErrors value
//
// This function is auto-generated.
func (res *recordBase) SetRecordValueErrors(argRecordValueErrors *types.RecordValueErrorSet) {
	res.recordValueErrors = argRecordValueErrors
}

// RecordValueErrors returns recordValueErrors
//
// This function is auto-generated.
func (res recordBase) RecordValueErrors() *types.RecordValueErrorSet {
	return res.recordValueErrors
}

// SetInvoker sets new invoker value
//
// This function is auto-generated.
func (res *recordBase) SetInvoker(argInvoker auth.Identifiable) {
	res.invoker = argInvoker
}

// Invoker returns invoker
//
// This function is auto-generated.
func (res recordBase) Invoker() auth.Identifiable {
	return res.invoker
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res recordBase) Encode() (args map[string][]byte, err error) {
	args = make(map[string][]byte)

	if args["record"], err = json.Marshal(res.record); err != nil {
		return nil, err
	}

	if args["oldRecord"], err = json.Marshal(res.oldRecord); err != nil {
		return nil, err
	}

	if args["module"], err = json.Marshal(res.module); err != nil {
		return nil, err
	}

	if args["namespace"], err = json.Marshal(res.namespace); err != nil {
		return nil, err
	}

	if args["recordValueErrors"], err = json.Marshal(res.recordValueErrors); err != nil {
		return nil, err
	}

	if args["invoker"], err = json.Marshal(res.invoker); err != nil {
		return nil, err
	}

	return
}

// Encode internal data to be passed as event params & arguments to workflow
func (res recordBase) EncodeVars() (vars *expr.Vars, err error) {
	var (
		rvars = expr.RVars{}
	)

	if rvars["record"], err = automation.NewComposeRecord(res.record); err != nil {
		return nil, err
	}

	if rvars["oldRecord"], err = automation.NewComposeRecord(res.oldRecord); err != nil {
		return nil, err
	}

	if rvars["module"], err = automation.NewComposeModule(res.module); err != nil {
		return nil, err
	}

	if rvars["namespace"], err = automation.NewComposeNamespace(res.namespace); err != nil {
		return nil, err
	}

	if rvars["recordValueErrors"], err = automation.NewComposeRecordValueErrorSet(res.recordValueErrors); err != nil {
		return nil, err
	}

	// Could not found expression-type counterpart for auth.Identifiable

	return rvars.Vars(), err
}

// Decode return values from Corredor script into struct props
func (res *recordBase) Decode(results map[string][]byte) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.record != nil {
		if r, ok := results["result"]; ok && len(results) == 1 {
			if err = json.Unmarshal(r, res.record); err != nil {
				return
			}
		}
	}

	if res.record != nil {
		if r, ok := results["record"]; ok {
			if err = json.Unmarshal(r, res.record); err != nil {
				return
			}
		}
	}

	// Do not decode oldRecord; marked as immutable

	// Do not decode module; marked as immutable

	// Do not decode namespace; marked as immutable

	if res.recordValueErrors != nil {
		if r, ok := results["recordValueErrors"]; ok {
			if err = json.Unmarshal(r, res.recordValueErrors); err != nil {
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

func (res *recordBase) DecodeVars(vars *expr.Vars) (err error) {
	if res.immutable {
		// Respect immutability
		return
	}
	if res.record != nil && vars.Has("record") {
		var aux *automation.ComposeRecord
		aux, err = automation.NewComposeRecord(expr.Must(vars.Select("record")))
		if err != nil {
			return
		}

		res.record = aux.GetValue()
	}
	// oldRecord marked as immutable
	// module marked as immutable
	// namespace marked as immutable
	if res.recordValueErrors != nil && vars.Has("recordValueErrors") {
		var aux *automation.ComposeRecordValueErrorSet
		aux, err = automation.NewComposeRecordValueErrorSet(expr.Must(vars.Select("recordValueErrors")))
		if err != nil {
			return
		}

		res.recordValueErrors = aux.GetValue()
	}
	// Could not find expression-type counterpart for auth.Identifiable

	return
}
