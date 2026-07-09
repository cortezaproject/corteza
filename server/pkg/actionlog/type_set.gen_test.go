package actionlog

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

func TestActionSetWalk(t *testing.T) {
	var (
		value = make(ActionSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Action) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Action) error { return fmt.Errorf("walk error") }))
}

func TestActionSetFilter(t *testing.T) {
	var (
		value = make(ActionSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Action) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*Action) (bool, error) {
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
		_, err := value.Filter(func(*Action) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestActionSetIDs(t *testing.T) {
	var (
		value = make(ActionSet, 3)
		req   = require.New(t)
	)

	value[0] = new(Action)
	value[1] = new(Action)
	value[2] = new(Action)
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
