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

	struct: {
		name:        { primaryKey: true, ignoreCase: true }
		owned_by:    { primaryKey: true, goType: "uint64", storeIdent: "rel_owner" }
		value:       { goType: "rawJson" }
		updated_by:  {                   goType: "uint64" }
		updated_at:  schema.SortableTimestampField
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
		settings: {
			rdbms: {
				table: "settings"
			}
		}

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
