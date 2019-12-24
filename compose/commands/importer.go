package commands

import (
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/compose/importer"
	"github.com/cortezaproject/corteza-server/compose/repository"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
)

func Importer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import",

		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx    = auth.SetSuperUserContext(cli.Context())
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

			if len(args) > 0 {
				ff = make([]io.Reader, len(args))
				for a, arg := range args {
					ff[a], err = os.Open(arg)
					cli.HandleError(err)
				}
				cli.HandleError(importer.Import(ctx, ns, ff...))
			} else {
				cli.HandleError(importer.Import(ctx, ns, os.Stdin))
			}
		},
	}

	cmd.Flags().String("namespace", "", "Import into namespace (by ID or string)")

	return cmd
}
