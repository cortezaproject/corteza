// +build unit

package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/markbates/goth"

	"github.com/crusttech/crust/internal/test"
	"github.com/crusttech/crust/system/internal/repository"
	repomock "github.com/crusttech/crust/system/internal/repository/mocks"
	"github.com/crusttech/crust/system/types"
)

func TestSocialSigninWithExistingCredentials(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var u = &types.User{ID: 300000, Email: "foo@example.tld"}
	var c = &types.Credentials{ID: 200000, OwnerID: u.ID}
	var p = goth.User{UserID: "some-profile-id", Provider: "gplus", Email: u.Email}

	crdRpoMock := repomock.NewMockCredentialsRepository(mockCtrl)
	crdRpoMock.EXPECT().
		FindByCredentials(types.CredentialsKindGPlus, p.UserID).
		Times(1).
		Return(types.CredentialsSet{c}, nil)

	usrRpoMock := repomock.NewMockUserRepository(mockCtrl)
	usrRpoMock.EXPECT().FindByID(u.ID).Times(1).Return(u, nil)

	svc := &auth{
		db:          &mockDB{},
		users:       usrRpoMock,
		credentials: crdRpoMock,
	}

	{
		auser, err := svc.Social(p)
		test.Assert(t, err == nil, "Auth.Social error: %+v", err)
		test.Assert(t, auser.ID == u.ID, "Did not receive expected user")
	}
}

func TestSocialSigninWithNewUserCredentials(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var u = &types.User{ID: 300000, Email: "foo@example.tld"}
	var c = &types.Credentials{ID: 200000, OwnerID: u.ID}
	var p = goth.User{UserID: "some-profile-id", Provider: "gplus", Email: u.Email}

	crdRpoMock := repomock.NewMockCredentialsRepository(mockCtrl)
	crdRpoMock.EXPECT().
		FindByCredentials(types.CredentialsKindGPlus, p.UserID).
		Times(1).
		Return(types.CredentialsSet{}, nil)

	crdRpoMock.EXPECT().
		Create(&types.Credentials{Kind: types.CredentialsKindGPlus, OwnerID: u.ID, Credentials: p.UserID}).
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

	svc := &auth{
		db:          &mockDB{},
		users:       usrRpoMock,
		credentials: crdRpoMock,
	}

	{
		auser, err := svc.Social(p)
		test.Assert(t, err == nil, "Auth.Social error: %+v", err)
		test.Assert(t, auser.ID == u.ID, "Did not receive expected user")
	}
}
