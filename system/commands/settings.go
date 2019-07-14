package commands

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/internal/settings"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/internal/service"
)

func Settings(ctx context.Context, c *cli.Config) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "settings",
			Short: "Settings management",
		}
	)

	list := &cobra.Command{
		Use:   "list",
		Short: "List all",
		Run: func(cmd *cobra.Command, args []string) {
			c.InitServices(ctx, c)

			prefix := cmd.Flags().Lookup("prefix").Value.String()
			if kv, err := service.DefaultIntSettings.FindByPrefix(prefix); err != nil {
				cli.HandleError(err)
			} else {
				var maxlen int
				for _, v := range kv {
					if l := len(v.Name); l > maxlen {
						maxlen = l
					}
				}

				for _, v := range kv {
					cmd.Printf("%s%s\t%v\n", v.Name, strings.Repeat(" ", maxlen-len(v.Name)), v.Value)
				}
			}
		},
	}

	list.Flags().String("prefix", "", "Filter settings by prefix")

	get := &cobra.Command{
		Use: "get [key to get, ...]",

		Short: "Get value (raw JSON) for a specific key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.InitServices(ctx, c)

			if v, err := service.DefaultIntSettings.Get(args[0], 0); err != nil {
				cli.HandleError(err)
			} else if v != nil {
				cmd.Printf("%v\n", v.Value)
			}
		},
	}

	set := &cobra.Command{
		Use:   "set [key to set] [value]",
		Short: "Set value (raw JSON) for a specific key",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			c.InitServices(ctx, c)

			value := args[1]
			v := &settings.Value{
				Name: args[0],
			}

			if err := v.SetRawValue(value); err != nil {
				cli.HandleError(err)
			}

			cli.HandleError(service.DefaultIntSettings.Set(v))
		},
	}

	imp := &cobra.Command{
		Use:   "import [file]",
		Short: "Import settings as JSON from stdin or file",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.InitServices(ctx, c)

			var (
				fh  *os.File
				err error
			)

			if len(args) > 0 {
				fh, err = os.Open(args[0])
				cli.HandleError(err)
			} else {
				fh = os.Stdin
			}

			var (
				decoder = json.NewDecoder(fh)
				input   = map[string]interface{}{}
				vv      settings.ValueSet
			)

			cli.HandleError(decoder.Decode(&input))

			for k, v := range input {
				val := &settings.Value{Name: k}

				cli.HandleError(val.SetValue(v))
				vv = append(vv, val)
			}

			if len(vv) > 0 {
				cli.HandleError(service.DefaultIntSettings.BulkSet(vv))
			}
		},
	}

	exp := &cobra.Command{
		Use:   "export [file]",
		Short: "Import settings as JSON to stdout or file",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.InitServices(ctx, c)

			var (
				fh  *os.File
				err error
			)

			if len(args) > 0 {
				fh, err = os.Create(args[0])
				cli.HandleError(err)
			} else {
				fh = os.Stdin
			}

			var (
				encoder = json.NewEncoder(fh)
			)

			encoder.SetIndent("", "  ")

			if vv, err := service.DefaultIntSettings.FindByPrefix(""); err != nil {
				cli.HandleError(err)
			} else {
				cli.HandleError(encoder.Encode(vv.KV()))
			}
		},
	}

	del := &cobra.Command{
		Use:   "delete [keys, ...]",
		Short: "Set value (raw JSON) for a specific key",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c.InitServices(ctx, c)

			for a := 0; a < len(args); a++ {
				cli.HandleError(service.DefaultIntSettings.Delete(args[a], 0))
			}
		},
	}

	cmd.AddCommand(
		list,
		get,
		set,
		del,
		imp,
		exp,
	)

	return cmd
}

func SettingsAutoConfigure(cmd *cobra.Command) {
	// @todo load
	// autoDiscoverAuthSettings()
	// @todo store
}
