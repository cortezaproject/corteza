package types

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/label/types/types.yaml

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLabelSetWalk(t *testing.T) {
	var (
		value = make(LabelSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*Label) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*Label) error { return fmt.Errorf("walk error") }))
}

func TestLabelSetFilter(t *testing.T) {
	var (
		value = make(LabelSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*Label) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*Label) (bool, error) {
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
		_, err := value.Filter(func(*Label) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}
