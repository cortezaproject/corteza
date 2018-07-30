package repository

import (
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestChannel(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := New()
	chn := &types.Channel{}

	var name1, name2 = "Test channel v1", "Test channel v2"

	var cc []*types.Channel

	chn.Name = name1

	chn, err = rpo.CreateChannel(chn)
	must(t, err)
	if chn.Name != name1 {
		t.Fatal("Changes were not stored")
	}

	chn.Name = name2

	chn, err = rpo.UpdateChannel(chn)
	must(t, err)
	if chn.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	chn, err = rpo.FindChannelByID(chn.ID)
	must(t, err)
	if chn.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	cc, err = rpo.FindChannels(&types.ChannelFilter{Query: name2})
	must(t, err)
	if len(cc) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.ArchiveChannelByID(chn.ID))
	must(t, rpo.UnarchiveChannelByID(chn.ID))
	must(t, rpo.DeleteChannelByID(chn.ID))
}

func TestChannelMembers(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := New()

	chn := &types.Channel{}
	chn, err = rpo.CreateChannel(chn)
	must(t, err)

	usr := &types.User{}
	usr, err = rpo.CreateUser(usr)
	must(t, err)

	must(t, rpo.AddChannelMember(chn.ID, usr.ID))
	must(t, rpo.RemoveChannelMember(chn.ID, usr.ID))
}
