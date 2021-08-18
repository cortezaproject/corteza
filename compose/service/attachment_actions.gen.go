package service

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// compose/service/attachment_actions.yaml

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/locale"
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
		namespace  *types.Namespace
		record     *types.Record
		page       *types.Page
		module     *types.Module
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

// setNamespace updates attachmentActionProps's namespace
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setNamespace(namespace *types.Namespace) *attachmentActionProps {
	p.namespace = namespace
	return p
}

// setRecord updates attachmentActionProps's record
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setRecord(record *types.Record) *attachmentActionProps {
	p.record = record
	return p
}

// setPage updates attachmentActionProps's page
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setPage(page *types.Page) *attachmentActionProps {
	p.page = page
	return p
}

// setModule updates attachmentActionProps's module
//
// Allows method chaining
//
// This function is auto-generated.
//
func (p *attachmentActionProps) setModule(module *types.Module) *attachmentActionProps {
	p.module = module
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
		m.Set("attachment.namespaceID", p.attachment.NamespaceID, true)
	}
	if p.filter != nil {
		m.Set("filter.filter", p.filter.Filter, true)
		m.Set("filter.kind", p.filter.Kind, true)
		m.Set("filter.sort", p.filter.Sort, true)
	}
	if p.namespace != nil {
		m.Set("namespace.name", p.namespace.Name, true)
		m.Set("namespace.slug", p.namespace.Slug, true)
		m.Set("namespace.ID", p.namespace.ID, true)
	}
	if p.record != nil {
		m.Set("record.ID", p.record.ID, true)
		m.Set("record.moduleID", p.record.ModuleID, true)
		m.Set("record.namespaceID", p.record.NamespaceID, true)
	}
	if p.page != nil {
		m.Set("page.handle", p.page.Handle, true)
		m.Set("page.title", p.page.Title, true)
		m.Set("page.ID", p.page.ID, true)
	}
	if p.module != nil {
		m.Set("module.handle", p.module.Handle, true)
		m.Set("module.name", p.module.Name, true)
		m.Set("module.ID", p.module.ID, true)
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
				p.attachment.NamespaceID,
			),
		)
		pairs = append(pairs, "{{attachment.name}}", fns(p.attachment.Name))
		pairs = append(pairs, "{{attachment.kind}}", fns(p.attachment.Kind))
		pairs = append(pairs, "{{attachment.url}}", fns(p.attachment.Url))
		pairs = append(pairs, "{{attachment.previewUrl}}", fns(p.attachment.PreviewUrl))
		pairs = append(pairs, "{{attachment.meta}}", fns(p.attachment.Meta))
		pairs = append(pairs, "{{attachment.ownerID}}", fns(p.attachment.OwnerID))
		pairs = append(pairs, "{{attachment.ID}}", fns(p.attachment.ID))
		pairs = append(pairs, "{{attachment.namespaceID}}", fns(p.attachment.NamespaceID))
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

	if p.namespace != nil {
		// replacement for "{{namespace}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{namespace}}",
			fns(
				p.namespace.Name,
				p.namespace.Slug,
				p.namespace.ID,
			),
		)
		pairs = append(pairs, "{{namespace.name}}", fns(p.namespace.Name))
		pairs = append(pairs, "{{namespace.slug}}", fns(p.namespace.Slug))
		pairs = append(pairs, "{{namespace.ID}}", fns(p.namespace.ID))
	}

	if p.record != nil {
		// replacement for "{{record}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{record}}",
			fns(
				p.record.ID,
				p.record.ModuleID,
				p.record.NamespaceID,
			),
		)
		pairs = append(pairs, "{{record.ID}}", fns(p.record.ID))
		pairs = append(pairs, "{{record.moduleID}}", fns(p.record.ModuleID))
		pairs = append(pairs, "{{record.namespaceID}}", fns(p.record.NamespaceID))
	}

	if p.page != nil {
		// replacement for "{{page}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{page}}",
			fns(
				p.page.Handle,
				p.page.Title,
				p.page.ID,
			),
		)
		pairs = append(pairs, "{{page.handle}}", fns(p.page.Handle))
		pairs = append(pairs, "{{page.title}}", fns(p.page.Title))
		pairs = append(pairs, "{{page.ID}}", fns(p.page.ID))
	}

	if p.module != nil {
		// replacement for "{{module}}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{{module}}",
			fns(
				p.module.Handle,
				p.module.Name,
				p.module.ID,
			),
		)
		pairs = append(pairs, "{{module.handle}}", fns(p.module.Handle))
		pairs = append(pairs, "{{module.name}}", fns(p.module.Name))
		pairs = append(pairs, "{{module.ID}}", fns(p.module.ID))
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

