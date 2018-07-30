package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"testing"
)

func TestOrganisation(t *testing.T) {
	var err error

	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	rpo := Organisation()
	ctx := context.Background()
	org := &types.Organisation{}

	var name1, name2 = "Test organisation v1", "Test organisation v2"

	var oo []*types.Organisation

	org.Name = name1

	org, err = rpo.CreateOrganisation(org)
	must(t, err)
	if org.Name != name1 {
		t.Fatal("Changes were not stored")
	}

	org.Name = name2

	org, err = rpo.UpdateOrganisation(org)
	must(t, err)
	if org.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	org, err = rpo.FindOrganisationByID(org.ID)
	must(t, err)
	if org.Name != name2 {
		t.Fatal("Changes were not stored")
	}

	oo, err = rpo.FindOrganisations(&types.OrganisationFilter{Query: name2})
	must(t, err)
	if len(oo) == 0 {
		t.Fatal("No results found")
	}

	must(t, rpo.ArchiveOrganisationByID(org.ID))
	must(t, rpo.UnarchiveOrganisationByID(org.ID))
	must(t, rpo.DeleteOrganisationByID(org.ID))
}
