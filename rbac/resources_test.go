package rbac_test

import (
	"testing"
)

func TestResources(t *testing.T) {
	rbac, err := getClient()
	if err != nil {
		t.Errorf("Unexpected error when creating RBAC instance: %+v", err)
	}
	rbac.Debug("debug")

	resources := rbac.Resources()
	resources.Delete("test-resource")

	if err := resources.Create("test-resource", []string{"view", "edit", "delete"}); err != nil {
		t.Errorf("Error when creating test-resource, %+v", err)
		return
	}

	// test get resources
	if false {
		_, err := resources.Get("test-resource")
		if err != nil {
			t.Errorf("Error when retrieving test-resource, %+v", err)
			return
		}
	}

	if err := resources.Delete("test-resource"); err != nil {
		t.Errorf("Unexpected error deleting a resource, %+v", err)
		return
	}

	if err := resources.Delete("test-resource"); err == nil {
		t.Errorf("Expected error when deleting unexistant resource, got none")
	}
}
