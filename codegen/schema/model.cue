package schema

import (
	"strings"
)

#Model: {
	ident: string
	attributes: {
		[name=_]: {"name": name} & #ModelAttribute
	}
}

// logic in struct fields is a bit different
#ModelAttribute: {
	name:   #ident
	_words: strings.Replace(strings.Replace(name, "_", " ", -1), ".", " ", -1)

	_ident: strings.ToCamel(strings.Replace(strings.ToTitle(_words), " ", "", -1))

	// Golang type (built-in or other)
	goType: string | *"string"

	// lowercase (unexported, golang) identifier
	ident: #ident | *_ident

	// uppercase (exported, golang) identifier
	expIdent: #expIdent | *strings.ToTitle(ident)

	// store identifier
	storeIdent: #ident | *name
	store: bool | *true

	unique:     bool | *false
	sortable:   bool | *false
	descending: bool | *false
	primaryKey: bool | *false
	ignoreCase: bool | *false

	// currently disabled since not used by anything
	// it adds more than 4x overhead to the time it takes to generate the store code!
	// #ModelAttributeJsonTag
}

IdField: {
	// Expecting ID field to always have name ID
	name:       "id"
	expIdent:   "ID"
	primaryKey: true
	unique:     true

	// @todo someday we'll replace this with the "ID" type
	goType: "uint64"
}

HandleField: {
	// Expecting ID field to always have name handle
	name:   "handle"
	unique: true
	ignoreCase: true

	goType: "string"
}

SortableTimestampField: {
	sortable: true
	goType: "time.Time"
}

SortableTimestampNilField: {
	sortable: true
	goType: "*time.Time"
}

#ModelAttributeJsonTag: {
	name: string

	_specs: {field: string | *name, omitEmpty: bool | *false, "string": bool | *false}

	json:     string | _specs | bool | *false
	jsonTag?: string

	// just wrap whatever we got in json
	if (json & string) != _|_ {
		jsonTag: "json:\"\(json)\""
	}

	// json enable,d wrap with ident as a JSON prop name
	if (json & bool) != _|_ && json {
		// generic json tag
		jsonTag: "json:\"\(name)\""
	}

	// full-specs
	if (json & bool) == _|_ && (json & _specs) != _|_ {
		_omitEmpty: string | *""
		if json.omitEmpty {
			_omitEmpty: ",omitempty"
		}
		_string: string | *""
		if json.string {
			_string: ",string"
		}

		jsonTag: "json:\"\(json.field)\(_omitEmpty)\(_string)\""
	}
}
