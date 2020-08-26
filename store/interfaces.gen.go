package store

// This file is auto-generated.
//
// Template:	pkg/codegen/assets/store_interfaces_joined.gen.go.tpl
// Definitions:
//  - store/actionlog.yaml
//  - store/applications.yaml
//  - store/attachments.yaml
//  - store/compose_charts.yaml
//  - store/compose_module_fields.yaml
//  - store/compose_modules.yaml
//  - store/compose_namespaces.yaml
//  - store/compose_pages.yaml
//  - store/compose_record_values.yaml
//  - store/compose_records.yaml
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
	"context"
)

type (
	Transactionable interface {
		Tx(context.Context, func(context.Context, Storable) error) error
	}

	// Sortable interface combines interfaces of all supported store interfaces
	Storable interface {
		Transactionable

		Actionlogs
		Applications
		Attachments
		ComposeCharts
		ComposeModuleFields
		ComposeModules
		ComposeNamespaces
		ComposePages
		ComposeRecordValues
		ComposeRecords
		Credentials
		RbacRules
		Reminders
		RoleMembers
		Roles
		Settings
		Users
	}
)
