package types

import "github.com/crusttech/crust/internal/rules"

/* File is generated from sam/types/permissions/3-channel.json with permissions.go */

func (c *Channel) Permissions() []rules.OperationGroup {
	return []rules.OperationGroup{
		rules.OperationGroup{
			Title: "General permissions",
			Operations: []rules.Operation{
				rules.Operation{
					Key:      "manage.webhooks",
					Title:    "Manage webhooks (@todo: implement webhooks)",
					Subtitle: "Members with this permission can create, edit and delete webhooks",
					Enabled:  false,
					Default:  rules.Inherit,
				},
			},
		}, rules.OperationGroup{
			Title: "Text Permissions",
			Operations: []rules.Operation{
				rules.Operation{
					Key:      "send",
					Title:    "Send Messages",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "embed",
					Title:    "Embed Links",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "attach",
					Title:    "Attach Files",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "manage.messages",
					Title:    "Manage messages",
					Subtitle: "Members with this permission can edit/delete messages inside this channel",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "react",
					Title:    "Manage reactions",
					Subtitle: "Members with this permission can add new reactions to a message",
					Enabled:  true,
					Default:  rules.Inherit,
				},
			},
		},
	}
}
