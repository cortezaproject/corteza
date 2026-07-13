package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}
{{/*
  Service CRUD methods: each body is wrapped in a recordAction closure (action
  logging) and may emit eventbus Before/After events. The method bodies and the
  loadXxx / toLabeledXxx helpers are always generated. The access-control
  interface is generated when genAccessController is set. The struct is always
  generated here; the constructor lives in the hand-written companion file
  (alongside any custom methods). When extraServices is set the struct gains a
  services *{ident}Services field and the companion provides that struct +
  scopeServices(); otherwise {ident}Services{scope,caps} is auto-generated.
  Ops listed in customBodyOps delegate to hand-written on<Op> handlers.
*/}}
import (
	"context"

	"github.com/crusttech/human/server/pkg/actionlog"
{{- if and .load .hasID }}
	"github.com/crusttech/human/server/pkg/errors"
{{- end }}
{{- if and .handle .update }}
	"github.com/crusttech/human/server/pkg/handle"
{{- end }}
{{- if .usesLabel }}
	"github.com/crusttech/human/server/pkg/label"
{{- end }}
	"github.com/crusttech/human/server/pkg/scope"
{{- if .events }}
	"{{ .eventImport }}"
{{- end }}
	"github.com/crusttech/human/server/store"
	types "{{ .typesImport }}"
{{- range .customFunctionImports }}
	{{ . }}
{{- end }}
)
{{- if .genAccessController }}

type {{ .ident }}AccessController interface {
{{- if .create }}
	{{ .ac.create }}(context.Context) bool
{{- end }}
{{- if .search }}
	{{ .ac.search }}(context.Context) bool
{{- end }}
{{- if or .lookup (and .search .checkFn) }}
	{{ .ac.read }}(context.Context, *{{ .goType }}) bool
{{- end }}
{{- if .update }}
	{{ .ac.update }}(context.Context, *{{ .goType }}) bool
{{- end }}
{{- if or .delete .undelete }}
	{{ .ac.delete }}(context.Context, *{{ .goType }}) bool
{{- end }}
}
{{- end }}

type {{ .ident }} struct {
	actionlog actionlog.Recorder
	store     store.Storer
	ac        {{ .ident }}AccessController
{{- if .events }}
	eventbus  eventDispatcher
{{- end }}
{{- if .extraServices }}
	services  *{{ .ident }}Services
{{- end }}
}

{{- if not .extraServices }}

type {{ .ident }}Services struct {
	scope scope.Scope
	caps  scope.Capabilities
}
{{- end }}
{{- if .lookup }}

func (svc *{{ .ident }}) {{ .lookupIdent }}(ctx context.Context, {{ .lookupParentParamsTyped }}ID uint64) (res *{{ .goType }}, err error) {
	var (
		aProps = &{{ .ident }}ActionProps{ {{ .actionProp }}: &{{ .goType }}{ {{- if .hasID }}ID: ID{{- end }}} }
	)

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

{{- if has "lookup" .customBodyOps }}
		res, err = svc.onLookup(ctx, {{ .lookupParentParams }}ID, aProps)
		return err
{{- else }}
{{- if .hooks.beforeLookup }}
		if err = svc.beforeLookup(ctx, {{ .lookupParentParams }}ID); err != nil {
			return err
		}

{{- end }}
		if res, err = load{{ .expIdent }}(ctx, svc.store, ID); err != nil {
			return {{ .expIdent }}ErrInvalidID().Wrap(err)
		}
{{- if .hooks.afterLookup }}

		if res, err = svc.afterLookup(ctx, res); err != nil {
			return err
		}
{{- end }}

		aProps.set{{ .actionPropExp }}(res)
{{- if not .skipGuard }}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}
{{- end }}
{{- if not (has "lookup" .customAccessOps) }}

		if !svc.ac.{{ .ac.read }}(ctx, res) {
			return {{ .expIdent }}ErrNotAllowedToRead()
		}
{{- end }}

		return nil
{{- end }}
	}()

	return res, svc.recordAction(ctx, aProps, {{ .expIdent }}ActionLookup, err)
}
{{- end }}
{{- if .search }}

