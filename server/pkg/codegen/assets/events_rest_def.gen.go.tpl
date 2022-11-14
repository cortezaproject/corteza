package rest

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.

func getEventTypeDefinitions() []eventTypeDef {
	return []eventTypeDef{
	{{ range .Definitions }}
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
		Constraints: []eventTypeConstraintDef{
		{{ range $.res.Constraints }}
			{
				Name: {{ printf "%q" .Name }},
			},
		{{ end }}
		},
	},
	{{ end }}
{{ end }}
