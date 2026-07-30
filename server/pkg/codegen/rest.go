package codegen

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type (
	// The following structure represents
	// legacy API definition (spec.json)
	//
	//

	// definitions are in one file
	restDef struct {
		App       string
		Source    string
		outputDir string

		Endpoints []*restEndpointDef
	}

	restEndpointDef struct {
		Title          string                `yaml:"title"`
		Path           string                `yaml:"path"`
		Entrypoint     string                `yaml:"entrypoint"`
		Authentication []interface{}         `yaml:"authentication,omitempty"`
		Apis           []*restEndpointApi    `yaml:"apis"`
		Imports        []string              `yaml:"imports"`
		Description    string                `yaml:"description,omitempty"`
		Params         restEndpointParamsDef `yaml:"parameters,omitempty"`

		// GenController opts the endpoint into the generated CRUD
		// controller (<entrypoint>.gen.go). The hand-written companion
		// (<entrypoint>.go) must provide the controller struct,
		// constructor, local service/access interfaces, makeFilter,
		// makePayload, makeFilterPayload and the beforeCreate/
		// beforeUpdate hooks. Off by default; enabled per-endpoint as
		// each companion is slimmed.
		GenController bool `yaml:"genController,omitempty"`

		// Resource overrides the controller struct name and the
		// types.<Resource> resource type used by the generated
		// controller. By default both are derived from the entrypoint
		// (export(entrypoint)). Set this when the entrypoint (and thus
		// the request-type prefix) differs from the singular resource
		// type, e.g. entrypoint "queues" but type "Queue".
		Resource string `yaml:"resource,omitempty"`

		// ServiceFieldName overrides the controller struct field the
		// generated controller calls into (ctrl.<field>.Search, ...).
		// Defaults to unexport(entrypoint). Set it when the companion's
		// service field is named differently, e.g. "svc" or "renderer".
		ServiceFieldName string `yaml:"serviceField,omitempty"`

		// Compound overrides automatic compound-resource detection for every
		// api under this endpoint (an api may still override it). See
		// restEndpointApi.Compound. Use `compound: false` when the route
		// nests the resource under a parent id (e.g. /projects/{projectID})
		// but the service is single-id. nil => auto-detect (default).
		Compound *bool `yaml:"compound,omitempty"`

		// GenAfterCreate makes the generated Create call an
		// ctrl.afterCreate(ctx, res, r) hook after svc.Create succeeds and
		// before makePayload. Use for post-service logic that needs the
		// created resource (e.g. its assigned id) -- such as a userGroup's
		// MemberAdd loop. Off by default; the hook is only emitted for
		// opted-in endpoints, so every other controller stays unchanged.
		// The companion must then provide
		// afterCreate(ctx, res *types.<Resource>, r *request.<Resource>Create) error.
		GenAfterCreate bool `yaml:"genAfterCreate,omitempty"`

		// GenAfterUpdate is the update counterpart of GenAfterCreate: the
		// generated Update calls ctrl.afterUpdate(ctx, res, r) after
		// svc.Update succeeds and before makePayload. The companion must
		// provide
		// afterUpdate(ctx, res *types.<Resource>, r *request.<Resource>Update) error.
		GenAfterUpdate bool `yaml:"genAfterUpdate,omitempty"`
	}

	restEndpointApi struct {
		Name   string                `yaml:"name"`
		Method string                `yaml:"method"`
		Title  string                `yaml:"title"`
		Path   string                `yaml:"path"`
		Params restEndpointParamsDef `yaml:"parameters,omitempty"`
		// Raw skips request/response envelope (api.Send / params.Fill) and
		// delegates straight to the controller method `func(w, r)`. Use for
		// SSE, file downloads, or any non-JSON transport.
		Raw bool `yaml:"raw,omitempty"`

		// GenSkip opts a single api out of controller generation even when
		// the endpoint has genController enabled. Use when a CRUD method
		// needs custom behavior the generator can't reproduce (e.g.
		// request-dependent post-load enrichment); the companion then
		// provides that method by hand.
		GenSkip bool `yaml:"genSkip,omitempty"`

		// Compound overrides the automatic compound-resource detection for
		// this api. The auto-detection treats every endpoint-level (scope)
		// uint64 path param as a leading scope id of the resource and so
		// emits e.g. FindByID(ctx, r.ProjectID, r.ProjectAiSystemID). Set
		// `compound: false` to force single-id ops -- only the resolved
		// resource id is passed (FindByID(ctx, r.ProjectAiSystemID)) and the
		// scope ids are ignored by the by-id methods and create/update.
		// Use when the route nests the resource under a parent (projectID)
		// but the service is single-id. nil => auto-detect (default).
		Compound *bool `yaml:"compound,omitempty"`

		// ScopeParams explicitly lists, by param name and in order, the
		// path params that form the leading scope ids passed to the by-id
		// service methods (before the resource id). It overrides both the
		// automatic endpoint-level scope detection and `compound`. Use for
		// resources whose by-id ops have non-uniform parent arity -- e.g.
		// compose page_layout, where read/FindByID takes only namespaceID
		// while delete/undelete take namespaceID+pageID. An empty list
		// forces single-id; nil falls back to auto-detection.
		ScopeParams []string `yaml:"scopeParams,omitempty"`

		// DeleteExtraArgs lists extra arguments forwarded to DeleteByID
		// AFTER the path ids (e.g. compose page's child-delete strategy:
		// DeleteByID(ctx, namespaceID, pageID, strategy)). Each entry names
		// a request field and, optionally, a Go type the field is cast to
		// (the request field type often differs, e.g. string request param
		// feeding a named-string service arg). Only meaningful on the
		// delete api. Empty => no change.
		DeleteExtraArgs []restDeleteExtraArg `yaml:"deleteExtraArgs,omitempty"`
	}

	// restDeleteExtraArg describes one extra argument forwarded to the
	// generated DeleteByID after the path ids. Field is the request struct
	// field name; Type, when set, is the Go type the field value is cast
	// to (e.g. types.PageChildrenDeleteStrategy).
	restDeleteExtraArg struct {
		Field string `yaml:"field"`
		Type  string `yaml:"type,omitempty"`
	}

	restEndpointParamsDef struct {
		Post []*restEndpointParamDef `yaml:"post"`
		Path []*restEndpointParamDef `yaml:"path"`
		Get  []*restEndpointParamDef `yaml:"get"`
	}

	restEndpointParamDef struct {
		Name     string `yaml:"name"`
		Type     string `yaml:"type"`
		Required bool   `yaml:"required"`
		Title    string `yaml:"title"`
		Origin   string

		// Scope marks a path param that is declared at the endpoint level
		// and therefore scopes every api under it (e.g. namespaceID for the
		// compose chart endpoint). The generated controller treats these as
		// the leading scope ids of compound resources: they are passed to
		// the FindByID/DeleteByID/... service methods and assigned onto the
		// resource struct (create/update), distinct from the per-api
		// resource id that maps to the struct's ID field. Set during
		// procRest, never from yaml.
		Scope bool `yaml:"-"`

		Sensitive bool `yaml:"sensitive"`

		DefinedParser string `yaml:"parser"`

		// GenHook forces the generated create/update controller to skip
		// the direct struct assignment for this param so the companion's
		// beforeCreate/beforeUpdate hook can fill it. Use when the
		// request type differs from the struct field type (e.g. string
		// request param feeding a named-string struct field that needs
		// an explicit conversion).
		GenHook bool `yaml:"genHook,omitempty"`
	}
)

