package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

{{- range .types }}

func Test{{ .expIdent }}SetWalk(t *testing.T) {
	var (
		value = make({{ .expIdent }}Set, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*{{ .expIdent }}) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*{{ .expIdent }}) error { return fmt.Errorf("walk error") }))
}

func Test{{ .expIdent }}SetFilter(t *testing.T) {
	var (
		value = make({{ .expIdent }}Set, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*{{ .expIdent }}) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*{{ .expIdent }}) (bool, error) {
			if !found {
				found = true
				return found, nil
			}
			return false, nil
		})
		req.NoError(err)
		req.Len(set, 1)
	}

	{
		_, err := value.Filter(func(*{{ .expIdent }}) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}
{{- if not .noIdField }}

func Test{{ .expIdent }}SetIDs(t *testing.T) {
	var (
		value = make({{ .expIdent }}Set, 3)
		req   = require.New(t)
	)

	value[0] = new({{ .expIdent }})
	value[1] = new({{ .expIdent }})
	value[2] = new({{ .expIdent }})
	value[0].ID = 1
	value[1].ID = 2
	value[2].ID = 3

	{
		val := value.FindByID(2)
		req.Equal(uint64(2), val.ID)
	}

	{
		val := value.FindByID(4)
		req.Nil(val)
	}

	{
		val := value.IDs()
		req.Equal(len(val), len(value))
	}
}
{{- end }}
{{- end }}
