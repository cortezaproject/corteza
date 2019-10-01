package messaging

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/messaging/importer"
	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	provision "github.com/cortezaproject/corteza-server/provision/messaging"
)

func provisionConfig(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	c.Log.Debug("running configuration provision")
	c.InitServices(ctx, c)

	// Make sure we have all full access for provisioning
	ctx = auth.SetSuperUserContext(ctx)

	if provisioned, err := isProvisioned(ctx); err != nil {
		return err
	} else if provisioned {
		c.Log.Debug("configuration already provisioned")
		return nil
	}

	readers, err := impAux.ReadStatic(provision.Asset)
	if err != nil {
		return err
	}

	return errors.Wrap(
		importer.Import(ctx, readers...),
		"could not provision configuration for messaging service",
	)
}

// Provision ONLY when there are no channels (even if we find delete channels we abort provisioning
func isProvisioned(ctx context.Context) (bool, error) {
	cc, err := service.DefaultChannel.With(ctx).Find(&types.ChannelFilter{IncludeDeleted: true})
	return len(cc) > 0, err
}
