package service

// This file is auto-generated from compose/service/notifications_actions.yaml
//

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
)

type (
	notificationActionProps struct {
		mail           *types.EmailNotification
		recipient      string
		attachmentURL  string
		attachmentSize int64
		attachmentType string
	}

	notificationAction struct {
		timestamp time.Time
		resource  string
		action    string
		log       string
		severity  actionlog.Severity

		// prefix for error when action fails
		errorMessage string

		props *notificationActionProps
	}

	notificationError struct {
		timestamp time.Time
		error     string
		resource  string
		action    string
		message   string
		log       string
		severity  actionlog.Severity

		wrap error

		props *notificationActionProps
	}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setMail updates notificationActionProps's mail
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *notificationActionProps) setMail(mail *types.EmailNotification) *notificationActionProps {
	p.mail = mail
	return p
}

// setRecipient updates notificationActionProps's recipient
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *notificationActionProps) setRecipient(recipient string) *notificationActionProps {
	p.recipient = recipient
	return p
}

// setAttachmentURL updates notificationActionProps's attachmentURL
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *notificationActionProps) setAttachmentURL(attachmentURL string) *notificationActionProps {
	p.attachmentURL = attachmentURL
	return p
}

// setAttachmentSize updates notificationActionProps's attachmentSize
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *notificationActionProps) setAttachmentSize(attachmentSize int64) *notificationActionProps {
	p.attachmentSize = attachmentSize
	return p
}

// setAttachmentType updates notificationActionProps's attachmentType
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *notificationActionProps) setAttachmentType(attachmentType string) *notificationActionProps {
	p.attachmentType = attachmentType
	return p
}

// serialize converts notificationActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p notificationActionProps) serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	if p.mail != nil {
		m.Set("mail.subject", p.mail.Subject, true)
		m.Set("mail.to", p.mail.To, true)
		m.Set("mail.cc", p.mail.Cc, true)
		m.Set("mail.replyTo", p.mail.ReplyTo, true)
		m.Set("mail.contentPlain", p.mail.ContentPlain, true)
		m.Set("mail.contentHTML", p.mail.ContentHTML, true)
		m.Set("mail.remoteAttachments", p.mail.RemoteAttachments, true)
	}
	m.Set("recipient", p.recipient, true)
	m.Set("attachmentURL", p.attachmentURL, true)
	m.Set("attachmentSize", p.attachmentSize, true)
	m.Set("attachmentType", p.attachmentType, true)

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p notificationActionProps) tr(in string, err error) string {
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

	if p.mail != nil {
		// replacement for "{mail}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{mail}",
			fns(
				p.mail.Subject,
				p.mail.To,
				p.mail.Cc,
				p.mail.ReplyTo,
				p.mail.ContentPlain,
				p.mail.ContentHTML,
				p.mail.RemoteAttachments,
			),
		)
		pairs = append(pairs, "{mail.subject}", fns(p.mail.Subject))
		pairs = append(pairs, "{mail.to}", fns(p.mail.To))
		pairs = append(pairs, "{mail.cc}", fns(p.mail.Cc))
		pairs = append(pairs, "{mail.replyTo}", fns(p.mail.ReplyTo))
		pairs = append(pairs, "{mail.contentPlain}", fns(p.mail.ContentPlain))
		pairs = append(pairs, "{mail.contentHTML}", fns(p.mail.ContentHTML))
		pairs = append(pairs, "{mail.remoteAttachments}", fns(p.mail.RemoteAttachments))
	}
	pairs = append(pairs, "{recipient}", fns(p.recipient))
	pairs = append(pairs, "{attachmentURL}", fns(p.attachmentURL))
	pairs = append(pairs, "{attachmentSize}", fns(p.attachmentSize))
	pairs = append(pairs, "{attachmentType}", fns(p.attachmentType))
	return strings.NewReplacer(pairs...).Replace(in)
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Action methods

// String returns loggable description as string
//
// This function is auto-generated.
//
func (a *notificationAction) String() string {
	var props = &notificationActionProps{}

	if a.props != nil {
		props = a.props
	}

	return props.tr(a.log, nil)
}

