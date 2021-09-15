package {{ .Package }}

{{ template "header-gentext.tpl" }}
{{ template "header-definitions.tpl" . }}

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	resourceTranslation struct {
		actionlog actionlog.Recorder
		locale    locale.Resource
		store     store.Storer
		ac        localeAccessController
	}

	localeAccessController interface {
		// CanManageResourceTranslation(context.Context) bool
	}

	ResourceTranslationService interface {
{{- range .Def }}
		{{ .Resource }}(ctx context.Context, {{ if .Locale.Resource }}{{ range .Locale.Resource.References }}{{ .Field }} uint64, {{ end }}{{ end }}) (locale.ResourceTranslationSet, error)
{{- end }}

		Upsert(context.Context, locale.ResourceTranslationSet) error
		Locale() locale.Resource
	}
)

func ResourceTranslation(ls locale.Resource) *resourceTranslation {
	return &resourceTranslation{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		locale:    ls,
	}
}

func (svc resourceTranslation) Upsert(ctx context.Context, rr locale.ResourceTranslationSet) (err error) {
	// @todo AC
	// @todo validation

	defer locale.Global().ReloadResourceTranslations(ctx)
	me := auth.GetIdentityFromContext(ctx)

	// - group by resource
	localeByRes := make(map[string]locale.ResourceTranslationSet)
	for _, r := range rr {
		localeByRes[r.Resource] = append(localeByRes[r.Resource], r)
	}

	// - for each resource, fetch the current state
	sysLocale := make(systemTypes.ResourceTranslationSet, 0, len(rr))
	for res, rr := range localeByRes {
		current, _, err := store.SearchResourceTranslations(ctx, svc.store, systemTypes.ResourceTranslationFilter{
			Resource: res,
		})
		if err != nil {
			return err
		}

		// get deltas and prepare upsert accordingly
		aux := current.New(rr)
		aux.Walk(func(cc *systemTypes.ResourceTranslation) error {
			cc.ID = nextID()
			cc.CreatedAt = *now()
			cc.CreatedBy = me.Identity()

			return nil
		})
		sysLocale = append(sysLocale, aux...)

		aux = current.Old(rr)
		aux.Walk(func(cc *systemTypes.ResourceTranslation) error {
			cc.UpdatedAt = now()
			cc.UpdatedBy = me.Identity()
			return nil
		})
		sysLocale = append(sysLocale, aux...)
	}

	err = store.UpsertResourceTranslation(ctx, svc.store, sysLocale...)
	if err != nil {
		return err
	}
	return nil
}

func (svc resourceTranslation) Locale() locale.Resource {
	return svc.locale
}

{{- range .Def }}
{{ $Resource := .Resource }}
{{ $GoType   := printf "types.%s" .Resource }}

func (svc resourceTranslation) {{ .Resource }}(ctx context.Context, {{ if .Locale.Resource }}{{ range .Locale.Resource.References }}{{ .Field }} uint64, {{ end }}{{ end }}) (locale.ResourceTranslationSet, error) {
	var (
		err       error
		out       locale.ResourceTranslationSet
		res       *{{ $GoType }}
		k         types.LocaleKey
	)

	res, err = svc.load{{$Resource}}(ctx, svc.store, {{ if .Locale.Resource }}{{ range .Locale.Resource.References }}{{ .Field }}, {{ end }}{{ end }})
	if err != nil {
		return nil, err
	}

	for _, tag := range svc.locale.Tags() {
{{- range .Locale.Keys}}
	{{- if not .Custom }}
		k = types.LocaleKey{{ $Resource }}{{coalesce (export .Name) (export .Path) }}
		out = append(out, &locale.ResourceTranslation{
			Resource: res.ResourceTranslation(),
			Lang:     tag.String(),
			Key:      k.Path,
			Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
		})
	{{ end }}
{{- end}}

{{- range .Locale.Keys}}
	{{- if and .Custom .CustomHandler }}
		auxResSet, err = svc.{{unexport $Resource}}{{export .CustomHandler}}Handler(ctx, tag, res, k.Path)
	{{- end}}
{{- end}}
	}

{{- if .Locale.Extended }}

	tmp, err := svc.{{unexport $Resource}}Extended(ctx, res)
	return append(out, tmp...), err
	{{- else }}
	return out, nil
{{- end }}
}

{{- end }}
