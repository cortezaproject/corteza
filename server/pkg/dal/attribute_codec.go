package dal

type (
	AttributeCodecType string

	// Codec defines how values for a specific model attribute
	// are retrieved or stored
	//
	// If attribute does not have store strategy set
	// store driver should use attribute name to determinate
	// source/destination of the value (table column, json document property)
	Codec interface {
		Type() AttributeCodecType
		SingleValueOnly() bool
	}

	CodecPlain struct{}

	// CodecRecordValueSetJSON defines that values are encoded/decoded into
	// a JSON simple document { [_: string]: Array<unknown> }
	//
	// This only solves
	// Attribute{Ident: "foo", Store: StoreCodecRecordValueJSON{ Ident: "bar" }
	// => "bar"->'foo'->0
	CodecRecordValueSetJSON struct {
		Ident string
	}

	// handling complex JSON documents
	//StoreCodecJSON struct {
	//	Ident string
	//	Path  []any
	//}
	//
	// { "@value": ... "@type": .... }
	// StoreCodecJSONLD struct { Ident  string; Path   []any }

	// StoreCodecXML
	//StoreCodecXML struct {}

	// CodecAlias defines case when value is not stored
	// under the same column (in case of an SQL table)
	//
	// Value of CodecAlias.Ident is used as base
	// and value of Attribute.Ident holding CodecAlias is used
	// as an alias!
	//
	// Attribute{Ident: "foo", Store: CodecAlias{ Ident: "bar" }
	// => "bar" as "foo"
	CodecAlias struct {
		Ident string
	}
)

const (
	AttributeCodecPlain              AttributeCodecType = "corteza::dal:attribute-codec:plain"
	AttributeCodecRecordValueSetJSON AttributeCodecType = "corteza::dal:attribute-codec:record-value-set-json"
	AttributeCodecAlias              AttributeCodecType = "corteza::dal:attribute-codec:alias"
)

func (*CodecPlain) Type() AttributeCodecType { return AttributeCodecPlain }
func (*CodecRecordValueSetJSON) Type() AttributeCodecType {
	return AttributeCodecRecordValueSetJSON
}
func (*CodecAlias) Type() AttributeCodecType { return AttributeCodecAlias }

func (*CodecPlain) SingleValueOnly() bool              { return true }
func (*CodecRecordValueSetJSON) SingleValueOnly() bool { return false }
func (*CodecAlias) SingleValueOnly() bool              { return true }
