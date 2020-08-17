package tests

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func testAllGenerated(t *testing.T, all interface{}) {

	// Run generated tests for Actionlog
	t.Run("Actionlog", func(t *testing.T) {
		var s = all.(actionlogsStore)
		require.New(t).NoError(s.TruncateActionlogs(context.Background()))
		testActionlog(t, s)
	})

	// Run generated tests for Applications
	t.Run("Applications", func(t *testing.T) {
		var s = all.(applicationsStore)
		require.New(t).NoError(s.TruncateApplications(context.Background()))
		testApplications(t, s)
	})

	// Run generated tests for Attachment
	t.Run("Attachment", func(t *testing.T) {
		var s = all.(attachmentsStore)
		require.New(t).NoError(s.TruncateAttachments(context.Background()))
		testAttachment(t, s)
	})

	// Run generated tests for ComposeCharts
	t.Run("ComposeCharts", func(t *testing.T) {
		var s = all.(composeChartsStore)
		require.New(t).NoError(s.TruncateComposeCharts(context.Background()))
		testComposeCharts(t, s)
	})

	// Run generated tests for ComposeModuleFields
	t.Run("ComposeModuleFields", func(t *testing.T) {
		var s = all.(composeModuleFieldsStore)
		require.New(t).NoError(s.TruncateComposeModuleFields(context.Background()))
		testComposeModuleFields(t, s)
	})

	// Run generated tests for ComposeModules
	t.Run("ComposeModules", func(t *testing.T) {
		var s = all.(composeModulesStore)
		require.New(t).NoError(s.TruncateComposeModules(context.Background()))
		testComposeModules(t, s)
	})

	// Run generated tests for ComposeNamespaces
	t.Run("ComposeNamespaces", func(t *testing.T) {
		var s = all.(composeNamespacesStore)
		require.New(t).NoError(s.TruncateComposeNamespaces(context.Background()))
		testComposeNamespaces(t, s)
	})

	// Run generated tests for ComposePages
	t.Run("ComposePages", func(t *testing.T) {
		var s = all.(composePagesStore)
		require.New(t).NoError(s.TruncateComposePages(context.Background()))
		testComposePages(t, s)
	})

	// Run generated tests for Credentials
	t.Run("Credentials", func(t *testing.T) {
		var s = all.(credentialsStore)
		require.New(t).NoError(s.TruncateCredentials(context.Background()))
		testCredentials(t, s)
	})

	// Run generated tests for RbacRules
	t.Run("RbacRules", func(t *testing.T) {
		var s = all.(rbacRulesStore)
		require.New(t).NoError(s.TruncateRbacRules(context.Background()))
		testRbacRules(t, s)
	})

	// Run generated tests for Reminders
	t.Run("Reminders", func(t *testing.T) {
		var s = all.(remindersStore)
		require.New(t).NoError(s.TruncateReminders(context.Background()))
		testReminders(t, s)
	})

	// Run generated tests for RoleMembers
	t.Run("RoleMembers", func(t *testing.T) {
		var s = all.(roleMembersStore)
		require.New(t).NoError(s.TruncateRoleMembers(context.Background()))
		testRoleMembers(t, s)
	})

	// Run generated tests for Roles
	t.Run("Roles", func(t *testing.T) {
		var s = all.(rolesStore)
		require.New(t).NoError(s.TruncateRoles(context.Background()))
		testRoles(t, s)
	})

	// Run generated tests for Settings
	t.Run("Settings", func(t *testing.T) {
		var s = all.(settingsStore)
		require.New(t).NoError(s.TruncateSettings(context.Background()))
		testSettings(t, s)
	})

	// Run generated tests for Users
	t.Run("Users", func(t *testing.T) {
		var s = all.(usersStore)
		require.New(t).NoError(s.TruncateUsers(context.Background()))
		testUsers(t, s)
	})

}
