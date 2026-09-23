package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//  - {{ .Source }}

// builtInFunctionNames lists every function the expression language registers.
var builtInFunctionNames = map[string]struct{}{
{{- range .Names }}
	{{ printf "%q" . }}: {},
{{- end }}
}
