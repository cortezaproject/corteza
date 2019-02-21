package service

import (
	"strings"

	"github.com/pkg/errors"
)

func validatePermission(resource string, operation string) error {
	var services = map[string]map[string]bool{
		"messaging:channel": map[string]bool{
			"manage.webhooks":    false,
			"message.send":       true,
			"message.embed":      true,
			"message.attach":     true,
			"message.update_own": true,
			"message.update_all": true,
			"message.react":      true,
		},
	}

	delimiter := ":"
	resourceParts := strings.Split(resource, delimiter)
	if len(resourceParts) < 2 {
		return errors.Errorf("Invalid resource format, expected >= 2, got %d", len(resourceParts))
	}

	resourceName := resourceParts[0] + delimiter + resourceParts[1]
	if service, ok := services[resourceName]; ok {
		if op := service[operation]; op {
			if len(resourceParts) == 3 {
				if val := resourceParts[2]; val != "" {
					return nil
				}
			}
			return errors.Errorf("Invalid resource format, missing resource ID")
		}
		return errors.Errorf("Unknown operation: '%s'", operation)
	}
	return errors.Errorf("Unknown resource name: '%s'", resourceName)
}
