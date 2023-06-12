package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/service/attachment_actions.yaml

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

	attachmentLogMetaKey   struct{}
	attachmentPropsMetaKey struct{}
)

var (
	// just a placeholder to cover template cases w/o fmt package use
	_ = fmt.Println
)

// *********************************************************************************************************************
// *********************************************************************************************************************
// Props methods
// setSize updates attachmentActionProps's size
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setSize(size int64) *attachmentActionProps {
	p.size = size
	return p
}

// setName updates attachmentActionProps's name
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setName(name string) *attachmentActionProps {
	p.name = name
	return p
}

// setMimetype updates attachmentActionProps's mimetype
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setMimetype(mimetype string) *attachmentActionProps {
	p.mimetype = mimetype
	return p
}

// setUrl updates attachmentActionProps's url
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setUrl(url string) *attachmentActionProps {
	p.url = url
	return p
}

// setAttachment updates attachmentActionProps's attachment
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setAttachment(attachment *types.Attachment) *attachmentActionProps {
	p.attachment = attachment
	return p
}

// setFilter updates attachmentActionProps's filter
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setFilter(filter *types.AttachmentFilter) *attachmentActionProps {
	p.filter = filter
	return p
}

// Serialize converts attachmentActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p attachmentActionProps) Serialize() actionlog.Meta {
	var (
		m = make(actionlog.Meta)
	)

	m.Set("size", p.size, true)
	m.Set("name", p.name, true)
	m.Set("mimetype", p.mimetype, true)
	m.Set("url", p.url, true)
	if p.attachment != nil {
		m.Set("attachment.name", p.attachment.Name, true)
		m.Set("attachment.kind", p.attachment.Kind, true)
		m.Set("attachment.url", p.attachment.Url, true)
		m.Set("attachment.previewUrl", p.attachment.PreviewUrl, true)
		m.Set("attachment.meta", p.attachment.Meta, true)
		m.Set("attachment.ownerID", p.attachment.OwnerID, true)
		m.Set("attachment.ID", p.attachment.ID, true)
	}
	if p.filter != nil {
		m.Set("filter.filter", p.filter.Filter, true)
		m.Set("filter.kind", p.filter.Kind, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}

	return m
}

