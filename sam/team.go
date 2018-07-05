package sam

import (
	"github.com/pkg/errors"
)

func (*Team) Edit(r *teamEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	t := Team{}.new()
	t.SetID(r.id).SetName(r.name).SetMemberIDs(r.members)
	if t.GetID() > 0 {
		return t, db.Replace("team", t)
	}
	c.SetID(factory.Sonyflake.NextID())
	return c, db.Insert("team", t)
}

func (*Team) Remove(r *teamRemoveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	return nil, func() error {
		_, err := db.Exec("delete from team where id=?", r.id)
		return err
	}()
}

func (*Team) Read(r *teamReadRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	t := Team{}.new()
	return t, db.Get(t, "select * from team where id=?", r.id)
}

func (*Team) Search(r *teamSearchRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	res := make([]Team, 0)
	err = db.Select(&res, "select * from team order by name asc")
	return res, err
}

func (*Team) Archive(r *teamArchiveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	return nil, func() error {
		_, err := db.Exec("delete from team where id=?", r.id)
		return err
	}()
}

func (*Team) Move(r *teamMoveRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.move")
}

func (*Team) Merge(r *teamMergeRequest) (interface{}, error) {
	return nil, errors.New("Not implemented: Team.merge")
}
