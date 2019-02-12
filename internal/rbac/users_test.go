package rbac_test

import (
	"testing"
)

func TestUsers(t *testing.T) {
	rbac, err := getClient()
	must(t, err, "Error when creating RBAC instance")

	users := rbac.Users()
	roles := rbac.Roles()

	// Cleanup data
	roles.Delete("test-role")
	// @todo until users.Get implements getting user by email, we need to delete users at end of the test successful and unsuccessful.

	must(t, roles.Create("test-role"), "Error when creating test-role")

	user, err := users.Create("test-user@crust.tech", "test-password")
	must(t, err, "Error when creating test-user")

	// check if we inherited some roles (should be empty)
	{
		u1, err := users.Get(user.ID)
		must(t, err, "Error when retrieving test-user 1")
		assert(t, len(u1.AssignedRoles) == 0, "Unexpected number of roles, expected empty, got %+v", u1.AssignedRoles)
	}

	must(t, users.AddRole(user.ID, "test-role"), "Error when assigning test-role to test-user")

	// check if we inherited some roles (should be empty)
	{
		u2, err := users.Get(user.ID)
		must(t, err, "Error when retrieving test-user 2")
		assert(t, len(u2.AssignedRoles) == 1, "Unexpected number of roles, expected 1, got %+v", u2.AssignedRoles)
		assert(t, u2.AssignedRoles[0] == "test-role", "Unexpected role name, test-role != '%s'", u2.AssignedRoles[0])
	}

	must(t, users.RemoveRole(user.ID, "test-role"), "Error when de-assigning test-role to test-user")

	// check roles are empty after de-assign
	{
		u3, err := users.Get(user.ID)
		must(t, err, "Error when retrieving test-user 3")
		assert(t, len(u3.AssignedRoles) == 0, "Unexpected number of roles, expected empty, got %+v", u3.AssignedRoles)
	}

	must(t, users.Delete(user.ID), "Error when deleting test-user")
	mustFail(t, func() error {
		_, err := users.Get(user.ID)
		return err
	}())
	mustFail(t, users.Delete(user.ID))
}