func (svc *{{ .ident }}) Search(ctx context.Context, filter {{ .goFilterType }}) (set {{ .goSetType }}, f {{ .goFilterType }}, err error) {
	var (
		aProps = &{{ .ident }}ActionProps{ {{ .filterProp }}: &filter }
	)
{{- if and .checkFn (not (has "search" .customBodyOps)) }}

	// For each fetched item, store backend will check if it is valid or not
	filter.Check = func(res *{{ .goType }}) (bool, error) {
		if !svc.ac.{{ .ac.read }}(ctx, res) {
			return false, nil
		}

		return true, nil
	}
{{- end }}

	err = func() error {
		if err = svc.checkScope(ctx, scope.CapRead); err != nil {
			return err
		}

{{- if has "search" .customBodyOps }}
{{- if not (has "search" .customAccessOps) }}
		if !svc.ac.{{ .ac.search }}(ctx) {
			return {{ .expIdent }}ErrNotAllowedToSearch()
		}
{{- end }}
		set, f, err = svc.onSearch(ctx, filter, aProps)
		return err
{{- else }}
{{- if .hooks.beforeSearch }}
		if err = svc.beforeSearch(ctx, &filter); err != nil {
			return err
		}

{{- end }}
		if !svc.ac.{{ .ac.search }}(ctx) {
			return {{ .expIdent }}ErrNotAllowedToSearch()
		}
{{- if .labels }}

		if len(filter.Labels) > 0 {
			filter.LabeledIDs, err = label.Search(
				ctx,
				svc.store,
				{{ .goType }}{}.LabelResourceKind(),
				filter.Labels,
			)
			if err != nil {
				return err
			}

			// labels specified but no labeled resources found
			if len(filter.LabeledIDs) == 0 {
				return nil
			}
		}
{{- end }}

		if set, f, err = store.Search{{ .storeExpIdentPlural }}(ctx, svc.store, filter); err != nil {
			return err
		}
{{- if .labels }}

		if err = label.Load(ctx, svc.store, toLabeled{{ .expIdentPlural }}(set)...); err != nil {
			return err
		}
{{- end }}
{{- if .hooks.afterSearch }}

		if err = svc.afterSearch(ctx, set, &filter); err != nil {
			return err
		}
{{- end }}

		return nil
{{- end }}
	}()

	return set, f, svc.recordAction(ctx, aProps, {{ .expIdent }}ActionSearch, err)
}
{{- end }}
{{- if .create }}

func (svc *{{ .ident }}) Create(ctx context.Context, new *{{ .goType }}) (res *{{ .goType }}, err error) {
	var (
		// set both the resource-named prop (action-log message templates
		// reference it by resource name) and the new prop
		aProps = &{{ .ident }}ActionProps{ {{ .actionProp }}: new{{ if .createProp }}, {{ .createProp }}: new{{ end }} }
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

{{- if has "create" .customBodyOps }}
{{- if not (has "create" .customAccessOps) }}
		if !svc.ac.{{ .ac.create }}(ctx) {
			return {{ .expIdent }}ErrNotAllowedToCreate()
		}
{{- end }}
		res = new
		return svc.onCreate(ctx, new)
{{- else }}
{{- if .hooks.validate }}
		if err = svc.validate(ctx, new); err != nil {
			return err
		}

{{- end }}
		if !svc.ac.{{ .ac.create }}(ctx) {
			return {{ .expIdent }}ErrNotAllowedToCreate()
		}
{{- if .hooks.beforeCreate }}

		if err = svc.beforeCreate(ctx, new); err != nil {
			return err
		}
{{- end }}

		new.ID = nextID()
{{- if .hasCreatedAt }}
		new.CreatedAt = *now()
{{- end }}

		if err = store.Create{{ .storeExpIdent }}(ctx, svc.store, new); err != nil {
			return
		}
{{- if .labels }}

		if err = label.Create(ctx, svc.store, new); err != nil {
			return
		}
{{- end }}

		res = new
{{- if .hooks.afterCreate }}

		if err = svc.afterCreate(ctx, res); err != nil {
			return err
		}
{{- end }}
		return nil
{{- end }}
	}()

	return res, svc.recordAction(ctx, aProps, {{ .expIdent }}ActionCreate, err)
}
{{- end }}
{{- if .update }}

func (svc *{{ .ident }}) Update(ctx context.Context, upd *{{ .goType }}) (res *{{ .goType }}, err error) {
	var (
		aProps = &{{ .ident }}ActionProps{ {{- .updateProp }}: upd}
		old    *{{ .goType }}
	)

{{- if has "update" .customBodyOps }}
{{- if not .hasID }}
	err = func() (err error) {
		res, err = svc.onUpdate(ctx, upd, aProps)
		return err
	}()
{{- else }}
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = load{{ .expIdent }}(ctx, s, upd.ID); err != nil {
			return
		}
{{- if .labels }}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}
{{- end }}

		aProps.set{{ .actionPropExp }}(res)
		aProps.set{{ .updateProp | title }}(res)
		old = res.Clone()
{{- if not .skipGuard }}

		if err = svc.guard(ctx, res); err != nil {
			return
		}
{{- end }}
{{- if .handle }}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return {{ .expIdent }}ErrInvalidHandle()
		}
{{- end }}
{{- if .stale }}

		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return {{ .expIdent }}ErrStaleData()
		}
{{- end }}
		before := func() error { return nil }
		after := func() error { return nil }

		if err = svc.onUpdate(ctx, s, upd, res, aProps, before, after); err != nil {
			return
		}
{{- range .settable }}
		res.{{ .expIdent }} = upd.{{ .expIdent }}
{{- end }}
{{- if .hasUpdatedAt }}
		res.UpdatedAt = now()
{{- end }}

		if err = store.Update{{ .storeExpIdent }}(ctx, s, res); err != nil {
			return err
		}
{{- if .labels }}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, s, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}
{{- end }}

		return nil
	})
{{- end }}
{{- else }}
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

