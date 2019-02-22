package service

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	permissionList = map[string]map[string]bool{
		"system": map[string]bool{
			"access":              true,
			"organisation.create": true,
			"role.create":         true,
		},
		"system:organisation": map[string]bool{
			"access": true,
		},
		"system:role": map[string]bool{
			"read":           true,
			"update":         true,
			"delete":         true,
			"members.manage": true,
		},
		"messaging": map[string]bool{
			"access":                 true,
			"channel.public.create":  true,
			"channel.private.create": true,
			"channel.direct.create":  true,
		},
		"messaging:channel": map[string]bool{
			"update":             true,
			"read":               true,
			"join":               true,
			"leave":              true,
			"members.manage":     true,
			"webhooks.manage":    true,
			"attachments.manage": true,
			"message.send":       true,
			"message.reply":      true,
			"message.embed":      true,
			"message.attach":     true,
			"message.update.own": true,
			"message.update.all": true,
			"message.react":      true,
		},
		"compose": map[string]bool{
			"access":           true,
			"namespace.create": true,
		},
		"compose:namespace": map[string]bool{
			"read":           true,
			"update":         true,
			"delete":         true,
			"module.create":  true,
			"chart.create":   true,
			"trigger.create": true,
			"page.create":    true,
		},
		"compose:module": map[string]bool{
			"read":          true,
			"update":        true,
			"delete":        true,
			"record.create": true,
			"record.read":   true,
			"record.update": true,
			"record.delete": true,
		},
		"compose:chart": map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		},
		"compose:trigger": map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		},
		"compose:page": map[string]bool{
			"read":   true,
			"update": true,
			"delete": true,
		},
	}
)

func validatePermission(resource string, operation string) error {
	delimiter := ":"
	resourceParts := strings.Split(resource, delimiter)
	if len(resourceParts) < 1 {
		return errors.Errorf("Invalid resource format, expected >= 1, got %d", len(resourceParts))
	}

	resourceName := resourceParts[0]
	if len(resourceParts) > 1 {
		resourceName = resourceParts[0] + delimiter + resourceParts[1]
	}

	if service, ok := permissionList[resourceName]; ok {
		if op := service[operation]; op {
			if len(resourceParts) == 3 {
				if val := resourceParts[2]; val != "" {
					return nil
				}
				return errors.Errorf("Invalid resource format, missing resource ID")
			}
			return nil
		}
		return errors.Errorf("Unknown operation: '%s'", operation)
	}
	return errors.Errorf("Unknown resource name: '%s'", resourceName)
}
