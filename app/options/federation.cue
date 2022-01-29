package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

federation: schema.#optionsGroup & {
	handle: "federation"

	imports: [
		"\"time\"",
	]

	options: {
		enabled: {
			type:        "bool"
			default:     "false"
			description: "Federation enabled on system, it toggles rest API endpoints, possibility to map modules in Compose and sync itself"
		}
		label: {
			type:        "string"
			default:     "\"federated\""
			description: "Federation label"
		}
		host: {
			type:        "string"
			default:     "\"local.cortezaproject.org\""
			description: "Host that is used during node pairing, also included in invitation"
		}
		structure_monitor_interval: {
			type:        "time.Duration"
			default:     "time.Minute * 2"
			description: "Delay in seconds for structure sync"
			env:         "FEDERATION_SYNC_STRUCTURE_MONITOR_INTERVAL"
		}
		structure_page_size: {
			type:        "int"
			default:     "1"
			description: "Bulk size in fetching for structure sync"
			env:         "FEDERATION_SYNC_STRUCTURE_PAGE_SIZE"
		}
		data_monitor_interval: {
			type:        "time.Duration"
			default:     "time.Second * 60"
			description: "Delay in seconds for data sync"
			env:         "FEDERATION_SYNC_DATA_MONITOR_INTERVAL"
		}
		data_page_size: {
			type:        "int"
			default:     "100"
			description: "Bulk size in fetching for data sync"
			env:         "FEDERATION_SYNC_DATA_PAGE_SIZE"
		}
	}
}
