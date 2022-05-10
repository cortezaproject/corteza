package data

type (
	// temp
	Type interface {
		Type() string
		IsNullable() bool
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

func (t TypeID) Type() string        { return "ID" }
func (t TypeRef) Type() string       { return "REF" }
func (t TypeTimestamp) Type() string { return "TIMESTAMP" }
func (t TypeTime) Type() string      { return "TIME" }
func (t TypeDate) Type() string      { return "DATE" }
func (t TypeNumber) Type() string    { return "NUMBER" }
func (t TypeText) Type() string      { return "TEXT" }
func (t TypeBoolean) Type() string   { return "BOOLEAN" }
func (t TypeEnum) Type() string      { return "ENUM" }
func (t TypeGeometry) Type() string  { return "GEOMETRY" }
func (t TypeJSON) Type() string      { return "JSON" }
func (t TypeBlob) Type() string      { return "BLOB" }
func (t TypeUUID) Type() string      { return "UUID" }
