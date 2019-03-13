package repository

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/titpetric/factory"

	"github.com/crusttech/crust/system/types"

	. "github.com/crusttech/crust/internal/test"
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
	Error(t, db.Transaction(func() error {
		db.Delete("sys_credentials", "1=1")

		cc := types.CredentialsSet{
			&types.Credentials{OwnerID: 10000, Kind: types.CredentialsKindLinkedin, Credentials: "linkedin-profile-id"},
			&types.Credentials{OwnerID: 10000, Kind: types.CredentialsKindGPlus, Credentials: "gplus-profile-id"},
			&types.Credentials{OwnerID: 20000, Kind: types.CredentialsKindFacebook, Credentials: "facebook-profile-id"},
		}

		for _, c := range cc {
			cNew, err := crepo.Create(c)
			assert(t, err == nil, "Credentials.Create error: %+v", err)
			assert(t, c.ID > 0, "Expecting credentials to have a valid ID")
			assert(t, c.Valid(), "Expecting credentials to be valid after creation")

			_, err = crepo.FindByID(cNew.ID)
			assert(t, err == nil, "Credentials.FindByID error: %+v", err)

			{
				r, err := crepo.FindByKind(c.OwnerID, c.Kind)
				assert(t, err == nil, "Credentials.FindByKind error: %+v", err)
				assert(t, len(r) == 1, "Expecting exactly 1 result from FindByKind, got: %v", len(r))
			}

			{
				r, err := crepo.FindByCredentials(c.Kind, c.Credentials)
				assert(t, err == nil, "Credentials.FindByKind error: %+v", err)
				assert(t, len(r) == 1, "Expecting exactly 1 result from FindByCredentials, got: %v", len(r))
			}
		}
		return errors.New("Rollback")
	}), "expected rollback error")
}
