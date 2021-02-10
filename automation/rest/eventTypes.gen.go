package rest

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// <no value>

func getEventTypeDefinitions() []eventTypeDef {
	return []eventTypeDef{

		// []string{"github.com/cortezaproject/corteza-server/compose/types", "github.com/cortezaproject/corteza-server/compose/automation", "github.com/cortezaproject/corteza-server/pkg/auth"}

		{
			ResourceType: "compose",
			EventType:    "onManual",
			Properties:   []eventTypePropertyDef{},
		},

		{
			ResourceType: "compose",
			EventType:    "onInterval",
			Properties:   []eventTypePropertyDef{},
		},

		{
			ResourceType: "compose",
			EventType:    "onTimestamp",
			Properties:   []eventTypePropertyDef{},
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
			},
		},

		// []string{"github.com/cortezaproject/corteza-server/system/types", "github.com/cortezaproject/corteza-server/system/automation", "github.com/cortezaproject/corteza-server/pkg/auth"}

		{
			ResourceType: "system",
			EventType:    "onManual",
			Properties:   []eventTypePropertyDef{},
		},

		{
			ResourceType: "system",
			EventType:    "onInterval",
			Properties:   []eventTypePropertyDef{},
		},

		{
			ResourceType: "system",
			EventType:    "onTimestamp",
			Properties:   []eventTypePropertyDef{},
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
		},
	}
}
