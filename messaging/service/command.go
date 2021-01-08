package service

import (
	"context"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	command struct {
	}

	CommandService interface {
		Do(ctx context.Context, channelID uint64, command, input string) (*types.Message, error)
	}
)

func Command() CommandService {
	return &command{}
}

func (svc command) Do(ctx context.Context, channelID uint64, command, input string) (*types.Message, error) {
	switch command {
	case "me":
		if input != "" {
			return DefaultMessage.Create(ctx, &types.Message{
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
		return DefaultMessage.Create(ctx, msg)
	}

	return nil, nil
}
