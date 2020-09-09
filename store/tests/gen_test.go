package tests

// This file is auto-generated.
//
// Template:	pkg/codegen/assets/store_test_all.gen.go.tpl
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
//  - store/credentials.yaml
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

import (
	"github.com/cortezaproject/corteza-server/store"
	"testing"
)

func testAllGenerated(t *testing.T, s store.Storer) {
	// Run generated tests for Actionlog
	t.Run("Actionlog", func(t *testing.T) {
		testActionlog(t, s)
	})

	// Run generated tests for Applications
	t.Run("Applications", func(t *testing.T) {
		testApplications(t, s)
	})

	// Run generated tests for Attachment
	t.Run("Attachment", func(t *testing.T) {
		testAttachment(t, s)
	})

	// Run generated tests for ComposeAttachments
	t.Run("ComposeAttachments", func(t *testing.T) {
		testComposeAttachments(t, s)
	})

	// Run generated tests for ComposeCharts
	t.Run("ComposeCharts", func(t *testing.T) {
		testComposeCharts(t, s)
	})

	// Run generated tests for ComposeModuleFields
	t.Run("ComposeModuleFields", func(t *testing.T) {
		testComposeModuleFields(t, s)
	})

	// Run generated tests for ComposeModules
	t.Run("ComposeModules", func(t *testing.T) {
		testComposeModules(t, s)
	})

	// Run generated tests for ComposeNamespaces
	t.Run("ComposeNamespaces", func(t *testing.T) {
		testComposeNamespaces(t, s)
	})

	// Run generated tests for ComposePages
	t.Run("ComposePages", func(t *testing.T) {
		testComposePages(t, s)
	})

	// Run generated tests for ComposeRecordValues
	t.Run("ComposeRecordValues", func(t *testing.T) {
		testComposeRecordValues(t, s)
	})

	// Run generated tests for ComposeRecords
	t.Run("ComposeRecords", func(t *testing.T) {
		testComposeRecords(t, s)
	})

	// Run generated tests for Credentials
	t.Run("Credentials", func(t *testing.T) {
		testCredentials(t, s)
	})

	// Run generated tests for MessagingAttachments
	t.Run("MessagingAttachments", func(t *testing.T) {
		testMessagingAttachments(t, s)
	})

	// Run generated tests for MessagingChannelMembers
	t.Run("MessagingChannelMembers", func(t *testing.T) {
		testMessagingChannelMembers(t, s)
	})

	// Run generated tests for MessagingChannels
	t.Run("MessagingChannels", func(t *testing.T) {
		testMessagingChannels(t, s)
	})

	// Run generated tests for MessagingFlags
	t.Run("MessagingFlags", func(t *testing.T) {
		testMessagingFlags(t, s)
	})

	// Run generated tests for MessagingMentions
	t.Run("MessagingMentions", func(t *testing.T) {
		testMessagingMentions(t, s)
	})

	// Run generated tests for MessagingMessageAttachments
	t.Run("MessagingMessageAttachments", func(t *testing.T) {
		testMessagingMessageAttachments(t, s)
	})

	// Run generated tests for MessagingMessages
	t.Run("MessagingMessages", func(t *testing.T) {
		testMessagingMessages(t, s)
	})

	// Run generated tests for MessagingUnread
	t.Run("MessagingUnread", func(t *testing.T) {
		testMessagingUnread(t, s)
	})

	// Run generated tests for RbacRules
	t.Run("RbacRules", func(t *testing.T) {
		testRbacRules(t, s)
	})

	// Run generated tests for Reminders
	t.Run("Reminders", func(t *testing.T) {
		testReminders(t, s)
	})

	// Run generated tests for RoleMembers
	t.Run("RoleMembers", func(t *testing.T) {
		testRoleMembers(t, s)
	})

	// Run generated tests for Roles
	t.Run("Roles", func(t *testing.T) {
		testRoles(t, s)
	})

	// Run generated tests for Settings
	t.Run("Settings", func(t *testing.T) {
		testSettings(t, s)
	})

	// Run generated tests for Users
	t.Run("Users", func(t *testing.T) {
		testUsers(t, s)
	})
}