func procRest(mm ...string) (dd []*restDef, err error) {
	dd = make([]*restDef, 0)

	for _, m := range mm {
		err = func() error {
			f, err := os.Open(m)
			if err != nil {
				return fmt.Errorf("%s read failed: %w", m, err)
			}

			defer f.Close()

			var d = &restDef{}

			if err := yaml.NewDecoder(f).Decode(d); err != nil {
				return err
			}

			d.outputDir = path.Dir(m)

			// Append params from endpoit to all apis
			for _, e := range d.Endpoints {
				// Endpoint-level path params scope every api beneath them;
				// mark them so the controller generator can treat them as
				// the leading scope ids of a compound resource.
				for _, p := range e.Params.Path {
					p.Scope = true
				}

				for _, a := range e.Apis {
					a.Params.Path = append(e.Params.Path, a.Params.Path...)
					a.Params.Post = append(e.Params.Post, a.Params.Post...)
					a.Params.Get = append(e.Params.Get, a.Params.Get...)

					// Inherit the endpoint-level compound override unless the
					// api sets its own; the api value always wins.
					if a.Compound == nil {
						a.Compound = e.Compound
					}
				}
			}

			dd = append(dd, d)
			return nil
		}()

		if err != nil {
			return nil, fmt.Errorf("failed to process %s: %w", m, err)
		}
	}

	return dd, nil
}

