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
