package messaging

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/messaging/internal/repository"
	"github.com/cortezaproject/corteza-server/messaging/internal/service"
	"github.com/cortezaproject/corteza-server/messaging/types"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

func accessControlSetup(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	c.InitServices(ctx, c)

	// Calling grant directly on internal permissions service to avoid AC check for "grant"
	var p = service.DefaultPermissions
	var ac = service.DefaultAccessControl
	return p.Grant(ctx, ac.Whitelist(), ac.DefaultRules()...)
}

// Add default channels when there are none
func makeDefaultChannels(ctx context.Context, cmd *cobra.Command, c *cli.Config) error {
	db, err := factory.Database.Get(messaging)
	if err != nil {
		return err
	}

	repo := repository.Channel(ctx, db)

	cc, err := repo.Find(nil)
	if err != nil {
		return err
	}

	if len(cc) == 0 {
		cc = types.ChannelSet{
			&types.Channel{Name: "General"},
			&types.Channel{Name: "Random"},
		}

		err = cc.Walk(func(c *types.Channel) error {
			_, err := repo.Create(c)
			return err
		})

		if err != nil {
			return err
		}
	}

	return nil
}
