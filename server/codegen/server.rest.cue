package codegen

import (
	"strings"
	"list"
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

// REST codegen loader.
//
// Ports the legacy pkg/codegen rest.go generator into the CUE codegen: it
// reads each component's `rest` block and emits the rest/handlers, rest/request
// and (opt-in) rest/<entrypoint>.gen.go controller files. All the per-param /
// per-api logic that lived as Go methods on the legacy structs (Parser,
// ExportedName, IsPlainValue, ResourceIDField, ScopeIDParams, PathIDArgs ...)
// is precomputed here so the templates stay dumb.

_goReservedIdents: [
	"break", "default", "func", "interface", "select", "case", "defer", "go",
	"map", "struct", "chan", "else", "goto", "package", "switch", "const",
	"fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var",
]

_restCrudNames: ["list", "create", "read", "update", "delete", "undelete"]

// _allParserTypes lists every type the Parser switch handles explicitly; the
// default branch fires for anything not in this set (and with no explicit parser).
_allParserTypes: [
	"[]uint64", "[]uint", "time.Time", "*time.Time", "sqlxTypes.JSONText",
	"int", "uint", "uint64", "int64", "float", "float64", "bool",
	"string", "[]string", "filter.State",
]

// _RestParser ports rest.go's (restEndpointParamDef).Parser(arg). Given a type
// and (possibly empty) explicit parser, returns the value-parsing expression
// the generated Fill() assigns: `r.X, err = <out>`.
_RestParser: {
	_type:   string
	_parser: string
	arg:     string

	out: [
		if _parser != "" {"\(_parser)(\(arg))"},
		if _parser == "" if _type == "[]uint64" {"payload.ParseUint64s(\(arg)), nil"},
		if _parser == "" if _type == "[]uint" {"payload.ParseUints(\(arg)), nil"},
		if _parser == "" if _type == "time.Time" {"payload.ParseISODateWithErr(\(arg))"},
		if _parser == "" if _type == "*time.Time" {"payload.ParseISODatePtrWithErr(\(arg))"},
		if _parser == "" if _type == "sqlxTypes.JSONText" {"payload.ParseJSONTextWithErr(\(arg))"},
		if _parser == "" if list.Contains(["int", "uint", "uint64", "int64", "float", "float64", "bool"], _type) {"payload.Parse\(strings.ToTitle(_type))(\(arg)), nil"},
		if _parser == "" if list.Contains(["string", "[]string"], _type) {"\(arg), nil"},
		if _parser == "" if _type == "filter.State" {"payload.ParseFilterState(\(arg)), nil"},
		if _parser == "" if !list.Contains(_allParserTypes, _type) {"\(_type)(\(arg)), nil"},
	][0]
}

// _RestParam resolves one raw param + origin into the fully-derived param the
// templates consume.
_RestParam: {
	_p:      schema.#restParam
	_origin: string
	scope:   bool | *false

	name:         _p.name
	exportedName: strings.ToTitle(_p.name)
	type:         _p.type
	title:        _p.title
	sensitive:    _p.sensitive
	origin:       _origin

	hasExplicitParser: _p.parser != ""
	isSlice:           strings.HasPrefix(_p.type, "[]") || strings.HasSuffix(_p.type, "Set")
	isUpload:          _p.type == "*multipart.FileHeader"
	fieldTag: [ if _p.type == "uint64" {"`json:\",string\"`"}, ""][0]

	parserVal:  (_RestParser & {_type: _p.type, _parser: _p.parser, arg: "val"}).out
	parserVal0: (_RestParser & {_type: _p.type, _parser: _p.parser, arg: "val[0]"}).out

	// IsPlainValue: can the create/update controller assign it straight into the
	// resource struct (or must a before-hook own it)?
	_t: strings.TrimPrefix(_p.type, "*")
	isPlain: [
		if isUpload {false},
		if _p.genHook {false},
		if strings.HasPrefix(_t, "map[") {true},
		if strings.Contains(_t, ".") {_t == "time.Time"},
		if list.Contains(["string", "[]string", "[]*string", "bool", "int", "int64", "uint", "uint64", "float", "float64"], _t) {true},
		false,
	][0]
}

// _RestApi resolves one api: merges endpoint-level params (scope), resolves all
// params, and derives the controller helpers (resource id, path id args, plain
// post params, isCRUD ...).
_RestApi: {
	_a:            schema.#restApi
	_ep:           schema.#restEndpoint
	entrypointExp: string
	endpointPath:  string

	out: {
		name:        _a.name
		nameExp:     strings.ToTitle(_a.name)
		reqIdent:    entrypointExp + strings.ToTitle(_a.name)
		method:      _a.method
		methodRoute: strings.ToTitle(strings.ToLower(_a.method))
		routePath:   endpointPath + _a.path
		raw:         _a.raw

		// Endpoint-level params are prepended (and, for path, scope every api).
		_path: [ for p in _ep.parameters.path {_RestParam & {_p: p, _origin: "PATH", scope: true}}] +
			[ for p in _a.parameters.path {_RestParam & {_p: p, _origin: "PATH"}}]
		_get: [ for p in _ep.parameters.get {_RestParam & {_p: p, _origin: "GET"}}] +
			[ for p in _a.parameters.get {_RestParam & {_p: p, _origin: "GET"}}]
		_post: [ for p in _ep.parameters.post {_RestParam & {_p: p, _origin: "POST"}}] +
			[ for p in _a.parameters.post {_RestParam & {_p: p, _origin: "POST"}}]

		params: {
			path: _path
			get:  _get
			post: _post
			all:  _path + _get + _post
		}

		// ResourceIDField: single path param, else the last one ending in "id".
		_idCandidates: [ for p in _path if strings.HasSuffix(strings.ToLower(p.name), "id") {p.exportedName}]
		resourceIDField: [
			if len(_path) == 0 {""},
			if len(_path) == 1 {_path[0].exportedName},
			if len(_path) > 1 if len(_idCandidates) > 0 {_idCandidates[len(_idCandidates)-1]},
			if len(_path) > 1 if len(_idCandidates) == 0 {_path[len(_path)-1].exportedName},
		][0]

		// id.ID request fields convert to uint64 with .Num().
		_idParam: [ for p in _path if p.exportedName == resourceIDField {p}]
		resourceIDArg: [
			if len(_idParam) == 0 {"r." + resourceIDField},
			if len(_idParam) > 0 if _idParam[0].type == "id.ID" {"r." + _idParam[0].exportedName + ".Num()"},
			if len(_idParam) > 0 if _idParam[0].type != "id.ID" {"r." + _idParam[0].exportedName},
		][0]

		// Leading scope ids (endpoint-level uint64 path params) for compound resources.
		scopeIDParams: [ for p in _path if p.scope if p.type == "uint64" {{exportedName: p.exportedName}}]
		_scopeArgs: [ for p in _path if p.scope if p.type == "uint64" {"r." + p.exportedName}]
		pathIDArgs: [
			if len(_scopeArgs) == 0 {resourceIDArg},
			if len(_scopeArgs) > 0 {strings.Join(_scopeArgs, ", ") + ", " + resourceIDArg},
		][0]

		hasUpdatedAt: len([ for p in _post if p.name == "updatedAt" {p}]) > 0
		plainPostParams: [ for p in _post if p.name != "updatedAt" if p.isPlain {{exportedName: p.exportedName}}]

		// No resource declares extra delete args yet; kept for parity.
		deleteExtraArgsCall: ""

		_isCrudName: list.Contains(_restCrudNames, _a.name)
		_needsID:    list.Contains(["read", "update", "delete", "undelete"], _a.name)
		isCRUD: [
			if _a.raw {false},
			if _a.genSkip {false},
			if !_isCrudName {false},
			if _needsID {resourceIDField != ""},
			true,
		][0]
	}
}

// _RestEndpoint resolves one endpoint into the template payload.
_RestEndpoint: {
	_e:  schema.#restEndpoint
	app: string

	entrypoint:    _e.entrypoint
	entrypointExp: strings.ToTitle(_e.entrypoint)
	path:          _e.path

	// normalizeImport: alias-aware quoting of import paths.
	importsNorm: [ for i in _e.imports {
		[
			if strings.Contains(i, " ") {"\(strings.SplitN(i, " ", 2)[0]) \"\(strings.Trim(strings.SplitN(i, " ", 2)[1], "\""))\""},
			if !strings.Contains(i, " ") {"\"\(i)\""},
		][0]
	}]

	// Controller struct name + service field.
	controllerName: [ if _e.resource != "" {strings.ToTitle(_e.resource)}, strings.ToTitle(_e.entrypoint)][0]
	_svcRaw:        strings.ToLower(entrypointExp[0:1]) + entrypointExp[1:]
	serviceField: [
		if _e.serviceField != "" {_e.serviceField},
		if list.Contains(_goReservedIdents, _svcRaw) {""},
		_svcRaw,
	][0]

	apis: [ for a in _e.apis {
		(_RestApi & {_a: a, _ep: _e, "entrypointExp": entrypointExp, endpointPath: _e.path}).out
	}]

	_anyCrud: len([ for a in apis if a.isCRUD {a}]) > 0
	hasCRUD:  _e.genController && serviceField != "" && _anyCrud

	hasCreate: len([ for a in apis if a.isCRUD if a.name == "create" {a}]) > 0
	hasUpdate: len([ for a in apis if a.isCRUD if a.name == "update" {a}]) > 0
	hasDelete: len([ for a in apis if a.isCRUD if a.name == "delete" || a.name == "undelete" {a}]) > 0

	genAfterCreate: _e.genAfterCreate
	genAfterUpdate: _e.genAfterUpdate
}

