// +build integration

package repository

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/system/types"
)

func TestOrganisation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	db := factory.Database.MustGet()

	// Run tests in transaction to maintain DB state.
	test.Error(t, db.Transaction(func() error {
		rpo := Organisation(context.Background(), db)
		org := &types.Organisation{
			Name: "Test organisation v1",
		}

		{
			oa, err := rpo.CreateOrganisation(org)
			test.Assert(t, err == nil, "CreateOrganisation error: %+v", err)
			test.Assert(t, oa.Name == org.Name, "Changes were not stored")
		}

		{
			org.Name = "Test organisation v2"

			oa, err := rpo.UpdateOrganisation(org)
			test.Assert(t, err == nil, "UpdateOrganisation error: %+v", err)
			test.Assert(t, oa.Name == org.Name, "Changes were not stored")
		}

		{
			oa, err := rpo.FindOrganisationByID(org.ID)
			test.Assert(t, err == nil, "FindOrganisationByID error: %+v", err)
			test.Assert(t, oa.Name == org.Name, "Changes were not stored")
		}

		{
			oa, err := rpo.FindOrganisations(&types.OrganisationFilter{Query: org.Name})
			test.Assert(t, err == nil, "FindOrganisations error: %+v", err)
			test.Assert(t, len(oa) != 0, "No results found")
		}

		{
			err := rpo.ArchiveOrganisationByID(org.ID)
			test.Assert(t, err == nil, "ArchiveOrganisationByID error: %+v", err)
		}

		{
			err := rpo.UnarchiveOrganisationByID(org.ID)
			test.Assert(t, err == nil, "UnarchiveOrganisationByID error: %+v", err)
		}

		{
			err := rpo.DeleteOrganisationByID(org.ID)
			test.Assert(t, err == nil, "DeleteOrganisationByID error: %+v", err)
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
