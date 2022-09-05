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

#ModelAttributeDalType:
	"ID" | "Ref" |
	"Timestamp" | "Time" | "Date" |
	"Number" |
	"Text" |
	"Boolean" |
	"Enum" |
	"Geometry" |
	"JSON" |
	"Blob" |
	"UUID"

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
	// @todo this should be moved to dal.ident
	storeIdent: #ident | *name

	// enable or disable store for this attribute
	// @todo we should use dal prop for this, and extend it to support bool "false"
	//       so that it can be disabled
	store: bool | *true

	unique:     bool | *false
	sortable:   bool | *false
	descending: bool | *false
	primaryKey: bool | *false
	ignoreCase: bool | *false

	// currently disabled since not used by anything
	// it adds more than 4x overhead to the time it takes to generate the store code!
	// #ModelAttributeJsonTag

	dal: {
		type: #ModelAttributeDalType | * "Text"

		fqType: "dal.Type\(type)"

		nullable: bool | *false

		default?: string

		if type == "ID" {
			generatedByStore: bool | *false
		}

		if type == "Ref" {
			refModelResType: #FQRT
			attribute: #handle | *"id"
		}

		if type == "Timestamp" {
			timezone: bool | *false
			precision: number | *0
		}

		if type == "Time" {
			timezone: bool | *false
			precision: number | *0
		}

		if type == "Date" {}

		if type == "Number" {
			precision: number | *0
			scale: number | *0
		}

		if type == "Text" {
			length: number | *0
		}

		if type == "Boolean" {}

		if type == "Enum" {
			values: []
		}

		if type == "Geometry" {}
		if type == "JSON" {}
		if type == "Blob" {}
		if type == "UUID" {}
	}
}

IdField: {
	// Expecting ID field to always have name ID
	name:       "id"
	expIdent:   "ID"
	primaryKey: true
	unique:     true

	// @todo someday we'll replace this with the "ID" type
	goType: "uint64"
	dal: { type: "ID" }
}

HandleField: {
	// Expecting ID field to always have name handle
	name:   "handle"
	unique: true
	ignoreCase: true

	goType: "string"
	dal: { type: "Text", length: 255 }
}

AttributeUserRef: {
	goType: "uint64"
	dal: { type: "Ref", refModelResType: "corteza::system:user" }
}

SortableTimestampField: {
	sortable: true
	goType: "time.Time"
	dal: { type: "Timestamp", precision: 0, timezone: true, nullable: false }
}

SortableTimestampNilField: {
	sortable: true
	goType: "*time.Time"
	dal: { type: "Timestamp", precision: 0, timezone: true, nullable: true }
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
