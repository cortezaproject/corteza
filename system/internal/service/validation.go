package service

import (
	"github.com/pkg/errors"

	internalRules "github.com/crusttech/crust/internal/rules"
)

var (
	permissionList = map[string]map[string]bool{
		"system": map[string]bool{
			"access":              true,
			"grant":               true,
			"settings.read":       true,
			"settings.manage":     true,
			"organisation.create": true,
			"role.create":         true,
			"application.create":  true,
		},
		"system:organisation:": map[string]bool{
			"access": true,
		},
		"system:role:": map[string]bool{
			"read":           true,
			"update":         true,
			"delete":         true,
			"members.manage": true,
		},
		"system:application:": map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		},
		"messaging": map[string]bool{
			"access":                 true,
			"grant":                  true,
			"channel.public.create":  true,
			"channel.private.create": true,
			"channel.group.create":   true,
		},
		"messaging:channel:": map[string]bool{
			"update":             true,
			"read":               true,
			"join":               true,
			"leave":              true,
			"delete":             true,
			"undelete":           true,
			"archive":            true,
			"unarchive":          true,
			"members.manage":     true,
			"webhooks.manage":    true,
			"attachments.manage": true,
			"message.send":       true,
			"message.reply":      true,
			"message.embed":      true,
			"message.attach":     true,
			"message.update.own": true,
			"message.update.all": true,
			"message.delete.own": true,
			"message.delete.all": true,
			"message.react":      true,
		},
		"compose": map[string]bool{
			"access":           true,
			"grant":            true,
			"namespace.create": true,
		},
		"compose:namespace:": map[string]bool{
			"read":           true,
			"update":         true,
			"delete":         true,
			"module.create":  true,
			"chart.create":   true,
			"trigger.create": true,
			"page.create":    true,
		},
		"compose:module:": map[string]bool{
			"read":          true,
			"update":        true,
			"delete":        true,
			"record.create": true,
			"record.read":   true,
			"record.update": true,
			"record.delete": true,
		},
		"compose:chart:": map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		},
		"compose:trigger:": map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		},
		"compose:page:": map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		},
	}
)

func validatePermission(resource internalRules.Resource, operation string) error {
	if !resource.IsValid() {
		return errors.Errorf("invalid resource format: %q", resource)
	}

	res := resource.TrimID().String()

	if service, ok := permissionList[res]; ok {
		if op := service[operation]; op {
			return nil
		}
		return errors.Errorf("Unknown operation: '%s'", operation)
	}
	return errors.Errorf("Unknown resource name: '%s'", resource)
}
