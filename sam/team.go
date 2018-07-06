package sam

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

var _ = errors.Wrap

const (
	sqlTeamScope  = "deleted_at IS NULL AND archived_at IS NULL"
	sqlTeamSelect = "SELECT * FROM teams WHERE " + sqlTeamScope
)

func (*Team) Edit(r *teamEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check if user can add/edit the team
	// @todo: make sure archived & deleted entries can not be edited

	t := Team{}.new()
	t.SetID(r.id).SetName(r.name).SetMemberIDs(r.members)
	if t.GetID() > 0 {
		return t, db.Replace("team", t)
	}
	t.SetID(factory.Sonyflake.NextID())
	return t, db.Insert("team", t)
}

func (*Team) Remove(r *teamRemoveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	stmt := "UPDATE teams SET deleted_at = NOW() WHERE deleted_at IS NULL AND id = ?"

	return nil, func() error {
		_, err := db.Exec(stmt, r.id)
		return err
	}()
}

func (*Team) Read(r *teamReadRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	t := Team{}.new()
	return t, db.Get(t, sqlTeamSelect+" AND id = ?", r.id)
}

func (*Team) Search(r *teamSearchRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	res := make([]Team, 0)
	err = db.Select(&res, sqlTeamSelect+" ORDER BY name ASC")
	return res, err
}

func (*Team) Archive(r *teamArchiveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	stmt := fmt.Sprintf(
		"UPDATE teams SET archived_at = NOW() WHERE %s AND id = ?",
		sqlTeamScope)

	return nil, func() error {
		_, err := db.Exec(stmt, r.id)
		return err
	}()
}

func (*Team) Move(r *teamMoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.move")
}

func (*Team) Merge(r *teamMergeRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.merge")
}
