package external

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"go.uber.org/zap"
)

var (
	log = zap.NewNop()
)

const (
	OIDC_PROVIDER_PREFIX = "openid-connect."
)

func Init(logger *zap.Logger, store sessions.Store) {
	log = logger.Named("external")
	gothic.Store = store
}
