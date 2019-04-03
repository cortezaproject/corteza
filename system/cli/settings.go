package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/crusttech/crust/internal/settings"
)

func settingsCmd(ctx context.Context, service settings.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Settings management",
	}

	list := &cobra.Command{
		Use:   "list",
		Short: "List all",
		Run: func(cmd *cobra.Command, args []string) {
			prefix := cmd.Flags().Lookup("prefix").Value.String()
			if kv, err := service.FindByPrefix(prefix); err != nil {
				exit(cmd, err)
			} else {
				for _, v := range kv {
					cmd.Printf("%s\t%v\n", v.Name, v.Value)
				}
			}
		},
	}

	list.Flags().String("prefix", "", "Filter settings by prefix")

	get := &cobra.Command{
		Use: "get [key to get]",

		Short: "Get value (raw JSON) for a specific key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if v, err := service.Get(args[0], 0); err != nil {
				exit(cmd, err)
			} else if v != nil {
				cmd.Printf("%v\n", v.Value)
			}
			exit(cmd, nil)
		},
	}

	set := &cobra.Command{
		Use:   "set [key to set] [value",
		Short: "Set value (raw JSON) for a specific key",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			value := args[1]
			v := &settings.Value{
				Name: args[0],
			}

			if err := v.SetValueAsString(value); err != nil {
				exit(cmd, err)
			}

			exit(cmd, service.Set(v))
		},
	}

	del := &cobra.Command{
		Use:   "delete [key to remove]",
		Short: "Set value (raw JSON) for a specific key",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			exit(cmd, service.Delete(args[0], 0))
		},
	}

	cmd.AddCommand(list, get, set, del)

	return cmd
}
