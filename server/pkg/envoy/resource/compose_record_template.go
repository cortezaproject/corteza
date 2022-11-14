package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	MappingTpl struct {
		// Index is like Cell, but the index
		Index uint
		Cell  string
		Field string

		// If specifies when the field should be mapped
		// @todo
		If string

		// Expr specifies how to aditionally manipulate the value
		// @todo
		Expr string
	}
	MappingTplSet []*MappingTpl

	ComposeRecordTemplate struct {
		// We'll do this so it conforms to the resource.Interface
		*base

		ModRef *Ref
		NsRef  *Ref

		// Name is the source name; topically file name
		Name string
		// Key determines what defines an identifier
		Key         []string
		FieldMap    MappingTplSet
		Defaultable bool
	}
)

// NewComposeRecordTemplate returns a record template based on the provided parameters.
//
// The template is only valid up until the resource shaping; it is not allowed after the fact.
func NewComposeRecordTemplate(modRef, nsRef, name string, defaultable bool, fieldMap MappingTplSet, key ...string) *ComposeRecordTemplate {
	r := &ComposeRecordTemplate{base: &base{}}
	r.Name = name
	r.Key = key
	r.FieldMap = fieldMap
	r.Defaultable = defaultable

	r.NsRef = r.AddRef(types.NamespaceResourceType, nsRef)
	r.ModRef = r.AddRef(types.ModuleResourceType, modRef).Constraint(r.NsRef)

	r.SetResourceType(DataSourceResourceType)
	r.AddIdentifier(identifiers(name)...)
	return r
}

func (r *ComposeRecordTemplate) Resource() interface{} {
	return nil
}

// MapToMappingTplSet converts the given string map to a propper MappingTplSet
func MapToMappingTplSet(b map[string]string) MappingTplSet {
	mp := make(MappingTplSet, 0, len(b))
	for c, f := range b {
		mp = append(mp, &MappingTpl{
			Cell:  c,
			Field: f,
		})
	}
	return mp
}

func (t *MappingTpl) IsIgnored() bool {
	// @todo expand this?
	return t.Field == "/"
}
