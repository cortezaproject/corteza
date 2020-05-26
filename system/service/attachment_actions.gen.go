package service

// This file is auto-generated from system/service/attachment_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	attachmentActionProps struct {
		size       int64
		name       string
		mimetype   string
		url        string
		attachment *types.Attachment
		filter     *types.AttachmentFilter
	}

	attachmentAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *attachmentActionProps
	}

	attachmentError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *attachmentActionProps
	}
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setSize updates attachmentActionProps's size
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setSize(size int64) *attachmentActionProps {
	p.size = size
	return p
}

// setName updates attachmentActionProps's name
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setName(name string) *attachmentActionProps {
	p.name = name
	return p
}

// setMimetype updates attachmentActionProps's mimetype
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setMimetype(mimetype string) *attachmentActionProps {
	p.mimetype = mimetype
	return p
}

// setUrl updates attachmentActionProps's url
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setUrl(url string) *attachmentActionProps {
	p.url = url
	return p
}

// setAttachment updates attachmentActionProps's attachment
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setAttachment(attachment *types.Attachment) *attachmentActionProps {
	p.attachment = attachment
	return p
}

// setFilter updates attachmentActionProps's filter
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setFilter(filter *types.AttachmentFilter) *attachmentActionProps {
	p.filter = filter
	return p
}

