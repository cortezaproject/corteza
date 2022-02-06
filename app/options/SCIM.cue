package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

SCIM: schema.#optionsGroup & {
	handle: "scim"
	title:  "SCIM Server"

	// @todo remove explicitly defined expIdent and adjust the code
	expIdent: "SCIM"

	options: {
		enabled: {
			type:        "bool"
			description: "Enable SCIM subsystem"
		}
		base_URL: {
			defaultValue: "/scim"
			description:  "Prefix for SCIM API endpoints"
		}
		secret: {
			description: "Secret to use to validate requests on SCIM API endpoints"
		}
		external_id_as_primary: {
			type:        "bool"
			description: "Use external IDs in SCIM API endpoints"
		}
		external_id_validation: {
			defaultValue: "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
			description:  "Validates format of external IDs. Defaults to UUID"
		}
	}
}
