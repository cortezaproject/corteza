package rbac_test

import (
	"testing"
)

func TestSessions(t *testing.T) {
	rbac, err := getClient()
	if err != nil {
		t.Errorf("Unexpected error when creating RBAC instance: %+v", err)
	}
	rbac.Debug(false)

	sessions := rbac.Sessions()
	users := rbac.Users()

	// clean up data
	users.Delete("test-user")
	sessions.Delete("test-session")

	if err := users.Create("test-user", "test-password"); err != nil {
		t.Errorf("Unexpected error when creating test-user, %+v", err)
		return
	}
	if err := sessions.Create("test-session", "test-user", []string{}); err != nil {
		t.Errorf("Unexpected error when creating test-session, %+v", err)
	}

	if _, err := sessions.Get("test-session"); err != nil {
		t.Errorf("Unexpected error when getting test-session, %+v", err)
	}

	if err := sessions.Delete("test-session"); err != nil {
		t.Errorf("Unexpected error when deleting test-session, %+v", err)
	}

	// Write tests (need users, roles)
}
