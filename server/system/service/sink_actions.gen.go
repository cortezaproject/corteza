package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/sink_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/system/types"
	"strings"
	"time"
)

type (
	sinkActionProps struct {
		url            string
		responseStatus int
		contentType    string
		sinkParams     *SinkRequestUrlParams
		mailHeader     *types.MailMessageHeader
	}

	sinkAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *sinkActionProps
	}

	sinkLogMetaKey   struct{}
	sinkPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setUrl updates sinkActionProps's url
//
// This function is auto-generated.
//
func (p *sinkActionProps) setUrl(url string) *sinkActionProps {
	p.url = url
	return p
}

// setResponseStatus updates sinkActionProps's responseStatus
//
// This function is auto-generated.
//
func (p *sinkActionProps) setResponseStatus(responseStatus int) *sinkActionProps {
	p.responseStatus = responseStatus
	return p
}

// setContentType updates sinkActionProps's contentType
//
// This function is auto-generated.
//
func (p *sinkActionProps) setContentType(contentType string) *sinkActionProps {
	p.contentType = contentType
	return p
}

// setSinkParams updates sinkActionProps's sinkParams
//
// This function is auto-generated.
//
func (p *sinkActionProps) setSinkParams(sinkParams *SinkRequestUrlParams) *sinkActionProps {
	p.sinkParams = sinkParams
	return p
}

// setMailHeader updates sinkActionProps's mailHeader
//
// This function is auto-generated.
//
func (p *sinkActionProps) setMailHeader(mailHeader *types.MailMessageHeader) *sinkActionProps {
	p.mailHeader = mailHeader
	return p
}

