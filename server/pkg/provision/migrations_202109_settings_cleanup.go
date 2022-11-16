package provision

import (
	"context"

	"github.com/cortezaproject/corteza/server/store"
	"go.uber.org/zap"
)

func cleanupPre202109Settings(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	log.Info("cleaning up pre-2021.9 settings")
	names := []string{
		"compose.ui.namespace-switcher.enabled",
		"compose.ui.namespace-switcher.default-open",
		"messaging.message.attachments.enabled",
		"messaging.message.attachments.max-size",
		"messaging.message.attachments.mimetypes",
		"messaging.message.attachments.source.camera.enabled",
		"messaging.message.attachments.source.gallery.enabled",
		"messaging.ui.browser-notifications.enabled",
		"messaging.ui.browser-notifications.header",
		"messaging.ui.browser-notifications.message-trim",
		"messaging.ui.emoji.enabled",
		"auth.external.session-store-secret",
		"auth.external.session-store-secure",
	}

	for _, name := range names {
		if err = store.DeleteSettingByNameOwnedBy(ctx, s, name, 0); err != nil {
			return
		}
	}

	return
}