// AttachmentActionSearch returns "compose:attachment.search" action
//
// This function is auto-generated.
//
func AttachmentActionSearch(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		action:    "search",
		log:       "searched for attachments",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionLookup returns "compose:attachment.lookup" action
//
// This function is auto-generated.
//
func AttachmentActionLookup(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		action:    "lookup",
		log:       "looked-up for a {{attachment}}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionCreate returns "compose:attachment.create" action
//
// This function is auto-generated.
//
func AttachmentActionCreate(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		action:    "create",
		log:       "created {{attachment}}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionDelete returns "compose:attachment.delete" action
//
// This function is auto-generated.
//
func AttachmentActionDelete(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "compose:attachment",
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

// AttachmentErrGeneric returns "compose:attachment.generic" as *errors.Error
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
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "{err}"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.generic"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotFound returns "compose:attachment.notFound" as *errors.Error
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
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNamespaceNotFound returns "compose:attachment.namespaceNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNamespaceNotFound(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("namespace not found", nil),

		errors.Meta("type", "namespaceNotFound"),
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.namespaceNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrModuleNotFound returns "compose:attachment.moduleNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrModuleNotFound(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("module not found", nil),

		errors.Meta("type", "moduleNotFound"),
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.moduleNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrPageNotFound returns "compose:attachment.pageNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrPageNotFound(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("page not found", nil),

		errors.Meta("type", "pageNotFound"),
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.pageNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrRecordNotFound returns "compose:attachment.recordNotFound" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrRecordNotFound(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("record not found", nil),

		errors.Meta("type", "recordNotFound"),
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.recordNotFound"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrInvalidID returns "compose:attachment.invalidID" as *errors.Error
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
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.invalidID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrInvalidNamespaceID returns "compose:attachment.invalidNamespaceID" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidNamespaceID(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid namespace ID", nil),

		errors.Meta("type", "invalidNamespaceID"),
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.invalidNamespaceID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrInvalidModuleID returns "compose:attachment.invalidModuleID" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidModuleID(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid module ID", nil),

		errors.Meta("type", "invalidModuleID"),
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.invalidModuleID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrInvalidPageID returns "compose:attachment.invalidPageID" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidPageID(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid page ID", nil),

		errors.Meta("type", "invalidPageID"),
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.invalidPageID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrInvalidRecordID returns "compose:attachment.invalidRecordID" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidRecordID(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("invalid record ID", nil),

		errors.Meta("type", "invalidRecordID"),
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.invalidRecordID"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToListAttachments returns "compose:attachment.notAllowedToListAttachments" as *errors.Error
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
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not list attachments; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToListAttachments"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToCreate returns "compose:attachment.notAllowedToCreate" as *errors.Error
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
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not create attachments; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToCreate"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToCreateEmptyAttachment returns "compose:attachment.notAllowedToCreateEmptyAttachment" as *errors.Error
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
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "failed to create attachment; empty file"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToCreateEmptyAttachment"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrFailedToExtractMimeType returns "compose:attachment.failedToExtractMimeType" as *errors.Error
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
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.failedToExtractMimeType"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrFailedToStoreFile returns "compose:attachment.failedToStoreFile" as *errors.Error
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
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.failedToStoreFile"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrFailedToProcessImage returns "compose:attachment.failedToProcessImage" as *errors.Error
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
		errors.Meta("resource", "compose:attachment"),

		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.failedToProcessImage"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToRead returns "compose:attachment.notAllowedToRead" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToRead(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this module", nil),

		errors.Meta("type", "notAllowedToRead"),
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not delete {{module}}; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToRead"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToSearch returns "compose:attachment.notAllowedToSearch" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToSearch(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to search or list modules", nil),

		errors.Meta("type", "notAllowedToSearch"),
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not search or list modules; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToSearch"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToReadNamespace returns "compose:attachment.notAllowedToReadNamespace" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToReadNamespace(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this namespace", nil),

		errors.Meta("type", "notAllowedToReadNamespace"),
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not delete {{namespace}}; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToReadNamespace"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToReadPage returns "compose:attachment.notAllowedToReadPage" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToReadPage(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this page", nil),

		errors.Meta("type", "notAllowedToReadPage"),
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not read {{page}}; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToReadPage"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToReadRecord returns "compose:attachment.notAllowedToReadRecord" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToReadRecord(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to read this record", nil),

		errors.Meta("type", "notAllowedToReadRecord"),
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not read {{record}}; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToReadRecord"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToUpdatePage returns "compose:attachment.notAllowedToUpdatePage" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToUpdatePage(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this page", nil),

		errors.Meta("type", "notAllowedToUpdatePage"),
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not update {{page}}; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToUpdatePage"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToCreateRecords returns "compose:attachment.notAllowedToCreateRecords" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToCreateRecords(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to create records", nil),

		errors.Meta("type", "notAllowedToCreateRecords"),
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not create records; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToCreateRecords"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToUpdateRecord returns "compose:attachment.notAllowedToUpdateRecord" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToUpdateRecord(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this record", nil),

		errors.Meta("type", "notAllowedToUpdateRecord"),
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not update {{record}}; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToUpdateRecord"),

		errors.StackSkip(1),
	)

	if len(mm) > 0 {
	}

	return e
}

// AttachmentErrNotAllowedToUpdateNamespace returns "compose:attachment.notAllowedToUpdateNamespace" as *errors.Error
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToUpdateNamespace(mm ...*attachmentActionProps) *errors.Error {
	var p = &attachmentActionProps{}
	if len(mm) > 0 {
		p = mm[0]
	}

	var e = errors.New(
		errors.KindInternal,

		p.Format("not allowed to update this namespace", nil),

		errors.Meta("type", "notAllowedToUpdateNamespace"),
		errors.Meta("resource", "compose:attachment"),

		// action log entry; no formatting, it will be applied inside recordAction fn.
		errors.Meta(attachmentLogMetaKey{}, "could not update {{namespace}}; insufficient permissions"),
		errors.Meta(attachmentPropsMetaKey{}, p),

		// translation namespace & key
		errors.Meta(locale.ErrorMetaNamespace{}, "compose"),
		errors.Meta(locale.ErrorMetaKey{}, "attachment.errors.notAllowedToUpdateNamespace"),

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
