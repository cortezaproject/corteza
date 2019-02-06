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
			Title: "Text Permissions",
			Operations: []rules.Operation{
				rules.Operation{
					Key:      "text.send",
					Title:    "Send Messages",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "text.embed",
					Title:    "Embed Links",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "text.attach",
					Title:    "Attach Files",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "text.edit_own",
					Title:    "Manage own messages",
					Subtitle: "Members with this permission can edit/delete their own messages inside channels",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "text.edit_all",
					Title:    "Manage messages",
					Subtitle: "Members with this permission can edit/delete messages inside channels",
					Enabled:  true,
					Default:  rules.Deny,
				}, rules.Operation{
					Key:      "text.react",
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
		"admin":               rules.Deny,
		"manage.organisation": rules.Deny,
		"manage.roles":        rules.Deny,
		"text.react":          rules.Allow,
		"text.embed":          rules.Allow,
		"text.attach":         rules.Allow,
		"text.edit_own":       rules.Allow,
		"text.edit_all":       rules.Deny,
		"audit":               rules.Deny,
		"manage.channels":     rules.Deny,
		"manage.webhooks":     rules.Deny,
		"text.send":           rules.Allow,
	}
	if value, ok := values[key]; ok {
		return value
	}
	return rules.Inherit
}
