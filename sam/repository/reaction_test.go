package repository

import (
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestReaction(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := New()
	react := &types.Reaction{}

	var reaction = ":laugh:"

	{
		react.Reaction = reaction
		react, err = rpo.CreateReaction(react)
		assert(t, err == nil, "CreateReaction error: %v", err)
		assert(t, react.Reaction == reaction, "Changes were not stored")

		{
			react, err = rpo.FindReactionByID(react.ID)
			assert(t, err == nil, "FindReactionByID error: %v", err)
			assert(t, react.Reaction == reaction, "Changes were not stored")
		}

		{
			err = rpo.DeleteReactionByID(react.ID)
			assert(t, err == nil, "DeleteReactionByID error: %v", err)
		}
	}
}
