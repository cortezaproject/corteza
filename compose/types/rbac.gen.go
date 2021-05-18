package types

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
	"strconv"
)

type (
	// Component struct serves as a virtual resource type for the compose component
	//
	// This struct is auto-generated
	Component struct{}
)

const (
	ChartRbacResourceSchema       = "corteza+compose.chart"
	ModuleFieldRbacResourceSchema = "corteza+compose.module-field"
	ModuleRbacResourceSchema      = "corteza+compose.module"
	NamespaceRbacResourceSchema   = "corteza+compose.namespace"
	PageRbacResourceSchema        = "corteza+compose.page"
	RecordRbacResourceSchema      = "corteza+compose.record"
	ComponentRbacResourceSchema   = "corteza+compose"
)

// RbacResource returns string representation of RBAC resource for Chart by calling ChartRbacResource fn
//
// RBAC resource is in the corteza+compose.chart:/... format
//
// This function is auto-generated
func (r Chart) RbacResource() string {
	return ChartRbacResource(r.NamespaceID, r.ID)
}

// ChartRbacResource returns string representation of RBAC resource for Chart
//
// RBAC resource is in the corteza+compose.chart:/... format
//
// This function is auto-generated
func ChartRbacResource(namespaceID uint64, iD uint64) string {
	out := ChartRbacResourceSchema + ":"
	out += "/"

	if namespaceID != 0 {
		out += strconv.FormatUint(namespaceID, 10)
	} else {
		out += "*"
	}
	out += "/"

	if iD != 0 {
		out += strconv.FormatUint(iD, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for ModuleField by calling ModuleFieldRbacResource fn
//
// RBAC resource is in the corteza+compose.module-field:/... format
//
// This function is auto-generated
func (r ModuleField) RbacResource() string {
	return ModuleFieldRbacResource(r.NamespaceID, r.ModuleID, r.ID)
}

// ModuleFieldRbacResource returns string representation of RBAC resource for ModuleField
//
// RBAC resource is in the corteza+compose.module-field:/... format
//
// This function is auto-generated
func ModuleFieldRbacResource(namespaceID uint64, moduleID uint64, iD uint64) string {
	out := ModuleFieldRbacResourceSchema + ":"
	out += "/"

	if namespaceID != 0 {
		out += strconv.FormatUint(namespaceID, 10)
	} else {
		out += "*"
	}
	out += "/"

	if moduleID != 0 {
		out += strconv.FormatUint(moduleID, 10)
	} else {
		out += "*"
	}
	out += "/"

	if iD != 0 {
		out += strconv.FormatUint(iD, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Module by calling ModuleRbacResource fn
//
// RBAC resource is in the corteza+compose.module:/... format
//
// This function is auto-generated
func (r Module) RbacResource() string {
	return ModuleRbacResource(r.NamespaceID, r.ID)
}

// ModuleRbacResource returns string representation of RBAC resource for Module
//
// RBAC resource is in the corteza+compose.module:/... format
//
// This function is auto-generated
func ModuleRbacResource(namespaceID uint64, iD uint64) string {
	out := ModuleRbacResourceSchema + ":"
	out += "/"

	if namespaceID != 0 {
		out += strconv.FormatUint(namespaceID, 10)
	} else {
		out += "*"
	}
	out += "/"

	if iD != 0 {
		out += strconv.FormatUint(iD, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Namespace by calling NamespaceRbacResource fn
//
// RBAC resource is in the corteza+compose.namespace:/... format
//
// This function is auto-generated
func (r Namespace) RbacResource() string {
	return NamespaceRbacResource(r.ID)
}

// NamespaceRbacResource returns string representation of RBAC resource for Namespace
//
// RBAC resource is in the corteza+compose.namespace:/... format
//
// This function is auto-generated
func NamespaceRbacResource(iD uint64) string {
	out := NamespaceRbacResourceSchema + ":"
	out += "/"

	if iD != 0 {
		out += strconv.FormatUint(iD, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Page by calling PageRbacResource fn
//
// RBAC resource is in the corteza+compose.page:/... format
//
// This function is auto-generated
func (r Page) RbacResource() string {
	return PageRbacResource(r.NamespaceID, r.ID)
}

// PageRbacResource returns string representation of RBAC resource for Page
//
// RBAC resource is in the corteza+compose.page:/... format
//
// This function is auto-generated
func PageRbacResource(namespaceID uint64, iD uint64) string {
	out := PageRbacResourceSchema + ":"
	out += "/"

	if namespaceID != 0 {
		out += strconv.FormatUint(namespaceID, 10)
	} else {
		out += "*"
	}
	out += "/"

	if iD != 0 {
		out += strconv.FormatUint(iD, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Record by calling RecordRbacResource fn
//
// RBAC resource is in the corteza+compose.record:/... format
//
// This function is auto-generated
func (r Record) RbacResource() string {
	return RecordRbacResource(r.NamespaceID, r.ModuleID, r.ID)
}

// RecordRbacResource returns string representation of RBAC resource for Record
//
// RBAC resource is in the corteza+compose.record:/... format
//
// This function is auto-generated
func RecordRbacResource(namespaceID uint64, moduleID uint64, iD uint64) string {
	out := RecordRbacResourceSchema + ":"
	out += "/"

	if namespaceID != 0 {
		out += strconv.FormatUint(namespaceID, 10)
	} else {
		out += "*"
	}
	out += "/"

	if moduleID != 0 {
		out += strconv.FormatUint(moduleID, 10)
	} else {
		out += "*"
	}
	out += "/"

	if iD != 0 {
		out += strconv.FormatUint(iD, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza+compose:/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza+compose:/... format
//
// This function is auto-generated
func ComponentRbacResource() string {
	out := ComponentRbacResourceSchema + ":"
	return out
}