// tr translates string and replaces meta value placeholder with values
//
// This function is auto-generated.
//
func (p attachmentActionProps) Format(in string, err error) string {
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
	pairs = append(pairs, "{{size}}", fns(p.size))
	pairs = append(pairs, "{{name}}", fns(p.name))
	pairs = append(pairs, "{{mimetype}}", fns(p.mimetype))
	pairs = append(pairs, "{{url}}", fns(p.url))

	if p.attachment != nil {
		// replacement for "{{attachment}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{attachment}}",
			fns(
				p.attachment.Name,
				p.attachment.Kind,
				p.attachment.Url,
				p.attachment.PreviewUrl,
				p.attachment.Meta,
				p.attachment.OwnerID,
				p.attachment.ID,
			),
		)
		pairs = append(pairs, "{{attachment.name}}", fns(p.attachment.Name))
		pairs = append(pairs, "{{attachment.kind}}", fns(p.attachment.Kind))
		pairs = append(pairs, "{{attachment.url}}", fns(p.attachment.Url))
		pairs = append(pairs, "{{attachment.previewUrl}}", fns(p.attachment.PreviewUrl))
		pairs = append(pairs, "{{attachment.meta}}", fns(p.attachment.Meta))
		pairs = append(pairs, "{{attachment.ownerID}}", fns(p.attachment.OwnerID))
		pairs = append(pairs, "{{attachment.ID}}", fns(p.attachment.ID))
	}

	if p.filter != nil {
		// replacement for "{{filter}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{filter}}",
			fns(
				p.filter.Filter,
				p.filter.Kind,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{{filter.filter}}", fns(p.filter.Filter))
		pairs = append(pairs, "{{filter.kind}}", fns(p.filter.Kind))
		pairs = append(pairs, "{{filter.sort}}", fns(p.filter.Sort))
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

	return props.Format(a.log, nil)
}

func (e *attachmentAction) ToAction() *actionlog.Action {
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

// AttachmentActionSearch returns "system:attachment.search" action
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

// AttachmentActionLookup returns "system:attachment.lookup" action
//
// This function is auto-generated.
//
func AttachmentActionLookup(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "system:attachment",
		action:    "lookup",
		log:       "looked-up for a {{attachment}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionCreate returns "system:attachment.create" action
//
// This function is auto-generated.
//
func AttachmentActionCreate(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "system:attachment",
		action:    "create",
		log:       "created {{attachment}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionDelete returns "system:attachment.delete" action
//
// This function is auto-generated.
//
func AttachmentActionDelete(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "system:attachment",
		action:    "delete",
		log:       "deleted {{attachment}}",
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

// AttachmentErrGeneric returns "system:attachment.generic" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrGeneric(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to complete request due to internal error", nil),

		errors.Meta("type", "generic"),
		errors.Meta("resource", "system:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "{err}"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotFound returns "system:attachment.notFound" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotFound(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("attachment not found", nil),

		errors.Meta("type", "notFound"),
		errors.Meta("resource", "system:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrInvalidID returns "system:attachment.invalidID" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidID(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid ID", nil),

		errors.Meta("type", "invalidID"),
		errors.Meta("resource", "system:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToListAttachments returns "system:attachment.notAllowedToListAttachments" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToListAttachments(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to list attachments", nil),

		errors.Meta("type", "notAllowedToListAttachments"),
		errors.Meta("resource", "system:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "failed to list attachment; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToListAttachments"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToCreate returns "system:attachment.notAllowedToCreate" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToCreate(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create attachments", nil),

		errors.Meta("type", "notAllowedToCreate"),
		errors.Meta("resource", "system:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "failed to create attachment; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToCreateEmptyAttachment returns "system:attachment.notAllowedToCreateEmptyAttachment" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToCreateEmptyAttachment(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create empty attachments", nil),

		errors.Meta("type", "notAllowedToCreateEmptyAttachment"),
		errors.Meta("resource", "system:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "failed to create attachment; empty file"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToCreateEmptyAttachment"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrFailedToExtractMimeType returns "system:attachment.failedToExtractMimeType" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrFailedToExtractMimeType(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("could not extract mime type", nil),

		errors.Meta("type", "failedToExtractMimeType"),
		errors.Meta("resource", "system:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.failedToExtractMimeType"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrFailedToStoreFile returns "system:attachment.failedToStoreFile" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrFailedToStoreFile(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("could not extract store file", nil),

		errors.Meta("type", "failedToStoreFile"),
		errors.Meta("resource", "system:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.failedToStoreFile"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrFailedToProcessImage returns "system:attachment.failedToProcessImage" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrFailedToProcessImage(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("could not process image", nil),

		errors.Meta("type", "failedToProcessImage"),
		errors.Meta("resource", "system:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.failedToProcessImage"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrInvalidAvatarFileType returns "system:attachment.invalidAvatarFileType" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidAvatarFileType(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid file type only JPEG and PNG are allowed", nil),

		errors.Meta("type", "invalidAvatarFileType"),
		errors.Meta("resource", "system:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.invalidAvatarFileType"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrInvalidAvatarFileSize returns "system:attachment.invalidAvatarFileSize" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidAvatarFileSize(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("file size is too large", nil),

		errors.Meta("type", "invalidAvatarFileSize"),
		errors.Meta("resource", "system:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.invalidAvatarFileSize"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrInvalidAvatarGenerateFontFile returns "system:attachment.invalidAvatarGenerateFontFile" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidAvatarGenerateFontFile(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("failed to generate avatar initial, font file configuration is invalid", nil),

		errors.Meta("type", "invalidAvatarGenerateFontFile"),
		errors.Meta("resource", "system:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "system"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.invalidAvatarGenerateFontFile"),

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
func (svc attachment) recordAction(ctx context.Context, props *attachmentActionProps, actionFn func(...*attachmentActionProps) *attachmentAction, err error) error {
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
		a.Description = props.Format(m.AsString(attachmentLogMetaKey{}), err)

		if p, has := m[attachmentPropsMetaKey{}]; has {
			a.Meta = p.(*attachmentActionProps).Serialize()
		}

		svc.actionlog.Record(ctx, a)
	default:
		svc.actionlog.Record(ctx, a)
	}

	// Original error is passed on
	return err
}
