package repository

import (
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestAttachment(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := New()
	att := &types.Attachment{}

	var aa []*types.Attachment

	att.ChannelID = 1

	{
		att, err = rpo.CreateAttachment(att)
		assert(t, err == nil, "CreateAttachment error: %v", err)
		assert(t, att.ChannelID == 1, "Changes were not stored")

		{
			att.ChannelID = 2

			att, err = rpo.UpdateAttachment(att)
			assert(t, err == nil, "UpdateAttachment error: %v", err)
			assert(t, att.ChannelID == 2, "Changes were not stored")
		}

		{
			att, err = rpo.FindAttachmentByID(att.ID)
			assert(t, err == nil, "FindAttachmentByID error: %v", err)
			assert(t, att.ChannelID == 2, "Changes were not stored")
		}

		{
			aa, err = rpo.FindAttachmentByRange(2, 0, att.ID)
			assert(t, err == nil, "FindAttachmentByRange error: %v", err)
			assert(t, len(aa) > 0, "No results found")
		}

		{
			err = rpo.DeleteAttachmentByID(att.ID)
			assert(t, err == nil, "DeleteAttachmentByID error: %v", err)
		}
	}
}