{{- if .hooks.validate }}
		if err = svc.validate(ctx, upd); err != nil {
			return err
		}

{{- end }}
		if res, err = load{{ .expIdent }}(ctx, svc.store, upd.ID); err != nil {
			return
		}

		old = res.Clone()
		aProps.set{{ .actionPropExp }}(res)
{{- if not .skipGuard }}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}
{{- end }}
{{- if .handle }}

		if upd.Handle != res.Handle && !handle.IsValid(upd.Handle) {
			return {{ .expIdent }}ErrInvalidHandle()
		}
{{- end }}

		if !svc.ac.{{ .ac.update }}(ctx, res) {
			return {{ .expIdent }}ErrNotAllowedToUpdate()
		}
{{- if .stale }}

		// Test if stale (update has an older version of data)
		if isStale(upd.UpdatedAt, res.UpdatedAt, res.CreatedAt) {
			return {{ .expIdent }}ErrStaleData()
		}
{{- end }}
{{- if .hooks.beforeUpdate }}

		if err = svc.beforeUpdate(ctx, upd, res); err != nil {
			return err
		}
{{- end }}

{{- range .settable }}
		res.{{ .expIdent }} = upd.{{ .expIdent }}
{{- end }}
{{- if .hasUpdatedAt }}
		res.UpdatedAt = now()
{{- end }}

		if err = store.Update{{ .storeExpIdent }}(ctx, svc.store, res); err != nil {
			return err
		}
{{- if .labels }}

		if label.Changed(res.Labels, upd.Labels) {
			if err = label.Update(ctx, svc.store, upd); err != nil {
				return
			}
			res.Labels = upd.Labels
		}
{{- end }}
{{- if .hooks.afterUpdate }}

		if err = svc.afterUpdate(ctx, res); err != nil {
			return err
		}
{{- end }}
		return nil
	}()
{{- end }}

	return res, svc.recordAction(ctx, aProps, {{ .expIdent }}ActionUpdate, err, old, res)
}
{{- end }}
{{- if .delete }}

func (svc *{{ .ident }}) {{ .deleteIdent }}(ctx context.Context, {{ .deleteParentParamsTyped }}ID uint64{{ .deleteExtraArgsTyped }}) (err error) {
	var (
		aProps = &{{ .ident }}ActionProps{}
		res    *{{ .goType }}
	)

{{- if has "delete" .customBodyOps }}
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = load{{ .expIdent }}(ctx, s, ID); err != nil {
			return
		}
{{- if .labels }}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}
{{- end }}

		aProps.set{{ .actionPropExp }}(res)
{{- if not .skipGuard }}

		if err = svc.guard(ctx, res); err != nil {
			return
		}
{{- end }}

		return svc.onDelete(ctx, s, {{ .deleteParentParams }}res, {{ .deleteExtraArgsCall }}aProps)
	})
{{- else }}
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = load{{ .expIdent }}(ctx, svc.store, ID); err != nil {
			return
		}
{{- if not .skipGuard }}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}
{{- end }}

		aProps.set{{ .actionPropExp }}(res)

		if !svc.ac.{{ .ac.delete }}(ctx, res) {
			return {{ .expIdent }}ErrNotAllowedToDelete()
		}
{{- if .events }}

		if err = svc.eventbus.WaitFor(ctx, event.{{ .expIdent }}BeforeDelete(nil, res)); err != nil {
			return
		}
{{- end }}
{{- if .hooks.beforeDelete }}

		if err = svc.beforeDelete(ctx, res); err != nil {
			return err
		}
{{- end }}

		res.DeletedAt = now()
		if err = store.Update{{ .storeExpIdent }}(ctx, svc.store, res); err != nil {
			return
		}
{{- if .events }}

		_ = svc.eventbus.WaitFor(ctx, event.{{ .expIdent }}AfterDelete(nil, res))
{{- end }}
{{- if .hooks.afterDelete }}

		if err = svc.afterDelete(ctx, res); err != nil {
			return err
		}
{{- end }}
		return nil
	}()
{{- end }}

	return svc.recordAction(ctx, aProps, {{ .expIdent }}ActionDelete, err)
}
{{- end }}
{{- if .undelete }}

