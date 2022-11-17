package commands

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/system/service"
	"time"

	"github.com/spf13/cobra"
)

func GenerateSyntheticUsers(ctx context.Context, app serviceInitializer) *cobra.Command {
	var (
		total uint
		faker = gofakeit.NewCrypto()

		base = &cobra.Command{
			Use: "users",
			PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
				return app.InitServices(ctx)
			},
		}

		gen = &cobra.Command{
			Use:     "generate",
			Aliases: []string{"gen"},
			Short:   "Generate synthetic users",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Printf("Generating %d users ...", total)
				bm := time.Now()

				ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())
				cli.HandleError(service.DefaultUser.CreateSynthetic(ctx, faker, total))

				cmd.Printf("done in %s", time.Since(bm).Round(time.Millisecond))
				cmd.Println()
			},
		}

		rem = &cobra.Command{
			Use:     "remove",
			Aliases: []string{"rm", "d", "delete", "del"},
			Short:   "Remove synthetic users",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Printf("Removing all synthetic users ...")
				bm := time.Now()

				ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())
				cli.HandleError(service.DefaultUser.RemoveSynthetic(ctx))

				cmd.Printf("done in %s", time.Since(bm).Round(time.Millisecond))
				cmd.Println()
			},
		}
	)

	gen.Flags().UintVarP(&total, "total", "t", 1, "Number of synthetic users generated")

	base.AddCommand(gen, rem)

	return base
}
