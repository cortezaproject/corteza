package types

/*
	Hello! This file is auto-generated from `docs/src/spec.json`.

	For development:
	In order to update the generated files, edit this file under the location,
	add your struct fields, imports, API definitions and whatever you want, and:

	1. run [spec](https://github.com/titpetric/spec) in the same folder,
	2. run `./_gen.php` in this folder.

	You may edit `field.go`, `field.util.go` or `field_test.go` to
	implement your API calls, helper functions and tests. The file `field.go`
	is only generated the first time, and will not be overwritten if it exists.
*/

type (
	// Fields
	Field struct {
		Name string `json:"name" db:"name"`
		Type string `json:"type" db:"type"`

		changed []string
	}
)

/* Constructors */
func (Field) New() *Field {
	return &Field{}
}

/* Getters/setters */
func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) SetName(value string) *Field {
	if f.Name != value {
		f.changed = append(f.changed, "Name")
		f.Name = value
	}
	return f
}
func (f *Field) GetType() string {
	return f.Type
}

func (f *Field) SetType(value string) *Field {
	if f.Type != value {
		f.changed = append(f.changed, "Type")
		f.Type = value
	}
	return f
}
