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
	sqlOrganisationScope  = "deleted_at IS NULL AND archived_at IS NULL"
	sqlOrganisationSelect = "SELECT * FROM organisations WHERE " + sqlOrganisationScope
)

type Organisation struct{}

func (Organisation) New() *Organisation {
	return &Organisation{}
}

func (*Organisation) Create(ctx context.Context, r *server.OrganisationCreateRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check if user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	o := types.Organisation{}.New().SetName(r.Name).SetID(factory.Sonyflake.NextID())
	return o, db.Insert("organisation", o)
}

func (*Organisation) Edit(ctx context.Context, r *server.OrganisationEditRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check if user can add/edit organisation
	// @todo: make sure archived & deleted entries can not be edited

	o := types.Organisation{}.New().SetID(r.ID).SetName(r.Name)
	return o, db.Replace("organisation", o)
}

func (*Organisation) Remove(ctx context.Context, r *server.OrganisationRemoveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check
	// @todo: don't actually delete organisation?!

	stmt := "UPDATE organisationss SET deleted_at = NOW() WHERE deleted_at IS NULL AND id = ?"

	return nil, func() error {
		_, err := db.Exec(stmt, r.ID)
		return err
	}()
}

func (*Organisation) Read(ctx context.Context, r *server.OrganisationReadRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permissions check

	o := types.Organisation{}.New()
	return o, db.Get(o, sqlOrganisationSelect+" AND id = ?", r.ID)
}

func (*Organisation) List(ctx context.Context, r *server.OrganisationListRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permissions check
	// @todo: actual search for org

	res := make([]Organisation, 0)
	err = db.Select(&res, sqlOrganisationSelect+" WHERE label LIKE = ? ORDER BY label ASC", r.Query+"%")
	return res, err
}

func (*Organisation) Archive(ctx context.Context, r *server.OrganisationArchiveRequest) (interface{}, error) {
	db, err := factory.Database.Get()
	if err != nil {
		return nil, err
	}

	// @todo: permission check

	stmt := fmt.Sprintf(
		"UPDATE organisation SET archived_at = NOW() WHERE %s AND id = ?",
		sqlChannelScope)

	return nil, func() error {
		_, err := db.Exec(stmt, r.ID)
		return err
	}()
}
