package messaging

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/importer"
	"github.com/cortezaproject/corteza-server/messaging/types"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

// Provision only where there are no channels
func hasChannels(ctx context.Context, s store.Storer) (bool, error) {
	if set, _, err := store.SearchMessagingChannels(ctx, s, types.ChannelFilter{}); err != nil {
		return false, err
	} else {
		return len(set) > 0, nil
	}
}

func Provision(ctx context.Context, log *zap.Logger, s store.Storer) error {
	if channelsExist, err := hasChannels(ctx, s); err != nil {
		return err
	} else if !channelsExist {
		log.Info("provisioning messaging")
		readers, err := impAux.ReadStatic(Asset)
		if err != nil {
			return err
		}

		return importer.Import(ctx, readers...)
	} else {
		log.Info("messaging already provisioned")
	}

	return nil
}
