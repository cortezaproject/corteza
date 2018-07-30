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
	team := &types.Team{}

	var name1, name2 = "Test user v1", "Test user v2"

	var aa []*types.Team

	team.Name = name1

	team, err = rpo.CreateTeam(team)
	must(t, err)
	if team.Name != name1 {
		t.Fatal("Changes were not stored")
	}

	team.Name = name2

	team, err = rpo.UpdateTeam(team)
	must(t, err)
	if team.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	team, err = rpo.FindTeamByID(team.ID)
	must(t, err)
	if team.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	aa, err = rpo.FindTeams(&types.TeamFilter{Query: name2})
	must(t, err)
	if len(aa) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.DeleteTeamByID(team.ID))
}
