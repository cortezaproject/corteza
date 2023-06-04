package dal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/modern-go/reflect2"
)

type (
	// ModelRef is used to retrieve a model from the DAL based on given params
	ModelRef struct {
		ConnectionID uint64

		ResourceID uint64

		ResourceType string
		Resource     string

		Refs map[string]any
	}

	// Model describes the underlying data and its shape
	Model struct {
		ConnectionID uint64
		Ident        string
		Label        string

		Resource     string
		ResourceID   uint64
		ResourceType string

		// Refs is an arbitrary map to identify a model
		// @todo consider reworking this; I'm not the biggest fan
		Refs map[string]any

		SensitivityLevelID uint64

		Attributes AttributeSet

		Constraints map[string][]any
		Indexes     IndexSet
	}

	ModelSet []*Model

	// Attribute describes a specific value of the dataset
	Attribute struct {
		Ident string
		Label string

		SensitivityLevelID uint64

		MultiValue bool

		PrimaryKey bool

		// If attribute has SoftDeleteFlag=true we use it
		// when filtering out deleted items
		SoftDeleteFlag bool

		// System indicates the attribute was defined by the system
		System bool

		// Is attribute sortable?
		// Note: all primary keys are sortable
		Sortable bool

		// Can attribute be used in query expression?
		Filterable bool

		// Store describes the strategy the underlying storage system should
		// apply to the underlying value
		Store Codec

		// Type describes what the value represents and how it should be
		// encoded/decoded
		Type Type
	}

	// auxAttribute is a helper struct used for marshaling/unmarshaling
	//
	// This is required since some fields are interfaces
	auxAttribute struct {
		Ident              string `json:"ident"`
		Label              string `json:"label"`
		SensitivityLevelID uint64 `json:"sensitivityLevelID"`
		MultiValue         bool   `json:"multiValue"`
		PrimaryKey         bool   `json:"primaryKey"`
		SoftDeleteFlag     bool   `json:"softDeleteFlag"`
		System             bool   `json:"system"`
		Sortable           bool   `json:"sortable"`
		Filterable         bool   `json:"filterable"`

		Store *auxAttributeStore `json:"store"`
		Type  *auxAttributeType  `json:"type"`
	}

	// auxAttributeStore is a helper struct used for marshaling/unmarshaling
	//
	// This is required since some fields are interfaces
	auxAttributeStore struct {
		Type string `json:"type"`

		Plain              *CodecPlain              `json:"plain,omitempty"`
		RecordValueSetJSON *CodecRecordValueSetJSON `json:"recordValueSetJSON,omitempty"`
		Alias              *CodecAlias              `json:"alias,omitempty"`
	}

	// auxAttributeType is a helper struct used for marshaling/unmarshaling
	//
	// This is required since some fields are interfaces
	auxAttributeType struct {
		Type string `json:"type"`

		ID        *TypeID        `json:"id,omitempty"`
		Ref       *TypeRef       `json:"ref,omitempty"`
		Timestamp *TypeTimestamp `json:"timestamp,omitempty"`
		Time      *TypeTime      `json:"time,omitempty"`
		Date      *TypeDate      `json:"date,omitempty"`
		Number    *TypeNumber    `json:"number,omitempty"`
		Text      *TypeText      `json:"text,omitempty"`
		Boolean   *TypeBoolean   `json:"boolean,omitempty"`
		Enum      *TypeEnum      `json:"enum,omitempty"`
		Geometry  *TypeGeometry  `json:"geometry,omitempty"`
		JSON      *TypeJSON      `json:"jSON,omitempty"`
		Blob      *TypeBlob      `json:"blob,omitempty"`
		UUID      *TypeUUID      `json:"uuid,omitempty"`
	}

	AttributeSet []*Attribute

	Index struct {
		Ident  string
		Type   string
		Unique bool

		Fields []*IndexField

		Predicate string
	}

	IndexField struct {
		AttributeIdent string
		Modifiers      []IndexFieldModifier
		Sort           IndexFieldSort
		Nulls          IndexFieldNulls
	}

	IndexSet []*Index

	IndexFieldModifier string
	IndexFieldSort     int
	IndexFieldNulls    int
)

