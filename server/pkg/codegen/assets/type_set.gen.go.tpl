package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}


{{ if .Imports }}
import (
	{{ range $i, $import := .Imports }}
	"{{ $import }}"
	{{ end }}
)
{{ end }}

type (
{{ range $name, $set := .Types }}
	// {{ $name }}Set slice of {{ $name }}
	//
	// This type is auto-generated.
	{{ $name }}Set []*{{ $name }}
{{ end }}
)

{{ range $name, $set := .Types }}
// Walk iterates through every slice item and calls w({{ $name }}) err
//
// This function is auto-generated.
func (set {{ $name }}Set) Walk(w func(*{{ $name }}) error) (err error) {
	for i := range set {
		if err = w(set[i]); err != nil {
			return
		}
	}

	return
}

// Filter iterates through every slice item, calls f({{ $name }}) (bool, err) and return filtered slice
//
// This function is auto-generated.
func (set {{ $name }}Set) Filter(f func(*{{ $name }}) (bool, error)) (out {{ $name }}Set, err error) {
	var ok bool
	out = {{ $name }}Set{}
	for i := range set {
		if ok, err = f(set[i]); err != nil {
			return
		} else if ok {
			out = append(out, set[i])
		}
	}

	return
}

{{ if not $set.NoIdField }}
// FindByID finds items from slice by its ID property
//
// This function is auto-generated.
func (set {{ $name }}Set) FindByID(ID uint64) *{{ $name }} {
	for i := range set {
		if set[i].ID == ID {
			return set[i]
		}
	}

	return nil
}

// IDs returns a slice of uint64s from all items in the set
//
// This function is auto-generated.
func (set {{ $name }}Set) IDs() (IDs []uint64) {
	IDs = make([]uint64, len(set))

	for i := range set {
		IDs[i] = set[i].ID
	}

	return
}
{{ end }}


{{ end }}
