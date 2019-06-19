package system

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/internal/service"
)

func accessControlSetup(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	// Calling grant directly on internal permissions service to avoid AC check for "grant"
	var p = service.DefaultPermissions
	var ac = service.DefaultAccessControl
	return p.Grant(ctx, ac.Whitelist(), ac.DefaultRules()...)
}
