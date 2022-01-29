package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

websocket: schema.#optionsGroup & {
	handle: "websocket"

	imports: [
		"\"time\"",
	]
	title: "Websocket server"

	options: {
		log_enabled: {
			type:        "bool"
			description: "Enable extra logging for authentication flows"
		}
		timeout: {
			type:        "time.Duration"
			default:     "15 * time.Second"
			description: "Time before `WsServer` gets timed out."
		}
		ping_timeout: {
			type:    "time.Duration"
			default: "120 * time.Second"
		}
		ping_period: {
			type:    "time.Duration"
			default: "((120 * time.Second) * 9) / 10"
		}
	}
}
