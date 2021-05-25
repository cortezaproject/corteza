package commands

import (
	"context"
	"io"
	"os"
	"path"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
)

func Export(ctx context.Context, storeInit func(ctx context.Context) (store.Storer, error)) *cobra.Command {
	var (
		output string
	)

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export",
		Long:  `Export data to YAML files.`,

		Run: func(cmd *cobra.Command, args []string) {
			var (
				f = su.NewDecodeFilter()
			)

			s, err := storeInit(ctx)
			cli.HandleError(err)

			f = f.FromResource(args...)

			sd := su.Decoder()
			nn, err := sd.Decode(ctx, s, f)
			cli.HandleError(err)

			ye := yaml.NewYamlEncoder(&yaml.EncoderConfig{
				MappedOutput: true,
				// CompactOutput: true,
			})
			bld := envoy.NewBuilder(ye)
			g, err := bld.Build(ctx, nn...)
			cli.HandleError(err)
			err = envoy.Encode(ctx, g, ye)
			cli.HandleError(err)
			ss := ye.Stream()
			_ = ss

			makeFN := func(base, res string) string {
				pp := strings.Split(strings.Trim(res, ":"), ":")
				name := strings.Join(pp, "_") + ".yaml"
				return path.Join(base, name)
			}

			for _, s := range ss {
				f, err := os.Create(makeFN(output, s.Resource))
				cli.HandleError(err)
				defer f.Close()

				io.Copy(f, s.Source)
			}
		},
	}

	cmd.Flags().StringVarP(&output, "out", "o", "./", "The directory to write output files to")

	return cmd
}
