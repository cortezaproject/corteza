package commands

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
)

func Import(ctx context.Context, storeInit func(ctx context.Context) (store.Storer, error)) *cobra.Command {
	var (
		replaceOnExisting    bool
		mergeLeftOnExisting  bool
		mergeRightOnExisting bool
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import data from yaml sources.",

		Run: func(cmd *cobra.Command, args []string) {
			s, err := storeInit(ctx)
			cli.HandleError(err)

			yd := yaml.Decoder()
			nn := make([]resource.Interface, 0, 200)

			if len(args) > 0 {
				for _, fn := range args {
					mm, err := directory.Decode(ctx, fn, yd)
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

			se := es.NewStoreEncoder(s, opt)
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

	return cmd
}
