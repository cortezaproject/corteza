package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"net/http"
	"strconv"
)

type (
	stdMessageResponse struct {
		Response struct {
			MessageID string `json:"messageID"`
		} `json:"response"`
	}
)

func (h helper) makeMessage(msg string, ch *types.Channel, u *sysTypes.User) *types.Message {
	m := &types.Message{
		ID:        id.Next(),
		Message:   msg,
		ChannelID: ch.ID,
		UserID:    u.ID,
	}
	err := store.CreateMessagingMessage(context.Background(), service.DefaultStore, m)

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

func (h helper) lookupMessageByID(ID uint64) (*types.Message, error) {
	return store.LookupMessagingMessageByID(context.Background(), service.DefaultStore, ID)
}

func (h helper) lookupFlagByMessageID(ID uint64) types.MessageFlagSet {
	ff, _, err := store.SearchMessagingFlags(context.Background(), service.DefaultStore, types.MessageFlagFilter{MessageID: []uint64{ID}})
	h.a.NoError(err)
	return ff
}
