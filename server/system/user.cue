package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_userDefs: {
	UserKind: {
		name: "UserKind"
		values: [
			{ident: "NormalUser", value: ""},
			{ident: "SystemUser", value: "sys"},
		]
	}
	UserMeta: {
		name: "UserMeta"
		fields: [
			{name: "AvatarID", type:          "uint64", json:                         "avatarID,string"},
			{name: "AvatarKind", type:        "string", json:                         "avatarKind,omitempty"},
			{name: "AvatarColor", type:       "string", json:                         "avatarColor,omitempty"},
			{name: "AvatarBgColor", type:     "string", json:                         "avatarBgColor,omitempty"},
			{name: "PreferredLanguage", type: "string", json:                         "preferredLanguage"},
			{name: "Theme", type:             "string", json:                         "theme"},
			{name: "SecurityPolicy", type:    _userDefs.UserMetaSecurityPolicy, json: "securityPolicy"},
		]
	}
	UserMetaSecurityPolicy: {
		name: "UserMetaSecurityPolicy"
		fields: [
			{name: "MFA", type: _userDefs.UserMetaSecurityPolicyMFA, json: "mfa"},
		]
	}
	UserMetaSecurityPolicyMFA: {
		name: "UserMetaSecurityPolicyMFA"
		fields: [
			{name: "EnforcedEmailOTP", type: "bool", json: "enforcedEmailOTP"},
			{name: "EnforcedTOTP", type:     "bool", json: "enforcedTOTP"},
		]
	}
	UserMetrics: {name: "UserMetrics", fields: [
				{name: "Total", type:          "uint", json:  "total"},
				{name: "Valid", type:          "uint", json:  "valid"},
				{name: "Deleted", type:        "uint", json:  "deleted"},
				{name: "Suspended", type:      "uint", json:  "suspended"},
				{name: "DailyCreated", type:   "uint", slice: true, json: "dailyCreated"},
				{name: "DailyDeleted", type:   "uint", slice: true, json: "dailyDeleted"},
				{name: "DailyUpdated", type:   "uint", slice: true, json: "dailyUpdated"},
				{name: "DailySuspended", type: "uint", slice: true, json: "dailySuspended"},
	]}
}

