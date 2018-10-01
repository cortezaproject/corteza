package types

import "github.com/crusttech/crust/internal/rbac"

/* File is generated from sam/types/permissions/1-organisation.json & main.go */

func (c *Organisation) Permissions() []rbac.OperationGroup {
	return []rbac.OperationGroup{
		rbac.OperationGroup{
			Name: "",
			Operations: []rbac.Operation{
				rbac.Operation{
					Key:     "admin",
					Name:    "Administrator",
					Title:   "Members with this permission have every permission and also bypass channel specific permissions. Granting this permission is dangerous",
					Enabled: true, Default: "deny",
				}, rbac.Operation{
					Key:     "audit",
					Name:    "View Audit Log (@todo: add audit logs)",
					Title:   "Members with this permission have access to view the servers audit logs",
					Enabled: false, Default: "deny",
				}, rbac.Operation{
					Key:     "manage.organisation",
					Name:    "Manage Organisation",
					Title:   "Members with this permission can change the organisation name and other organisation details",
					Enabled: true, Default: "deny",
				}, rbac.Operation{
					Key:     "manage.roles",
					Name:    "Manage Roles",
					Title:   "Members with this permission can create/edit/delete roles inside this organisation",
					Enabled: true, Default: "deny",
				}, rbac.Operation{
					Key:     "manage.channels",
					Name:    "Manage channels",
					Title:   "Members with this permission can create/edit/delete channels inside this organisation",
					Enabled: true, Default: "deny",
				}, rbac.Operation{
					Key:     "manage.webhooks",
					Name:    "Manage webhooks (@todo: implement webhooks)",
					Title:   "Members with this permission can create, edit and delete webhooks",
					Enabled: false, Default: "deny",
				},
			},
		}, rbac.OperationGroup{
			Name: "",
			Operations: []rbac.Operation{
				rbac.Operation{
					Key:     "send",
					Name:    "Send Messages",
					Title:   "",
					Enabled: true, Default: "allow",
				}, rbac.Operation{
					Key:     "embed",
					Name:    "Embed Links",
					Title:   "",
					Enabled: true, Default: "allow",
				}, rbac.Operation{
					Key:     "attach",
					Name:    "Attach Files",
					Title:   "",
					Enabled: true, Default: "allow",
				}, rbac.Operation{
					Key:     "manage.messages",
					Name:    "Manage messages",
					Title:   "Members with this permission can edit/delete messages inside channels",
					Enabled: true, Default: "deny",
				}, rbac.Operation{
					Key:     "react",
					Name:    "Manage reactions",
					Title:   "Members with this permission can add new reactions to a message",
					Enabled: true, Default: "allow",
				},
			},
		},
	}
}
