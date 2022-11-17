package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

import (
	"context"
	"fmt"
	"sync"
{{- range .Imports }}
  {{ normalizeImport . }}
{{- end }}
{{- if ne .Package "expr" }}
	. "github.com/cortezaproject/corteza/server/pkg/expr"
{{- end }}
)

var _ = context.Background
var _ = fmt.Errorf

{{ range $exprType, $def := .Types }}
{{ if not $def.CustomType }}
// {{ $exprType }} is an expression type, wrapper for {{ $def.As }} type
type {{ $exprType }} struct{
	value {{ $def.As }}
	mux sync.RWMutex
}

// New{{ $exprType }} creates new instance of {{ $exprType }} expression type
func New{{ $exprType }}(val interface{}) (*{{ $exprType }}, error) {
	if c, err := {{ export "CastTo" $exprType }}(val); err != nil {
		return nil, fmt.Errorf("unable to create {{ $exprType }}: %w", err)
	} else {
		return &{{ $exprType }}{value: c}, nil
	}
}


// Get return underlying value on {{ $exprType }}
func (t *{{ $exprType }}) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on {{ $exprType }}
func (t *{{ $exprType }}) GetValue()  {{ $def.As }} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func ({{ $exprType }}) Type() string                               { return "{{ $exprType }}" }

// Cast converts value to {{ $def.As }}
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

{{ if $def.Comparable }}
{{ if not $def.CustomComparator }}
// Compare the two {{ $exprType }} values
func (t {{ $exprType }}) Compare(to TypedValue) (int, error) {
	c, err := New{{ $exprType }}(to)
	if err != nil {
		return 0, fmt.Errorf("cannot compare %s and %s: %s", t.Type(), c.Type(), err.Error())
	}

	switch {
	case t.value == c.value:
		return 0, nil
	case t.value < c.value:
		return -1, nil
	case t.value > c.value:
		return 1, nil
	default:
		return 0, fmt.Errorf("cannot compare %s and %s: unknown state", t.Type(), c.Type())
	}
}
{{ else }}
// Compare the two {{ $exprType }} values
func (t {{ $exprType }}) Compare(to TypedValue) (int, error) {
	return compareTo{{ $exprType }}(t, to)
}
{{ end }}
{{ end }}

{{ if $def.Struct }}
{{ if not $def.CustomFieldAssigner }}
func (t *{{ $exprType }}) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return {{ $def.AssignerFn }}(t.value, key, val)
}
{{ end }}

{{ if not $def.CustomGValSelector }}
// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access {{ $exprType }}'s underlying value ({{ $def.As }})
// and it's fields
//
func (t *{{ $exprType }}) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return {{ unexport $exprType "GValSelector" }}(t.value, k)
}
{{ end }}

{{ if not $def.CustomSelector }}
// Select is field accessor for {{ $def.As }}
//
// Similar to SelectGVal but returns typed values
func (t *{{ $exprType }}) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return {{ unexport $exprType "TypedValueSelector" }}(t.value, k)
}
{{ end }}

func (t *{{ $exprType }}) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
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
	{{- if hasPtr $def.As }}
	if res == nil {
		return nil, nil
	}

	{{- end }}
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
	{{- if hasPtr $def.As }}
	if res == nil {
		return nil, nil
	}

	{{- end }}
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
{{ end }} {{/* if $def.Struct */}}
{{ end }} {{/* if not $def.CustomType */}}
{{ end }} {{/* types loop */}}

