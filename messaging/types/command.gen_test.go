package types

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/require"
)

// 	Hello! This file is auto-generated.

func TestCommandSetWalk(t *testing.T) {
	var (
		value = make(CommandSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*Command) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*Command) error { return errors.New("walk error") }))

}

func TestCommandSetFilter(t *testing.T) {
	var (
		value = make(CommandSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*Command) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*Command) (bool, error) {
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
		_, err := value.Filter(func(*Command) (bool, error) {
			return false, errors.New("filter error")
		})
		req.Error(err)
	}
}
