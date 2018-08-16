package rbac_test

import (
	"testing"
)

func TestResources(t *testing.T) {
	rbac, err := getClient()
	assert(t, err == nil, "Error when creating RBAC instance: %+v", err)
	rbac.Debug("debug")

	roles := rbac.Roles()
	resources := rbac.Resources()

	roles.Delete("test-role")
	resources.Delete("test-resource")

	must(t, roles.Create("test-role"), "Error when creating test-role")
	must(t, resources.Create("test-resource", []string{"view", "edit", "delete"}), "Error when creating test-resource")
	must(t, resources.Grant("test-resource", "test-role", []string{"view", "edit"}), "Error when granting permissions to role on resource")

	// test get resources (not implemented) @todo
	if false {
		res, err := resources.Get("test-resource")
		must(t, err, "Error when retrieving test-resource")
		assert(t, res != nil, "Expected non-nil test-resource")
	}

	must(t, resources.Delete("test-resource"), "Error deleting a resource")
	mustFail(t, resources.Delete("test-resource"))
}
