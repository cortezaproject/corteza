package {{ .Package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

import (
    "fmt"
    "github.com/stretchr/testify/require"
    "testing"

	{{ range $i, $import := .Imports }}
	"{{ $import }}"
	{{ end }}
)


{{ range $name, $set := .Types }}

func Test{{ $name }}SetWalk(t *testing.T) {
	var (
		value = make({{ $name }}Set, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*{{ $name }}) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*{{ $name }}) error { return fmt.Errorf("walk error") }))
}

func Test{{ $name }}SetFilter(t *testing.T) {
	var (
		value = make({{ $name }}Set, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*{{ $name }}) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*{{ $name }}) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	// filter error
	{
		_, err := value.Filter(func(*{{ $name }}) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

{{ if not $set.NoIdField }}
func Test{{ $name }}SetIDs(t *testing.T) {
	var (
		value = make({{ $name }}Set, 3)
		req   = require.New(t)
	)

	// construct objects
	value[0] = new({{ $name }})
	value[1] = new({{ $name }})
	value[2] = new({{ $name }})
	// set ids
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	// Find existing
	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	// Find non-existing
	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	// List IDs from set
	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}
{{ end }}

{{ end }}
