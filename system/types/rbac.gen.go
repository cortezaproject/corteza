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
	AttachmentResourceType          = "corteza::system:attachment"
	ApplicationResourceType         = "corteza::system:application"
	ApigwRouteResourceType          = "corteza::system:apigw-route"
	ApigwFilterResourceType         = "corteza::system:apigw-filter"
	AuthClientResourceType          = "corteza::system:auth-client"
	AuthConfirmedClientResourceType = "corteza::system:auth-confirmed-client"
	AuthSessionResourceType         = "corteza::system:auth-session"
	AuthOa2tokenResourceType        = "corteza::system:auth-oa2token"
	CredentialResourceType          = "corteza::system:credential"
	QueueResourceType               = "corteza::system:queue"
	QueueMessageResourceType        = "corteza::system:queue_message"
	ReminderResourceType            = "corteza::system:reminder"
	ReportResourceType              = "corteza::system:report"
	ResourceTranslationResourceType = "corteza::system:resource-translation"
	RoleResourceType                = "corteza::system:role"
	RoleMemberResourceType          = "corteza::system:role_member"
	SettingValueResourceType        = "corteza::system:settings"
	TemplateResourceType            = "corteza::system:template"
	UserResourceType                = "corteza::system:user"
	ComponentResourceType           = "corteza::system"
)

// RbacResource returns string representation of RBAC resource for Attachment by calling AttachmentRbacResource fn
//
// RBAC resource is in the corteza::system:attachment/... format
//
// This function is auto-generated
func (r Attachment) RbacResource() string {
	return AttachmentRbacResource(r.ID)
}

