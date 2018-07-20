package rbac_test

import (
	"testing"
)

func TestUsers(t *testing.T) {
	rbac, err := getClient()
	if err != nil {
		t.Errorf("Unexpected error when creating RBAC instance: %+v", err)
	}
	rbac.Debug(false)

	users := rbac.Users()
	users.Delete("test-user")

	if err := users.Create("test-user", "test-password"); err != nil {
		t.Errorf("Error when creating test-user: %+v", err)
	}

	if err := users.Delete("test-user"); err != nil {
		t.Errorf("Error when deleting test-user: %+v", err)
	}

	if _, err := users.Get("non-existant"); err == nil {
		t.Errorf("Expected error on retrieving a non-existant user")
	}

	if err := users.Delete("non-existant"); err == nil {
		t.Errorf("Expected error on deleting a non-existant user")
	}
}
