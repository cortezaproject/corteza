package codegen

import (
	"strings"
	"list"
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

// Builds the per-resource payload for the CRUD service skeleton template.
// Mirrors the _StoreResource pattern in server.store.cue.
_ServiceResource: {
	res         = "res":         schema.#Resource
	typesImport = "typesImport": string

	// soft-delete / timestamp capabilities, derived from the model
	_hasDeletedAt: res.model.attributes["deleted_at"] != _|_
	_hasCreatedAt: res.model.attributes["created_at"] != _|_
	_hasUpdatedAt: res.model.attributes["updated_at"] != _|_

	// whether the model declares an `id` primary-key attribute. Composite-keyed
	// mapping resources (e.g. federation module_mapping) have none, so the
	// generated FindByID must not initialise the action-prop struct with {ID: ID}.
	_hasID: res.model.attributes["id"] != _|_

	// resource RBAC operations (empty when the resource declares none)
	_rbacOps: {}
	if res.rbac != _|_ {
		_rbacOps: res.rbac.operations
	}

	_delete:   res.service.delete && _hasDeletedAt
	_undelete: res.service.undelete && _hasDeletedAt

	// "standard" = generated AND not delegated to a custom on<Op> handler.
	// loadXxx / toLabeledXxx helpers and the store/errors/label imports are
	// only emitted when a standard-body op actually uses them.
	_cb:          res.service.customBodyOps
	_stdLookup:   res.service.lookup && !list.Contains(_cb, "lookup")
	_stdSearch:   res.service.search && !list.Contains(_cb, "search")
	_stdCreate:   res.service.create && !list.Contains(_cb, "create")
	_stdUpdate:   res.service.update && !list.Contains(_cb, "update")
	_stdDelete:   _delete && !list.Contains(_cb, "delete")
	_stdUndelete: _undelete && !list.Contains(_cb, "undelete")

	_load:        _stdLookup || _stdUpdate || _stdDelete || _stdUndelete
	_usesStore:   _load || _stdCreate || _stdUpdate || _stdDelete || _stdUndelete || _stdSearch
	_usesErrors:  _load
	// toLabeledXxx + the label import are emitted for any labelled resource --
	// companion (manual) Search methods rely on the helper too.
	_genToLabeled: res.features.labels
	_usesLabel:    res.features.labels

	// receiver prefix -- always pointer (unified)
	_recv: "*"

	// action-log resource prop name (+ its exported setter suffix)
	_actionProp: res.ident
	if res.service.actionProp != "" {_actionProp: res.service.actionProp}
	_actionPropExp: strings.ToTitle(_actionProp)

	// Update field-copy list: explicit when provided, otherwise derived.
	// The derived form is best-effort only (omitSetter is a DAL concern, not a
	// reliable service-settability signal) -- prefer setting updateFields.
	_derivedSettable: [
		for attr in res.model.attributes
		if attr.store
		if !attr.omitSetter
		if !list.Contains(res.service.omitUpdateFields, attr.name)
		if attr.dal != _|_
		if attr.dal.type != "ID"
		if attr.dal.type != "Timestamp" {
			{expIdent: attr.expIdent}
		},
	]

	result: {
		ident:          res.ident
		expIdent:       res.expIdent
		expIdentPlural: res.expIdentPlural
		fileBase:       strings.Replace(res.handle, "-", "_", -1)

		// generation mode
		recv: _recv
		events: res.service.events

		// public method names (unified across all generated services)
		lookupIdent:   "FindByID"
		deleteIdent:   "DeleteByID"
		undeleteIdent: "UndeleteByID"

		// action-log prop (field name + setter suffix)
		actionProp:    _actionProp
		actionPropExp: _actionPropExp
		filterProp:    res.service.filterProp
		updateProp:    res.service.updateProp

		// createProp: the secondary action-log prop set by Create alongside
		// actionProp. Empty when omitCreateProp -- the template skips it then.
		if res.service.omitCreateProp {
			createProp: ""
		}
		if !res.service.omitCreateProp {
			createProp: res.service.createProp
		}

		"typesImport": typesImport
		eventImport:   "\(strings.Replace(typesImport, "/types", "/service/event", 1))"
		goType:        "types.\(res.expIdent)"
		goSetType:     "types.\(res.expIdent)Set"
		goFilterType:  "types.\(res.filter.expIdent)"

		// store function name stems (store.Create<Stem>, store.Search<StemPlural>, ...)
		storeExpIdent:       res.store.expIdent
		storeExpIdentPlural: res.store.expIdentPlural

		// generated method toggles
		lookup:   res.service.lookup
		search:   res.service.search
		create:   res.service.create
		update:   res.service.update
		delete:   _delete
		undelete: _undelete

		// helper / import gating
		load:        _load
		usesStore:   _usesStore
		usesErrors:  _usesErrors
		usesLabel:   _usesLabel
		genToLabeled: _genToLabeled

		// capabilities used by Create / Update bodies
		hasCreatedAt: _hasCreatedAt
		hasUpdatedAt: _hasUpdatedAt
		stale:        _hasUpdatedAt && _hasCreatedAt

		// whether the model has an `id` attribute -- gates the {ID: ID} action-prop
		// initialiser in the generated FindByID body.
		hasID: _hasID

		// features
		labels:  res.features.labels
		flags:   res.features.flags
		checkFn: res.features.checkFn

		hooks: res.service.hooks
		guard: res.service.guard

		// generate the AC interface / struct+constructor in the .gen.go (opt-in)
		genAccessController: res.service.genAccessController
		genConstructor:      res.service.genConstructor

		// scoped (compound-id) support: namespace-scoped resources (compose) prepend
		// their model parent ids as leading args on the by-id methods and on<Op> hooks.
		scoped: res.service.scoped

		// parents derived from the resource model -- each entry carries the call/arg
		// param (e.g. "namespaceID") and the exported ref field (e.g. "NamespaceID").
		parents: [ for p in res.parents {{param: p.param, refField: p.refField}} ]

		// per-op parent overrides: by-id ops may use a subset of the declared
		// model parents (non-uniform arity, e.g. compose page_layout). The helper
		// below resolves, per op, the ordered parent list -- the per-op override
		// when present, otherwise the full scoped parents list.
		_opParentList: {
			lookup: [
				if res.service.scoped if res.service.opParents.lookup != _|_
				for h in res.service.opParents.lookup
				for p in res.parents if p.handle == h {p},
				if res.service.scoped if res.service.opParents.lookup == _|_
				for p in res.parents {p},
			]
			delete: [
				if res.service.scoped if res.service.opParents.delete != _|_
				for h in res.service.opParents.delete
				for p in res.parents if p.handle == h {p},
				if res.service.scoped if res.service.opParents.delete == _|_
				for p in res.parents {p},
			]
			undelete: [
				if res.service.scoped if res.service.opParents.undelete != _|_
				for h in res.service.opParents.undelete
				for p in res.parents if p.handle == h {p},
				if res.service.scoped if res.service.opParents.undelete == _|_
				for p in res.parents {p},
			]
		}

		// per-op typed signature prefix (e.g. "namespaceID uint64, ") and call-arg
		// prefix (e.g. "namespaceID, "). Empty when the op has no parents.
		_opParentParamsTyped: {
			for op, pl in _opParentList {
				(op): [ if len(pl) > 0 {strings.Join([ for p in pl {"\(p.param) uint64"} ], ", ") + ", "}, if len(pl) == 0 {""} ][0]
			}
		}
		_opParentParams: {
			for op, pl in _opParentList {
				(op): [ if len(pl) > 0 {strings.Join([ for p in pl {p.param} ], ", ") + ", "}, if len(pl) == 0 {""} ][0]
			}
		}

		// DeleteByID extra args: appended after the ID on the signature, forwarded to
		// onDelete after the ID. Empty for every existing resource.
		deleteExtraArgsTyped: strings.Join([ for a in res.service.deleteExtraArgs {", \(a.name) \(a.goType)"} ], "")
		deleteExtraArgsCall:  strings.Join([ for a in res.service.deleteExtraArgs {"\(a.name), "} ], "")

		// per-op parent prefixes consumed by the template -- default to the uniform
		// scoped parents so non-overridden resources stay byte-identical.
		lookupParentParamsTyped:   _opParentParamsTyped.lookup
		lookupParentParams:        _opParentParams.lookup
		deleteParentParamsTyped:   _opParentParamsTyped.delete
		deleteParentParams:        _opParentParams.delete
		undeleteParentParamsTyped: _opParentParamsTyped.undelete
		undeleteParentParams:      _opParentParams.undelete

		// per-op custom handler overrides
		customAccessOps: res.service.customAccessOps
		customBodyOps:   res.service.customBodyOps

		// access-control interface methods this service depends on.
		// Standard CRUD names are fixed; everything else (e.g. members.manage)
		// is pulled from the resource RBAC operations so the hand-written
		// companion methods that reference them still satisfy the interface.
		ac: {
			create: "CanCreate\(res.expIdent)"
			search: "CanSearch\(res.expIdentPlural)"
			read:   "CanRead\(res.expIdent)"
			update: "CanUpdate\(res.expIdent)"
			delete: "CanDelete\(res.expIdent)"

			extra: [
				for op in _rbacOps
				if !list.Contains(["read", "update", "delete", "create", "search"], op.handle) {
					{checkFuncName: op.checkFuncName}
				},
			]
		}

		// fields copied verbatim by Update -- explicit list wins, else derived.
		if len(res.service.updateFields) > 0 {
			settable: [ for f in res.service.updateFields {{expIdent: f}} ]
		}
		if len(res.service.updateFields) == 0 {
			settable: _derivedSettable
		}
	}
}

[...schema.#codegen] &
[
	for cmp in app.human.components
	for res in cmp.resources
	if res.service != _|_ if res.store != _|_ {
		template: "gocode/service/$component_service.go.tpl"
		output:   "\(cmp.ident)/service/\(strings.Replace(res.handle, "-", "_", -1)).gen.go"
		payload: {
			package: "service"
			(_ServiceResource & {
				"res":         res
				"typesImport": "github.com/crusttech/human/server/\(cmp.ident)/types"
			}).result
		}
	},
]
