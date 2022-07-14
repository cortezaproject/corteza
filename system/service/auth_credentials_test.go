package service

import (
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func Test_isValidPassword(t *testing.T) {
	var (
		pwdPlain      = " ... plain password ... "
		pwdHashedB, _ = bcrypt.GenerateFromPassword([]byte(pwdPlain), bcrypt.DefaultCost)
		pwdHashed     = string(pwdHashedB)
		pwdUnknown    = "$2a$10$8sOZxfZinxnu3bAtpkqEx.wBBwOfci6aG1szgUyxm5.BL2WiLu.ni"
	)

	cases := []struct {
		name     string
		password string
		cc       types.CredentialsSet
		rval     bool
	}{
		{
			name: "empty set",
			rval: false,
		},
		{
			name:     "bad pwd",
			rval:     false,
			password: " foo ",
			cc:       types.CredentialsSet{&types.Credentials{ID: 1, Credentials: pwdHashed}},
		},
		{
			name:     "invalid credentials",
			rval:     false,
			password: " foo ",
			cc:       types.CredentialsSet{&types.Credentials{ID: 0, Credentials: pwdHashed}},
		},
		{
			name:     "ok",
			rval:     true,
			password: pwdPlain,
			cc:       types.CredentialsSet{&types.Credentials{ID: 1, Credentials: pwdHashed}},
		},
		{
			name:     "multipass",
			rval:     true,
			password: pwdPlain,
			cc: types.CredentialsSet{
				&types.Credentials{ID: 0, Credentials: pwdHashed},
				&types.Credentials{ID: 1, Credentials: pwdUnknown},
				&types.Credentials{ID: 2, Credentials: pwdHashed},
				&types.Credentials{ID: 3, Credentials: ""},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				req = require.New(t)
				rsp = isValidPassword(c.cc, c.password)
			)

			if c.rval {
				req.True(rsp)
			} else {
				req.False(rsp)
			}
		})
	}
}

func Test_isPasswordReused(t *testing.T) {
	var (
		pwdPlain      = " ... plain password ... "
		pwdHashedB, _ = bcrypt.GenerateFromPassword([]byte(pwdPlain), bcrypt.DefaultCost)
		pwdHashed     = string(pwdHashedB)
		pwdUnknown    = "$2a$10$8sOZxfZinxnu3bAtpkqEx.wBBwOfci6aG1szgUyxm5.BL2WiLu.ni"
	)

	cases := []struct {
		name     string
		password string
		window   time.Duration
		cc       types.CredentialsSet
		rval     bool
	}{
		{
			name:     "no credentials, not reused",
			rval:     false,
			password: pwdPlain,
			cc:       types.CredentialsSet{},
		},
		{
			name:     "not reused",
			rval:     false,
			password: pwdPlain,
			cc: types.CredentialsSet{
				&types.Credentials{ID: 1, Credentials: pwdUnknown},
				&types.Credentials{ID: 2, Credentials: ""},
			},
		},
		{
			name:     "present, valid, first",
			rval:     true,
			password: pwdPlain,
			cc: types.CredentialsSet{
				&types.Credentials{ID: 1, Credentials: pwdHashed},
				&types.Credentials{ID: 2, Credentials: pwdUnknown},
				&types.Credentials{ID: 3, Credentials: ""},
			},
		},
		{
			name:     "present, but within time window",
			rval:     false,
			password: pwdPlain,
			window:   5 * time.Minute,
			cc: types.CredentialsSet{
				&types.Credentials{ID: 1, Credentials: pwdHashed, CreatedAt: *now()},
				&types.Credentials{ID: 2, Credentials: pwdUnknown},
				&types.Credentials{ID: 3, Credentials: ""},
			},
		},
		{
			name:     "present, invalid, last",
			rval:     true,
			password: pwdPlain,
			cc: types.CredentialsSet{
				&types.Credentials{ID: 2, Credentials: pwdUnknown},
				&types.Credentials{ID: 3, Credentials: ""},
				&types.Credentials{ID: 1, Credentials: pwdHashed, DeletedAt: now()},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				req = require.New(t)
				rsp = isPasswordReused(c.cc, c.password, c.window)
			)

			if c.rval {
				req.True(rsp)
			} else {
				req.False(rsp)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name            string
		args            args
		wantID          uint64
		wantCredentials string
	}{
		{
			name:            "empty",
			wantID:          0,
			wantCredentials: "",
			args:            args{token: ""}},
		{
			name:            "foo",
			wantID:          0,
			wantCredentials: "",
			args:            args{token: "foo1"}},
		{
			name:            "semivalid",
			wantID:          0,
			wantCredentials: "",
			args:            args{token: "foofoofoofoofoofoofoofoofoofoofo0"}},
		{
			name:            "valid",
			wantID:          1,
			wantCredentials: "foofoofoofoofoofoofoofoofoofoofo",
			args:            args{token: "foofoofoofoofoofoofoofoofoofoofo1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotCredentials := validateToken(tt.args.token)

			if gotID != tt.wantID {
				t.Errorf("auth.validateToken() gotID = %v, want %v", gotID, tt.wantID)
			}
			if gotCredentials != tt.wantCredentials {
				t.Errorf("auth.validateToken() gotCredentials = %v, want %v", gotCredentials, tt.wantCredentials)
			}
		})
	}
}
