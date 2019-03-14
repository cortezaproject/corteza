package websocket

import (
	"context"
	"log"
	"time"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/payload"
	"github.com/crusttech/crust/internal/payload/incoming"
	"github.com/crusttech/crust/internal/payload/outgoing"
	"github.com/crusttech/crust/messaging/types"
	systemService "github.com/crusttech/crust/system/service"
)

func (s *Session) execCommand(ctx context.Context, c *incoming.ExecCommand) error {
	// @todo: check access / can we join this channel (should be done by service layer)

	log.Printf("Received command '%s(%v)", c.Command, c.Params)

	switch c.Command {
	case "echo":
		if c.Input != "" {
			user, err := systemService.User(ctx).FindByID(s.user.Identity())
			if err != nil {
				return err
			}

			return s.sendReply(&outgoing.Message{
				ID:        factory.Sonyflake.NextID(),
				User:      payload.User(user),
				CreatedAt: time.Now(),
				Type:      "hallucination",
				ChannelID: c.ChannelID,
				Message:   c.Input})
		}

	case "shrug":
		msg := &types.Message{
			ChannelID: payload.ParseUInt64(c.ChannelID),
			Message:   "¯\\_(ツ)_/¯",
		}

		if c.Input != "" {
			msg.Message = c.Input + " " + msg.Message
		}
		_, err := s.svc.msg.With(ctx).Create(msg)

		return err
	}

	return nil
}
