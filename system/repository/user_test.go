package repository

import (
	"context"
	"testing"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/types"
)

func TestUser(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	userRepo := User(context.Background(), factory.Database.MustGet())
	user := &types.User{
		Name:     "John Doe",
		Username: "johndoe",
		SatosaID: "1234",
	}
	user.GeneratePassword("johndoe")

	{
		uu, err := userRepo.Create(user)
		assert(t, err == nil, "User.Create error: %+v", err)
		assert(t, user.ID == uu.ID, "Changes were not stored")
	}

	teamRepo := Team(context.Background(), factory.Database.MustGet())
	team := &types.Team{
		Name: "Test team v1",
	}

	{
		t1, err := teamRepo.Create(team)
		assert(t, err == nil, "Team.Create error: %+v", err)
		assert(t, team.Name == t1.Name, "Changes were not stored")

		err = teamRepo.MemberAddByID(t1.ID, user.ID)
		assert(t, err == nil, "Team.MemberAddByID error: %+v", err)
	}

	{
		uu, err := userRepo.FindByID(user.ID)
		assert(t, err == nil, "User.FindByID error: %+v", err)
		assert(t, len(uu.Teams) == 1, "Expected 1 team, got %d", len(uu.Teams))
	}

	{
		users, err := userRepo.Find(&types.UserFilter{Query: ""})
		assert(t, err == nil, "User.Find error: %+v", err)
		assert(t, len(users) == 1, "User.Find: expected 1 user, got %d", len(users))
		assert(t, len(users[0].Teams) == 1, "User.Find: expected 1 team, got %d", len(users[0].Teams))
	}
}
