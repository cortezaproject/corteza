package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - store/actionlog.yaml
//  - store/compose_charts.yaml
//  - store/compose_module_fields.yaml
//  - store/compose_modules.yaml
//  - store/compose_namespaces.yaml
//  - store/compose_pages.yaml
//  - store/rbac_rules.yaml

type (
	// Interface combines interfaces of all supported store interfaces
	storeInterface interface {
		actionlogsStore
		composeChartsStore
		composeModuleFieldsStore
		composeModulesStore
		composeNamespacesStore
		composePagesStore
		rbacRulesStore
	}
)
