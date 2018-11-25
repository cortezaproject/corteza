package repository

import (
	"context"

	"github.com/titpetric/factory"

	"testing"

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
		u1, err := userRepo.Create(user)
		assert(t, err == nil, "User.Create error: %+v", err)
		assert(t, user.ID == u1.ID, "Changes were not stored")
	}

	teamRepo := Team(context.Background(), factory.Database.MustGet())
	team := &types.Team{
		Name: "Test team v1",
	}

	{
		t1, err := teamRepo.Create(team)
		assert(t, err == nil, "Team.Create error: %+v", err)
		assert(t, team.Name == t1.Name, "Changes were not stored")
	}

	{
		err := teamRepo.MemberAddByID(team.ID, user.ID)
		assert(t, err == nil, "Team.MemberAddByID error: %+v", err)
	}

	{
		users, err := userRepo.Find(&types.UserFilter{Query: ""})
		assert(t, err == nil, "User.Find error: %+v", err)
		assert(t, len(users) > 0, "No user results found")
		assert(t, len(users[0].Teams) > 0, "No team results found")
	}
}
