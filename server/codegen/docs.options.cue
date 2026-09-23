package codegen

import (
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

[...schema.#codegen] &
[
	{
		template: "docs/environment.md.tpl"
		output:   "reference/environment.gen.md"
		payload: {
			groups: [
				for g in app.human.options {
					title: g.title
					if g.intro != _|_ {
						intro: g.intro
					}

					options: g.options
				},
			]
		}
	},
]
