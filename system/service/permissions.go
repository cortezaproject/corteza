package service

import (
	"strings"

	"github.com/pkg/errors"
)

func validateResource(resource string) error {
	delimiter := ":"

	resources := map[string]bool{
		"messaging:channel": true,
	}

	if result := resources[resource]; result {
		return nil
	}

	resourceParts := strings.Split(resource, delimiter)
	if len(resourceParts) != 3 {
		return errors.Errorf("Invalid resource format, expected 3, got %d", len(resourceParts))
	}

	for prefix := range resources {
		prefix = prefix + delimiter
		if len(resource) < len(prefix)+1 {
			continue
		}
		if resource[:len(prefix)] == prefix {
			return nil
		}
	}

	return errors.Errorf("Invalid resource name, '%s'", resource)
}

func validatePermission(service string, resource string, operation string) error {
	var services = map[string]map[string]bool{
		"messaging": map[string]bool{
			"manage.webhooks":    false,
			"message.send":       true,
			"message.embed":      true,
			"message.attach":     true,
			"message.update_own": true,
			"message.update_all": true,
			"message.react":      true,
		},
	}
	if service, ok := services[service]; ok {
		if op := service[operation]; op {
			return validateResource(resource)
		}
		return errors.Errorf("Unknown operation: '%s'", operation)
	}
	return errors.Errorf("Unknown service name: '%s'", service)
}
