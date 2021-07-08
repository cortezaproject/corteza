package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - compose.chart.yaml
// - compose.module-field.yaml
// - compose.module.yaml
// - compose.namespace.yaml
// - compose.page.yaml
// - compose.record.yaml
// - compose.yaml

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

// ComposeChartRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeChartRbacReferences(namespace string, chart string) (res *Ref, pp []*Ref, err error) {
	if namespace != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)})
	}
	if chart != "*" {
		res = &Ref{ResourceType: types.ChartResourceType, Identifiers: MakeIdentifiers(chart)}
	}

	return
}

// ComposeModuleFieldRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeModuleFieldRbacReferences(namespace string, module string, moduleField string) (res *Ref, pp []*Ref, err error) {
	if namespace != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)})
	}
	if module != "*" {
		pp = append(pp, &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(module)})
	}
	if moduleField != "*" {
		res = &Ref{ResourceType: types.ModuleFieldResourceType, Identifiers: MakeIdentifiers(moduleField)}
	}

	return
}

// ComposeModuleRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeModuleRbacReferences(namespace string, module string) (res *Ref, pp []*Ref, err error) {
	if namespace != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)})
	}
	if module != "*" {
		res = &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(module)}
	}

	return
}

// ComposeNamespaceRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeNamespaceRbacReferences(namespace string) (res *Ref, pp []*Ref, err error) {
	if namespace != "*" {
		res = &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)}
	}

	return
}

// ComposePageRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposePageRbacReferences(namespace string, page string) (res *Ref, pp []*Ref, err error) {
	if namespace != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)})
	}
	if page != "*" {
		res = &Ref{ResourceType: types.PageResourceType, Identifiers: MakeIdentifiers(page)}
	}

	return
}

// ComposeRecordRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeRecordRbacReferences(namespace string, module string, record string) (res *Ref, pp []*Ref, err error) {
	if namespace != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespace)})
	}
	if module != "*" {
		pp = append(pp, &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(module)})
	}
	if record != "*" {
		res = &Ref{ResourceType: types.RecordResourceType, Identifiers: MakeIdentifiers(record)}
	}

	return
}
