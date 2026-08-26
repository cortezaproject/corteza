package codegen

import (
	"strings"
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

[...schema.#codegen] &
[
	{
		template: "gocode/options/options.go.tpl"
		output:   "pkg/options/options.gen.go"
		payload: {
			package: "options"

			// make unique list of packages we'll import
			imports: [ for i in {for g in app.human.options for i in g.imports {"\(i)": i}} {i}]

			groups: [
				for g in app.human.options {
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
			// Value each type falls back to when an option declares no default
			_zeroValue: {
				"bool":          "false"
				"int":           "0"
				"int64":         "0"
				"float64":       "0"
				"time.Duration": "0s"
				"string":        ""
			}

			groups: [
				for g in app.human.options {
					title: "# " + strings.Join(strings.Split(g.title, "\n"), "\n# ")

					if (g.intro != _|_) {
						intro: "# " + strings.Join(strings.Split(g.intro, "\n"), "\n# ")
					}

					options: [
						for o in g.options {
							handle: o.handle
							env:    o.env
							type:   o.type

							// An option with a Go default has to say what that default is
							if (o.defaultGoExpr != _|_ && o.defaultValue == _|_ && o.defaultNote == _|_) {
								_documented: "\(o.env) sets defaultGoExpr, so it must also set defaultValue or defaultNote" & ""
							}

							// value on the commented assignment line
							defaultValue: string
							if (o.defaultValue != _|_) {
								defaultValue: o.defaultValue
							}
							if (o.defaultValue == _|_) {
								defaultValue: _zeroValue[o.type]
							}

							// what the "Default:" line says
							defaultDoc: string
							if (o.defaultNote != _|_) {
								defaultDoc: o.defaultNote
							}
							if (o.defaultNote == _|_) {
								defaultDoc: defaultValue
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
