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

resources: { [key=_]: {"handle": key, "component": "system", "platform": "corteza" } & schema.#PkgResource } & {
	"rbac-rule": {
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
				role_id:   {
					goType: "uint64",
					ident: "roleID",
					storeIdent: "rel_role"
					dal: { type: "Ref", refModelResType: "corteza::system:role" }
				}
				resource:  {
					dal: { length: 512 }
				}
				operation: {
					dal: { length: 50 }
				}
				access:    {
					goType: "types.Access"
					dal: { type: "Number" }
				}
			}

			indexes: {
				"primary": { attributes: [ "role_id", "resource", "operation" ] }
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

	"label": {
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
				kind: {
			 		dal: { length: 64 }
				}
				resource_id: {
			 		goType: "uint64",
			 		ident: "resourceID",
			 		storeIdent: "rel_resource"
			 		dal: { type: "ID" }
			 	}
				name: {
			 		ignoreCase: true
			 		dal: { length: 512 }
				}
				value: {
			 		dal: {}
				}
			}

			indexes: {
				"unique_kind_res_name": {
					fields: [
					  { attribute: "kind" },
						{ attribute: "resource_id" },
					 	{ attribute: "name", modifiers: [ "LOWERCASE" ] },
					]

				}
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

	"flag": {
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
				kind:        {
					dal: {}
				}
				resource_id: {
					goType: "uint64",
					ident: "resourceID",
					storeIdent: "rel_resource"
					dal: { type: "ID" }
				}
				owned_by:   schema.AttributeUserRef
		  	name:        {
		  		ignoreCase: true
					dal: {}
				}
				active: {
					goType: "bool"
					dal: { type: "Boolean" }
				}
			}

			indexes: {
				"unique_kind_res_owner_name": {
					 fields: [
						 { attribute: "kind" },
						 { attribute: "resource_id" },
						 { attribute: "owned_by" },
						 { attribute: "name", modifiers: [ "LOWERCASE" ] },
					 ]
				}
			}
		}

		filter: {
			expIdent: "FlagFilter"
			struct: {
				kind: {}
				resource_id: { goType: "[]uint64",
				ident: "resourceID", storeIdent: "rel_resource" }
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

	"actionlog": {
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
  			id:        schema.IdField
				timestamp: schema.SortableTimestampField & { storeIdent: "ts" }
				actor_ip_addr: {
					ident: "actorIPAddr"
					dal: { type: "Text", length: 64 }
				}
				actor_id: {
					goType: "uint64",
					ident: "actorID"
					dal: { type: "Ref", refModelResType: "corteza::system:user" }
				}
				request_origin: {
					dal: { type: "Text", length: 32 }
				}
				request_id:  {
					ident: "requestID"
					dal: { type: "Text", length: 256 }
				}
				resource:  {
					dal: { type: "Text", length: 512 }
				}
				action: {
					dal: { type: "Text", length: 64 }
				}
				error: {
					dal: {}
				}
				severity: {
					goType: "types.Severity"
					dal: { type: "Number", default: 0 }
				}
				description: {
					dal: {}
				}
				meta:  {
					goType: "types.Meta"
					dal: { type: "JSON", defaultEmptyObject: true }
				}
			}

			indexes: {
				"primary": { attribute: "id" }
				"action": { attribute: "action"}
				"actor_id": { attribute: "actor_id"}
				"rel_resource": { attribute: "resource"}
				"ts": { attribute: "timestamp"}
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

	"resource-activity": {
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
				id: schema.IdField
				timestamp: schema.SortableTimestampField & { storeIdent: "ts" }
				resource_type: {
					dal: {}
				}
				resource_action: {
					dal: {}
				}
				resource_id: {
					goType: "uint64",
					ident: "resourceID",
					storeIdent: "rel_resource"
					dal: { type: "Number" }
				}
				meta: {
					goType: "rawJson"
					dal: { type: "JSON", defaultEmptyObject: true }
				}
			}

			indexes: {
				"primary": { attribute: "id" }
				"resource": { attribute: "resource_id" }
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
