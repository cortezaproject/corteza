// +build unit

package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/markbates/goth"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/system/internal/repository"
	repomock "github.com/crusttech/crust/system/internal/repository/mocks"
	"github.com/crusttech/crust/system/types"
)

// @todo this mockDB will be probably be used by other tests, move it to some common place
type mockDB struct{}

func (mockDB) Transaction(callback func() error) error { return callback() }

// Mock auth service with nil for current time, dummy provider validator and mock db
func makeMockAuthService(u repository.UserRepository, c repository.CredentialsRepository) *auth {
	return &auth{
		db:          &mockDB{},
		users:       u,
		credentials: c,

		providerValidator: func(s string) error {
			// All providers are valid.
			return nil
		},

		now: func() *time.Time {
			return nil
		},
	}
}

func TestAuth_External_Existing(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create some virtual user and credentials
	var u = &types.User{ID: 300000, Email: "foo@example.tld"}
	var c = &types.Credentials{ID: 200000, OwnerID: u.ID}

	// Profile to be used. make sure email matches
	var p = goth.User{UserID: "some-profile-id", Provider: "gplus", Email: u.Email}

	crdRpoMock := repomock.NewMockCredentialsRepository(mockCtrl)
	crdRpoMock.EXPECT().
		FindByCredentials("gplus", p.UserID).
		Times(1).
		Return(types.CredentialsSet{c}, nil)

	usrRpoMock := repomock.NewMockUserRepository(mockCtrl)
	usrRpoMock.EXPECT().FindByID(u.ID).Times(1).Return(u, nil)

	svc := makeMockAuthService(usrRpoMock, crdRpoMock)

	{
		auser, err := svc.External(p)
		test.NoError(t, err, "unexpected error from auth.External", err)
		test.Assert(t, auser.ID == u.ID, "Did not receive expected user")
	}
}

func TestAuth_External_NonExisting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var u = &types.User{ID: 300000, Email: "foo@example.tld"}
	var c = &types.Credentials{ID: 200000, OwnerID: u.ID}
	var p = goth.User{UserID: "some-profile-id", Provider: "gplus", Email: u.Email}

	crdRpoMock := repomock.NewMockCredentialsRepository(mockCtrl)
	crdRpoMock.EXPECT().
		FindByCredentials("gplus", p.UserID).
		Times(1).
		Return(types.CredentialsSet{}, nil)

	crdRpoMock.EXPECT().
		Create(&types.Credentials{Kind: "gplus", OwnerID: u.ID, Credentials: p.UserID}).
		Times(1).
		Return(c, nil)

	usrRpoMock := repomock.NewMockUserRepository(mockCtrl)
	usrRpoMock.EXPECT().
		FindByEmail(u.Email).
		Times(1).
		Return(nil, repository.ErrUserNotFound)

	usrRpoMock.EXPECT().
		Create(&types.User{Email: "foo@example.tld"}).
		Times(1).
		Return(u, nil)

	svc := makeMockAuthService(usrRpoMock, crdRpoMock)

	{
		auser, err := svc.External(p)
		test.NoError(t, err, "unexpected error from auth.External", err)
		test.Assert(t, auser.ID == u.ID, "Did not receive expected user")
	}
}
