package external

import (
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/internal/service"
)

func Init() {
	setupGoth(service.DefaultAuthSettings)
}

func log() *zap.Logger {
	return logger.Default().Named("auth.external")
}
