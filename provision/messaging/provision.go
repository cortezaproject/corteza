package messaging

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/provision/util"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

// provision only where there are no channels
func hasChannels(ctx context.Context, s store.Storer) (bool, error) {
	if set, _, err := store.SearchMessagingChannels(ctx, s, types.ChannelFilter{IncludeDeleted: true}); err != nil {
		return false, err
	} else {
		return len(set) > 0, nil
	}
}

func Provision(ctx context.Context, log *zap.Logger, s store.Storer) error {
	log.Info("provisioning messaging")
	if channelsExist, err := hasChannels(ctx, s); err != nil {
		return err
	} else if !channelsExist {
		// Provision from YAML files
		// - access control
		// - channels
		if err = util.EncodeStatik(ctx, s, Asset, "/"); err != nil {
			return err
		}
	}

	log.Info("messaging provisioned")
	return nil
}
