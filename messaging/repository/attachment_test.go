package repository

import (
	"context"

	"github.com/titpetric/factory"

	"testing"

	"github.com/crusttech/crust/messaging/types"
)

func TestAttachment(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Attachment(context.Background(), factory.Database.MustGet())
	att := &types.Attachment{}

	att.UserID = 1

	{
		att, err = rpo.CreateAttachment(att)
		assert(t, err == nil, "CreateAttachment error: %v", err)
		assert(t, att.UserID == 1, "Changes were not stored")

		{
			att, err = rpo.FindAttachmentByID(att.ID)
			assert(t, err == nil, "FindAttachmentByID error: %v", err)
			assert(t, att.UserID == 1, "Changes were not stored")
		}

		{
			att, err = rpo.FindAttachmentByID(att.ID)
			assert(t, err == nil, "FindAttachmentByMessageID error: %v", err)
			assert(t, att != nil, "No results found")
		}

		{
			err = rpo.DeleteAttachmentByID(att.ID)
			assert(t, err == nil, "DeleteAttachmentByID error: %v", err)
		}
	}
}
