package locale

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	locale "github.com/cortezaproject/corteza-locale"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type (
	service struct {
		l sync.RWMutex

		// logger facility for the locale service
		log *zap.Logger

		// options
		opt options.LocaleOpt

		// configured sources where we can load languages from
		src []string

		// all configured tags
		tags []language.Tag

		// set of all known languages
		set map[language.Tag]*Language

		// default language
		def *Language
	}
)

var (
	// Global locales service
	global *service
)

// Global returns global RBAC service
func Global() *service {
	return global
}

// SetGlobal re-sets global service
func SetGlobal(ll *service) {
	global = ll
}

func Static(ll ...*Language) *service {
	return &service{}
}

func Service(log *zap.Logger, opt options.LocaleOpt) (*service, error) {
	svc := &service{
		opt: opt,
		src: strings.Split(opt.Path, ":"),
		log: log.Named("locale"),
	}

	for _, lang := range strings.Split(opt.Languages, ",") {
		lang = strings.TrimSpace(lang)
		tag := language.Make(lang)
		if tag.IsRoot() {
			return nil, fmt.Errorf("failed to parse language '%s' (did you use comma to separate them?)", lang)
		}

		svc.tags = append(svc.tags, tag)
	}

	return svc, svc.Reload()
}

// Default language
func (svc *service) Default() *Language {
	return svc.def
}

// Tags of all loaded languages
func (svc *service) Tags() (tt []language.Tag) {
	tt = make([]language.Tag, 0, len(svc.set))
	for t := range svc.set {
		tt = append(tt, t)
	}

	return
}

// Reload all embedded (via github.com/cortezaproject/corteza-locale package)
// and imported (from one or more paths found in LOCALE_PATH) translations
func (svc *service) Reload() (err error) {
	var (
		i       int
		lang    *Language
		ll, aux []*Language

		logFields = make([]zap.Field, 0)
	)

	svc.l.RLock()
	defer svc.l.RUnlock()

	svc.log.Info("reloading",
		zap.Strings("path", svc.src),
		zap.Strings("tags", tagsToStrings(svc.tags)),
	)

	svc.set = make(map[language.Tag]*Language)

	// load embedded locales from the corteza-locale package
	if ll, err = loadConfigs(locale.Languages()); err != nil {
		return fmt.Errorf("could not load embedded locales: %w", err)
	} else {
		// cleanup src and effectively mark the
		// loaded languages as embedded
		for _, l := range ll {
			l.src = ""
		}
	}

	// load imported locales from the configured (LOCALE_PATH) paths
	for _, p := range svc.src {
		aux, err = loadConfigs(os.DirFS(p))
		if err != nil {
			return fmt.Errorf("could not load imported locales: %w", err)
		}

		for _, l := range aux {
			l.src = p + "/" + l.src
		}

		ll = append(ll, aux...)
	}

	for i, lang = range ll {
		logFields = []zap.Field{
			zap.Stringer("tag", lang.Tag),
		}

		if !hasTag(lang.Tag, svc.tags) {
			svc.log.Debug(
				"language skipped (not in LOCALE_LANGUAGES)",
				zap.Stringer("tag", lang.Tag),
			)

			continue
		}

		if !lang.Extends.IsRoot() {
			logFields = append(logFields, zap.Stringer("extends", lang.Extends))
		}

		if lang.src != "" {
			logFields = append(logFields, zap.String("imported", lang.src))
		} else {
			logFields = append(logFields, zap.Bool("embedded", true))
		}

		if err = loadTranslations(lang); err != nil {
			return err
		}

		if i == 0 && svc.def == nil {
			// set first one as default
			svc.def = lang
		}

		if svc.set[lang.Tag] != nil {
			svc.log.Info(
				"language overloaded",
				logFields...,
			)
		} else {
			svc.log.Info(
				"language loaded",
				logFields...,
			)
		}

		svc.set[lang.Tag] = lang
	}

	// Do another pass and link all extended languages
	for _, lang = range svc.set {
		if lang.Extends.IsRoot() {
			continue
		}

		if svc.set[lang.Extends] == nil {
			return fmt.Errorf("could not extend langage %q from an unknown language %q", lang.Tag, lang.Extends)
		}

		lang.extends = svc.set[lang.Extends]
	}

	return nil
}

// List returns list of all languages
func (svc *service) List() []*Language {
	svc.l.RLock()
	defer svc.l.RUnlock()

	ll := make([]*Language, 0, len(svc.set))
	for _, tag := range svc.tags {
		if svc.set[tag] != nil {
			ll = append(ll, svc.set[tag])
		}
	}

	return ll
}

// List localized list of all languages
func (svc *service) LocalizedList(ctx context.Context) []*Language {
	svc.l.RLock()
	defer svc.l.RUnlock()

	var (
		l       *Language
		ll      = make([]*Language, 0, len(svc.set))
		reqLang = GetLanguageFromContext(ctx)
	)

	for _, tag := range svc.tags {
		l = svc.set[tag]
		if l == nil {
			continue
		}

		l = &Language{
			Tag:           l.Tag,
			Name:          l.Name,
			LocalizedName: display.Languages(reqLang).Name(l.Tag),
		}

		if l.Name == "" {
			// for cases when language name is not configured
			l.Name = display.Languages(l.Tag).Name(l.Tag)
		}

		ll = append(ll, l)
	}

	return ll
}

func (svc *service) HasLanguage(lang language.Tag) bool {
	return svc.set[lang] != nil
}

func (svc *service) HasApplication(lang language.Tag, app string) bool {
	return svc.set[lang] != nil && svc.set[lang].external[app] != nil
}

// EncodeExternal writes into
func (svc *service) EncodeExternal(w io.Writer, app string, ll ...language.Tag) (err error) {
	svc.l.RLock()
	defer svc.l.RUnlock()

	_, _ = fmt.Fprint(w, "{")
	for i, lang := range ll {
		if i > 0 {
			_, _ = fmt.Fprint(w, ",")
		}

		_, _ = fmt.Fprintf(w, "%q:", lang)
		if svc.HasApplication(lang, app) {
			svc.set[lang].external[app].Seek(0, 0)
			_, _ = io.Copy(w, svc.set[lang].external[app])
		} else {
			_, _ = fmt.Fprint(w, "null")
		}
	}
	_, _ = fmt.Fprint(w, "}")

	return err
}

// NS returns translator function for the current language (from context) and a specific namespace
//
// Language is picked from the context
func (svc *service) NS(ctx context.Context, ns string) func(key string, rr ...string) string {
	var (
		code = GetLanguageFromContext(ctx)
	)

	return func(key string, rr ...string) string {
		return svc.t(code, ns, key, rr...)
	}
}

// T returns translated key from namespaces using list of replacement pairs
//
// Language is picked from the context
func (svc *service) T(ctx context.Context, ns, key string, rr ...string) string {
	return svc.t(GetLanguageFromContext(ctx), ns, key, rr...)
}

// Finds language and uses it to translate the given key
func (svc *service) t(code language.Tag, ns, key string, rr ...string) string {
	if svc != nil && svc.set != nil {
		if l, has := svc.set[code]; has {
			return l.t(ns, key, rr...)
		}
	}

	return key
}

func hasTag(t language.Tag, tt []language.Tag) bool {
	for _, tag := range tt {
		if tag.String() == t.String() {
			return true
		}
	}

	return false
}

func tagsToStrings(tt []language.Tag) []string {
	out := make([]string, len(tt))
	for t := range tt {
		out[t] = tt[t].String()
	}

	return out
}
