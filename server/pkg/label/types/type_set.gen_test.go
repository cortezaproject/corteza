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

func TestLabelSetWalk(t *testing.T) {
	var (
		value = make(LabelSet, 3)
		req   = require.New(t)
	)

	{
		err := value.Walk(func(*Label) error { return nil })
		req.NoError(err)
	}

	req.Error(value.Walk(func(*Label) error { return fmt.Errorf("walk error") }))
}

func TestLabelSetFilter(t *testing.T) {
	var (
		value = make(LabelSet, 3)
		req   = require.New(t)
	)

	{
		set, err := value.Filter(func(*Label) (bool, error) { return true, nil })
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

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

	{
		_, err := value.Filter(func(*Label) (bool, error) {
			return false, fmt.Errorf("filter error")
		})
		req.Error(err)
	}
}
