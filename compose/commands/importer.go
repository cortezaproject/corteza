package commands

import (
	"context"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/compose/importer"
	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
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
				aux    interface{}
				ff     []io.Reader
				nsFlag = cmd.Flags().Lookup("namespace").Value.String()
				ns     *types.Namespace
				err    error
			)

			if nsFlag != "" {
				if namespaceID, _ := strconv.ParseUint(nsFlag, 10, 64); namespaceID > 0 {
					ns, err = service.DefaultNamespace.FindByID(namespaceID)
					if err != repository.ErrNamespaceNotFound {
						cli.HandleError(err)
					}
				} else if ns, err = service.DefaultNamespace.FindByHandle(nsFlag); err != nil {
					if err != repository.ErrNamespaceNotFound {
						cli.HandleError(err)
					}
				}
			}

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

			// Initialize importer
			imp := importer.NewImporter(
				service.DefaultNamespace.With(ctx),
				service.DefaultModule.With(ctx),
				service.DefaultChart.With(ctx),
				service.DefaultPage.With(ctx),
				service.DefaultInternalAutomationManager,
				permissions.NewImporter(service.DefaultAccessControl.Whitelist()),
			)

			for i, f := range ff {
				cmd.Printf("Importing from %s\n", args[i])

				cli.HandleError(yaml.NewDecoder(f).Decode(&aux))

				if ns != nil {
					// If we're importing with --namespace switch,
					// we're going to import all into one NS

					cli.HandleError(imp.GetNamespaceImporter().Cast(ns.Slug, aux))
				} else {
					// importing one or more namespaces
					cli.HandleError(imp.Cast(aux))
				}
			}

			// Store all imported
			cli.HandleError(imp.Store(
				ctx,
				service.DefaultNamespace.With(ctx),
				service.DefaultModule.With(ctx),
				service.DefaultChart.With(ctx),
				service.DefaultPage.With(ctx),
				service.DefaultInternalAutomationManager,
				service.DefaultAccessControl,
			))
		},
	}

	cmd.Flags().String("namespace", "", "Import into namespace (by ID or string)")

	return cmd
}