func (e *notificationAction) LoggableAction() *actionlog.Action {
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
func (e *notificationError) String() string {
	var props = &notificationActionProps{}

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
func (e *notificationError) Error() string {
	var props = &notificationActionProps{}

	if e.props != nil {
		props = e.props
	}

	return props.tr(e.message, e.wrap)
}

// Is fn for error equality check
//
// This function is auto-generated.
//
func (e *notificationError) Is(Resource error) bool {
	t, ok := Resource.(*notificationError)
	if !ok {
		return false
	}

	return t.resource == e.resource && t.error == e.error
}

// Wrap wraps notificationError around another error
//
// This function is auto-generated.
//
func (e *notificationError) Wrap(err error) *notificationError {
	e.wrap = err
	return e
}

// Unwrap returns wrapped error
//
// This function is auto-generated.
//
func (e *notificationError) Unwrap() error {
	return e.wrap
}

func (e *notificationError) LoggableAction() *actionlog.Action {
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

// NotificationActionSend returns "compose:notification.send" error
//
// This function is auto-generated.
//
func NotificationActionSend(props ...*notificationActionProps) *notificationAction {
	a := &notificationAction{
		timestamp: time.Now(),
		resource:  "compose:notification",
		action:    "send",
		log:       "notification sent",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// NotificationActionAttachmentDownload returns "compose:notification.attachmentDownload" error
//
// This function is auto-generated.
//
func NotificationActionAttachmentDownload(props ...*notificationActionProps) *notificationAction {
	a := &notificationAction{
		timestamp: time.Now(),
		resource:  "compose:notification",
		action:    "attachmentDownload",
		log:       "attachment downloaded",
		severity:  actionlog.Debug,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// *********************************************************************************************************************
// *********************************************************************************************************************
// Error constructors

// NotificationErrGeneric returns "compose:notification.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NotificationErrGeneric(props ...*notificationActionProps) *notificationError {
	var e = &notificationError{
		timestamp: time.Now(),
		resource:  "compose:notification",
		error:     "generic",
		action:    "error",
		message:   "failed to complete request due to internal error",
		log:       "{err}",
		severity:  actionlog.Error,
		props: func() *notificationActionProps {
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

// NotificationErrFailedToLoadUser returns "compose:notification.failedToLoadUser" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NotificationErrFailedToLoadUser(props ...*notificationActionProps) *notificationError {
	var e = &notificationError{
		timestamp: time.Now(),
		resource:  "compose:notification",
		error:     "failedToLoadUser",
		action:    "error",
		message:   "could not load user for {recipient}",
		log:       "could not load user for {recipient}",
		severity:  actionlog.Error,
		props: func() *notificationActionProps {
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

// NotificationErrInvalidReceipientFormat returns "compose:notification.invalidReceipientFormat" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NotificationErrInvalidReceipientFormat(props ...*notificationActionProps) *notificationError {
	var e = &notificationError{
		timestamp: time.Now(),
		resource:  "compose:notification",
		error:     "invalidReceipientFormat",
		action:    "error",
		message:   "invalid recipient format ({recipient})",
		log:       "invalid recipient format ({recipient})",
		severity:  actionlog.Error,
		props: func() *notificationActionProps {
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

// NotificationErrNoRecipients returns "compose:notification.noRecipients" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NotificationErrNoRecipients(props ...*notificationActionProps) *notificationError {
	var e = &notificationError{
		timestamp: time.Now(),
		resource:  "compose:notification",
		error:     "noRecipients",
		action:    "error",
		message:   "can not send email message without recipients",
		log:       "can not send email message without recipients",
		severity:  actionlog.Error,
		props: func() *notificationActionProps {
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

// NotificationErrFailedToDownloadAttachment returns "compose:notification.failedToDownloadAttachment" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func NotificationErrFailedToDownloadAttachment(props ...*notificationActionProps) *notificationError {
	var e = &notificationError{
		timestamp: time.Now(),
		resource:  "compose:notification",
		error:     "failedToDownloadAttachment",
		action:    "error",
		message:   "could not download attachment from {attachmentURL}: {err}",
		log:       "could not download attachment from {attachmentURL}: {err}",
		severity:  actionlog.Error,
		props: func() *notificationActionProps {
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
// action (optional) fn will be used to construct notificationAction struct from given props (and error)
// err is any error that occurred while action was happening
//
// Action has success and fail (error) state:
//  - when recorded without an error (4th param), action is recorded as successful.
//  - when an additional error is given (4th param), action is used to wrap
//    the additional error
//
// This function is auto-generated.
//
func (svc notification) recordAction(ctx context.Context, props *notificationActionProps, action func(...*notificationActionProps) *notificationAction, err error) error {
	var (
		ok bool

		// Return error
		retError *notificationError

		// Recorder error
		recError *notificationError
	)

	if err != nil {
		if retError, ok = err.(*notificationError); !ok {
			// got non-notification error, wrap it with NotificationErrGeneric
			retError = NotificationErrGeneric(props).Wrap(err)

			if action != nil {
				// copy action to returning and recording error
				retError.action = action().action
			}

			// we'll use NotificationErrGeneric for recording too
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

				// update recError ONLY of wrapped error is of type notificationError
				if unwrappedSinkError, ok := unwrappedError.(*notificationError); ok {
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
