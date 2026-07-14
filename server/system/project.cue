package system

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_projectDefs: {
	ProjectStatus: {name: "ProjectStatus", values: [
				{ident: "ProjectStatusDraft", value:      "draft"},
				{ident: "ProjectStatusActive", value:     "active"},
				{ident: "ProjectStatusPublished", value:  "published"},
				{ident: "ProjectStatusArchived", value:   "archived"},
				{ident: "ProjectStatusSuspended", value:  "suspended"},
				{ident: "ProjectStatusDeprecated", value: "deprecated"},
	]}

	ProjectVisibility: {name: "ProjectVisibility", values: [
					{ident: "ProjectVisibilityOpen", value:       "open"},
					{ident: "ProjectVisibilityInviteOnly", value: "invite-only"},
	]}

	ProjectMode: {name: "ProjectMode", values: [
				{ident: "ProjectModeFree", value:  "free"},
				{ident: "ProjectModeGated", value: "gated"},
	]}

	ProjectMemberRole: {name: "ProjectMemberRole", values: [
					{ident: "ProjectRoleGovernanceOwner", value:             "governance-owner"},
					{ident: "ProjectRoleSecurityOwner", value:               "security-owner"},
					{ident: "ProjectRoleDeveloper", value:                   "developer"},
					{ident: "ProjectRoleJuniorDeveloper", value:             "junior-developer"},
					{ident: "ProjectRoleMember", value:                      "member"},
					{ident: "ProjectRoleExecutiveAuthority", value:          "executive-authority"},
					{ident: "ProjectRoleInfrastructureAdministrator", value: "infrastructure-administrator"},
	]}

	ProjectConfig: {name: "ProjectConfig", fields: [
				{name: "Visibility", type:         _projectDefs.ProjectVisibility, json:         "visibility,omitempty"},
				{name: "DefaultMemberRole", type:  _projectDefs.ProjectMemberRole, json:         "defaultMemberRole,omitempty"},
				{name: "FeatureFlags", goType:     "map[string]bool", json:                      "featureFlags,omitempty"},
				{name: "NamespaceID", type:        "uint64", json:                               "namespaceID,string,omitempty"},
				{name: "DeployerCategories", type: _projectDefs.ProjectDeployerCategories, json: "deployerCategories,omitempty"},
				{name: "FriaRequired", type:       "bool", json:                                 "friaRequired,omitempty"},
				{name: "ResourceManagement", type: _projectDefs.ProjectResourceManagement, json: "resourceManagement,omitempty"},
	]}

	ProjectDeployerCategories: {name: "ProjectDeployerCategories", fields: [
						{name: "PublicAuthorityAnnex3", type:    "bool", json: "publicAuthorityAnnex3,omitempty"},
						{name: "PrivateEssentialServices", type: "bool", json: "privateEssentialServices,omitempty"},
						{name: "InsuranceBanking", type:         "bool", json: "insuranceBanking,omitempty"},
	]}

	ProjectResourceManagement: {name: "ProjectResourceManagement", fields: [
						{name: "AI", goType:          "map[string]any", json:                "ai,omitempty"},
						{name: "Infra", goType:       "map[string]any", json:                "infra,omitempty"},
						{name: "Connections", goType: "[]*ProjectPermittedConnection", json: "connections,omitempty"},
	]}

	ProjectMeta: {name: "ProjectMeta", fields: [
				{name: "Short", type:       "string", json: "short,omitempty"},
				{name: "Description", type: "string", json: "description,omitempty"},
				{name: "Icon", type:        "string", json: "icon,omitempty"},
				{name: "Color", type:       "string", json: "color,omitempty"},
				{name: "Tags", slice:       true, type:     "string", json: "tags,omitempty"},
	]}

	ProjectGovernanceStatus: {name: "ProjectGovernanceStatus", values: [
					{ident: "ProjectGovernanceStatusDraft", value:            "draft"},
					{ident: "ProjectGovernanceStatusSubmitted", value:        "submitted"},
					{ident: "ProjectGovernanceStatusApproved", value:         "approved"},
					{ident: "ProjectGovernanceStatusChangesRequested", value: "changes-requested"},
	]}

	ProjectGovernanceStep: {name: "ProjectGovernanceStep", fields: [
					{name: "Values", goType:   "map[string]any", json:                     "values,omitempty"},
					{name: "Status", type:     _projectDefs.ProjectGovernanceStatus, json: "status"},
					{name: "ReviewNote", type: "string", json:                             "reviewNote,omitempty"},
	]}

	ProjectGovernance: {name: "ProjectGovernance", key: "string", value: _projectDefs.ProjectGovernanceStep, valuePtr: true}

	ProjectPermittedConnection: {name: "ProjectPermittedConnection", fields: [
						{name: "ID", type:                  "string", json: "id"},
						{name: "Name", type:                "string", json: "name"},
						{name: "Connector", type:           "string", json: "connector,omitempty"},
						{name: "Type", type:                "string", json: "type,omitempty"},
						{name: "Description", type:         "string", json: "description,omitempty"},
						{name: "ActionIfUnavailable", type: "string", json: "actionIfUnavailable,omitempty"},
						{name: "Replacement", type:         "string", json: "replacement,omitempty"},
						{name: "IsAiSystem", type:          "string", json: "isAiSystem,omitempty"},
	]}
}

