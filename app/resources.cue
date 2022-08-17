package app

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

_allFeaturesDisabled: {
	labels: false
	paging: false
	sorting: false
	checkFn: false
}

resources: {
	"rbac-rule": schema.#PkgResource & {
		package: {
			ident: "rbac"
			import: "github.com/cortezaproject/corteza-server/pkg/rbac"
		}

		ident: "rule"
		identPlural: "rules"
		expIdent: "Rule"

		features: _allFeaturesDisabled


		model: {
			ident: "rbac_rules"
			attributes: {
				role_id:   { primaryKey: true, goType: "uint64", ident: "roleID", storeIdent: "rel_role" }
				resource:  { primaryKey: true }
				operation: { primaryKey: true }
				access:    {                   goType: "types.Access" }
			}
		}

		store: {
			ident: "rbacRule"

			api: {
				functions: [
					{
						expIdent: "TransferRbacRules"
						args: [
							{ident: "src", goType: "uint64"},
							{ident: "dst", goType: "uint64"},
						]
						return: []
					},
				]
			}
		}
	}

	"label": schema.#PkgResource & {
		package: {
			ident: "labels"
			import: "github.com/cortezaproject/corteza-server/pkg/label/types"
		}

		ident: "label"
		identPlural: "labels"
		expIdent: "Label"

		features: _allFeaturesDisabled

		model: {
			attributes: {
				kind:        { primaryKey: true }
				resource_id: { primaryKey: true, goType: "uint64", ident: "resourceID", storeIdent: "rel_resource" }
				name:        { primaryKey: true, ignoreCase: true  }
				value:       {}
			}
		}

		filter: {
			expIdent: "LabelFilter"
			struct: {
				kind: {}
				rel_resource: { goType: "[]uint64", ident: "resourceID" }
				limit: { goType: "uint" }
			}

			byValue: ["kind" , "rel_resource", ]
		}

		store: {
			api: {
				lookups: [
					{
						fields: ["kind", "resource_id", "name"]
						description: """
							searches for label by kind, resource ID and name
							"""
					},
				]

				functions: [
					{
						expIdent: "DeleteExtraLabels"
						args: [
							{ident: "kind", goType: "string"},
							{ident: "resourceId", goType: "uint64"},
							{ident: "name", goType: "string", spread: true},
						]
						return: []
					},
				]
			}
		}
	}

	"flag": schema.#PkgResource & {
		package: {
			ident: "flag"
			import: "github.com/cortezaproject/corteza-server/pkg/flag/types"
		}

		ident: "flag"
		identPlural: "flags"
		expIdent: "Flag"

		features: _allFeaturesDisabled

		model: {
			attributes: {
				kind:        { primaryKey: true }
				resource_id: { primaryKey: true, goType: "uint64", ident: "resourceID", storeIdent: "rel_resource" }
				owned_by:    { primaryKey: true, goType: "uint64" }
		  	name:        { primaryKey: true, ignoreCase: true }
				active:      {                   goType: "bool"}
			}
		}

		filter: {
			expIdent: "FlagFilter"
			struct: {
				kind: {}
				resource_id: { goType: "[]uint64", ident: "resourceID", storeIdent: "rel_resource" }
				owned_by: { goType: "[]uint64", ident: "ownedBy" }
				name: { goType: "[]string", ident: "name" }
			}

			byValue: ["kind", "resource_id", "owned_by", "name", ]
		}

		store: {
			api: {
				lookups: [
					{
						fields: ["kind", "resource_id", "owned_by", "name"]
						description: """
							searches for flag by kind, resource ID, owner and name
							"""
					},
				]
			}
		}
	}

	"actionlog": schema.#PkgResource & {
		package: {
			ident: "actionlog"
			import: "github.com/cortezaproject/corteza-server/pkg/actionlog"
		}

		ident: "action"
		identPlural: "action"
		expIdent: "Action"

		features: {
			labels: false
			paging: false
			checkFn: false
		}

		model: {
			ident: "actionlog"
			attributes: {
  			id:           schema.IdField
				timestamp:    schema.SortableTimestampField & { storeIdent: "ts" }
				request_origin:  {}
				request_id:  { ident: "requestID" }
				actor_ip_addr:  { ident: "actorIPAddr"}
				actor_id:  { goType: "uint64", ident: "actorID" }
				resource:  {}
				action:  {}
				error:  {}
				severity:  { goType: "types.Severity" }
				description:  {}
				meta:  { goType: "types.Meta" }
			}
		}

		filter: {
			expIdent: "Filter"
			struct: {
				from_timestamp: { goType: "*time.Time" }
				to_timestamp: { goType: "*time.Time" }
				before_action_id: { goType: "uint64", ident: "beforeActionID" }
				actor_id: { goType: "[]uint64", ident: "actorID" }
				origin: {}
				resource: {}
				action: {}
				limit: { goType: "uint" }
			}

			byValue: ["action", "resource", "origin", "actor_id" ]
		}

		store: {
			ident: "actionlog"

			api: {
				lookups: [
					{
						fields: ["id"]
						description: """
							searches for action log by ID
							"""
					},
				]
			}
		}
	}

	"resource-activity": schema.#PkgResource & {
		package: {
			ident: "discovery"
			import: "github.com/cortezaproject/corteza-server/pkg/discovery/types"
		}

		ident: "resourceActivity"
		identPlural: "resourceActivities"
		expIdent: "ResourceActivity"
		expIdentPlural: "ResourceActivities"

		features: _allFeaturesDisabled

		model: {
			ident: "resource_activity_log"
			attributes: {
				id:           schema.IdField
				timestamp:    schema.SortableTimestampField & { storeIdent: "ts" }
				resource_type:   {}
				resource_action: {}
				resource_id:     { goType: "uint64", ident: "resourceID", storeIdent: "rel_resource" }
				meta:            { goType: "rawJson" }
			}
		}

		filter: {
			expIdent: "ResourceActivityFilter"
			struct: {
				from_timestamp: { goType: "*time.Time" }
				to_timestamp: { goType: "*time.Time" }
			}
		}

		store: {
			ident: "resourceActivity"

			api: {
				lookups: []
			}
		}
	}
}
