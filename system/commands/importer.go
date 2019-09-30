package commands

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/importer"
)

func Importer(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import",

		Run: func(cmd *cobra.Command, args []string) {

			c.InitServices(ctx, c)

			var (
				ff  []io.Reader
				err error
			)

			ctx = auth.SetSuperUserContext(ctx)

			if len(args) > 0 {
				ff = make([]io.Reader, len(args))
				for a, arg := range args {
					ff[a], err = os.Open(arg)
					cli.HandleError(err)
				}
				cli.HandleError(importer.Import(ctx, ff...))
			} else {
				cli.HandleError(importer.Import(ctx, os.Stdin))
			}
		},
	}

	return cmd
}
