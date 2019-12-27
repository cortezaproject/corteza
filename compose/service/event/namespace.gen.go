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
	// namespaceBase
	//
	// This type is auto-generated.
	namespaceBase struct {
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
)

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
