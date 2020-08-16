package {{ .Package }}

// This file is auto-generated.
//
// Template: pkg/store_interfaces_joined.gen.go.tpl
// Definitions:
{{- range .Definitions }}
//  - {{ .Source }}
{{- end }}
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

type (
	// Interface combines interfaces of all supported store interfaces
	storeGeneratedInterfaces interface {
	{{ range .Definitions -}}
		{{ unpubIdent .Types.Plural }}Store
	{{ end }}
	}
)
