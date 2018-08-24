package repository

import (
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestUser(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := New()
	team := &types.User{}

	var name1, name2 = "Test user v1", "Test user v2"

	var aa []*types.User

	{
		team.Username = name1
		team, err = rpo.CreateUser(team)
		assert(t, err == nil, "CreateUser error: %v", err)
		assert(t, team.Username == name1, "Changes were not stored")

		{
			team.Username = name2
			team, err = rpo.UpdateUser(team)
			assert(t, err == nil, "UpdateUser error: %v", err)
			assert(t, team.Username == name2, "Changes were not stored")
		}

		{
			team, err = rpo.FindUserByID(team.ID)
			assert(t, err == nil, "FindUserByID error: %v", err)
			assert(t, team.Username == name2, "Changes were not stored")
		}

		{
			aa, err = rpo.FindUsers(&types.UserFilter{Query: name2})
			assert(t, err == nil, "FindUsers error: %v", err)
			assert(t, len(aa) > 0, "No results found")
		}

		{
			err = rpo.DeleteUserByID(team.ID)
			assert(t, err == nil, "DeleteUserByID error: %v", err)
		}
	}
}