func genRest(tpl *template.Template, dd ...*restDef) (err error) {
	var (
		// Will only be generated if file does not exist previously
		tplHandler    = tpl.Lookup("rest_handler.go.tpl")
		tplRequest    = tpl.Lookup("rest_request.go.tpl")
		tplController = tpl.Lookup("rest_controller.go.tpl")

		dst string
	)

	for _, d := range dd {
		for _, e := range d.Endpoints {

			// Generic code, every event goes into one file (per app)
			dst = path.Join(d.outputDir, "rest", "handlers", e.Entrypoint+".go")
			err = goTemplate(dst, tplHandler, map[string]interface{}{
				"Source":   d.Source,
				"Endpoint": e,
				"App":      path.Base(d.outputDir),
			})
			if err != nil {
				return
			}

			// Generic code, every event goes into one file (per app)
			dst = path.Join(d.outputDir, "rest", "request", e.Entrypoint+".go")
			err = goTemplate(dst, tplRequest, map[string]interface{}{
				"Source":   d.Source,
				"Endpoint": e,
				"Imports":  e.Imports,
			})
			if err != nil {
				return
			}

			// Generated CRUD controller; opt-in per endpoint via
			// genController. The hand-written companion provides the
			// controller struct, constructor and hooks.
			if e.hasCRUD() {
				dst = path.Join(d.outputDir, "rest", e.Entrypoint+".gen.go")
				err = goTemplate(dst, tplController, map[string]interface{}{
					"Source":   d.Source,
					"Endpoint": e,
					"App":      path.Base(d.outputDir),
				})
				if err != nil {
					return
				}
			}
		}

	}

	return nil
}

// restCRUDNames are the api names the controller generator can implement.
var restCRUDNames = map[string]bool{
	"list":     true,
	"create":   true,
	"read":     true,
	"update":   true,
	"delete":   true,
	"undelete": true,
}

// goReservedIdents are Go keywords that cannot be used verbatim as the
// generated controller's service field identifier (unexport(entrypoint)).
var goReservedIdents = map[string]bool{
	"break": true, "default": true, "func": true, "interface": true, "select": true,
	"case": true, "defer": true, "go": true, "map": true, "struct": true,
	"chan": true, "else": true, "goto": true, "package": true, "switch": true,
	"const": true, "fallthrough": true, "if": true, "range": true, "type": true,
	"continue": true, "for": true, "import": true, "return": true, "var": true,
}

// ControllerName returns the exported name used for the generated
// controller struct receiver and the types.<name> resource type. It is
// the explicit Resource override when set, otherwise export(entrypoint).
func (e *restEndpointDef) ControllerName() string {
	if e.Resource != "" {
		return export(e.Resource)
	}

	return export(e.Entrypoint)
}