// Serialize converts sinkActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p sinkActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	m.Set("url", p.url, true)
	m.Set("responseStatus", p.responseStatus, true)
	m.Set("contentType", p.contentType, true)
	m.Set("sinkParams", p.sinkParams, true)
	if p.mailHeader != nil {
		m.Set("mailHeader.to", p.mailHeader.To, true)
		m.Set("mailHeader.CC", p.mailHeader.CC, true)
		m.Set("mailHeader.BCC", p.mailHeader.BCC, true)
		m.Set("mailHeader.from", p.mailHeader.From, true)
		m.Set("mailHeader.replyTo", p.mailHeader.ReplyTo, true)
		m.Set("mailHeader.raw", p.mailHeader.Raw, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p sinkActionProps) Format(in string, err error) string {
	var (
		pairs = []string{"{{err}}"}
		// first non-empty string
		fns = func(ii ...interface{}) string {
			for _, i := range ii {
				if s := fmt.Sprintf("%v", i); len(s) > 0 {
					return s
				}
			}

			return ""
		}
	)

	if err != nil {
		pairs = append(pairs, err.Error())
	} else {
		pairs = append(pairs, "nil")
	}
	pairs = append(pairs, "{{url}}", fns(p.url))
	pairs = append(pairs, "{{responseStatus}}", fns(p.responseStatus))
	pairs = append(pairs, "{{contentType}}", fns(p.contentType))
	pairs = append(pairs, "{{sinkParams}}", fns(p.sinkParams))

	if p.mailHeader != nil {
		// replacement for "{{mailHeader}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{mailHeader}}",
			fns(
				p.mailHeader.To,
				p.mailHeader.CC,
				p.mailHeader.BCC,
				p.mailHeader.From,
				p.mailHeader.ReplyTo,
				p.mailHeader.Raw,
			),
		)
		pairs = append(pairs, "{{mailHeader.to}}", fns(p.mailHeader.To))
		pairs = append(pairs, "{{mailHeader.CC}}", fns(p.mailHeader.CC))
		pairs = append(pairs, "{{mailHeader.BCC}}", fns(p.mailHeader.BCC))
		pairs = append(pairs, "{{mailHeader.from}}", fns(p.mailHeader.From))
		pairs = append(pairs, "{{mailHeader.replyTo}}", fns(p.mailHeader.ReplyTo))
		pairs = append(pairs, "{{mailHeader.raw}}", fns(p.mailHeader.Raw))
	}
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
//
func (a *sinkAction) String() string {
	var props = &sinkActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.Format(a.log, nil)
}

func (e *sinkAction) ToAction() *actionlog.Action {
	return &actionlog.Action{
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.Serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// SinkActionSign returns "system:sink.sign" action
//
// This function is auto-generated.
//
func SinkActionSign(props ...*sinkActionProps) *sinkAction {
	a := &sinkAction{
		timestamp: time.Now(),
		resource:  "system:sink",
		action:    "sign",
		log:       "signed sink request URL",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SinkActionPreprocess returns "system:sink.preprocess" action
//
// This function is auto-generated.
//
func SinkActionPreprocess(props ...*sinkActionProps) *sinkAction {
	a := &sinkAction{
		timestamp: time.Now(),
		resource:  "system:sink",
		action:    "preprocess",
		log:       "preprocess",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SinkActionRequest returns "system:sink.request" action
//
// This function is auto-generated.
//
func SinkActionRequest(props ...*sinkActionProps) *sinkAction {
	a := &sinkAction{
		timestamp: time.Now(),
		resource:  "system:sink",
		action:    "request",
		log:       "sink request processed",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// SinkErrGeneric returns "system:sink.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrGeneric(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:sink"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(sinkLogMetaKey{}, "{err}"),
		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrFailedToSign returns "system:sink.failedToSign" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrFailedToSign(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("could not sign request params: {{err}}", nil),

		errors.Meta("type", "failedToSign"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.failedToSign"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrFailedToCreateEvent returns "system:sink.failedToCreateEvent" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrFailedToCreateEvent(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to create sink event from request", nil),

		errors.Meta("type", "failedToCreateEvent"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.failedToCreateEvent"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrFailedToProcess returns "system:sink.failedToProcess" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrFailedToProcess(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to process request", nil),

		errors.Meta("type", "failedToProcess"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.failedToProcess"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrFailedToRespond returns "system:sink.failedToRespond" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrFailedToRespond(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to respond to request", nil),

		errors.Meta("type", "failedToRespond"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.failedToRespond"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrMissingSignature returns "system:sink.missingSignature" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrMissingSignature(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("missing sink signature parameter", nil),

		errors.Meta("type", "missingSignature"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.missingSignature"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrInvalidSignatureParam returns "system:sink.invalidSignatureParam" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrInvalidSignatureParam(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid sink signature parameter", nil),

		errors.Meta("type", "invalidSignatureParam"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.invalidSignatureParam"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrBadSinkParamEncoding returns "system:sink.badSinkParamEncoding" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrBadSinkParamEncoding(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("bad encoding of sink parameters", nil),

		errors.Meta("type", "badSinkParamEncoding"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.badSinkParamEncoding"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrInvalidSignature returns "system:sink.invalidSignature" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrInvalidSignature(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid signature", nil),

		errors.Meta("type", "invalidSignature"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.invalidSignature"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrInvalidSinkRequestUrlParams returns "system:sink.invalidSinkRequestUrlParams" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrInvalidSinkRequestUrlParams(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid sink request url params", nil),

		errors.Meta("type", "invalidSinkRequestUrlParams"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.invalidSinkRequestUrlParams"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrInvalidHttpMethod returns "system:sink.invalidHttpMethod" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrInvalidHttpMethod(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid HTTP method", nil),

		errors.Meta("type", "invalidHttpMethod"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.invalidHttpMethod"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrInvalidContentType returns "system:sink.invalidContentType" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrInvalidContentType(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid content-type header", nil),

		errors.Meta("type", "invalidContentType"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.invalidContentType"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrInvalidPath returns "system:sink.invalidPath" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrInvalidPath(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid path", nil),

		errors.Meta("type", "invalidPath"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.invalidPath"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrMisplacedSignature returns "system:sink.misplacedSignature" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrMisplacedSignature(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("signature misplaced", nil),

		errors.Meta("type", "misplacedSignature"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.misplacedSignature"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrSignatureExpired returns "system:sink.signatureExpired" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrSignatureExpired(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("signature expired", nil),

		errors.Meta("type", "signatureExpired"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.signatureExpired"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrContentLengthExceedsMaxAllowedSize returns "system:sink.contentLengthExceedsMaxAllowedSize" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrContentLengthExceedsMaxAllowedSize(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("content length exceeds max size limit", nil),

		errors.Meta("type", "contentLengthExceedsMaxAllowedSize"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.contentLengthExceedsMaxAllowedSize"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// SinkErrProcessingError returns "system:sink.processingError" as *errors.Error
//
//
// This function is auto-generated.
//
func SinkErrProcessingError(mm ...*sinkActionProps) *errors.Error {
	var p = &sinkActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("sink request process error", nil),

		errors.Meta("type", "processingError"),
		errors.Meta("resource", "system:sink"),

		errors.Meta(sinkPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "sink.errors.processingError"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// It will wrap unrecognized/internal errors with generic errors.
//
// This function is auto-generated.
//
func (svc sink) recordAction(ctx context.Context, props *sinkActionProps, actionFn func(...*sinkActionProps) *sinkAction, err error) error {
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
		a.Description = props.Format(m.AsString(sinkLogMetaKey{}), err)

		if p, has := m[sinkPropsMetaKey{}]; has {
			a.Meta = p.(*sinkActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
