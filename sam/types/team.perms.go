package types

import "github.com/crusttech/crust/internal/rbac"

/* File is generated from sam/types/permissions/2-team.json  with permissions.go */

func (c *Team) Permissions() []rbac.OperationGroup {
	return []rbac.OperationGroup{
		rbac.OperationGroup{
			Title: "General permissions",
			Operations: []rbac.Operation{
				rbac.Operation{
					Key:      "manage.webhooks",
					Title:    "Manage webhooks (@todo: implement webhooks)",
					Subtitle: "Members with this permission can create, edit and delete webhooks",
					Enabled:  false,
					Default:  "",
				},
			},
		}, rbac.OperationGroup{
			Title: "Text Permissions",
			Operations: []rbac.Operation{
				rbac.Operation{
					Key:      "send",
					Title:    "Send Messages",
					Subtitle: "",
					Enabled:  true,
					Default:  "",
				}, rbac.Operation{
					Key:      "embed",
					Title:    "Embed Links",
					Subtitle: "",
					Enabled:  true,
					Default:  "",
				}, rbac.Operation{
					Key:      "attach",
					Title:    "Attach Files",
					Subtitle: "",
					Enabled:  true,
					Default:  "",
				}, rbac.Operation{
					Key:      "manage.messages",
					Title:    "Manage messages",
					Subtitle: "Members with this permission can edit/delete messages inside channels",
					Enabled:  true,
					Default:  "",
				}, rbac.Operation{
					Key:      "react",
					Title:    "Manage reactions",
					Subtitle: "Members with this permission can add new reactions to a message",
					Enabled:  true,
					Default:  "",
				},
			},
		},
	}
}
