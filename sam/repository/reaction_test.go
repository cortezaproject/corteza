package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestReaction(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Reaction()
	ctx := context.Background()
	react := &types.Reaction{}

	var reaction = ":laugh:"

	react.Reaction = reaction

	react, err = rpo.Create(ctx, react)
	must(t, err)
	if react.Reaction != reaction {
		t.Fatal("Changes were not stored")
	}

	react, err = rpo.FindByID(ctx, react.ID)
	must(t, err)
	if react.Reaction != reaction {
		t.Fatal("Changes were not stored")
	}

	must(t, rpo.Delete(ctx, react.ID))
}
