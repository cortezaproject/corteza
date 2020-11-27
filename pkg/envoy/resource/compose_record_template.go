package resource

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
		Key      []string
		FieldMap MappingTplSet
	}
)

// NewComposeRecordTemplate returns a record template based on the provided parameters.
//
// The template is only valid up until the resource shaping; it is not allowed after the fact.
func NewComposeRecordTemplate(modRef, nsRef, name string, fieldMap MappingTplSet, key ...string) *ComposeRecordTemplate {
	r := &ComposeRecordTemplate{base: &base{}}
	r.Name = name
	r.Key = key
	r.FieldMap = fieldMap

	r.ModRef = r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, modRef)
	r.NsRef = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)

	r.SetResourceType(DATA_SOURCE_RESOURCE_TYPE)
	r.AddIdentifier(identifiers(name)...)
	return r
}
