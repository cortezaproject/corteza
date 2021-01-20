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
	"context"
	"fmt"
{{- range .Imports }}
  {{ normalizeImport . }}
{{- end }}
{{- if ne .Package "expr" }}
	. "github.com/cortezaproject/corteza-server/pkg/expr"
{{- end }}
)
{{ end }}


var _ = context.Background
var _ = fmt.Errorf


{{ range $exprType, $def := .Types }}
// {{ $exprType }} is an expression type, wrapper for {{ $def.As }} type
type {{ $exprType }} struct{ value {{ $def.As }} }

// New{{ $exprType }} creates new instance of {{ $exprType }} expression type
func New{{ $exprType }}(val interface{}) (*{{ $exprType }}, error) {
	if c, err := {{ export "CastTo" $exprType }}(val); err != nil {
		return nil, fmt.Errorf("unable to create {{ $exprType }}: %w", err)
	} else {
		return &{{ $exprType }}{value: c}, nil
	}
}


// Return underlying value on {{ $exprType }}
func (t {{ $exprType }}) Get() interface{}                         { return t.value }

// Return type name
func ({{ $exprType }}) Type() string                               { return "{{ $.Prefix }}{{ $exprType }}" }

// Convert value to {{ $def.As }}
func ({{ $exprType }}) Cast(val interface{}) (TypedValue, error) {
	return New{{ $exprType }}(val)
}

// Assign new value to {{ $exprType }}
//
// value is first passed through {{ export "CastTo" $exprType }}
func (t *{{ $exprType }}) Assign(val interface{}) (error) {
	if c, err := {{ export "CastTo" $exprType }}(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}


{{ if $def.Struct }}
func (t *{{ $exprType }}) AssignFieldValue(key string, val interface{}) error {
	return {{ $def.AssignerFn }}(t.value, key, val)
}

{{ if not $def.CustomGValSelector }}
// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access {{ $exprType }}'s underlying value ({{ $def.As }})
// and it's fields
//
func (t {{ $exprType }}) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return {{ unexport $exprType "GValSelector" }}(t.value, k)
}
{{ end }}

// Select is field accessor for {{ $def.As }}
//
// Similar to SelectGVal but returns typed values
func (t {{ $exprType }}) Select(k string) (TypedValue, error) {
	return {{ unexport $exprType "TypedValueSelector" }}(t.value, k)
}

func (t {{ $exprType }}) Has(k string) bool {
	switch k {
	{{- range $def.Struct }}
		{{- if .ExprType }}
		case {{ printf "%q" .Name }}{{ if .Alias }}, {{ printf "%q" .Alias }}{{ end }}:
			return true
		{{- end }}
	{{- end }}
	}
	return false
}

// {{ unexport $exprType "GValSelector" }} is field accessor for {{ $def.As }}
func {{ unexport $exprType "GValSelector" }}(res {{ $def.As }}, k string) (interface{}, error) {
	switch k {
	{{- range $def.Struct }}
		{{- if .ExprType }}
		case {{ printf "%q" .Name }}{{ if .Alias }}, {{ printf "%q" .Alias }}{{ end }}:
			return res.{{ export .Name }}, nil
		{{- end }}
	{{- end }}
	}

	return nil, fmt.Errorf("unknown field '%s'", k)}

// {{ unexport $exprType "TypedValueSelector" }} is field accessor for {{ $def.As }}
func {{ unexport $exprType "TypedValueSelector" }}(res {{ $def.As }}, k string) (TypedValue, error) {
	switch k {
	{{- range $def.Struct }}
		{{- if .ExprType }}
		case {{ printf "%q" .Name }}{{ if .Alias }}, {{ printf "%q" .Alias }}{{ end }}:
			return {{ export "New" .ExprType }}(res.{{ export .Name }})
		{{- end }}
	{{- end }}
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

{{ if $def.BuiltInAssignerFn }}
// {{ $def.AssignerFn }} is field value setter for {{ $def.As }}
func {{ $def.AssignerFn }}(res {{ $def.As }}, k string, val interface{}) (error) {
	switch k {
{{- range $def.Struct }}
	case {{ printf "%q" .Name }}{{ if .Alias }}, {{ printf "%q" .Alias }}{{ end }}:
	{{- if .Readonly }}
		return fmt.Errorf("field '%s' is read-only", k)
	{{- else }}
		aux, err := {{ export "CastTo" .ExprType }}(val)
		if err != nil {
			return err
		}

		res.{{ export .Name }} = aux
		return nil
	{{- end }}
{{- end }}
	}

	return fmt.Errorf("unknown field '%s'", k)
}
{{ end }}
{{ end }}

{{ end }}

