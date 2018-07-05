package sam

import (
	"github.com/pkg/errors"
)

func (*Organisation) Edit(r *organisationEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check if user can add/edit organisation

	o := Organisation{}.new().SetID(r.id).SetName(r.name)
	if o.GetID() > 0 {
		return o, db.Replace("organisation", o)
	}
	o.SetID(factory.Sonyflake.NextID())
	return o, db.Insert("organisation", o)
}

func (*Organisation) Remove(r *organisationRemoveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check
	// @todo: don't actually delete organisation?!

	return nil, func() error {
		_, err := db.Exec("delete from organisation where id=?", r.id)
		return err
	}()
}

func (*Organisation) Read(r *organisationReadRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permissions check

	o := Organisation{}.new()
	return o, db.Get(o, "select * from organisation where id=?", r.id)
}

func (*Organisation) Search(r *organisationSearchRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permissions check
	// @todo: actual search for org

	res := make([]Organisation, 0)
	err = db.Select(&res, "select * from organisation order by name asc")
	return res, err
}

func (*Organisation) Archive(r *organisationArchiveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check
	// @todo: don't actually delete organisation?!

	return nil, func() error {
		_, err := db.Exec("delete from organisation where id=?", r.id)
		return err
	}()
}