user: {
	features: {
		labelResourceType: "user"
	}

	model: {
		attributes: {
			id:         schema.IdField
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			email: {
				sortable:   true
				unique:     true
				ignoreCase: true
				dal: {length: 254}
			}
			email_confirmed: {
				goType: "bool"
				dal: {type: "Boolean"}
			}
			user_group_id: {
				ident:      "userGroupID"
				goType:     "uint64"
				storeIdent: "rel_user_group"
				dal: {
					type:            "Ref"
					refModelResType: "corteza::system:user-group"
					nullable:        true
					default:         0
				}

				envoy: {
					yaml: {
						identKeyAlias: ["usergroup", "user_group", "user_group_id", "group"]
					}
				}
			}
			username: {
				sortable:   true
				unique:     true
				ignoreCase: true
				dal: {}
			}
			roles: {
				goType:     "[]uint64"
				store:      false
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						customDecoder: true
						omitEncoder:   true
					}
				}
			}
			name: {
				sortable: true
				dal: {}
			}
			handle: schema.HandleField
			kind: {
				sortable: true
				type:     _userDefs.UserKind
				dal: {length: 8}
				omitSetter: true
				omitGetter: true
			}
			meta: {
				type: _userDefs.UserMeta
				ptr:  true
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			suspended_at: schema.SortableTimestampNilField
			created_at:   schema.SortableTimestampNowField
			updated_at:   schema.SortableTimestampNilField
			deleted_at:   schema.SortableTimestampNilField
		}

		indexes: {
			"primary": {attribute: "id"}
			"unique_email": {
				fields: [{attribute: "email", modifiers: ["LOWERCASE"]}]
				predicate: "email != '' AND deleted_at IS NULL"
			}
			"unique_handle": {
				fields: [{attribute: "handle", modifiers: ["LOWERCASE"]}]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
			"unique_username": {
				fields: [{attribute: "username", modifiers: ["LOWERCASE"]}]
				predicate: "username != '' AND deleted_at IS NULL"
			}
		}
	}

	types: {
		defs: _userDefs
	}

	filter: {
		struct: {
			user_id: {goType: "[]uint64", ident: "userID", storeIdent: "id"}
			tenant_id:  schema.TenantFilterField
			project_id: schema.ProjectFilterField
			role_id: {goType: "[]uint64", ident: "roleID"}
			user_group_id: {goType: "uint64", ident: "userGroupID"}
			email: {goType: "string"}
			name: {goType: "string"}
			username: {goType: "string"}
			handle: {goType: "string"}
			kind: {goType: "types.UserKind"}
			allKinds: {goType: "bool"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
			suspended: {goType: "filter.State", storeIdent: "suspended_at"}
		}

		query: ["email", "username", "handle", "name"]
		byValue: ["user_id", "email", "username", "handle"]
		byNilState: ["deleted", "suspended"]
	}

	envoy: {
		yaml: {
			supportMappedInput: true
			mappedField:        "Handle"
			identKeyAlias: ["users", "usr"]
		}
		store: {}
	}

	rbac: {
		operations: {
			"read": description:               "Read user"
			"update": description:             "Update user"
			"delete": description:             "Delete user"
			"suspend": description:            "Suspend user"
			"unsuspend": description:          "Unsuspend user"
			"email.unmask": description:       "Unmask email"
			"name.unmask": description:        "Unmask name"
			"impersonate": description:        "Impersonate user"
			"credentials.manage": description: "Manage user's credentials"
		}
	}

	service: {

		undelete: true

		customBodyOps: ["lookup", "search", "create", "update", "delete", "undelete"]

		customFunctions: [
			{
				name:   "FindByEmail"
				cap:    "read"
				action: "Lookup"
				args: [{name: "email", goType: "string"}]
				results: [{name: "u", goType: "*types.User"}, {name: "err", goType: "error"}]
			},
			{
				name:   "FindByHandle"
				cap:    "read"
				action: "Lookup"
				args: [{name: "handle", goType: "string"}]
				results: [{name: "u", goType: "*types.User"}, {name: "err", goType: "error"}]
			},
			{
				name:   "ToggleEmailConfirmation"
				cap:    "write"
				action: "Update"
				args: [{name: "userID", goType: "uint64"}, {name: "confirmed", goType: "bool"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "Suspend"
				cap:    "write"
				action: "Suspend"
				args: [{name: "userID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "Unsuspend"
				cap:    "write"
				action: "Unsuspend"
				args: [{name: "userID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "SetPassword"
				cap:    "write"
				action: "SetPassword"
				args: [{name: "userID", goType: "uint64"}, {name: "newPassword", goType: "string"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "DeleteAuthTokensByUserID"
				cap:    "write"
				action: "DeleteAuthTokens"
				args: [{name: "userID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "DeleteAuthSessionsByUserID"
				cap:    "write"
				action: "DeleteAuthSessions"
				args: [{name: "userID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:  "CreateSynthetic"
				cap:   "write"
				ac:    "CanCreateUser"
				acErr: "ErrNotAllowedToCreate"
				args: [{name: "src", goType: "synteticUserDataGen"}, {name: "total", goType: "uint"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:  "RemoveSynthetic"
				cap:   "write"
				ac:    "CanCreateUser"
				acErr: "ErrNotAllowedToCreate"
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "UploadAvatar"
				cap:    "write"
				action: "UploadAvatar"
				args: [{name: "userID", goType: "uint64"}, {name: "upload", goType: "*multipart.FileHeader"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "DeleteAvatar"
				cap:    "write"
				action: "DeleteAvatar"
				args: [{name: "userID", goType: "uint64"}]
				results: [{name: "err", goType: "error"}]
			},
			{
				name:   "GenerateAvatar"
				cap:    "write"
				action: "GenerateAvatar"
				args: [{name: "userID", goType: "uint64"}, {name: "bgColor", goType: "string"}, {name: "initialColor", goType: "string"}]
				results: [{name: "err", goType: "error"}]
			},
		]
		customFunctionImports: ["\"mime/multipart\""]
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for user by ID

						It returns user even if deleted or suspended
						"""
				}, {
					fields: ["email"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for user by email

						It returns only valid user (not deleted, not suspended)
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for user by handle

						It returns only valid user (not deleted, not suspended)
						"""
				}, {
					fields: ["username"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for user by username

						It returns only valid user (not deleted, not suspended)
						"""
				},
			]

			functions: [
				{
					expIdent: "CountUsers"
					args: [ {ident: "u", goType: "types.UserFilter"}]
					return: [ "uint"]
				}, {
					expIdent: "UserMetrics"
					return: [ "*types.UserMetrics"]
				},
			]
		}
	}
}
