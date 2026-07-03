package schema

// #GoBuiltin is the allow-list of scalar Go types a #TypeField may use directly.
// CUE can't verify these resolve to real Go types, so they're enumerated here;
// composite field types must reference a #StructType/#EnumType instead.
#GoBuiltin: "string" | "bool" | "int" | "uint" | "uint64" | "int64" |
	"float64" | "[]byte" | "any" | "time.Time" | "*time.Time"

// #TypeField is one field of a #StructType. Set exactly one of:
//   type   — a builtin token or a value-reference to another registry entry.
//            Refs are existence-enforced at `cue export`; the render adds ptr/slice.
//   goType — verbatim escape hatch for shapes the registry can't express
//            (named slices/Sets, maps, external-pkg types). No enforcement, written
//            in full (e.g. "DeDupRuleSet", "map[string][]any", "language.Tag").
#TypeField: {
	name:   #expIdent
	type?:  #GoBuiltin | {name: #expIdent, ...}
	goType?: string
	ptr:   bool | *false
	slice: bool | *false

	// json is the full struct-tag content, verbatim ("avatarID,string"). When
	// absent the types loader derives it by convention from name/type.
	json?: string
	doc?:  string

	// in-package Go type the decl template renders (no `types.` qualifier — nested
	// types live in the same package as the struct that references them).
	if goType != _|_ {_goType: goType}
	if goType == _|_ {
		_base: string
		if (type & string) != _|_ {_base: type}
		if (type & string) == _|_ {_base: type.name}
		_goType: [ if ptr {"*"}, "" ][0] + [ if slice {"[]"}, "" ][0] + _base
	}
}

#StructType: {
	name: #expIdent
	fields: [...#TypeField]
	doc?: string
}

#EnumType: {
	name: #expIdent
	kind: "string" | "int" | *"string"
	values: [...{ident: #expIdent, value: string | int, doc?: string}]
	doc?: string
}

// #SliceType is a named slice type (`type <name> []<elem>` / `[]*<elem>`) — e.g. a
// corteza Set (`NgAutomationStepSet []*NgAutomationStep`). elem references a builtin
// or another registry entry; elemPtr adds the pointer.
#SliceType: {
	name: #expIdent
	// element type: a builtin/registry ref, OR elemGoType verbatim for an
	// external/unexpressible element (e.g. "wfexec.Frame"). Set exactly one.
	elem?:       #GoBuiltin | {name: #expIdent, ...}
	elemGoType?: string
	elemPtr: bool | *false
	doc?: string
}

// #MapType is a named map type (`type <name> map[<key>]<value>`) — e.g.
// `ModuleFieldOptions map[string]interface{}`. value/valueGoType mirror the
// slice element escape.
#MapType: {
	name: #expIdent
	key:  #GoBuiltin
	value?:       #GoBuiltin | {name: #expIdent, ...}
	valueGoType?: string
	valuePtr: bool | *false
	doc?: string
}
