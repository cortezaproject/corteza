package types

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/require"
)

// 	Hello! This file is auto-generated.

func TestMessageAttachmentSetWalk(t *testing.T) {
	var (
		value = make(MessageAttachmentSet, 3)
		req   = require.New(t)
	)

	// check walk with no errors
	{
		err := value.Walk(func(*MessageAttachment) error {
			return nil
		})
		req.NoError(err)
	}

	// check walk with error
	req.Error(value.Walk(func(*MessageAttachment) error { return errors.New("walk error") }))

}

func TestMessageAttachmentSetFilter(t *testing.T) {
	var (
		value = make(MessageAttachmentSet, 3)
		req   = require.New(t)
	)

	// filter nothing
	{
		set, err := value.Filter(func(*MessageAttachment) (bool, error) {
			return true, nil
		})
		req.NoError(err)
		req.Equal(len(set), len(value))
	}

	// filter one item
	{
		found := false
		set, err := value.Filter(func(*MessageAttachment) (bool, error) {
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
		_, err := value.Filter(func(*MessageAttachment) (bool, error) {
			return false, errors.New("filter error")
		})
		req.Error(err)
	}
}

func TestMessageAttachmentSetIDs(t *testing.T) {
	var (
		value = make(MessageAttachmentSet, 3)
		req   = require.New(t)
	)

	// construct objects
	value[0] = new(MessageAttachment)
	value[1] = new(MessageAttachment)
	value[2] = new(MessageAttachment)
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
