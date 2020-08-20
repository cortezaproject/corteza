package tests

// This file is auto-generated.
//
// Template:	pkg/codegen/assets/store_test_all.gen.go
// Definitions:
//  - store/actionlog.yaml
//  - store/applications.yaml
//  - store/attachments.yaml
//  - store/compose_charts.yaml
//  - store/compose_module_fields.yaml
//  - store/compose_modules.yaml
//  - store/compose_namespaces.yaml
//  - store/compose_pages.yaml
//  - store/credentials.yaml
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

func testAllGenerated(t *testing.T, s store.Storable) {
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

	// Run generated tests for Credentials
	t.Run("Credentials", func(t *testing.T) {
		testCredentials(t, s)
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
