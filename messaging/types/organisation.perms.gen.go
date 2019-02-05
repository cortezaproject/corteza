package types

import "github.com/crusttech/crust/internal/rules"

/* File is generated from sam/types/permissions/1-organisation.json with permissions.go */

func (c *Organisation) Permissions() []rules.OperationGroup {
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
					Key:      "send",
					Title:    "Send Messages",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "embed",
					Title:    "Embed Links",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "attach",
					Title:    "Attach Files",
					Subtitle: "",
					Enabled:  true,
					Default:  rules.Allow,
				}, rules.Operation{
					Key:      "manage.messages",
					Title:    "Manage messages",
					Subtitle: "Members with this permission can edit/delete messages inside channels",
					Enabled:  true,
					Default:  rules.Deny,
				}, rules.Operation{
					Key:      "react",
					Title:    "Manage reactions",
					Subtitle: "Members with this permission can add new reactions to a message",
					Enabled:  true,
					Default:  rules.Allow,
				},
			},
		},
	}
}
