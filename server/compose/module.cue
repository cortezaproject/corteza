package compose

import (
	"github.com/crusttech/human/server/codegen/schema"
)

_moduleDefs: {
			ModuleConfig: { name: "ModuleConfig", fields: [
				{ name: "DAL", type: _moduleDefs.ModuleConfigDAL, json: "dal" },
				{ name: "Privacy", type: _moduleDefs.ModuleConfigDataPrivacy, json: "privacy" },
				{ name: "Discovery", type: _moduleDefs.ModuleConfigDiscovery, json: "discovery" },
				{ name: "RecordRevisions", type: _moduleDefs.ModuleConfigRecordRevisions, json: "recordRevisions" },
				{ name: "RecordDeDup", type: _moduleDefs.ModuleConfigRecordDeDup, json: "recordDeDup" },
			]}
			ModuleConfigDAL: { name: "ModuleConfigDAL", fields: [
				{ name: "ConnectionID", type: "uint64", json: "connectionID,string" },
				{ name: "Constraints", goType: "map[string][]any", json: "constraints" },
				{ name: "Ident", type: "string", json: "ident" },
				{ name: "SystemFieldEncoding", type: _moduleDefs.SystemFieldEncoding, json: "systemFieldEncoding" },
			]}
			ModuleConfigDataPrivacy: { name: "ModuleConfigDataPrivacy", fields: [
				{ name: "SensitivityLevelID", type: "uint64", json: "sensitivityLevelID,string,omitempty" },
				{ name: "UsageDisclosure", type: "string", json: "usageDisclosure" },
			]}
			ModuleConfigDiscovery: { name: "ModuleConfigDiscovery", fields: [
				{ name: "Public", type: _moduleDefs.DiscoveryResult, json: "public" },
				{ name: "Private", type: _moduleDefs.DiscoveryResult, json: "private" },
				{ name: "Protected", type: _moduleDefs.DiscoveryResult, json: "protected" },
			]}
			ModuleConfigRecordRevisions: { name: "ModuleConfigRecordRevisions", fields: [
				{ name: "Enabled", type: "bool", json: "enabled" },
				{ name: "Ident", type: "string", json: "ident" },
			]}
			ModuleConfigRecordDeDup: { name: "ModuleConfigRecordDeDup", fields: [
				{ name: "Strict", type: "bool", json: "-" },
				{ name: "Rules", goType: "DeDupRuleSet", json: "rules,omitempty" },
			]}
			DiscoveryResult: { name: "DiscoveryResult", fields: [
				{ name: "Result", slice: true, type: _moduleDefs.DiscoveryResultResult, json: "result" },
			]}
			DiscoveryResultResult: { name: "DiscoveryResultResult", fields: [
				{ name: "Lang", type: "string", json: "lang" },
				{ name: "Fields", slice: true, type: "string", json: "fields" },
			]}
			SystemFieldEncoding: { name: "SystemFieldEncoding", fields: [
				{ name: "ID", goType: "*EncodingStrategy", json: "id" },
				{ name: "ModuleID", goType: "*EncodingStrategy", json: "moduleID" },
				{ name: "NamespaceID", goType: "*EncodingStrategy", json: "namespaceID" },
				{ name: "Revision", goType: "*EncodingStrategy", json: "revision" },
				{ name: "Meta", goType: "*EncodingStrategy", json: "meta" },
				{ name: "OwnedBy", goType: "*EncodingStrategy", json: "ownedBy" },
				{ name: "CreatedAt", goType: "*EncodingStrategy", json: "createdAt" },
				{ name: "CreatedBy", goType: "*EncodingStrategy", json: "createdBy" },
				{ name: "UpdatedAt", goType: "*EncodingStrategy", json: "updatedAt" },
				{ name: "UpdatedBy", goType: "*EncodingStrategy", json: "updatedBy" },
				{ name: "DeletedAt", goType: "*EncodingStrategy", json: "deletedAt" },
				{ name: "DeletedBy", goType: "*EncodingStrategy", json: "deletedBy" },
			]}
		}

