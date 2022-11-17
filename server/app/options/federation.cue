package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

federation: schema.#optionsGroup & {
	handle: "federation"

	imports: [
		"\"time\"",
	]

	options: {
		enabled: {
			type:        "bool"
			description: "Federation enabled on system, it toggles rest API endpoints, possibility to map modules in Compose and sync itself"
		}
		label: {
			type:         "string"
			defaultValue: "federated"
			description:  "Federation label"
		}
		host: {
			type:         "string"
			defaultValue: "local.cortezaproject.org"
			description:  "Host that is used during node pairing, also included in invitation"
		}
		structure_monitor_interval: {
			type:          "time.Duration"
			defaultGoExpr: "time.Minute * 2"
			defaultValue:  "2m"
			description:   "Delay in seconds for structure sync"
			env:           "FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL"
		}
		structure_page_size: {
			type:          "int"
			defaultGoExpr: "1"
			defaultValue:  "1"
			description:   "Bulk size in fetching for structure sync"
			env:           "FEDERATION_SYNC_STRUCTURE_PAGE_SIZE"
		}
		data_monitor_interval: {
			type:          "time.Duration"
			defaultGoExpr: "time.Minute "
			defaultValue:  "1m"
			description:   "Delay in seconds for data sync"
			env:           "FEDERATION_SYNC_DATA_MONITOR_INTERVAL"
		}
		data_page_size: {
			type:          "int"
			defaultGoExpr: "100"
			defaultValue:  "100"
			description:   "Bulk size in fetching for data sync"
			env:           "FEDERATION_SYNC_DATA_PAGE_SIZE"
		}
	}
}
