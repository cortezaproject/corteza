// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/messaging/types"
)

func TestChannel(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Channel(context.Background(), factory.Database.MustGet())
	chn := &types.Channel{}

	var name1, name2 = "Test channel v1", "Test channel v2"

	var cc []*types.Channel

	{
		chn.Name = name1
		chn, err = rpo.Create(chn)
		test.Assert(t, err == nil, "CreateChannel error: %+v", err)
		test.Assert(t, chn.Name == name1, "Changes were not stored")

		{
			chn.Name = name2

			chn, err = rpo.Update(chn)
			test.Assert(t, err == nil, "UpdateChannel error: %+v", err)
			test.Assert(t, chn.Name == name2, "Changes were not stored")
		}

		{
			chn, err = rpo.FindByID(chn.ID)
			test.Assert(t, err == nil, "FindByID error: %+v", err)
			test.Assert(t, chn.Name == name2, "Changes were not stored")
		}

		{
			cc, err = rpo.Find(&types.ChannelFilter{Query: name2})
			test.Assert(t, err == nil, "FindChannels error: %+v", err)
			test.Assert(t, len(cc) > 0, "No results found")
		}

		{
			err = rpo.ArchiveByID(chn.ID)
			test.Assert(t, err == nil, "ArchiveByID error: %+v", err)
		}

		{
			err = rpo.UnarchiveByID(chn.ID)
			test.Assert(t, err == nil, "UnarchiveByID error: %+v", err)
		}

		{
			err = rpo.DeleteByID(chn.ID)
			test.Assert(t, err == nil, "DeleteByID error: %+v", err)
		}

		{
			err = rpo.UndeleteByID(chn.ID)
			test.Assert(t, err == nil, "UndeleteByID error: %+v", err)
		}
	}
}
