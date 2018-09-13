package rbac_test

import (
	"testing"
)

func TestRoles(t *testing.T) {
	rbac, err := getClient()
	must(t, err, "Error when creating RBAC instance")
	rbac.Debug("info")

	roles := rbac.Roles()
	roles.Delete("test-role")

	mustFail(t, roles.CreateNested())
	must(t, roles.Create("test-role"), "Error when creating test-role")
	mustFail(t, roles.Create("test-role/nested/role"))
	must(t, roles.Create("test-role/nested"), "Error when creating deep nested role")
	must(t, roles.CreateNested("test-role", "nested", "role"), "Error when creating deep nested role")

	{
		role, err := roles.Get("test-role")
		must(t, err, "Error when getting role")
		assert(t, role.Name == "test-role", "Unexpected role name, test-role != '%s'", role.Name)
	}

	{
		role, err := roles.Get("test-role/nested/role")
		must(t, err, "Error when getting role")
		assert(t, role.Name == "test-role/nested/role", "Unexpected role name, test != '%s'", role.Name)
	}

	{
		role, err := roles.GetNested()
		mustFail(t, err)
		assert(t, role == nil, "Expected role=nil, got %+v", role)
	}

	{
		role, err := roles.GetNested("test-role", "nested")
		must(t, err, "Error when getting role")
		assert(t, role.Name == "test-role/nested", "Unexpected role name, test != '%s'", role.Name)
	}

	must(t, roles.Delete("test-role"), "Error when deleting test-role")
	mustFail(t, roles.Delete("non-existant"))
	mustFail(t, roles.Delete("test-role"))
}
