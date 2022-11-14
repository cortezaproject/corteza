package actionlog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	metaFoo struct{ val string }
	metaBar struct{ val string }
)

func (f metaFoo) ActionLogMetaValue() (interface{}, bool) {
	return strings.ToUpper(f.val), len(f.val) == 0
}

func TestMeta_Set(t *testing.T) {
	tests := []struct {
		name      string
		val       interface{}
		omitempty bool
		out       Meta
	}{
		{"string value, omit", "str", true, Meta{"t": "str"}},
		{"string value, keep", "str", false, Meta{"t": "str"}},
		{"string empty, omit", "", true, Meta{}},
		{"string empty, keep", "", false, Meta{"t": ""}},

		{"byte value, omit", []byte{'b', 'y', 't', 'e'}, true, Meta{"t": []byte{'b', 'y', 't', 'e'}}},
		{"byte value, keep", []byte{'b', 'y', 't', 'e'}, false, Meta{"t": []byte{'b', 'y', 't', 'e'}}},
		{"byte empty, omit", []byte{}, true, Meta{}},
		{"byte empty, keep", []byte{}, false, Meta{"t": []byte{}}},

		{"bool true", true, false, Meta{"t": true}},
		{"bool false", false, false, Meta{"t": false}},

		{"int value, omit", 1, true, Meta{"t": "1"}},
		{"int value, keep", 1, false, Meta{"t": "1"}},
		{"int empty, omit", 0, true, Meta{}},
		{"int empty, keep", 0, false, Meta{"t": "0"}},
		{"uint64 empty, omit", uint64(0), true, Meta{}},
		{"uint64 empty, keep", uint64(0), false, Meta{"t": "0"}},
		{"bigint", 244783268048994492, false, Meta{"t": "244783268048994492"}},
		{"int64", int64(244783268048994492), false, Meta{"t": "244783268048994492"}},
		{"uint64()", uint64(244783268048994492), false, Meta{"t": "244783268048994492"}},
		{"int", 42, false, Meta{"t": "42"}},
		{"float64", float64(42), false, Meta{"t": float64(42)}},
		{"float64", float64(244783268048994492), false, Meta{"t": 2.447832680489945e+17}},

		{"slice", []int{0, 1}, false, Meta{"t": []int{0, 1}}},

		{"meta foo", metaFoo{"value"}, false, Meta{"t": "VALUE"}},
		{"meta foo, omit", metaFoo{""}, true, Meta{}},
		// @todo raw
		{"meta bar", metaBar{"value"}, false, Meta{"t": metaBar{"value"}}},
		{"meta bar", metaBar{}, true, Meta{"t": metaBar{""}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			m := make(Meta)
			m.Set("t", tt.val, tt.omitempty)
			if !a.Equal(tt.out, m) {
				t.Fail()
			}
		})
	}
}
