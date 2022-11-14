package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

provision: schema.#optionsGroup & {
	handle: "provision"
	options: {
		always: {
			type:          "bool"
			defaultGoExpr: "true"
			description:   "Controls if provision should run when the server starts."
		}
		path: {
			defaultValue: "provision/*"
			description:  "Colon seperated paths to config files for provisioning."
		}
	}
	title: "Provisioning"
	intro: """
		Provisioning allows you to configure a {PRODUCT_NAME} instance when deployed.
		It occurs automatically after the {PRODUCT_NAME} server starts.

		[IMPORTANT]
		====
		We recommend you to keep provisioning enabled as it simplifies version updates by updating the database and updating settings.

		If you're doing local development or some debugging, you can disable this.
		====
		"""
}
