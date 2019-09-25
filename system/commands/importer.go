package commands

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/importer"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

func Importer(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import",

		Run: func(cmd *cobra.Command, args []string) {

			c.InitServices(ctx, c)

			var (
				aux interface{}
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

			roles, err := service.DefaultRole.With(ctx).Find(&types.RoleFilter{})
			cli.HandleError(err)

			for i, f := range ff {
				cmd.Printf("Importing from %s\n", args[i])

				cli.HandleError(yaml.NewDecoder(f).Decode(&aux))

				perm := permissions.NewImporter(service.DefaultAccessControl.Whitelist())

				imp := importer.NewImporter(perm,
					importer.NewRoleImport(perm, roles),
				)

				cli.HandleError(imp.Store(
					ctx,
					service.DefaultRole.With(ctx),
					service.DefaultAccessControl,
				))
			}
		},
	}

	return cmd
}