module: {
	handle: "module"
	features: {
		labelResourceType: "compose:module"
	}

	parents: [
		{handle: "namespace"},
	]

	model: {
		ident: "compose_module"
		attributes: {
			id: schema.IdField & {
				envoy: {
					yaml: {
						identKeyEncode: "moduleID"
					}
				}
			}
			tenant_id:  schema.TenantRefField
			project_id: schema.ProjectRefField
			namespace_id: {
				ident: "namespaceID",
				goType: "uint64",
				storeIdent: "rel_namespace"
				dal: { type: "Ref", refModelResType: "corteza::compose:namespace" }

				envoy: {
					yaml: {
						identKeyAlias: ["namespace", "namespace_id", "ns", "ns_id"]
					}
				}
			}
			handle: schema.HandleField
			name: {
				sortable: true
				dal: {}
			}
			meta: {
				goType: "rawJson"
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			config: {
				type: _moduleDefs.ModuleConfig
				dal: { type: "JSON", defaultEmptyObject: true }
				omitSetter: true
				omitGetter: true
			}
			fields: {
				goType: "types.ModuleFieldSet",
				store: false
				omitSetter: true
				omitGetter: true
				envoy: {
					yaml: {
						omitEncoder: true
					}
				}
			}
			created_at: schema.SortableTimestampNowField
			updated_at: schema.SortableTimestampNilField
			deleted_at: schema.SortableTimestampNilField
			created_by_agent: schema.AttributeAgentRef
		}

		indexes: {
			"primary": { attribute: "id" }
			"namespace": { attribute: "namespace_id" },
			"unique_handle": {
				fields: [{ attribute: "handle", modifiers: ["LOWERCASE"] }, { attribute: "namespace_id" }]
				predicate: "handle != '' AND deleted_at IS NULL"
			}
		}
	}

	types: {
		defs: _moduleDefs
	}

	filter: {
		struct: {
			module_id:    { goType: "[]uint64", ident: "moduleID", storeIdent: "id" }
			tenant_id:    schema.TenantFilterField
			project_id:   schema.ProjectFilterField
			namespace_id: { goType: "uint64", ident: "namespaceID", storeIdent: "rel_namespace" }
			handle: { goType: "string" }
			name: { goType: "string" }
			deleted: { goType: "filter.State", storeIdent: "deleted_at" }
		}

		query: ["handle", "name"]
		byValue: ["handle", "module_id", "namespace_id"]
		byNilState: ["deleted"]
	}

	envoy: {
		scoped: true
		yaml: {
			supportMappedInput: true
			mappedField: "Handle"
			identKeyAlias: ["modules", "mod"]

			extendedResourcePostProcess: true
			extendedResourceDecoders: [{
				ident: "source"
				expIdent: "Source"
				// @deprecated records is what the old version used
				identKeys: ["source", "datasource", "records"]
				supportMappedInput: false
			}]
		}
		store: {
			postSetEncoder: true
			extendedEncoder: true
			extendedFilterBuilder: true
			extendedDecoder: true
		}
	}

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"record.create": description:  "Create record"
			"owned-record.create": description:  "Create record with custom owner"
			"records.search": description: "List, search or filter records"
		}
	}

	refs: {
		items: [
			{path: "Config.DAL.ConnectionID", kind: "KindDalConnection", reason: "ReasonModuleConnection"},
		]
		extended: true
	}

	service: {
		scoped: true

		undelete: true

		customBodyOps: ["lookup", "search", "create", "update", "delete", "undelete"]

		customFunctions: [
			{
				name: "ReloadDALModels"
				cap:  "write"
				results: [{name: "err", goType: "error"}]
			},
		]

		customAccessOps: ["search", "create"]

		omitCreateProp: true
		updateProp:     "changed"
	}

	store: {
		ident: "composeModule"

		api: {
			lookups: [
				{
					fields: ["namespace_id", "handle"]
					constraintCheck: true
					nullConstraint: ["deleted_at"]
					description: """
						searches for compose module by handle (case-insensitive)
						"""
				}, {
					fields: ["namespace_id", "name"]
					nullConstraint: ["deleted_at"]
					description: """
						searches for compose module by name (case-insensitive)
						"""
				}, {
					fields: ["id"]
					description: """
						searches for compose module by ID

						It returns compose module even if deleted
						"""
				},
			]
		}
	}

	locale: {
		extended: true

		keys: {
			"name": {}
		}
	}

	//locale:
	//  resource:
	//    references: [ namespace, ID ]
	//
	//  extended: true
	//  keys:
	//    - name
}
