package rbac_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func TestSessions(t *testing.T) {
	rbac, err := getClient()
	must(t, err, "Error when creating RBAC instance")

	sessions := rbac.Sessions()
	users := rbac.Users()
	roles := rbac.Roles()
	resources := rbac.Resources()

	// @todo until users.Get implements getting user by email, we need to delete users at end of the test successful and unsuccessful.
	sessions.Delete("test-session")
	roles.Delete("test-role")
	resources.Delete("test-resource")
	resources.Delete("team-1", "team-2", "team-3")

	must(t, roles.Create("test-role"), "Error when creating test-role")

	user, err := users.Create("test-user@crust.tech", "test-password")
	must(t, err, "Error when creating test-user@crust.tech")
	assert(t, user != nil, "%+v", errors.New("Expected non-nil user"))
	assert(t, user.ID != "", "%+v", errors.New("Expected non-empty user.ID"))
	assert(t, user.Username == "test-user@crust.tech", "%+v", errors.Errorf("Expected test-user@crust.tech == %s", user.Username))

	must(t, users.AddRole(user.ID, "test-role"), "Error when assigning test-role to test-user@crust.tech")
	must(t, sessions.Create("test-session", user.ID, "test-role"), "Error when creating test-session")
	must(t, resources.Create("test-resource", []string{"view", "edit", "delete"}), "Error when creating test-resource")
	must(t, resources.Grant("test-resource", "test-role", []string{"view", "edit"}), "Error when granting permissions to role on resource")

	// check role is created
	{
		session, err := sessions.Get("test-session")
		must(t, err, "Error when getting test-session")
		assert(t, session.ID == "test-session", "Unexpected Session ID, test-session != '%s'", session.ID)
		assert(t, len(session.Roles) == 1, "Expected one session role, got %+v", session.Roles)
		assert(t, session.Roles[0] == "test-role", "Unexpected session role, test-role != '%s'", session.Roles[0])
	}

	// check user has permissions from role
	{
		must(t, resources.CheckAccess("test-resource", "view", "test-session"), "Owner has permission, but CheckAccess reports error")
		mustFail(t, resources.CheckAccess("test-resource", "delete", "test-session"))
	}

	// check multi access
	{
		for i := 1; i <= 5; i++ {
			resources.Delete(fmt.Sprintf("team:%d", i))
			must(t, resources.Create(fmt.Sprintf("team:%d", i), []string{"edit"}), fmt.Sprintf("Error when creating team:%d", i))
		}
		mustFail(t, resources.CheckAccessMulti("team:*", "edit", "test-session"))
		resources.Grant("team:4", "test-role", []string{"edit"})
		must(t, resources.CheckAccess("team:4", "edit", "test-session"))
		must(t, resources.CheckAccessMulti("team:*", "edit", "test-session"))
	}

	must(t, sessions.DeactivateRole("test-session", "test-role"), "Error when deactivating session role")

	// check role is deactivated
	{
		session, err := sessions.Get("test-session")
		must(t, err, "Error when getting test-session")
		assert(t, session.ID == "test-session", "Unexpected Session ID, test-session != '%s'", session.ID)
		// assert(t, session.Username == "test-user", "Unexpected user, test-user != '%s'", session.Username)
		assert(t, len(session.Roles) == 0, "Expected one session role, got %+v", session.Roles)
	}

	must(t, sessions.ActivateRole("test-session", "test-role"), "Error when deactivating session role")

	// check role is activated
	{
		session, err := sessions.Get("test-session")
		must(t, err, "Error when getting test-session")
		assert(t, session.ID == "test-session", "Unexpected Session ID, test-session != '%s'", session.ID)
		// assert(t, session.Username == "test-user", "Unexpected user, test-user != '%s'", session.Username)
		assert(t, len(session.Roles) == 1, "Expected one session role, got %+v", session.Roles)
		assert(t, session.Roles[0] == "test-role", "Unexpected session role, test-role != '%s'", session.Roles[0])
	}

	must(t, sessions.Delete("test-session"), "Error when deleting test-session")
	mustFail(t, func() error {
		_, err := sessions.Get("test-session")
		return err
	}())
	mustFail(t, sessions.Delete("test-session"))

	must(t, users.Delete(user.ID), "Error when deleting test-user")
	mustFail(t, func() error {
		_, err := users.Get(user.ID)
		return err
	}())
	mustFail(t, users.Delete(user.ID))
}
