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
	usr := &types.User{}

	var name1, name2 = "Test user v1", "Test user v2"

	var aa []*types.User

	usr.Username = name1

	usr, err = rpo.Create(ctx, usr)
	must(t, err)
	if usr.Username != name1 {
		t.Fatal("Changes were not stored")
	}

	usr.Username = name2

	usr, err = rpo.Update(ctx, usr)
	must(t, err)
	if usr.Username != name2 {
		t.Fatal("Changes were not stored")
	}

	usr, err = rpo.FindByID(ctx, usr.ID)
	must(t, err)
	if usr.Username != name2 {
		t.Fatal("Changes were not stored")
	}

	aa, err = rpo.Find(ctx, &types.UserFilter{Query: name2})
	must(t, err)
	if len(aa) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.Delete(ctx, usr.ID))
}
