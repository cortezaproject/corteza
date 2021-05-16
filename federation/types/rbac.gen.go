package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - federation.exposed-module.yaml
// - federation.node.yaml
// - federation.shared-module.yaml
// - federation.yaml

import (
	"strconv"
)

type (
	// Component struct serves as a virtual resource type for the federation component
	//
	// This struct is auto-generated
	Component struct{}
)

const (
	ExposedModuleRbacResourceSchema = "corteza+federation.exposed-module"
	NodeRbacResourceSchema          = "corteza+federation.node"
	SharedModuleRbacResourceSchema  = "corteza+federation.shared-module"
	ComponentRbacResourceSchema     = "corteza+federation"
)

// RbacResource returns string representation of RBAC resource for ExposedModule by calling ExposedModuleRbacResource fn
//
// RBAC resource is in the corteza+federation.exposed-module:/... format
//
// This function is auto-generated
func (r ExposedModule) RbacResource() string {
	return ExposedModuleRbacResource(r.NodeID, r.ID)
}

// ExposedModuleRbacResource returns string representation of RBAC resource for ExposedModule
//
// RBAC resource is in the corteza+federation.exposed-module:/... format
//
// This function is auto-generated
func ExposedModuleRbacResource(NodeID uint64, ID uint64) string {
	out := ExposedModuleRbacResourceSchema + ":"
	out += "/"

	if NodeID != 0 {
		out += strconv.FormatUint(NodeID, 10)
	} else {
		out += "*"
	}
	out += "/"

	if ID != 0 {
		out += strconv.FormatUint(ID, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Node by calling NodeRbacResource fn
//
// RBAC resource is in the corteza+federation.node:/... format
//
// This function is auto-generated
func (r Node) RbacResource() string {
	return NodeRbacResource(r.ID)
}

// NodeRbacResource returns string representation of RBAC resource for Node
//
// RBAC resource is in the corteza+federation.node:/... format
//
// This function is auto-generated
func NodeRbacResource(ID uint64) string {
	out := NodeRbacResourceSchema + ":"
	out += "/"

	if ID != 0 {
		out += strconv.FormatUint(ID, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for SharedModule by calling SharedModuleRbacResource fn
//
// RBAC resource is in the corteza+federation.shared-module:/... format
//
// This function is auto-generated
func (r SharedModule) RbacResource() string {
	return SharedModuleRbacResource(r.NodeID, r.ID)
}

// SharedModuleRbacResource returns string representation of RBAC resource for SharedModule
//
// RBAC resource is in the corteza+federation.shared-module:/... format
//
// This function is auto-generated
func SharedModuleRbacResource(NodeID uint64, ID uint64) string {
	out := SharedModuleRbacResourceSchema + ":"
	out += "/"

	if NodeID != 0 {
		out += strconv.FormatUint(NodeID, 10)
	} else {
		out += "*"
	}
	out += "/"

	if ID != 0 {
		out += strconv.FormatUint(ID, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza+federation:/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza+federation:/... format
//
// This function is auto-generated
func ComponentRbacResource() string {
	out := ComponentRbacResourceSchema + ":"
	return out
}
