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

	{
		chn.Name = name1
		chn, err = rpo.CreateChannel(chn)
		assert(t, err == nil, "CreateChannel error: %v", err)
		assert(t, chn.Name == name1, "Changes were not stored")

		{
			chn.Name = name2

			chn, err = rpo.UpdateChannel(chn)
			assert(t, err == nil, "UpdateChannel error: %v", err)
			assert(t, chn.Name == name2, "Changes were not stored")
		}

		{
			chn, err = rpo.FindChannelByID(chn.ID)
			assert(t, err == nil, "FindChannelByID error: %v", err)
			assert(t, chn.Name == name2, "Changes were not stored")
		}

		{
			cc, err = rpo.FindChannels(&types.ChannelFilter{Query: name2})
			assert(t, err == nil, "FindChannels error: %v", err)
			assert(t, len(cc) > 0, "No results found")
		}

		{
			err = rpo.ArchiveChannelByID(chn.ID)
			assert(t, err == nil, "ArchiveChannelByID error: %v", err)
		}

		{
			err = rpo.UnarchiveChannelByID(chn.ID)
			assert(t, err == nil, "UnarchiveChannelByID error: %v", err)
		}

		{
			err = rpo.DeleteChannelByID(chn.ID)
			assert(t, err == nil, "DeleteChannelByID error: %v", err)
		}
	}
}

func TestChannelMembers(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := New()
	chn := &types.Channel{}
	usr := &types.User{}

	{
		chn, err = rpo.CreateChannel(chn)
		assert(t, err == nil, "CreateChannel: %v", err)

		{
			usr, err = rpo.CreateUser(usr)
			assert(t, err == nil, "CreateUser error: %v", err)

			{
				_, err = rpo.AddChannelMember(&types.ChannelMember{ChannelID: chn.ID, UserID: usr.ID})
				assert(t, err == nil, "AddChannelMember error: %v", err)
			}
		}
	}
}
