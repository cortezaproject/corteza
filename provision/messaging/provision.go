package messaging

import (
	"context"
	"github.com/cortezaproject/corteza-server/messaging/importer"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	"go.uber.org/zap"
)

func Provision(ctx context.Context, log *zap.Logger) error {
	if provisioned, err := notProvisioned(ctx); err != nil {
		return err
	} else if !provisioned {
		log.Info("provisioning messaging")
		readers, err := impAux.ReadStatic(Asset)
		if err != nil {
			return err
		}

		return importer.Import(ctx, readers...)
	}

	return nil
}

// Provision only where there are no channels
func notProvisioned(ctx context.Context) (bool, error) {
	cc, _, err := service.DefaultChannel.With(ctx).Find(types.ChannelFilter{IncludeDeleted: true})
	return len(cc) == 0, err
}
