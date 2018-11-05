package repository

import (
	"context"

	"github.com/titpetric/factory"

	"testing"

	"github.com/crusttech/crust/sam/types"
)

func TestTeam(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Team(context.Background(), factory.Database.MustGet())
	team := &types.Team{}

	var name1, name2 = "Test team v1", "Test team v2"

	var aa []*types.Team

	{
		team.Name = name1
		team, err = rpo.CreateTeam(team)
		assert(t, err == nil, "CreateTeam error: %v", err)
		assert(t, team.Name == name1, "Changes were not stored")

		{
			team.Name = name2
			team, err = rpo.UpdateTeam(team)
			assert(t, err == nil, "UpdateTeam error: %v", err)
			assert(t, team.Name == name2, "Changes were not stored")
		}

		{
			team, err = rpo.FindTeamByID(team.ID)
			assert(t, err == nil, "FindTeamByID error: %v", err)
			assert(t, team.Name == name2, "Changes were not stored")
		}

		{
			aa, err = rpo.FindTeams(&types.TeamFilter{Query: name2})
			assert(t, err == nil, "FindTeams error: %v", err)
			assert(t, len(aa) > 0, "No results found")
		}

		{
			err = rpo.ArchiveTeamByID(team.ID)
			assert(t, err == nil, "ArchiveTeamByID error: %v", err)
		}

		{
			err = rpo.UnarchiveTeamByID(team.ID)
			assert(t, err == nil, "UnarchiveTeamByID error: %v", err)
		}

		{
			err = rpo.DeleteTeamByID(team.ID)
			assert(t, err == nil, "DeleteTeamByID error: %v", err)
		}
	}

}