// AttachmentRbacResource returns string representation of RBAC resource for Attachment
//
// RBAC resource is in the corteza::system:attachment/... format
//
// This function is auto-generated
func AttachmentRbacResource(id uint64) string {
	cpts := []interface{}{AttachmentResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(AttachmentRbacResourceTpl(), cpts...)

}

func AttachmentRbacResourceTpl() string {
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

func ApplicationRbacResourceTpl() string {
	return "%s/%s"
}

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

func ApigwRouteRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ApigwFilter by calling ApigwFilterRbacResource fn
//
// RBAC resource is in the corteza::system:apigw-filter/... format
//
// This function is auto-generated
func (r ApigwFilter) RbacResource() string {
	return ApigwFilterRbacResource(r.ID)
}

// ApigwFilterRbacResource returns string representation of RBAC resource for ApigwFilter
//
// RBAC resource is in the corteza::system:apigw-filter/... format
//
// This function is auto-generated
func ApigwFilterRbacResource(id uint64) string {
	cpts := []interface{}{ApigwFilterResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ApigwFilterRbacResourceTpl(), cpts...)

}

func ApigwFilterRbacResourceTpl() string {
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

func AuthClientRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for AuthConfirmedClient by calling AuthConfirmedClientRbacResource fn
//
// RBAC resource is in the corteza::system:auth-confirmed-client/... format
//
// This function is auto-generated
func (r AuthConfirmedClient) RbacResource() string {
	return AuthConfirmedClientRbacResource(r.ID)
}

// AuthConfirmedClientRbacResource returns string representation of RBAC resource for AuthConfirmedClient
//
// RBAC resource is in the corteza::system:auth-confirmed-client/... format
//
// This function is auto-generated
func AuthConfirmedClientRbacResource(id uint64) string {
	cpts := []interface{}{AuthConfirmedClientResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(AuthConfirmedClientRbacResourceTpl(), cpts...)

}

func AuthConfirmedClientRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for AuthSession by calling AuthSessionRbacResource fn
//
// RBAC resource is in the corteza::system:auth-session/... format
//
// This function is auto-generated
func (r AuthSession) RbacResource() string {
	return AuthSessionRbacResource(r.ID)
}

// AuthSessionRbacResource returns string representation of RBAC resource for AuthSession
//
// RBAC resource is in the corteza::system:auth-session/... format
//
// This function is auto-generated
func AuthSessionRbacResource(id uint64) string {
	cpts := []interface{}{AuthSessionResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(AuthSessionRbacResourceTpl(), cpts...)

}

func AuthSessionRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for AuthOa2token by calling AuthOa2tokenRbacResource fn
//
// RBAC resource is in the corteza::system:auth-oa2token/... format
//
// This function is auto-generated
func (r AuthOa2token) RbacResource() string {
	return AuthOa2tokenRbacResource(r.ID)
}

// AuthOa2tokenRbacResource returns string representation of RBAC resource for AuthOa2token
//
// RBAC resource is in the corteza::system:auth-oa2token/... format
//
// This function is auto-generated
func AuthOa2tokenRbacResource(id uint64) string {
	cpts := []interface{}{AuthOa2tokenResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(AuthOa2tokenRbacResourceTpl(), cpts...)

}

func AuthOa2tokenRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Credential by calling CredentialRbacResource fn
//
// RBAC resource is in the corteza::system:credential/... format
//
// This function is auto-generated
func (r Credential) RbacResource() string {
	return CredentialRbacResource(r.ID)
}

// CredentialRbacResource returns string representation of RBAC resource for Credential
//
// RBAC resource is in the corteza::system:credential/... format
//
// This function is auto-generated
func CredentialRbacResource(id uint64) string {
	cpts := []interface{}{CredentialResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(CredentialRbacResourceTpl(), cpts...)

}

func CredentialRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Queue by calling QueueRbacResource fn
//
// RBAC resource is in the corteza::system:queue/... format
//
// This function is auto-generated
func (r Queue) RbacResource() string {
	return QueueRbacResource(r.ID)
}

// QueueRbacResource returns string representation of RBAC resource for Queue
//
// RBAC resource is in the corteza::system:queue/... format
//
// This function is auto-generated
func QueueRbacResource(id uint64) string {
	cpts := []interface{}{QueueResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(QueueRbacResourceTpl(), cpts...)

}

func QueueRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for QueueMessage by calling QueueMessageRbacResource fn
//
// RBAC resource is in the corteza::system:queue_message/... format
//
// This function is auto-generated
func (r QueueMessage) RbacResource() string {
	return QueueMessageRbacResource(r.ID)
}

// QueueMessageRbacResource returns string representation of RBAC resource for QueueMessage
//
// RBAC resource is in the corteza::system:queue_message/... format
//
// This function is auto-generated
func QueueMessageRbacResource(id uint64) string {
	cpts := []interface{}{QueueMessageResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(QueueMessageRbacResourceTpl(), cpts...)

}

func QueueMessageRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for Reminder by calling ReminderRbacResource fn
//
// RBAC resource is in the corteza::system:reminder/... format
//
// This function is auto-generated
func (r Reminder) RbacResource() string {
	return ReminderRbacResource(r.ID)
}

// ReminderRbacResource returns string representation of RBAC resource for Reminder
//
// RBAC resource is in the corteza::system:reminder/... format
//
// This function is auto-generated
func ReminderRbacResource(id uint64) string {
	cpts := []interface{}{ReminderResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ReminderRbacResourceTpl(), cpts...)

}

func ReminderRbacResourceTpl() string {
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

func ReportRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for ResourceTranslation by calling ResourceTranslationRbacResource fn
//
// RBAC resource is in the corteza::system:resource-translation/... format
//
// This function is auto-generated
func (r ResourceTranslation) RbacResource() string {
	return ResourceTranslationRbacResource(r.ID)
}

// ResourceTranslationRbacResource returns string representation of RBAC resource for ResourceTranslation
//
// RBAC resource is in the corteza::system:resource-translation/... format
//
// This function is auto-generated
func ResourceTranslationRbacResource(id uint64) string {
	cpts := []interface{}{ResourceTranslationResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(ResourceTranslationRbacResourceTpl(), cpts...)

}

func ResourceTranslationRbacResourceTpl() string {
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

func RoleRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for RoleMember by calling RoleMemberRbacResource fn
//
// RBAC resource is in the corteza::system:role_member/... format
//
// This function is auto-generated
func (r RoleMember) RbacResource() string {
	return RoleMemberRbacResource(r.ID)
}

// RoleMemberRbacResource returns string representation of RBAC resource for RoleMember
//
// RBAC resource is in the corteza::system:role_member/... format
//
// This function is auto-generated
func RoleMemberRbacResource(id uint64) string {
	cpts := []interface{}{RoleMemberResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(RoleMemberRbacResourceTpl(), cpts...)

}

func RoleMemberRbacResourceTpl() string {
	return "%s/%s"
}

// RbacResource returns string representation of RBAC resource for SettingValue by calling SettingValueRbacResource fn
//
// RBAC resource is in the corteza::system:settings/... format
//
// This function is auto-generated
func (r SettingValue) RbacResource() string {
	return SettingValueRbacResource(r.ID)
}

// SettingValueRbacResource returns string representation of RBAC resource for SettingValue
//
// RBAC resource is in the corteza::system:settings/... format
//
// This function is auto-generated
func SettingValueRbacResource(id uint64) string {
	cpts := []interface{}{SettingValueResourceType}
	if id != 0 {
		cpts = append(cpts, strconv.FormatUint(id, 10))
	} else {
		cpts = append(cpts, "*")
	}

	return fmt.Sprintf(SettingValueRbacResourceTpl(), cpts...)

}

func SettingValueRbacResourceTpl() string {
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

func ComponentRbacResourceTpl() string {
	return "%s"
}
