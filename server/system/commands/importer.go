package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza/server/compose/service"
	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/csv"
	"github.com/cortezaproject/corteza/server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza/server/pkg/envoy/store"
	"github.com/cortezaproject/corteza/server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza/server/store"
)

func Import(ctx context.Context, storeInit func(ctx context.Context) (store.Storer, error)) *cobra.Command {
	var (
		replaceOnExisting    bool
		mergeLeftOnExisting  bool
		mergeRightOnExisting bool
		defaultResTr         bool
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import data from yaml sources.",

		Run: func(cmd *cobra.Command, args []string) {
			s, err := storeInit(ctx)
			cli.HandleError(err)

			err = service.DalModelReload(ctx, s, dal.Service())
			cli.HandleError(err)

			yd := yaml.Decoder()
			cd := csv.Decoder()
			nn := make([]resource.Interface, 0, 200)

			if len(args) > 0 {
				for _, fn := range args {
					mm, err := directory.Decode(ctx, fn, yd, cd)
					cli.HandleError(err)
					nn = append(nn, mm...)
				}
			} else {
				do := &envoy.DecoderOpts{
					Name: "stdin.yaml",
					Path: "",
				}
				mm, err := yd.Decode(ctx, os.Stdin, do)
				cli.HandleError(err)
				nn = append(nn, mm...)
			}

			if !defaultResTr {
				nn = pruneResTr(nn)
			}

			nn, err = resource.Shape(nn, resource.ComposeRecordShaper())
			if err != nil {
				cli.HandleError(err)
				return
			}

			opt := &es.EncoderConfig{
				OnExisting: resource.Skip,
			}

			if replaceOnExisting {
				opt.OnExisting = resource.Replace
			}
			if mergeLeftOnExisting {
				opt.OnExisting = resource.MergeLeft
			}
			if mergeRightOnExisting {
				opt.OnExisting = resource.MergeRight
			}

			se := es.NewStoreEncoder(s, dal.Service(), opt)
			bld := envoy.NewBuilder(se)
			g, err := bld.Build(ctx, nn...)
			cli.HandleError(err)

			cli.HandleError(envoy.Encode(ctx, g, se))
		},
	}

	cmd.Flags().BoolVar(
		&replaceOnExisting,
		"replace-existing",
		false,
		"Replace any existing values. Default skips.",
	)
	cmd.Flags().BoolVar(
		&mergeLeftOnExisting,
		"merge-left-existing",
		false,
		"Update any existing values; existing data takes priority. Default skips.",
	)
	cmd.Flags().BoolVar(
		&mergeRightOnExisting,
		"merge-right-existing",
		false,
		"Update any existing values; new data takes priority. Default skips.",
	)
	cmd.Flags().BoolVar(
		&defaultResTr,
		"resource-translationsDefaults",
		false,
		"Automatically extract and determine resource translations for the provided resources.",
	)

	return cmd
}

func pruneResTr(nn []resource.Interface) (mm []resource.Interface) {
	mm = make([]resource.Interface, 0, len(nn))
	for _, n := range nn {
		if n.ResourceType() != resource.ResourceTranslationType {
			mm = append(mm, n)
			continue
		}

		r := n.(*resource.ResourceTranslation)
		if !r.IsDefault() {
			mm = append(mm, n)
			continue
		}
	}
	return mm
}
