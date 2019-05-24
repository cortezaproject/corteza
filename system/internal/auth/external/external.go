package external

import (
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/settings"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func Init(settingsService settings.Service) {
	if eas, err := ExternalAuthSettings(settingsService); err != nil {
		log().Error("failed load external authentication settings", zap.Error(err))
	} else {
		setupGoth(eas)
	}
}

func log() *zap.Logger {
	return logger.Default().Named("auth.external")
}
