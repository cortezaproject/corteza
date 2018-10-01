package types

import "github.com/crusttech/crust/internal/rbac"

/* File is generated from sam/types/permissions/3-channel.json & main.go */

func (c *Channel) Permissions() []rbac.OperationGroup {
	return []rbac.OperationGroup{
		rbac.OperationGroup{
			Name: "",
			Operations: []rbac.Operation{
				rbac.Operation{
					Key:     "manage.webhooks",
					Name:    "Manage webhooks (@todo: implement webhooks)",
					Title:   "Members with this permission can create, edit and delete webhooks",
					Enabled: false, Default: "",
				},
			},
		}, rbac.OperationGroup{
			Name: "",
			Operations: []rbac.Operation{
				rbac.Operation{
					Key:     "send",
					Name:    "Send Messages",
					Title:   "",
					Enabled: true, Default: "",
				}, rbac.Operation{
					Key:     "embed",
					Name:    "Embed Links",
					Title:   "",
					Enabled: true, Default: "",
				}, rbac.Operation{
					Key:     "attach",
					Name:    "Attach Files",
					Title:   "",
					Enabled: true, Default: "",
				}, rbac.Operation{
					Key:     "manage.messages",
					Name:    "Manage messages",
					Title:   "Members with this permission can edit/delete messages inside this channel",
					Enabled: true, Default: "",
				}, rbac.Operation{
					Key:     "react",
					Name:    "Manage reactions",
					Title:   "Members with this permission can add new reactions to a message",
					Enabled: true, Default: "",
				},
			},
		},
	}
}
