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
	resourceStore interface {
		TransformResource(context.Context, language.Tag) (map[string]map[string]*ResourceTranslation, error)
	}

	Locale interface {
		T(ctx context.Context, ns, key string, rr ...string) string
		TFor(tag language.Tag, ns, key string, rr ...string) string
		Tags() []language.Tag
	}

	Resource interface {
		TResource(ctx context.Context, ns, key string, rr ...string) string
		TResourceFor(tag language.Tag, ns, key string, rr ...string) string
		Tags() []language.Tag
		SupportedLang(language.Tag) bool
		ResourceTranslations(code language.Tag, resource string) ResourceTranslationIndex
		Default() *Language
	}

	service struct {
		l sync.RWMutex
		s resourceStore

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

func Static(ll ...*Language) (svc *service) {
	svc = &service{
		log: zap.NewNop(),
		set: make(map[language.Tag]*Language),
	}

	for _, l := range ll {
		svc.set[l.Tag] = l
		svc.tags = append(svc.tags, l.Tag)
	}

	return svc
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

	return svc, svc.ReloadStatic()
}

func (svc *service) BindStore(s resourceStore) {
	svc.s = s
}

// Default language
func (svc *service) Default() *Language {
	return svc.def
}

// Tags of all loaded languages
func (svc *service) Tags() (tt []language.Tag) {
	svc.l.RLock()
	defer svc.l.RUnlock()
	return svc.tags
}

func (svc *service) SupportedLang(tag language.Tag) bool {
	return svc.set[tag] != nil
}

// ReloadStatic all language configurations (as configured via path options) and
// all translation files
func (svc *service) ReloadStatic() (err error) {
	var (
		i       int
		lang    *Language
		ll, aux []*Language

		logFields = make([]zap.Field, 0)
	)

	svc.l.Lock()
	defer svc.l.Unlock()

	svc.log.Info("reloading static",
		zap.Strings("path", svc.src),
		zap.Strings("tags", tagsToStrings(svc.tags)),
	)

	// In case where this is already initialized, we need to make sure to
	// preserve any resource translations.
	// The preservation is handled below, before we update svc.set.
	if svc.set == nil {
		svc.set = make(map[language.Tag]*Language)
	}

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

		// Preserve current resource translations.
		// Resource translations are updated when they change so these are
		// already up to date.
		if existing, ok := svc.set[lang.Tag]; ok {
			lang.resources = existing.resources
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

// ReloadResourceTranslations all language configurations (as configured via path options) and
// all translation files
func (svc *service) ReloadResourceTranslations(ctx context.Context) (err error) {
	if svc.s == nil {
		return fmt.Errorf("store for locale service not set")
	}

	svc.l.RLock()
	defer svc.l.RUnlock()

	svc.log.Info("reloading resource translations",
		zap.Strings("tags", tagsToStrings(svc.tags)),
	)

	for _, tag := range svc.tags {
		lang, ok := svc.set[tag]
		if !ok {
			lang = &Language{
				Tag: tag,
				// @todo find a better name for this, because it's not the
				//       language that it's loaded from the store, the resource translations are
				src: "store",
			}

			svc.set[tag] = lang
		}

		if err = svc.loadResourceTranslations(ctx, lang, lang.Tag); err != nil {
			return err
		}

		svc.log.Info(
			"resource translations loaded",
			zap.Stringer("tag", lang.Tag),
			zap.String("src", lang.src),
			zap.Int("translations", len(lang.resources)),
		)
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
		reqLang = GetAcceptLanguageFromContext(ctx)
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
		tag = GetAcceptLanguageFromContext(ctx)
	)

	return func(key string, rr ...string) string {
		return svc.t(tag, ns, key, rr...)
	}
}

// T returns translated key from namespaces using list of replacement pairs
//
// Language is picked from the context
func (svc *service) T(ctx context.Context, ns, key string, rr ...string) string {
	return svc.t(GetAcceptLanguageFromContext(ctx), ns, key, rr...)
}

// T returns translated key from namespaces using list of replacement pairs
//
// Language is specified
func (svc *service) TFor(tag language.Tag, ns, key string, rr ...string) string {
	return svc.t(tag, ns, key, rr...)
}

// TResource returns translated key for resource using list of replacement pairs
//
// Language is picked from the context
func (svc *service) TResource(ctx context.Context, ns, key string, rr ...string) string {

	return svc.tResource(GetAcceptLanguageFromContext(ctx), ns, key, rr...)
}

// TResourceFor returns translated key for resource using list of replacement pairs
//
// Language is picked from the context
func (svc *service) TResourceFor(tag language.Tag, ns, key string, rr ...string) string {
	return svc.tResource(tag, ns, key, rr...)
}

// ResourceTranslations returns all translations for the given language for the
// given resource.
//
// The response is indexed by translation key for nicer lookups.
func (svc *service) ResourceTranslations(tag language.Tag, resource string) ResourceTranslationIndex {
	out := make(ResourceTranslationIndex)

	if svc != nil && svc.set != nil {
		if l, has := svc.set[tag]; has {
			return l.resourceTranslations(resource)
		}
	}

	return out
}

// Finds language and uses it to translate the given key
func (svc *service) t(tag language.Tag, ns, key string, rr ...string) string {
	if svc != nil && svc.set != nil {
		if l, has := svc.set[tag]; has {
			return l.t(ns, key, rr...)
		}
	}

	// In case of a missing language, translation
	// key is returned
	return key
}

// Finds language and uses it to translate the given key for resource
func (svc *service) tResource(tag language.Tag, ns, key string, rr ...string) string {
	if svc != nil && svc.set != nil {
		if l, has := svc.set[tag]; has {
			return l.tResource(ns, key, rr...)
		}
	}

	// In case of a missing language, only empty string is
	// returned (in contrast to static translations
	// where a string with key is returned)
	return ""
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
