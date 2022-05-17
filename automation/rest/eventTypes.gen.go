package rest

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

func getEventTypeDefinitions() []eventTypeDef {
	return []eventTypeDef{

		{
			ResourceType: "compose",
			EventType:    "onManual",
			Properties:   []eventTypePropertyDef{},
			Constraints:  []eventTypeConstraintDef{},
		},

		{
			ResourceType: "compose",
			EventType:    "onInterval",
			Properties:   []eventTypePropertyDef{},
			Constraints:  []eventTypeConstraintDef{},
		},

		{
			ResourceType: "compose",
			EventType:    "onTimestamp",
			Properties:   []eventTypePropertyDef{},
			Constraints:  []eventTypeConstraintDef{},
		},

		{
			ResourceType: "compose:module",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: false,
				},

				{
					Name:      "oldModule",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},
			},
		},

		{
			ResourceType: "compose:module",
			EventType:    "beforeCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: false,
				},

				{
					Name:      "oldModule",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},
			},
		},

		{
			ResourceType: "compose:module",
			EventType:    "beforeUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: false,
				},

				{
					Name:      "oldModule",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},
			},
		},

		{
			ResourceType: "compose:module",
			EventType:    "beforeDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: false,
				},

				{
					Name:      "oldModule",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},
			},
		},

		{
			ResourceType: "compose:module",
			EventType:    "afterCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: false,
				},

				{
					Name:      "oldModule",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},
			},
		},

		{
			ResourceType: "compose:module",
			EventType:    "afterUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: false,
				},

				{
					Name:      "oldModule",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},
			},
		},

		{
			ResourceType: "compose:module",
			EventType:    "afterDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: false,
				},

				{
					Name:      "oldModule",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},
			},
		},

		{
			ResourceType: "compose:namespace",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: false,
				},

				{
					Name:      "oldNamespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},
			},
		},

		{
			ResourceType: "compose:namespace",
			EventType:    "beforeCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: false,
				},

				{
					Name:      "oldNamespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},
			},
		},

		{
			ResourceType: "compose:namespace",
			EventType:    "beforeUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: false,
				},

				{
					Name:      "oldNamespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},
			},
		},

		{
			ResourceType: "compose:namespace",
			EventType:    "beforeDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: false,
				},

				{
					Name:      "oldNamespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},
			},
		},

		{
			ResourceType: "compose:namespace",
			EventType:    "afterCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: false,
				},

				{
					Name:      "oldNamespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},
			},
		},

		{
			ResourceType: "compose:namespace",
			EventType:    "afterUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: false,
				},

				{
					Name:      "oldNamespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},
			},
		},

		{
			ResourceType: "compose:namespace",
			EventType:    "afterDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: false,
				},

				{
					Name:      "oldNamespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},
			},
		},

		{
			ResourceType: "compose:page",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "page",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldPage",
					Type:      "",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "page.handle",
				},

				{
					Name: "page.name",
				},
			},
		},

		{
			ResourceType: "compose:page",
			EventType:    "beforeCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "page",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldPage",
					Type:      "",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "page.handle",
				},

				{
					Name: "page.name",
				},
			},
		},

		{
			ResourceType: "compose:page",
			EventType:    "beforeUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "page",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldPage",
					Type:      "",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "page.handle",
				},

				{
					Name: "page.name",
				},
			},
		},

		{
			ResourceType: "compose:page",
			EventType:    "beforeDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "page",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldPage",
					Type:      "",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "page.handle",
				},

				{
					Name: "page.name",
				},
			},
		},

		{
			ResourceType: "compose:page",
			EventType:    "afterCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "page",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldPage",
					Type:      "",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "page.handle",
				},

				{
					Name: "page.name",
				},
			},
		},

		{
			ResourceType: "compose:page",
			EventType:    "afterUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "page",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldPage",
					Type:      "",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "page.handle",
				},

				{
					Name: "page.name",
				},
			},
		},

		{
			ResourceType: "compose:page",
			EventType:    "afterDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "page",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldPage",
					Type:      "",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "page.handle",
				},

				{
					Name: "page.name",
				},
			},
		},

		{
			ResourceType: "compose:record",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "record",
					Type:      "ComposeRecord",
					Immutable: false,
				},

				{
					Name:      "oldRecord",
					Type:      "ComposeRecord",
					Immutable: true,
				},

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "recordValueErrors",
					Type:      "ComposeRecordValueErrorSet",
					Immutable: false,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},

				{
					Name: "record.created-at",
				},

				{
					Name: "record.updated-at",
				},

				{
					Name: "record.deleted-at",
				},

				{
					Name: "record.values.*",
				},
			},
		},

		{
			ResourceType: "compose:record",
			EventType:    "onIteration",
			Properties: []eventTypePropertyDef{

				{
					Name:      "record",
					Type:      "ComposeRecord",
					Immutable: false,
				},

				{
					Name:      "oldRecord",
					Type:      "ComposeRecord",
					Immutable: true,
				},

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "recordValueErrors",
					Type:      "ComposeRecordValueErrorSet",
					Immutable: false,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},

				{
					Name: "record.created-at",
				},

				{
					Name: "record.updated-at",
				},

				{
					Name: "record.deleted-at",
				},

				{
					Name: "record.values.*",
				},
			},
		},

		{
			ResourceType: "compose:record",
			EventType:    "beforeCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "record",
					Type:      "ComposeRecord",
					Immutable: false,
				},

				{
					Name:      "oldRecord",
					Type:      "ComposeRecord",
					Immutable: true,
				},

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "recordValueErrors",
					Type:      "ComposeRecordValueErrorSet",
					Immutable: false,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},

				{
					Name: "record.created-at",
				},

				{
					Name: "record.updated-at",
				},

				{
					Name: "record.deleted-at",
				},

				{
					Name: "record.values.*",
				},
			},
		},

		{
			ResourceType: "compose:record",
			EventType:    "beforeUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "record",
					Type:      "ComposeRecord",
					Immutable: false,
				},

				{
					Name:      "oldRecord",
					Type:      "ComposeRecord",
					Immutable: true,
				},

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "recordValueErrors",
					Type:      "ComposeRecordValueErrorSet",
					Immutable: false,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},

				{
					Name: "record.created-at",
				},

				{
					Name: "record.updated-at",
				},

				{
					Name: "record.deleted-at",
				},

				{
					Name: "record.values.*",
				},
			},
		},

		{
			ResourceType: "compose:record",
			EventType:    "beforeDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "record",
					Type:      "ComposeRecord",
					Immutable: false,
				},

				{
					Name:      "oldRecord",
					Type:      "ComposeRecord",
					Immutable: true,
				},

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "recordValueErrors",
					Type:      "ComposeRecordValueErrorSet",
					Immutable: false,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},

				{
					Name: "record.created-at",
				},

				{
					Name: "record.updated-at",
				},

				{
					Name: "record.deleted-at",
				},

				{
					Name: "record.values.*",
				},
			},
		},

		{
			ResourceType: "compose:record",
			EventType:    "afterCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "record",
					Type:      "ComposeRecord",
					Immutable: false,
				},

				{
					Name:      "oldRecord",
					Type:      "ComposeRecord",
					Immutable: true,
				},

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "recordValueErrors",
					Type:      "ComposeRecordValueErrorSet",
					Immutable: false,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},

				{
					Name: "record.created-at",
				},

				{
					Name: "record.updated-at",
				},

				{
					Name: "record.deleted-at",
				},

				{
					Name: "record.values.*",
				},
			},
		},

		{
			ResourceType: "compose:record",
			EventType:    "afterUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "record",
					Type:      "ComposeRecord",
					Immutable: false,
				},

				{
					Name:      "oldRecord",
					Type:      "ComposeRecord",
					Immutable: true,
				},

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "recordValueErrors",
					Type:      "ComposeRecordValueErrorSet",
					Immutable: false,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},

				{
					Name: "record.created-at",
				},

				{
					Name: "record.updated-at",
				},

				{
					Name: "record.deleted-at",
				},

				{
					Name: "record.values.*",
				},
			},
		},

		{
			ResourceType: "compose:record",
			EventType:    "afterDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "record",
					Type:      "ComposeRecord",
					Immutable: false,
				},

				{
					Name:      "oldRecord",
					Type:      "ComposeRecord",
					Immutable: true,
				},

				{
					Name:      "module",
					Type:      "ComposeModule",
					Immutable: true,
				},

				{
					Name:      "namespace",
					Type:      "ComposeNamespace",
					Immutable: true,
				},

				{
					Name:      "recordValueErrors",
					Type:      "ComposeRecordValueErrorSet",
					Immutable: false,
				},

				{
					Name:      "selected",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "namespace.handle",
				},

				{
					Name: "namespace.name",
				},

				{
					Name: "module.handle",
				},

				{
					Name: "module.name",
				},

				{
					Name: "record.created-at",
				},

				{
					Name: "record.updated-at",
				},

				{
					Name: "record.deleted-at",
				},

				{
					Name: "record.values.*",
				},
			},
		},

		{
			ResourceType: "system",
			EventType:    "onManual",
			Properties:   []eventTypePropertyDef{},
			Constraints:  []eventTypeConstraintDef{},
		},

		{
			ResourceType: "system",
			EventType:    "onInterval",
			Properties:   []eventTypePropertyDef{},
			Constraints:  []eventTypeConstraintDef{},
		},

		{
			ResourceType: "system",
			EventType:    "onTimestamp",
			Properties:   []eventTypePropertyDef{},
			Constraints:  []eventTypeConstraintDef{},
		},

		{
			ResourceType: "system:application",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "application",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldApplication",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "application.name",
				},
			},
		},

		{
			ResourceType: "system:application",
			EventType:    "beforeCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "application",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldApplication",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "application.name",
				},
			},
		},

		{
			ResourceType: "system:application",
			EventType:    "beforeUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "application",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldApplication",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "application.name",
				},
			},
		},

		{
			ResourceType: "system:application",
			EventType:    "beforeDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "application",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldApplication",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "application.name",
				},
			},
		},

		{
			ResourceType: "system:application",
			EventType:    "afterCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "application",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldApplication",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "application.name",
				},
			},
		},

		{
			ResourceType: "system:application",
			EventType:    "afterUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "application",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldApplication",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "application.name",
				},
			},
		},

		{
			ResourceType: "system:application",
			EventType:    "afterDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "application",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldApplication",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "application.name",
				},
			},
		},

		{
			ResourceType: "system:auth",
			EventType:    "beforeLogin",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "provider",
					Type:      "",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:auth",
			EventType:    "beforeSignup",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "provider",
					Type:      "",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:auth",
			EventType:    "afterLogin",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "provider",
					Type:      "",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:auth",
			EventType:    "afterSignup",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "provider",
					Type:      "",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:auth-client",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "authClient",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldAuthClient",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "auth-client.handle",
				},
			},
		},

		{
			ResourceType: "system:auth-client",
			EventType:    "beforeCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "authClient",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldAuthClient",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "auth-client.handle",
				},
			},
		},

		{
			ResourceType: "system:auth-client",
			EventType:    "beforeUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "authClient",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldAuthClient",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "auth-client.handle",
				},
			},
		},

		{
			ResourceType: "system:auth-client",
			EventType:    "beforeDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "authClient",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldAuthClient",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "auth-client.handle",
				},
			},
		},

		{
			ResourceType: "system:auth-client",
			EventType:    "afterCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "authClient",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldAuthClient",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "auth-client.handle",
				},
			},
		},

		{
			ResourceType: "system:auth-client",
			EventType:    "afterUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "authClient",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldAuthClient",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "auth-client.handle",
				},
			},
		},

		{
			ResourceType: "system:auth-client",
			EventType:    "afterDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "authClient",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldAuthClient",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "auth-client.handle",
				},
			},
		},

		{
			ResourceType: "system:data-privacy-request",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "dataPrivacyRequest",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldDataPrivacyRequest",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:data-privacy-request",
			EventType:    "beforeCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "dataPrivacyRequest",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldDataPrivacyRequest",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:data-privacy-request",
			EventType:    "beforeUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "dataPrivacyRequest",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldDataPrivacyRequest",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:data-privacy-request",
			EventType:    "beforeDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "dataPrivacyRequest",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldDataPrivacyRequest",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:data-privacy-request",
			EventType:    "afterCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "dataPrivacyRequest",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldDataPrivacyRequest",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:data-privacy-request",
			EventType:    "afterUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "dataPrivacyRequest",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldDataPrivacyRequest",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:data-privacy-request",
			EventType:    "afterDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "dataPrivacyRequest",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "oldDataPrivacyRequest",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:mail",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "message",
					Type:      "",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "message.header.subject",
				},

				{
					Name: "message.header.from",
				},

				{
					Name: "message.header.to",
				},

				{
					Name: "message.header.reply-to",
				},

				{
					Name: "message.header.cc",
				},

				{
					Name: "message.header.bcc",
				},
			},
		},

		{
			ResourceType: "system:mail",
			EventType:    "onReceive",
			Properties: []eventTypePropertyDef{

				{
					Name:      "message",
					Type:      "",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "message.header.subject",
				},

				{
					Name: "message.header.from",
				},

				{
					Name: "message.header.to",
				},

				{
					Name: "message.header.reply-to",
				},

				{
					Name: "message.header.cc",
				},

				{
					Name: "message.header.bcc",
				},
			},
		},

		{
			ResourceType: "system:mail",
			EventType:    "onSend",
			Properties: []eventTypePropertyDef{

				{
					Name:      "message",
					Type:      "",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "message.header.subject",
				},

				{
					Name: "message.header.from",
				},

				{
					Name: "message.header.to",
				},

				{
					Name: "message.header.reply-to",
				},

				{
					Name: "message.header.cc",
				},

				{
					Name: "message.header.bcc",
				},
			},
		},

		{
			ResourceType: "system:queue",
			EventType:    "onMessage",
			Properties: []eventTypePropertyDef{

				{
					Name:      "payload",
					Type:      "QueueMessage",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "payload.queue",
				},
			},
		},

		{
			ResourceType: "system:role",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},

				{
					Name:      "oldRole",
					Type:      "Role",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:role",
			EventType:    "beforeCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},

				{
					Name:      "oldRole",
					Type:      "Role",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:role",
			EventType:    "beforeUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},

				{
					Name:      "oldRole",
					Type:      "Role",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:role",
			EventType:    "beforeDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},

				{
					Name:      "oldRole",
					Type:      "Role",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:role",
			EventType:    "afterCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},

				{
					Name:      "oldRole",
					Type:      "Role",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:role",
			EventType:    "afterUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},

				{
					Name:      "oldRole",
					Type:      "Role",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:role",
			EventType:    "afterDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},

				{
					Name:      "oldRole",
					Type:      "Role",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},
			},
		},

		{
			ResourceType: "system:role:member",
			EventType:    "beforeAdd",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:role:member",
			EventType:    "beforeRemove",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:role:member",
			EventType:    "afterAdd",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:role:member",
			EventType:    "afterRemove",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "role",
					Type:      "Role",
					Immutable: false,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "role.handle",
				},

				{
					Name: "role.name",
				},

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:sink",
			EventType:    "onRequest",
			Properties: []eventTypePropertyDef{

				{
					Name:      "response",
					Type:      "",
					Immutable: false,
				},

				{
					Name:      "request",
					Type:      "",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "request.host",
				},

				{
					Name: "request.remote-address",
				},

				{
					Name: "request.method",
				},

				{
					Name: "request.path",
				},

				{
					Name: "request.username",
				},

				{
					Name: "request.password",
				},

				{
					Name: "request.content-type",
				},

				{
					Name: "request.get.*",
				},

				{
					Name: "request.post.*",
				},

				{
					Name: "request.header.*",
				},
			},
		},

		{
			ResourceType: "system:user",
			EventType:    "onManual",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "oldUser",
					Type:      "User",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:user",
			EventType:    "beforeCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "oldUser",
					Type:      "User",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:user",
			EventType:    "beforeUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "oldUser",
					Type:      "User",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:user",
			EventType:    "beforeDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "oldUser",
					Type:      "User",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:user",
			EventType:    "beforeSuspend",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "oldUser",
					Type:      "User",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:user",
			EventType:    "afterCreate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "oldUser",
					Type:      "User",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:user",
			EventType:    "afterUpdate",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "oldUser",
					Type:      "User",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:user",
			EventType:    "afterDelete",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "oldUser",
					Type:      "User",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},

		{
			ResourceType: "system:user",
			EventType:    "afterSuspend",
			Properties: []eventTypePropertyDef{

				{
					Name:      "user",
					Type:      "User",
					Immutable: false,
				},

				{
					Name:      "oldUser",
					Type:      "User",
					Immutable: true,
				},
			},
			Constraints: []eventTypeConstraintDef{

				{
					Name: "user.handle",
				},

				{
					Name: "user.email",
				},
			},
		},
	}
}
