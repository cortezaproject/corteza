package resource

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

// ComposeAttachmentRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeAttachmentRbacReferences(attachment string) (res *Ref, pp []*Ref, err error) {
	if attachment != "*" {
		res = &Ref{ResourceType: types.AttachmentResourceType, Identifiers: MakeIdentifiers(attachment)}
	}

	return
}

// ComposeChartRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeChartRbacReferences(namespaceID string, chart string) (res *Ref, pp []*Ref, err error) {
	if namespaceID != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespaceID)})
	}
	if chart != "*" {
		res = &Ref{ResourceType: types.ChartResourceType, Identifiers: MakeIdentifiers(chart)}
	}

	return
}

// ComposeModuleRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeModuleRbacReferences(namespaceID string, module string) (res *Ref, pp []*Ref, err error) {
	if namespaceID != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespaceID)})
	}
	if module != "*" {
		res = &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(module)}
	}

	return
}

// ComposeModuleFieldRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeModuleFieldRbacReferences(namespaceID string, moduleID string, moduleField string) (res *Ref, pp []*Ref, err error) {
	if namespaceID != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespaceID)})
	}
	if moduleID != "*" {
		pp = append(pp, &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(moduleID)})
	}
	if moduleField != "*" {
		res = &Ref{ResourceType: types.ModuleFieldResourceType, Identifiers: MakeIdentifiers(moduleField)}
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
func ComposePageRbacReferences(namespaceID string, page string) (res *Ref, pp []*Ref, err error) {
	if namespaceID != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespaceID)})
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
func ComposeRecordRbacReferences(namespaceID string, moduleID string, record string) (res *Ref, pp []*Ref, err error) {
	if namespaceID != "*" {
		pp = append(pp, &Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers(namespaceID)})
	}
	if moduleID != "*" {
		pp = append(pp, &Ref{ResourceType: types.ModuleResourceType, Identifiers: MakeIdentifiers(moduleID)})
	}
	if record != "*" {
		res = &Ref{ResourceType: types.RecordResourceType, Identifiers: MakeIdentifiers(record)}
	}

	return
}

// ComposeRecordValueRbacReferences generates RBAC references
//
// Resources with "envoy: false" are skipped
//
// This function is auto-generated
func ComposeRecordValueRbacReferences(recordValue string) (res *Ref, pp []*Ref, err error) {
	if recordValue != "*" {
		res = &Ref{ResourceType: types.RecordValueResourceType, Identifiers: MakeIdentifiers(recordValue)}
	}

	return
}
