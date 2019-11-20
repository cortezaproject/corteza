package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cortezaproject/corteza-server/messaging/repository"
	"github.com/cortezaproject/corteza-server/messaging/types"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	stdMessageResponse struct {
		Response struct {
			MessageID string `json:"messageID"`
		} `json:"response"`
	}
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

func (h helper) apiMessageCreateReply(msg string, o *types.Message) *types.Message {
	rsp := h.apiInit().
		Post(fmt.Sprintf("/channels/%d/messages/%d/replies", o.ChannelID, o.ID)).
		JSON(fmt.Sprintf(`{"message": "%s"}`, msg)).
		Expect(h.t).
		Status(http.StatusOK).
		End()

	thrMsg := stdMessageResponse{}
	if err := json.NewDecoder(rsp.Response.Body).Decode(&thrMsg); err != nil {
		h.a.Fail(err.Error())
	}

	msgID, err := strconv.ParseInt(thrMsg.Response.MessageID, 10, 64)
	h.a.Nil(err)
	return &types.Message{
		ID:   uint64(msgID),
		Type: types.MessageTypeSimpleMessage,
	}
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
