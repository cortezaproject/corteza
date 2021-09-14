package locale

import (
	"fmt"
	"io"
	"io/fs"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

type (
	// Internal translations (application == corteza-server)
	// map contains namespace => key => text
	internal map[string]map[string]string

	// External translations (webapps)
	// map structure:
	//   application => namespace => buffered JSON docs
	external map[string]io.ReadSeeker

	// resource => key => ResourceTranslation
	resources map[string]map[string]*ResourceTranslation

	Language struct {
		l sync.RWMutex

		// location of the language files
		// this is mainly for logging/debugging purposes
		//
		// This value is empty when loading
		// embedded languages
		src string

		// pointer to the place we loaded files from
		fs fs.FS

		Tag           language.Tag `json:"tag,string"`
		Name          string       `json:"name"`
		LocalizedName string       `json:"localizedName"`

		// Language (tag) that this language is extending
		// @todo
		Extends language.Tag `json:"-"`

		// Points to extended language
		// @todo
		extends *Language

		// Internal translations (application == corteza-server)
		internal internal

		// External translations (webapps)
		external  external
		resources resources
	}

	ErrorMetaNamespace struct{}
	ErrorMetaKey       struct{}
)

// t returns the translated string for internal content
func (l *Language) t(ns, key string, rr ...string) string {
	l.l.RLock()
	defer l.l.RUnlock()

	for r := 0; r < len(rr); r += 2 {
		rr[r] = fmt.Sprintf("{{%s}}", rr[r])
	}

	msg, has := l.internal[ns][key]
	if has {
		return strings.NewReplacer(rr...).Replace(msg)
	}

	if l.extends != nil {
		return l.extends.t(ns, key, rr...)
	}

	return key
}

// tr returns the translated string for resource translations
func (l *Language) tr(ns, key string, rr ...string) string {
	l.l.RLock()
	defer l.l.RUnlock()

	for r := 0; r < len(rr); r += 2 {
		rr[r] = fmt.Sprintf("{{%s}}", rr[r])
	}

	rt, has := l.resources[ns][key]
	if has {
		return strings.NewReplacer(rr...).Replace(rt.Msg)
	}

	if l.extends != nil {
		return l.extends.tr(ns, key, rr...)
	}

	return key
}

// resourceTranslations returns all resource translations for the specified resource
func (l *Language) resourceTranslations(resource string) ResourceTranslationIndex {
	out := make(ResourceTranslationIndex)
	if l.resources == nil {
		return out
	}

	for k, rt := range l.resources[resource] {
		out[k] = rt
	}

	return out
}
