package codegen

import (
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
	"strings"
)

localeService:
	[...schema.#codegen] &
	[
		// wrapped with additional for loop to trim out templates with empty types list
		for tpl in [
			for cmp in app.corteza.components {
				template: "gocode/locale/$component_service.go.tpl"
				output:   "\(cmp.ident)/service/locale.gen.go"
				payload: {
					package: "service"
					imports: [
						"\"github.com/cortezaproject/corteza-server/\(cmp.ident)/types\"",
					]

					resources: [
						for res in cmp.resources if (res.locale != _|_) if (!res.locale.skipSvc) {
							expIdent: res.expIdent
							ident:    res.ident

							references: [ for field in res.locale.resource.references {strings.ToTitle(field)}]

							extended: res.locale.extended

							keys: [ for key in res.locale.keys if key.handlerFunc == _|_ {
								struct: key.struct

								"extended":    extended
								customHandler: key.customHandler
								if key.serviceFunc != _|_ {serviceFunc: key.serviceFunc}
							}]
						},
					]
				}
			},
			// skip empty type lists
		] if len(tpl.payload.resources) > 0 {tpl}]
