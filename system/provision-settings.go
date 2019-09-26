package system

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/service"
)

func settingsAutoDiscovery(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	return service.DefaultSettings.With(ctx).AutoDiscovery()
}
