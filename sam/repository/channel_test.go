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

	chn, err = rpo.CreateChannel(ctx, chn)
	must(t, err)
	if chn.Name != name1 {
		t.Fatal("Changes were not stored")
	}

	chn.Name = name2

	chn, err = rpo.UpdateChannel(ctx, chn)
	must(t, err)
	if chn.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	chn, err = rpo.FindChannelByID(ctx, chn.ID)
	must(t, err)
	if chn.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	cc, err = rpo.FindChannels(ctx, &types.ChannelFilter{Query: name2})
	must(t, err)
	if len(cc) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.ArchiveChannel(ctx, chn.ID))
	must(t, rpo.UnarchiveChannel(ctx, chn.ID))
	must(t, rpo.DeleteChannelByID(ctx, chn.ID))
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
	chn, err = rpo.CreateChannel(ctx, chn)
	must(t, err)

	usr := &types.User{}
	usr, err = User().Create(ctx, usr)
	must(t, err)

	must(t, rpo.AddChannelMember(ctx, chn.ID, usr.ID))
	must(t, rpo.RemoveChannelMember(ctx, chn.ID, usr.ID))
}
