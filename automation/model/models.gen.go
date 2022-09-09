package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
)

var Session = &dal.Model{
	Ident:        "automation_sessions",
	ResourceType: types.SessionResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "WorkflowID", Sortable: true,
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::automation:workflow",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_workflow"},
		},

		&dal.Attribute{
			Ident: "Status", Sortable: true,
			Type: &dal.TypeNumber{HasDefault: true,
				DefaultValue: 0,
				Precision:    -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"},
			},
			Store: &dal.CodecAlias{Ident: "status"},
		},

		&dal.Attribute{
			Ident: "EventType", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "event_type"},
		},

		&dal.Attribute{
			Ident: "ResourceType", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "resource_type"},
		},

		&dal.Attribute{
			Ident: "Input",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "input"},
		},

		&dal.Attribute{
			Ident: "Output",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "output"},
		},

		&dal.Attribute{
			Ident: "Stacktrace",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "stacktrace"},
		},

		&dal.Attribute{
			Ident: "CreatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "created_by"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "PurgeAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "purge_at"},
		},

		&dal.Attribute{
			Ident: "SuspendedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "suspended_at"},
		},

		&dal.Attribute{
			Ident: "CompletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "completed_at"},
		},

		&dal.Attribute{
			Ident: "Error",
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "error"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "completedAt",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "CompletedAt",
				},
			},
		},

		&dal.Index{
			Ident: "createdAt",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "CreatedAt",
				},
			},
		},

		&dal.Index{
			Ident: "eventType",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "EventType",
				},
			},
		},

		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},

		&dal.Index{
			Ident: "resourceType",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ResourceType",
				},
			},
		},

		&dal.Index{
			Ident: "status",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "Status",
				},
			},
		},

		&dal.Index{
			Ident: "suspendedAt",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "SuspendedAt",
				},
			},
		},
	},
}

var Trigger = &dal.Model{
	Ident:        "automation_triggers",
	ResourceType: types.TriggerResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "WorkflowID", Sortable: true,
			Type: &dal.TypeRef{
				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::automation:workflow",
				},
			},
			Store: &dal.CodecAlias{Ident: "rel_workflow"},
		},

		&dal.Attribute{
			Ident: "StepID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "rel_step"},
		},

		&dal.Attribute{
			Ident: "Enabled", Sortable: true,
			Type: &dal.TypeBoolean{HasDefault: true,
				DefaultValue: true,
			},
			Store: &dal.CodecAlias{Ident: "enabled"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "ResourceType", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "resource_type"},
		},

		&dal.Attribute{
			Ident: "EventType", Sortable: true,
			Type:  &dal.TypeText{},
			Store: &dal.CodecAlias{Ident: "event_type"},
		},

		&dal.Attribute{
			Ident: "Constraints",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "constraints"},
		},

		&dal.Attribute{
			Ident: "Input",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "input"},
		},

		&dal.Attribute{
			Ident: "OwnedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "owned_by"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},

		&dal.Attribute{
			Ident: "CreatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "created_by"},
		},

		&dal.Attribute{
			Ident: "UpdatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "updated_by"},
		},

		&dal.Attribute{
			Ident: "DeletedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "deleted_by"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

var Workflow = &dal.Model{
	Ident:        "automation_workflows",
	ResourceType: types.WorkflowResourceType,

	Attributes: dal.AttributeSet{
		&dal.Attribute{
			Ident: "ID",
			Type:  &dal.TypeID{},
			Store: &dal.CodecAlias{Ident: "id"},
		},

		&dal.Attribute{
			Ident: "Handle",
			Type:  &dal.TypeText{Length: 64},
			Store: &dal.CodecAlias{Ident: "handle"},
		},

		&dal.Attribute{
			Ident: "Meta",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "meta"},
		},

		&dal.Attribute{
			Ident: "Enabled", Sortable: true,
			Type: &dal.TypeBoolean{HasDefault: true,
				DefaultValue: true,
			},
			Store: &dal.CodecAlias{Ident: "enabled"},
		},

		&dal.Attribute{
			Ident: "Trace",
			Type: &dal.TypeBoolean{HasDefault: true,
				DefaultValue: false,
			},
			Store: &dal.CodecAlias{Ident: "trace"},
		},

		&dal.Attribute{
			Ident: "KeepSessions",
			Type: &dal.TypeNumber{HasDefault: true,
				DefaultValue: 0,
				Precision:    -1, Scale: -1, Meta: map[string]interface{}{"rdbms:type": "integer"},
			},
			Store: &dal.CodecAlias{Ident: "keep_sessions"},
		},

		&dal.Attribute{
			Ident: "Scope",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "scope"},
		},

		&dal.Attribute{
			Ident: "Steps",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "steps"},
		},

		&dal.Attribute{
			Ident: "Paths",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "paths"},
		},

		&dal.Attribute{
			Ident: "Issues",
			Type: &dal.TypeJSON{
				DefaultValue: "{}",
			},
			Store: &dal.CodecAlias{Ident: "issues"},
		},

		&dal.Attribute{
			Ident: "RunAs",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "run_as"},
		},

		&dal.Attribute{
			Ident: "OwnedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "owned_by"},
		},

		&dal.Attribute{
			Ident: "CreatedAt", Sortable: true,
			Type: &dal.TypeTimestamp{
				DefaultCurrentTimestamp: true, Timezone: true, Precision: -1,
			},
			Store: &dal.CodecAlias{Ident: "created_at"},
		},

		&dal.Attribute{
			Ident: "UpdatedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "updated_at"},
		},

		&dal.Attribute{
			Ident: "DeletedAt", Sortable: true,
			Type:  &dal.TypeTimestamp{Nullable: true, Timezone: true, Precision: -1},
			Store: &dal.CodecAlias{Ident: "deleted_at"},
		},

		&dal.Attribute{
			Ident: "CreatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "created_by"},
		},

		&dal.Attribute{
			Ident: "UpdatedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "updated_by"},
		},

		&dal.Attribute{
			Ident: "DeletedBy",
			Type: &dal.TypeRef{HasDefault: true,
				DefaultValue: 0,

				RefAttribute: "id",
				RefModel: &dal.ModelRef{
					ResourceType: "corteza::system:user",
				},
			},
			Store: &dal.CodecAlias{Ident: "deleted_by"},
		},
	},

	Indexes: dal.IndexSet{
		&dal.Index{
			Ident: "PRIMARY",
			Type:  "BTREE",

			Fields: []*dal.IndexField{
				{
					AttributeIdent: "ID",
				},
			},
		},
	},
}

func init() {
	models = append(
		models,
		Session,
		Trigger,
		Workflow,
	)
}
