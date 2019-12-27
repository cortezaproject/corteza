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
	// recordBase
	//
	// This type is auto-generated.
	recordBase struct {
		record    *types.Record
		oldRecord *types.Record
		module    *types.Module
		namespace *types.Namespace
		invoker   auth.Identifiable
	}

	// recordOnManual
	//
	// This type is auto-generated.
	recordOnManual struct {
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
) *recordOnManual {
	return &recordOnManual{
		recordBase: &recordBase{
			record:    argRecord,
			oldRecord: argOldRecord,
			module:    argModule,
			namespace: argNamespace,
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
) *recordBeforeCreate {
	return &recordBeforeCreate{
		recordBase: &recordBase{
			record:    argRecord,
			oldRecord: argOldRecord,
			module:    argModule,
			namespace: argNamespace,
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
) *recordBeforeUpdate {
	return &recordBeforeUpdate{
		recordBase: &recordBase{
			record:    argRecord,
			oldRecord: argOldRecord,
			module:    argModule,
			namespace: argNamespace,
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
) *recordBeforeDelete {
	return &recordBeforeDelete{
		recordBase: &recordBase{
			record:    argRecord,
			oldRecord: argOldRecord,
			module:    argModule,
			namespace: argNamespace,
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
) *recordAfterCreate {
	return &recordAfterCreate{
		recordBase: &recordBase{
			record:    argRecord,
			oldRecord: argOldRecord,
			module:    argModule,
			namespace: argNamespace,
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
) *recordAfterUpdate {
	return &recordAfterUpdate{
		recordBase: &recordBase{
			record:    argRecord,
			oldRecord: argOldRecord,
			module:    argModule,
			namespace: argNamespace,
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
) *recordAfterDelete {
	return &recordAfterDelete{
		recordBase: &recordBase{
			record:    argRecord,
			oldRecord: argOldRecord,
			module:    argModule,
			namespace: argNamespace,
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
