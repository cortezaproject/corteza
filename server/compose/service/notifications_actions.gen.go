package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/service/notifications_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/actionlog"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"strings"
	"time"
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

	notificationLogMetaKey   struct{}
	notificationPropsMetaKey struct{}
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

// Serialize converts notificationActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p notificationActionProps) Serialize() actionlog.Meta {
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
func (p notificationActionProps) Format(in string, err error) string {
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

	if p.mail != nil {
		// replacement for "{{mail}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{mail}}",
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
		pairs = append(pairs, "{{mail.subject}}", fns(p.mail.Subject))
		pairs = append(pairs, "{{mail.to}}", fns(p.mail.To))
		pairs = append(pairs, "{{mail.cc}}", fns(p.mail.Cc))
		pairs = append(pairs, "{{mail.replyTo}}", fns(p.mail.ReplyTo))
		pairs = append(pairs, "{{mail.contentPlain}}", fns(p.mail.ContentPlain))
		pairs = append(pairs, "{{mail.contentHTML}}", fns(p.mail.ContentHTML))
		pairs = append(pairs, "{{mail.remoteAttachments}}", fns(p.mail.RemoteAttachments))
	}
	pairs = append(pairs, "{{recipient}}", fns(p.recipient))
	pairs = append(pairs, "{{attachmentURL}}", fns(p.attachmentURL))
	pairs = append(pairs, "{{attachmentSize}}", fns(p.attachmentSize))
	pairs = append(pairs, "{{attachmentType}}", fns(p.attachmentType))
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

	return props.Format(a.log, nil)
}

func (e *notificationAction) ToAction() *actionlog.Action {
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

// NotificationActionSend returns "compose:notification.send" action
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

// NotificationActionAttachmentDownload returns "compose:notification.attachmentDownload" action
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

// NotificationErrGeneric returns "compose:notification.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func NotificationErrGeneric(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "compose:notification"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(notificationLogMetaKey{}, "{err}"),
		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrFailedToLoadUser returns "compose:notification.failedToLoadUser" as *errors.Error
//
//
// This function is auto-generated.
//
func NotificationErrFailedToLoadUser(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("could not load user for {{recipient}}", nil),

		errors.Meta("type", "failedToLoadUser"),
		errors.Meta("resource", "compose:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.failedToLoadUser"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrInvalidReceipientFormat returns "compose:notification.invalidReceipientFormat" as *errors.Error
//
//
// This function is auto-generated.
//
func NotificationErrInvalidReceipientFormat(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid recipient format ({{recipient}})", nil),

		errors.Meta("type", "invalidReceipientFormat"),
		errors.Meta("resource", "compose:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.invalidReceipientFormat"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrNoRecipients returns "compose:notification.noRecipients" as *errors.Error
//
//
// This function is auto-generated.
//
func NotificationErrNoRecipients(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("cannot send email message without recipients", nil),

		errors.Meta("type", "noRecipients"),
		errors.Meta("resource", "compose:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.noRecipients"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// NotificationErrFailedToDownloadAttachment returns "compose:notification.failedToDownloadAttachment" as *errors.Error
//
//
// This function is auto-generated.
//
func NotificationErrFailedToDownloadAttachment(mm ...*notificationActionProps) *errors.Error {
	var p = &notificationActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("could not download attachment from {{attachmentURL}}: {{err}}", nil),

		errors.Meta("type", "failedToDownloadAttachment"),
		errors.Meta("resource", "compose:notification"),

		errors.Meta(notificationPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "notification.errors.failedToDownloadAttachment"),

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
func (svc notification) recordAction(ctx context.Context, props *notificationActionProps, actionFn func(...*notificationActionProps) *notificationAction, err error) error {
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
		a.Description = props.Format(m.AsString(notificationLogMetaKey{}), err)

		if p, has := m[notificationPropsMetaKey{}]; has {
			a.Meta = p.(*notificationActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