const (
	IndexFieldSortDesc IndexFieldSort = -1
	IndexFieldSortAsc  IndexFieldSort = 1

	IndexFieldNullsLast  IndexFieldNulls = -1
	IndexFieldNullsFirst IndexFieldNulls = 1

	IndexFieldModifierLower = "LOWERCASE"
)

func PrimaryAttribute(ident string, codec Codec) *Attribute {
	out := FullAttribute(ident, TypeID{}, codec)
	out.Type = &TypeID{}
	out.PrimaryKey = true
	return out
}

func FullAttribute(ident string, at Type, codec Codec) *Attribute {
	return &Attribute{
		Ident:      ident,
		Label:      ident,
		Sortable:   true,
		Filterable: true,
		Store:      codec,
		Type:       at,
	}
}

func (a *Attribute) WithSoftDelete() *Attribute {
	a.SoftDeleteFlag = true
	return a
}

func (a *Attribute) WithMultiValue() *Attribute {
	a.MultiValue = true
	return a
}

func (a *Attribute) StoreIdent() string {
	switch s := a.Store.(type) {
	case *CodecRecordValueSetJSON:
		return s.Ident

	case *CodecAlias:
		return s.Ident

	default:
		return a.Ident

	}
}

func (mm ModelSet) FindByResourceID(resourceID uint64) *Model {
	for _, m := range mm {
		if m.ResourceID == resourceID {
			return m
		}
	}
	return nil
}

func (mm ModelSet) FindByResourceIdent(resourceType, resourceIdent string) *Model {
	for _, m := range mm {
		if m.ResourceType != resourceType {
			continue
		}

		if m.Resource != resourceIdent {
			continue
		}

		return m
	}
	return nil
}

func (mm ModelSet) FindByIdent(ident string) *Model {
	for _, m := range mm {
		if m.Ident == ident {
			return m
		}
	}

	return nil
}

// FindByRefs returns the first Model that matches the given refs
func (mm ModelSet) FindByRefs(refs map[string]any) *Model {
	for _, model := range mm {
		for k, v := range refs {
			ref, ok := model.Refs[k]
			if !ok {
				goto skip
			}
			if v != ref {
				goto skip
			}
		}

		return model

	skip:
	}
	return nil
}

// FilterByReferenced returns all of the models that reference b
func (aa ModelSet) FilterByReferenced(b *Model) (out ModelSet) {
	for _, aModel := range aa {
		if aModel.Resource == b.Resource {
			continue
		}

		for _, aAttribute := range aModel.Attributes {
			switch casted := aAttribute.Type.(type) {
			case *TypeRef:
				if casted.RefModel.Resource == b.Resource {
					out = append(out, aModel)
				}
			}
		}
	}

	return
}

func (m Model) ToFilter() ModelRef {
	return ModelRef{
		ConnectionID: m.ConnectionID,

		ResourceID: m.ResourceID,

		ResourceType: m.ResourceType,
		Resource:     m.Resource,
	}
}

// HasAttribute returns true when the model includes the specified attribute
func (m Model) HasAttribute(ident string) bool {
	return m.Attributes.FindByIdent(ident) != nil
}

func (aa AttributeSet) FindByIdent(ident string) *Attribute {
	for _, a := range aa {
		if strings.EqualFold(a.Ident, ident) {
			return a
		}
	}

	return nil
}

func (aa AttributeSet) FindByStoreIdent(ident string) *Attribute {
	for _, a := range aa {
		if strings.EqualFold(a.StoreIdent(), ident) {
			return a
		}
	}

	return nil
}

// Validate performs a base model validation before it is passed down
func (m Model) Validate() error {
	if m.Resource == "" {
		return fmt.Errorf("resource not defined")
	}

	seen := make(map[string]bool)
	for _, attr := range m.Attributes {
		if attr.Ident == "" {
			return fmt.Errorf("invalid attribute ident: ident must not be empty")
		}

		if !handle.IsValid(attr.Ident) {
			return fmt.Errorf("invalid attribute ident: %s is not a valid handle", attr.Ident)
		}

		if seen[attr.Ident] {
			return fmt.Errorf("invalid attribute %s: duplicate attributes are not allowed", attr.Ident)
		}
		seen[attr.Ident] = true

		if reflect2.IsNil(attr.Type) {
			return fmt.Errorf("attribute does not define a type: %s", attr.Ident)
		}
	}

	return nil
}

