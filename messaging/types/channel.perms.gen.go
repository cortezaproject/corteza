package types

import "github.com/crusttech/crust/internal/rules"

/* File is generated from messaging/types/permissions/3-channel.json with permissions.go */

func (*Channel) Permissions() []rules.OperationGroup {
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
					Key:      "text.send",
					Title:    "Send Messages",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "text.embed",
					Title:    "Embed Links",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "text.attach",
					Title:    "Attach Files",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "text.edit_own",
					Title:    "Manage own messages",
					Subtitle: "Members with this permission can edit/delete their own messages inside this channel",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "text.edit_all",
					Title:    "Manage messages",
					Subtitle: "Members with this permission can edit/delete messages inside this channel",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "text.react",
					Title:    "Manage reactions",
					Subtitle: "Members with this permission can add new reactions to a message",
					Enabled:  true,
					Default:  rules.Inherit,
				},
			},
		},
	}
}

func (*Channel) PermissionDefault(key string) rules.Access {
	values := map[string]rules.Access{
		"text.edit_all":   rules.Inherit,
		"text.react":      rules.Inherit,
		"manage.webhooks": rules.Inherit,
		"text.send":       rules.Inherit,
		"text.embed":      rules.Inherit,
		"text.attach":     rules.Inherit,
		"text.edit_own":   rules.Inherit,
	}
	if value, ok := values[key]; ok {
		return value
	}
	return rules.Inherit
}
