package repository

import (
	"context"

	"github.com/titpetric/factory"

	"testing"

	"github.com/crusttech/crust/system/types"
)

func TestTeam(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Team(context.Background(), factory.Database.MustGet())
	team := &types.Team{
		Name: "Test team v1",
	}

	{
		t1, err := rpo.CreateTeam(team)
		assert(t, err == nil, "CreateTeam error: %v", err)
		assert(t, team.Name == t1.Name, "Changes were not stored")
	}
	{
		team.Name = "Test team v2"
		t1, err := rpo.UpdateTeam(team)
		assert(t, err == nil, "UpdateTeam error: %v", err)
		assert(t, team.Name == t1.Name, "Changes were not stored")
	}

	{
		t1, err := rpo.FindTeamByID(team.ID)
		assert(t, err == nil, "FindTeamByID error: %v", err)
		assert(t, team.Name == t1.Name, "Changes were not stored")
	}

	{
		aa, err := rpo.FindTeams(&types.TeamFilter{Query: team.Name})
		assert(t, err == nil, "FindTeams error: %v", err)
		assert(t, len(aa) > 0, "No results found")
	}

	{
		err := rpo.ArchiveTeamByID(team.ID)
		assert(t, err == nil, "ArchiveTeamByID error: %v", err)
	}

	{
		err := rpo.UnarchiveTeamByID(team.ID)
		assert(t, err == nil, "UnarchiveTeamByID error: %v", err)
	}

	{
		err := rpo.DeleteTeamByID(team.ID)
		assert(t, err == nil, "DeleteTeamByID error: %v", err)
	}
}
