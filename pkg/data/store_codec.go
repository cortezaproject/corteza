package data

type (
	// StoreStrategy defines how values for a specific model attribute
	// are retrieved or stored
	//
	// If attribute does not have store strategy set
	// store driver should use attribute name to determinate
	// source/destination of the value (table column, json doc property)
	StoreStrategy any

	// StoreCodecJSON defines that values are encoded/decoded into
	// a JSON document (stored under ident)
	// Path defines exact location inside the doc.
	StoreCodecJSON struct {
		Ident string
		Path  []any
	}

	// { "@value": ... "@type": .... }
	// StoreCodecJSONLD struct { Ident  string; Path   []any }

	// StoreCodecXML
	//StoreCodecXML struct {}

	// StoreCodecAlias defines case when value is not stored
	// under the same column (in case of an SQL table)
	//
	// Example: attribute ident is "foo" and database table column is "bar"
	StoreCodecAlias struct {
		Ident string
	}
)