// ServiceField returns the controller's service field name (the
// unexported entrypoint). Empty when that would collide with a Go
// keyword, in which case the endpoint is not auto-generatable.
func (e *restEndpointDef) ServiceField() string {
	if e.ServiceFieldName != "" {
		return e.ServiceFieldName
	}

	f := unexport(e.Entrypoint)
	if goReservedIdents[f] {
		return ""
	}

	return f
}

// hasCRUD reports whether the endpoint exposes at least one CRUD api the
// controller generator can implement.
func (e *restEndpointDef) hasCRUD() bool {
	// Opt-in only: the companion file must be slimmed to the generated
	// contract before enabling this.
	if !e.GenController {
		return false
	}

	// A keyword service field can't be referenced (ctrl.<field>); skip
	// the whole endpoint so we never emit invalid Go.
	if e.ServiceField() == "" {
		return false
	}

	for _, a := range e.Apis {
		if a.IsCRUD() {
			return true
		}
	}

	return false
}

// IsCRUD reports whether this api is one of the standard CRUD endpoints
// that the controller generator can emit. Raw apis are never CRUD. apis
// whose semantics require a resource id in the path (read/update/delete/
// undelete) are only generatable when that id can be resolved.
func (a *restEndpointApi) IsCRUD() bool {
	if a.Raw {
		return false
	}

	// Explicitly opted out of generation; the companion provides this
	// method by hand.
	if a.GenSkip {
		return false
	}

	if !restCRUDNames[a.Name] {
		return false
	}

	switch a.Name {
	case "read", "update", "delete", "undelete":
		// These require a resolvable resource-id path param.
		return a.ResourceIDField() != ""
	}

	return true
}

// ResourceIDField returns the exported name of the path param that
// identifies the resource (e.g. ClientID, ApplicationID). When there is
// a single path param it is used; when there are several, the one whose
// name ends in "ID" (case-insensitive) is chosen. Returns "" when none
// can be determined.
func (a *restEndpointApi) ResourceIDField() string {
	pp := a.Params.Path
	if len(pp) == 0 {
		return ""
	}

	if len(pp) == 1 {
		return export(pp[0].Name)
	}

	// Multiple path params: the resource id is the last one ending in ID
	for i := len(pp) - 1; i >= 0; i-- {
		if strings.HasSuffix(strings.ToLower(pp[i].Name), "id") {
			return export(pp[i].Name)
		}
	}

	return export(pp[len(pp)-1].Name)
}

// ScopeIDParams returns the leading scope-id path params that precede the
// resource id in the FindByID/DeleteByID/... service-method args (and that
// create/update assign onto the resource struct), in declared order.
//
// Resolution order:
//  1. explicit ScopeParams override -- exactly the named path params (an
//     empty override list forces single-id);
//  2. `compound: false` override -- no scope ids (single-id);
//  3. auto-detection -- the endpoint-level (scope) uint64 path params.
//
// Empty for single-id (system) resources, which keeps their output
// byte-identical. The override knobs are nil/unset for every existing
// resource, so they are inert until opted into.
func (a *restEndpointApi) ScopeIDParams() []*restEndpointParamDef {
	// (1) explicit, ordered override by param name.
	if a.ScopeParams != nil {
		out := make([]*restEndpointParamDef, 0, len(a.ScopeParams))
		for _, name := range a.ScopeParams {
			for _, p := range a.Params.Path {
				if p.Name == name {
					out = append(out, p)
					break
				}
			}
		}
		return out
	}

	// (2) compound override forces single-id.
	if a.Compound != nil && !*a.Compound {
		return nil
	}

	// (3) auto-detection: endpoint-level uint64 scope params.
	out := make([]*restEndpointParamDef, 0, len(a.Params.Path))
	for _, p := range a.Params.Path {
		if p.Scope && p.Type == "uint64" {
			out = append(out, p)
		}
	}

	return out
}

