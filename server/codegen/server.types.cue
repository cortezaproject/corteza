package codegen

import (
	"strings"
	"list"
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

// Builds the per-resource payload for the type-struct template. Emits the
// resource struct fields (model.attributes + a Labels field when labelled);
// the JSON tag is derived by convention unless the attribute sets `json`.
_TypeResource: {
	res = "res": schema.#Resource

	// virtual sortableJSON attributes (e.g. a `name` mapped onto meta->>'name')
	// are sort keys extracted from a JSON column, NOT real struct fields -- exclude
	// them so codegen does not emit a phantom field.
	_attrs: [ for attr in res.model.attributes if attr.sortableJSON == _|_ {attr} ]

	// JSON-column nested types: every attribute stored as a JSON dal column whose
	// goType is an in-package exported identifier (strip leading "*" and the local
	// "types." qualifier) gets Scan/Value/Parse generated. Excludes builtins
	// (uint64), pointers to other packages (*expr.Vars -> "expr.Vars" has a dot),
	// slices/maps, and anything the resource lists in types.jsonTypesSkip (custom
	// Scan/Value stay hand-written).
	_jsonDerived: [
		for attr in res.model.attributes
		if attr.dal != _|_ if attr.dal.type == "JSON"
		let _b = strings.Replace(strings.TrimPrefix(attr.goType, "*"), "types.", "", -1)
		if _b =~ "^[A-Z][A-Za-z0-9_]*$"
		if !list.Contains(res.types.jsonTypesSkip, _b) {_b},
	]
	// derived + any manual extras, minus skips
	_jsonTypes: [ for t in (_jsonDerived + [ for t in res.types.jsonTypes if !list.Contains(_jsonDerived, t) {t} ]) {t} ]

	// types.defs split by kind. #StructType/#EnumType are closed definitions, so
	// unifying with {kind: string} succeeds only for enums (they declare kind);
	// structs violate closedness and collapse to _|_.
	// discriminate registry entries by shape: structs have fields, enums have
	// values. #SliceType/#MapType have neither and are excluded here (their decls
	// aren't generated yet — Sets still come from the legacy type_set.gen.go).
	_structDefs: [ for d in res.types.defs if (d & {fields: _}) != _|_ {d} ]
	_enumDefs:   [ for d in res.types.defs if (d & {values: _}) != _|_ {d} ]

	result: {
		expIdent: res.expIdent
		fileBase: strings.Replace(res.handle, "-", "_", -1)
		labels:   res.features.labels

		fields: [ for attr in _attrs {
			expIdent: attr.expIdent
			// strip the resource's own types-package prefix for in-package use
			goType: strings.Replace(attr.goType, "types.", "", -1)

			// dal is an optional field; dot into it only inside a guarded
			// comprehension so struct-only attributes (store:false, no dal block)
			// don't trip "cannot reference optional field".
			_dalType: [ if attr.dal != _|_ {attr.dal.type}, "" ][0]
			_isID:   _dalType == "ID" && attr.name == "id"
			_isRef:  !_isID && attr.goType == "uint64" && (_dalType == "ID" || _dalType == "Ref")
			_isTime: attr.goType == "time.Time" || attr.goType == "*time.Time"

			jsonTag: [
				// explicit overrides
				if attr.json != _|_ if (attr.json & string) != _|_ {"json:\"\(attr.json)\""},
				if attr.json != _|_ if (attr.json & bool) != _|_ if (attr.json & bool) {"json:\"\(attr.name)\""},
				if attr.json != _|_ if (attr.json & bool) == _|_ if (attr.json & string) == _|_ {
					// fill optional struct fields with defaults so they can be referenced
					let _j = {field: string | *attr.ident, omitEmpty: bool | *false, "string": bool | *false} & attr.json
					let _oe = [ if _j.omitEmpty {",omitempty"}, "" ][0]
					let _st = [ if _j.string {",string"}, "" ][0]
					"json:\"\(_j.field)\(_oe)\(_st)\""
				},
				// convention
				if attr.json == _|_ if _isID {"json:\"\(res.ident)ID,string\""},
				if attr.json == _|_ if _isRef {"json:\"\(attr.ident),string,omitempty\""},
				if attr.json == _|_ if _isTime {"json:\"\(attr.ident),omitempty\""},
				if attr.json == _|_ if !_isID && !_isRef && !_isTime {"json:\"\(attr.ident)\""},
			][0]
		}]

		// import gating
		needsTime:  len([ for attr in _attrs if attr.goType == "time.Time" || attr.goType == "*time.Time" {attr} ]) > 0
		needsLabel: res.features.labels
		// feature-injected struct fields (not model attributes), mirroring Labels
		flags:   res.features.flags
		imports: res.types.imports

		// id-as-string []uint64 helper types to generate (decl + Marshal/Unmarshal)
		idLists: res.types.idLists

		// JSON-column nested types to generate Scan/Value/Parse helpers for
		// (auto-derived from JSON dal attributes, minus jsonTypesSkip). ptr marks
		// types whose Parse returns *name.
		jsonTypes: [ for t in _jsonTypes {{name: t, ptr: list.Contains(res.types.jsonTypesPtr, t)}} ]

		// complex-field registry (types.defs): struct + enum decls authored in .cue,
		// generated into this component's types pkg. A struct stored in a JSON dal
		// column also gets Scan/Value/Parse via jsonTypes above (decl here, methods there).
		structDefs: [ for d in _structDefs {
			name: d.name
			doc:  [ if d.doc != _|_ {d.doc}, "" ][0]
			fields: [ for f in d.fields {
				// schema's _goType/_base are package-hidden (unreadable here), so the
				// in-package Go type is recomputed: the verbatim goType escape hatch
				// wins; otherwise a builtin token verbatim or the referenced def's
				// name, with ptr/slice prefixes.
				let _base = [ if f.type == _|_ {""}, if (f.type & string) != _|_ {f.type}, if (f.type & string) == _|_ {f.type.name} ][0]
				let _ptr = [ if f.ptr {"*"}, "" ][0]
				let _slice = [ if f.slice {"[]"}, "" ][0]
				name: f.name
				goType: [ if f.goType != _|_ {f.goType}, if f.goType == _|_ {"\(_ptr)\(_slice)\(_base)"} ][0]
				jsonTag: [
					if f.json != _|_ {"json:\"\(f.json)\""},
					if f.json == _|_ {"json:\"\(strings.ToCamel(f.name))\""},
				][0]
				doc: [ if f.doc != _|_ {f.doc}, "" ][0]
			}]
		}]

		enumDefs: [ for d in _enumDefs {
			name:   d.name
			goKind: d.kind
			doc:    [ if d.doc != _|_ {d.doc}, "" ][0]
			values: [ for v in d.values {
				// idents authored bare ("NormalUser") get the type name prepended;
				// idents authored already-qualified ("RequestKindCorrect") are kept
				// as-is. Both converge on the single-prefixed const the doc strings
				// assume, so mixed authoring across resources stays valid.
				constIdent: [ if strings.HasPrefix(v.ident, d.name) {v.ident}, d.name + v.ident ][0]
				constType:  d.name
				lit: [ if d.kind == "string" {"\"\(v.value)\""}, if d.kind == "int" {"\(v.value)"} ][0]
				doc: [ if v.doc != _|_ {v.doc}, "" ][0]
			}]
		}]
	}
}

[...schema.#codegen] &
[
	for cmp in app.human.components
	for res in cmp.resources
	if res.types != _|_ if res.types.gen {
		template: "gocode/types/$component_type.go.tpl"
		output:   "\(cmp.ident)/types/\(strings.Replace(res.handle, "-", "_", -1)).gen.go"
		payload: {
			package: "types"
			(_TypeResource & {"res": res}).result
		}
	},
	for cmp in app.human.components {
		template: "gocode/types/$component_resources.go.tpl"
		output:   "\(cmp.ident)/types/resources.gen.go"
		payload: {
			package: "types"

			cmpIdent: cmp.ident
			// Operation/resource validators, grouped by resource
			types: [
				for res in cmp.resources {
					const:   "\(res.expIdent)ResourceType"
					type:    res.fqrt
				},
				{
					const:     "ComponentResourceType"
					type:      cmp.fqrt
				},
			]
		}
	},
	for cmp in app.human.components {
		template: "gocode/types/$component_getters_setters.go.tpl"
		output:   "\(cmp.ident)/types/getters_setters.gen.go"
		payload: {
			package: "types"

			cmpIdent: cmp.ident
			resources: [ for res in cmp.resources { res }]
		}
	},
]
