package locale

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
	"golang.org/x/text/language"
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
			return nil, fmt.Errorf("failed to parse language '%s'", lang)
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

// Reload all language configurations (as configured via path options) and
// all translation files
func (svc *service) Reload() error {
	svc.l.RLock()
	defer svc.l.RUnlock()

	svc.log.Info("reloading",
		zap.Strings("path", svc.src),
		zap.Strings("tags", tagsToStrings(svc.tags)),
	)

	svc.set = make(map[language.Tag]*Language)
	configs, err := loadConfigs(svc.src...)
	if err != nil {
		return err
	}

	for i, lang := range configs {
		if !hasTag(lang.Tag, svc.tags) {
			svc.log.Info(
				"language skipped (see LOCALE_LANGUAGES and LOCALE_PATH)",
				zap.Stringer("tag", lang.Tag),
			)

			continue
		}

		if err = loadTranslations(lang, lang.src); err != nil {
			return err
		}

		svc.log.Info(
			"language loaded",
			zap.Stringer("tag", lang.Tag),
			zap.String("src", lang.src),
			zap.Stringer("extends", lang.Extends),
		)

		if i == 0 && svc.def == nil {
			// set first one as default
			svc.def = lang
		}

		svc.set[lang.Tag] = lang
	}

	return nil
}

// List returns list of all languages
func (svc *service) List() []*Language {
	svc.l.RLock()
	defer svc.l.RUnlock()

	ll := make([]*Language, 0, len(svc.set))
	for _, tag := range svc.tags {
		ll = append(ll, svc.set[tag])
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
func (svc *service) EncodeExternal(w io.Writer, lang language.Tag, app string) (err error) {
	svc.l.RLock()
	defer svc.l.RUnlock()

	if svc.HasApplication(lang, app) {
		_, err = io.Copy(w, svc.set[lang].external[app])
	} else {
		err = fmt.Errorf("application or language missing")
	}

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