// IsCompound reports whether the api addresses a compound, namespace-scoped
// resource: it carries at least one endpoint-level (scope) uint64 path param
// in addition to its own resource id (e.g. namespaceID + chartID).
func (a *restEndpointApi) IsCompound() bool {
	return len(a.ScopeIDParams()) > 0
}

// resourceIDParam returns the path-param def that identifies the resource,
// mirroring ResourceIDField's selection. Returns nil when none resolves.
func (a *restEndpointApi) resourceIDParam() *restEndpointParamDef {
	field := a.ResourceIDField()
	if field == "" {
		return nil
	}

	for _, p := range a.Params.Path {
		if p.ExportedName() == field {
			return p
		}
	}

	return nil
}

// idArg renders a single id path param as the request-field accessor the
// service method expects. When the request field is an id.ID (its parser
// produced an id.ID, e.g. system userGroup) but the service takes a plain
// uint64, the generated controller calls .Num() to convert.
func (p *restEndpointParamDef) idArg() string {
	if p.Type == "id.ID" {
		return "r." + p.ExportedName() + ".Num()"
	}

	return "r." + p.ExportedName()
}

// ResourceIDArg renders the resolved resource-id request field as the value
// the generated create/update assigns onto the resource struct's ID field.
// Like idArg, it converts id.ID request fields to uint64 with .Num().
func (a *restEndpointApi) ResourceIDArg() string {
	if p := a.resourceIDParam(); p != nil {
		return p.idArg()
	}

	return "r." + a.ResourceIDField()
}

// PathIDArgs returns the leading id arguments the service methods take,
// rendered as request-field accessors joined by ", " (e.g.
// "r.NamespaceID, r.ChartID"). For single-id resources this is just the
// resolved resource-id field; for compound resources the leading scope ids
// are prepended in declared order. id.ID-typed path params are converted
// with .Num() so they line up with uint64 service args.
func (a *restEndpointApi) PathIDArgs() string {
	scope := a.ScopeIDParams()
	idParam := a.resourceIDParam()

	idArg := "r." + a.ResourceIDField()
	if idParam != nil {
		idArg = idParam.idArg()
	}

	if len(scope) == 0 {
		return idArg
	}

	parts := make([]string, 0, len(scope)+1)
	for _, p := range scope {
		parts = append(parts, p.idArg())
	}
	parts = append(parts, idArg)

	return strings.Join(parts, ", ")
}

// DeleteExtraArgsCall renders the extra delete arguments the controller
// forwards to DeleteByID after the path ids, each prefixed with ", " so it
// can be appended directly after PathIDArgs (e.g.
// ", types.PageChildrenDeleteStrategy(r.Strategy)"). Empty when the api
// declares no extra delete args, which keeps single-id deletes unchanged.
func (a *restEndpointApi) DeleteExtraArgsCall() string {
	if len(a.DeleteExtraArgs) == 0 {
		return ""
	}

	parts := make([]string, 0, len(a.DeleteExtraArgs))
	for _, ea := range a.DeleteExtraArgs {
		acc := "r." + export(ea.Field)
		if ea.Type != "" {
			acc = ea.Type + "(" + acc + ")"
		}
		parts = append(parts, acc)
	}

	return ", " + strings.Join(parts, ", ")
}

// HasUpdatedAt reports whether the api has a POST param named updatedAt,
// which the update controller threads through for optimistic locking.
func (a *restEndpointApi) HasUpdatedAt() bool {
	for _, p := range a.Params.Post {
		if p.Name == "updatedAt" {
			return true
		}
	}

	return false
}

// PlainPostParams returns the POST params that can be assigned directly
// into the resource struct, excluding updatedAt (handled separately by
// update) and any complex/hook-filled params.
func (a *restEndpointApi) PlainPostParams() []*restEndpointParamDef {
	out := make([]*restEndpointParamDef, 0, len(a.Params.Post))
	for _, p := range a.Params.Post {
		if p.Name == "updatedAt" {
			continue
		}
		if !p.IsPlainValue() {
			continue
		}
		out = append(out, p)
	}

	return out
}

