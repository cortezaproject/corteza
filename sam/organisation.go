package sam

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/titpetric/factory"
)

var _ = errors.Wrap

const (
	sqlOrganisationScope  = "deleted_at IS NULL AND archived_at IS NULL"
	sqlOrganisationSelect = "SELECT * FROM organisations WHERE " + sqlOrganisationScope
)

func (*Organisation) Edit(r *organisationEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check if user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

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

	stmt := "UPDATE organisationss SET deleted_at = NOW() WHERE deleted_at IS NULL AND id = ?"

	return nil, func() error {
		_, err := db.Exec(stmt, r.id)
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
	return o, db.Get(o, sqlOrganisationSelect+" AND id = ?", r.id)
}

func (*Organisation) Search(r *organisationSearchRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permissions check
	// @todo: actual search for org

	res := make([]Organisation, 0)
	err = db.Select(&res, sqlOrganisationSelect+" WHERE label LIKE = ? ORDER BY label ASC", r.query+"%")
	return res, err
}

func (*Organisation) Archive(r *organisationArchiveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check

	stmt := fmt.Sprintf(
		"UPDATE organisation SET archived_at = NOW() WHERE %s AND id = ?",
		sqlChannelScope)

	return nil, func() error {
		_, err := db.Exec(stmt, r.id)
		return err
	}()
}
