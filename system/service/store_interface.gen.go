package service

// This file is auto-generated.
//
// Template: pkg/store_interfaces_joined.gen.go.tpl
// Definitions:
//  - store/actionlog.yaml
//  - store/applications.yaml
//  - store/credentials.yaml
//  - store/rbac_rules.yaml
//  - store/reminders.yaml
//  - store/roles.yaml
//  - store/settings.yaml
//  - store/system_attachments.yaml
//  - store/users.yaml
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	// Interface combines interfaces of all supported store interfaces
	storeGeneratedInterfaces interface {
		actionlogsStore
		applicationsStore
		credentialsStore
		rbacRulesStore
		remindersStore
		rolesStore
		settingsStore
		attachmentsStore
		usersStore
	}
)
