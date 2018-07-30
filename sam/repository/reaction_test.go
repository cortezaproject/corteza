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

	react.Reaction = reaction

	react, err = rpo.CreateReaction(react)
	must(t, err)
	if react.Reaction != reaction {
		t.Fatal("Changes were not stored")
	}

	react, err = rpo.FindReactionByID(react.ID)
	must(t, err)
	if react.Reaction != reaction {
		t.Fatal("Changes were not stored")
	}

	must(t, rpo.DeleteReactionByID(react.ID))
}
