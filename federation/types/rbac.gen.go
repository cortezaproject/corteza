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
	// Component struct serves as a virtual resource type for the federation component
	//
	// This struct is auto-generated
	Component struct{}
)

var (
	_ = fmt.Printf
	_ = strconv.FormatUint
)

const (
	NodeResourceType          = "corteza::federation:node"
	NodeSyncResourceType      = "corteza::federation:node-sync"
	ExposedModuleResourceType = "corteza::federation:exposed-module"
	SharedModuleResourceType  = "corteza::federation:shared-module"
	ModuleMappingResourceType = "corteza::federation:module-mapping"
	ComponentResourceType     = "corteza::federation"
)

// RbacResource returns string representation of RBAC resource for Node by calling NodeRbacResource fn
//
// RBAC resource is in the corteza::federation:node/... format
//
// This function is auto-generated
func (r Node) RbacResource() string {
	return NodeRbacResource(r.ID)
}

// NodeRbacResource returns string representation of RBAC resource for Node
//
// RBAC resource is in the corteza::federation:node/... format
//
// This function is auto-generated
func NodeRbacResource(id uint64) string {
	cpts := []interface{}{NodeResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(NodeRbacResourceTpl(), cpts...)

}

func NodeRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for NodeSync by calling NodeSyncRbacResource fn
//
// RBAC resource is in the corteza::federation:node-sync/... format
//
// This function is auto-generated
func (r NodeSync) RbacResource() string {
	return NodeSyncRbacResource(r.ID)
}

// NodeSyncRbacResource returns string representation of RBAC resource for NodeSync
//
// RBAC resource is in the corteza::federation:node-sync/... format
//
// This function is auto-generated
func NodeSyncRbacResource(id uint64) string {
	cpts := []interface{}{NodeSyncResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(NodeSyncRbacResourceTpl(), cpts...)

}

func NodeSyncRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ExposedModule by calling ExposedModuleRbacResource fn
//
// RBAC resource is in the corteza::federation:exposed-module/... format
//
// This function is auto-generated
func (r ExposedModule) RbacResource() string {
	return ExposedModuleRbacResource(r.NodeID, r.ID)
}

// ExposedModuleRbacResource returns string representation of RBAC resource for ExposedModule
//
// RBAC resource is in the corteza::federation:exposed-module/... format
//
// This function is auto-generated
func ExposedModuleRbacResource(nodeID uint64, id uint64) string {
	cpts := []interface{}{ExposedModuleResourceType}
	if nodeID != 0 {
		cpts = append(cpts, strconv.FormatUint(nodeID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ExposedModuleRbacResourceTpl(), cpts...)

}

func ExposedModuleRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for SharedModule by calling SharedModuleRbacResource fn
//
// RBAC resource is in the corteza::federation:shared-module/... format
//
// This function is auto-generated
func (r SharedModule) RbacResource() string {
	return SharedModuleRbacResource(r.NodeID, r.ID)
}

// SharedModuleRbacResource returns string representation of RBAC resource for SharedModule
//
// RBAC resource is in the corteza::federation:shared-module/... format
//
// This function is auto-generated
func SharedModuleRbacResource(nodeID uint64, id uint64) string {
	cpts := []interface{}{SharedModuleResourceType}
	if nodeID != 0 {
		cpts = append(cpts, strconv.FormatUint(nodeID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(SharedModuleRbacResourceTpl(), cpts...)

}

func SharedModuleRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for ModuleMapping by calling ModuleMappingRbacResource fn
//
// RBAC resource is in the corteza::federation:module-mapping/... format
//
// This function is auto-generated
func (r ModuleMapping) RbacResource() string {
	return ModuleMappingRbacResource(r.NodeID, r.ID)
}

// ModuleMappingRbacResource returns string representation of RBAC resource for ModuleMapping
//
// RBAC resource is in the corteza::federation:module-mapping/... format
//
// This function is auto-generated
func ModuleMappingRbacResource(nodeID uint64, id uint64) string {
	cpts := []interface{}{ModuleMappingResourceType}
	if nodeID != 0 {
		cpts = append(cpts, strconv.FormatUint(nodeID, 10))
	} else {
		cpts = append(cpts, "*")
	}

	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ModuleMappingRbacResourceTpl(), cpts...)

}

func ModuleMappingRbacResourceTpl() string {
	return "%s/%s/%s"
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza::federation/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza::federation/ format
//
// This function is auto-generated
func ComponentRbacResource() string {
	return ComponentResourceType + "/"

}

func ComponentRbacResourceTpl() string {
	return "%s"
}
