package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/crusttech/human/server/compose/service"
	composeTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/auth"
	"github.com/crusttech/human/server/pkg/cli"
	"github.com/crusttech/human/server/pkg/filter"
	"github.com/crusttech/human/server/pkg/options"
	"github.com/crusttech/human/server/store"
)

func Namespaces(ctx context.Context, app serviceInitializer) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:     "namespaces",
		Aliases: []string{"ns", "namespace"},
	}

	cmd.AddCommand(NamespacesPurge(ctx, app))

	return
}

// NamespacesPurge hard-deletes SOFT-DELETED namespaces (matched by slug
// prefix) together with their module/field/page/layout/chart store rows.
//
// Development hygiene tool: repeated create/delete test cycles accumulate
// soft-deleted rows that share slugs with live data and pollute lookups.
// Refused in production environments.
func NamespacesPurge(ctx context.Context, app serviceInitializer) *cobra.Command {
	var (
		prefix string
	)

	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Hard-delete SOFT-DELETED namespaces (and their compose resources) by slug prefix — development only",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.InitServices(ctx)
		},

		Run: func(cmd *cobra.Command, args []string) {
			if options.Environment().IsProduction() {
				cli.HandleError(fmt.Errorf("purge is refused when ENVIRONMENT=production"))
			}
			if strings.TrimSpace(prefix) == "" {
				cli.HandleError(fmt.Errorf("a non-empty --prefix is required"))
			}

			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())
			s := service.DefaultStore

			nn, _, err := store.SearchComposeNamespaces(ctx, s, composeTypes.NamespaceFilter{
				Deleted: filter.StateExclusive,
			})
			cli.HandleError(err)

			purged := 0
			for _, ns := range nn {
				if !strings.HasPrefix(ns.Slug, prefix) {
					continue
				}

				err = store.Tx(ctx, s, func(ctx context.Context, s store.Storer) (err error) {
					// children are often not soft-deleted themselves, so
					// include everything under the namespace
					mm, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{
						NamespaceID: ns.ID, Deleted: filter.StateInclusive,
					})
					if err != nil {
						return
					}
					for _, m := range mm {
						ff, _, err := store.SearchComposeModuleFields(ctx, s, composeTypes.ModuleFieldFilter{
							ModuleID: []uint64{m.ID}, Deleted: filter.StateInclusive,
						})
						if err != nil {
							return err
						}
						for _, f := range ff {
							if err = store.DeleteComposeModuleFieldByID(ctx, s, f.ID); err != nil {
								return err
							}
						}
						if err = store.DeleteComposeModuleByID(ctx, s, m.ID); err != nil {
							return err
						}
					}

					ll, _, err := store.SearchComposePageLayouts(ctx, s, composeTypes.PageLayoutFilter{
						NamespaceID: ns.ID, Deleted: filter.StateInclusive,
					})
					if err != nil {
						return
					}
					for _, l := range ll {
						if err = store.DeleteComposePageLayoutByID(ctx, s, l.ID); err != nil {
							return err
						}
					}

					pp, _, err := store.SearchComposePages(ctx, s, composeTypes.PageFilter{
						NamespaceID: ns.ID, Deleted: filter.StateInclusive,
					})
					if err != nil {
						return
					}
					for _, p := range pp {
						if err = store.DeleteComposePageByID(ctx, s, p.ID); err != nil {
							return err
						}
					}

					cc, _, err := store.SearchComposeCharts(ctx, s, composeTypes.ChartFilter{
						NamespaceID: ns.ID, Deleted: filter.StateInclusive,
					})
					if err != nil {
						return
					}
					for _, c := range cc {
						if err = store.DeleteComposeChartByID(ctx, s, c.ID); err != nil {
							return err
						}
					}

					return store.DeleteComposeNamespaceByID(ctx, s, ns.ID)
				})
				cli.HandleError(err)

				cmd.Printf("purged namespace %s (ID %d)\n", ns.Slug, ns.ID)
				purged++
			}

			cmd.Printf("purged %d soft-deleted namespace(s) with slug prefix %q\n", purged, prefix)
			if purged > 0 {
				cmd.Println("note: record rows in DAL storage are not purged (unreachable once their module rows are gone)")
			}
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "agent-", "only purge namespaces whose slug starts with this prefix")

	return cmd
}
