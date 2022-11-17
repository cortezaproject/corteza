package codegen

import (
	"strings"
	"github.com/cortezaproject/corteza/server/app"
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

[...schema.#codegen] &
[
	{
		template: "gocode/options/options.go.tpl"
		output:   "pkg/options/options.gen.go"
		payload: {
			package: "options"

			// make unique list of packages we'll import
			imports: [ for i in {for g in app.corteza.options for i in g.imports {"\(i)": i}} {i}]

			groups: [
				for g in app.corteza.options {
					func:   g.expIdent
					struct: g.expIdent + "Opt"
					options: [
						for o in g.options {
							o

							default?: string
							if (o.defaultGoExpr != _|_) {
								default: o.defaultGoExpr
							}

							if (o.defaultGoExpr == _|_ && o.defaultValue != _|_) {
								default: "\"" + o.defaultValue + "\""
							}
						},
					]
				},
			]
		}
	},
]+
[
	{
		template: "docs/.env.example.tpl"
		output:   ".env.example"
		syntax:   ".env"
		payload: {
			groups: [
				for g in app.corteza.options {
					title: "# " + strings.Join(strings.Split(g.title, "\n"), "\n# ")

					if (g.intro != _|_) {
						intro: "# " + strings.Join(strings.Split(g.intro, "\n"), "\n# ")
					}

					options: [
						for o in g.options {
							handle: o.handle
							env:    o.env
							type:   o.type

							defaultGoExpr?: string
							if (o.defaultGoExpr != _|_) {
								defaultGoExpr: o.defaultGoExpr
							}

							defaultValue?: string
							if (o.defaultValue != _|_) {
								defaultValue: o.defaultValue
							}

							if (o.description != _|_) {
								description: "# " + strings.Join(strings.Split(o.description, "\n"), "\n# ")
							}
						},
					]
				},
			]
		}
	},
]
