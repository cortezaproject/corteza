package outgoing

import (
	"github.com/crusttech/crust/sam/types"
	"strconv"
	"time"
)

type WsMessage struct {
	Error *Error `json:"error,omitempty"`

	*Message `json:"m"`

	// @todo: implement outgoing message types
	timestamp time.Time
}

//func (WsMessage) New() *WsMessage {
//	return &WsMessage{
//		//id:        factory.Sonyflake.NextID(),
//		timestamp: time.Now().UTC(),
//	}
//}

func NewError(err error) *WsMessage {
	return &WsMessage{Error: &Error{Message: err.Error()}}
}

func FromMessage(msg *types.Message) *WsMessage {
	return &WsMessage{Message: &Message{
		Message:   msg.Message,
		Id:        strconv.FormatUint(msg.ID, 10),
		ChannelId: strconv.FormatUint(msg.ChannelId, 10),
		Type:      msg.Type,
		UserId:    strconv.FormatUint(msg.UserId, 10),
		ReplyTo:   strconv.FormatUint(msg.ReplyTo, 10),
	}}
}
