package repository

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/types"

	. "github.com/crusttech/crust/internal/test"
)

func TestOrganisation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	db := factory.Database.MustGet()

	// Run tests in transaction to maintain DB state.
	Error(t, db.Transaction(func() error {
		rpo := Organisation(context.Background(), db)
		org := &types.Organisation{
			Name: "Test organisation v1",
		}

		{
			oa, err := rpo.CreateOrganisation(org)
			assert(t, err == nil, "CreateOrganisation error: %+v", err)
			assert(t, oa.Name == org.Name, "Changes were not stored")
		}

		{
			org.Name = "Test organisation v2"

			oa, err := rpo.UpdateOrganisation(org)
			assert(t, err == nil, "UpdateOrganisation error: %+v", err)
			assert(t, oa.Name == org.Name, "Changes were not stored")
		}

		{
			oa, err := rpo.FindOrganisationByID(org.ID)
			assert(t, err == nil, "FindOrganisationByID error: %+v", err)
			assert(t, oa.Name == org.Name, "Changes were not stored")
		}

		{
			oa, err := rpo.FindOrganisations(&types.OrganisationFilter{Query: org.Name})
			assert(t, err == nil, "FindOrganisations error: %+v", err)
			assert(t, len(oa) != 0, "No results found")
		}

		{
			err := rpo.ArchiveOrganisationByID(org.ID)
			assert(t, err == nil, "ArchiveOrganisationByID error: %+v", err)
		}

		{
			err := rpo.UnarchiveOrganisationByID(org.ID)
			assert(t, err == nil, "UnarchiveOrganisationByID error: %+v", err)
		}

		{
			err := rpo.DeleteOrganisationByID(org.ID)
			assert(t, err == nil, "DeleteOrganisationByID error: %+v", err)
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
