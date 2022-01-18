package app

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
	"github.com/cortezaproject/corteza-server/app/options"
	"github.com/cortezaproject/corteza-server/system"
	"github.com/cortezaproject/corteza-server/compose"
	"github.com/cortezaproject/corteza-server/automation"
	"github.com/cortezaproject/corteza-server/federation"
)

corteza: schema.#platform & {
	handle: "corteza"

	"options": [
		options.DB,
		options.HTTPClient,
		options.HTTPServer,
		options.RBAC,
		options.SCIM,
		options.SMTP,
		options.actionLog,
		options.apigw,
		options.auth,
		options.corredor,
		options.environment,
		options.eventbus,
		options.federation,
		options.limit,
		options.locale,
		options.log,
		options.messagebus,
		options.monitor,
		options.objectStore,
		options.plugins,
		options.provision,
		options.seeder,
		options.sentry,
		options.template,
		options.upgrade,
		options.waitFor,
		options.websocket,
		options.workflow,
	]

	components: [
		system.component,
		compose.component,
		automation.component,
		federation.component,
	]

	gig: {
		decoders: _decoders

		preprocessors: _noopPreprocessors +
			_attachmentPreprocessors +
			_envoyGenericPreprocessors +
			_importPreprocessors +
			_exportPreprocessors +
			[]

		postprocessors: _postprocessors

		workers: [
			{
				ident: "noop",
				description: "Noop worker has no predefined operations -- it proxies decoder results into postprocessor input.",
				tasks:       _noopPreprocessors
			},
			{
				ident: "attachment",
				description: "@todo not implemented.",
				tasks: _attachmentPreprocessors + _noopPreprocessors
			},
			{
				ident: "import",
				description: "Import worker is used to import external data into Corteza.",
				tasks:      _envoyGenericPreprocessors + _importPreprocessors + _noopPreprocessors
			},
			{
				ident: "export",
				description: "Export worker is used to export internal data into a predefined format.",
				tasks:      _envoyGenericPreprocessors + _exportPreprocessors + _noopPreprocessors
			},
		]
	}
}
