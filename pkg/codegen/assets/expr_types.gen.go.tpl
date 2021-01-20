package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

{{ if .Imports }}
import (
{{- range .Imports }}
  {{ normalizeImport . }}
{{- end }}
{{- if ne .Package "expr" }}
	"github.com/cortezaproject/corteza-server/pkg/expr"
{{- end }}
)
{{ end }}


{{ $TypedValue := "TypedValue" }}
{{ if ne .Package "expr" }}
	{{ $TypedValue = "expr.TypedValue" }}
{{ end }}

{{ range $exprType, $nativeType := .Types }}
// {{ $exprType }} is an expression type, wrapper for {{ $nativeType }} type
type {{ $exprType }} struct{ value {{ $nativeType }} }

// New{{ $exprType }} creates new instance of {{ $exprType }} expression type
func New{{ $exprType }}(new interface{}) ({{ $TypedValue }}, error) {
	t := &{{ $exprType }}{}
	return t, t.Set(new)
}

// Returns underlying value on {{ $exprType }}
func (t {{ $exprType }}) Get() interface{}                         { return t.value }

// Returns type name
func ({{ $exprType }}) Type() string                               { return "{{ $exprType }}" }

// Casts value to {{ $nativeType }}
func ({{ $exprType }}) Cast(value interface{}) ({{ $TypedValue }}, error) { return New{{ $exprType }}(value) }

{{ end }}
