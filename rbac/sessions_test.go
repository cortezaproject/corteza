package rbac_test

import (
	"testing"
)

func TestSessions(t *testing.T) {
	rbac, err := getClient()
	if err != nil {
		t.Errorf("Unexpected error when creating RBAC instance: %+v", err)
	}
	rbac.Debug("info")

	sessions := rbac.Sessions()
	users := rbac.Users()
	roles := rbac.Roles()

	// clean up data
	users.Delete("test-user")
	sessions.Delete("test-session")
	roles.Delete("test-role")

	if err := roles.Create("test-role"); err != nil {
		t.Errorf("Unexpected error when creating test-role, %+v", err)
		return
	}

	if err := users.Create("test-user", "test-password"); err != nil {
		t.Errorf("Unexpected error when creating test-user, %+v", err)
		return
	}

	if err := users.AddRole("test-user", "test-role"); err != nil {
		t.Errorf("Unexpected error when assigning test-role to test-user, %+v", err)
		return
	}

	if err := sessions.Create("test-session", "test-user", "test-role"); err != nil {
		t.Errorf("Unexpected error when creating test-session, %+v", err)
		return
	}

	// check role is created
	{
		session, err := sessions.Get("test-session")
		if !assert(t, err == nil, "Unexpected error when getting test-session, %+v", err) {
			return
		}
		// @todo: DAASI should return session ID from a get-query as well
		// assert(t, session.ID == "test-session", "Unexpected returned Session ID, test-session != '%s'", session.ID)
		assert(t, session.Username == "test-user", "Unexpected returned user, test-user != '%s'", session.Username)
		if !assert(t, len(session.Roles) == 1, "Expected one session role, got %+v", session.Roles) {
			return
		}
		assert(t, session.Roles[0] == "test-role", "Unexpected session role, test-role != '%s'", session.Roles[0])
	}

	if err := sessions.DeactivateRole("test-session", "test-role"); err != nil {
		t.Errorf("Unexpected error when deactivating session role, %+v", err)
		return
	}

	// check role is deactivated
	{
		session, err := sessions.Get("test-session")
		if !assert(t, err == nil, "Unexpected error when getting test-session, %+v", err) {
			return
		}
		// @todo: DAASI should return session ID from a get-query as well
		// assert(t, session.ID == "test-session", "Unexpected returned Session ID, test-session != '%s'", session.ID)
		assert(t, session.Username == "test-user", "Unexpected returned user, test-user != '%s'", session.Username)
		if !assert(t, len(session.Roles) == 0, "Expected one session role, got %+v", session.Roles) {
			return
		}
	}

	if err := sessions.ActivateRole("test-session", "test-role"); err != nil {
		t.Errorf("Unexpected error when deactivating session role, %+v", err)
		return
	}

	// check role is activated
	{
		session, err := sessions.Get("test-session")
		if !assert(t, err == nil, "Unexpected error when getting test-session, %+v", err) {
			return
		}
		// @todo: DAASI should return session ID from a get-query as well
		// assert(t, session.ID == "test-session", "Unexpected returned Session ID, test-session != '%s'", session.ID)
		assert(t, session.Username == "test-user", "Unexpected returned user, test-user != '%s'", session.Username)
		if !assert(t, len(session.Roles) == 1, "Expected one session role, got %+v", session.Roles) {
			return
		}
		assert(t, session.Roles[0] == "test-role", "Unexpected session role, test-role != '%s'", session.Roles[0])
	}

	if err := sessions.Delete("test-session"); err != nil {
		t.Errorf("Unexpected error when deleting test-session, %+v", err)
	}

	// Write tests (need users, roles)
}
