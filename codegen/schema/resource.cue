package schema

import (
	"strings"
)

#Resource: {
	#_base

	imports: [...{ import: string }]

	// copy field values from #_base
	handle: handle, ident: ident, expIdent: expIdent

	component: #baseHandle | *"component"
	platform:  #baseHandle | *"corteza"

	// Fully qualified resource name
	fqrn: string | *(platform + "::" + component + ":" + handle)

	model: #Model
	filter: {
		"expIdent": #expIdent | *"\(expIdent)Filter"

		model: #Model

		// generate filtering by-nil-state for the specified fields
		"byNilState": [...string]

		// generate filtering by-false-state for the specified fields
		"byFalseState": [...string]

		// generate query filter for the specified fields
		"query": [...string]

		// filter resources by fields (eq)
		"byValue": [...string]
	}
	// operations: #Operations

	features: {
		// filtering by label
		labels:   bool | *true

		// filtering by flag
		flags:   bool | *false

		// support pagination
		paging:   bool | *true

		// support sorting
		sorting:   bool | *true

		// support resource check function
		checkFn:   bool | *true
	}

	// All parent resources
	parents: [... #_base & {
		// copy field values from #_base
		handle: handle, ident: ident, expIdent: expIdent

		refField: #expIdent | *(expIdent + "ID")
		param:    #ident | *(ident + "ID")
	}]

	// All known RBAC operations for this resource
	rbac?: #rbacResource & {
		resourceExpIdent: expIdent
	}

	locale?: #locale & {
		resourceExpIdent: expIdent
		resource: {
			// @todo can we merge this with RBAC type (FQRN?)
			type: component + ":" + handle
		}
	}

	store?: {
		// how is this resource represented (prefixed/suffixed functions) in the store
		"ident": #ident | *ident
		"identPlural": #ident | *"\(store.ident)s"
		"expIdent": #expIdent | *strings.ToTitle(store.ident)
		"expIdentPlural": #expIdent | *"\(store.expIdent)s"

		api?: {
			lookups: [...{
				_expFields: [ for f in fields {strings.ToTitle(model[f].expIdent)}]

				"expIdent":  "Lookup\(store.expIdent)By" + strings.Join(_expFields, "")
				description: string | *""

				// fields used for the lookup (must exist in the struct)
				fields: [...string]

				// Skip null constraints
				nullConstraint: [...string]
				constraintCheck: bool | *false
			}]

			functions: [...{
				expIdent: string

				description: string | *""

				args: [...{ident: #ident, goType: string, spread: bool | *false}]
				return: [...string]
			}]
		}

		settings: {
			defaultOrder: [...{ field: string, descending: bool | *false }]

			rdbms: {
				// use resource handle (plural) as RDBMS table name as default
				table: string | *"\(strings.Replace(handle, "-", "_", -1))s"
			}
		}
	}
}

#storeFunction: {
	expIdent: #expIdent
	args: [...string]
	return: [...string]
}

#PkgResource: #Resource & {
	package: {
		ident: #ident
		import: string
	}
}
