package commands

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	sysExporter "github.com/cortezaproject/corteza-server/system/exporter"
	"github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

func Exporter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export",
		Long:  `Export system resources`,

		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx = auth.SetSuperUserContext(cli.Context())

				sFlag = cmd.Flags().Lookup("settings").Changed
				pFlag = cmd.Flags().Lookup("permissions").Changed

				out = &System{
					Settings: yaml.MapSlice{},
				}
			)

			if !sFlag && !pFlag {
				cli.HandleError(errors.New("Specify setting or permissions flag"))
			}

			if pFlag {
				permissionExporter(ctx, out)
			}

			if sFlag {
				settingExporter(ctx, out)
			}

			y := yaml.NewEncoder(cmd.OutOrStdout())
			cli.HandleError(y.Encode(out))
		},
	}

	cmd.Flags().BoolP("settings", "s", false, "Export settings")
	cmd.Flags().BoolP("permissions", "p", false, "Export system permissions")

	return cmd
}

func permissionExporter(ctx context.Context, out *System) {
	roles := sysTypes.RoleSet{
		&sysTypes.Role{ID: permissions.EveryoneRoleID, Handle: "everyone"},
		&sysTypes.Role{ID: permissions.AdminsRoleID, Handle: "admins"},
	}

	out.Allow = sysExporter.ExportableServicePermissions(roles, service.DefaultPermissions, permissions.Allow)
	out.Deny = sysExporter.ExportableServicePermissions(roles, service.DefaultPermissions, permissions.Deny)
}

func settingExporter(ctx context.Context, out *System) {
	var (
		err error
	)

	ss, err := service.DefaultSettings.FindByPrefix(ctx)
	cli.HandleError(err)

	out.Settings = settings.Export(ss)
}

// This is PoC for exporting system resources
//

type (
	System struct {
		Settings yaml.MapSlice `yaml:",omitempty"`

		Allow map[string]map[string][]string `yaml:",omitempty"`
		Deny  map[string]map[string][]string `yaml:",omitempty"`
	}
)
