package locale

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

type (
	// Internal translations (application == corteza-server)
	// map contains namespace => key => text
	internal map[string]map[string]string

	// External translations (webapps)
	// map structure:
	//   application => namespace => buffered JSON docs
	external map[string]io.Reader

	Language struct {
		l sync.RWMutex

		Tag  language.Tag
		Name string

		// @todo
		Extends string
		// @todo
		base *Language

		// Internal translations (application == corteza-server)
		internal internal

		// External translations (webapps)
		external external
	}

	Languages struct {
		l sync.RWMutex

		log *zap.Logger

		opt options.LocaleOpt

		// sources
		src []string
		ll  map[language.Tag]*Language

		def language.Tag
	}

	ErrorMetaNamespace struct{}
	ErrorMetaKey       struct{}
)

func New(log *zap.Logger, opt options.LocaleOpt) (*Languages, error) {
	ll := &Languages{
		opt: opt,
		src: strings.Split(opt.Path, ":"),
		log: log.Named("locale"),
	}
	return ll, ll.Reload()
}

func (set *Languages) Default() language.Tag {
	return set.def
}

func (set *Languages) Tags() (tt []language.Tag) {
	tt = make([]language.Tag, 0, len(set.ll))
	for t := range set.ll {
		tt = append(tt, t)
	}

	return
}

func (set *Languages) Reload() error {
	set.l.RLock()
	defer set.l.RUnlock()

	set.log.Info("reloading", zap.Strings("src", set.src))

	set.ll = make(map[language.Tag]*Language)
	ll, err := load(set.log, set.src...)
	if err != nil {
		return err
	}

	for i, l := range ll {
		if i == 0 {
			// set first one as default
			set.def = l.Tag
		}

		set.ll[l.Tag] = l
	}

	return nil
}

// List returns list of all languages
func (set *Languages) List() []*Language {
	set.l.RLock()
	defer set.l.RUnlock()

	ll := make([]*Language, 0, len(set.ll))
	for _, l := range set.ll {
		ll = append(ll, &Language{
			Tag:  l.Tag,
			Name: l.Name,
		})
	}

	sort.SliceStable(ll, func(i, j int) bool {
		return strings.Compare(ll[i].Name, ll[j].Name) < 0
	})

	return ll
}

func (set *Languages) EncodeExternal(w io.Writer, lang language.Tag, app string) error {
	set.l.RLock()
	defer set.l.RUnlock()

	if set.ll[lang] == nil {
		return fmt.Errorf("language %q not found", lang)
	}

	if set.ll[lang].external[app] == nil {
		return fmt.Errorf("application %q not found", app)
	}

	_, err := io.Copy(w, set.ll[lang].external[app])
	return err
}

func (set *Languages) GetNS(ctx context.Context, ns string) func(key string, rr ...string) string {
	var (
		code = GetLanguageFromContext(ctx)
	)

	return func(key string, rr ...string) string {
		return set.get(code, ns, key, rr...)
	}
}

func (set *Languages) Current(ctx context.Context) language.Tag {
	return GetLanguageFromContext(ctx)
}

func (set *Languages) Get(ctx context.Context, ns, key string, rr ...string) string {
	return set.get(GetLanguageFromContext(ctx), ns, key, rr...)
}

func (set *Languages) get(code language.Tag, ns, key string, rr ...string) string {
	if set != nil && set.ll != nil {
		if l, has := set.ll[code]; has {
			return l.get(ns, key, rr...)
		}
	}

	return key
}

func (l *Language) get(ns, key string, rr ...string) string {
	l.l.RLock()
	defer l.l.RUnlock()

	for r := 0; r < len(rr); r += 2 {
		rr[r] = fmt.Sprintf("{{%s}}", rr[r])
	}

	msg, has := l.internal[ns][key]
	if has {
		return strings.NewReplacer(rr...).Replace(msg)
	}

	if l.base != nil {
		return l.base.get(ns, key, rr...)
	}

	return key
}
