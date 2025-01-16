package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

messagebus: schema.#optionsGroup & {
	handle: "messagebus"
	options: {
		Enabled: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Enable messagebus"
		}
		log_enabled: {
			type:        "bool"
			description: "Enable extra logging for messagebus watchers"
		}
		servicebus_connection_string: {
			type:        "string"
			description: "Service Bus connection string for the namespace or for an entity, which contains a SharedAccessKeyName and SharedAccessKey properties"
		}
	}
	title: "Messaging queue"
}
