package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestAttachment(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Attachment()
	ctx := context.Background()
	att := &types.Attachment{}

	var aa []*types.Attachment

	att.ChannelID = 1

	att, err = rpo.CreateAttachment(ctx, att)
	must(t, err)
	if att.ChannelID != 1 {
		t.Fatal("Changes were not stored")
	}

	att.ChannelID = 2

	att, err = rpo.UpdateAttachment(ctx, att)
	must(t, err)
	if att.ChannelID != 2 {
		t.Fatal("Changes were not stored")
	}

	att, err = rpo.FindAttachmentByID(ctx, att.ID)
	must(t, err)
	if att.ChannelID != 2 {
		t.Fatal("Changes were not stored")
	}

	aa, err = rpo.FindAttachmentByRange(ctx, 2, 0, att.ID)
	must(t, err)
	if len(aa) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.DeleteAttachmentByID(ctx, att.ID))
}
