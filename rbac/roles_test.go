package rbac_test

import (
	"testing"
)

func TestRoles(t *testing.T) {
	rbac, err := getClient()
	if err != nil {
		t.Errorf("Unexpected error when creating RBAC instance: %+v", err)
	}
	rbac.Debug("info")

	roles := rbac.Roles()
	roles.Delete("test-role")

	if err := roles.Create("test-role"); err != nil {
		t.Errorf("Error when creating test-role: %+v", err)
	}

	if err := roles.Create("test-role/nested/role"); err == nil {
		t.Errorf("Expected error when creating deep nested role, got nil")
		return
	}

	if err := roles.Create("test-role/nested"); err != nil {
		t.Errorf("Expected error when creating deep nested role, got nil")
		return
	}

	{
		role, err := roles.Get("test-role")
		assert(t, err == nil, "Unexpected error when getting role, %+v", err)
		assert(t, role.Name == "test-role", "Unexpected role name, test-role != '%s'", role.Name)
	}

	if err := roles.Delete("test-role"); err != nil {
		t.Errorf("Error when deleting test-role: %+v", err)
	}

	if err := roles.Delete("non-existant"); err == nil {
		t.Errorf("Expected error on deleting a non-existant role")
	}
}
