package rest

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

func getEventTypeDefinitions() []eventTypeDef {
	return []eventTypeDef{
	{{ range .Definitions }}
		// {{ printf "%#v" .Imports }}
		{{ range $r := .Resources }}
		{{ template "eventTypeDefinitions" dict "res" $r "type" "on"     "types" .On          }}
		{{ template "eventTypeDefinitions" dict "res" $r "type" "before" "types" .BeforeAfter }}
		{{ template "eventTypeDefinitions" dict "res" $r "type" "after"  "types" .BeforeAfter }}
		{{ end }}
	{{ end }}
	}
}

{{ define "eventTypeDefinitions" }}
	{{ range $ev := $.types }}
	{
		ResourceType: {{ printf "%q" $.res.ResourceString }},
		EventType: {{ printf "%q"  (camelCase $.type $ev) }},
		Properties: []eventTypePropertyDef{
		{{ range $.res.Properties }}
		{{ if not .Internal }}
			{
				Name: {{ printf "%q" .Name }},
				Type: {{ printf "%q" .ExprType }},
				Immutable: {{ printf "%v" .Immutable }},
			},
		{{ end }}
		{{ end }}
		},
	},
	{{ end }}
{{ end }}
