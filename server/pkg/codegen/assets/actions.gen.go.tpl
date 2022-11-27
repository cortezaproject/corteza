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
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
{{- range .Import }}
    {{ normalizeImport . }}
{{- end }}
	"strings"
	"time"

)

type (
	{{ $.Service }}ActionProps struct {
	{{- range $prop := $.Props }}
		{{ $prop.Name }} {{ $prop.Type }}
	{{- end }}
	}

{{ if $.Actions }}
	{{ $.Service }}Action struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *{{ $.Service }}ActionProps
	}
{{ end }}

	{{ $.Service }}LogMetaKey struct {}
	{{ $.Service }}PropsMetaKey struct {}
)

var (
    // just a placeholder to cover template cases w/o fmt package use
    _ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods

{{- range $prop := $.Props }}
// {{ camelCase "set" $prop.Name }} updates {{ $.Service }}ActionProps's {{ $prop.Name }}
//
// This function is auto-generated.
//
func (p *{{ $.Service }}ActionProps) {{ camelCase "set" $prop.Name }}({{ $prop.Name }} {{ $prop.Type }}) *{{ $.Service }}ActionProps {
    p.{{ $prop.Name }} = {{ $prop.Name }}
    return p
}
{{ end }}


// Serialize converts {{ $.Service }}ActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p {{ $.Service }}ActionProps) Serialize() actionlog.Meta {
	var (
		m   = make(actionlog.Meta)
	)

{{ range $prop := $.Props }}
	{{- if $prop.Builtin }}
		m.Set("{{ $prop.Name }}", p.{{ $prop.Name }}, true)
	{{- else }}
	if p.{{ $prop.Name }} != nil {
    {{- range $f := $prop.Fields }}
		m.Set("{{ $prop.Name }}.{{ $f }}", p.{{ $prop.Name }}.{{ camelCase " " $f }}, true)
    {{- end }}
    }
    {{- end }}
{{- end }}

    return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p {{ $.Service }}ActionProps) Format(in string, err error) string {
    var (
        pairs = []string{"{{"{{"}}err}}"}

{{- if $.Props }}
        // first non-empty string
        fns = func(ii ... interface{}) string {
			for _, i:= range ii {
				if s :=fmt.Sprintf("%v", i); len(s) > 0 {
					return s
				}
			}

			return ""
		}
{{- end }}
    )

	if err != nil {
	    pairs = append(pairs, err.Error())
	} else {
	    pairs = append(pairs, "nil")
	}

{{- range $prop := $.Props }}
	{{- if $prop.Builtin }}
	    pairs = append(pairs, "{{"{{"}}{{ $prop.Name }}}}", fns(p.{{ $prop.Name }}))
	{{- else }}

	if p.{{ $prop.Name }} != nil {
	    // replacement for "{{"{{"}}{{ $prop.Name }}}}" (in order how fields are defined)
        pairs = append(
            pairs,
            "{{"{{"}}{{ $prop.Name }}}}",
            fns(
            {{- range $f := $prop.Fields }}
                p.{{ $prop.Name }}.{{ camelCase " " $f }},
            {{- end }}
            ),
        )

    {{- range $f := $prop.Fields }}
        pairs = append(pairs, "{{"{{"}}{{ $prop.Name }}.{{ $f }}}}", fns(p.{{ $prop.Name }}.{{ camelCase " " $f }}))
    {{- end }}
    }
    {{- end }}
{{- end }}
    return strings.NewReplacer(pairs...).Replace(in)
}

{{ if $.Actions }}
// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
//
func (a *{{ $.Service }}Action) String() string {
    var props = &{{ $.Service }}ActionProps{}

    if a.props != nil {
        props = a.props
    }

    return props.Format(a.log, nil)
}

func (e *{{ $.Service }}Action) ToAction() *actionlog.Action {
	return &actionlog.Action{
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.Serialize(),
	}
}
{{ end }}

{{ if $.Actions }}
// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

{{ range $a := $.Actions }}
// {{ camelCase "" $.Service "Action" $a.Action  }} returns "{{ $.Resource }}.{{ $a.Action  }}" action
//
// This function is auto-generated.
//
func {{ camelCase "" $.Service "Action" $a.Action  }}(props ... *{{ $.Service }}ActionProps) *{{ $.Service }}Action {
	a := &{{ $.Service }}Action{
		timestamp:    time.Now(),
		resource:     "{{ $.Resource }}",
		action:       "{{ $a.Action }}",
		log:          "{{ $a.Log }}",
		severity:     {{ $a.SeverityConstName }},
	}

	if len(props) > 0 {
	    a.props = props[0]
	}

	return a
}
{{ end }}
{{ end }}

{{ if $.Errors }}
// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

{{ range $e := $.Errors }}
// {{ camelCase "" $.Service "Err" $e.Error }} returns "{{ $.Resource }}.{{ $e.Error  }}" as *errors.Error
//
{{- if $e.MaskedWith }}
// Note: This error will be wrapped with safe ({{ $.Resource }}.{{ $e.MaskedWith }}) error!
{{- end }}
//
// This function is auto-generated.
//
func {{ camelCase "" $.Service "Err" $e.Error }}(mm ... *{{ $.Service }}ActionProps) *errors.Error {
	var p = &{{ $.Service }}ActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		{{ if $e.Message }}p.Format({{ printf "%q" $e.Message }}, nil){{ else }}{{ printf "%q" $e.Error }}{{ end }},

		errors.Meta("type",          {{ printf "%q" $e.Error         }}),
		errors.Meta("resource",      {{ printf "%q" $.Resource       }}),

		{{ if $e.Documentation }}
		// link to documentation; formatting applies in case we need some special link formatting
		errors.Meta("documentation", p.Format({{ printf "%q" $e.Documentation }}, nil)),
		{{ end -}}

		{{ if $e.Details }}
		// details, used in detailed eror reporting
		errors.Meta("details", p.Format({{ printf "%q" $e.Details  }}, nil)),
		{{ end -}}


		{{- if $e.Log }}
		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta({{ $.Service }}LogMetaKey{},           {{ printf "%q" $e.Log           }}),
		{{ end -}}
		errors.Meta({{ $.Service }}PropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, {{ printf "%q" $.Component }}),
		errors.Meta(locale.ErrorMetaKey{}, "{{ kebabCase $.Service }}.errors.{{ $e.Error }}"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	{{ if $e.MaskedWith }}
	// Wrap with safe error
	e = {{ camelCase "" $.Service "Err" $e.MaskedWith }}().Wrap(e)
	{{ end }}

	return e
}
{{ end }}
{{ end }}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// It will wrap unrecognized/internal errors with generic errors.
//
// This function is auto-generated.
//
func (svc {{ $.Service }}) recordAction(ctx context.Context, props *{{ $.Service }}ActionProps, actionFn func(... *{{ $.Service }}ActionProps) *{{ $.Service }}Action, err error) error {
	if svc.actionlog == nil || actionFn == nil {
		// action log disabled or no action fn passed, return error as-is
		return err
	} else if err == nil {
		// action completed w/o error, record it
		svc.actionlog.Record(ctx, actionFn(props).ToAction())
		return nil
	}

	a := actionFn(props).ToAction()


	// Extracting error information and recording it as action
	a.Error = err.Error()

	switch c := err.(type) {
	case *errors.Error:
		m := c.Meta()

		a.Error = err.Error()
		a.Severity = actionlog.Severity(m.AsInt("severity"))
		a.Description = props.Format(m.AsString({{ $.Service }}LogMetaKey{}), err)

		if p, has := m[{{ $.Service }}PropsMetaKey{}]; has {
			a.Meta = p.(*{{ $.Service }}ActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}


	// Original error is passed on
	return err
}
