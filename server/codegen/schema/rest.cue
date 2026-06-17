package schema

// REST endpoint definitions -- a 1:1 CUE port of the legacy per-component
// rest.yaml (pkg/codegen). The new codegen (server.rest.cue) consumes this to
// generate the rest/handlers, rest/request and rest/<entrypoint>.gen.go files
// that pkg/codegen used to emit. Field names/semantics mirror the legacy
// restEndpointDef/restEndpointApi/restEndpointParamDef structs.
#rest: {
	endpoints: [...#restEndpoint]
}

#restEndpoint: {
	title:          string | *""
	path:           string
	entrypoint:     string
	authentication: [...string] | *[]
	imports: [...string] | *[]
	description: string | *""

	// GenController opts the endpoint into the generated CRUD controller
	// (<entrypoint>.gen.go). The hand-written companion provides the struct,
	// constructor, interfaces, makeFilter/makePayload and before-hooks.
	genController: bool | *false

	// Resource overrides the controller struct name + types.<Resource> when the
	// entrypoint differs from the singular resource type.
	resource: string | *""

	// serviceField overrides the controller struct field the generated methods
	// call into (ctrl.<field>.Search ...). Defaults to unexport(entrypoint).
	serviceField: string | *""

	genAfterCreate: bool | *false
	genAfterUpdate: bool | *false

	// Endpoint-level params scope every api beneath them (rare).
	parameters: #restParams

	apis: [...#restApi]
}

#restApi: {
	name:   string
	method: string
	title:  string | *""
	path:   string | *""

	// Raw skips the request/response envelope and delegates straight to the
	// controller method func(w, r).
	raw: bool | *false

	// genSkip opts a single api out of controller generation.
	genSkip: bool | *false

	parameters: #restParams
}

#restParams: {
	post: [...#restParam] | *[]
	path: [...#restParam] | *[]
	get: [...#restParam] | *[]
}

#restParam: {
	name:      string
	type:      string
	required:  bool | *false
	title:     string | *""
	sensitive: bool | *false

	// parser overrides the value parser used in the generated Fill().
	parser: string | *""

	// genHook routes create/update struct assignment through the before-hook
	// (the param is not a plain struct-assignable value).
	genHook: bool | *false
}
