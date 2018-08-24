package repository

import (
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestOrganisation(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := New()
	org := &types.Organisation{}

	var name1, name2 = "Test organisation v1", "Test organisation v2"

	var oo []*types.Organisation

	{
		org.Name = name1

		org, err = rpo.CreateOrganisation(org)
		assert(t, err == nil, "CreateOrganisation error: %v", err)
		assert(t, org.Name == name1, "Changes were not stored")

		{
			org.Name = name2

			org, err = rpo.UpdateOrganisation(org)
			assert(t, err == nil, "UpdateOrganisation error: %v", err)
			assert(t, org.Name == name2, "Changes were not stored")
		}

		{
			org, err = rpo.FindOrganisationByID(org.ID)
			assert(t, err == nil, "FindOrganisationByID error: %v", err)
			assert(t, org.Name == name2, "Changes were not stored")
		}

		{
			oo, err = rpo.FindOrganisations(&types.OrganisationFilter{Query: name2})
			assert(t, err == nil, "FindOrganisations error: %v", err)
			assert(t, len(oo) != 0, "No results found")
		}

		{
			err = rpo.ArchiveOrganisationByID(org.ID)
			assert(t, err == nil, "ArchiveOrganisationByID error: %v", err)
		}

		{
			err = rpo.UnarchiveOrganisationByID(org.ID)
			assert(t, err == nil, "UnarchiveOrganisationByID error: %v", err)
		}

		{
			err = rpo.DeleteOrganisationByID(org.ID)
			assert(t, err == nil, "DeleteOrganisationByID error: %v", err)
		}
	}
}
