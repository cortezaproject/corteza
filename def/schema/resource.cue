package schema

import (
	"strings"
)

#resource: {
	handle:    #baseHandle | *"unknown-resource"

	_words:   strings.Replace(strings.Replace(strings.Replace(handle, "-", " ", -1), "_", " ", -1), ".", " ", -1)

	ident:     #ident    | *strings.ToCamel(strings.Replace(strings.ToTitle(_words), " ", "", -1))
	expIdent:  #expIdent | *strings.Replace(strings.ToTitle(_words), " ", "", -1)
	platform:  #baseHandle | *"unknown-platform"
	component: string | *"unknown-component"

	// Fully qualified resource name
	fqrn: string | *(platform + "::" + component + ":" + handle)

	goType: string | *("types." + expIdent)

	// fields: #Fields
	// operations: #Operations

	// All known RBAC operations for this resource
	rbac: #rbacResource & {
		resource: {
			type:       fqrn
			"expIdent": expIdent
		}
	}

	// List of known keys for resource translation
	// locale?: {
	//  [Name=_]: {
	//   name:   Name & #Handle
	//   path:   string
	//   custom: bool | *false
	//  }
	// }
}

#fields: {
	// Each field can be
	[key=_]: #fields | *({name: key} & #field)
}

#field: {
	name:   #expIdent
	unique: bool | *false

	// Golang type (built-in or other)
	type: string | *"string"

	// System fields,
	system: bool | *false

	if name =~ "At$" {
		type: string | *"*time.Time"
	}
}

//#Operations: {
// [Operation=_]: {operation: Operation} & #Operation
//}

//#Operation: {
// name: #ExpIdent
// description: string
// can: string | false | *"\(name)"
//}

idField: {
	// Expecting ID field to allways have name ID
	name:   "ID"
	unique: true

	// Service fields,
	// @todo We might want to have a better name for this
	// service: true

	// @todo someday we'll replace this with the "ID" type
	type: "uint64"
}

handleField: {
	// Expecting ID field to allways have name ID
	name:   "handle"
	unique: true

	// @todo someday we'll replace this with the "ID" type
	type: "string" & #handle
}
