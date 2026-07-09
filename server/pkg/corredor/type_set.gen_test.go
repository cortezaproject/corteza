package corredor

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScriptSetWalk(t *testing.T) {
	var (
		value = make(ScriptSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Script) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Script) error { return fmt.Errorf("walk error") }))
}

func TestScriptSetFilter(t *testing.T) {
	var (
		value = make(ScriptSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Script) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Script) (bool, error) {
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
		_, err := value.Filter(func(*Script) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}
