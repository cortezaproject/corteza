package codegen

import (
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

[...schema.#codegen] &
[
	{
		template: "docs/options.adoc.tpl"
		output:   "src/modules/generated/partials/env-options.gen.adoc"
		payload: {
			groups: [
				for g in app.human.options {
					title: g.title
					intro?: g.intro

					options: g.options
				},
			]
		}
	},
]
