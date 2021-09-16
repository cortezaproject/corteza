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
		external external
	}

	ErrorMetaNamespace struct{}
	ErrorMetaKey       struct{}
)

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
