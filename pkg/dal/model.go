package dal

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

		PrimaryKey bool

		// If attribute has SoftDeleteFlag=true we use it
		// when filtering out deleted items
		SoftDeleteFlag bool

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

	AttributeSet []*Attribute
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
			case *TypeRef:
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
