package service

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
)

type goneAvatar struct {
	AttachmentService
	created bool
}

func (goneAvatar) FindByID(context.Context, uint64) (*types.Attachment, error) {
	return nil, store.ErrNotFound
}

func (a *goneAvatar) CreateAvatarInitialsAttachment(context.Context, string, string, string) (*types.Attachment, error) {
	a.created = true
	att := &types.Attachment{ID: 99}
	att.Meta.Original.Image = &types.AttachmentImageMeta{}
	return att, nil
}

// A user whose avatar attachment is gone can still be saved: the avatar is
// made again rather than the save failing on the lookup.
func TestAvatarInitialReplacesAMissingAttachment(t *testing.T) {
	att := &goneAvatar{}
	svc := user{services: &userServices{att: att}}
	u := &types.User{Name: "Plain Home", Meta: &types.UserMeta{AvatarID: 42}}

	require.NoError(t, svc.generateUserAvatarInitial(context.Background(), u))
	require.True(t, att.created)
	require.Equal(t, uint64(99), u.Meta.AvatarID)
}