[...schema.#codegen] & (
	// rest/handlers/<entrypoint>.go -- one per endpoint
	[ for cmp in app.human.components if cmp.rest != _|_
		for e in cmp.rest.endpoints {
			template: "gocode/rest/$component_handler.go.tpl"
			output:   "\(cmp.ident)/rest/handlers/\(e.entrypoint).go"
			payload: {package: "handlers"} & (_RestEndpoint & {_e: e, app: cmp.ident})
		},
	] +
	// rest/request/<entrypoint>.go -- one per endpoint
	[ for cmp in app.human.components if cmp.rest != _|_
		for e in cmp.rest.endpoints {
			template: "gocode/rest/$component_request.go.tpl"
			output:   "\(cmp.ident)/rest/request/\(e.entrypoint).go"
			payload: {package: "request"} & (_RestEndpoint & {_e: e, app: cmp.ident})
		},
	] +
	// rest/<entrypoint>.gen.go controller -- only for genController endpoints
	[ for cmp in app.human.components if cmp.rest != _|_
		for e in cmp.rest.endpoints
		if (_RestEndpoint & {_e: e, app: cmp.ident}).hasCRUD {
			template: "gocode/rest/$component_controller.go.tpl"
			output:   "\(cmp.ident)/rest/\(e.entrypoint).gen.go"
			payload: {package: "rest"} & (_RestEndpoint & {_e: e, app: cmp.ident})
		},
	]
)
