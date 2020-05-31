package service

// This file is auto-generated from system/service/sink_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/system/types"
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

	sinkError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *sinkActionProps

		httpStatusCode int
	}
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
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sinkActionProps) setUrl(url string) *sinkActionProps {
	p.url = url
	return p
}

// setResponseStatus updates sinkActionProps's responseStatus
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sinkActionProps) setResponseStatus(responseStatus int) *sinkActionProps {
	p.responseStatus = responseStatus
	return p
}

// setContentType updates sinkActionProps's contentType
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sinkActionProps) setContentType(contentType string) *sinkActionProps {
	p.contentType = contentType
	return p
}

// setSinkParams updates sinkActionProps's sinkParams
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sinkActionProps) setSinkParams(sinkParams *SinkRequestUrlParams) *sinkActionProps {
	p.sinkParams = sinkParams
	return p
}

// setMailHeader updates sinkActionProps's mailHeader
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *sinkActionProps) setMailHeader(mailHeader *types.MailMessageHeader) *sinkActionProps {
	p.mailHeader = mailHeader
	return p
}

// serialize converts sinkActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p sinkActionProps) serialize() actionlog.Meta {
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
func (p sinkActionProps) tr(in string, err error) string {
	var (
		pairs = []string{"{err}"}
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
	pairs = append(pairs, "{url}", fns(p.url))
	pairs = append(pairs, "{responseStatus}", fns(p.responseStatus))
	pairs = append(pairs, "{contentType}", fns(p.contentType))
	pairs = append(pairs, "{sinkParams}", fns(p.sinkParams))

	if p.mailHeader != nil {
		// replacement for "{mailHeader}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{mailHeader}",
			fns(
				p.mailHeader.To,
				p.mailHeader.CC,
				p.mailHeader.BCC,
				p.mailHeader.From,
				p.mailHeader.ReplyTo,
				p.mailHeader.Raw,
			),
		)
		pairs = append(pairs, "{mailHeader.to}", fns(p.mailHeader.To))
		pairs = append(pairs, "{mailHeader.CC}", fns(p.mailHeader.CC))
		pairs = append(pairs, "{mailHeader.BCC}", fns(p.mailHeader.BCC))
		pairs = append(pairs, "{mailHeader.from}", fns(p.mailHeader.From))
		pairs = append(pairs, "{mailHeader.replyTo}", fns(p.mailHeader.ReplyTo))
		pairs = append(pairs, "{mailHeader.raw}", fns(p.mailHeader.Raw))
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

	return props.tr(a.log, nil)
}

func (e *sinkAction) LoggableAction() *actionlog.Action {
	return &actionlog.Action{
		Timestamp:   e.timestamp,
		Resource:    e.resource,
		Action:      e.action,
		Severity:    e.severity,
		Description: e.String(),
		Meta:        e.props.serialize(),
	}
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error methods

// String returns loggable description as string
//
// It falls back to message if log is not set
//
// This function is auto-generated.
//
func (e *sinkError) String() string {
	var props = &sinkActionProps{}

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
func (e *sinkError) Error() string {
	var props = &sinkActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *sinkError) Is(Resource error) bool {
	t, ok := Resource.(*sinkError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps sinkError around another error
//
// This function is auto-generated.
//
func (e *sinkError) Wrap(err error) *sinkError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *sinkError) Unwrap() error {
	return e.wrap
}

func (e *sinkError) LoggableAction() *actionlog.Action {
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

func (e *sinkError) HttpResponse(w http.ResponseWriter) {
	var code = e.httpStatusCode
	if code == 0 {
		code = http.StatusInternalServerError
	}

	http.Error(w, e.message, code)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// SinkActionSign returns "system:sink.sign" error
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

// SinkActionPreprocess returns "system:sink.preprocess" error
//
// This function is auto-generated.
//
func SinkActionPreprocess(props ...*sinkActionProps) *sinkAction {
	a := &sinkAction{
		timestamp: time.Now(),
		resource:  "system:sink",
		action:    "preprocess",
		log:       "preprocess",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// SinkActionRequest returns "system:sink.request" error
//
// This function is auto-generated.
//
func SinkActionRequest(props ...*sinkActionProps) *sinkAction {
	a := &sinkAction{
		timestamp: time.Now(),
		resource:  "system:sink",
		action:    "request",
		log:       "sink request processed",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// SinkErrGeneric returns "system:sink.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SinkErrGeneric(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrFailedToSign returns "system:sink.failedToSign" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SinkErrFailedToSign(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "failedToSign",
		action:    "error",
		message:   "could not sign request params: {err}",
		log:       "could not sign request params: {err}",
		severity:  actionlog.Error,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrFailedToCreateEvent returns "system:sink.failedToCreateEvent" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SinkErrFailedToCreateEvent(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "failedToCreateEvent",
		action:    "error",
		message:   "failed to create sink event from request",
		log:       "failed to create sink event from request",
		severity:  actionlog.Error,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrFailedToProcess returns "system:sink.failedToProcess" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SinkErrFailedToProcess(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "failedToProcess",
		action:    "error",
		message:   "failed to process request",
		log:       "failed to process request",
		severity:  actionlog.Error,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrFailedToRespond returns "system:sink.failedToRespond" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func SinkErrFailedToRespond(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "failedToRespond",
		action:    "error",
		message:   "failed to respond to request",
		log:       "failed to respond to request",
		severity:  actionlog.Error,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrMissingSignature returns "system:sink.missingSignature" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrMissingSignature(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "missingSignature",
		action:    "error",
		message:   "missing sink signature parameter",
		log:       "missing sink signature parameter",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusBadRequest,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrInvalidSignatureParam returns "system:sink.invalidSignatureParam" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrInvalidSignatureParam(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "invalidSignatureParam",
		action:    "error",
		message:   "invalid sink signature parameter",
		log:       "invalid sink signature parameter",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusUnauthorized,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrBadSinkParamEncoding returns "system:sink.badSinkParamEncoding" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrBadSinkParamEncoding(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "badSinkParamEncoding",
		action:    "error",
		message:   "bad encoding of sink parameters",
		log:       "bad encoding of sink parameters",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusBadRequest,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrInvalidSignature returns "system:sink.invalidSignature" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrInvalidSignature(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "invalidSignature",
		action:    "error",
		message:   "invalid signature",
		log:       "invalid signature",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusUnauthorized,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrInvalidSinkRequestUrlParams returns "system:sink.invalidSinkRequestUrlParams" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrInvalidSinkRequestUrlParams(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "invalidSinkRequestUrlParams",
		action:    "error",
		message:   "invalid sink request url params",
		log:       "invalid sink request url params",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusInternalServerError,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrInvalidHttpMethod returns "system:sink.invalidHttpMethod" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrInvalidHttpMethod(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "invalidHttpMethod",
		action:    "error",
		message:   "invalid HTTP method",
		log:       "invalid HTTP method",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusUnauthorized,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrInvalidContentType returns "system:sink.invalidContentType" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrInvalidContentType(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "invalidContentType",
		action:    "error",
		message:   "invalid content-type header",
		log:       "invalid content-type header",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusUnauthorized,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrInvalidPath returns "system:sink.invalidPath" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrInvalidPath(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "invalidPath",
		action:    "error",
		message:   "invalid path",
		log:       "invalid path",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusUnauthorized,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrMisplacedSignature returns "system:sink.misplacedSignature" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrMisplacedSignature(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "misplacedSignature",
		action:    "error",
		message:   "signature misplaced",
		log:       "signature misplaced",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusBadRequest,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrSignatureExpired returns "system:sink.signatureExpired" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrSignatureExpired(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "signatureExpired",
		action:    "error",
		message:   "signature expired",
		log:       "signature expired",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusGone,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrContentLengthExceedsMaxAllowedSize returns "system:sink.contentLengthExceedsMaxAllowedSize" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrContentLengthExceedsMaxAllowedSize(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "contentLengthExceedsMaxAllowedSize",
		action:    "error",
		message:   "content length exceeds max size limit",
		log:       "content length exceeds max size limit",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusRequestEntityTooLarge,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// SinkErrProcessingError returns "system:sink.processingError" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func SinkErrProcessingError(props ...*sinkActionProps) *sinkError {
	var e = &sinkError{
		timestamp: time.Now(),
		resource:  "system:sink",
		error:     "processingError",
		action:    "error",
		message:   "sink request process error",
		log:       "sink request process error",
		severity:  actionlog.Alert,
		props: func() *sinkActionProps {
			if len(props) > 0 {
				return props[0]
			}
			return nil
		}(),

		httpStatusCode: http.StatusInternalServerError,
	}

	if len(props) > 0 {
		e.props = props[0]
	}

	return e

}

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// context is used to enrich audit log entry with current user info, request ID, IP address...
// props are collected action/error properties
// action (optional) fn will be used to construct sinkAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc sink) recordAction(ctx context.Context, props *sinkActionProps, action func(...*sinkActionProps) *sinkAction, err error) error {
	var (
		ok bool

		// Return error
		retError *sinkError

		// Recorder error
		recError *sinkError
	)

	if err != nil {
		if retError, ok = err.(*sinkError); !ok {
			// got non-sink error, wrap it with SinkErrGeneric
			retError = SinkErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use SinkErrGeneric for recording too
			// because it can hold more info
			recError = retError
		} else if retError != nil {
			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}
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

				// update recError ONLY of wrapped error is of type sinkError
				if unwrappedSinkError, ok := unwrappedError.(*sinkError); ok {
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
