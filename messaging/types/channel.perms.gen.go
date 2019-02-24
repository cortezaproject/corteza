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
			Title: "Message Permissions",
			Operations: []rules.Operation{
				rules.Operation{
					Key:      "message.send",
					Title:    "Send Messages",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "message.embed",
					Title:    "Embed Links",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "message.attach",
					Title:    "Attach Files",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "message.update_own",
					Title:    "Update own messages",
					Subtitle: "Members with this permission can update/delete their own messages inside this channel",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "message.update_all",
					Title:    "Update messages",
					Subtitle: "Members with this permission can update/delete messages inside this channel",
					Enabled:  true,
					Default:  rules.Inherit,
				}, rules.Operation{
					Key:      "message.react",
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
		"message.embed":      rules.Inherit,
		"message.attach":     rules.Inherit,
		"message.update_own": rules.Inherit,
		"message.update_all": rules.Inherit,
		"message.react":      rules.Inherit,
		"manage.webhooks":    rules.Inherit,
		"message.send":       rules.Inherit,
	}
	if value, ok := values[key]; ok {
		return value
	}
	return rules.Inherit
}
