package service

// This file is auto-generated from compose/service/attachment_actions.yaml
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

// serialize converts attachmentActionProps to actionlog.Meta
//
// This function is auto-generated.
//
func (p attachmentActionProps) serialize() actionlog.Meta {
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
func (p attachmentActionProps) tr(in string, err error) string {
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
	pairs = append(pairs, "{size}", fns(p.size))
	pairs = append(pairs, "{name}", fns(p.name))
	pairs = append(pairs, "{mimetype}", fns(p.mimetype))
	pairs = append(pairs, "{url}", fns(p.url))

	if p.attachment != nil {
		// replacement for "{attachment}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{attachment}",
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
		pairs = append(pairs, "{attachment.name}", fns(p.attachment.Name))
		pairs = append(pairs, "{attachment.kind}", fns(p.attachment.Kind))
		pairs = append(pairs, "{attachment.url}", fns(p.attachment.Url))
		pairs = append(pairs, "{attachment.previewUrl}", fns(p.attachment.PreviewUrl))
		pairs = append(pairs, "{attachment.meta}", fns(p.attachment.Meta))
		pairs = append(pairs, "{attachment.ownerID}", fns(p.attachment.OwnerID))
		pairs = append(pairs, "{attachment.ID}", fns(p.attachment.ID))
		pairs = append(pairs, "{attachment.namespaceID}", fns(p.attachment.NamespaceID))
	}

	if p.filter != nil {
		// replacement for "{filter}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{filter}",
			fns(
				p.filter.Filter,
				p.filter.Kind,
				p.filter.Sort,
			),
		)
		pairs = append(pairs, "{filter.filter}", fns(p.filter.Filter))
		pairs = append(pairs, "{filter.kind}", fns(p.filter.Kind))
		pairs = append(pairs, "{filter.sort}", fns(p.filter.Sort))
	}

	if p.namespace != nil {
		// replacement for "{namespace}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{namespace}",
			fns(
				p.namespace.Name,
				p.namespace.Slug,
				p.namespace.ID,
			),
		)
		pairs = append(pairs, "{namespace.name}", fns(p.namespace.Name))
		pairs = append(pairs, "{namespace.slug}", fns(p.namespace.Slug))
		pairs = append(pairs, "{namespace.ID}", fns(p.namespace.ID))
	}

	if p.record != nil {
		// replacement for "{record}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{record}",
			fns(
				p.record.ID,
				p.record.ModuleID,
				p.record.NamespaceID,
			),
		)
		pairs = append(pairs, "{record.ID}", fns(p.record.ID))
		pairs = append(pairs, "{record.moduleID}", fns(p.record.ModuleID))
		pairs = append(pairs, "{record.namespaceID}", fns(p.record.NamespaceID))
	}

	if p.page != nil {
		// replacement for "{page}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{page}",
			fns(
				p.page.Handle,
				p.page.Title,
				p.page.ID,
			),
		)
		pairs = append(pairs, "{page.handle}", fns(p.page.Handle))
		pairs = append(pairs, "{page.title}", fns(p.page.Title))
		pairs = append(pairs, "{page.ID}", fns(p.page.ID))
	}

	if p.module != nil {
		// replacement for "{module}" (in order how fields are defined)
		pairs = append(
			pairs,
			"{module}",
			fns(
				p.module.Handle,
				p.module.Name,
				p.module.ID,
			),
		)
		pairs = append(pairs, "{module.handle}", fns(p.module.Handle))
		pairs = append(pairs, "{module.name}", fns(p.module.Name))
		pairs = append(pairs, "{module.ID}", fns(p.module.ID))
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

// AttachmentActionSearch returns "compose:attachment.search" error
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

// AttachmentActionLookup returns "compose:attachment.lookup" error
//
// This function is auto-generated.
//
func AttachmentActionLookup(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		action:    "lookup",
		log:       "looked-up for a {attachment}",
		severity:  actionlog.Info,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionCreate returns "compose:attachment.create" error
//
// This function is auto-generated.
//
func AttachmentActionCreate(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		action:    "create",
		log:       "created {attachment}",
		severity:  actionlog.Notice,
	}

	if len(props) > 0 {
		a.props = props[0]
	}

	return a
}

// AttachmentActionDelete returns "compose:attachment.delete" error
//
// This function is auto-generated.
//
func AttachmentActionDelete(props ...*attachmentActionProps) *attachmentAction {
	a := &attachmentAction{
		timestamp: time.Now(),
		resource:  "compose:attachment",
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

// AttachmentErrGeneric returns "compose:attachment.generic" audit event as actionlog.Error
//
//
// This function is auto-generated.
//
func AttachmentErrGeneric(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
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

// AttachmentErrNotFound returns "compose:attachment.notFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrNotFound(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notFound",
		action:    "error",
		message:   "attachment not found",
		log:       "attachment not found",
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

// AttachmentErrNamespaceNotFound returns "compose:attachment.namespaceNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrNamespaceNotFound(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "namespaceNotFound",
		action:    "error",
		message:   "namespace not found",
		log:       "namespace not found",
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

// AttachmentErrModuleNotFound returns "compose:attachment.moduleNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrModuleNotFound(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "moduleNotFound",
		action:    "error",
		message:   "module not found",
		log:       "module not found",
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

// AttachmentErrPageNotFound returns "compose:attachment.pageNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrPageNotFound(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "pageNotFound",
		action:    "error",
		message:   "page not found",
		log:       "page not found",
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

// AttachmentErrRecordNotFound returns "compose:attachment.recordNotFound" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrRecordNotFound(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "recordNotFound",
		action:    "error",
		message:   "record not found",
		log:       "record not found",
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

// AttachmentErrInvalidID returns "compose:attachment.invalidID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidID(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
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

// AttachmentErrInvalidNamespaceID returns "compose:attachment.invalidNamespaceID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidNamespaceID(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "invalidNamespaceID",
		action:    "error",
		message:   "invalid namespace ID",
		log:       "invalid namespace ID",
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

// AttachmentErrInvalidModuleID returns "compose:attachment.invalidModuleID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidModuleID(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "invalidModuleID",
		action:    "error",
		message:   "invalid module ID",
		log:       "invalid module ID",
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

// AttachmentErrInvalidPageID returns "compose:attachment.invalidPageID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidPageID(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "invalidPageID",
		action:    "error",
		message:   "invalid page ID",
		log:       "invalid page ID",
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

// AttachmentErrInvalidRecordID returns "compose:attachment.invalidRecordID" audit event as actionlog.Warning
//
//
// This function is auto-generated.
//
func AttachmentErrInvalidRecordID(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "invalidRecordID",
		action:    "error",
		message:   "invalid record ID",
		log:       "invalid record ID",
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

// AttachmentErrNotAllowedToListAttachments returns "compose:attachment.notAllowedToListAttachments" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToListAttachments(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notAllowedToListAttachments",
		action:    "error",
		message:   "not allowed to list attachments",
		log:       "could not list attachments; insufficient permissions",
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

// AttachmentErrNotAllowedToCreate returns "compose:attachment.notAllowedToCreate" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToCreate(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notAllowedToCreate",
		action:    "error",
		message:   "not allowed to create attachments",
		log:       "could not create attachments; insufficient permissions",
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

// AttachmentErrFailedToExtractMimeType returns "compose:attachment.failedToExtractMimeType" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrFailedToExtractMimeType(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
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

// AttachmentErrFailedToStoreFile returns "compose:attachment.failedToStoreFile" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrFailedToStoreFile(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
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

// AttachmentErrFailedToProcessImage returns "compose:attachment.failedToProcessImage" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrFailedToProcessImage(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
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

// AttachmentErrNotAllowedToReadModule returns "compose:attachment.notAllowedToReadModule" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToReadModule(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notAllowedToReadModule",
		action:    "error",
		message:   "not allowed to read this module",
		log:       "could not delete {module}; insufficient permissions",
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

// AttachmentErrNotAllowedToReadNamespace returns "compose:attachment.notAllowedToReadNamespace" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToReadNamespace(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notAllowedToReadNamespace",
		action:    "error",
		message:   "not allowed to read this namespace",
		log:       "could not delete {namespace}; insufficient permissions",
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

// AttachmentErrNotAllowedToReadPage returns "compose:attachment.notAllowedToReadPage" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToReadPage(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notAllowedToReadPage",
		action:    "error",
		message:   "not allowed to read this page",
		log:       "could not read {page}; insufficient permissions",
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

// AttachmentErrNotAllowedToReadRecord returns "compose:attachment.notAllowedToReadRecord" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToReadRecord(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notAllowedToReadRecord",
		action:    "error",
		message:   "not allowed to read this record",
		log:       "could not read {record}; insufficient permissions",
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

// AttachmentErrNotAllowedToUpdatePage returns "compose:attachment.notAllowedToUpdatePage" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToUpdatePage(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notAllowedToUpdatePage",
		action:    "error",
		message:   "not allowed to update this page",
		log:       "could not update {page}; insufficient permissions",
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

// AttachmentErrNotAllowedToCreateRecords returns "compose:attachment.notAllowedToCreateRecords" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToCreateRecords(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notAllowedToCreateRecords",
		action:    "error",
		message:   "not allowed to create records",
		log:       "could not create records; insufficient permissions",
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

// AttachmentErrNotAllowedToUpdateRecord returns "compose:attachment.notAllowedToUpdateRecord" audit event as actionlog.Alert
//
//
// This function is auto-generated.
//
func AttachmentErrNotAllowedToUpdateRecord(props ...*attachmentActionProps) *attachmentError {
	var e = &attachmentError{
		timestamp: time.Now(),
		resource:  "compose:attachment",
		error:     "notAllowedToUpdateRecord",
		action:    "error",
		message:   "not allowed to update this record",
		log:       "could not update {record}; insufficient permissions",
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
