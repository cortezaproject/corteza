package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

RBAC: schema.#optionsGroup & {
	handle: "rbac"
	title:  "RBAC options"

	options: {
		log: {
			type:        "bool"
			description: "Log RBAC related events and actions"
		}
		service_user: {}
		bypass_roles: {
			defaultValue: "super-admin"
			description: """
				Space delimited list of role handles.
				These roles causes short-circuiting access control check and allowing all operations.
				System will refuse to start if check-bypassing roles are also listed as authenticated or anonymous auto-assigned roles.
				"""
		}
		authenticated_roles: {
			defaultValue: "authenticated"
			description: """
				Space delimited list of role handles.
				These roles are automatically assigned to authenticated user.
				Memberships can not be managed for these roles.
				System will refuse to start if roles listed here are also listed under anonymous roles
				"""
		}
		anonymous_roles: {
			defaultValue: "anonymous"
			description: """
				Space delimited list of role handles.
				These roles are automatically assigned to anonymous user.
				Memberships can not be managed for these roles.
				"""
		}
	}
}