// serialize converts attachmentActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p attachmentActionProps) serialize() actionlog.Meta {
	var (
		m   = make(actionlog.Meta)
		str = func(i interface{}) string { return fmt.Sprintf("%v", i) }
	)

	// avoiding declared but not used
	_ = str

	m["size"] = str(p.size)
	m["name"] = str(p.name)
	m["mimetype"] = str(p.mimetype)
	m["url"] = str(p.url)
	if p.attachment != nil {
		m["attachment.name"] = str(p.attachment.Name)
		m["attachment.kind"] = str(p.attachment.Kind)
		m["attachment.url"] = str(p.attachment.Url)
		m["attachment.previewUrl"] = str(p.attachment.PreviewUrl)
		m["attachment.meta"] = str(p.attachment.Meta)
		m["attachment.ownerID"] = str(p.attachment.OwnerID)
		m["attachment.ID"] = str(p.attachment.ID)
	}
	if p.filter != nil {
		m["filter.filter"] = str(p.filter.Filter)
		m["filter.kind"] = str(p.filter.Kind)
		m["filter.sort"] = str(p.filter.Sort)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p attachmentActionProps) tr(in string, err error) string {
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
	pairs = append(pairs, "{size}", fmt.Sprintf("%v", p.size))
	pairs = append(pairs, "{name}", fmt.Sprintf("%v", p.name))
	pairs = append(pairs, "{mimetype}", fmt.Sprintf("%v", p.mimetype))
	pairs = append(pairs, "{url}", fmt.Sprintf("%v", p.url))

	if p.attachment != nil {
		pairs = append(pairs, "{attachment}", fmt.Sprintf("%v", p.attachment.Name))
		pairs = append(pairs, "{attachment.name}", fmt.Sprintf("%v", p.attachment.Name))
		pairs = append(pairs, "{attachment.kind}", fmt.Sprintf("%v", p.attachment.Kind))
		pairs = append(pairs, "{attachment.url}", fmt.Sprintf("%v", p.attachment.Url))
		pairs = append(pairs, "{attachment.previewUrl}", fmt.Sprintf("%v", p.attachment.PreviewUrl))
		pairs = append(pairs, "{attachment.meta}", fmt.Sprintf("%v", p.attachment.Meta))
		pairs = append(pairs, "{attachment.ownerID}", fmt.Sprintf("%v", p.attachment.OwnerID))
		pairs = append(pairs, "{attachment.ID}", fmt.Sprintf("%v", p.attachment.ID))
	}

	if p.filter != nil {
		pairs = append(pairs, "{filter}", fmt.Sprintf("%v", p.filter.Filter))
		pairs = append(pairs, "{filter.filter}", fmt.Sprintf("%v", p.filter.Filter))
		pairs = append(pairs, "{filter.kind}", fmt.Sprintf("%v", p.filter.Kind))
		pairs = append(pairs, "{filter.sort}", fmt.Sprintf("%v", p.filter.Sort))
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
func (a *attachmentAction) String() string {
	var props = &attachmentActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *attachmentAction) LoggableAction() *actionlog.Action {
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
func (e *attachmentError) String() string {
	var props = &attachmentActionProps{}

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
func (e *attachmentError) Error() string {
	var props = &attachmentActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *attachmentError) Is(Resource error) bool {
	t, ok := Resource.(*attachmentError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps attachmentError around another error
//
// This function is auto-generated.
//
func (e *attachmentError) Wrap(err error) *attachmentError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *attachmentError) Unwrap() error {
	return e.wrap
}

func (e *attachmentError) LoggableAction() *actionlog.Action {
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

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action constructors

// AttachmentActionSearch returns "system:attachment.search" error
//
// This function is auto-generated.
//
func AttachmentActionSearch(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "system:attachment",
		action:    "search",
		log:       "searched for attachments",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionLookup returns "system:attachment.lookup" error
//
// This function is auto-generated.
//
func AttachmentActionLookup(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "system:attachment",
		action:    "lookup",
		log:       "looked-up for a {attachment}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionCreate returns "system:attachment.create" error
//
// This function is auto-generated.
//
func AttachmentActionCreate(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "system:attachment",
		action:    "create",
		log:       "created {attachment}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionDelete returns "system:attachment.delete" error
//
// This function is auto-generated.
//
func AttachmentActionDelete(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "system:attachment",
		action:    "delete",
		log:       "deleted {attachment}",
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

// AttachmentErrGeneric returns "system:attachment.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func AttachmentErrGeneric(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "system:attachment",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *attachmentActionProps {
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

// AttachmentErrNonexistent returns "system:attachment.nonexistent" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrNonexistent(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "system:attachment",
		error:     "nonexistent",
		action:    "error",
		message:   "attachment does not exist",
		log:       "attachment does not exist",
		severity:  actionlog.Warning,
		props: func() *attachmentActionProps {
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

// AttachmentErrInvalidID returns "system:attachment.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidID(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "system:attachment",
		error:     "invalidID",
		action:    "error",
		message:   "invalid ID",
		log:       "invalid ID",
		severity:  actionlog.Warning,
		props: func() *attachmentActionProps {
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

// AttachmentErrNotAllowedToListAttachments returns "system:attachment.notAllowedToListAttachments" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToListAttachments(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "system:attachment",
		error:     "notAllowedToListAttachments",
		action:    "error",
		message:   "not allowed to list attachments",
		log:       "failed to list attachment; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *attachmentActionProps {
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

// AttachmentErrNotAllowedToCreate returns "system:attachment.notAllowedToCreate" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToCreate(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "system:attachment",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create attachment",
		log:       "failed to create attachment; insufficient permissions",
		severity:  actionlog.Alert,
		props: func() *attachmentActionProps {
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

// AttachmentErrFailedToExtractMimeType returns "system:attachment.failedToExtractMimeType" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrFailedToExtractMimeType(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "system:attachment",
		error:     "failedToExtractMimeType",
		action:    "error",
		message:   "could not extract mime type",
		log:       "could not extract mime type",
		severity:  actionlog.Alert,
		props: func() *attachmentActionProps {
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

// AttachmentErrFailedToStoreFile returns "system:attachment.failedToStoreFile" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrFailedToStoreFile(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "system:attachment",
		error:     "failedToStoreFile",
		action:    "error",
		message:   "could not extract store file",
		log:       "could not extract store file",
		severity:  actionlog.Alert,
		props: func() *attachmentActionProps {
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

// AttachmentErrFailedToProcessImage returns "system:attachment.failedToProcessImage" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrFailedToProcessImage(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "system:attachment",
		error:     "failedToProcessImage",
		action:    "error",
		message:   "could not process image",
		log:       "could not process image",
		severity:  actionlog.Alert,
		props: func() *attachmentActionProps {
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

// *********************************************************************************************************************
// *********************************************************************************************************************

// recordAction is a service helper function wraps function that can return error
//
// context is used to enrich audit log entry with current user info, request ID, IP address...
// props are collected action/error properties
// action (optional) fn will be used to construct attachmentAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc attachment) recordAction(ctx context.Context, props *attachmentActionProps, action func(...*attachmentActionProps) *attachmentAction, err error) error {
	var (
		ok bool

		// Return error
		retError *attachmentError

		// Recorder error
		recError *attachmentError
	)

	if err != nil {
		if retError, ok = err.(*attachmentError); !ok {
			// got non-attachment error, wrap it with AttachmentErrGeneric
			retError = AttachmentErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use AttachmentErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type attachmentError
				if unwrappedSinkError, ok := unwrappedError.(*attachmentError); ok {
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
