package types

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

func TestResourceActivitySetWalk(t *testing.T) {
	var (
		value = make(ResourceActivitySet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*ResourceActivity) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*ResourceActivity) error { return fmt.Errorf("walk error") }))
}

func TestResourceActivitySetFilter(t *testing.T) {
	var (
		value = make(ResourceActivitySet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*ResourceActivity) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	{
		found := false
		set, err := value.Filter(func(*ResourceActivity) (bool, error) {
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
		_, err := value.Filter(func(*ResourceActivity) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}

func TestResourceActivitySetIDs(t *testing.T) {
	var (
		value = make(ResourceActivitySet, 3)
		req   = require.New(t)
	)

	value[0] = new(ResourceActivity)
	value[1] = new(ResourceActivity)
	value[2] = new(ResourceActivity)
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
