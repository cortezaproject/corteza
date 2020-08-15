package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}

type (
	// Interface combines interfaces of all supported store interfaces
	storeInterface interface {
	{{ range .Definitions -}}
		{{ unpubIdent .Types.Plural }}Store
	{{ end }}
	}
)
