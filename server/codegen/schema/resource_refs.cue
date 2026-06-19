package schema

// Iteration strategy for a declared resource reference.
//
//   "direct"     — single uint64 field: Make(kind, r.<path>, reason, "<path>")
//   "sliceID"    — slice of uint64: for i, id := range r.<path> { Make(..., id, ...) }
//   "sliceField" — slice of structs: for i, item := range r.<path> { Make(..., item.<field>, ...) }
#ResourceRefIter: "direct" | "sliceID" | "sliceField"

// #ResourceRef declares a single cross-resource reference emitted by ResourceRefs().
#ResourceRef: {
	// dot-path to the holding field on the receiver (e.g. "WorkflowID",
	// "Execution.Model.LLMProviderID", "Behavior.KnowledgeBases")
	path: string

	// resourceref.Kind* constant WITHOUT the "resourceref." prefix,
	// e.g. "KindAutomationWorkflow"
	kind: string

	// resourceref.Reason* constant WITHOUT the "resourceref." prefix,
	// e.g. "ReasonTriggerWorkflow"
	reason: string

	// iteration strategy; defaults to single-field
	iter: #ResourceRefIter | *"direct"

	// for iter="sliceField": exported field on each struct element that holds
	// the uint64 ID (e.g. "ID"). Ignored for other iter values.
	field: string | *"ID"
}

// #ResourceRefs is the full resource-reference configuration for a resource.
#ResourceRefs: {
	// statically declarable refs, emitted in declaration order
	items: [...#ResourceRef] | *[]

	// when true the generated ResourceRefs() ends with
	//   return r.resourceRefsExt(out)
	// and the companion file MUST implement resourceRefsExt. Use when the
	// resource has additional refs that cannot be expressed as simple items
	// (nested loops, conditional logic, string parsing, etc.).
	extended: bool | *false
}
