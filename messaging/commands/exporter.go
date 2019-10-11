package commands

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/cortezaproject/corteza-server/messaging/service"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/settings"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

func Exporter(ctx context.Context, c *cli.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export",
		Long:  `Export messaging resources`,

		Run: func(cmd *cobra.Command, args []string) {

			c.InitServices(ctx, c)
			ctx = auth.SetSuperUserContext(ctx)

			var (
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
	cmd.Flags().BoolP("permissions", "p", false, "Export permission")

	return cmd
}

func permissionExporter(ctx context.Context, out *Messaging) {
	roles = sysTypes.RoleSet{
		&sysTypes.Role{ID: permissions.EveryoneRoleID, Handle: "everyone"},
		&sysTypes.Role{ID: permissions.AdminsRoleID, Handle: "admins"},
	}

	out.Allow = expServicePermissions(permissions.Allow)
	out.Deny = expServicePermissions(permissions.Deny)
}

func settingExporter(ctx context.Context, out *Messaging) {
	var (
		err error
	)

	ss, err := service.DefaultSettings.FindByPrefix("")
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

var (
	roles sysTypes.RoleSet
)

// @todo move to pkg/permissions
func expServicePermissions(access permissions.Access) map[string]map[string][]string {
	var (
		has   bool
		res   string
		rules permissions.RuleSet
		sp    = make(map[string]map[string][]string)
	)

	for _, r := range roles {
		rules = service.DefaultPermissions.FindRulesByRoleID(r.ID)

		if len(rules) == 0 {
			continue
		}

		for _, rule := range rules {
			if rule.Resource.GetService() != rule.Resource && !rule.Resource.HasWildcard() {
				continue
			}

			res = strings.TrimRight(rule.Resource.String(), ":*")

			if _, has = sp[r.Handle]; !has {
				sp[r.Handle] = map[string][]string{}
			}

			if _, has = sp[r.Handle][res]; !has {
				sp[r.Handle][res] = make([]string, 0)
			}

			sp[r.Handle][res] = append(sp[r.Handle][res], rule.Operation.String())
		}
	}

	return sp
}
