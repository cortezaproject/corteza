package schema

// #service generates CRUD service methods (FindByID/Search/Create/Update/
// DeleteByID/UndeleteByID + loadXxx). Opt in via a `service:` block.
//
// Each method body wraps a recordAction closure and may emit eventbus events;
// depends on generated *_actions.gen.go. Only bodies + helpers are generated --
// struct, access interface, constructor, custom methods live in the companion
// file. Non-standard CRUD: use customBodyOps + on<Op> handlers.
#service: {
	resourceExpIdent: #expIdent

	// Emit eventbus Before/After events in Create/Update/Delete. Needs an
	// events.yaml entry (event.<Resource>BeforeCreate, ...).
	events: bool | *false

	// cbEvents: when true, the customBodyOps Update scaffold generates real
	// Before/After event closures passed to onUpdate; false → noop closures.
	cbEvents: bool | *false

	// templateUpdate: when true, the customBodyOps Update method generates
	// field copy (from settable), UpdatedAt, store.Update and label.Update
	// AFTER calling onUpdate. onUpdate then owns only business logic (ACL,
	// validation). When false (default), onUpdate receives before/after
	// func() error callbacks and is responsible for the full update body.
	templateUpdate: bool | *false

	// Action-log resource prop name from <resource>_actions.yaml. Usually the
	// ident; some differ (configured_connection uses "connection"). Empty => ident.
	actionProp: #ident | *""

	// Action-log filter prop name. Usually "filter"; some use "search".
	filterProp: #ident | *"filter"

	// Action-log prop set (alongside actionProp) by the generated Create body.
	createProp: #ident | *"new"

	// Omit the createProp from the Create body's action props -- use when the
	// resource's action props have no dedicated "new" field and actionProp already
	// covers it (compose chart).
	omitCreateProp: bool | *false

	// Action-log prop set by the generated Update body. Usually "update"; some
	// resources name it differently (compose chart uses "changed").
	updateProp: #ident | *"update"

	// Standard CRUD toggles. Off => hand-write that op in the companion file.
	lookup:   bool | *true
	search:   bool | *true
	create:   bool | *true
	update:   bool | *true

	// Soft-delete ops; need a `deleted_at` attribute or the loader sets false.
	delete:   bool | *true
	undelete: bool | *false

	// Per-op overrides. Generated method always owns the action-log scaffold and,
	// unless in customAccessOps, the access check. customBodyOps delegates the
	// rest to svc.on<Op>(ctx, ...) -- for bespoke CRUD/validation/event order
	// (e.g. user). Ops: "lookup" "search" "create" "update" "delete" "undelete".
	customAccessOps: [...string] | *[]
	customBodyOps:   [...string] | *[]

	// customFunctions generates, per entry, a wrapper method that owns the standard
	// boilerplate -- tenant+project scope check, optional access-control guard and
	// optional action-log -- and delegates the actual body to a hand-written
	// svc.on<Name>(...) handler. Use for bespoke ops (Reorder, MarkAsRead,
	// Impersonate, ...) that still want the CRUD scaffold. The companion keeps only
	// the on<Name> handler; the public method itself is generated.
	customFunctions: [...{
		// method name; exported or unexported Go identifier. Delegates to svc.on<Name>.
		name: =~"^[a-zA-Z_][a-zA-Z0-9_]*$"

		// scope capability enforced by checkScope (tenant membership + capability).
		cap: *"write" | "read"

		// optional access-control guard: `if !svc.ac.<ac>(ctx) { return <Resource><acErr>() }`.
		// acErr is the error constructor (defaults to ErrNotAllowedTo<Name>).
		ac?:    #expIdent
		acErr?: #expIdent

		// action-log action stem -> <Resource>Action<action>. Defaults to the
		// title-cased function name. Override when the action name differs.
		action?: #expIdent

		// method signature after ctx. results' last entry must be the error return.
		// arg/result names are Go identifiers (may be exported-looking, e.g. ID).
		args?: [...{name: =~"^[a-zA-Z_][a-zA-Z0-9_]*$", goType: string}]
		results?: [...{name?: =~"^[a-zA-Z_][a-zA-Z0-9_]*$", goType: string}]
	}] | *[]

	// extra import specs the customFunctions signatures need beyond the service's
	// own types pkg (which is imported implicitly). Each entry is a full Go import
	// spec, optionally aliased, e.g. "expr \"github.com/.../pkg/expr\"".
	customFunctionImports: [...string] | *[]

	// Exported attr idents copied verbatim by Update (e.g. ["Handle","Type"]).
	// Set explicitly when non-trivial -- omitSetter is a DAL flag and unreliable
	// here. Empty => loader derives (stored, !omitSetter, non-ID/Timestamp, minus
	// omitUpdateFields).
	updateFields: [...#expIdent] | *[]

	// Subtractive override for the derived list (ignored when updateFields set).
	omitUpdateFields: [...#ident] | *[]

	// Optional hooks; when on, the generated method calls them and the companion
	// file must implement them.
	//
	//   validate      (ctx context.Context, res *T) error           // top of Create & Update, on the incoming resource
	//   beforeCreate  (ctx context.Context, new *T) error           // after access check, before id/timestamps
	//   afterCreate   (ctx context.Context, res *T) error           // after store create (+ label create)
	//   beforeUpdate  (ctx context.Context, upd, existing *T) error // after stale check, before field copy
	//   afterUpdate   (ctx context.Context, res *T) error           // after store update (+ label update)
	//   beforeDelete  (ctx context.Context, res *T) error           // after access check, before soft-delete
	//   afterDelete   (ctx context.Context, res *T) error           // after store delete -- post-store side effects (reload, cascade, sync)
	//   beforeUndelete(ctx context.Context, res *T) error           // after access check, before clearing deleted_at
	//   afterUndelete (ctx context.Context, res *T) error           // after store undelete -- post-store side effects
	//   beforeLookup  (ctx context.Context, ID uint64) error        // before the by-id load (guard/short-circuit)
	//   afterLookup   (ctx context.Context, res *T) (*T, error)     // post-process the loaded record (built-in/computed processor)
	//   beforeSearch  (ctx context.Context, f *TFilter) error       // mutate/validate the filter before access check + store search
	//   afterSearch   (ctx context.Context, set TSet, f *TFilter) error // post-process the result set
	hooks: {
		validate:       bool | *false
		beforeCreate:   bool | *false
		afterCreate:    bool | *false
		beforeUpdate:   bool | *false
		afterUpdate:    bool | *false
		beforeDelete:   bool | *false
		afterDelete:    bool | *false
		beforeUndelete: bool | *false
		afterUndelete:  bool | *false
		beforeLookup:   bool | *false
		afterLookup:    bool | *false
		beforeSearch:   bool | *false
		afterSearch:    bool | *false
	}

	// genAccessController generates the <resource>AccessController interface from the
	// standard CRUD access checks, gated by the enabled ops (CanCreate when create,
	// CanSearch when search, CanRead when lookup, CanUpdate when update, CanDelete
	// when delete||undelete). Only safe when the interface is exactly that derivable
	// set -- services whose access checks are cross-resource (e.g.
	// CanCreateChartOnNamespace), grant-based, component-scoped or undelete-specific
	// must keep the interface hand-written.
	genAccessController: bool | *false

	// extraServices suppresses generation of the {ident}Services struct and its
	// scopeServices() method. Set when the companion file hand-writes
	// {ident}Services to carry extra dependencies beyond scope/caps.
	// Default true: most services carry extra deps. Set false only for simple
	// resources whose Services struct needs only scope+caps.
	extraServices: bool | *true

	// guard, when set, suppresses the generated no-op guard() in the .gen.go file,
	// signalling that the companion file provides a real svc.guard(ctx, res) error
	// implementation. The guard call itself is always emitted on all standard and
	// customBodyOps paths (lookup, update, delete, undelete).
	guard: bool | *false

	// skipGuard, when true, omits the svc.guard(ctx, res) call from all generated
	// standard and customBodyOps paths for this resource. Use for resources where
	// guarding is never needed and the call overhead is unwanted.
	// Default false: guard is called on every op.
	skipGuard: bool | *false

	// scoped emits the resource's model parent ids (e.g. namespaceID) as leading
	// arguments on the by-id methods (FindByID/DeleteByID/UndeleteByID/loadXxx) and
	// their on<Op> hooks, for namespace-scoped compound-id resources (compose).
	scoped: bool | *false

	// deleteExtraArgs appends extra typed arguments AFTER the ID on the generated
	// DeleteByID signature (and forwards them to onDelete after the ID, before
	// aProps). Used by resources whose public delete carries an out-of-band option
	// -- e.g. compose page's child-delete strategy:
	//   DeleteByID(ctx, namespaceID, pageID, strategy types.PageChildrenDeleteStrategy)
	// Only meaningful when "delete" is in customBodyOps (the extra args are handled
	// by the hand-written onDelete). Empty => no change.
	deleteExtraArgs: [...{name: #ident, goType: string}] | *[]

	// opParents overrides, per by-id op, WHICH of the resource's declared model
	// parents are emitted as leading args. The value is the subset of parent
	// handles (in order) used by that op. Use for resources whose by-id methods
	// have non-uniform parent arity -- e.g. compose page_layout declares
	// [namespace, page] but FindByID takes only namespace while
	// DeleteByID/UndeleteByID take both. An unset op falls back to the full
	// parents list (scoped) or none (non-scoped). Keys: "lookup" "delete"
	// "undelete" (Search/Create/Update derive their parents from the resource).
	opParents: {
		lookup?: [...#baseHandle]
		delete?: [...#baseHandle]
		undelete?: [...#baseHandle]
	}
}
