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
	b, _ := json.Marshal(f)
	return string(b)
}
