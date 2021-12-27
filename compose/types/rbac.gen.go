package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	"strconv"
)

type (
	// Component struct serves as a virtual resource type for the compose component
	//
	// This struct is auto-generated
	Component struct{}
)

var (
	_ = fmt.Printf
	_ = strconv.FormatUint
)

const (
	NamespaceResourceType   = "corteza::compose:namespace"
	ModuleResourceType      = "corteza::compose:module"
	ModuleFieldResourceType = "corteza::compose:module-field"
	RecordResourceType      = "corteza::compose:record"
	PageResourceType        = "corteza::compose:page"
	ChartResourceType       = "corteza::compose:chart"
	ComponentResourceType   = "corteza::compose"
)

// RbacResource returns string representation of RBAC resource for Namespace by calling NamespaceRbacResource fn
//
// RBAC resource is in the corteza::compose:namespace/... format
//
// This function is auto-generated
func (r Namespace) RbacResource() string {
	return NamespaceRbacResource(r.ID)
}

// NamespaceRbacResource returns string representation of RBAC resource for Namespace
//
// RBAC resource is in the corteza::compose:namespace/... format
//
// This function is auto-generated
func NamespaceRbacResource(ID uint64) string {
	cpts := []interface{}{NamespaceResourceType}
	if ID != 0 {
		cpts = append(cpts, strconv.FormatUint(ID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(NamespaceRbacResourceTpl(), cpts...)

}

func NamespaceRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Module by calling ModuleRbacResource fn
//
// RBAC resource is in the corteza::compose:module/... format
//
// This function is auto-generated
func (r Module) RbacResource() string {
	return ModuleRbacResource(r.NamespaceID, r.ID)
}

// ModuleRbacResource returns string representation of RBAC resource for Module
//
// RBAC resource is in the corteza::compose:module/... format
//
// This function is auto-generated
func ModuleRbacResource(NamespaceID uint64, ID uint64) string {
	cpts := []interface{}{ModuleResourceType}
	if NamespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(NamespaceID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if ID != 0 {
		cpts = append(cpts, strconv.FormatUint(ID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ModuleRbacResourceTpl(), cpts...)

}

func ModuleRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for ModuleField by calling ModuleFieldRbacResource fn
//
// RBAC resource is in the corteza::compose:module-field/... format
//
// This function is auto-generated
func (r ModuleField) RbacResource() string {
	return ModuleFieldRbacResource(r.NamespaceID, r.ModuleID, r.ID)
}

// ModuleFieldRbacResource returns string representation of RBAC resource for ModuleField
//
// RBAC resource is in the corteza::compose:module-field/... format
//
// This function is auto-generated
func ModuleFieldRbacResource(NamespaceID uint64, ModuleID uint64, ID uint64) string {
	cpts := []interface{}{ModuleFieldResourceType}
	if NamespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(NamespaceID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if ModuleID != 0 {
		cpts = append(cpts, strconv.FormatUint(ModuleID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if ID != 0 {
		cpts = append(cpts, strconv.FormatUint(ID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ModuleFieldRbacResourceTpl(), cpts...)

}

func ModuleFieldRbacResourceTpl() string {
	return "%s/%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Record by calling RecordRbacResource fn
//
// RBAC resource is in the corteza::compose:record/... format
//
// This function is auto-generated
func (r Record) RbacResource() string {
	return RecordRbacResource(r.NamespaceID, r.ModuleID, r.ID)
}

// RecordRbacResource returns string representation of RBAC resource for Record
//
// RBAC resource is in the corteza::compose:record/... format
//
// This function is auto-generated
func RecordRbacResource(NamespaceID uint64, ModuleID uint64, ID uint64) string {
	cpts := []interface{}{RecordResourceType}
	if NamespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(NamespaceID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if ModuleID != 0 {
		cpts = append(cpts, strconv.FormatUint(ModuleID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if ID != 0 {
		cpts = append(cpts, strconv.FormatUint(ID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(RecordRbacResourceTpl(), cpts...)

}

func RecordRbacResourceTpl() string {
	return "%s/%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Page by calling PageRbacResource fn
//
// RBAC resource is in the corteza::compose:page/... format
//
// This function is auto-generated
func (r Page) RbacResource() string {
	return PageRbacResource(r.NamespaceID, r.ID)
}

// PageRbacResource returns string representation of RBAC resource for Page
//
// RBAC resource is in the corteza::compose:page/... format
//
// This function is auto-generated
func PageRbacResource(NamespaceID uint64, ID uint64) string {
	cpts := []interface{}{PageResourceType}
	if NamespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(NamespaceID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if ID != 0 {
		cpts = append(cpts, strconv.FormatUint(ID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(PageRbacResourceTpl(), cpts...)

}

func PageRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Chart by calling ChartRbacResource fn
//
// RBAC resource is in the corteza::compose:chart/... format
//
// This function is auto-generated
func (r Chart) RbacResource() string {
	return ChartRbacResource(r.NamespaceID, r.ID)
}

// ChartRbacResource returns string representation of RBAC resource for Chart
//
// RBAC resource is in the corteza::compose:chart/... format
//
// This function is auto-generated
func ChartRbacResource(NamespaceID uint64, ID uint64) string {
	cpts := []interface{}{ChartResourceType}
	if NamespaceID != 0 {
		cpts = append(cpts, strconv.FormatUint(NamespaceID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if ID != 0 {
		cpts = append(cpts, strconv.FormatUint(ID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ChartRbacResourceTpl(), cpts...)

}

func ChartRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza::compose/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza::compose/ format
//
// This function is auto-generated
func ComponentRbacResource() string {
	return ComponentResourceType + "/"

}

func ComponentRbacResourceTpl() string {
	return "%s"
}
