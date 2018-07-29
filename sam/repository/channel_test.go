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
	chn := &types.Channel{}

	var name1, name2 = "Test channel v1", "Test channel v2"

	var cc []*types.Channel

	chn.Name = name1

	chn, err = rpo.Create(ctx, chn)
	must(t, err)
	if chn.Name != name1 {
		t.Fatal("Changes were not stored")
	}

	chn.Name = name2

	chn, err = rpo.Update(ctx, chn)
	must(t, err)
	if chn.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	chn, err = rpo.FindByID(ctx, chn.ID)
	must(t, err)
	if chn.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	cc, err = rpo.Find(ctx, &types.ChannelFilter{Query: name2})
	must(t, err)
	if len(cc) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.Archive(ctx, chn.ID))
	must(t, rpo.Unarchive(ctx, chn.ID))
	must(t, rpo.Delete(ctx, chn.ID))
}

func TestChannelMembers(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Channel()
	ctx := context.Background()

	chn := &types.Channel{}
	chn, err = rpo.Create(ctx, chn)
	must(t, err)

	usr := &types.User{}
	usr, err = User().Create(ctx, usr)
	must(t, err)

	must(t, rpo.AddMember(ctx, chn.ID, usr.ID))
	must(t, rpo.RemoveMember(ctx, chn.ID, usr.ID))
}
