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
			_envoyPreprocessors +
			[]

		postprocessors: _postprocessors

		workers: [
			{ident: "noop", tasks:       _noopPreprocessors},
			{ident: "attachment", tasks: _attachmentPreprocessors + _noopPreprocessors},
			{ident: "envoy", tasks:      _envoyPreprocessors + _noopPreprocessors},
		]
	}
}
