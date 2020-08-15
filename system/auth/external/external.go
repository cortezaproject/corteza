package external

import (
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/service"
)

const (
	OIDC_PROVIDER_PREFIX = "openid-connect."
)

func Init() {
	ase := service.CurrentSettings.Auth.External

	if !ase.Enabled {
		log().Info("external authentication disabled")
		return
	}

	setupGoth(ase.SessionStoreSecure, []byte(ase.SessionStoreSecret))
	setupGothProviders(ase.Providers, ase.RedirectUrl)

}

func log() *zap.Logger {
	return logger.Default().Named("auth.external")
}
