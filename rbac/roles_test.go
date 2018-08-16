package rbac_test

import (
	"testing"
)

func TestRoles(t *testing.T) {
	rbac, err := getClient()
	assert(t, err == nil, "Error when creating RBAC instance: %+v", err)
	rbac.Debug("info")

	roles := rbac.Roles()
	roles.Delete("test-role")

	err = roles.Create("test-role")
	assert(t, err == nil, "Error when creating test-role: %+v", err)

	err = roles.Create("test-role/nested/role")
	assert(t, err != nil, "Expected error when creating deep nested role, got nil")

	err = roles.Create("test-role/nested")
	assert(t, err == nil, "Error when creating deep nested role, got %+v", err)

	err = roles.CreateNested("test-role", "nested", "role")
	assert(t, err == nil, "Error when creating deep nested role, got %+v", err)

	err = roles.CreateNested()
	assert(t, err != nil, "Expected non-nil error")

	{
		role, err := roles.Get("test-role")
		assert(t, err == nil, "Error when getting role, %+v", err)
		assert(t, role.Name == "test-role", "Unexpected role name, test-role != '%s'", role.Name)
	}

	{
		role, err := roles.Get("test-role/nested/role")
		assert(t, err == nil, "Error when getting role, %+v", err)
		assert(t, role.Name == "test-role/nested/role", "Unexpected role name, test != '%s'", role.Name)
	}

	{
		role, err := roles.GetNested()
		assert(t, role == nil, "Expected role=nil, got %+v", role)
		assert(t, err != nil, "Expected non-nil error")
	}

	{
		role, err := roles.GetNested("test-role", "nested")
		assert(t, err == nil, "Error when getting role, %+v", err)
		assert(t, role.Name == "test-role/nested", "Unexpected role name, test != '%s'", role.Name)
	}

	err = roles.Delete("test-role")
	assert(t, nil == err, "Error when deleting test-role: %+v", err)

	err = roles.Delete("non-existant")
	assert(t, err != nil, "Expected error on deleting a non-existant role")
}
