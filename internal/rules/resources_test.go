package rules_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/auth"
	"github.com/crusttech/crust/internal/rules"
	"github.com/crusttech/crust/system/types"
)

func TestRules(t *testing.T) {
	user := &types.User{ID: 1337}
	ctx := auth.SetIdentityToContext(context.Background(), user)

	db := factory.Database.MustGet()
	mustFail(t, db.Transaction(func() error {
		db.Insert("sys_user", user)
		var i uint64 = 0
		for i < 5 {
			db.Insert("sys_team", types.Team{ID: i, Name: fmt.Sprintf("Team %d", i)})
			i++
		}
		db.Insert("sys_team_member", types.TeamMember{TeamID: 1, UserID: user.ID})
		db.Insert("sys_team_member", types.TeamMember{TeamID: 2, UserID: user.ID})

		resources := rules.NewResources(ctx, db)

		// default (unset=deny)
		{
			mustFail(t, resources.CheckAccess("channel:1", "edit"))
			mustFail(t, resources.CheckAccessMulti("channel:*", "edit"))
		}

		// allow channel:2 group:2 (default deny, multi=allow)
		{
			resources.Grant("channel:2", 2, []string{"edit", "delete"}, rules.Allow)
			mustFail(t, resources.CheckAccess("channel:1", "edit"))
			must(t, resources.CheckAccess("channel:2", "edit"))
			must(t, resources.CheckAccessMulti("channel:*", "edit"))
		}

		// deny channel:1 group:1 (explicit deny, multi=deny)
		{
			resources.Grant("channel:1", 1, []string{"edit"}, rules.Deny)
			mustFail(t, resources.CheckAccess("channel:1", "edit"))
			must(t, resources.CheckAccess("channel:2", "edit"))
			mustFail(t, resources.CheckAccessMulti("channel:*", "edit"))
		}

		// reset (unset=deny)
		{
			resources.Grant("channel:2", 2, []string{"edit", "delete"}, rules.Inherit)
			resources.Grant("channel:1", 1, []string{"edit", "delete"}, rules.Inherit)
			mustFail(t, resources.CheckAccess("channel:1", "edit"))
			mustFail(t, resources.CheckAccessMulti("channel:*", "edit"))
		}
		return errors.New("Rollback")
	}))
}
