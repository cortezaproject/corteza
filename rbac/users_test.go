package rbac_test

import (
	"testing"
)

func TestUsers(t *testing.T) {
	rbac, err := getClient()
	assert(t, err == nil, "Error when creating RBAC instance: %+v", err)
	rbac.Debug("info")

	users := rbac.Users()
	roles := rbac.Roles()

	users.Delete("test-user")
	roles.Delete("test-role")

	if err := roles.Create("test-role"); err != nil {
		t.Fatalf("Error when creating test-role, %+v", err)
		return
	}

	if err := users.Create("test-user", "test-password"); err != nil {
		t.Fatalf("Error when creating test-user: %+v", err)
		return
	}

	// check if we inherited some roles (should be empty)
	{
		user, err := users.Get("test-user")
		if !assert(t, err == nil, "Error when retrieving test-user 1, %+v", err) {
			return
		}
		assert(t, user.Username == "test-user", "Unexpected username, test-user != '%s'", user.Username)
		assert(t, len(user.AssignedRoles) == 0, "Unexpected number of roles, expected empty, got %+v", user.AssignedRoles)
	}

	if err := users.AddRole("test-user", "test-role"); err != nil {
		t.Fatalf("Error when assigning test-role to test-user 2, %+v", err)
		return
	}

	// check if we inherited some roles (should be empty)
	{
		user, err := users.Get("test-user")
		if !assert(t, err == nil, "Error when retrieving test-user 3, %+v", err) {
			return
		}
		assert(t, user.Username == "test-user", "Unexpected username, test-user != '%s'", user.Username)
		if !assert(t, len(user.AssignedRoles) == 1, "Unexpected number of roles, expected 1, got %+v", user.AssignedRoles) {
			return
		}
		assert(t, user.AssignedRoles[0] == "test-role", "Unexpected role name, test-role != '%s'", user.AssignedRoles[0])
	}

	if err := users.RemoveRole("test-user", "test-role"); err != nil {
		t.Fatalf("Error when deassigning test-role to test-user, %+v", err)
		return
	}

	// check roles are empty after de-assign
	{
		user, err := users.Get("test-user")
		if !assert(t, err == nil, "Error when retrieving test-user 4, %+v", err) {
			return
		}
		assert(t, user.Username == "test-user", "Unexpected username, test-user != '%s'", user.Username)
		assert(t, len(user.AssignedRoles) == 0, "Unexpected number of roles, expected empty, got %+v", user.AssignedRoles)
	}

	if err := users.Delete("test-user"); err != nil {
		t.Fatalf("Error when deleting test-user: %+v", err)
		return
	}

	if _, err := users.Get("test-user"); err == nil {
		t.Fatalf("Expected error on retrieving a non-existant user")
		return
	}

	if err := users.Delete("test-user"); err == nil {
		t.Fatalf("Expected error on deleting a non-existant user")
	}
}
