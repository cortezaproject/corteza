package schema

#resource: #_base & {
	// copy field values from #_base
	handle: handle, ident: ident, expIdent: expIdent

	component: #baseHandle | *"component"
	platform:  #baseHandle | *"corteza"

	// Fully qualified resource name
	fqrn: string | *(platform + "::" + component + ":" + handle)

	// fields: #Fields
	// operations: #Operations

	// All parent resources
	parents: [... #_base & {
		// copy field values from #_base
		handle: handle, ident: ident, expIdent: expIdent

		refField: #expIdent | *(expIdent + "ID")
		param:    #ident | *(ident + "ID")
	}]

	// All known RBAC operations for this resource
	rbac: #rbacResource & {
		resourceExpIdent: expIdent
	}

	locale?: #locale & {
		resourceExpIdent: expIdent
		resource: {
			// @todo can we merge this with RBAC type (FQRN?)
			type: component + ":" + handle
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

//#fields: {
// // Each field can be
// [key=_]: #fields | *({name: key} & #field)
//}
//
//#field: {
// name:   #expIdent
// unique: bool | *false
//
// // Golang type (built-in or other)
// type: string | *"string"
//
// // System fields,
// system: bool | *false
//
// if name =~ "At$" {
//  type: string | *"*time.Time"
// }
//}

//#Operations: {
// [Operation=_]: {operation: Operation} & #Operation
//}

//#Operation: {
// name: #ExpIdent
// description: string
// can: string | false | *"\(name)"
//}

//idField: {
// // Expecting ID field to allways have name ID
// name:   "ID"
// unique: true
//
// // Service fields,
// // @todo We might want to have a better name for this
// // service: true
//
// // @todo someday we'll replace this with the "ID" type
// type: "uint64"
//}
//
//handleField: {
// // Expecting ID field to allways have name ID
// name:   "handle"
// unique: true
//
// // @todo someday we'll replace this with the "ID" type
// type: "string" & #handle
//}
