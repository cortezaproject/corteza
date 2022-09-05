package model

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
)

type (
	modelReplacer interface {
		ReplaceModel(ctx context.Context, model *dal.Model) (err error)
	}
)

var (
	Workflow = &dal.Model{
		Ident:        "automation_workflows",
		ResourceType: types.WorkflowResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident: "Handle",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "handle"},
			},

			&dal.Attribute{
				Ident: "Meta",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident:    "Enabled",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "enabled"},
			},

			&dal.Attribute{
				Ident: "Trace",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "trace"},
			},

			&dal.Attribute{
				Ident: "KeepSessions",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "keep_sessions"},
			},

			&dal.Attribute{
				Ident: "Scope",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "scope"},
			},

			&dal.Attribute{
				Ident: "Steps",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "steps"},
			},

			&dal.Attribute{
				Ident: "Paths",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "paths"},
			},

			&dal.Attribute{
				Ident: "Issues",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "issues"},
			},

			&dal.Attribute{
				Ident: "RunAs",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "run_as"},
			},

			&dal.Attribute{
				Ident:    "CreatedAt",
				Sortable: true,
				Type:     &dal.TypeTimestamp{},
				Store:    &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident:    "UpdatedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident:    "DeletedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},

			&dal.Attribute{
				Ident: "OwnedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "owned_by"},
			},

			&dal.Attribute{
				Ident: "CreatedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "created_by"},
			},

			&dal.Attribute{
				Ident: "UpdatedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "updated_by"},
			},

			&dal.Attribute{
				Ident: "DeletedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "deleted_by"},
			},
		},
	}

	Session = &dal.Model{
		Ident:        "automation_sessions",
		ResourceType: types.SessionResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident:    "WorkflowID",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "rel_workflow"},
			},

			&dal.Attribute{
				Ident:    "EventType",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "event_type"},
			},

			&dal.Attribute{
				Ident:    "ResourceType",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "resource_type"},
			},

			&dal.Attribute{
				Ident:    "Status",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "status"},
			},

			&dal.Attribute{
				Ident: "Input",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "input"},
			},

			&dal.Attribute{
				Ident: "Output",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "output"},
			},

			&dal.Attribute{
				Ident: "Stacktrace",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "stacktrace"},
			},

			&dal.Attribute{
				Ident: "CreatedBy",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "created_by"},
			},

			&dal.Attribute{
				Ident:    "CreatedAt",
				Sortable: true,
				Type:     &dal.TypeTimestamp{},
				Store:    &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident:    "PurgeAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "purge_at"},
			},

			&dal.Attribute{
				Ident:    "CompletedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "completed_at"},
			},

			&dal.Attribute{
				Ident:    "SuspendedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "suspended_at"},
			},

			&dal.Attribute{
				Ident: "Error",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "error"},
			},
		},
	}

	Trigger = &dal.Model{
		Ident:        "automation_triggers",
		ResourceType: types.TriggerResourceType,

		Attributes: dal.AttributeSet{
			&dal.Attribute{
				Ident:      "ID",
				PrimaryKey: true,
				Type:       &dal.TypeID{},
				Store:      &dal.CodecAlias{Ident: "id"},
			},

			&dal.Attribute{
				Ident:    "WorkflowID",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "rel_workflow"},
			},

			&dal.Attribute{
				Ident: "StepID",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "rel_step"},
			},

			&dal.Attribute{
				Ident:    "Enabled",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "enabled"},
			},

			&dal.Attribute{
				Ident:    "ResourceType",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "resource_type"},
			},

			&dal.Attribute{
				Ident:    "EventType",
				Sortable: true,
				Type:     &dal.TypeText{},
				Store:    &dal.CodecAlias{Ident: "event_type"},
			},

			&dal.Attribute{
				Ident: "Meta",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "meta"},
			},

			&dal.Attribute{
				Ident: "Constraints",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "constraints"},
			},

			&dal.Attribute{
				Ident: "Input",

				Type:  &dal.TypeText{},
				Store: &dal.CodecAlias{Ident: "input"},
			},

			&dal.Attribute{
				Ident:    "CreatedAt",
				Sortable: true,
				Type:     &dal.TypeTimestamp{},
				Store:    &dal.CodecAlias{Ident: "created_at"},
			},

			&dal.Attribute{
				Ident:    "UpdatedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "updated_at"},
			},

			&dal.Attribute{
				Ident:    "DeletedAt",
				Sortable: true,
				Type: &dal.TypeTimestamp{
					Nullable: true},
				Store: &dal.CodecAlias{Ident: "deleted_at"},
			},

			&dal.Attribute{
				Ident: "OwnedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "owned_by"},
			},

			&dal.Attribute{
				Ident: "CreatedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "created_by"},
			},

			&dal.Attribute{
				Ident: "UpdatedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "updated_by"},
			},

			&dal.Attribute{
				Ident: "DeletedBy",

				Type: &dal.TypeRef{

					RefAttribute: "id",
					RefModel: &dal.ModelRef{
						ResourceType: "corteza::system:user",
					},
				},
				Store: &dal.CodecAlias{Ident: "deleted_by"},
			},
		},
	}
)

func All() dal.ModelSet {
	return dal.ModelSet{
		Workflow,
		Session,
		Trigger,
	}
}

func Register(ctx context.Context, mr modelReplacer) (err error) {
	if err = mr.ReplaceModel(ctx, Workflow); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Session); err != nil {
		return
	}

	if err = mr.ReplaceModel(ctx, Trigger); err != nil {
		return
	}

	return
}
