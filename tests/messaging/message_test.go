package messaging

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

func (h helper) repoMessage() repository.MessageRepository {
	return repository.Message(context.Background(), db())
}
func (h helper) repoMessageFlag() repository.MessageFlagRepository {
	return repository.MessageFlag(context.Background(), db())
}

func (h helper) repoMakeMessage(msg string, ch *types.Channel, u *sysTypes.User) *types.Message {
	m, err := h.repoMessage().Create(&types.Message{
		Message:   msg,
		ChannelID: ch.ID,
		UserID:    u.ID,
	})

	h.a.NoError(err)
	return m
}

func (h helper) repoMsgExistingLoad(ID uint64) *types.Message {
	m, err := h.repoMessage().FindByID(ID)
	h.a.NoError(err)
	h.a.NotNil(m)
	return m
}

func (h helper) repoMsgFlagLoad(ID uint64) types.MessageFlagSet {
	ff, err := h.repoMessageFlag().FindByMessageIDs(ID)
	h.a.NoError(err)
	return ff
}
