package data

type (
	StoreCodecType string

	// StoreCodec defines how values for a specific model attribute
	// are retrieved or stored
	//
	// If attribute does not have store strategy set
	// store driver should use attribute name to determinate
	// source/destination of the value (table column, json document property)
	StoreCodec interface {
		Type() StoreCodecType
	}

	StoreCodecPlain struct{}

	// StoreCodecStdRecordValueJSON defines that values are encoded/decoded into
	// a JSON simple document { [_: string]: Array<unknown> }
	//
	// This only solves
	// Attribute{Ident: "foo", Store: StoreCodecRecordValueJSON{ Ident: "bar" }
	// => "bar"->'foo'->0
	StoreCodecStdRecordValueJSON struct {
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

	// StoreCodecAlias defines case when value is not stored
	// under the same column (in case of an SQL table)
	//
	// Value of StoreCodecAlias.Ident is used as base
	// and value of Attribute.Ident holding StoreCodecAlias is used
	// as an alias!
	//
	// Attribute{Ident: "foo", Store: StoreCodecAlias{ Ident: "bar" }
	// => "bar" as "foo"
	StoreCodecAlias struct {
		Ident string
	}
)

const (
	StoreCodecPlainType              StoreCodecType = "corteza::data-store-codec:plain"
	StoreCodecStdRecordValueJSONType StoreCodecType = "corteza::data-store-codec:record-store-json"
	StoreCodecAliasType              StoreCodecType = "corteza::data-store-codec:alias"
)

func (t StoreCodecType) Is(c StoreCodec) bool { return t == c.Type() }

func (StoreCodecPlain) Type() StoreCodecType              { return StoreCodecPlainType }
func (StoreCodecStdRecordValueJSON) Type() StoreCodecType { return StoreCodecStdRecordValueJSONType }
func (StoreCodecAlias) Type() StoreCodecType              { return StoreCodecAliasType }