project: {
	features: {
		labels:            true
		labelResourceType: "project"
	}

	types: {
		gen:  true
		defs: _projectDefs
	}

	model: {
		attributes: {
			id:        schema.IdField
			tenant_id: schema.TenantRefField
			handle:    schema.HandleField
			status: {
				sortable: true
				type:     _projectDefs.ProjectStatus
				dal: {length: 32}
				omitSetter: true
				omitGetter: true
			}
			config: {
				type: _projectDefs.ProjectConfig
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			// Top-level, sortable. Mode is chosen at creation and immutable after;
			// it lives here (not in config JSON) so it's a plain sortable column.
			mode: {
				sortable: true
				type:     _projectDefs.ProjectMode
				dal: {length: 32}
				omitSetter: true
				omitGetter: true
			}
			meta: {
				type: _projectDefs.ProjectMeta
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			governance: {
				type: _projectDefs.ProjectGovernance
				dal: {type: "JSON", defaultEmptyObject: true}
				omitSetter: true
				omitGetter: true
			}
			// Revision chain: root project ID, immediate parent revision, and the
			// revision number. Column names are explicit (root_project_id) and differ
			// from the exported idents; no getters/setters are generated.
			root_project_id: {
				ident:      "projectID"
				expIdent:   "ProjectID"
				goType:     "uint64"
				storeIdent: "root_project_id"
				json:       "rootProjectID,string,omitempty"
				dal: {type: "ID", default: 0}
				omitSetter: true
				omitGetter: true
			}
			parent_revision_id: {
				ident:      "parentRevisionID"
				expIdent:   "ParentRevisionID"
				goType:     "uint64"
				storeIdent: "parent_revision_id"
				json:       "parentRevisionID,string,omitempty"
				dal: {type: "ID", default: 0}
				omitSetter: true
				omitGetter: true
			}
			revision: {
				ident:      "revision"
				expIdent:   "Revision"
				goType:     "int"
				storeIdent: "revision"
				json:       "revision,omitempty"
				dal: {type: "Number", default: 0, precision: 0, scale: 0}
				omitSetter: true
				omitGetter: true
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by: schema.AttributeUserRef & {
				json: {field: "createdBy", string: true}
			}
			updated_by: schema.AttributeUserRef
			deleted_by: schema.AttributeUserRef
		}

		indexes: {
			"primary": {attribute: "id"}
			"unique_handle": {
				fields: [{attribute: "handle", modifiers: ["LOWERCASE"]}]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	filter: {
		struct: {
			project_id: {goType: "[]uint64", ident: "projectID", storeIdent: "id"}
			root_project_id: {goType: "uint64", ident: "rootProjectID", storeIdent: "root_project_id"}
			tenant_id: schema.TenantFilterField
			handle: {goType: "string"}
			status: {goType: "types.ProjectStatus"}

			deleted: {goType: "filter.State", storeIdent: "deleted_at"}
		}

		query: ["handle"]
		byValue: ["project_id", "root_project_id", "handle"]
		byNilState: ["deleted"]
	}

	rbac: {
		operations: {
			read: description:             "Read project"
			update: description:           "Update project"
			delete: description:           "Delete project"
			"members.manage": description: "Manage project members"
		}
	}

	envoy: {
		omit: true
	}

	service: {
		customFunctions: [
			{
				name: "SearchMembers"
				cap:  "read"
				args: [{name: "filter", goType: "types.ProjectMemberFilter"}]
				results: [
					{name: "set", goType: "types.ProjectMemberSet"},
					{name: "f", goType:   "types.ProjectMemberFilter"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "AddMember"
				cap:  "write"
				args: [{name: "m", goType: "*types.ProjectMember"}]
				results: [
					{name: "res", goType: "*types.ProjectMember"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "UpdateMember"
				cap:  "write"
				args: [{name: "m", goType: "*types.ProjectMember"}]
				results: [
					{name: "res", goType: "*types.ProjectMember"},
					{name: "err", goType: "error"},
				]
			},
			{
				name: "RemoveMember"
				cap:  "write"
				args: [
					{name: "projectID", goType: "uint64"},
					{name: "userID", goType:    "uint64"},
				]
				results: [
					{name: "err", goType: "error"},
				]
			},
		]

		// The project struct carries extra revision dependencies (nsSvc/dalSvc/
		// dalConns/recordSvc), so the struct + constructor are hand-written in
		// project.go rather than generated (genConstructor is only safe for
		// dependency-free services).

		// projectServices is hand-written in project.go with all extra deps.

		// Create/delete/undelete delegate to hand-written on<Op> handlers
		// (namespace cascade, member seeding, Tx).  Update uses beforeUpdate.
		// Lookup uses afterLookup for label loading.
		lookup:   true
		create:   true
		update:   true
		delete:   true
		undelete: true

		customBodyOps: ["create", "delete", "undelete"]

		// Explicit list because status/config/meta carry omitSetter:true (a DAL
		// concern) which would wrongly exclude them from the derived settable list.
		// updated_by is included so beforeUpdate can stamp the caller identity.
		updateFields: ["Handle", "Status", "Config", "Meta", "UpdatedBy"]

		hooks: {
			afterLookup:  true
			beforeUpdate: true
		}
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["id"]
					description: """
						searches for project by ID

						It also returns deleted projects.
						"""
				}, {
					fields: ["handle"]
					nullConstraint: ["deleted_at"]
					constraintCheck: true
					description: """
						searches for project by handle

						It returns only valid projects (not deleted)
						"""
				},
			]
		}
	}
}
