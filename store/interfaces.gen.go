package store

import "context"

// This file is auto-generated.
//
// Template:	pkg/codegen/assets/store_interfaces_joined.gen.go.tpl
// Definitions:
//  - store/actionlog.yaml
//  - store/applications.yaml
//  - store/attachments.yaml
//  - store/compose_attachments.yaml
//  - store/compose_charts.yaml
//  - store/compose_module_fields.yaml
//  - store/compose_modules.yaml
//  - store/compose_namespaces.yaml
//  - store/compose_pages.yaml
//  - store/compose_record_values.yaml
//  - store/compose_records.yaml
//  - store/credentials.yaml
//  - store/labels.yaml
//  - store/messaging_attachments.yaml
//  - store/messaging_channel_members.yaml
//  - store/messaging_channels.yaml
//  - store/messaging_flags.yaml
//  - store/messaging_mentions.yaml
//  - store/messaging_message_attachments.yaml
//  - store/messaging_messages.yaml
//  - store/messaging_unread.yaml
//  - store/rbac_rules.yaml
//  - store/reminders.yaml
//  - store/role_members.yaml
//  - store/roles.yaml
//  - store/settings.yaml
//  - store/users.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	Transactioner interface {
		Tx(context.Context, func(context.Context, Storer) error) error
	}

	// Storer interface combines interfaces of all supported store interfaces
	Storer interface {
		Transactioner

		Actionlogs
		Applications
		Attachments
		ComposeAttachments
		ComposeCharts
		ComposeModuleFields
		ComposeModules
		ComposeNamespaces
		ComposePages
		ComposeRecordValues
		ComposeRecords
		Credentials
		Labels
		MessagingAttachments
		MessagingChannelMembers
		MessagingChannels
		MessagingFlags
		MessagingMentions
		MessagingMessageAttachments
		MessagingMessages
		MessagingUnreads
		RbacRules
		Reminders
		RoleMembers
		Roles
		Settings
		Users
	}
)
