package compose

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/compose/importer"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	provision "github.com/cortezaproject/corteza-server/provision/compose"
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
		importer.Import(ctx, nil, readers...),
		"could not provision configuration for compose service",
	)
}

// Provision only where there are no namespaces
func isProvisioned(ctx context.Context) (bool, error) {
	_, f, err := service.DefaultNamespace.With(ctx).Find(types.NamespaceFilter{})
	return f.Count > 0, err
}
