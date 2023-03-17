package schema

import (
	"strings"
)

#Model: {
	ident: string
	attributes: {
		[name=_]: { "name": name }
	} & {
		[string]: #ModelAttribute
	}

	omitGetterSetter: bool | *false
	defaultGetter: bool | *false
	defaultSetter: bool | *false

	indexes: ({
		[name=_]: { "name": name, "modelIdent": ident } & #ModelIndex
	} & {
		[string]: #ModelIndex
	}) | *({})
}


// logic in struct fields is a bit different
#ModelAttribute: {
	name:   #ident
	_words: strings.Replace(strings.Replace(name, "_", " ", -1), ".", " ", -1)

	_ident: strings.ToCamel(strings.Replace(strings.ToTitle(_words), " ", "", -1))

	// Golang type (built-in or other)
	goType: string | *"string"

	goCastFnc: string | *strings.ToTitle(goType)

	if goType == "*time.Time" {
		goCastFnc: "TimePtr"
	}
	if goType == "time.Time" {
		goCastFnc: "Time"
	}
	if goType == "map[string]any" {
		goCastFnc: "Meta"
	}

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
	ignoreCase: bool | *false

	// currently disabled since not used by anything
	// it adds more than 4x overhead to the time it takes to generate the store code!
	// #ModelAttributeJsonTag

	dal?: #ModelAttributeDal

	envoy: {
		$attrIdent: ident
	} & #attributeEnvoy

	// specifies all of the identifiers this attribute may define when using getters/setters
	identAlias: [...string] | *[ident, expIdent]

	// enable or disable GetValue and SetValue for this attribute
	omitSetter: bool | *false
	omitGetter: bool | *false
}

#attributeEnvoy: {
	$attrIdent: string

	// controls if this attribute can be used as an identifier
	identifier: bool | *false

	// YAML decode/encode configs
	yaml: {
		identKeyLabel: string | *strings.ToLower($attrIdent)
		identKeyAlias: [...string] | *[]
		// identKeys defines what identifiers this attribute supports when
		// decoding yaml documents
		identKeys: ([identKeyLabel]+identKeyAlias)

		customDecoder: bool | *false
		customEncoder: bool | *false

		identKeyEncode: string | *$attrIdent

		omitEncoder: bool | *false
	}

	// store decode/encode configs
	store: {
		// defines a custom field identifier when constructing
		// resource filters and assigning reference constraints
		filterRefField: string | *""

		omitRefFilter: bool | *false
	}
}

#ModelAttributeDal: {
	type: #ModelAttributeDalType | *"Text"

	fqType: "dal.Type\(type)"

	nullable: bool | *false


	if type == "ID" {
		generatedByStore: bool | *false
		default?: 0
	}

	if type == "Ref" {
		refModelResType: #FQRT
		attribute: #handle | *"id"
		default?: 0
	}

	if type == "Timestamp" {
		timezone: bool | *false
		precision: number | *(-1)
		defaultCurrentTimestamp?: true
	}

	if type == "Time" {
		timezone: bool | *false
		precision: number | *(-1)
		defaultCurrentTimestamp?: true
	}

	if type == "Date" {
		defaultCurrentTimestamp?: true
	}

	if type == "Number" {
		precision: number | *(-1)
		scale: number | *(-1)
		default?: number
		meta?: { [string]: _ }
	}

	if type == "Text" {
		length: number | *0
		default?: string
	}

	if type == "Boolean" {
		default?: bool
	}

	if type == "Enum" {
		values: []
		default?: string
	}

	if type == "Geometry" {}

	if type == "JSON" {
		default?: string | bytes
		defaultEmptyObject?: true
	}

	if type == "Blob" {
		default?: bytes
	}

	if type == "UUID" {}
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


IdField: {
	// Expecting ID field to always have name ID
	name:       "id"
	expIdent:   "ID"
	unique:     true

	// @todo someday we'll replace this with the "ID" type
	goType: "uint64"
	dal: { type: "ID" }

	envoy: #attributeEnvoy & {
		identifier: true
	}
}
HandleField: {
	// Expecting ID field to always have name handle
	name:   "handle"
	unique: true
	ignoreCase: true

	goType: "string"
	dal: { type: "Text", length: 64 }

	envoy: #attributeEnvoy & {
		identifier: true
	}
}

AttributeUserRef: {
	goType: "uint64"
	dal: { type: "Ref", refModelResType: "corteza::system:user", default: 0 }
	envoy: {
		store: {
			omitRefFilter: true
		}
	}
}

SortableTimestampField: {
	sortable: true
	goType: "time.Time"
	dal: { type: "Timestamp", timezone: true, nullable: false }
}

SortableTimestampNowField: {
	sortable: true
	goType: "time.Time"
	dal: { type: "Timestamp", timezone: true, nullable: false, defaultCurrentTimestamp: true }
}

SortableTimestampNilField: {
	sortable: true
	goType: "*time.Time"
	dal: { type: "Timestamp", timezone: true, nullable: true }
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

#ModelIndex: close({
	name:   #ident
	modelIdent: #ident
	_attributes: { [_]: #ModelAttribute }
	_words: strings.Replace(strings.Replace(name, "_", " ", -1), ".", " ", -1)

	_ident: strings.ToCamel(strings.Replace(strings.ToTitle(_words), " ", "", -1))

	// lowercase (unexported, golang) identifier
	ident: #ident | *"\(modelIdent)_\(_ident)"

	primary: bool | *(strings.ToLower(name) == "primary")
	unique: bool  | *(strings.Contains(name, "unique") || primary)

	type: "BTREE" | *"BTREE"

 	// index predicate,
 	// condition that must be met for the index to be used
	predicate?: string

	attribute?: string
	attributes?: [string, ...]
	fields: [#ModelIndexField, ...]

	if attribute != _|_ {
		attributes: [string, ...] & [attribute]
	}

	if fields != _|_ && attributes != _|_ {
		fields: [
			for a in attributes {
				{"attribute": a} & #ModelIndexField
			}
		]
	}
})

#IndexFieldModifier: "LOWERCASE"

#ModelIndexField: close({
  attribute: string
	modifiers?: [#IndexFieldModifier, ...]
	length?: number
	sort?: "DESC" | "ASC"
	nulls?: "FIRST" | "LAST"
})
