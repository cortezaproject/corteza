package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - system.application.yaml
// - system.auth-client.yaml
// - system.role.yaml
// - system.template.yaml
// - system.user.yaml
// - system.yaml

import (
	"strconv"
)

type (
	// Component struct serves as a virtual resource type for the system component
	//
	// This struct is auto-generated
	Component struct{}
)

const (
	ApplicationRbacResourceSchema = "corteza+system.application"
	AuthClientRbacResourceSchema  = "corteza+system.auth-client"
	RoleRbacResourceSchema        = "corteza+system.role"
	TemplateRbacResourceSchema    = "corteza+system.template"
	UserRbacResourceSchema        = "corteza+system.user"
	ComponentRbacResourceSchema   = "corteza+system"
)

// RbacResource returns string representation of RBAC resource for Application by calling ApplicationRbacResource fn
//
// RBAC resource is in the corteza+system.application:/... format
//
// This function is auto-generated
func (r Application) RbacResource() string {
	return ApplicationRbacResource(r.ID)
}

// ApplicationRbacResource returns string representation of RBAC resource for Application
//
// RBAC resource is in the corteza+system.application:/... format
//
// This function is auto-generated
func ApplicationRbacResource(ID uint64) string {
	out := ApplicationRbacResourceSchema + ":"
	out += "/"

	if ID != 0 {
		out += strconv.FormatUint(ID, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for AuthClient by calling AuthClientRbacResource fn
//
// RBAC resource is in the corteza+system.auth-client:/... format
//
// This function is auto-generated
func (r AuthClient) RbacResource() string {
	return AuthClientRbacResource(r.ID)
}

// AuthClientRbacResource returns string representation of RBAC resource for AuthClient
//
// RBAC resource is in the corteza+system.auth-client:/... format
//
// This function is auto-generated
func AuthClientRbacResource(ID uint64) string {
	out := AuthClientRbacResourceSchema + ":"
	out += "/"

	if ID != 0 {
		out += strconv.FormatUint(ID, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Role by calling RoleRbacResource fn
//
// RBAC resource is in the corteza+system.role:/... format
//
// This function is auto-generated
func (r Role) RbacResource() string {
	return RoleRbacResource(r.ID)
}

// RoleRbacResource returns string representation of RBAC resource for Role
//
// RBAC resource is in the corteza+system.role:/... format
//
// This function is auto-generated
func RoleRbacResource(ID uint64) string {
	out := RoleRbacResourceSchema + ":"
	out += "/"

	if ID != 0 {
		out += strconv.FormatUint(ID, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for Template by calling TemplateRbacResource fn
//
// RBAC resource is in the corteza+system.template:/... format
//
// This function is auto-generated
func (r Template) RbacResource() string {
	return TemplateRbacResource(r.ID)
}

// TemplateRbacResource returns string representation of RBAC resource for Template
//
// RBAC resource is in the corteza+system.template:/... format
//
// This function is auto-generated
func TemplateRbacResource(ID uint64) string {
	out := TemplateRbacResourceSchema + ":"
	out += "/"

	if ID != 0 {
		out += strconv.FormatUint(ID, 10)
	} else {
		out += "*"
	}
	return out
}

// RbacResource returns string representation of RBAC resource for User by calling UserRbacResource fn
//
// RBAC resource is in the corteza+system.user:/... format
//
// This function is auto-generated
func (r User) RbacResource() string {
	return UserRbacResource(r.ID)
}

// UserRbacResource returns string representation of RBAC resource for User
//
// RBAC resource is in the corteza+system.user:/... format
//
// This function is auto-generated
func UserRbacResource(ID uint64) string {
	out := UserRbacResourceSchema + ":"
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
// RBAC resource is in the corteza+system:/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza+system:/... format
//
// This function is auto-generated
func ComponentRbacResource() string {
	out := ComponentRbacResourceSchema + ":"
	return out
}
