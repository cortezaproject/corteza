package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestChannel(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Channel()
	ctx := context.Background()
	att := types.Channel{}.New()

	var name1, name2 = "Test channel v1", "Test channel v2"

	var aa []*types.Channel

	att.SetName(name1)

	att, err = rpo.Create(ctx, att)
	must(t, err)
	if att.Name != name1 {
		t.Fatal("Changes were not stored")
	}

	att.SetName(name2)

	att, err = rpo.Update(ctx, att)
	must(t, err)
	if att.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	att, err = rpo.FindByID(ctx, att.ID)
	must(t, err)
	if att.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	aa, err = rpo.Find(ctx, &types.ChannelFilter{Query: name2})
	must(t, err)
	if len(aa) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.Archive(ctx, att.ID))
	must(t, rpo.Unarchive(ctx, att.ID))
	must(t, rpo.Delete(ctx, att.ID))
}
