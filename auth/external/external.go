package external

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

const (
	OIDC_PROVIDER_PREFIX = "openid-connect."
)

func Init(store sessions.Store) {
	gothic.Store = store
}
