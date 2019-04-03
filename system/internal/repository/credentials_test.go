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

func TestCredentials(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	db := factory.Database.MustGet()

	// Create credentials repository.
	crepo := Credentials(context.Background(), db)

	// Run tests in transaction to maintain DB state.
	test.Error(t, db.Transaction(func() error {
		db.Exec("DELETE FROM sys_credentials WHERE 1=1")

		cc := types.CredentialsSet{
			&types.Credentials{OwnerID: 10000, Kind: types.CredentialsKindLinkedin, Credentials: "linkedin-profile-id"},
			&types.Credentials{OwnerID: 10000, Kind: types.CredentialsKindGPlus, Credentials: "gplus-profile-id"},
			&types.Credentials{OwnerID: 20000, Kind: types.CredentialsKindFacebook, Credentials: "facebook-profile-id"},
		}

		for _, c := range cc {
			cNew, err := crepo.Create(c)
			test.Assert(t, err == nil, "Credentials.Create error: %+v", err)
			test.Assert(t, c.ID > 0, "Expecting credentials to have a valid ID")
			test.Assert(t, c.Valid(), "Expecting credentials to be valid after creation")

			_, err = crepo.FindByID(cNew.ID)
			test.Assert(t, err == nil, "Credentials.FindByID error: %+v", err)

			{
				r, err := crepo.FindByKind(c.OwnerID, c.Kind)
				test.Assert(t, err == nil, "Credentials.FindByKind error: %+v", err)
				test.Assert(t, len(r) == 1, "Expecting exactly 1 result from FindByKind, got: %d", len(r))
			}

			{
				r, err := crepo.FindByCredentials(c.Kind, c.Credentials)
				test.Assert(t, err == nil, "Credentials.FindByKind error: %+v", err)
				test.Assert(t, len(r) == 1, "Expecting exactly 1 result from FindByCredentials, got: %d", len(r))
			}
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
