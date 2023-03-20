package commands

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
)

func Import(ctx context.Context, storeInit storeInitFnc, dalInit dalInitFnc, envoyInit envoyInitFnc) *cobra.Command {
	var (
		replaceOnConflict bool
		skipOnConflict    bool
		panicOnConflict   bool
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import data from yaml sources.",

		Run: func(cmd *cobra.Command, args []string) {
			// Init all of the sub services
			s, err := storeInit(ctx)
			cli.HandleError(err)

			dalSvc, err := dalInit(ctx)
			cli.HandleError(err)

			envoySvc, err := envoyInit(ctx)
			cli.HandleError(err)

			var (
				nodes     envoyx.NodeSet
				providers []envoyx.Provider

				auxNodes     envoyx.NodeSet
				auxProviders []envoyx.Provider
			)

			if len(args) > 0 {
				for _, a := range args {
					auxNodes, auxProviders, err = envoySvc.Decode(ctx, envoyx.DecodeParams{
						Type: envoyx.DecodeTypeURI,
						Params: map[string]any{
							"uri": "file://" + a,
						},
					})
					cli.HandleError(err)

					nodes = append(nodes, auxNodes...)
					providers = append(providers, auxProviders...)
				}
			} else {
				auxNodes, auxProviders, err = envoySvc.Decode(ctx, envoyx.DecodeParams{
					Type: envoyx.DecodeTypeIO,
					Params: map[string]any{
						"reader": os.Stdin,
						"mime":   "text/yaml",
					},
				})
				cli.HandleError(err)

				// @todo consider changing this up
				if len(auxProviders) > 0 {
					cli.HandleError(errors.New("cannot define providers when importing from stdin"))
				}

				nodes = append(nodes, auxNodes...)
				// providers = append(providers, auxProviders...)
			}

			ep := envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": s,
					"dal":    dalSvc,
				},
			}

			if replaceOnConflict {
				ep.Envoy.MergeAlg = envoyx.OnConflictReplace
			}
			if skipOnConflict {
				ep.Envoy.MergeAlg = envoyx.OnConflictSkip
			}
			if panicOnConflict {
				ep.Envoy.MergeAlg = envoyx.OnConflictPanic
			}

			gg, err := envoySvc.Bake(ctx, ep,
				providers,
				nodes...,
			)
			cli.HandleError(err)

			err = envoySvc.Encode(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": s,
					"dal":    dalSvc,
				},
			}, gg)
			cli.HandleError(err)
		},
	}

	cmd.Flags().BoolVar(
		&replaceOnConflict,
		"replace-existing",
		false,
		"replace on conflict existing resources. Default skips",
	)

	cmd.Flags().BoolVar(
		&skipOnConflict,
		"skip-existing",
		false,
		"skip on conflict existing resources. Default skips",
	)

	cmd.Flags().BoolVar(
		&panicOnConflict,
		"panic-existing",
		false,
		"panic on conflict existing resources. Default skips",
	)

	return cmd
}
