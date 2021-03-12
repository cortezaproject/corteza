package commands

import (
	"context"
	"io"

	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
)

func Export(storeInit func(ctx context.Context) (store.Storer, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export",
		Long:  `Export data to YAML files.`,

		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx = auth.SetSuperUserContext(cli.Context())

				f = su.NewDecodeFilter()
			)

			s, err := storeInit(ctx)
			cli.HandleError(err)

			f = f.FromResource(args...)

			// Nothing to do here...

			sd := su.Decoder()
			nn, err := sd.Decode(ctx, s, f)
			cli.HandleError(err)

			ye := yaml.NewYamlEncoder(&yaml.EncoderConfig{
				MappedOutput:  true,
				CompactOutput: true,
			})
			bld := envoy.NewBuilder(ye)
			g, err := bld.Build(ctx, nn...)
			cli.HandleError(err)
			err = envoy.Encode(ctx, g, ye)
			cli.HandleError(err)
			ss := ye.Stream()

			// Std out
			// @todo write to directory?
			w := cmd.OutOrStdout()
			for _, s := range ss {
				io.Copy(w, s.Source)
			}
		},
	}

	return cmd
}
