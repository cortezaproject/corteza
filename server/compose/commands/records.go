package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/cli"

	"github.com/spf13/cobra"
)

func Records(ctx context.Context, app serviceInitializer) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     "records",
		Aliases: []string{"rec", "record"},
	}

	cmd.AddCommand(RecordsSynthetic(ctx, app))

	return
}

func RecordsSynthetic(ctx context.Context, app serviceInitializer) *cobra.Command {
	var (
		total uint
		faker = gofakeit.NewCrypto()

		namespace string
		module    string

		synth = &cobra.Command{
			Use:     "synthetic",
			Aliases: []string{"synth"},

			PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
				if err = app.InitServices(ctx); err != nil {
					return
				}

				if err = service.DefaultModule.ReloadDALModels(ctx); err != nil {
					return
				}

				return
			},
		}

		gen = &cobra.Command{
			Use:     "generate",
			Aliases: []string{"gen"},
			Short:   "Create synthetic records",
			Args:    cobra.MaximumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if len(namespace) == 0 || len(module) == 0 {
					cli.HandleError(fmt.Errorf("specifiy ID and handle for both, module and namespace"))
				}

				ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())
				ns, mod, err := resolveModule(ctx, service.DefaultNamespace, service.DefaultModule, namespace, module)
				cli.HandleError(err)

				_ = faker
				_ = ns
				_ = mod
				_ = err
				cmd.Printf("Generating %d compose records (module: %s) ...", total, mod.Name)
				bm := time.Now()

				cli.HandleError(service.DefaultRecord.CreateSynthetic(ctx, faker, mod, total))

				cmd.Printf("done in %s", time.Since(bm).Round(time.Millisecond))
				cmd.Println()
			},
		}

		rem = &cobra.Command{
			Use:     "remove",
			Aliases: []string{"rm", "d", "delete", "del"},
			Short:   "Remove synthetic records",
			Args:    cobra.MaximumNArgs(0),
			Run: func(cmd *cobra.Command, args []string) {
				if len(namespace) == 0 || len(module) == 0 {
					cli.HandleError(fmt.Errorf("specifiy ID and handle for both, module and namespace"))
				}

				ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())
				ns, mod, err := resolveModule(ctx, service.DefaultNamespace, service.DefaultModule, namespace, module)
				cli.HandleError(err)

				_ = ns
				_ = mod
				_ = err

				cmd.Printf("Removing all synthetic compose records (module: %s) ...", mod.Name)
				bm := time.Now()

				cli.HandleError(service.DefaultRecord.RemoveSynthetic(ctx, mod))

				cmd.Printf("done in %s", time.Since(bm).Round(time.Millisecond))
				cmd.Println()
			},
		}
	)

	for _, cmd := range []*cobra.Command{gen, rem} {
		cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace IS or handle for recode generation")
		cmd.Flags().StringVarP(&module, "module", "m", "", "module IS or handle for recode generation")
	}

	gen.Flags().UintVarP(&total, "total", "t", 1, "Number of synthetic records generated")

	synth.AddCommand(gen, rem)
	return synth
}

func resolveModule(ctx context.Context, nsSvc service.NamespaceService, modSvc service.ModuleService, nsIdent, modIdent string) (ns *types.Namespace, mod *types.Module, err error) {
	if ns, err = nsSvc.FindByAny(ctx, nsIdent); err != nil {
		return
	}

	mod, err = modSvc.FindByAny(ctx, ns.ID, modIdent)
	return
}