func (svc *{{ .ident }}) {{ .undeleteIdent }}(ctx context.Context, {{ .undeleteParentParamsTyped }}ID uint64) (err error) {
	var (
		aProps = &{{ .ident }}ActionProps{}
		res    *{{ .goType }}
	)

{{- if has "undelete" .customBodyOps }}
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
		if res, err = load{{ .expIdent }}(ctx, s, ID); err != nil {
			return
		}
{{- if .labels }}

		if err = label.Load(ctx, svc.store, res); err != nil {
			return err
		}
{{- end }}

		aProps.set{{ .actionPropExp }}(res)
{{- if not .skipGuard }}

		if err = svc.guard(ctx, res); err != nil {
			return
		}
{{- end }}

		return svc.onUndelete(ctx, s, {{ .undeleteParentParams }}res, aProps)
	})
{{- else }}
	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.CapWrite); err != nil {
			return err
		}

		if res, err = load{{ .expIdent }}(ctx, svc.store, ID); err != nil {
			return
		}
{{- if not .skipGuard }}

		if err = svc.guard(ctx, res); err != nil {
			return err
		}
{{- end }}

		aProps.set{{ .actionPropExp }}(res)

		if !svc.ac.{{ .ac.delete }}(ctx, res) {
			return {{ .expIdent }}ErrNotAllowedToUndelete()
		}
{{- if .hooks.beforeUndelete }}

		if err = svc.beforeUndelete(ctx, res); err != nil {
			return err
		}
{{- end }}

		res.DeletedAt = nil
		if err = store.Update{{ .storeExpIdent }}(ctx, svc.store, res); err != nil {
			return
		}
{{- if .hooks.afterUndelete }}

		if err = svc.afterUndelete(ctx, res); err != nil {
			return err
		}
{{- end }}

		return nil
	}()
{{- end }}

	return svc.recordAction(ctx, aProps, {{ .expIdent }}ActionUndelete, err)
}
{{- end }}
{{- if and .load .hasID }}

{{/* load looks up by globally-unique ID only -- scoped (compound-id) resources
     enforce their parent scope via the access check on the loaded record, so the
     parent ids are not threaded through here. */}}
func load{{ .expIdent }}(ctx context.Context, s store.{{ .storeExpIdentPlural }}, ID uint64) (res *{{ .goType }}, err error) {
	if ID == 0 {
		return nil, {{ .expIdent }}ErrInvalidID()
	}

	if res, err = store.Lookup{{ .storeExpIdent }}ByID(ctx, s, ID); errors.IsNotFound(err) {
		return nil, {{ .expIdent }}ErrNotFound()
	}

	return
}
{{- end }}
{{- if .genToLabeled }}

// toLabeled{{ .expIdentPlural }} converts to []label.LabeledResource
func toLabeled{{ .expIdentPlural }}(set []*{{ .goType }}) []label.LabeledResource {
	if len(set) == 0 {
		return nil
	}

	ll := make([]label.LabeledResource, len(set))
	for i := range set {
		ll[i] = set[i]
	}

	return ll
}
{{- end }}


{{- if and (not .guard) (not .skipGuard) }}
func (svc *{{ .ident }}) guard(_ context.Context, _ *{{ .goType }}) error { return nil }
{{- end }}

func (svc *{{ .ident }}) checkScope(ctx context.Context, cap scope.Capability) error {
	if err := scope.RequireTenantMembership(ctx); err != nil {
		return err
	}
	return scope.RequireCapability(ctx, cap)
}

{{- if not .extraServices }}
func (svc *{{ .ident }}) scopeServices(ctx context.Context) *{{ .ident }}Services {
	return &{{ .ident }}Services{
		scope: scope.GetScopeFromContext(ctx),
		caps:  scope.GetCapabilitiesFromContext(ctx),
	}
}
{{- end }}
{{- range .customFunctions }}

func (svc *{{ $.ident }}) {{ .name }}({{ .sigParams }}) {{ .resultsSig }} {
	var (
		aProps = &{{ $.ident }}ActionProps{}
	)

	err = func() (err error) {
		if err = svc.checkScope(ctx, scope.{{ .capConst }}); err != nil {
			return err
		}
{{- if .hasAc }}

		if !svc.ac.{{ .ac }}(ctx) {
			return {{ $.expIdent }}{{ .acErr }}()
		}
{{- end }}

		{{ .resultPrefix }}err = svc.{{ .handler }}(ctx, aProps{{ .callArgs }}{{ .variadicCallArg }})
		return err
	}()

	return {{ .resultPrefix }}svc.recordAction(ctx, aProps, {{ $.expIdent }}Action{{ .action }}, err)
}
{{- end }}
