package rbac_test

import (
	"os"
	"testing"

	"github.com/crusttech/crust/rbac"
	"github.com/namsral/flag"
)

var _ = os.Setenv

func TestUsers(t *testing.T) {
	rbac.Flags()
	flag.Parse()

	rbac, err := rbac.New()
	if err != nil {
		t.Errorf("Unexpected error when creating RBAC instance: %+v", err)
	}
	rbac.Debug(false)

	users := rbac.Users()

	// clean up env
	{
		// just making sure we don't have one
		// and we're ignoring possible err's
		users.Delete("test-user")
	}

	// create a user
	{
		err := users.Create("test-user", "test-password")
		if err != nil {
			t.Errorf("Error when creating test-user: %+v", err)
		}
	}

	// delete a user
	{
		err := users.Delete("test-user")
		if err != nil {
			t.Errorf("Error when deleting test-user: %+v", err)
		}
	}

	// check getting a non-existant user fails
	{
		_, err := users.Get("non-existant")
		if err == nil {
			t.Errorf("Expected error on retrieving a non-existant user")
		}
	}
}
