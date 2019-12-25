package {{ .Package }}

{{ if $.Events.Properties }}
import (
	"encoding/json"
)
{{ end }}

// Match returns false if given conditions do not match event & resource internals
func (res {{ camelCase .ResourceIdent "base" }}) Match(name string, op string, values ...string) bool {
	// By default we match no mather what kind of constraints we receive
	//
	// Function will be called multiple times - once for every trigger constraint
	// All should match (return true):
	//   constraint#1 AND constraint#2 AND constraint#3 ...
	//
	// When there are multiple values, Match() can decide how to treat them (OR, AND...)
	return true
}

// Encode internal data to be passed as event params & arguments to triggered Corredor script
func (res {{ camelCase .ResourceIdent "base" }}) Encode() (args map[string][]byte, err error) {
	{{- if $.Events.Properties }}
	args = make(map[string][]byte)

	{{ range $prop := $.Events.Properties }}
	if args["{{ $prop.Name }}"], err = json.Marshal(res.{{ $prop.Name }}); err != nil {
		return nil, err
	}
	{{ end }}
	{{ else }}
	// Handle argument encoding
	{{ end -}}
	return
}

// Decode return values from Corredor script into struct props
func (res *{{ camelCase .ResourceIdent "base" }}) Decode(results map[string][]byte)( err error) {
	{{- if $.Events.Result }}
	if r, ok := results["result"]; ok && len(results) == 1 {
		if err = json.Unmarshal(r, res.{{ $.Events.Result }}); err != nil {
			return
		}
	}
	{{ end -}}

	{{- range $prop := $.Events.Properties }}
	{{- if not $prop.Immutable }}
	if r, ok := results["{{ $prop.Name }}"]; ok && len(results) == 1 {
		if err = json.Unmarshal(r, res.{{ $prop.Name }}); err != nil {
			return
		}
	}
	{{ end -}}
	{{ end -}}

	return
}
