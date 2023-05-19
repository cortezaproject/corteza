package tests

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"testing"

	"github.com/cortezaproject/corteza/server/store"
)

func testAllGenerated(t *testing.T, s store.Storer) {

	t.Run("actionlog", func(t *testing.T) {
		testActionlogs(t, s)
	})
	t.Run("apigwFilter", func(t *testing.T) {
		testApigwFilters(t, s)
	})
	t.Run("apigwRoute", func(t *testing.T) {
		testApigwRoutes(t, s)
	})
	t.Run("application", func(t *testing.T) {
		testApplications(t, s)
	})
	t.Run("attachment", func(t *testing.T) {
		testAttachments(t, s)
	})
	t.Run("authClient", func(t *testing.T) {
		testAuthClients(t, s)
	})
	t.Run("authConfirmedClient", func(t *testing.T) {
		testAuthConfirmedClients(t, s)
	})
	t.Run("authOa2token", func(t *testing.T) {
		testAuthOa2tokens(t, s)
	})
	t.Run("authSession", func(t *testing.T) {
		testAuthSessions(t, s)
	})
	t.Run("automationSession", func(t *testing.T) {
		testAutomationSessions(t, s)
	})
	t.Run("automationTrigger", func(t *testing.T) {
		testAutomationTriggers(t, s)
	})
	t.Run("automationWorkflow", func(t *testing.T) {
		testAutomationWorkflows(t, s)
	})
	t.Run("composeAttachment", func(t *testing.T) {
		testComposeAttachments(t, s)
	})
	t.Run("composeChart", func(t *testing.T) {
		testComposeCharts(t, s)
	})
	t.Run("composeModule", func(t *testing.T) {
		testComposeModules(t, s)
	})
	t.Run("composeModuleField", func(t *testing.T) {
		testComposeModuleFields(t, s)
	})
	t.Run("composeNamespace", func(t *testing.T) {
		testComposeNamespaces(t, s)
	})
	t.Run("composePage", func(t *testing.T) {
		testComposePages(t, s)
	})
	t.Run("composePageLayout", func(t *testing.T) {
		testComposePageLayouts(t, s)
	})
	t.Run("credential", func(t *testing.T) {
		testCredentials(t, s)
	})
	t.Run("dalConnection", func(t *testing.T) {
		testDalConnections(t, s)
	})
	t.Run("dalSchemaAlteration", func(t *testing.T) {
		testDalSchemaAlterations(t, s)
	})
	t.Run("dalSensitivityLevel", func(t *testing.T) {
		testDalSensitivityLevels(t, s)
	})
	t.Run("dataPrivacyRequest", func(t *testing.T) {
		testDataPrivacyRequests(t, s)
	})
	t.Run("dataPrivacyRequestComment", func(t *testing.T) {
		testDataPrivacyRequestComments(t, s)
	})
	t.Run("federationExposedModule", func(t *testing.T) {
		testFederationExposedModules(t, s)
	})
	t.Run("federationModuleMapping", func(t *testing.T) {
		testFederationModuleMappings(t, s)
	})
	t.Run("federationNode", func(t *testing.T) {
		testFederationNodes(t, s)
	})
	t.Run("federationNodeSync", func(t *testing.T) {
		testFederationNodeSyncs(t, s)
	})
	t.Run("federationSharedModule", func(t *testing.T) {
		testFederationSharedModules(t, s)
	})
	t.Run("flag", func(t *testing.T) {
		testFlags(t, s)
	})
	t.Run("label", func(t *testing.T) {
		testLabels(t, s)
	})
	t.Run("queue", func(t *testing.T) {
		testQueues(t, s)
	})
	t.Run("queueMessage", func(t *testing.T) {
		testQueueMessages(t, s)
	})
	t.Run("rbacRule", func(t *testing.T) {
		testRbacRules(t, s)
	})
	t.Run("reminder", func(t *testing.T) {
		testReminders(t, s)
	})
	t.Run("report", func(t *testing.T) {
		testReports(t, s)
	})
	t.Run("resourceActivity", func(t *testing.T) {
		testResourceActivitys(t, s)
	})
	t.Run("resourceTranslation", func(t *testing.T) {
		testResourceTranslations(t, s)
	})
	t.Run("role", func(t *testing.T) {
		testRoles(t, s)
	})
	t.Run("roleMember", func(t *testing.T) {
		testRoleMembers(t, s)
	})
	t.Run("settingValue", func(t *testing.T) {
		testSettingValues(t, s)
	})
	t.Run("template", func(t *testing.T) {
		testTemplates(t, s)
	})
	t.Run("user", func(t *testing.T) {
		testUsers(t, s)
	})
}
