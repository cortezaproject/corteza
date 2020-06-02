package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	command struct {
		ctx context.Context
	}

	CommandService interface {
		With(context.Context) CommandService
		Do(channelID uint64, command, input string) (*types.Message, error)
	}
)

func Command(ctx context.Context) CommandService {
	return (&command{}).With(ctx)
}

func (svc command) With(ctx context.Context) CommandService {
	return &command{
		ctx: ctx,
	}
}

func (svc command) Do(channelID uint64, command, input string) (*types.Message, error) {
	switch command {
	case "me":
		if input != "" {
			return DefaultMessage.With(svc.ctx).Create(&types.Message{
				Type:      types.MessageTypeIlleism,
				ChannelID: channelID,
				Message:   input,
			})
		}
		return nil, nil
	case "tableflip":
		fallthrough
	case "unflip":
		fallthrough
	case "shrug":
		messages := map[string]string{
			"tableflip": `(╯°□°）╯︵ ┻━┻`,
			"unflip":    `┬─┬ ノ( ゜-゜ノ)`,
			"shrug":     `¯\\_(ツ)_/¯`,
		}
		msg := &types.Message{
			ChannelID: channelID,
			Message:   messages[command],
		}

		if input != "" {
			msg.Message = input + " " + msg.Message
		}
		return DefaultMessage.With(svc.ctx).Create(msg)
	}

	return nil, nil
}
