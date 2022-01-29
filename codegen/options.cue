package codegen

import (
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

options:
	[...schema.#codegen] &
	[
		//  for g in app.corteza.options {
		//   template: "gocode/options/$options_group.go.tpl"
		//   output:   "pkg/options/\(g.ident).gen.go"
		//   payload: {
		//    package: "options"
		//    imports: g.imports
		//    func:    g.expIdent
		//    struct:  g.expIdent + "Opt"
		//    options: g.options
		//   }
		//  },
		{
			template: "gocode/options/options.go.tpl"
			output:   "pkg/options/options.gen.go"
			payload: {
				package: "options"

				// make unique list of packages we'll import
				imports: [ for i in {for g in app.corteza.options for i in g.imports {"\(i)": i}} {i}]

				groups: [
				    for g in app.corteza.options {
				     func:    g.expIdent
				     struct:  g.expIdent + "Opt"
				     options: g.options
				    }

				]

			}
		},
	]
