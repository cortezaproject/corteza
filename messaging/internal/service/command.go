package service

import (
	"context"

	"github.com/crusttech/crust/messaging/types"
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

func (svc *command) With(ctx context.Context) CommandService {
	return &command{
		ctx: ctx,
	}
}

func (svc *command) Do(channelID uint64, command, input string) (*types.Message, error) {
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
	default:
		webhookSvc := DefaultWebhook.With(svc.ctx)
		webhooks, err := webhookSvc.Find(&types.WebhookFilter{
			ChannelID:       channelID,
			OutgoingTrigger: command,
		})
		if err != nil || len(webhooks) == 0 {
			return nil, err
		}
		return webhookSvc.Do(webhooks[0], input)
	}
	return nil, ErrUnknownCommand.new()
}
