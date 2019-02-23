package types

import "github.com/crusttech/crust/internal/rules"

/* File is generated from messaging/types/permissions/1-organisation.json with permissions.go */

func (*Organisation) Permissions() []rules.OperationGroup {
	return []rules.OperationGroup{
		rules.OperationGroup{
			Title: "General permissions",
			Operations: []rules.Operation{
				rules.Operation{
					Key:      "admin",
					Title:    "Administrator",
					Subtitle: "Members with this permission have every permission and also bypass channel specific permissions. Granting this permission is dangerous",
					Enabled:  true,
					Default:  rules.Deny,
				}, rules.Operation{
					Key:      "audit",
					Title:    "View Audit Log (@todo: add audit logs)",
					Subtitle: "Members with this permission have access to view the servers audit logs",
					Enabled:  false,
					Default:  rules.Deny,
				}, rules.Operation{
					Key:      "manage.organisation",
					Title:    "Manage Organisation",
					Subtitle: "Members with this permission can change the organisation name and other organisation details",
					Enabled:  true,
					Default:  rules.Deny,
				}, rules.Operation{
					Key:      "manage.roles",
					Title:    "Manage Roles",
					Subtitle: "Members with this permission can create/edit/delete roles inside this organisation",
					Enabled:  true,
					Default:  rules.Deny,
				}, rules.Operation{
					Key:      "manage.channels",
					Title:    "Manage channels",
					Subtitle: "Members with this permission can create/edit/delete channels inside this organisation",
					Enabled:  true,
					Default:  rules.Deny,
				}, rules.Operation{
					Key:      "manage.webhooks",
					Title:    "Manage webhooks (@todo: implement webhooks)",
					Subtitle: "Members with this permission can create, edit and delete webhooks",
					Enabled:  false,
					Default:  rules.Deny,
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
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "message.embed",
					Title:    "Embed Links",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "message.attach",
					Title:    "Attach Files",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "message.update_own",
					Title:    "Update own messages",
					Subtitle: "Members with this permission can update/delete their own messages inside channels",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "message.update_all",
					Title:    "Update messages",
					Subtitle: "Members with this permission can update/delete messages inside channels",
					Enabled:  true,
					Default:  rules.Deny,
				}, rules.Operation{
					Key:      "message.react",
					Title:    "Manage reactions",
					Subtitle: "Members with this permission can add new reactions to a message",
					Enabled:  true,
					Default:  rules.Allow,
				},
			},
		},
	}
}

func (*Organisation) PermissionDefault(key string) rules.Access {
	values := map[string]rules.Access{
		"message.send":        rules.Allow,
		"message.update_own":  rules.Allow,
		"message.update_all":  rules.Deny,
		"message.react":       rules.Allow,
		"admin":               rules.Deny,
		"manage.organisation": rules.Deny,
		"manage.channels":     rules.Deny,
		"message.embed":       rules.Allow,
		"message.attach":      rules.Allow,
		"audit":               rules.Deny,
		"manage.roles":        rules.Deny,
		"manage.webhooks":     rules.Deny,
	}
	if value, ok := values[key]; ok {
		return value
	}
	return rules.Inherit
}
