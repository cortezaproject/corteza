package app

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
	"github.com/cortezaproject/corteza-server/system"
	"github.com/cortezaproject/corteza-server/compose"
)

corteza: schema.#platform & {
	ident: "corteza"

	components: [
		system.component,
		compose.component,
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
