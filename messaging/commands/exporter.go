package commands

import (
	"context"
	"errors"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	sysExporter "github.com/cortezaproject/corteza-server/system/exporter"
	sysService "github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func Exporter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export",
		Long:  `Export messaging resources`,

		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx   = auth.SetSuperUserContext(cli.Context())
				sFlag = cmd.Flags().Lookup("settings").Changed
				pFlag = cmd.Flags().Lookup("permissions").Changed

				out = &Messaging{
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

func permissionExporter(ctx context.Context, out *Messaging) {
	roles := sysTypes.RoleSet{
		&sysTypes.Role{ID: rbac.EveryoneRoleID, Handle: "everyone"},
		&sysTypes.Role{ID: rbac.AdminsRoleID, Handle: "admins"},
	}

	out.Allow = sysExporter.ExportableServicePermissions(roles, rbac.Global(), rbac.Allow)
	out.Deny = sysExporter.ExportableServicePermissions(roles, rbac.Global(), rbac.Deny)
}

func settingExporter(ctx context.Context, out *Messaging) {
	var (
		err error
	)

	ss, err := sysService.DefaultSettings.FindByPrefix(ctx)
	cli.HandleError(err)

	out.Settings = settings.Export(ss)
}

// This is PoC for exporting messaging resources
//

type (
	Messaging struct {
		Settings yaml.MapSlice `yaml:",omitempty"`

		Allow map[string]map[string][]string `yaml:",omitempty"`
		Deny  map[string]map[string][]string `yaml:",omitempty"`
	}
)
