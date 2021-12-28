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
	// Component struct serves as a virtual resource type for the system component
	//
	// This struct is auto-generated
	Component struct{}
)

var (
	_ = fmt.Printf
	_ = strconv.FormatUint
)

const (
	ApigwRouteResourceType  = "corteza::system:apigw-route"
	ApplicationResourceType = "corteza::system:application"
	AuthClientResourceType  = "corteza::system:auth-client"
	QueueResourceType       = "corteza::system:queue"
	ReportResourceType      = "corteza::system:report"
	RoleResourceType        = "corteza::system:role"
	TemplateResourceType    = "corteza::system:template"
	UserResourceType        = "corteza::system:user"
	ComponentResourceType   = "corteza::system"
)

// RbacResource returns string representation of RBAC resource for ApigwRoute by calling ApigwRouteRbacResource fn
//
// RBAC resource is in the corteza::system:apigw-route/... format
//
// This function is auto-generated
func (r ApigwRoute) RbacResource() string {
	return ApigwRouteRbacResource()
}

// ApigwRouteRbacResource returns string representation of RBAC resource for ApigwRoute
//
// RBAC resource is in the corteza::system:apigw-route/ format
//
// This function is auto-generated
func ApigwRouteRbacResource() string {
	return ApigwRouteResourceType + "/"

}

func ApigwRouteRbacResourceTpl() string {
	return "%s"
}

// RbacResource returns string representation of RBAC resource for Application by calling ApplicationRbacResource fn
//
// RBAC resource is in the corteza::system:application/... format
//
// This function is auto-generated
func (r Application) RbacResource() string {
	return ApplicationRbacResource()
}

// ApplicationRbacResource returns string representation of RBAC resource for Application
//
// RBAC resource is in the corteza::system:application/ format
//
// This function is auto-generated
func ApplicationRbacResource() string {
	return ApplicationResourceType + "/"

}

func ApplicationRbacResourceTpl() string {
	return "%s"
}

// RbacResource returns string representation of RBAC resource for AuthClient by calling AuthClientRbacResource fn
//
// RBAC resource is in the corteza::system:auth-client/... format
//
// This function is auto-generated
func (r AuthClient) RbacResource() string {
	return AuthClientRbacResource()
}

// AuthClientRbacResource returns string representation of RBAC resource for AuthClient
//
// RBAC resource is in the corteza::system:auth-client/ format
//
// This function is auto-generated
func AuthClientRbacResource() string {
	return AuthClientResourceType + "/"

}

func AuthClientRbacResourceTpl() string {
	return "%s"
}

// RbacResource returns string representation of RBAC resource for Queue by calling QueueRbacResource fn
//
// RBAC resource is in the corteza::system:queue/... format
//
// This function is auto-generated
func (r Queue) RbacResource() string {
	return QueueRbacResource()
}

// QueueRbacResource returns string representation of RBAC resource for Queue
//
// RBAC resource is in the corteza::system:queue/ format
//
// This function is auto-generated
func QueueRbacResource() string {
	return QueueResourceType + "/"

}

func QueueRbacResourceTpl() string {
	return "%s"
}

// RbacResource returns string representation of RBAC resource for Report by calling ReportRbacResource fn
//
// RBAC resource is in the corteza::system:report/... format
//
// This function is auto-generated
func (r Report) RbacResource() string {
	return ReportRbacResource()
}

// ReportRbacResource returns string representation of RBAC resource for Report
//
// RBAC resource is in the corteza::system:report/ format
//
// This function is auto-generated
func ReportRbacResource() string {
	return ReportResourceType + "/"

}

func ReportRbacResourceTpl() string {
	return "%s"
}

// RbacResource returns string representation of RBAC resource for Role by calling RoleRbacResource fn
//
// RBAC resource is in the corteza::system:role/... format
//
// This function is auto-generated
func (r Role) RbacResource() string {
	return RoleRbacResource()
}

// RoleRbacResource returns string representation of RBAC resource for Role
//
// RBAC resource is in the corteza::system:role/ format
//
// This function is auto-generated
func RoleRbacResource() string {
	return RoleResourceType + "/"

}

func RoleRbacResourceTpl() string {
	return "%s"
}

// RbacResource returns string representation of RBAC resource for Template by calling TemplateRbacResource fn
//
// RBAC resource is in the corteza::system:template/... format
//
// This function is auto-generated
func (r Template) RbacResource() string {
	return TemplateRbacResource()
}

// TemplateRbacResource returns string representation of RBAC resource for Template
//
// RBAC resource is in the corteza::system:template/ format
//
// This function is auto-generated
func TemplateRbacResource() string {
	return TemplateResourceType + "/"

}

func TemplateRbacResourceTpl() string {
	return "%s"
}

// RbacResource returns string representation of RBAC resource for User by calling UserRbacResource fn
//
// RBAC resource is in the corteza::system:user/... format
//
// This function is auto-generated
func (r User) RbacResource() string {
	return UserRbacResource()
}

// UserRbacResource returns string representation of RBAC resource for User
//
// RBAC resource is in the corteza::system:user/ format
//
// This function is auto-generated
func UserRbacResource() string {
	return UserResourceType + "/"

}

func UserRbacResourceTpl() string {
	return "%s"
}

// RbacResource returns string representation of RBAC resource for Component by calling ComponentRbacResource fn
//
// RBAC resource is in the corteza::system/... format
//
// This function is auto-generated
func (r Component) RbacResource() string {
	return ComponentRbacResource()
}

// ComponentRbacResource returns string representation of RBAC resource for Component
//
// RBAC resource is in the corteza::system/ format
//
// This function is auto-generated
func ComponentRbacResource() string {
	return ComponentResourceType + "/"

}

func ComponentRbacResourceTpl() string {
	return "%s"
}
