package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/messaging/types"
)

type (
	command struct {
		ctx    context.Context
		logger *zap.Logger
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
		ctx:    ctx,
		logger: DefaultLogger.Named("command"),
	}
}

// log() returns zap's logger with requestID from current context and fields.
// func (svc command) log(fields ...zapcore.Field) *zap.Logger {
// 	return logger.AddRequestID(svc.ctx, svc.logger).With(fields...)
// }

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
}