// IsPlainValue reports whether the param can be assigned straight into
// the resource struct by the generated create/update controller. Plain
// values are primitives, time pointers, string slices and maps. Domain
// types (types.*, sqlxTypes.*), uploads and other complex values are not
// plain and must be handled by a beforeCreate/beforeUpdate hook.
func (d *restEndpointParamDef) IsPlainValue() bool {
	if d.IsUpload() {
		return false
	}

	// Explicitly routed through a before-hook by the definition.
	if d.GenHook {
		return false
	}

	t := strings.TrimPrefix(d.Type, "*")

	// Maps are assigned directly (e.g. labels), even when the value type
	// is package-qualified.
	if strings.HasPrefix(t, "map[") {
		return true
	}

	// Package-qualified domain/JSON types need conversion or pointer
	// wrapping; the before-hook owns them.
	if strings.Contains(t, ".") {
		// time pointers/values are directly assignable.
		return t == "time.Time"
	}

	switch t {
	case "string", "[]string", "[]*string",
		"bool", "int", "int64", "uint", "uint64", "float", "float64":
		return true
	}

	return false
}

func (d *restEndpointParamsDef) All() []*restEndpointParamDef {
	var pp = make([]*restEndpointParamDef, 0)

	for _, p := range d.Path {
		p.Origin = "PATH"
		pp = append(pp, p)
	}

	for _, p := range d.Get {
		p.Origin = "GET"
		pp = append(pp, p)
	}

	for _, p := range d.Post {
		p.Origin = "POST"
		pp = append(pp, p)
	}

	return pp
}

// ExportedName returns the exported Go identifier for the param, used as
// both the request struct field name and (by the generated controller)
// the resource struct field name.
func (d *restEndpointParamDef) ExportedName() string {
	return export(d.Name)
}

func (d *restEndpointParamDef) IsUpload() bool {
	switch d.Type {
	case "*multipart.FileHeader":
		return true
	}

	return false
}

func (d *restEndpointParamDef) IsSlice() bool {
	return strings.HasPrefix(d.Type, "[]") || strings.HasSuffix(d.Type, "Set")
}

func (d *restEndpointParamDef) IsString() bool {
	switch d.Type {
	case "string", "[]string", "[]*string":
		return true
	}

	return false
}

func (d *restEndpointParamDef) FieldTag() string {
	switch d.Type {
	case "uint64":
		return "`json:\",string\"`"
	}

	return ""
}

func (d *restEndpointParamDef) HasExplicitParser() bool {
	return d.DefinedParser != ""
}

func (d *restEndpointParamDef) Parser(arg string) string {
	if d.HasExplicitParser() {
		return fmt.Sprintf("%s(%s)", d.DefinedParser, arg)
	}

	switch d.Type {
	case "[]uint64":
		return fmt.Sprintf("payload.ParseUint64s(%s), nil", arg)
	case "[]uint":
		return fmt.Sprintf("payload.ParseUints(%s), nil", arg)
	case "time.Time":
		return fmt.Sprintf("payload.ParseISODateWithErr(%s)", arg)
	case "*time.Time":
		return fmt.Sprintf("payload.ParseISODatePtrWithErr(%s)", arg)
	case "sqlxTypes.JSONText":
		return fmt.Sprintf("payload.ParseJSONTextWithErr(%s)", arg)
	case "int", "uint", "uint64", "int64", "float", "float64", "bool":
		return fmt.Sprintf("payload.Parse%s(%s), nil", export(d.Type), arg)
	case "string", "[]string":
		return fmt.Sprintf("%s, nil", arg)
	case "filter.State":
		return fmt.Sprintf("payload.ParseFilterState(%s), nil", arg)
	default:
		return fmt.Sprintf("%s(%s), nil", d.Type, arg)
	}

}
