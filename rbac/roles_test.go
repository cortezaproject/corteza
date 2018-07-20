package rbac_test

import (
	"testing"
)

func TestRoles(t *testing.T) {
	rbac, err := getClient()
	if err != nil {
		t.Errorf("Unexpected error when creating RBAC instance: %+v", err)
	}
	rbac.Debug(false)

	roles := rbac.Roles()
	roles.Delete("test-role")

	if err := roles.Create("test-role"); err != nil {
		t.Errorf("Error when creating test-role: %+v", err)
	}

	if err := roles.Delete("test-role"); err != nil {
		t.Errorf("Error when deleting test-role: %+v", err)
	}

	if err := roles.Delete("non-existant"); err == nil {
		t.Errorf("Expected error on deleting a non-existant role")
	}
}
