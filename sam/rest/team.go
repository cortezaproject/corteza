package rest

import (
	"fmt"
	"context"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/sam/rest/server"
	"github.com/crusttech/crust/sam/types"
)

var _ = errors.Wrap

const (
	sqlTeamScope  = "deleted_at IS NULL AND archived_at IS NULL"
	sqlTeamSelect = "SELECT * FROM teams WHERE " + sqlTeamScope
)

type Team struct{}

func (Team) New() *Team {
	return &Team{}
}

func (*Team) Create(ctx context.Context, r *server.TeamCreateRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check if user can add/edit the team
	// @todo: make sure archived & deleted entries can not be edited

	t := types.Team{}.New()
	t.SetName(r.Name).SetMemberIDs(r.Members).SetID(factory.Sonyflake.NextID())
	return t, db.Insert("team", t)
}

func (*Team) Edit(ctx context.Context, r *server.TeamEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check if user can add/edit the team
	// @todo: make sure archived & deleted entries can not be edited

	t := types.Team{}.New()
	t.SetID(r.ID).SetName(r.Name).SetMemberIDs(r.Members)
	return t, db.Replace("team", t)
}

func (*Team) Remove(ctx context.Context, r *server.TeamRemoveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	stmt := "UPDATE teams SET deleted_at = NOW() WHERE deleted_at IS NULL AND id = ?"

	return nil, func() error {
		_, err := db.Exec(stmt, r.ID)
		return err
	}()
}

func (*Team) Read(ctx context.Context, r *server.TeamReadRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	t := types.Team{}.New()
	return t, db.Get(t, sqlTeamSelect+" AND id = ?", r.ID)
}

func (*Team) List(ctx context.Context, r *server.TeamListRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	res := make([]Team, 0)
	err = db.Select(&res, sqlTeamSelect+" ORDER BY name ASC")
	return res, err
}

func (*Team) Archive(ctx context.Context, r *server.TeamArchiveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf(
		"UPDATE teams SET archived_at = NOW() WHERE %s AND id = ?",
		sqlTeamScope)

	return nil, func() error {
		_, err := db.Exec(stmt, r.ID)
		return err
	}()
}

func (*Team) Move(ctx context.Context, r *server.TeamMoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.move")
}

func (*Team) Merge(ctx context.Context, r *server.TeamMergeRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.merge")
}
