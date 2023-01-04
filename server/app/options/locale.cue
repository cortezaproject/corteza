package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

locale: schema.#optionsGroup & {
	handle: "locale"
	options: {
		languages: {
			defaultValue: "en"
			description: """
				List of compa delimited languages (language tags) to enable.
				In case when an enabled language can not be loaded, error is logged.

				When loading language configurations (config.xml) from the configured path(s).

				"""
		}

		path: {
			description: """
				One or more paths to locale config and translation files, separated by colon

				When with LOCALE_DEVELOPMENT_MODE=true, default value for path is ../../locale
			"""
		}

		query_string_param: {
			defaultValue: "lng"
			description: """
				Name of the query string parameter used to pass the language tag (it overrides Accept-Language header).
				Set it to empty string to disable detection from the query string.
				This parameter is ignored if only one language is enabled

				"""
		}

		resource_translations_enabled: {
			type:        "bool"
			description: """
        When enabled, an editor for resource translations is enabled in UI
        """
		}

		log: {
			type:        "bool"
			description: "Log locale related events and actions"
		}

		development_mode: {
			type: "bool"
			description: """
				When enabled, Corteza reloads language files on every request
				Enable this for debugging or developing.
				"""
		}
	}
}
