package service

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	sink struct {
		logger   *zap.Logger
		eventbus sinkEventDispatcher
	}

	sinkEventDispatcher interface {
		WaitFor(ctx context.Context, ev eventbus.Event) (err error)
	}
)

const (
	ErrSinkContentTypeUnsupported  serviceError = "SinkUnsupportedContentType"
	ErrSinkContentProcessingFailed serviceError = "SinkProcessFailed"

	SinkContentTypeMail = "message/rfc822"
)

func Sink() *sink {
	return &sink{
		logger: DefaultLogger,
	}
}

// Processes sink request, casts it and forwards it to
func (svc *sink) Process(ctx context.Context, contentType string, r *http.Request) (err error) {
	switch strings.ToLower(contentType) {
	case SinkContentTypeMail, "rfc822", "email", "mail":
		// this is handled by dedicated event that parses raw payload from HTTP request
		// as rfc882 message.
		var msg *types.MailMessage
		msg, err = types.NewMailMessage(r.Body)
		if err != nil {
			return
		}

		return svc.eventbus.WaitFor(ctx, event.MailOnReceive(msg))

	default:
		var sr *types.SinkRequest
		sr, err = types.NewSinkRequest(r)
		if err != nil {
			return
		}

		return svc.eventbus.WaitFor(ctx, event.SinkOnRequest(sr))
	}
}

func (svc sink) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}
