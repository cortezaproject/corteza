package errors

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cast"
)

type (
	meta map[interface{}]interface{}
)

// StringKeys returns max length of (string) keys and slice of all strings
func (m meta) StringKeys() (int, []string) {
	var (
		l, ml int
		key   string
		kk    = make([]string, 0, len(m))
	)

	// collecting keys so that we can sort them to ensure
	// stable order and find out max length of the key
	for k := range m {
		switch c := k.(type) {
		case string:
			key = c
		case fmt.Stringer:
			key = c.String()
		default:
			continue
		}

		kk = append(kk, key)
		if l = len(key); l > ml {
			ml = l
		}
	}

	sort.Strings(kk)
	return ml, kk
}

func (m meta) AsString(key interface{}) string {
	if _, has := m[key]; has {
		if s, ok := m[key].(string); ok {
			return s
		}
	}

	return ""
}

func (m meta) AsInt(key interface{}) int {
	if _, has := m[key]; has {
		if s, ok := m[key].(int); ok {
			return s
		}
	}

	return 0
}

func (m meta) MarshalJSON() ([]byte, error) {
	o := make(map[string]interface{})
	_, kk := m.StringKeys()
	for _, k := range kk {
		o[k] = m[k]
	}

	return json.Marshal(o)
}

// returns string pairs (both key & value are strings=
func (m meta) pairs() []string {
	out := make([]string, 0, len(m)*2)

	for k := range m {
		if key, ok := k.(string); ok {
			out = append(out, key, cast.ToString(m[k]))
		}
	}

	return out
}
