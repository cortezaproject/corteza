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
	att := types.Reaction{}.New()

	var reaction = ":laugh:"

	att.SetReaction(reaction)

	att, err = rpo.Create(ctx, att)
	must(t, err)
	if att.Reaction != reaction {
		t.Fatal("Changes were not stored")
	}

	att, err = rpo.FindByID(ctx, att.ID)
	must(t, err)
	if att.Reaction != reaction {
		t.Fatal("Changes were not stored")
	}

	must(t, rpo.Delete(ctx, att.ID))
}
