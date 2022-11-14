package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"context"

{{- range .imports }}
    {{ . }}
{{- end }}
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	intAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"golang.org/x/text/language"
)

type (
	localeAccessControl interface {
		CanManageResourceTranslations(ctx context.Context) bool
	}

	resourceTranslationsManager struct {
		actionlog actionlog.Recorder
		locale    locale.Resource
		store     store.Storer
		ac        localeAccessController
	}

	localeAccessController interface {
		CanManageResourceTranslations(context.Context) bool
	}

	ResourceTranslationsManagerService interface {
{{- range .resources }}
		{{ .expIdent }}(ctx context.Context, {{ range .references }}{{ .param }} uint64, {{ end }}) (locale.ResourceTranslationSet, error)
{{- end }}

		Upsert(context.Context, locale.ResourceTranslationSet) error
		Locale() locale.Resource
	}
)

var ErrNotAllowedToManageResourceTranslations = errors.Unauthorized("not allowed to manage resource translations")

func ResourceTranslationsManager(ls locale.Resource) *resourceTranslationsManager {
	return &resourceTranslationsManager{
		actionlog: DefaultActionlog,
		store:     DefaultStore,
		ac:        DefaultAccessControl,
		locale:    ls,
	}
}

func (svc resourceTranslationsManager) Upsert(ctx context.Context, rr locale.ResourceTranslationSet) (err error) {
	// User is allowed to manage resource translations when:
	//  - managed resource translation strings are all for default language
	//  or
	//  - user is allowed to manage resource translations
	err = store.Tx(ctx, svc.store, func(ctx context.Context, s store.Storer) (err error) {
        if rr.ContainsForeign(svc.Locale().Default().Tag) {
            if !svc.ac.CanManageResourceTranslations(ctx) {
                return ErrNotAllowedToManageResourceTranslations
            }
        }

        for _, r := range rr {
            r.Msg = locale.SanitizeMessage(r.Msg)
        }

        // @todo validation

        me := intAuth.GetIdentityFromContext(ctx)

        // - group by resource
        localeByRes := make(map[string]locale.ResourceTranslationSet)
        for _, r := range rr {
            localeByRes[r.Resource] = append(localeByRes[r.Resource], r)
        }

        // - for each resource, fetch the current state
        sysLocale := make(systemTypes.ResourceTranslationSet, 0, len(rr))
        for res, rr := range localeByRes {
            current, _, err := store.SearchResourceTranslations(ctx, s, systemTypes.ResourceTranslationFilter{
                Resource: res,
                Deleted:  filter.StateInclusive,
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

            for _, diff := range current.Old(rr) {
                old := diff[0]
                new := diff[1]

                // soft delete; restore old message
                if new.Message == "" {
                    new.Message = old.Message
                    if new.DeletedAt == nil {
                        new.DeletedAt = now()
                        new.DeletedBy = me.Identity()
                    }
                } else {
                    new.UpdatedAt = now()
                    new.UpdatedBy = me.Identity()

                    new.DeletedAt = nil
                    new.DeletedBy = 0
                }

                sysLocale = append(sysLocale, new)
            }
        }

        err = store.UpsertResourceTranslation(ctx, s, sysLocale...)
        if err != nil {
            return err
        }

        return nil
    })
		if err != nil {
			return err
		}

		// Reload ALL resource translations
		// @todo we could probably do this more selectively and refresh only updated resources?
		return locale.Global().ReloadResourceTranslations(ctx)
}

func (svc resourceTranslationsManager) Locale() locale.Resource {
	return svc.locale
}

{{- range .resources }}

func (svc resourceTranslationsManager) {{ .expIdent }}(ctx context.Context, {{ range .references }}{{ .param }} uint64, {{ end }}) (locale.ResourceTranslationSet, error) {
	var (
		err       error
		out       locale.ResourceTranslationSet
		res       *types.{{ .expIdent }}
	)

	res, err = svc.load{{ .expIdent }}(ctx, svc.store, {{ range .references }}{{ .param }}, {{ end }})
	if err != nil {
		return nil, err
	}

    {{ $customHandlerExist := false }}
    {{range .keys}}
        {{if not .customHandler}}
            {{ $customHandlerExist = true }}
        {{end}}
    {{end}}

    {{ if $customHandlerExist }}
        var k types.LocaleKey
        for _, tag := range svc.locale.Tags() {
        {{- range .keys}}
            {{- if not .customHandler }}
                k = types.{{ .struct }}
                out = append(out, &locale.ResourceTranslation{
                    Resource: res.ResourceTranslation(),
                    Lang:     tag.String(),
                    Key:      k.Path,
                    Msg:      svc.locale.TResourceFor(tag, res.ResourceTranslation(), k.Path),
                })
            {{ end }}
        {{- end}}
        }
    {{ end }}

	{{ if .extended }}
		tmp, err := svc.{{ .ident }}Extended(ctx, res)
		return append(out, tmp...), err
		{{- else }}
		return out, nil
	{{- end }}
}

{{- end }}

func updateTranslations(ctx context.Context, ac localeAccessControl, lsvc ResourceTranslationsManagerService, tt ...*locale.ResourceTranslation) error {
	if lsvc == nil || lsvc.Locale() == nil || lsvc.Locale().Default() == nil {
		// gracefully handle partial initializations
		return nil
	}

	var (
		// assuming options will not change after start
		contentLang = lsvc.Locale().Default().Tag
	)

	if options.Locale().ResourceTranslationsEnabled {
		contentLang = locale.GetContentLanguageFromContext(ctx)
		// Resource translations enabled
		if contentLang == language.Und {
			// If no content-language meta (HTTP header) info was
			// used, do not run update translations - we do not know
			// what is the language that we're sending in
			return nil
		}

		if !lsvc.Locale().SupportedLang(contentLang) {
			// unsupported language
			return errors.InvalidData("unsupported language")
		}

		if !ac.CanManageResourceTranslations(ctx) {
			return errors.Unauthorized("not allowed to manage resource translations")
		}
	}

	locale.ResourceTranslationSet(tt).SetLanguage(contentLang)
	if err := lsvc.Upsert(ctx, tt); err != nil {
		return err
	}

	return nil
}

