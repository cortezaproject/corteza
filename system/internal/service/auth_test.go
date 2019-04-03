// +-build unit

package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"

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

func Test_auth_validateCredentials(t *testing.T) {
	type args struct {
		email    string
		password []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "no email", args: args{"", []byte("")}, wantErr: true},
		{name: "bad email", args: args{"test", []byte("")}, wantErr: true},
		{name: "no pass", args: args{"test@domain.tld", []byte("")}, wantErr: true},
		{name: "all good", args: args{"test@domain.tld", []byte("password")}, wantErr: false},
	}
	svc := auth{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := svc.validateCredentials(tt.args.email, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("auth.validateCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_auth_checkPassword(t *testing.T) {
	plainPassword := []byte(" ... plain password ... ")
	hashedPassword, _ := bcrypt.GenerateFromPassword(plainPassword, bcrypt.DefaultCost)
	type args struct {
		password []byte
		cc       types.CredentialsSet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty set",
			wantErr: true,
			args:    args{}},
		{
			name:    "bad pwd",
			wantErr: true,
			args: args{
				password: []byte(" foo "),
				cc:       types.CredentialsSet{&types.Credentials{ID: 1, Credentials: string(hashedPassword)}}}},
		{
			name:    "invalid credentials",
			wantErr: true,
			args: args{
				password: []byte(" foo "),
				cc:       types.CredentialsSet{&types.Credentials{ID: 0, Credentials: string(hashedPassword)}}}},
		{
			name:    "ok",
			wantErr: false,
			args: args{
				password: plainPassword,
				cc:       types.CredentialsSet{&types.Credentials{ID: 1, Credentials: string(hashedPassword)}}}},
		{
			name:    "multipass",
			wantErr: false,
			args: args{
				password: plainPassword,
				cc: types.CredentialsSet{
					&types.Credentials{ID: 0, Credentials: string(hashedPassword)},
					&types.Credentials{ID: 1, Credentials: "$2a$10$8sOZxfZinxnu3bAtpkqEx.wBBwOfci6aG1szgUyxm5.BL2WiLu.ni"},
					&types.Credentials{ID: 2, Credentials: string(hashedPassword)},
					&types.Credentials{ID: 3, Credentials: ""},
				}}},
	}
	svc := auth{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := svc.checkPassword(tt.args.password, tt.args.cc); (err != nil) != tt.wantErr {
				t.Errorf("auth.checkPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
