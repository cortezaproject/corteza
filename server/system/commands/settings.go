package commands

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/system/types"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza/server/pkg/cli"
	"github.com/cortezaproject/corteza/server/system/service"
)

func Settings(ctx context.Context, app serviceInitializer) *cobra.Command {
	var (
		cmd = &cobra.Command{
			Use:   "settings",
			Short: "Settings management",
		}
	)

	list := &cobra.Command{
		Use:     "list",
		Short:   "List all",
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

			prefix := cmd.Flags().Lookup("prefix").Value.String()
			if kv, err := service.DefaultSettings.FindByPrefix(ctx, prefix); err != nil {
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

	list.Flags().String("prefix", "", "SettingsFilter settings by prefix")

	get := &cobra.Command{
		Use: "get [key to get, ...]",

		Short:   "Get value (raw JSON) for a specific key",
		Args:    cobra.ExactArgs(1),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

			if v, err := service.DefaultSettings.Get(ctx, args[0], 0); err != nil {
				cli.HandleError(err)
			} else if v != nil {
				cmd.Printf("%v\n", v.Value)
			}
		},
	}

	set := &cobra.Command{
		Use:     "set [key to set] [value]",
		Short:   "Set value (raw JSON or string) for a specific key",
		Args:    cobra.ExactArgs(2),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

			value := args[1]

			v := &types.SettingValue{
				Name: args[0],
			}

			if cmd.Flags().Lookup("as-string").Changed {
				cli.HandleError(v.SetValue(value))
			} else {
				err := v.SetRawValue(value)
				if _, is := err.(*json.SyntaxError); is {
					// Quote the raw value and re-parse
					err = v.SetRawValue(`"` + value + `"`)
				}

				cli.HandleError(err)
			}

			cli.HandleError(service.DefaultSettings.Set(ctx, v))
		},
	}

	set.Flags().BoolP("as-string", "s", false, "Treat input value as string (to avoid wrapping in quotes)")

	imp := &cobra.Command{
		Use:     "import [file]",
		Short:   "Import settings as JSON from stdin or file",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

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
				vv      types.SettingValueSet
			)

			cli.HandleError(decoder.Decode(&input))

			for k, v := range input {
				val := &types.SettingValue{Name: k}

				cli.HandleError(val.SetValue(v))
				vv = append(vv, val)
			}

			if len(vv) > 0 {
				cli.HandleError(service.DefaultSettings.BulkSet(ctx, vv))
			}
		},
	}

	exp := &cobra.Command{
		Use:     "export [file]",
		Short:   "Import settings as JSON to stdout or file",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

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

			if vv, err := service.DefaultSettings.FindByPrefix(ctx); err != nil {
				cli.HandleError(err)
			} else {
				cli.HandleError(encoder.Encode(vv.KV()))
			}
		},
	}

	del := &cobra.Command{
		Use:     "delete [keys, ...]",
		Short:   "Set value (raw JSON) for a specific key (or by prefix)",
		Args:    cobra.MinimumNArgs(0),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			ctx = auth.SetIdentityToContext(ctx, auth.ServiceUser())

			var (
				names = []string{}
			)

			if prefix := cmd.Flags().Lookup("prefix").Value.String(); len(prefix) > 0 {
				if vv, err := service.DefaultSettings.FindByPrefix(ctx); err != nil {
					cli.HandleError(err)
				} else {
					_ = vv.Walk(func(v *types.SettingValue) error {
						names = append(names, v.Name)
						return nil
					})
				}
			} else if len(args) > 0 {
				names = args
			}

			for a := 0; a < len(names); a++ {
				cli.HandleError(service.DefaultSettings.Delete(ctx, names[a], 0))
			}
		},
	}

	del.Flags().String("prefix", "", "SettingsFilter settings by prefix")

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
