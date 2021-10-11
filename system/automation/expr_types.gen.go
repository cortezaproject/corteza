package automation

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// system/automation/expr_types.yaml

import (
	"context"
	"fmt"
	. "github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"
	"sync"
)

var _ = context.Background
var _ = fmt.Errorf

// DocumentType is an expression type, wrapper for types.DocumentType type
type DocumentType struct {
	value types.DocumentType
	mux   sync.RWMutex
}

// NewDocumentType creates new instance of DocumentType expression type
func NewDocumentType(val interface{}) (*DocumentType, error) {
	if c, err := CastToDocumentType(val); err != nil {
		return nil, fmt.Errorf("unable to create DocumentType: %w", err)
	} else {
		return &DocumentType{value: c}, nil
	}
}

// Get return underlying value on DocumentType
func (t *DocumentType) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on DocumentType
func (t *DocumentType) GetValue() types.DocumentType {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (DocumentType) Type() string { return "DocumentType" }

// Cast converts value to types.DocumentType
func (DocumentType) Cast(val interface{}) (TypedValue, error) {
	return NewDocumentType(val)
}

// Assign new value to DocumentType
//
// value is first passed through CastToDocumentType
func (t *DocumentType) Assign(val interface{}) error {
	if c, err := CastToDocumentType(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// QueueMessage is an expression type, wrapper for *types.QueueMessage type
type QueueMessage struct {
	value *types.QueueMessage
	mux   sync.RWMutex
}

// NewQueueMessage creates new instance of QueueMessage expression type
func NewQueueMessage(val interface{}) (*QueueMessage, error) {
	if c, err := CastToQueueMessage(val); err != nil {
		return nil, fmt.Errorf("unable to create QueueMessage: %w", err)
	} else {
		return &QueueMessage{value: c}, nil
	}
}

// Get return underlying value on QueueMessage
func (t *QueueMessage) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on QueueMessage
func (t *QueueMessage) GetValue() *types.QueueMessage {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (QueueMessage) Type() string { return "QueueMessage" }

// Cast converts value to *types.QueueMessage
func (QueueMessage) Cast(val interface{}) (TypedValue, error) {
	return NewQueueMessage(val)
}

// Assign new value to QueueMessage
//
// value is first passed through CastToQueueMessage
func (t *QueueMessage) Assign(val interface{}) error {
	if c, err := CastToQueueMessage(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *QueueMessage) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return assignToQueueMessage(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access QueueMessage's underlying value (*types.QueueMessage)
// and it's fields
//
func (t *QueueMessage) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return queueMessageGValSelector(t.value, k)
}

// Select is field accessor for *types.QueueMessage
//
// Similar to SelectGVal but returns typed values
func (t *QueueMessage) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return queueMessageTypedValueSelector(t.value, k)
}

func (t *QueueMessage) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch k {
	case "Queue":
		return true
	case "Payload":
		return true
	}
	return false
}

// queueMessageGValSelector is field accessor for *types.QueueMessage
func queueMessageGValSelector(res *types.QueueMessage, k string) (interface{}, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "Queue":
		return res.Queue, nil
	case "Payload":
		return res.Payload, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// queueMessageTypedValueSelector is field accessor for *types.QueueMessage
func queueMessageTypedValueSelector(res *types.QueueMessage, k string) (TypedValue, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "Queue":
		return NewString(res.Queue)
	case "Payload":
		return NewBytes(res.Payload)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToQueueMessage is field value setter for *types.QueueMessage
func assignToQueueMessage(res *types.QueueMessage, k string, val interface{}) error {
	switch k {
	case "Queue":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Queue = aux
		return nil
	case "Payload":
		aux, err := CastToBytes(val)
		if err != nil {
			return err
		}

		res.Payload = aux
		return nil
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// RbacResource is an expression type, wrapper for rbac.Resource type
type RbacResource struct {
	value rbac.Resource
	mux   sync.RWMutex
}

// NewRbacResource creates new instance of RbacResource expression type
func NewRbacResource(val interface{}) (*RbacResource, error) {
	if c, err := CastToRbacResource(val); err != nil {
		return nil, fmt.Errorf("unable to create RbacResource: %w", err)
	} else {
		return &RbacResource{value: c}, nil
	}
}

// Get return underlying value on RbacResource
func (t *RbacResource) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on RbacResource
func (t *RbacResource) GetValue() rbac.Resource {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (RbacResource) Type() string { return "RbacResource" }

// Cast converts value to rbac.Resource
func (RbacResource) Cast(val interface{}) (TypedValue, error) {
	return NewRbacResource(val)
}

// Assign new value to RbacResource
//
// value is first passed through CastToRbacResource
func (t *RbacResource) Assign(val interface{}) error {
	if c, err := CastToRbacResource(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// RenderOptions is an expression type, wrapper for map[string]string type
type RenderOptions struct {
	value map[string]string
	mux   sync.RWMutex
}

// NewRenderOptions creates new instance of RenderOptions expression type
func NewRenderOptions(val interface{}) (*RenderOptions, error) {
	if c, err := CastToRenderOptions(val); err != nil {
		return nil, fmt.Errorf("unable to create RenderOptions: %w", err)
	} else {
		return &RenderOptions{value: c}, nil
	}
}

// Get return underlying value on RenderOptions
func (t *RenderOptions) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on RenderOptions
func (t *RenderOptions) GetValue() map[string]string {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (RenderOptions) Type() string { return "RenderOptions" }

// Cast converts value to map[string]string
func (RenderOptions) Cast(val interface{}) (TypedValue, error) {
	return NewRenderOptions(val)
}

// Assign new value to RenderOptions
//
// value is first passed through CastToRenderOptions
func (t *RenderOptions) Assign(val interface{}) error {
	if c, err := CastToRenderOptions(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

// RenderedDocument is an expression type, wrapper for *renderedDocument type
type RenderedDocument struct {
	value *renderedDocument
	mux   sync.RWMutex
}

// NewRenderedDocument creates new instance of RenderedDocument expression type
func NewRenderedDocument(val interface{}) (*RenderedDocument, error) {
	if c, err := CastToRenderedDocument(val); err != nil {
		return nil, fmt.Errorf("unable to create RenderedDocument: %w", err)
	} else {
		return &RenderedDocument{value: c}, nil
	}
}

// Get return underlying value on RenderedDocument
func (t *RenderedDocument) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on RenderedDocument
func (t *RenderedDocument) GetValue() *renderedDocument {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (RenderedDocument) Type() string { return "RenderedDocument" }

// Cast converts value to *renderedDocument
func (RenderedDocument) Cast(val interface{}) (TypedValue, error) {
	return NewRenderedDocument(val)
}

// Assign new value to RenderedDocument
//
// value is first passed through CastToRenderedDocument
func (t *RenderedDocument) Assign(val interface{}) error {
	if c, err := CastToRenderedDocument(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *RenderedDocument) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return assignToRenderedDocument(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access RenderedDocument's underlying value (*renderedDocument)
// and it's fields
//
func (t *RenderedDocument) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return renderedDocumentGValSelector(t.value, k)
}

// Select is field accessor for *renderedDocument
//
// Similar to SelectGVal but returns typed values
func (t *RenderedDocument) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return renderedDocumentTypedValueSelector(t.value, k)
}

func (t *RenderedDocument) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch k {
	case "document":
		return true
	case "name":
		return true
	case "type":
		return true
	}
	return false
}

// renderedDocumentGValSelector is field accessor for *renderedDocument
func renderedDocumentGValSelector(res *renderedDocument, k string) (interface{}, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "document":
		return res.Document, nil
	case "name":
		return res.Name, nil
	case "type":
		return res.Type, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// renderedDocumentTypedValueSelector is field accessor for *renderedDocument
func renderedDocumentTypedValueSelector(res *renderedDocument, k string) (TypedValue, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "document":
		return NewReader(res.Document)
	case "name":
		return NewString(res.Name)
	case "type":
		return NewString(res.Type)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToRenderedDocument is field value setter for *renderedDocument
func assignToRenderedDocument(res *renderedDocument, k string, val interface{}) error {
	switch k {
	case "document":
		aux, err := CastToReader(val)
		if err != nil {
			return err
		}

		res.Document = aux
		return nil
	case "name":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Name = aux
		return nil
	case "type":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Type = aux
		return nil
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// Role is an expression type, wrapper for *types.Role type
type Role struct {
	value *types.Role
	mux   sync.RWMutex
}

// NewRole creates new instance of Role expression type
func NewRole(val interface{}) (*Role, error) {
	if c, err := CastToRole(val); err != nil {
		return nil, fmt.Errorf("unable to create Role: %w", err)
	} else {
		return &Role{value: c}, nil
	}
}

// Get return underlying value on Role
func (t *Role) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Role
func (t *Role) GetValue() *types.Role {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Role) Type() string { return "Role" }

// Cast converts value to *types.Role
func (Role) Cast(val interface{}) (TypedValue, error) {
	return NewRole(val)
}

// Assign new value to Role
//
// value is first passed through CastToRole
func (t *Role) Assign(val interface{}) error {
	if c, err := CastToRole(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *Role) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return assignToRole(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Role's underlying value (*types.Role)
// and it's fields
//
func (t *Role) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return roleGValSelector(t.value, k)
}

// Select is field accessor for *types.Role
//
// Similar to SelectGVal but returns typed values
func (t *Role) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return roleTypedValueSelector(t.value, k)
}

func (t *Role) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch k {
	case "ID", "roleID":
		return true
	case "name":
		return true
	case "handle":
		return true
	case "labels":
		return true
	case "createdAt":
		return true
	case "updatedAt":
		return true
	case "archivedAt":
		return true
	case "deletedAt":
		return true
	}
	return false
}

// roleGValSelector is field accessor for *types.Role
func roleGValSelector(res *types.Role, k string) (interface{}, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "ID", "roleID":
		return res.ID, nil
	case "name":
		return res.Name, nil
	case "handle":
		return res.Handle, nil
	case "labels":
		return res.Labels, nil
	case "createdAt":
		return res.CreatedAt, nil
	case "updatedAt":
		return res.UpdatedAt, nil
	case "archivedAt":
		return res.ArchivedAt, nil
	case "deletedAt":
		return res.DeletedAt, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// roleTypedValueSelector is field accessor for *types.Role
func roleTypedValueSelector(res *types.Role, k string) (TypedValue, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "ID", "roleID":
		return NewID(res.ID)
	case "name":
		return NewString(res.Name)
	case "handle":
		return NewHandle(res.Handle)
	case "labels":
		return NewKV(res.Labels)
	case "createdAt":
		return NewDateTime(res.CreatedAt)
	case "updatedAt":
		return NewDateTime(res.UpdatedAt)
	case "archivedAt":
		return NewDateTime(res.ArchivedAt)
	case "deletedAt":
		return NewDateTime(res.DeletedAt)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToRole is field value setter for *types.Role
func assignToRole(res *types.Role, k string, val interface{}) error {
	switch k {
	case "ID", "roleID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "name":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Name = aux
		return nil
	case "handle":
		aux, err := CastToHandle(val)
		if err != nil {
			return err
		}

		res.Handle = aux
		return nil
	case "labels":
		aux, err := CastToKV(val)
		if err != nil {
			return err
		}

		res.Labels = aux
		return nil
	case "createdAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "updatedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "archivedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "deletedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// Template is an expression type, wrapper for *types.Template type
type Template struct {
	value *types.Template
	mux   sync.RWMutex
}

// NewTemplate creates new instance of Template expression type
func NewTemplate(val interface{}) (*Template, error) {
	if c, err := CastToTemplate(val); err != nil {
		return nil, fmt.Errorf("unable to create Template: %w", err)
	} else {
		return &Template{value: c}, nil
	}
}

// Get return underlying value on Template
func (t *Template) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on Template
func (t *Template) GetValue() *types.Template {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (Template) Type() string { return "Template" }

// Cast converts value to *types.Template
func (Template) Cast(val interface{}) (TypedValue, error) {
	return NewTemplate(val)
}

// Assign new value to Template
//
// value is first passed through CastToTemplate
func (t *Template) Assign(val interface{}) error {
	if c, err := CastToTemplate(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *Template) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return assignToTemplate(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access Template's underlying value (*types.Template)
// and it's fields
//
func (t *Template) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return templateGValSelector(t.value, k)
}

// Select is field accessor for *types.Template
//
// Similar to SelectGVal but returns typed values
func (t *Template) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return templateTypedValueSelector(t.value, k)
}

func (t *Template) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch k {
	case "ID", "templateID":
		return true
	case "handle":
		return true
	case "language":
		return true
	case "type":
		return true
	case "partial":
		return true
	case "meta":
		return true
	case "template":
		return true
	case "labels":
		return true
	case "ownerID":
		return true
	case "createdAt":
		return true
	case "updatedAt":
		return true
	case "deletedAt":
		return true
	case "lastUsedAt":
		return true
	}
	return false
}

// templateGValSelector is field accessor for *types.Template
func templateGValSelector(res *types.Template, k string) (interface{}, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "ID", "templateID":
		return res.ID, nil
	case "handle":
		return res.Handle, nil
	case "language":
		return res.Language, nil
	case "type":
		return res.Type, nil
	case "partial":
		return res.Partial, nil
	case "meta":
		return res.Meta, nil
	case "template":
		return res.Template, nil
	case "labels":
		return res.Labels, nil
	case "ownerID":
		return res.OwnerID, nil
	case "createdAt":
		return res.CreatedAt, nil
	case "updatedAt":
		return res.UpdatedAt, nil
	case "deletedAt":
		return res.DeletedAt, nil
	case "lastUsedAt":
		return res.LastUsedAt, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// templateTypedValueSelector is field accessor for *types.Template
func templateTypedValueSelector(res *types.Template, k string) (TypedValue, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "ID", "templateID":
		return NewID(res.ID)
	case "handle":
		return NewHandle(res.Handle)
	case "language":
		return NewString(res.Language)
	case "type":
		return NewDocumentType(res.Type)
	case "partial":
		return NewBoolean(res.Partial)
	case "meta":
		return NewTemplateMeta(res.Meta)
	case "template":
		return NewString(res.Template)
	case "labels":
		return NewKV(res.Labels)
	case "ownerID":
		return NewID(res.OwnerID)
	case "createdAt":
		return NewDateTime(res.CreatedAt)
	case "updatedAt":
		return NewDateTime(res.UpdatedAt)
	case "deletedAt":
		return NewDateTime(res.DeletedAt)
	case "lastUsedAt":
		return NewDateTime(res.LastUsedAt)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToTemplate is field value setter for *types.Template
func assignToTemplate(res *types.Template, k string, val interface{}) error {
	switch k {
	case "ID", "templateID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "handle":
		aux, err := CastToHandle(val)
		if err != nil {
			return err
		}

		res.Handle = aux
		return nil
	case "language":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Language = aux
		return nil
	case "type":
		aux, err := CastToDocumentType(val)
		if err != nil {
			return err
		}

		res.Type = aux
		return nil
	case "partial":
		aux, err := CastToBoolean(val)
		if err != nil {
			return err
		}

		res.Partial = aux
		return nil
	case "meta":
		aux, err := CastToTemplateMeta(val)
		if err != nil {
			return err
		}

		res.Meta = aux
		return nil
	case "template":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Template = aux
		return nil
	case "labels":
		aux, err := CastToKV(val)
		if err != nil {
			return err
		}

		res.Labels = aux
		return nil
	case "ownerID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "createdAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "updatedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "deletedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "lastUsedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// TemplateMeta is an expression type, wrapper for types.TemplateMeta type
type TemplateMeta struct {
	value types.TemplateMeta
	mux   sync.RWMutex
}

// NewTemplateMeta creates new instance of TemplateMeta expression type
func NewTemplateMeta(val interface{}) (*TemplateMeta, error) {
	if c, err := CastToTemplateMeta(val); err != nil {
		return nil, fmt.Errorf("unable to create TemplateMeta: %w", err)
	} else {
		return &TemplateMeta{value: c}, nil
	}
}

// Get return underlying value on TemplateMeta
func (t *TemplateMeta) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on TemplateMeta
func (t *TemplateMeta) GetValue() types.TemplateMeta {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (TemplateMeta) Type() string { return "TemplateMeta" }

// Cast converts value to types.TemplateMeta
func (TemplateMeta) Cast(val interface{}) (TypedValue, error) {
	return NewTemplateMeta(val)
}

// Assign new value to TemplateMeta
//
// value is first passed through CastToTemplateMeta
func (t *TemplateMeta) Assign(val interface{}) error {
	if c, err := CastToTemplateMeta(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *TemplateMeta) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return assignToTemplateMeta(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access TemplateMeta's underlying value (types.TemplateMeta)
// and it's fields
//
func (t *TemplateMeta) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return templateMetaGValSelector(t.value, k)
}

// Select is field accessor for types.TemplateMeta
//
// Similar to SelectGVal but returns typed values
func (t *TemplateMeta) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return templateMetaTypedValueSelector(t.value, k)
}

func (t *TemplateMeta) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch k {
	case "short":
		return true
	case "description":
		return true
	}
	return false
}

// templateMetaGValSelector is field accessor for types.TemplateMeta
func templateMetaGValSelector(res types.TemplateMeta, k string) (interface{}, error) {
	switch k {
	case "short":
		return res.Short, nil
	case "description":
		return res.Description, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// templateMetaTypedValueSelector is field accessor for types.TemplateMeta
func templateMetaTypedValueSelector(res types.TemplateMeta, k string) (TypedValue, error) {
	switch k {
	case "short":
		return NewString(res.Short)
	case "description":
		return NewString(res.Description)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToTemplateMeta is field value setter for types.TemplateMeta
func assignToTemplateMeta(res types.TemplateMeta, k string, val interface{}) error {
	switch k {
	case "short":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Short = aux
		return nil
	case "description":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Description = aux
		return nil
	}

	return fmt.Errorf("unknown field '%s'", k)
}

// User is an expression type, wrapper for *types.User type
type User struct {
	value *types.User
	mux   sync.RWMutex
}

// NewUser creates new instance of User expression type
func NewUser(val interface{}) (*User, error) {
	if c, err := CastToUser(val); err != nil {
		return nil, fmt.Errorf("unable to create User: %w", err)
	} else {
		return &User{value: c}, nil
	}
}

// Get return underlying value on User
func (t *User) Get() interface{} {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// GetValue returns underlying value on User
func (t *User) GetValue() *types.User {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.value
}

// Type return type name
func (User) Type() string { return "User" }

// Cast converts value to *types.User
func (User) Cast(val interface{}) (TypedValue, error) {
	return NewUser(val)
}

// Assign new value to User
//
// value is first passed through CastToUser
func (t *User) Assign(val interface{}) error {
	if c, err := CastToUser(val); err != nil {
		return err
	} else {
		t.value = c
		return nil
	}
}

func (t *User) AssignFieldValue(key string, val TypedValue) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	return assignToUser(t.value, key, val)
}

// SelectGVal implements gval.Selector requirements
//
// It allows gval lib to access User's underlying value (*types.User)
// and it's fields
//
func (t *User) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return userGValSelector(t.value, k)
}

// Select is field accessor for *types.User
//
// Similar to SelectGVal but returns typed values
func (t *User) Select(k string) (TypedValue, error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return userTypedValueSelector(t.value, k)
}

func (t *User) Has(k string) bool {
	t.mux.RLock()
	defer t.mux.RUnlock()
	switch k {
	case "ID", "userID":
		return true
	case "username":
		return true
	case "email":
		return true
	case "name":
		return true
	case "handle":
		return true
	case "emailConfirmed":
		return true
	case "labels":
		return true
	case "createdAt":
		return true
	case "updatedAt":
		return true
	case "suspendedAt":
		return true
	case "deletedAt":
		return true
	}
	return false
}

// userGValSelector is field accessor for *types.User
func userGValSelector(res *types.User, k string) (interface{}, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "ID", "userID":
		return res.ID, nil
	case "username":
		return res.Username, nil
	case "email":
		return res.Email, nil
	case "name":
		return res.Name, nil
	case "handle":
		return res.Handle, nil
	case "emailConfirmed":
		return res.EmailConfirmed, nil
	case "labels":
		return res.Labels, nil
	case "createdAt":
		return res.CreatedAt, nil
	case "updatedAt":
		return res.UpdatedAt, nil
	case "suspendedAt":
		return res.SuspendedAt, nil
	case "deletedAt":
		return res.DeletedAt, nil
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// userTypedValueSelector is field accessor for *types.User
func userTypedValueSelector(res *types.User, k string) (TypedValue, error) {
	if res == nil {
		return nil, nil
	}
	switch k {
	case "ID", "userID":
		return NewID(res.ID)
	case "username":
		return NewString(res.Username)
	case "email":
		return NewString(res.Email)
	case "name":
		return NewString(res.Name)
	case "handle":
		return NewHandle(res.Handle)
	case "emailConfirmed":
		return NewBoolean(res.EmailConfirmed)
	case "labels":
		return NewKV(res.Labels)
	case "createdAt":
		return NewDateTime(res.CreatedAt)
	case "updatedAt":
		return NewDateTime(res.UpdatedAt)
	case "suspendedAt":
		return NewDateTime(res.SuspendedAt)
	case "deletedAt":
		return NewDateTime(res.DeletedAt)
	}

	return nil, fmt.Errorf("unknown field '%s'", k)
}

// assignToUser is field value setter for *types.User
func assignToUser(res *types.User, k string, val interface{}) error {
	switch k {
	case "ID", "userID":
		return fmt.Errorf("field '%s' is read-only", k)
	case "username":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Username = aux
		return nil
	case "email":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Email = aux
		return nil
	case "name":
		aux, err := CastToString(val)
		if err != nil {
			return err
		}

		res.Name = aux
		return nil
	case "handle":
		aux, err := CastToHandle(val)
		if err != nil {
			return err
		}

		res.Handle = aux
		return nil
	case "emailConfirmed":
		aux, err := CastToBoolean(val)
		if err != nil {
			return err
		}

		res.EmailConfirmed = aux
		return nil
	case "labels":
		aux, err := CastToKV(val)
		if err != nil {
			return err
		}

		res.Labels = aux
		return nil
	case "createdAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "updatedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "suspendedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	case "deletedAt":
		return fmt.Errorf("field '%s' is read-only", k)
	}

	return fmt.Errorf("unknown field '%s'", k)
}