func (a *Attribute) MarshalJSON() ([]byte, error) {
	aux := &auxAttribute{
		Ident:              a.Ident,
		Label:              a.Label,
		SensitivityLevelID: a.SensitivityLevelID,
		MultiValue:         a.MultiValue,
		PrimaryKey:         a.PrimaryKey,
		SoftDeleteFlag:     a.SoftDeleteFlag,
		System:             a.System,
		Sortable:           a.Sortable,
		Filterable:         a.Filterable,

		Store: &auxAttributeStore{},
		Type:  &auxAttributeType{},
	}

	switch s := a.Store.(type) {
	case *CodecPlain:
		aux.Store.Type = "plain"
		aux.Store.Plain = s

	case *CodecRecordValueSetJSON:
		aux.Store.Type = "recordValueSetJSON"
		aux.Store.RecordValueSetJSON = s

	case *CodecAlias:
		aux.Store.Type = "alias"
		aux.Store.Alias = s

	default:
		return nil, fmt.Errorf("unknown store codec type: %T", s)
	}

	switch t := a.Type.(type) {
	case *TypeID:
		aux.Type.Type = "ID"
		aux.Type.ID = t

	case *TypeRef:
		aux.Type.Type = "Ref"
		aux.Type.Ref = t

	case *TypeTimestamp:
		aux.Type.Type = "Timestamp"
		aux.Type.Timestamp = t

	case *TypeTime:
		aux.Type.Type = "Time"
		aux.Type.Time = t

	case *TypeDate:
		aux.Type.Type = "Date"
		aux.Type.Date = t

	case *TypeNumber:
		aux.Type.Type = "Number"
		aux.Type.Number = t

	case *TypeText:
		aux.Type.Type = "Text"
		aux.Type.Text = t

	case *TypeBoolean:
		aux.Type.Type = "Boolean"
		aux.Type.Boolean = t

	case *TypeEnum:
		aux.Type.Type = "Enum"
		aux.Type.Enum = t

	case *TypeGeometry:
		aux.Type.Type = "Geometry"
		aux.Type.Geometry = t

	case *TypeJSON:
		aux.Type.Type = "JSON"
		aux.Type.JSON = t

	case *TypeBlob:
		aux.Type.Type = "Blob"
		aux.Type.Blob = t

	case *TypeUUID:
		aux.Type.Type = "UUID"
		aux.Type.UUID = t

	default:
		return nil, fmt.Errorf("unknown attribute type type: %T", t)
	}

	return json.Marshal(aux)
}

func (a *Attribute) UnmarshalJSON(data []byte) (err error) {
	aux := &auxAttribute{}
	err = json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if a == nil {
		*a = Attribute{}
	}

	a.Ident = aux.Ident
	a.Label = aux.Label
	a.SensitivityLevelID = aux.SensitivityLevelID
	a.MultiValue = aux.MultiValue
	a.PrimaryKey = aux.PrimaryKey
	a.SoftDeleteFlag = aux.SoftDeleteFlag
	a.System = aux.System
	a.Sortable = aux.Sortable
	a.Filterable = aux.Filterable

	switch aux.Store.Type {
	case "plain":
		a.Store = aux.Store.Plain

	case "recordValueSetJSON":
		a.Store = aux.Store.RecordValueSetJSON

	case "alias":
		a.Store = aux.Store.Alias
	}

	switch aux.Type.Type {
	case "ID":
		a.Type = aux.Type.ID

	case "Ref":
		a.Type = aux.Type.Ref

	case "Timestamp":
		a.Type = aux.Type.Timestamp

	case "Time":
		a.Type = aux.Type.Time

	case "Date":
		a.Type = aux.Type.Date

	case "Number":
		a.Type = aux.Type.Number

	case "Text":
		a.Type = aux.Type.Text

	case "Boolean":
		a.Type = aux.Type.Boolean

	case "Enum":
		a.Type = aux.Type.Enum

	case "Geometry":
		a.Type = aux.Type.Geometry

	case "JSON":
		a.Type = aux.Type.JSON

	case "Blob":
		a.Type = aux.Type.Blob

	case "UUID":
		a.Type = aux.Type.UUID

	default:
		return fmt.Errorf("unknown attribute type type: %s", aux.Type.Type)
	}

	return
}
