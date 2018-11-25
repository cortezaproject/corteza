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

	userRepo := User(context.Background(), factory.Database.MustGet())
	user := &types.User{
		Name:     "John Doe",
		Username: "johndoe",
	}
	user.GeneratePassword("johndoe")

	{
		u1, err := userRepo.Create(user)
		assert(t, err == nil, "User.Create error: %v", err)
		assert(t, user.ID == u1.ID, "Changes were not stored")
	}

	teamRepo := Team(context.Background(), factory.Database.MustGet())
	team := &types.Team{
		Name: "Test team v1",
	}

	{
		t1, err := teamRepo.Create(team)
		assert(t, err == nil, "Team.Create error: %v", err)
		assert(t, team.Name == t1.Name, "Changes were not stored")
	}

	{
		team.Name = "Test team v2"
		t1, err := teamRepo.Update(team)
		assert(t, err == nil, "Team.Update error: %v", err)
		assert(t, team.Name == t1.Name, "Changes were not stored")
	}

	{
		t1, err := teamRepo.FindByID(team.ID)
		assert(t, err == nil, "Team.FindByID error: %v", err)
		assert(t, team.Name == t1.Name, "Changes were not stored")
	}

	{
		aa, err := teamRepo.Find(&types.TeamFilter{Query: team.Name})
		assert(t, err == nil, "Team.Find error: %v", err)
		assert(t, len(aa) > 0, "No results found")
	}

	{
		err := teamRepo.ArchiveByID(team.ID)
		assert(t, err == nil, "Team.ArchiveByID error: %v", err)
	}

	{
		err := teamRepo.UnarchiveByID(team.ID)
		assert(t, err == nil, "Team.UnarchiveByID error: %v", err)
	}

	{
		err := teamRepo.MemberAddByID(team.ID, user.ID)
		assert(t, err == nil, "Team.MemberAddByID error: %v", err)
	}

	{
		err := teamRepo.MemberRemoveByID(team.ID, user.ID)
		assert(t, err == nil, "Team.MemberRemoveByID error: %v", err)
	}

	{
		err := teamRepo.DeleteByID(team.ID)
		assert(t, err == nil, "Team.DeleteByID error: %v", err)
	}

	{
		err := userRepo.DeleteByID(user.ID)
		assert(t, err == nil, "User.DeleteByID error: %v", err)
	}
}
