package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

settings: {
	ident: "settingValue"
	expIdent: "SettingValue"

	features: {
		labels: false
		paging: false
		sorting: false
		checkFn: false
	}

	model: {
		ident: "settings"
		attributes: {
			owned_by:    {
				goType: "uint64",
				storeIdent: "rel_owner"
				dal: { type: "Ref", refModelResType: "corteza::system:user" }
			}
			name:        {
				ignoreCase: true
				dal: { type: "Text", length: 512 }
			}
			value:       {
				goType: "rawJson"
				dal: { type: "JSON" }
			}
			updated_by:  schema.AttributeUserRef
			updated_at:  schema.SortableTimestampField
		}

		indexes: {
			"unique_kind_res_name": {
				fields: [
					{ attribute: "owned_by" },
				 	{ attribute: "name", modifiers: [ "LOWERCASE" ] },
				]
			}
		}
	}

	filter: {
		expIdent: "SettingsFilter"

		struct: {
			prefix: {}
			owned_by: { goType: "uint64", storeIdent: "rel_owner" }
		}

		byValue: [ "owned_by" ]
	}

	store: {
		api: {
			lookups: [
				{
					fields: ["name", "owned_by"]
					description: """
						searches for settings by name and owner
						"""
				}
			]
		}
	}
}
