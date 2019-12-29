package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/system/repository"
	repomock "github.com/cortezaproject/corteza-server/system/repository/mocks"
	"github.com/cortezaproject/corteza-server/system/types"
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

		logger: zap.NewNop(),

		settings: &types.Settings{},

		eventbus: eventbus.New(),

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
	var p = goth.User{UserID: "some-profile-id", Provider: "google", Email: u.Email}

	crdRpoMock := repomock.NewMockCredentialsRepository(mockCtrl)
	crdRpoMock.EXPECT().
		FindByCredentials("google", p.UserID).
		Times(1).
		Return(types.CredentialsSet{c}, nil)

	crdRpoMock.EXPECT().
		Update(gomock.Any()).
		Times(1).
		Return(c, nil)

	usrRpoMock := repomock.NewMockUserRepository(mockCtrl)
	usrRpoMock.EXPECT().FindByID(u.ID).Times(1).Return(u, nil)

	svc := makeMockAuthService(usrRpoMock, crdRpoMock)
	svc.settings.Auth.External.Enabled = true

	{
		auser, err := svc.External(p)
		require.NoError(t, err)
		require.True(t, auser.ID == u.ID, "Did not receive expected user")
	}
}

func TestAuth_External_NonExisting(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	var u = &types.User{ID: 300000, Email: "foo@example.tld"}
	var c = &types.Credentials{ID: 200000, OwnerID: u.ID}
	var p = goth.User{UserID: "some-profile-id", Provider: "google", Email: u.Email}

	crdRpoMock := repomock.NewMockCredentialsRepository(mockCtrl)
	crdRpoMock.EXPECT().
		FindByCredentials("google", p.UserID).
		Times(1).
		Return(types.CredentialsSet{}, nil)

	crdRpoMock.EXPECT().
		Create(&types.Credentials{Kind: "google", OwnerID: u.ID, Credentials: p.UserID}).
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

	usrRpoMock.EXPECT().
		Total().
		Times(1).
		Return(uint(0))

	svc := makeMockAuthService(usrRpoMock, crdRpoMock)
	svc.settings.Auth.External.Enabled = true

	{
		auser, err := svc.External(p)
		require.NoError(t, err)
		require.True(t, auser.ID == u.ID, "Did not receive expected user")
	}
}

func Test_auth_validateInternalLogin(t *testing.T) {
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "no email", args: args{"", ""}, wantErr: true},
		{name: "bad email", args: args{"test", ""}, wantErr: true},
		{name: "no pass", args: args{"test@domain.tld", ""}, wantErr: true},
		{name: "all good", args: args{"test@domain.tld", "password"}, wantErr: false},
	}

	svc := auth{
		logger:   zap.NewNop(),
		settings: &types.Settings{},
	}

	svc.settings.Auth.Internal.Enabled = true

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := svc.validateInternalLogin(tt.args.email, tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("auth.validateInternalLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_auth_checkPassword(t *testing.T) {
	plainPassword := " ... plain password ... "
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	type args struct {
		password string
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
				password: " foo ",
				cc:       types.CredentialsSet{&types.Credentials{ID: 1, Credentials: string(hashedPassword)}}}},
		{
			name:    "invalid credentials",
			wantErr: true,
			args: args{
				password: " foo ",
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

	svc := auth{
		logger:   zap.NewNop(),
		settings: &types.Settings{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := svc.checkPassword(tt.args.password, tt.args.cc); (err != nil) != tt.wantErr {
				t.Errorf("auth.checkPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_auth_validateToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name            string
		args            args
		wantID          uint64
		wantCredentials string
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name:            "empty",
			wantID:          0,
			wantCredentials: "",
			wantErr:         true,
			args:            args{token: ""}},
		{
			name:            "foo",
			wantID:          0,
			wantCredentials: "",
			wantErr:         true,
			args:            args{token: "foo1"}},
		{
			name:            "semivalid",
			wantID:          0,
			wantCredentials: "",
			wantErr:         true,
			args:            args{token: "foofoofoofoofoofoofoofoofoofoofo0"}},
		{
			name:            "valid",
			wantID:          1,
			wantCredentials: "foofoofoofoofoofoofoofoofoofoofo",
			wantErr:         false,
			args:            args{token: "foofoofoofoofoofoofoofoofoofoofo1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := auth{}
			svc.logger = zap.NewNop()
			gotID, gotCredentials, err := svc.validateToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth.validateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("auth.validateToken() gotID = %v, want %v", gotID, tt.wantID)
			}
			if gotCredentials != tt.wantCredentials {
				t.Errorf("auth.validateToken() gotCredentials = %v, want %v", gotCredentials, tt.wantCredentials)
			}
		})
	}
}
