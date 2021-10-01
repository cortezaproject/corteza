package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

// Definitions file that controls how this file is generated:
// - system.apigw-route.yaml
// - system.application.yaml
// - system.auth-client.yaml
// - system.report.yaml
// - system.role.yaml
// - system.template.yaml
// - system.user.yaml
// - system.yaml

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

const (
	ApigwRouteResourceType  = "corteza::system:apigw-route"
	ApplicationResourceType = "corteza::system:application"
	AuthClientResourceType  = "corteza::system:auth-client"
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
	return ApigwRouteRbacResource(r.ID)
}

// ApigwRouteRbacResource returns string representation of RBAC resource for ApigwRoute
//
// RBAC resource is in the corteza::system:apigw-route/... format
//
// This function is auto-generated
func ApigwRouteRbacResource(id uint64) string {
	cpts := []interface{}{ApigwRouteResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ApigwRouteRbacResourceTpl(), cpts...)

}

// @todo template
func ApigwRouteRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Application by calling ApplicationRbacResource fn
//
// RBAC resource is in the corteza::system:application/... format
//
// This function is auto-generated
func (r Application) RbacResource() string {
	return ApplicationRbacResource(r.ID)
}

// ApplicationRbacResource returns string representation of RBAC resource for Application
//
// RBAC resource is in the corteza::system:application/... format
//
// This function is auto-generated
func ApplicationRbacResource(id uint64) string {
	cpts := []interface{}{ApplicationResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ApplicationRbacResourceTpl(), cpts...)

}

// @todo template
func ApplicationRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for AuthClient by calling AuthClientRbacResource fn
//
// RBAC resource is in the corteza::system:auth-client/... format
//
// This function is auto-generated
func (r AuthClient) RbacResource() string {
	return AuthClientRbacResource(r.ID)
}

// AuthClientRbacResource returns string representation of RBAC resource for AuthClient
//
// RBAC resource is in the corteza::system:auth-client/... format
//
// This function is auto-generated
func AuthClientRbacResource(id uint64) string {
	cpts := []interface{}{AuthClientResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(AuthClientRbacResourceTpl(), cpts...)

}

// @todo template
func AuthClientRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Report by calling ReportRbacResource fn
//
// RBAC resource is in the corteza::system:report/... format
//
// This function is auto-generated
func (r Report) RbacResource() string {
	return ReportRbacResource(r.ID)
}

// ReportRbacResource returns string representation of RBAC resource for Report
//
// RBAC resource is in the corteza::system:report/... format
//
// This function is auto-generated
func ReportRbacResource(id uint64) string {
	cpts := []interface{}{ReportResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ReportRbacResourceTpl(), cpts...)

}

// @todo template
func ReportRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Role by calling RoleRbacResource fn
//
// RBAC resource is in the corteza::system:role/... format
//
// This function is auto-generated
func (r Role) RbacResource() string {
	return RoleRbacResource(r.ID)
}

// RoleRbacResource returns string representation of RBAC resource for Role
//
// RBAC resource is in the corteza::system:role/... format
//
// This function is auto-generated
func RoleRbacResource(id uint64) string {
	cpts := []interface{}{RoleResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(RoleRbacResourceTpl(), cpts...)

}

// @todo template
func RoleRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Template by calling TemplateRbacResource fn
//
// RBAC resource is in the corteza::system:template/... format
//
// This function is auto-generated
func (r Template) RbacResource() string {
	return TemplateRbacResource(r.ID)
}

// TemplateRbacResource returns string representation of RBAC resource for Template
//
// RBAC resource is in the corteza::system:template/... format
//
// This function is auto-generated
func TemplateRbacResource(id uint64) string {
	cpts := []interface{}{TemplateResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(TemplateRbacResourceTpl(), cpts...)

}

// @todo template
func TemplateRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for User by calling UserRbacResource fn
//
// RBAC resource is in the corteza::system:user/... format
//
// This function is auto-generated
func (r User) RbacResource() string {
	return UserRbacResource(r.ID)
}

// UserRbacResource returns string representation of RBAC resource for User
//
// RBAC resource is in the corteza::system:user/... format
//
// This function is auto-generated
func UserRbacResource(id uint64) string {
	cpts := []interface{}{UserResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(UserRbacResourceTpl(), cpts...)

}

// @todo template
func UserRbacResourceTpl() string {
	return "%s/%s"
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

// @todo template
func ComponentRbacResourceTpl() string {
	return "%s"
}
