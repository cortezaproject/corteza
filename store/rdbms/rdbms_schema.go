package rdbms

import (
	"context"
	"fmt"

	. "github.com/cortezaproject/corteza-server/store/rdbms/ddl"
)

type (
	Schema struct{}

	// schemaUpgrader provides procedures to upgrade rdbms store tables
	schemaUpgrader interface {
		Before(context.Context) error
		CreateTable(context.Context, *Table) error
		After(context.Context) error
	}
)

const (
	handleLength   = 64
	resourceLength = 512

	// https://www.rfc-editor.org/errata/eid1690
	emailLength = 254
)

// Upgrades schema
func (s *Schema) Upgrade(ctx context.Context, g schemaUpgrader) (err error) {
	if err = g.Before(ctx); err != nil {
		return fmt.Errorf("could not run \"before\" upgrade procedures: %w", err)
	}

	for _, t := range s.Tables() {
		if err = g.CreateTable(ctx, t); err != nil {
			return fmt.Errorf("could not create table %s: %w", t.Name, err)
		}
	}

	if err = g.After(ctx); err != nil {
		return fmt.Errorf("could not run \"after\" upgrade procedures: %w", err)
	}

	return nil
}

func (s Schema) Tables() []*Table {
	return []*Table{
		s.Users(),
		s.Credentials(),
		s.Roles(),
		s.RoleMembers(),
		s.Applications(),
		s.Reminders(),
		s.Attachments(),
		s.ActionLog(),
		s.RbacRules(),
		s.Settings(),
		s.Labels(),
		s.Templates(),
		s.ComposeAttachment(),
		s.ComposeChart(),
		s.ComposeModule(),
		s.ComposeModuleField(),
		s.ComposeNamespace(),
		s.ComposePage(),
		s.ComposeRecord(),
		s.ComposeRecordValue(),
		s.MessagingAttachment(),
		s.MessagingChannel(),
		s.MessagingChannelMember(),
		s.MessagingMention(),
		s.MessagingMessage(),
		s.MessagingMessageAttachment(),
		s.MessagingMessageFlag(),
		s.MessagingUnread(),
		s.FederationModuleShared(),
		s.FederationModuleExposed(),
		s.FederationModuleMapping(),
		s.FederationNodes(),
		s.FederationNodesSync(),
	}
}

