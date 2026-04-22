package schema

import (
	"strings"
)

// fully qualified resource type
#FQRT: =~ "^corteza::(compose|system|federation|automation):[a-z][a-z0-9-]*$"

#Resource: {
	#_base

	// type: #resourceType | *""

	imports: [...{ import: string }]

	// copy field values from #_base
	handle: handle, ident: ident, expIdent: expIdent

	component: #baseHandle | *"component"
	platform:  #baseHandle | *"corteza"

	// Fully qualified resource name
	fqrt: #FQRT | *(platform + "::" + component + ":" + handle)

	model: #Model & {
		// use resource handle (plural) as model ident as default
		// model ident represents a db table or a container name
		ident: string | *"\(strings.Replace(handle, "-", "_", -1))s"
	}

	filter: {
		"expIdent": #expIdent | *"\(expIdent)Filter"

		struct: {
			[name=_]: {"name": name} & #ModelAttribute
		}

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

	envoy?: #resourceEnvoy & {
		omit: bool | *false
		$resourceIdent: ident
	}

	store?: {
		// how is this resource represented (prefixed/suffixed functions) in the store
		"ident": #ident | *ident
		"identPlural": #ident | *"\(store.ident)s"
		"expIdent": #expIdent | *strings.ToTitle(store.ident)
		"expIdentPlural": #expIdent | *"\(store.expIdent)s"

		api?: {
			lookups: [...{
				_expFields: [ for f in fields {strings.ToTitle(model.attributes[f].expIdent)}]

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

#resourceEnvoy: {
	$resourceIdent: string

	omit: bool

	// Scoped resources prioritize matching with resources in the same scope.
	// This is useful when we want to import multiple namespaces at the same time.
	scoped: bool | *false

	// YAML decode/encode configs
	yaml: {
		// supportMappedInput controls whether the resource can be presented in
		// a mapping node where the key of the map is an identifier.
		//
		// An example:
		//
		// modules:
		//   module1:
		//     name: module 1 name
		//
		// Where a resource with no mapping would look like:
		//
		// modules:
		//   - handle: module1
		//     name: module 1 name
		supportMappedInput: bool | *true
		// mappedField controls what identifier the map key represents
		// @todo this can probably be inferred so consider removing it.
		mappedField: string | *""

		identKeyLabel: string | *strings.ToLower($resourceIdent)
		identKeyAlias: [...string] | *[]
		// identKeys defines all of the identifiers that can be used when
		// referencing this resource
		identKeys: [...string] | *([identKeyLabel]+identKeyAlias)

		omitEncoder: bool | *false

		extendedResourcePostProcess: bool | *false
		extendedResourceDecoders: [...{
			ident: string
			expIdent: string
			identKeys: [...string]

			supportMappedInput: bool | *true
			mappedField: string | *""
		}] | *[]
		extendedResourceRefIdent: string | *""
		// enable or disable offloading unhandled nodes onto a custom default decoder
		defaultDecoder: bool | *false

		extendedResourceEncoders: [...{
			ident: string
			expIdent: string
			identKey: string | *""
		}] | *[]
	}

	// store decode/encode configs
	store: {
		// enable or disable custom logic after the resource is imported
		extendedEncoder: bool | *false
		extendedSubResources: bool | *false
		// temporary until I figure something better.
		// the idea here is that after all of the modules (or resources x) are imported
		// only then we should do something over them (like import records).
		postSetEncoder: bool | *false


		// enable or disable additional custom processing for determining 
		// resource references
		extendedRefDecoder: bool | *false
		handleField: string | *"Handle"

		customFilterBuilder: bool | *false
		// extendedFilterBuilder is called after the built-in which you can use
		// to append additional constraints to.
		//
		// Does nothing if customFilterBuilder is set to true
		extendedFilterBuilder: bool | *false

		// extendedDecoder is called after the built-in which you can use
		// to append additional nodes into.
		extendedDecoder: bool | *false

		// the resource calls a sanitizer function before saving to the database
		sanitizeBeforeSave: bool | *false
	}
}
