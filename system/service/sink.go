package service

import (
	"context"
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/logger"
)

type (
	sink struct {
		// processors
		proc map[string]sinkContentProc

		logger *zap.Logger
	}

	sinkContentProc interface {
		ContentProcessor(context.Context, io.Reader) error
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
		proc: map[string]sinkContentProc{
			SinkContentTypeMail: Mailproc(),
		},
	}
}

// Finds appropriate sink processor
func (svc *sink) Process(ctx context.Context, contentType string, r io.Reader) (err error) {
	switch strings.ToLower(contentType) {
	case SinkContentTypeMail, "rfc822", "email", "mail":
		if err = svc.proc[SinkContentTypeMail].ContentProcessor(ctx, r); err != nil {
			return ErrSinkContentProcessingFailed
		}

	default:
		return ErrSinkContentTypeUnsupported
	}

	return
}

func (svc sink) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}
