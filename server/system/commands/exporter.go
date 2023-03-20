package commands

import (
	"context"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	storeInitFnc func(ctx context.Context) (store.Storer, error)
	dalInitFnc   func(ctx context.Context) (dal.FullService, error)
	envoyInitFnc func(ctx context.Context) (*envoyx.Service, error)
)

var (
	omitModules bool
	omitPages   bool
	omitCharts  bool

	inclLocale bool
	inclRbac   bool
)

func Export(ctx context.Context, storeInit storeInitFnc, dalInit dalInitFnc, envoyInit envoyInitFnc) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export",
		Long:  `Export data`,
	}

	exportNsCommand := &cobra.Command{
		Use:   "compose-namespace",
		Short: "Export compose namespace",
		Args:  cobra.MinimumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

			// Init all of the sub services
			s, err := storeInit(ctx)
			cli.HandleError(err)

			dalSvc, err := dalInit(ctx)
			cli.HandleError(err)

			envoySvc, err := envoyInit(ctx)
			cli.HandleError(err)

			scp := envoyx.Scope{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  envoyx.MakeIdentifiers(args[0]),
			}

			f := map[string]envoyx.ResourceFilter{
				types.NamespaceResourceType: {
					Identifiers: envoyx.MakeIdentifiers(args[0]),
					Scope:       scp,
				},
			}

			if !omitModules {
				f[types.ModuleResourceType] = envoyx.ResourceFilter{
					Scope: scp,
					Refs: map[string]envoyx.Ref{
						"NamespaceID": {
							ResourceType: types.NamespaceResourceType,
							Identifiers:  envoyx.MakeIdentifiers(args[0]),
							Scope:        scp,
						},
					},
				}
			}

			if !omitPages {
				f[types.PageResourceType] = envoyx.ResourceFilter{
					Scope: scp,
					Refs: map[string]envoyx.Ref{
						"NamespaceID": {
							ResourceType: types.NamespaceResourceType,
							Identifiers:  envoyx.MakeIdentifiers(args[0]),
							Scope:        scp,
						},
					},
				}
			}

			if !omitCharts {
				f[types.ChartResourceType] = envoyx.ResourceFilter{
					Scope: scp,
					Refs: map[string]envoyx.Ref{
						"NamespaceID": {
							ResourceType: types.NamespaceResourceType,
							Identifiers:  envoyx.MakeIdentifiers(args[0]),
							Scope:        scp,
						},
					},
				}
			}

			nodes, _, err := envoySvc.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeStore,
				Params: map[string]any{
					"storer": s,
					"dal":    dalSvc,
				},
				Filter: f,
			})
			cli.HandleError(err)

			if inclRbac {
				var aux envoyx.NodeSet
				aux, err = encodeRbacRules(ctx, s, nodes)
				cli.HandleError(err)

				nodes = append(nodes, aux...)
			}

			if inclLocale {
				var aux envoyx.NodeSet
				aux, err = encodeLocale(ctx, s, nodes)
				cli.HandleError(err)
				nodes = append(nodes, aux...)
			}

			rf := map[string]envoyx.ResourceFilter{
				systemTypes.RoleResourceType: {},
			}
			refNodes, _, err := envoySvc.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeStore,
				Params: map[string]any{
					"storer": s,
					"dal":    dalSvc,
				},
				Filter: rf,
			})
			cli.HandleError(err)
			for _, n := range refNodes {
				n.Placeholder = true
			}

			nodes = append(nodes, refNodes...)

			gg, err := envoySvc.Bake(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeStore,
				Params: map[string]any{
					"storer": s,
					"dal":    dalSvc,
				},
			}, nil, nodes...)
			cli.HandleError(err)

			err = envoySvc.Encode(ctx, envoyx.EncodeParams{
				Type: envoyx.EncodeTypeIo,
				Params: map[string]any{
					"writer": cmd.OutOrStdout(),
				},
			}, gg)
			cli.HandleError(err)
		},
	}

	exportNsCommand.Flags().BoolVar(
		&omitModules,
		"omit-modules",
		false,
		"Omit modules from output",
	)
	exportNsCommand.Flags().BoolVar(
		&omitPages,
		"omit-pages",
		false,
		"Omit pages from output",
	)
	exportNsCommand.Flags().BoolVar(
		&omitCharts,
		"omit-charts",
		false,
		"Omit charts from output",
	)

	exportNsCommand.Flags().BoolVar(
		&inclLocale,
		"include-locale",
		false,
		"Include resource translations",
	)
	exportNsCommand.Flags().BoolVar(
		&inclRbac,
		"include-rbac",
		false,
		"Include RBAC rules",
	)

	cmd.AddCommand(
		exportNsCommand,
	)

	return cmd
}

func encodeRbacRules(ctx context.Context, s store.Storer, nn envoyx.NodeSet) (rules envoyx.NodeSet, err error) {
	rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
	if err != nil {
		return
	}

	rules, err = envoyx.RBACRulesForNodes(rr, nn...)
	if err != nil {
		return
	}

	return
}

func encodeLocale(ctx context.Context, s store.Storer, nn envoyx.NodeSet) (translations envoyx.NodeSet, err error) {
	tt, _, err := store.SearchResourceTranslations(ctx, s, systemTypes.ResourceTranslationFilter{})
	if err != nil {
		return
	}

	translations, err = envoyx.ResourceTranslationsForNodes(tt, nn...)
	if err != nil {
		return
	}

	return
}
