package {{ .Package }}

// This file is auto-generated from {{ .YAML }}
//

import (
	"context"
	"fmt"
	"strings"
	"errors"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
{{- range $import := $.Import }}
    {{ $import }}
{{- end }}
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

{{ if $.Errors }}
	{{ $.Service }}Error struct {
		timestamp time.Time
		error       string
		resource    string
		action      string
		message 	string
		log         string
		severity    actionlog.Severity

		wrap        error

		props       *{{ $.Service }}ActionProps
	}
{{ end }}
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods

{{- range $prop := $.Props }}
// {{ camelCase "set" $prop.Name }} updates {{ $.Service }}ActionProps's {{ $prop.Name }}
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *{{ $.Service }}ActionProps) {{ camelCase "set" $prop.Name }}({{ $prop.Name }} {{ $prop.Type }}) *{{ $.Service }}ActionProps {
    p.{{ $prop.Name }} = {{ $prop.Name }}
    return p
}
{{ end }}


// serialize converts {{ $.Service }}ActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p {{ $.Service }}ActionProps) serialize() actionlog.Meta {
	var (
		m   = make(actionlog.Meta)
		str = func(i interface{}) string { return fmt.Sprintf("%v", i) }
	)

{{ range $prop := $.Props }}
	{{- if $prop.Builtin }}
	    m["{{ $prop.Name }}"] = str(p.{{ $prop.Name }})
	{{- else }}
	if p.{{ $prop.Name }} != nil {
    {{- range $f := $prop.Fields }}
        m["{{ $prop.Name }}.{{ $f }}"] = str(p.{{ $prop.Name }}.{{ camelCase " " $f }})
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
func (p {{ $.Service }}ActionProps) tr(in string, err error) string {
    var pairs = []string{"{err}"}

	if err != nil {
		for {
			// Unwrap errors
			ue := errors.Unwrap(err)
			if ue == nil {
				break
			}

			err = ue
		}

	    pairs = append(pairs, err.Error())
	} else {
	    pairs = append(pairs, "nil")
	}

{{- range $prop := $.Props }}
	{{- if $prop.Builtin }}
	    pairs = append(pairs, "{{"{"}}{{ $prop.Name }}}", fmt.Sprintf("%v", p.{{ $prop.Name }}))
	{{- else }}

	if p.{{ $prop.Name }} != nil {
	{{- if $prop.DefaultField }}
        pairs = append(pairs, "{{"{"}}{{ $prop.Name }}}", fmt.Sprintf("%v", p.{{ $prop.Name }}.{{ camelCase " " $prop.DefaultField }}))
	{{- end }}

    {{- range $f := $prop.Fields }}
        pairs = append(pairs, "{{"{"}}{{ $prop.Name }}.{{ $f }}}", fmt.Sprintf("%v", p.{{ $prop.Name }}.{{ camelCase " " $f }}))

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

    return props.tr(a.log, nil)
}

func (e *{{ $.Service }}Action) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.serialize(),
	}
}
{{ end }}

{{ if $.Errors }}
// *********************************************************************************************************************
// *********************************************************************************************************************
// Error methods

// String returns loggable description as string
//
// It falls back to message if log is not set
//
// This function is auto-generated.
//
func (e *{{ $.Service }}Error) String() string {
    var props = &{{ $.Service }}ActionProps{}

    if e.props != nil {
        props = e.props
    }


    if e.wrap != nil && !strings.Contains(e.log, "{err}") {
        // Suffix error log with {err} to ensure
        // we log the cause for this error
        e.log += ": {err}"
    }

	return props.tr(e.log, e.wrap)
}

// Error satisfies
//
// This function is auto-generated.
//
func (e *{{ $.Service }}Error) Error() string {
    var props = &{{ $.Service }}ActionProps{}

    if e.props != nil {
        props = e.props
    }

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *{{ $.Service }}Error) Is(Resource error) bool {
	t, ok := Resource.(*{{ $.Service }}Error)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps {{ $.Service }}Error around another error
//
// This function is auto-generated.
//
func (e *{{ $.Service }}Error) Wrap(err error) *{{ $.Service }}Error {
    e.wrap = err
    return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *{{ $.Service }}Error) Unwrap() error {
	return e.wrap
}

func (e *{{ $.Service }}Error) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Error:       e.Error(),
		Meta:        e.props.serialize(),
	}
}
{{ end }}

{{ if $.Actions }}
// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

{{ range $a := $.Actions }}
// {{ camelCase "" $.Service "Action" $a.Action  }} returns "{{ $.Resource }}.{{ $a.Action  }}" error
//
// This function is auto-generated.
//
func {{ camelCase "" $.Service "Action" $a.Action  }}(props ... *{{ $.Service }}ActionProps) *{{ $.Service }}Action {
	a := &{{ $.Service }}Action{
		timestamp:    time.Now(),
		resource:     "{{ $.Resource }}",
		action:       "{{ $a.Action }}",
		log:   "{{ $a.Log }}",
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
{{- if $e.Safe }}
// {{ camelCase "" $.Service "Err" $e.Error }} returns "{{ $.Resource }}.{{ $e.Safe  }}" audit event as {{ $e.SeverityConstName }}
{{- else }}
// {{ camelCase "" $.Service "Err" $e.Error }} returns "{{ $.Resource }}.{{ $e.Error  }}" audit event as {{ $e.SeverityConstName }}
{{- end }}
//
{{- if $e.Safe }}
// Note: This error will be wrapped with safe ({{ $e.Safe }}) error!
{{- end }}
//
// This function is auto-generated.
//
func {{ camelCase "" $.Service "Err" $e.Error }}(props ... *{{ $.Service }}ActionProps) *{{ $.Service }}Error {
	var e = &{{ $.Service }}Error{
		timestamp: time.Now(),
		resource:  "{{ $.Resource }}",
		error:     "{{ $e.Error }}",
		action:    "error",
		message:   "{{ $e.Message }}",
		log:       "{{ $e.Log }}",
		severity:  {{ $e.SeverityConstName }},
		props:     func() *{{ $.Service }}ActionProps { if len(props) > 0 { return props[0] }; return nil}(),
	}

	if len(props) > 0 {
	    e.props = props[0]
	}

	{{ if $e.Safe }}
	// Wrap with safe error
	return {{ camelCase "" $.Service "Err" $e.Safe }}().Wrap(e)
	{{ else }}
	return e
	{{ end }}
}
{{ end }}
{{ end }}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// context is used to enrich audit log entry with current user info, request ID, IP address...
// props are collected action/error properties
// action (optional) fn will be used to construct {{ $.Service }}Action struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc {{ $.Service }}) recordAction(ctx context.Context, props *{{ $.Service }}ActionProps, action func(... *{{ $.Service }}ActionProps) *{{ $.Service }}Action, err error) error {
	var (
		ok bool

		// Return error
		retError *{{ $.Service }}Error

		// Recorder error
		recError *{{ $.Service }}Error
	)

	if err != nil {
		if retError, ok = err.(*{{ $.Service }}Error); !ok {
			// got non-{{ $.Service }} error, wrap it with {{ camelCase "" $.Service "err" "generic" }}
			retError = {{ camelCase "" $.Service "err" "generic" }}(props).Wrap(err)

			// copy action to returning and recording error
			retError.action = action().action

			// we'll use {{ camelCase "" $.Service "err" "generic" }} for recording too
			// because it can hold more info
			recError = retError
		} else if retError != nil {
			// copy action to returning and recording error
			retError.action = action().action
			// start with copy of return error for recording
			// this will be updated with tha root cause as we try and
			// unwrap the error
			recError = retError

			// find the original recError for this error
			// for the purpose of logging
			var unwrappedError error = retError
			for {
				if unwrappedError = errors.Unwrap(unwrappedError); unwrappedError == nil {
					// nothing wrapped
					break
				}

				// update recError ONLY of wrapped error is of type {{ $.Service }}Error
				if unwrappedSinkError, ok := unwrappedError.(*{{ $.Service }}Error); ok {
					recError = unwrappedSinkError
				}
			}

			if retError.props == nil {
				// set props on returning error if empty
				retError.props = props
			}

			if recError.props == nil {
				// set props on recording error if empty
				recError.props = props
			}
		}
	}

	if svc.actionlog != nil {
		if retError != nil {
			// failed action, log error
			svc.actionlog.Record(ctx, recError)
		} else if action != nil {
			// successful
			svc.actionlog.Record(ctx, action(props))
		}
	}

	if err == nil {
		// retError not an interface and that WILL (!!) cause issues
		// with nil check (== nil) when it is not explicitly returned
		return nil
	}

	return retError
}
