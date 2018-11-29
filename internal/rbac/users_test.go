package rbac_test

import (
	"testing"
)

func TestUsers(t *testing.T) {
	rbac, err := getClient()
	must(t, err, "Error when creating RBAC instance")

	users := rbac.Users()
	roles := rbac.Roles()

	users.Delete("test-user")
	roles.Delete("test-role")

	must(t, roles.Create("test-role"), "Error when creating test-role")
	must(t, users.Create("test-user", "test-password"), "Error when creating test-user")

	// check if we inherited some roles (should be empty)
	{
		user, err := users.Get("test-user")
		must(t, err, "Error when retrieving test-user 1")
		assert(t, user.Username == "test-user", "Unexpected username, test-user != '%s'", user.Username)
		assert(t, len(user.AssignedRoles) == 0, "Unexpected number of roles, expected empty, got %+v", user.AssignedRoles)
	}

	must(t, users.AddRole("test-user", "test-role"), "Error when assigning test-role to test-user")

	// check if we inherited some roles (should be empty)
	{
		user, err := users.Get("test-user")
		must(t, err, "Error when retrieving test-user 3")
		assert(t, user.Username == "test-user", "Unexpected username, test-user != '%s'", user.Username)
		assert(t, len(user.AssignedRoles) == 1, "Unexpected number of roles, expected 1, got %+v", user.AssignedRoles)
		assert(t, user.AssignedRoles[0] == "test-role", "Unexpected role name, test-role != '%s'", user.AssignedRoles[0])
	}

	must(t, users.RemoveRole("test-user", "test-role"), "Error when deassigning test-role to test-user")

	// check roles are empty after de-assign
	{
		user, err := users.Get("test-user")
		must(t, err, "Error when retrieving test-user 4")
		assert(t, user.Username == "test-user", "Unexpected username, test-user != '%s'", user.Username)
		assert(t, len(user.AssignedRoles) == 0, "Unexpected number of roles, expected empty, got %+v", user.AssignedRoles)
	}

	must(t, users.Delete("test-user"), "Error when deleting test-user")
	mustFail(t, func() error {
		_, err := users.Get("test-user")
		return err
	}())
	mustFail(t, users.Delete("test-user"))
}
