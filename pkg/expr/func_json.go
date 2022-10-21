package expr

import (
	"encoding/json"
	"github.com/PaesslerAG/gval"
)

func JsonFunctions() []gval.Language {
	return []gval.Language{
		// TODO: Json decoding, parsing, stringify, JQ, JSONPath
		gval.Function("toJSON", toJSON),
	}
}

func toJSON(f interface{}) string {
	if _, is := f.(json.Marshaler); !is {
		f = UntypedValue(f)
	}
	b, _ := json.Marshal(f)
	return string(b)
}
