package commands

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/compose/importer"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/pkg/cli"
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
			} else {
				args = []string{"STDIN"}
				ff = []io.Reader{os.Stdin}
			}

			for i, f := range ff {
				cmd.Printf("Importing from %s\n", args[i])

				imp := importer.NewImporter(
					service.DefaultNamespace.With(ctx),
					service.DefaultModule.With(ctx),
					service.DefaultChart.With(ctx),
					service.DefaultPage.With(ctx),
					permissions.NewImporter(service.DefaultAccessControl.Whitelist()),
				)

				cli.HandleError(imp.YAML(f))
				cli.HandleError(imp.Store(
					ctx,
					service.DefaultNamespace.With(ctx),
					service.DefaultModule.With(ctx),
					service.DefaultChart.With(ctx),
					service.DefaultPage.With(ctx),
					service.DefaultAccessControl,
				))
			}
		},
	}

	return cmd
}
