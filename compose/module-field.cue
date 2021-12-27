package compose

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

moduleField: schema.#resource & {
	rbac: {
		resource: references: [ "namespaceID", "moduleID", "ID"]

		operations: {
			"recod.value.read": description:   "Read field value on records"
			"recod.value.update": description: "Update field value on records"
		}
	}

	//locale:
	//  resource:
	//    references: [ namespace, module, ID ]
	//
	//  skipSvc: true
	//  keys:
	//    - label
	//    - { name: descriptionView, path: meta.description.view, custom: true, customHandler: descriptionView }
	//    - { name: descriptionEdit, path: meta.description.edit, custom: true, customHandler: descriptionEdit }
	//    - { name: hintView,        path: meta.hint.view,        custom: true, customHandler: hintView }
	//    - { name: hintEdit,        path: meta.hint.edit,        custom: true, customHandler: hintEdit }
	//    - { name: validatorError, path: "expression.validator.{{validatorID}}.error", custom: true, customHandler: validatorError }
	//    - { name: optionsOptionTexts,
	//        path: "meta.options.{{value}}.text",
	//        custom: true,
	//        customHandler: optionsOptionTexts
	//        }
}
