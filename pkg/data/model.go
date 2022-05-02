package data

import (
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/minions"
)

type (
	// Model describes the underlying data and its shape
	Model struct {
		StoreID uint64
		Ident   string

		ResourceID   uint64
		ResourceType string

		Attributes AttributeSet
	}
	ModelSet []*Model

	// Attribute describes a specific value of the dataset
	Attribute struct {
		Ident string

		MultiValue bool

		PrimaryKey     bool
		SoftDeleteFlag bool
		Sortable       bool

		// Store describes the strategy the underlying storage system should
		// apply to the underlying value
		Store StoreCodec
		// Type describes what the value represents and how it should be
		// encoded/decoded
		Type attributeType
	}
	AttributeSet []*Attribute

	AttributeType string

	// temp
	Type interface {
		attributeType
	}

	attributeType interface {
		Type() AttributeType
	}

	// TypeID handles ID (uint64) coding
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeID struct {
		// @todo need to figure out how to support when IDs
		//       generated/provided by store (SERIAL/AUTOINCREMENT)
		GeneratedByStore bool
		Nullable         bool
	}

	// TypeRef handles ID (uint64) coding + reference info
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeRef struct {
		RefModel     *Model
		RefAttribute *Attribute
		Nullable     bool
	}

	// TypeTimestamp handles timestamp coding
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeTimestamp struct {
		Timezone  bool
		Precision uint
		Nullable  bool
	}

	// TypeTime handles time coding
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeTime struct {
		Timezone  bool
		Precision uint
		Nullable  bool
	}

	// TypeDate handles date coding
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeDate struct {
		//
		Nullable bool
	}

	// TypeNumber handles number coding
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeNumber struct {
		Precision uint
		Scale     uint
		Nullable  bool
	}

	// TypeText handles string coding
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeText struct {
		Length   uint
		Nullable bool
	}

	// TypeBoolean
	TypeBoolean struct {
		//
		Nullable bool
	}

	// TypeEnum
	TypeEnum struct {
		Values   []string
		Nullable bool
	}

	// TypeGeometry
	TypeGeometry struct {
		//
		Nullable bool
	}

	// TypeJSON handles coding of arbitrary data into JSON structure
	// NOT TO BE CONFUSED with encodedField
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeJSON struct {
		//
		Nullable bool
	}

	// TypeBlob store/return data as
	TypeBlob struct {
		//
		Nullable bool
	}

	TypeUUID struct {
		//
		Nullable bool
	}
)

const (
	typeID        AttributeType = "id"
	typeRef       AttributeType = "ref"
	typeTimestamp AttributeType = "timestamp"
	typeTime      AttributeType = "time"
	typeDate      AttributeType = "date"
	typeNumber    AttributeType = "number"
	typeText      AttributeType = "text"
	typeBoolean   AttributeType = "boolean"
	typeEnum      AttributeType = "enum"
	typeGeometry  AttributeType = "geometry"
	typeJSON      AttributeType = "json"
	typeBlob      AttributeType = "blob"
	typeUUID      AttributeType = "uuid"
)

const (
	SysID          = "ID"
	SysNamespaceID = "namespaceID"
	SysModuleID    = "moduleID"
	SysCreatedAt   = "createdAt"
	SysCreatedBy   = "createdBy"
	SysUpdatedAt   = "updatedAt"
	SysUpdatedBy   = "updatedBy"
	SysDeletedAt   = "deletedAt"
	SysDeletedBy   = "deletedBy"
	SysOwnedBy     = "ownedBy"
)

// FindByIdent returns the model that matches the ident
func (mm ModelSet) FindByIdent(ident string) *Model {
	for _, m := range mm {
		if m.Ident == ident {
			return m
		}
	}

	return nil
}

// FilterByReferenced returns all of the models that reference b
func (aa ModelSet) FilterByReferenced(b *Model) (out ModelSet) {
	for _, aModel := range aa {
		if aModel.Ident == b.Ident {
			continue
		}

		for _, aAttribute := range aModel.Attributes {
			switch casted := aAttribute.Type.(type) {
			case TypeRef:
				if casted.RefModel.Ident == b.Ident {
					out = append(out, aModel)
				}
			}
		}
	}

	return
}

// HasAttribute returns true when the model includes the specified ident
func (m Model) HasAttribute(ident string) bool {
	return m.Attributes.FindByIdent(ident) != nil
}

func (aa AttributeSet) FindByIdent(ident string) *Attribute {
	for _, a := range aa {
		if strings.ToLower(a.Ident) == strings.ToLower(ident) {
			return a
		}
	}

	return nil
}

// Validate performs a base model validation before it is passed down
func (m Model) Validate() error {
	if m.Ident == "" {
		return fmt.Errorf("ident not defined")
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

		if minions.IsNil(attr.Type) {
			return fmt.Errorf("attribute does not define a type: %s", attr.Ident)
		}
	}

	return nil
}

//// AttributeGroups returns attributes grouped by ident
//func (m Model) AttributeGroups() (gg map[string][]*Attribute) {
//	gg = make(map[string][]*Attribute)
//
//	for _, attr := range m.Attributes {
//		// @todo properly check attribute integrity:
//		//    is embeddable
//		//    only same types use the same ident
//		if gIdent, embeddable := attr.Store.Embeddable(); !embeddable && len(gg[attr.Ident]) > 0 {
//			panic("attribute " + attr.Ident + " is not embeddable")
//		} else if len(gIdent) > 0 {
//			gg[gIdent] = append(gg[gIdent], attr)
//		} else {
//			gg[attr.Ident] = append(gg[attr.Ident], attr)
//		}
//	}
//
//	return
//}

// Receivers to conform to the interface

func (t TypeID) Type() AttributeType {
	return typeID
}

func (t TypeRef) Type() AttributeType {
	return typeRef
}

func (t TypeTimestamp) Type() AttributeType {
	return typeTimestamp
}

func (t TypeTime) Type() AttributeType {
	return typeTime
}

func (t TypeDate) Type() AttributeType {
	return typeDate
}

func (t TypeNumber) Type() AttributeType {
	return typeNumber
}

func (t TypeText) Type() AttributeType {
	return typeText
}

func (t TypeBoolean) Type() AttributeType {
	return typeBoolean
}

func (t TypeEnum) Type() AttributeType {
	return typeEnum
}

func (t TypeGeometry) Type() AttributeType {
	return typeGeometry
}

func (t TypeJSON) Type() AttributeType {
	return typeJSON
}

func (t TypeBlob) Type() AttributeType {
	return typeBlob
}

func (t TypeUUID) Type() AttributeType {
	return typeUUID
}
