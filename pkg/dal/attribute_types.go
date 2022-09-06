package dal

type (
	AttributeType string

	// temp
	Type interface {
		Type() AttributeType
		IsNullable() bool
	}

	// TypeID handles ID (uint64) coding
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	//
	// It is always Corteza ID
	TypeID struct {
		GeneratedByStore bool
		Nullable         bool
	}

	// TypeRef handles ID (uint64) coding + reference info
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeRef struct {
		// defaults to ID
		RefAttribute string
		RefModel     *ModelRef

		Nullable bool
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

		DefaultCurrentTimestamp bool
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

		DefaultCurrentTimestamp bool
	}

	// TypeDate handles date coding
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeDate struct {
		//
		Nullable bool

		DefaultCurrentTimestamp bool
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

		HasDefault   bool
		DefaultValue float64
	}

	// TypeText handles string coding
	//
	// Encoding/decoding might be different depending on
	//  1) underlying store (and dialect)
	//  2) value codec (raw, json ...)
	TypeText struct {
		Length   uint
		Nullable bool

		HasDefault   bool
		DefaultValue string
	}

	// TypeBoolean
	TypeBoolean struct {
		//
		Nullable bool

		HasDefault   bool
		DefaultValue bool
	}

	// TypeEnum
	TypeEnum struct {
		Values   []string
		Nullable bool

		HasDefault   bool
		DefaultValue string
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

		HasDefault   bool
		DefaultValue any
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
	AttributeTypeID        AttributeType = "corteza::dal:attribute-type:id"
	AttributeTypeRef       AttributeType = "corteza::dal:attribute-type:ref"
	AttributeTypeTimestamp AttributeType = "corteza::dal:attribute-type:timestamp"
	AttributeTypeTime      AttributeType = "corteza::dal:attribute-type:time"
	AttributeTypeDate      AttributeType = "corteza::dal:attribute-type:date"
	AttributeTypeNumber    AttributeType = "corteza::dal:attribute-type:number"
	AttributeTypeText      AttributeType = "corteza::dal:attribute-type:text"
	AttributeTypeBoolean   AttributeType = "corteza::dal:attribute-type:boolean"
	AttributeTypeEnum      AttributeType = "corteza::dal:attribute-type:enum"
	AttributeTypeGeometry  AttributeType = "corteza::dal:attribute-type:geometry"
	AttributeTypeJSON      AttributeType = "corteza::dal:attribute-type:json"
	AttributeTypeBlob      AttributeType = "corteza::dal:attribute-type:blob"
	AttributeTypeUUID      AttributeType = "corteza::dal:attribute-type:uuid"
)

func (t TypeID) IsNullable() bool        { return t.Nullable }
func (t TypeRef) IsNullable() bool       { return t.Nullable }
func (t TypeTimestamp) IsNullable() bool { return t.Nullable }
func (t TypeTime) IsNullable() bool      { return t.Nullable }
func (t TypeDate) IsNullable() bool      { return t.Nullable }
func (t TypeNumber) IsNullable() bool    { return t.Nullable }
func (t TypeText) IsNullable() bool      { return t.Nullable }
func (t TypeBoolean) IsNullable() bool   { return t.Nullable }
func (t TypeEnum) IsNullable() bool      { return t.Nullable }
func (t TypeGeometry) IsNullable() bool  { return t.Nullable }
func (t TypeJSON) IsNullable() bool      { return t.Nullable }
func (t TypeBlob) IsNullable() bool      { return t.Nullable }
func (t TypeUUID) IsNullable() bool      { return t.Nullable }

func (t TypeID) Type() AttributeType        { return AttributeTypeID }
func (t TypeRef) Type() AttributeType       { return AttributeTypeRef }
func (t TypeTimestamp) Type() AttributeType { return AttributeTypeTimestamp }
func (t TypeTime) Type() AttributeType      { return AttributeTypeTime }
func (t TypeDate) Type() AttributeType      { return AttributeTypeDate }
func (t TypeNumber) Type() AttributeType    { return AttributeTypeNumber }
func (t TypeText) Type() AttributeType      { return AttributeTypeText }
func (t TypeBoolean) Type() AttributeType   { return AttributeTypeBoolean }
func (t TypeEnum) Type() AttributeType      { return AttributeTypeEnum }
func (t TypeGeometry) Type() AttributeType  { return AttributeTypeGeometry }
func (t TypeJSON) Type() AttributeType      { return AttributeTypeJSON }
func (t TypeBlob) Type() AttributeType      { return AttributeTypeBlob }
func (t TypeUUID) Type() AttributeType      { return AttributeTypeUUID }
