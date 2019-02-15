package service

import (
	"context"
	"testing"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"
	internalRules "github.com/crusttech/crust/internal/rules"
	. "github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/messaging/types"
	systemRepos "github.com/crusttech/crust/system/repository"
	systemTypes "github.com/crusttech/crust/system/types"
)

func TestPermissions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}
	ctx := context.TODO()

	// Create user for test.
	userRepo := systemRepos.User(ctx, factory.Database.MustGet())
	user := &systemTypes.User{
		Name:     "John Doe",
		Username: "johndoe",
		SatosaID: "1234",
	}
	err := user.GeneratePassword("johndoe")
	NoError(t, err, "expected no error generating password, got %v", err)

	_, err = userRepo.Create(user)
	NoError(t, err, "expected no error creating user, got %v", err)

	// Create team for test and add user
	teamRepo := systemRepos.Team(ctx, factory.Database.MustGet())
	team := &systemTypes.Team{
		Name: "Test team v1",
	}
	_, err = teamRepo.Create(team)
	NoError(t, err, "expected no error creating team, got %v", err)

	err = teamRepo.MemberAddByID(team.ID, user.ID)
	NoError(t, err, "expected no error adding user to team, got %v", err)

	// Set Identity.
	ctx = auth.SetIdentityToContext(ctx, user)

	// Create scopes.
	scopes := internalRules.NewScope()
	scopes.Add(&types.Organisation{})
	scopes.Add(&types.Team{})
	scopes.Add(&types.Channel{})

	permissionsSvc := Permissions(scopes).With(ctx)

	// Get all available scopes and items
	list, err := permissionsSvc.List()

	scopeItems := list.([]internalRules.ScopeItem)
	NoError(t, err, "expected no error, receiving scopes")
	Assert(t, len(scopeItems) == 3, "expected 3 scopes, got %v", len(scopeItems))

	// Setup nothing for organization
	organisationScope := scopeItems[0]
	Assert(t, organisationScope.Scope == "organisation", "expected scope 'organisation', got %s", organisationScope.Scope)

	// Setup everything allow for team
	teamScope := scopeItems[1]
	Assert(t, teamScope.Scope == "team", "expected scope 'team', got %s", teamScope.Scope)

	rules := make([]internalRules.Rules, 0)
	for _, group := range teamScope.Permissions {
		for _, op := range group.Operations {
			r := internalRules.Rules{
				TeamID:    team.ID,
				Resource:  "team:1",
				Operation: op.Key,
				Value:     internalRules.Allow,
			}
			rules = append(rules, r)
		}
	}
	_, err = permissionsSvc.Set(team.ID, rules)
	NoError(t, err, "expected no error, setting rules")

	// Deny all permissions for scope channel:1
	channelScope := scopeItems[2]
	Assert(t, channelScope.Scope == "channel", "expected scope 'channel', got %s", channelScope.Scope)

	rules = make([]internalRules.Rules, 0)
	for _, group := range channelScope.Permissions {
		for _, op := range group.Operations {
			r := internalRules.Rules{
				TeamID:    team.ID,
				Resource:  "channel:1",
				Operation: op.Key,
				Value:     internalRules.Deny,
			}
			rules = append(rules, r)
		}
	}
	_, err = permissionsSvc.Set(team.ID, rules)
	NoError(t, err, "expected no error, setting rules")

	// Check permission on channel and team level to test inheritance.
	rls := Rules().With(ctx)

	canManageChannels := rls.canManageChannels()
	Assert(t, canManageChannels == false, "expected canManageChannels == false, got %v", canManageChannels)

	canSendMessage := rls.canSendMessages(&types.Channel{ID: 1})
	Assert(t, canSendMessage == false, "expected canSendMessage == false, got %v", canSendMessage)

	canSendMessage = rls.canSendMessages(&types.Channel{ID: 2})
	Assert(t, canSendMessage == true, "expected canSendMessage == true, got %v", canSendMessage)
}
