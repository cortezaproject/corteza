// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/test"
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
		test.Assert(t, err == nil, "CreateAttachment error: %+v", err)
		test.Assert(t, att.UserID == 1, "Changes were not stored")

		{
			att, err = rpo.FindAttachmentByID(att.ID)
			test.Assert(t, err == nil, "FindAttachmentByID error: %+v", err)
			test.Assert(t, att.UserID == 1, "Changes were not stored")
		}

		{
			att, err = rpo.FindAttachmentByID(att.ID)
			test.Assert(t, err == nil, "FindAttachmentByMessageID error: %+v", err)
			test.Assert(t, att != nil, "No results found")
		}

		{
			err = rpo.DeleteAttachmentByID(att.ID)
			test.Assert(t, err == nil, "DeleteAttachmentByID error: %+v", err)
		}
	}
}
