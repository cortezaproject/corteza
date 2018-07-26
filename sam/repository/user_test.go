package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestUser(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := User()
	ctx := context.Background()
	att := types.User{}.New()

	var name1, name2 = "Test user v1", "Test user v2"

	var aa []*types.User

	att.SetUsername(name1)

	att, err = rpo.Create(ctx, att)
	must(t, err)
	if att.Username != name1 {
		t.Fatal("Changes were not stored")
	}

	att.SetUsername(name2)

	att, err = rpo.Update(ctx, att)
	must(t, err)
	if att.Username != name2 {
		t.Fatal("Changes were not stored")
	}

	att, err = rpo.FindByID(ctx, att.ID)
	must(t, err)
	if att.Username != name2 {
		t.Fatal("Changes were not stored")
	}

	aa, err = rpo.Find(ctx, &types.UserFilter{Query: name2})
	must(t, err)
	if len(aa) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.Delete(ctx, att.ID))
}
