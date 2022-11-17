package codegen

import (
	"github.com/cortezaproject/corteza/server/app"
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

[...schema.#codegen] &
[
	{
		template: "docs/options.adoc.tpl"
		output:   "src/modules/generated/partials/env-options.gen.adoc"
		payload: {
			groups: [
				for g in app.corteza.options {
					title: g.title
					intro?: g.intro

					options: g.options
				},
			]
		}
	},
]