func (Schema) Users() *Table {
	return TableDef("users",
		ID,
		ColumnDef("email", ColumnTypeVarchar, ColumnTypeLength(emailLength)),
		ColumnDef("email_confirmed", ColumnTypeBoolean, DefaultValue("false")),
		ColumnDef("username", ColumnTypeText),
		ColumnDef("name", ColumnTypeText),
		ColumnDef("handle", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("kind", ColumnTypeVarchar, ColumnTypeLength(8)),
		ColumnDef("meta", ColumnTypeJson),
		ColumnDef("suspended_at", ColumnTypeTimestamp, Null),
		CUDTimestamps,

		AddIndex("unique_email", IExpr("LOWER(email)"), IWhere("LENGTH(email) > 0 AND deleted_at IS NULL AND suspended_at IS NULL")),
		AddIndex("unique_username", IExpr("LOWER(username)"), IWhere("LENGTH(username) > 0 AND deleted_at IS NULL AND suspended_at IS NULL")),
		AddIndex("unique_handle", IExpr("LOWER(handle)"), IWhere("LENGTH(handle) > 0 AND deleted_at IS NULL AND suspended_at IS NULL")),
	)
}

func (Schema) Credentials() *Table {
	return TableDef(`credentials`,
		ID,
		ColumnDef("rel_owner", ColumnTypeIdentifier),
		ColumnDef("label", ColumnTypeText),
		ColumnDef("kind", ColumnTypeVarchar, ColumnTypeLength(128)),
		ColumnDef("credentials", ColumnTypeText),
		ColumnDef("meta", ColumnTypeJson),
		ColumnDef("expires_at", ColumnTypeTimestamp, Null),
		ColumnDef("last_used_at", ColumnTypeTimestamp, Null),
		CUDTimestamps,

		AddIndex("owner_kind", IColumn("rel_owner", "kind"), IWhere("deleted_at IS NULL")),
	)
}

func (Schema) Roles() *Table {
	return TableDef(`roles`,
		ID,
		ColumnDef("name", ColumnTypeText),
		ColumnDef("handle", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("archived_at", ColumnTypeTimestamp, Null),
		CUDTimestamps,

		AddIndex("unique_handle", IExpr("LOWER(handle)"), IWhere("LENGTH(handle) > 0 AND deleted_at IS NULL AND archived_at IS NULL")),
	)
}

func (Schema) RoleMembers() *Table {
	return TableDef(`role_members`,
		ColumnDef("rel_role", ColumnTypeIdentifier),
		ColumnDef("rel_user", ColumnTypeIdentifier),

		AddIndex("unique_membership", IColumn("rel_role", "rel_user")),
	)
}

func (Schema) Applications() *Table {
	return TableDef("applications",
		ID,
		ColumnDef("name", ColumnTypeText),
		ColumnDef("enabled", ColumnTypeBoolean, DefaultValue("true")),
		ColumnDef("unify", ColumnTypeJson),
		ColumnDef("rel_owner", ColumnTypeIdentifier),
		CUDTimestamps,
	)
}

func (Schema) Reminders() *Table {
	return TableDef("reminders",
		ID,
		ColumnDef("resource", ColumnTypeVarchar, ColumnTypeLength(resourceLength)),
		ColumnDef("payload", ColumnTypeJson),
		ColumnDef("snooze_count", ColumnTypeInteger, DefaultValue("0")),
		ColumnDef("assigned_to", ColumnTypeIdentifier, DefaultValue("0")),
		ColumnDef("assigned_by", ColumnTypeIdentifier, DefaultValue("0")),
		ColumnDef("assigned_at", ColumnTypeTimestamp),
		ColumnDef("remind_at", ColumnTypeTimestamp, Null),
		ColumnDef("dismissed_at", ColumnTypeTimestamp, Null),
		ColumnDef("dismissed_by", ColumnTypeIdentifier, DefaultValue("0")),
		CUDTimestamps,

		AddIndex("resource", IColumn("resource")),
		AddIndex("assignee", IColumn("assigned_to")),
	)
}

func (Schema) Attachments() *Table {
	return TableDef("attachments",
		ID,
		ColumnDef("rel_owner", ColumnTypeIdentifier),
		ColumnDef("kind", ColumnTypeText),
		ColumnDef("url", ColumnTypeText),
		ColumnDef("preview_url", ColumnTypeText),
		ColumnDef("name", ColumnTypeText),
		ColumnDef("meta", ColumnTypeJson),
		CUDTimestamps,
	)
}

func (Schema) ActionLog() *Table {
	return TableDef("actionlog",
		ID,
		ColumnDef("ts", ColumnTypeTimestamp),
		ColumnDef("actor_ip_addr", ColumnTypeVarchar, ColumnTypeLength(15)),
		ColumnDef("actor_id", ColumnTypeIdentifier),
		ColumnDef("request_origin", ColumnTypeVarchar, ColumnTypeLength(32)),
		ColumnDef("request_id", ColumnTypeVarchar, ColumnTypeLength(256)),
		ColumnDef("resource", ColumnTypeVarchar, ColumnTypeLength(resourceLength)),
		ColumnDef("action", ColumnTypeVarchar, ColumnTypeLength(64)),
		ColumnDef("error", ColumnTypeText),
		ColumnDef("severity", ColumnTypeInteger),
		ColumnDef("description", ColumnTypeText),
		ColumnDef("meta", ColumnTypeJson),

		AddIndex("ts", IColumn("ts")),
		AddIndex("request_origin", IColumn("request_origin")),
		AddIndex("actor_id", IColumn("actor_id")),
		AddIndex("resource", IColumn("resource")),
		AddIndex("action", IColumn("action")),
	)
}

func (Schema) RbacRules() *Table {
	return TableDef("rbac_rules",
		ColumnDef("rel_role", ColumnTypeIdentifier),
		ColumnDef("resource", ColumnTypeVarchar, ColumnTypeLength(resourceLength)),
		ColumnDef("operation", ColumnTypeVarchar, ColumnTypeLength(50)),
		ColumnDef("access", ColumnTypeInteger),

		PrimaryKey(IColumn("rel_role", "resource", "operation")),
	)
}

func (Schema) Settings() *Table {
	return TableDef("settings",
		ColumnDef("rel_owner", ColumnTypeIdentifier),
		ColumnDef("name", ColumnTypeVarchar, ColumnTypeLength(resourceLength)),
		ColumnDef("value", ColumnTypeJson),
		ColumnDef("updated_by", ColumnTypeIdentifier),
		ColumnDef("updated_at", ColumnTypeTimestamp),

		AddIndex("unique_name", IExpr("LOWER(name)"), IColumn("rel_owner")),
	)
}

func (Schema) Labels() *Table {
	return TableDef("labels",
		ColumnDef("kind", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("rel_resource", ColumnTypeIdentifier),
		ColumnDef("name", ColumnTypeVarchar, ColumnTypeLength(resourceLength)),
		ColumnDef("value", ColumnTypeText),

		AddIndex("unique_kind_res_name", IColumn("kind", "rel_resource"), IExpr("LOWER(name)")),
	)
}

func (Schema) Templates() *Table {
	return TableDef("templates",
		ID,
		ColumnDef("rel_owner", ColumnTypeIdentifier),
		ColumnDef("handle", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("language", ColumnTypeText),
		ColumnDef("type", ColumnTypeText),
		ColumnDef("partial", ColumnTypeBoolean),
		ColumnDef("meta", ColumnTypeJson),
		ColumnDef("template", ColumnTypeText),
		CUDTimestamps,
		ColumnDef("last_used_at", ColumnTypeTimestamp, Null),

		AddIndex("unique_language_handle", IColumn("language"), IExpr("LOWER(handle)"), IWhere("LENGTH(handle) > 0 AND deleted_at IS NULL")),
	)
}

func (Schema) ComposeAttachment() *Table {
	// @todo merge with general attachment table

	return TableDef("compose_attachment",
		ID,
		ColumnDef("rel_namespace", ColumnTypeIdentifier),
		ColumnDef("rel_owner", ColumnTypeIdentifier),
		ColumnDef("kind", ColumnTypeText),
		ColumnDef("url", ColumnTypeText),
		ColumnDef("preview_url", ColumnTypeText),
		ColumnDef("name", ColumnTypeText),
		ColumnDef("meta", ColumnTypeJson),
		CUDTimestamps,

		AddIndex("namespace", IColumn("rel_namespace")),
	)
}

//func (Schema) ComposeRecordAttachment() *Table {
//	return TableDef("compose_record_attachments",
//		ColumnDef("rel_attachment", ColumnTypeIdentifier),
//		ColumnDef("rel_namespace", ColumnTypeIdentifier),
//		ColumnDef("rel_module", ColumnTypeIdentifier),
//		ColumnDef("rel_record", ColumnTypeIdentifier),
//		ColumnDef("field", ColumnTypeIdentifier),
//		ColumnDef("owned_by", ColumnTypeIdentifier),
//		CUDTimestamps,
//		CUDUsers,
//
//		PrimaryKey(IColumn("rel_attachment", "rel_namespace", "rel_module", "field")),
//	)
//}

//func (Schema) ComposePageAttachment() *Table {
//	return TableDef("compose_record_attachments",
//		ColumnDef("rel_attachment", ColumnTypeIdentifier),
//		ColumnDef("rel_namespace", ColumnTypeIdentifier),
//		ColumnDef("rel_page", ColumnTypeIdentifier),
//		ColumnDef("owned_by", ColumnTypeIdentifier),
//		CUDTimestamps,
//		CUDUsers,
//
//		PrimaryKey(IColumn("rel_attachment", "rel_namespace", "rel_page")),
//	)
//}

func (Schema) ComposeChart() *Table {
	return TableDef("compose_chart",
		ID,
		ColumnDef("handle", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("rel_namespace", ColumnTypeIdentifier),
		ColumnDef("name", ColumnTypeText),
		ColumnDef("config", ColumnTypeJson),
		CUDTimestamps,

		AddIndex("namespace", IColumn("rel_namespace")),
		AddIndex("unique_handle", IColumn("rel_namespace"), IExpr("LOWER(handle)"), IWhere("LENGTH(handle) > 0 AND deleted_at IS NULL")),
	)
}

func (Schema) ComposeModule() *Table {
	return TableDef("compose_module",
		ID,
		ColumnDef("rel_namespace", ColumnTypeIdentifier),
		ColumnDef("handle", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("name", ColumnTypeText),
		ColumnDef("meta", ColumnTypeJson),
		CUDTimestamps,

		AddIndex("namespace", IColumn("rel_namespace")),
		AddIndex("unique_handle", IColumn("rel_namespace"), IExpr("LOWER(handle)"), IWhere("LENGTH(handle) > 0 AND deleted_at IS NULL")),
	)
}

func (Schema) ComposeModuleField() *Table {
	return TableDef("compose_module_field",
		ID,

		ColumnDef("rel_module", ColumnTypeIdentifier),
		ColumnDef("place", ColumnTypeInteger),
		ColumnDef("kind", ColumnTypeText),
		ColumnDef("options", ColumnTypeJson),
		ColumnDef("default_value", ColumnTypeJson),
		ColumnDef("expressions", ColumnTypeJson),
		ColumnDef("name", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("label", ColumnTypeText),
		ColumnDef("is_private", ColumnTypeBoolean),
		ColumnDef("is_required", ColumnTypeBoolean),
		ColumnDef("is_visible", ColumnTypeBoolean),
		ColumnDef("is_multi", ColumnTypeBoolean),

		CUDTimestamps,

		AddIndex("unique_name", IColumn("rel_module"), IExpr("LOWER(name)"), IWhere("LENGTH(name) > 0 AND deleted_at IS NULL")),
		AddIndex("module", IColumn("rel_module")),
	)
}

func (Schema) ComposeNamespace() *Table {
	return TableDef("compose_namespace",
		ID,
		ColumnDef("name", ColumnTypeText),
		ColumnDef("slug", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("enabled", ColumnTypeBoolean),
		ColumnDef("meta", ColumnTypeJson),
		CUDTimestamps,

		AddIndex("unique_slug", IExpr("LOWER(slug)"), IWhere("LENGTH(slug) > 0 AND deleted_at IS NULL")),
	)
}

func (Schema) ComposePage() *Table {
	return TableDef("compose_page",
		ID,
		ColumnDef("title", ColumnTypeText),
		ColumnDef("handle", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("description", ColumnTypeText),
		ColumnDef("rel_namespace", ColumnTypeIdentifier),
		ColumnDef("rel_module", ColumnTypeIdentifier),
		ColumnDef("self_id", ColumnTypeIdentifier),
		ColumnDef("blocks", ColumnTypeJson),
		ColumnDef("visible", ColumnTypeBoolean),
		ColumnDef("weight", ColumnTypeInteger),
		CUDTimestamps,

		AddIndex("self", IColumn("self_id")),
		AddIndex("namespace", IColumn("rel_namespace")),
		AddIndex("module", IColumn("rel_module")),
		AddIndex("unique_handle", IColumn("rel_namespace"), IExpr("LOWER(handle)"), IWhere("LENGTH(handle) > 0 AND deleted_at IS NULL")),
	)
}

func (Schema) ComposeRecord() *Table {
	return TableDef("compose_record",
		ID,
		ColumnDef("rel_namespace", ColumnTypeIdentifier),
		ColumnDef("module_id", ColumnTypeIdentifier),
		ColumnDef("owned_by", ColumnTypeIdentifier),
		CUDTimestamps,
		CUDUsers,

		AddIndex("namespace", IColumn("rel_namespace")),
		AddIndex("module", IColumn("module_id")),
		AddIndex("owner", IColumn("owned_by")),
	)
}

func (Schema) ComposeRecordValue() *Table {
	return TableDef("compose_record_value",
		ColumnDef("record_id", ColumnTypeIdentifier),
		ColumnDef("name", ColumnTypeVarchar, ColumnTypeLength(64)),
		ColumnDef("value", ColumnTypeText, ColumnTypeFlag("mysqlLongText", true)),
		ColumnDef("ref", ColumnTypeIdentifier),
		ColumnDef("place", ColumnTypeInteger),
		ColumnDef("deleted_at", ColumnTypeTimestamp, Null),

		PrimaryKey(IColumn("record_id", "name", "place")),
		AddIndex("ref", IColumn("ref"), IWhere("ref > 0")),
	)
}

func (Schema) MessagingAttachment() *Table {
	// @todo merge with general attachment table
	return TableDef("messaging_attachment",
		ID,
		ColumnDef("rel_owner", ColumnTypeIdentifier),
		ColumnDef("url", ColumnTypeText),
		ColumnDef("preview_url", ColumnTypeText),
		ColumnDef("name", ColumnTypeText),
		ColumnDef("meta", ColumnTypeJson),
		CUDTimestamps,
	)
}

func (Schema) MessagingChannel() *Table {
	return TableDef("messaging_channel",
		ID,
		ColumnDef("name", ColumnTypeText),
		ColumnDef("topic", ColumnTypeText),
		ColumnDef("meta", ColumnTypeJson),
		ColumnDef("type", ColumnTypeText),
		ColumnDef("membership_policy", ColumnTypeText),
		ColumnDef("rel_creator", ColumnTypeIdentifier), // @todo rename => created_by
		ColumnDef("archived_at", ColumnTypeTimestamp, Null),
		ColumnDef("rel_last_message", ColumnTypeIdentifier),
		CUDTimestamps,
	)
}

func (Schema) MessagingChannelMember() *Table {
	return TableDef("messaging_channel_member",
		ColumnDef("rel_channel", ColumnTypeIdentifier),
		ColumnDef("rel_user", ColumnTypeIdentifier),
		ColumnDef("type", ColumnTypeText),
		ColumnDef("flag", ColumnTypeText),
		CUDTimestamps,
		PrimaryKey(IColumn("rel_channel", "rel_user")),
	)
}

func (Schema) MessagingMention() *Table {
	return TableDef("messaging_mention",
		ID,
		ColumnDef("rel_channel", ColumnTypeIdentifier),
		ColumnDef("rel_message", ColumnTypeIdentifier),
		ColumnDef("rel_user", ColumnTypeIdentifier),
		ColumnDef("rel_mentioned_by", ColumnTypeIdentifier),
		ColumnDef("created_at", ColumnTypeTimestamp),
	)
}

func (Schema) MessagingMessage() *Table {
	return TableDef("messaging_message",
		ID,
		ColumnDef("type", ColumnTypeText),
		ColumnDef("message", ColumnTypeText),
		ColumnDef("meta", ColumnTypeJson),
		ColumnDef("rel_channel", ColumnTypeIdentifier),
		ColumnDef("rel_user", ColumnTypeIdentifier),
		ColumnDef("reply_to", ColumnTypeIdentifier, DefaultValue("0")),
		ColumnDef("replies", ColumnTypeInteger, DefaultValue("0")),

		CUDTimestamps,
	)
}

func (Schema) MessagingMessageAttachment() *Table {
	return TableDef("messaging_message_attachment",
		ColumnDef("rel_message", ColumnTypeIdentifier),
		ColumnDef("rel_attachment", ColumnTypeIdentifier),
		PrimaryKey(IColumn("rel_message")),
	)
}

func (Schema) MessagingMessageFlag() *Table {
	return TableDef("messaging_message_flag",
		ID,
		ColumnDef("rel_channel", ColumnTypeIdentifier),
		ColumnDef("rel_message", ColumnTypeIdentifier),
		ColumnDef("rel_user", ColumnTypeIdentifier),
		ColumnDef("flag", ColumnTypeText),
		ColumnDef("created_at", ColumnTypeTimestamp),
	)
}

func (Schema) MessagingUnread() *Table {
	return TableDef("messaging_unread",
		ColumnDef("rel_channel", ColumnTypeIdentifier),
		ColumnDef("rel_reply_to", ColumnTypeIdentifier),
		ColumnDef("rel_user", ColumnTypeIdentifier),
		ColumnDef("count", ColumnTypeInteger),
		ColumnDef("rel_last_message", ColumnTypeIdentifier),
		PrimaryKey(IColumn("rel_channel", "rel_reply_to", "rel_user")),
	)
}

func (Schema) FederationModuleShared() *Table {
	return TableDef("federation_module_shared",
		ID,
		ColumnDef("handle", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("name", ColumnTypeText),
		ColumnDef("rel_node", ColumnTypeIdentifier),
		ColumnDef("xref_module", ColumnTypeIdentifier),
		ColumnDef("fields", ColumnTypeText),
		CUDTimestamps,
		CUDUsers,
	)
}

func (Schema) FederationModuleExposed() *Table {
	return TableDef("federation_module_exposed",
		ID,
		ColumnDef("handle", ColumnTypeVarchar, ColumnTypeLength(handleLength)),
		ColumnDef("name", ColumnTypeText),
		ColumnDef("rel_node", ColumnTypeIdentifier),
		ColumnDef("rel_compose_module", ColumnTypeIdentifier),
		ColumnDef("rel_compose_namespace", ColumnTypeIdentifier),
		ColumnDef("fields", ColumnTypeText),
		CUDTimestamps,
		CUDUsers,

		AddIndex("unique_node_compose_module", IColumn("rel_node", "rel_compose_module", "rel_compose_namespace")),
	)
}

func (Schema) FederationModuleMapping() *Table {
	return TableDef("federation_module_mapping",
		ColumnDef("rel_federation_module", ColumnTypeIdentifier),
		ColumnDef("rel_compose_module", ColumnTypeIdentifier),
		ColumnDef("rel_compose_namespace", ColumnTypeIdentifier),
		ColumnDef("field_mapping", ColumnTypeText),

		AddIndex("unique_module_compose_module", IColumn("rel_federation_module", "rel_compose_module", "rel_compose_namespace")),
	)
}

func (Schema) FederationNodes() *Table {
	return TableDef("federation_nodes",
		ID,
		ColumnDef("shared_node_id", ColumnTypeIdentifier),
		ColumnDef("name", ColumnTypeText),
		ColumnDef("base_url", ColumnTypeText),
		ColumnDef("status", ColumnTypeText),
		ColumnDef("contact", ColumnTypeVarchar, ColumnTypeLength(emailLength)),
		ColumnDef("pair_token", ColumnTypeText),
		ColumnDef("auth_token", ColumnTypeText),

		CUDTimestamps,
		CUDUsers,
	)
}

func (Schema) FederationNodesSync() *Table {
	return TableDef("federation_nodes_sync",
		ColumnDef("rel_node", ColumnTypeIdentifier),
		ColumnDef("rel_module", ColumnTypeIdentifier),
		ColumnDef("sync_type", ColumnTypeText),
		ColumnDef("sync_status", ColumnTypeText),
		ColumnDef("time_action", ColumnTypeTimestamp),
	)
}
